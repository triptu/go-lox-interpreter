import { snippetCompletion } from "@codemirror/autocomplete";
import { javascriptLanguage } from "@codemirror/lang-javascript";
import { LanguageSupport } from "@codemirror/language";

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

const codeLocalStorageKey = "savedCode";
export const defaultCode = `// You can edit this code!
// You can also try other examples from the select box above
print("Hello World!");
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
