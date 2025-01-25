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
	color = "blue",
	...props
}: {
	color?: "blue" | "red";
	className?: string;
	children: ComponentChildren;
} & JSX.IntrinsicElements["button"]) => {
	const colorClass =
		color === "blue"
			? "bg-blue-600 hover:bg-blue-700"
			: "bg-red-600 hover:bg-red-700";
	return (
		<button
			class={clsx(
				className,
				colorClass,
				"cursor-default font-bold h-9 px-4 text-sm rounded flex items-center gap-1 disabled:bg-gray-400 text-white",
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
			<div class="flex flex-col">
				<h1 class="text-2xl font-bold">
					GoLox <span class="hidden sm:inline">Playground</span>
				</h1>
				<span class="text-gray-300 text-sm hidden md:inline">
					A local in browser playground for Lox interpreter written in Golang{" "}
					<span class="hidden lg:inline">and compiled to WASM</span>
				</span>
				<span class="md:hidden text-gray-300 text-sm">
					Lox interpreter written in Golang
				</span>
			</div>

			<div class="flex gap-4 items-center">
				<a
					href="https://craftinginterpreters.com/the-lox-language.html"
					target="_blank"
					rel="noopener noreferrer"
					class="text-gray-100 hover:text-white flex gap-2 items-center"
				>
					<span>Lox</span>
					<ExternalLinkIcon />
				</a>

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
