import { javascript } from "@codemirror/lang-javascript";
import { ViewPlugin } from "@codemirror/view";
import { EditorView, basicSetup } from "codemirror";
import {
	codeStorage,
	debounce,
	getOutputLogger,
	runCode,
	throttle,
} from "./utils";

const editorParent = document.getElementById("editor");
const runButton = document.getElementById("runBtn");
// biome-ignore lint/style/noNonNullAssertion: <explanation>
const outputElement = document.getElementById("output")!;
if (!editorParent || !outputElement || !runButton)
	throw new Error("Run, Input or Output element not found");

let autoRun = false;

const outputLogger = getOutputLogger(outputElement);

const autoRunCodePlugin = ViewPlugin.fromClass(
	class {
		code = "";

		constructor(view) {
			this.code = view.state.doc.toString();
			runCode(this.code, outputLogger); // initial run
		}

		update = debounce((viewUpdate) => {
			if (!autoRun) return;
			const newCode = viewUpdate.state.doc.toString();
			if (this.code !== newCode) {
				this.code = newCode;
				codeStorage.set(this.code);
				runCode(this.code, outputLogger);
			}
		}, 1000);
	},
);

const view = new EditorView({
	doc: codeStorage.get(),
	extensions: [basicSetup, javascript(), autoRunCodePlugin],
	parent: editorParent,
});

const onRunButtonClick = throttle(() => {
	runCode(view.state.doc.toString(), outputLogger);
}, 1000);

const onToggleAutoRunClick = (e: Event) => {
	const state = e.target as HTMLInputElement;
	autoRun = state.checked;
};

runButton.addEventListener("click", onRunButtonClick);
