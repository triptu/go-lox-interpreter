import { javascript } from "@codemirror/lang-javascript";
import { ViewPlugin } from "@codemirror/view";
import { signal } from "@preact/signals";
import { EditorView, basicSetup } from "codemirror";
import { render } from "preact";
import { useEffect, useRef } from "preact/hooks";
import { Button, Header } from "./components";
import { RunIcon, SpinnerIcon } from "./icons";
import { type OutputLogger, codeStorage, debounce, runCode } from "./utils";

const isRunning = signal(false);
const isAutoRunEnabled = signal(false);
const editorView = signal<EditorView | null>(null);
const outputLines = signal<{ text: string; isError: boolean }[]>([]);

const outputLogger: OutputLogger = {
	clear: () => {
		outputLines.value = [];
	},
	log: (arg) => {
		outputLines.value.push({ text: arg, isError: false });
	},
	error: (errMsg: string) => {
		outputLines.value.push({ text: errMsg, isError: true });
	},
};

async function runCodeWithStateStuff(code: string) {
	try {
		isRunning.value = true;
		await runCode(code, outputLogger);
	} finally {
		isRunning.value = false;
	}
}

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

			<div class="flex items-center">
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
			extensions: [basicSetup, javascript(), autoRunCodePlugin],
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
