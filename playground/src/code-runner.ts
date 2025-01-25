import { codeStorage } from "./utils";

export interface OutputLogger {
	clear: () => void;
	log: (arg: string) => void;
	error: (arg: string) => void;
}

let worker: Worker | null = null;
function initWorker() {
	if (worker) return;
	worker = new window.Worker("./code-runner-worker.js");
}

function terminateWorker() {
	if (!worker) return;
	worker.terminate();
	worker = null;
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
	initWorker();
	outputLogger.clear();
	if (!skipSave) {
		codeStorage.set(code);
	}

	const promise = new Promise<void>((resolve, reject) => {
		worker.onmessage = (event) => {
			const { type, data } = event.data;
			switch (type) {
				case "log":
					outputLogger.log(data);
					break;
				case "input": {
					// use prompt to get input from user
					const value = prompt(data);
					worker.postMessage({ type: "inputResult", data: value });
					break;
				}
				case "error": {
					const msgText = data.replace("Expect", "Expected");
					outputLogger.error(msgText);
					break;
				}
				case "fatal":
					outputLogger.error(data);
					reject(new Error(data));
					break;
				case "done":
					resolve();
					break;
				default:
					console.error(`Unknown event type from wasm: ${type}`);
			}
		};
	});

	worker.postMessage({ type: "run", code });

	return promise;
}

export function stopCurrentRun() {
	worker?.postMessage({ type: "stop" });
}
