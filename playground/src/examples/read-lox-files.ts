// @ts-expect-error - runs at build time
import path from "node:path";
// @ts-expect-error - runs at build time
import { Glob } from "bun";

export async function readLoxFiles(): Promise<{ [filename: string]: string }> {
	try {
		const filesMap: { [filename: string]: string } = {};
		const glob = new Glob("./*.lox");
		// @ts-expect-error - runs at build time
		const currDir = import.meta.dir;
		for (const fileName of glob.scanSync({ cwd: currDir })) {
			const fullPath = path.join(currDir, fileName);
			const justFileName = path.basename(fullPath);
			// @ts-expect-error - runs at build time, we've bun
			filesMap[justFileName] = await Bun.file(fullPath).text();
		}
		return filesMap;
	} catch (dirError) {
		console.error("Error reading directory:", dirError);
		throw dirError;
	}
}
