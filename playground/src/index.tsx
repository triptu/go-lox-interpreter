import { Prec } from "@codemirror/state";
import { ViewPlugin, keymap } from "@codemirror/view";
import { computed, effect, signal } from "@preact/signals";
import { EditorView, basicSetup } from "codemirror";
import { render } from "preact";
import { useEffect, useRef } from "preact/hooks";

import {
	errorLineDecorationField,
	keywordDecorationPlugin,
	updateErrorLinesEffect,
} from "./cm-extensions";
import { Button, Header } from "./components";
import { readLoxFiles } from "./examples/read-lox-files" with { type: "macro" };
import { RunIcon, SpinnerIcon, StopIcon } from "./icons";
import {
	type OutputLogger,
	codeStorage,
	debounce,
	jsLox,
	runCode,
	stopRun,
} from "./utils";

const isRunning = signal(false);
const isAutoRunEnabled = signal(false);
const editorView = signal<EditorView | null>(null);
const outputLines = signal<{ text: string; isError: boolean }[]>([]);

const reErrorLine = /\[line (\d+)(:\d+)?\] (Error.+)/;

export const sampleLoxFiles = await readLoxFiles();
console.log("sample lox files", sampleLoxFiles);

const $errorLinesToHighlight = computed(() => {
	if (!outputLines.value) return [];
	const errLines = [];
	for (const line of outputLines.value) {
		if (!line.isError) continue;
		const match = reErrorLine.exec(line.text);
		if (match) {
			let lineNum = Number.parseInt(match[1]);
			const msgText = match[3];
			if (msgText.endsWith("Expected ';' after previous expression.")) {
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
		outputLines.value = [...outputLines.value, { text: arg, isError: false }];
	},
	error: (errMsg: string) => {
		outputLines.value = [...outputLines.value, { text: errMsg, isError: true }];
	},
};

async function runCodeWithStateStuff(code: string) {
	try {
		isRunning.value = true;
		await runCode(code, outputLogger);
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
];

function SelectSampleCode() {
	return (
		<div class="flex items-center">
			<span class="text-gray-900 mr-2 hidden md:inline-block">Example:</span>
			<select class="px-2 w-36 sm:w-48 md:w-60 h-9 bg-gray-200 rounded-md shadow-xs text-sm appearance-none hover:bg-gray-300">
				<option value="hello">Hello World</option>
				<option value="fibonacci">Fibonacci</option>
				<option value="inheritance">Inheritance</option>
			</select>
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
					onClick={() => stopRun()}
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
			runCodeWithStateStuff(this.code);
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

	return (
		<div class="h-4/5 lg:h-auto lg:flex-1 flex flex-col gap-2">
			<div class="h-full" ref={editorParent} />
		</div>
	);
}

function Output() {
	return (
		<div class="flex-1 flex flex-col gap-2">
			<div class="flex-1 bg-gray-100 ring-1 ring-gray-200 text-gray-900 text-md rounded p-4 whitespace-pre-wrap font-mono overflow-auto">
				{outputLines.value.map((line, index) => (
					// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
					<pre key={index} class={line.isError && "text-red-700"}>
						{line.text}
					</pre>
				))}
			</div>
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
				<CodeEditor />
				<Output />
			</div>
		</div>
	);
}

render(<App />, document.getElementById("app"));
