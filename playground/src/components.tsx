import type { ComponentChildren } from "preact";
import type { JSX } from "preact/jsx-runtime";
import { ExternalLinkIcon, GithubIcon, WebsiteLinkIcon } from "./icons";

// https://github.com/lukeed/clsx/blob/master/src/lite.js
function clsx(...classes: string[]) {
	let i = 0;
	let className: string;
	let merged = "";
	for (; i < classes.length; i++) {
		className = classes[i];
		if (className && typeof className === "string") {
			merged += (merged && " ") + className;
		}
	}
	return merged;
}

export const Button = ({
	className,
	children,
	...props
}: {
	className?: string;
	children: ComponentChildren;
} & JSX.IntrinsicElements["button"]) => {
	return (
		<button
			class={clsx(
				className,
				"cursor-default bg-blue-600 hover:bg-blue-700 text-white font-bold h-9 px-4 text-sm rounded flex items-center gap-1 disabled:bg-gray-400",
			)}
			type="button"
			{...props}
		>
			{children}
		</button>
	);
};

export function Header() {
	return (
		<div class="rounded-t p-4 bg-gray-900 text-gray-50 flex items-center justify-between">
			<h1 class="text-2xl font-bold">GoLox Playground</h1>
			<div class="flex gap-4 items-center">
				<a
					href="https://github.com/triptu/go-lox-interpreter"
					target="_blank"
					rel="noopener noreferrer"
					class="text-gray-100 hover:text-white flex gap-2 items-center"
				>
					<GithubIcon class="sm:hidden" />
					<span class="hidden sm:inline-block">Github</span>
					<ExternalLinkIcon class="hidden sm:inline-block" />
				</a>

				<a
					href="https://tushartripathi.me/"
					target="_blank"
					rel="noopener noreferrer"
					class="text-gray-100 hover:text-white flex gap-2 items-center"
				>
					<WebsiteLinkIcon class="sm:hidden" />
					<span class="hidden sm:inline-block">Tushar</span>
					<ExternalLinkIcon class="hidden sm:inline-block" />
				</a>
			</div>
		</div>
	);
}
