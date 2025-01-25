import { codeStorage } from "./utils";

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

export interface OutputLogger {
	clear: () => void;
	log: (arg: string) => void;
	error: (arg: string) => void;
}

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

export async function runCode(
	code: string,
	outputLogger: OutputLogger,
	skipSave = false,
) {
	if (!code) {
		console.warn("No code to run");
		return;
	}
	await initWasm();
	outputLogger.clear();
	if (!skipSave) {
		codeStorage.set(code);
	}

	return new Promise<void>((resolve, reject) => {
		window.loxrun?.("run", code, ({ type, data }) => {
			switch (type) {
				case "log":
					outputLogger.log(data);
					break;
				case "error": {
					const msgText = data.replace("Expect", "Expected");
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
