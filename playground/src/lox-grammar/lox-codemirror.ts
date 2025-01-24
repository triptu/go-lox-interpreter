import { completeFromList, snippetCompletion } from "@codemirror/autocomplete";
import {
	LRLanguage,
	LanguageSupport,
	delimitedIndent,
	foldInside,
	foldNodeProp,
	indentNodeProp,
} from "@codemirror/language";
import { styleTags, tags as t } from "@lezer/highlight";
import { parser } from "./lox-parser";

export const loxLanguage = LRLanguage.define({
	parser: parser.configure({
		props: [
			styleTags({
				// Keywords
				"if else while for return and or": t.controlKeyword,
				"var fun class": t.definitionKeyword,
				"this super": t.special(t.keyword),

				// Constants
				"true false nil": t.bool,

				// Identifiers and Types
				"Identifier ClassDeclaration/identifier": t.definition(t.name),
				"FunctionDeclaration/identifier": t.function(t.definition(t.name)),
				"Parameters/identifier": t.variableName,

				// Literals
				number: t.number,
				string: t.string,

				// Comments
				LineComment: t.lineComment,

				// Operators
				"== != < <= > >= + - * /": t.operator,
				"= ! .": t.punctuation,
			}),

			indentNodeProp.add({
				Block: (context) =>
					context.column(context.node.parent.from) + context.unit,
				ClassBody: (context) =>
					context.column(context.node.parent.from) + context.unit,
			}),

			foldNodeProp.add({
				Block: foldInside,
				ClassBody: foldInside,
				ClassDeclaration: foldInside,
			}),
		],
	}),

	languageData: {
		commentTokens: { line: "//" },
		closeBrackets: { brackets: ["(", "[", "{", '"'] },
		indentOnInput: /^\s*(}|class\b|fun\b)/,
	},
});

const keywords = [
	"if",
	"else",
	"while",
	"for",
	"return",
	"and",
	"or",
	"var",
	"fun",
	"class",
	"this",
	"super",
	"true",
	"false",
	"nil",
	"print",
];

export function lox() {
	return new LanguageSupport(loxLanguage, [
		loxLanguage.data.of({
			autocomplete: [
				snippetCompletion("class ${name} {\n  init(args) {\n}}", {
					label: "class",
				}),
				completeFromList(keywords),
			],
		}),
	]);
}
