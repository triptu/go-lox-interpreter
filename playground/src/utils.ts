import { snippetCompletion } from "@codemirror/autocomplete";
import { javascriptLanguage } from "@codemirror/lang-javascript";
import { LanguageSupport } from "@codemirror/language";

declare class Go {
	argv: string[];
	env: { [envKey: string]: string };
	exit: (code: number) => void;
	importObject: WebAssembly.Imports;
	exited: boolean;
	mem: DataView;
	run(instance: WebAssembly.Instance): Promise<void>;
}

declare global {
	interface Window {
		loxrun: (
			command: string,
			code: string,
			onEvent: (event: { type: string; data: string }) => void,
		) => void;
		loxstop: () => void;
	}
}

// biome-ignore lint/suspicious/noExplicitAny: we don't args type here
type Callable = (...args: any[]) => void;

export function debounce(fn: Callable, delayMs: number) {
	let timeoutId: number | undefined;
	// biome-ignore lint/suspicious/noExplicitAny: <explanation>
	return (...args: any[]) => {
		if (timeoutId) {
			clearTimeout(timeoutId);
		}
		timeoutId = setTimeout(() => {
			timeoutId = undefined;
			fn(...args);
		}, delayMs);
	};
}

export interface OutputLogger {
	clear: () => void;
	log: (arg: string) => void;
	error: (arg: string) => void;
}

const codeLocalStorageKey = "savedCode";
const defaultCode = `print("Hello World!");
var a = 2;

fun sum(a,b) {
  return a + b;
}

print(sum(a, 5));`;

export const codeStorage = {
	get: (): string => {
		const savedCode = localStorage.getItem(codeLocalStorageKey);
		return savedCode || defaultCode;
	},
	set: (code: string) => {
		localStorage.setItem(codeLocalStorageKey, code);
	},
};

let initDone = false;
let initPromise: Promise<void> | null = null;
async function initWasm() {
	if (initDone) {
		return;
	}
	if (initPromise) {
		return initPromise;
	}
	const go = new Go();
	initPromise = new Promise((resolve, reject) => {
		WebAssembly.instantiateStreaming(fetch("lox.wasm"), go.importObject)
			.then((obj) => {
				go.run(obj.instance); // run the main method in go
				console.log("wasm loaded");
				initDone = true;
				initPromise = null;
				resolve();
			})
			.catch((err) => {
				console.error("failed to load wasm");
				reject(err);
			});
	});
	return initPromise;
}

export async function runCode(code: string, outputLogger: OutputLogger) {
	if (!code) {
		console.warn("No code to run");
		return;
	}
	await initWasm();
	outputLogger.clear();
	codeStorage.set(code);

	return new Promise<void>((resolve, reject) => {
		window.loxrun?.("run", code, ({ type, data }) => {
			switch (type) {
				case "log":
					outputLogger.log(data);
					break;
				case "error": {
					let msgText = data.replace("Expect", "Expected");
					if (msgText.endsWith("Expected ';' after expression.")) {
						msgText = msgText.replace(
							"after expression.",
							"after previous expression.",
						);
					}
					outputLogger.error(msgText);
					break;
				}
				case "done":
					resolve();
					break;
				default:
					console.error(`Unknown event type from wasm: ${type}`);
			}
		});
	});
}

export function stopRun() {
	window.loxstop();
}

export const jsLox = new LanguageSupport(javascriptLanguage, [
	javascriptLanguage.data.of({
		autocomplete: [
			snippetCompletion("class ${name} {\n  init(${args}) {\n    ${}\n  }\n}", {
				label: "class",
			}),
			snippetCompletion("fun ${name}(${args}) {\n  ${}\n}", {
				label: "fun",
			}),
			snippetCompletion("var ${name} = ${value};", {
				label: "var",
			}),
			snippetCompletion("if (${condition}) {\n  ${}\n}", {
				label: "if",
			}),
			snippetCompletion("while (${condition}) {\n  ${}\n}", {
				label: "while",
			}),
			snippetCompletion("for (${expr}; ${condition}; ${step}) {\n  ${}\n}", {
				label: "for",
			}),
			snippetCompletion("return ${value};", {
				label: "return",
			}),
			snippetCompletion("print(${value});", {
				label: "print",
			}),
		],
	}),
]);
