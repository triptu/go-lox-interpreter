import { RangeSetBuilder, StateEffect, StateField } from "@codemirror/state";
import { Decoration, MatchDecorator, ViewPlugin } from "@codemirror/view";
import { EditorView } from "codemirror";

const lineDeco = Decoration.line({ class: "cm-errorline" });

// A custom state effect for updating error lines
export const updateErrorLinesEffect =
	StateEffect.define<{ line: number; col?: number; text: string }[]>();

// This state field listens to the above effect to apply the decorations for error lines
export const errorLineDecorationField = StateField.define({
	create() {
		return Decoration.none;
	},
	update(oldDecorations, transaction) {
		// Check if error lines have been updated
		const updateEffect = transaction.effects.find((e) =>
			e.is(updateErrorLinesEffect),
		);
		if (!updateEffect) {
			// if doc has changed, clear the error lines
			if (transaction.docChanged) {
				return Decoration.none;
			}
			return oldDecorations;
		}
		const builder = new RangeSetBuilder<Decoration>();
		const errorLines = updateEffect.value;
		for (const errorLine of errorLines) {
			builder.add(
				transaction.newDoc.line(errorLine.line).from,
				transaction.newDoc.line(errorLine.line).from,
				lineDeco,
			);
		}
		return builder.finish();
	},
	provide: (f) => EditorView.decorations.from(f),
});

// syntax highlighting for lox keywords
const functionDeco = Decoration.mark({ class: "cm-lox-fun" });
const printDeco = Decoration.mark({ class: "cm-lox-print" });
const keywordDeco = Decoration.mark({ class: "cm-lox-keyword" });
const keywordDecorator = new MatchDecorator({
	regexp: /\b(fun|print|or|and)\b/g,
	decoration: (match) => {
		const keyword = match[1];
		switch (keyword) {
			case "fun":
				return functionDeco;
			case "print":
				return printDeco;
			default:
				return keywordDeco;
		}
	},
});
export const keywordDecorationPlugin = ViewPlugin.define(
	(view) => ({
		decorations: keywordDecorator.createDeco(view),
		update(u) {
			this.decorations = keywordDecorator.updateDeco(u, this.decorations);
		},
	}),
	{
		decorations: (v) => v.decorations,
	},
);
