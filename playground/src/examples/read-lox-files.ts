// @ts-expect-error - runs at build time
import path from "node:path";
// @ts-expect-error - runs at build time
import { Glob } from "bun";

// some examples are from - https://github.com/nanmu42/bluelox/tree/master/web/toy

export async function readLoxFiles(): Promise<{ [filename: string]: string }> {
	try {
		const filesMap: { [filename: string]: string } = {};
		const glob = new Glob("./*.lox");
		// @ts-expect-error - runs at build time
		const currDir = import.meta.dir;
		for (const fileName of glob.scanSync({ cwd: currDir })) {
			console.log("fileName", fileName);
			const fullPath = path.join(currDir, fileName);
			const justFileName = path.basename(fullPath);

			// https://github.com/oven-sh/bun/issues/7611
			// @ts-expect-error - runs at build time, we've bun
			filesMap[justFileName] = await Bun.readableStreamToText(
				// @ts-expect-error - runs at build time, we've bun
				Bun.file(fullPath).stream(),
			);
		}
		return filesMap;
	} catch (dirError) {
		console.error("Error reading directory:", dirError);
		throw dirError;
	}
}
