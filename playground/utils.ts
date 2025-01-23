// biome-ignore lint/suspicious/noExplicitAny: we don't args type here
type Callable = (...args: any[]) => void;

export function throttle(fn: Callable, delayMs: number) {
	let lastCall = 0;
	// biome-ignore lint/suspicious/noExplicitAny: <explanation>
	return (...args: any[]) => {
		const now = performance.now();
		if (now - lastCall >= delayMs) {
			lastCall = now;
			fn(...args);
		}
	};
}

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

interface OutputLogger {
	clear: () => void;
	// biome-ignore lint/suspicious/noExplicitAny: <explanation>
	log: (...args: any[]) => void;
	// biome-ignore lint/suspicious/noExplicitAny: <explanation>
	error: (...args: any[]) => void;
}

export function getOutputLogger(outputElement: HTMLElement): OutputLogger {
	return {
		clear: () => {
			outputElement.innerHTML = "";
		},
		// biome-ignore lint/suspicious/noExplicitAny: <explanation>
		log: (args: any[]) => {
			outputElement.innerHTML += `${args
				.map((arg) =>
					typeof arg === "object" ? JSON.stringify(arg, null, 2) : arg,
				)
				.join(" ")}\n`;
		},
		error: (errMsg: string) => {
			outputElement.innerHTML += `<span class="text-red-700">${errMsg}</span>\n`;
		},
	};
}

const codeLocalStorageKey = "savedCode";
const defaultCode = "console.log('Hello World!')";
export const codeStorage = {
	get: (): string => {
		const savedCode = localStorage.getItem(codeLocalStorageKey);
		return savedCode || defaultCode;
	},
	set: (code: string) => {
		localStorage.setItem(codeLocalStorageKey, code);
	},
};

export function runCode(code: string, outputLogger: OutputLogger) {
	outputLogger.clear();

	// Use a Web Worker for safer execution
	const worker = new Worker(
		URL.createObjectURL(
			new Blob(
				[
					`
              // Redirect console.log to main thread
              console.log = (...args) => {
                postMessage({type: 'log', data: args});
              };
              
              onmessage = function(e) {
                try {
                  eval(e.data);
                  postMessage({type: 'done'});
                } catch (error) {
                  postMessage({type: 'error', error: error.toString()});
                }
              }
            `,
				],
				{ type: "application/javascript" },
			),
		),
	);

	worker.onmessage = (e) => {
		if (e.data.type === "log") {
			outputLogger.log(e.data.data);
		} else if (e.data.type === "error") {
			outputLogger.error(e.data.error);
		}
	};

	worker.postMessage(code);
}
