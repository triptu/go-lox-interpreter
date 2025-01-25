import { Prec } from "@codemirror/state";
import { ViewPlugin, keymap } from "@codemirror/view";
import { computed, effect, signal } from "@preact/signals";
import { EditorView, basicSetup } from "codemirror";
import { Fragment, render } from "preact";
import { useEffect, useRef, useState } from "preact/hooks";

import { indentWithTab } from "@codemirror/commands";
import {
	errorLineDecorationField,
	keywordDecorationPlugin,
	updateErrorLinesEffect,
} from "./cm-extensions";
import { type OutputLogger, runCode, stopCurrentRun } from "./code-runner";
import { Button, Header } from "./components";
import { readLoxFiles } from "./examples/read-lox-files" with { type: "macro" };
import { RunIcon, SpinnerIcon, StopIcon } from "./icons";
import { codeStorage, debounce, defaultCode, jsLox } from "./utils";

const isRunning = signal(false);
const isAutoRunEnabled = signal(false);
const editorView = signal<EditorView | null>(null);
const outputLines = signal<{ text: string; isError: boolean }[]>([]);

const reErrorLine = /\[line (\d+)(:\d+)?\] (Error.+)/;

export const sampleLoxFiles: { [filename: string]: string } =
	await readLoxFiles();

const $errorLinesToHighlight = computed(() => {
	if (!outputLines.value) return [];
	const errLines = [];
	for (const line of outputLines.value) {
		if (!line.isError) continue;
		const match = reErrorLine.exec(line.text);
		if (match) {
			let lineNum = Number.parseInt(match[1]);
			const msgText = match[3];
			if (
				msgText.endsWith("Expected ';' after expression.") &&
				!msgText.includes("Error at end")
			) {
				lineNum -= 1; // fix the highlighted line
			}
			errLines.push({
				line: lineNum,
				col: match[2] ? Number.parseInt(match[2].slice(1)) : undefined,
			});
		}
	}
	return errLines;
});

effect(() => {
	if (editorView.value) {
		editorView.value.dispatch({
			effects: updateErrorLinesEffect.of($errorLinesToHighlight.value),
		});
	}
});

const outputLogger: OutputLogger = {
	clear: () => {
		outputLines.value = [];
	},
	log: (arg) => {
		if (arg === "\f") {
			outputLines.value = [];
		}
		outputLines.value = [...outputLines.value, { text: arg, isError: false }];
	},
	error: (errMsg: string) => {
		outputLines.value = [...outputLines.value, { text: errMsg, isError: true }];
	},
};

async function runCodeWithStateStuff(code: string, skipSave = false) {
	try {
		isRunning.value = true;
		await runCode(code, outputLogger, skipSave);
	} catch (err) {
		console.warn(err);
	} finally {
		isRunning.value = false;
	}
}

const keymapExtension = [
	Prec.highest(
		keymap.of([
			{
				key: "Mod-s",
				run: ({ state }) => {
					runCodeWithStateStuff(state.doc.toString());
					return true;
				},
			},
		]),
	),
	keymap.of([indentWithTab]),
];

function SelectSampleCode() {
	const [selectedFilename, setSelectedFilename] = useState<string>("choose");
	const onChange = (e: Event) => {
		const target = e.target as HTMLSelectElement;
		const filename = target.value;
		setSelectedFilename(filename);
		let code = "";
		if (filename === "choose") {
			code = codeStorage.get();
		} else if (filename === "default") {
			code = defaultCode;
		} else {
			code = sampleLoxFiles[filename];
		}
		if (code) {
			// update editorview
			editorView.value?.dispatch({
				changes: {
					from: 0,
					to: editorView.value.state.doc.length,
					insert: code,
				},
			});
			runCodeWithStateStuff(code, true);
		}
	};
	return (
		<div class="flex items-center">
			<span class="text-gray-900 mr-2 hidden md:inline-block">Example:</span>
			<div class="px-2 bg-gray-200 hover:bg-gray-300 rounded-md shadow-xs">
				<select
					value={selectedFilename}
					onChange={onChange}
					class="w-36 truncate sm:w-48 md:w-60 h-9 bg-gray-200 text-sm hover:bg-gray-300 outline-none focus:border-none"
				>
					<option value="choose">Choose an example</option>
					<option value="default">Default Code</option>
					{Object.keys(sampleLoxFiles).map((filename) => (
						<option key={filename} value={filename}>
							{filename.replace(".lox", "")}
						</option>
					))}
				</select>
			</div>
		</div>
	);
}

function Toolbar() {
	return (
		<div class="w-full flex items-center justify-end md:justify-between gap-4">
			<SelectSampleCode />

			<div class="flex items-center gap-2">
				<Button
					className="h-9"
					type="button"
					onClick={() =>
						runCodeWithStateStuff(editorView.value?.state.doc.toString())
					}
					disabled={isRunning.value}
				>
					{isRunning.value ? <SpinnerIcon class="animate-spin" /> : <RunIcon />}
					Run
				</Button>
				<Button
					className="h-9"
					type="button"
					color="red"
					onClick={() => stopCurrentRun()}
					disabled={!isRunning.value}
				>
					<StopIcon />
					Stop
				</Button>
			</div>
		</div>
	);
}

const autoRunCodePlugin = ViewPlugin.fromClass(
	class {
		code = "";
		constructor(view) {
			this.code = view.state.doc.toString();
		}
		update = debounce((viewUpdate) => {
			if (!isAutoRunEnabled.value) return;
			const newCode = viewUpdate.state.doc.toString();
			if (this.code !== newCode) {
				this.code = newCode;
				runCodeWithStateStuff(this.code);
			}
		}, 1000);
	},
);

function CodeEditor() {
	const editorParent = useRef(null);

	useEffect(() => {
		editorView.value = new EditorView({
			doc: codeStorage.get(),
			extensions: [
				basicSetup,
				jsLox,
				autoRunCodePlugin,
				keymapExtension,
				errorLineDecorationField,
				keywordDecorationPlugin,
			],
			parent: editorParent.current,
		});
	}, []);

	return <div class="h-full" ref={editorParent} />;
}

function OutputLine({ text }: { text: string }) {
	return (
		<>
			{text.split("\\n").map((line, index) => (
				// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
				<Fragment key={index}>
					{line}
					<br />
				</Fragment>
			))}
		</>
	);
}

function Output() {
	return (
		<div class="h-full bg-gray-100 ring-1 ring-gray-200 text-gray-900 text-md rounded p-4 whitespace-pre-wrap font-mono overflow-auto">
			{outputLines.value.map((line, index) => (
				<span
					// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
					key={index}
					class={
						line.isError
							? "text-red-700 whitespace-pre-wrap"
							: "whitespace-pre-wrap"
					}
				>
					<OutputLine text={line.text} />
				</span>
			))}
		</div>
	);
}

function App() {
	return (
		<div class="bg-white rounded shadow flex flex-col flex-1">
			<Header />

			<div class="flex px-4 my-2">
				<Toolbar />
			</div>

			<div class="flex gap-4 px-4 mb-4 flex-col lg:flex-row flex-1">
				<div class="h-4/5 lg:h-auto lg:w-3/5 flex flex-col gap-2">
					<CodeEditor />
				</div>
				<div class="lg:w-2/5  h-1/5 lg:h-auto flex flex-col gap-2">
					<Output />
				</div>
			</div>
		</div>
	);
}

render(<App />, document.getElementById("app"));
