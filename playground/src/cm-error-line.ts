import { RangeSetBuilder, StateEffect, StateField } from "@codemirror/state";
import { Decoration } from "@codemirror/view";
import { EditorView } from "codemirror";

const lineDeco = Decoration.line({ class: "cm-errorline" });

// Create a custom state effect for updating error lines
export const updateErrorLinesEffect =
	StateEffect.define<{ line: number; col?: number; text: string }[]>();

// Create a state field to manage error line decorations
export const errorLineDecorationField = StateField.define({
	create() {
		return Decoration.none;
	},
	update(oldDecorations, transaction) {
		// Check if error lines have been updated
		const updateEffect = transaction.effects.find((e) =>
			e.is(updateErrorLinesEffect),
		);
		if (!updateEffect) return oldDecorations;
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
