import { createRequire } from "node:module";
import { UNICODE, colors, consola, ensureFile, execa, execaCommand, findMonorepoRoot, formatFile, generatorContentHash, getPackages, getStagedFiles, gitAdd, outputJSON, readJSON, toPosixPath } from "@vben/node-utils";
import { cac } from "cac";
import { access, mkdtemp, readFile, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { basename, dirname, extname, join, relative } from "node:path";
import depcheck from "depcheck";
import { publint } from "publint";
import { formatMessage } from "publint/utils";
//#region package.json
var version = "5.7.0";
//#endregion
//#region src/check-circular/index.ts
const circularScannerCli = createRequire(import.meta.url).resolve("circular-dependency-scanner/dist/cli.js");
const DEFAULT_CONFIG$1 = {
	allowedExtensions: [
		".cjs",
		".js",
		".jsx",
		".mjs",
		".ts",
		".tsx",
		".vue"
	],
	ignoreDirs: [
		"dist",
		".turbo",
		"output",
		".cache",
		"scripts",
		"internal",
		"packages/effects/request/src/",
		"packages/@core/ui-kit/menu-ui/src/",
		"packages/@core/ui-kit/popup-ui/src/"
	],
	threshold: 0
};
const cache = /* @__PURE__ */ new Map();
async function detectCircularDependencies({ cwd, ignorePattern, staged }) {
	const tempDir = await mkdtemp(join(tmpdir(), "vsh-check-circular-"));
	const outputFile = join(tempDir, "circles.json");
	try {
		const args = [
			circularScannerCli,
			cwd,
			"--output",
			outputFile
		];
		if (staged) args.push("--absolute");
		args.push("--ignore", ignorePattern);
		await execa(process.execPath, args, { cwd });
		await access(outputFile);
		const output = await readFile(outputFile, "utf8");
		return JSON.parse(output);
	} catch (error) {
		if (error?.code === "ENOENT") return [];
		throw error;
	} finally {
		await rm(tempDir, {
			force: true,
			recursive: true
		});
	}
}
/**
* 格式化循环依赖的输出
* @param circles - 循环依赖结果
*/
function formatCircles(circles) {
	if (circles.length === 0) {
		console.log("✅ No circular dependencies found");
		return;
	}
	console.log("⚠️ Circular dependencies found:");
	circles.forEach((circle, index) => {
		console.log(`\nCircular dependency #${index + 1}:`);
		circle.forEach((file) => console.log(`  → ${file}`));
	});
}
/**
* 检查项目中的循环依赖
* @param options - 检查选项
* @param options.staged - 是否只检查暂存区文件
* @param options.verbose - 是否显示详细信息
* @param options.config - 自定义配置
* @returns Promise<void>
*/
async function checkCircular({ config = {}, staged, verbose }) {
	try {
		const finalConfig = {
			...DEFAULT_CONFIG$1,
			...config
		};
		const ignorePattern = `**/{${finalConfig.ignoreDirs.join(",")}}/**`;
		const cacheKey = `${staged}-${process.cwd()}-${ignorePattern}`;
		if (cache.has(cacheKey)) {
			const cachedResults = cache.get(cacheKey);
			if (cachedResults && verbose) formatCircles(cachedResults);
			return;
		}
		const results = await detectCircularDependencies({
			cwd: process.cwd(),
			ignorePattern,
			staged
		});
		if (staged) {
			let files = await getStagedFiles();
			const allowedExtensions = new Set(finalConfig.allowedExtensions);
			files = files.filter((file) => allowedExtensions.has(extname(file)));
			const circularFiles = [];
			for (const file of files) for (const result of results) if (result.flat().includes(file)) circularFiles.push(result);
			cache.set(cacheKey, circularFiles);
			if (verbose) formatCircles(circularFiles);
		} else {
			cache.set(cacheKey, results);
			if (verbose) formatCircles(results);
		}
		if (results.length > 0) console.log("\n⚠️ Warning: Circular dependencies found, please check and fix");
	} catch (error) {
		console.error("❌ Error checking circular dependencies:", error instanceof Error ? error.message : error);
	}
}
/**
* 定义检查循环依赖的命令
* @param cac - CAC实例
*/
function defineCheckCircularCommand(cac) {
	cac.command("check-circular").option("--staged", "Only check staged files").option("--verbose", "Show detailed information").option("--threshold <number>", "Threshold for circular dependencies", { default: 0 }).option("--ignore-dirs <dirs>", "Directories to ignore, comma separated").usage("Analyze project circular dependencies").action(async ({ ignoreDirs, staged, threshold, verbose }) => {
		await checkCircular({
			config: {
				threshold: Number(threshold),
				...ignoreDirs && { ignoreDirs: ignoreDirs.split(",") }
			},
			staged,
			verbose: verbose ?? true
		});
	});
}
//#endregion
//#region src/check-dep/index.ts
const DEFAULT_CONFIG = {
	ignoreMatches: [
		"vite",
		"vitest",
		"tsdown",
		"@vben/tailwind-config",
		"@vben/tsconfig",
		"@vben/vite-config",
		"@types/*",
		"@vben-core/design"
	],
	ignorePackages: [
		"@vben/commitlint-config",
		"@vben/eslint-config",
		"@vben/node-utils",
		"@vben/oxfmt-config",
		"@vben/oxlint-config",
		"@vben/stylelint-config",
		"@vben/tsconfig",
		"@vben/vite-config",
		"@vben/vsh"
	],
	ignorePatterns: [
		"dist",
		"node_modules",
		"public"
	]
};
/**
* 清理依赖检查结果
* @param unused - 依赖检查结果
*/
function cleanDepcheckResult(unused) {
	Reflect.deleteProperty(unused.missing, "file:");
	Object.keys(unused.missing).forEach((key) => {
		unused.missing[key] = (unused.missing[key] || []).filter((item) => !item.startsWith("/"));
		if (unused.missing[key].length === 0) Reflect.deleteProperty(unused.missing, key);
	});
}
/**
* 格式化依赖检查结果
* @param pkgName - 包名
* @param unused - 依赖检查结果
*/
function formatDepcheckResult(pkgName, unused) {
	if (!(Object.keys(unused.missing).length > 0 || unused.dependencies.length > 0 || unused.devDependencies.length > 0)) return;
	console.log("\n📦 Package:", pkgName);
	if (Object.keys(unused.missing).length > 0) {
		console.log("❌ Missing dependencies:");
		Object.entries(unused.missing).forEach(([dep, files]) => {
			console.log(`  - ${dep}:`);
			files.forEach((file) => console.log(`    → ${file}`));
		});
	}
	if (unused.dependencies.length > 0) {
		console.log("⚠️ Unused dependencies:");
		unused.dependencies.forEach((dep) => console.log(`  - ${dep}`));
	}
	if (unused.devDependencies.length > 0) {
		console.log("⚠️ Unused devDependencies:");
		unused.devDependencies.forEach((dep) => console.log(`  - ${dep}`));
	}
}
/**
* 运行依赖检查
* @param config - 配置选项
*/
async function runDepcheck(config = {}) {
	try {
		const finalConfig = {
			...DEFAULT_CONFIG,
			...config
		};
		const { packages } = await getPackages();
		let hasIssues = false;
		await Promise.all(packages.map(async (pkg) => {
			if (finalConfig.ignorePackages.includes(pkg.packageJson.name)) return;
			const unused = await depcheck(pkg.dir, {
				ignoreMatches: finalConfig.ignoreMatches,
				ignorePatterns: finalConfig.ignorePatterns
			});
			cleanDepcheckResult(unused);
			if (Object.keys(unused.missing).length > 0 || unused.dependencies.length > 0 || unused.devDependencies.length > 0) {
				hasIssues = true;
				formatDepcheckResult(pkg.packageJson.name, unused);
			}
		}));
		if (!hasIssues) console.log("\n✅ Dependency check completed, no issues found");
	} catch (error) {
		console.error("❌ Dependency check failed:", error instanceof Error ? error.message : error);
	}
}
/**
* 定义依赖检查命令
* @param cac - CAC实例
*/
function defineDepcheckCommand(cac) {
	cac.command("check-dep").option("--ignore-packages <packages>", "Packages to ignore, comma separated").option("--ignore-matches <matches>", "Dependency patterns to ignore, comma separated").option("--ignore-patterns <patterns>", "File patterns to ignore, comma separated").usage("Analyze project dependencies").action(async ({ ignoreMatches, ignorePackages, ignorePatterns }) => {
		await runDepcheck({
			...ignorePackages && { ignorePackages: ignorePackages.split(",") },
			...ignoreMatches && { ignoreMatches: ignoreMatches.split(",") },
			...ignorePatterns && { ignorePatterns: ignorePatterns.split(",") }
		});
	});
}
//#endregion
//#region src/code-workspace/index.ts
const CODE_WORKSPACE_FILE = join("vben-admin.code-workspace");
async function createCodeWorkspace({ autoCommit = false, spaces = 2 }) {
	const { packages, rootDir } = await getPackages();
	let folders = packages.map((pkg) => {
		const { dir, packageJson } = pkg;
		return {
			name: packageJson.name,
			path: toPosixPath(relative(rootDir, dir))
		};
	});
	folders = folders.filter(Boolean);
	const monorepoRoot = findMonorepoRoot();
	const outputPath = join(monorepoRoot, CODE_WORKSPACE_FILE);
	await outputJSON(outputPath, { folders }, spaces);
	await formatFile(outputPath);
	if (autoCommit) await gitAdd(CODE_WORKSPACE_FILE, monorepoRoot);
}
async function runCodeWorkspace({ autoCommit, spaces }) {
	await createCodeWorkspace({
		autoCommit,
		spaces
	});
	if (autoCommit) return;
	consola.log("");
	consola.success(colors.green(`${CODE_WORKSPACE_FILE} is updated!`));
	consola.log("");
}
function defineCodeWorkspaceCommand(cac) {
	cac.command("code-workspace").usage("Update the `.code-workspace` file").option("--spaces [number]", ".code-workspace JSON file spaces.", { default: 2 }).option("--auto-commit", "auto commit .code-workspace JSON file.", { default: false }).action(runCodeWorkspace);
}
//#endregion
//#region src/lint/index.ts
async function runLint({ format }) {
	if (format) {
		await execaCommand(`stylelint "**/*.{vue,css,less,scss}" --cache --fix`, { stdio: "inherit" });
		await execaCommand(`oxfmt`, { stdio: "inherit" });
		await execaCommand(`oxlint --fix`, { stdio: "inherit" });
		await execaCommand(`eslint . --cache --fix`, { stdio: "inherit" });
		return;
	}
	await Promise.all([
		execaCommand(`oxfmt --check`, { stdio: "inherit" }),
		execaCommand(`oxlint`, { stdio: "inherit" }),
		execaCommand(`eslint . --cache`, { stdio: "inherit" }),
		execaCommand(`stylelint "**/*.{vue,css,less,scss}" --cache`, { stdio: "inherit" })
	]);
}
function defineLintCommand(cac) {
	cac.command("lint").usage("Batch execute project lint check.").option("--format", "Format lint problem.").action(runLint);
}
//#endregion
//#region src/publint/index.ts
const CACHE_FILE = join("node_modules", ".cache", "publint", ".pkglintcache.json");
/**
* Get files that require lint
* @param files
*/
async function getLintFiles(files = []) {
	const lintFiles = [];
	if (files?.length > 0) return files.filter((file) => basename(file) === "package.json");
	const { packages } = await getPackages();
	for (const { dir } of packages) lintFiles.push(join(dir, "package.json"));
	return lintFiles;
}
function getCacheFile() {
	return join(findMonorepoRoot(), CACHE_FILE);
}
async function readCache(cacheFile) {
	try {
		await ensureFile(cacheFile);
		return await readJSON(cacheFile);
	} catch {
		return {};
	}
}
async function runPublint(files, { check }) {
	const lintFiles = await getLintFiles(files);
	const cacheFile = getCacheFile();
	const cache = await readCache(cacheFile);
	const results = await Promise.all(lintFiles.map(async (file) => {
		try {
			const pkgJson = await readJSON(file);
			if (pkgJson.private) return null;
			Reflect.deleteProperty(pkgJson, "dependencies");
			Reflect.deleteProperty(pkgJson, "devDependencies");
			Reflect.deleteProperty(pkgJson, "peerDependencies");
			const hash = generatorContentHash(JSON.stringify(pkgJson));
			const publintResult = cache?.[file]?.hash === hash ? cache?.[file]?.result ?? [] : await publint({
				level: "suggestion",
				pkgDir: dirname(file),
				strict: true
			});
			cache[file] = {
				hash,
				result: publintResult
			};
			return {
				pkgJson,
				pkgPath: file,
				publintResult
			};
		} catch {
			return null;
		}
	}));
	await outputJSON(cacheFile, cache);
	printResult(results, check);
}
function printResult(results, check) {
	let errorCount = 0;
	let warningCount = 0;
	let suggestionsCount = 0;
	for (const result of results) {
		if (!result) continue;
		const { pkgJson, pkgPath, publintResult } = result;
		const messages = publintResult?.messages ?? [];
		if (messages?.length < 1) continue;
		consola.log("");
		consola.log(pkgPath);
		for (const message of messages) {
			switch (message.type) {
				case "error":
					errorCount++;
					break;
				case "suggestion":
					suggestionsCount++;
					break;
				case "warning":
					warningCount++;
					break;
			}
			const ruleUrl = `https://publint.dev/rules#${message.code.toLocaleLowerCase()}`;
			consola.log(`  ${formatMessage(message, pkgJson)}${colors.dim(` ${ruleUrl}`)}`);
		}
	}
	const totalCount = warningCount + errorCount + suggestionsCount;
	if (totalCount > 0) {
		consola.error(colors.red(`${UNICODE.FAILURE} ${totalCount} problem (${errorCount} errors, ${warningCount} warnings, ${suggestionsCount} suggestions)`));
		!check && process.exit(1);
	} else consola.log(colors.green(`${UNICODE.SUCCESS} No problem`));
}
function definePubLintCommand(cac) {
	cac.command("publint [...files]").usage("Check if the monorepo package conforms to the publint standard.").option("--check", "Only errors are checked, no program exit is performed.").action(runPublint);
}
//#endregion
//#region src/index.ts
const COMMAND_DESCRIPTIONS = {
	"check-circular": "Check for circular dependencies",
	"check-dep": "Check for unused dependencies",
	"code-workspace": "Manage VS Code workspace settings",
	lint: "Run linting on the project",
	publint: "Check package.json files for publishing standards"
};
/**
* Initialize and run the CLI
*/
async function main() {
	try {
		const vsh = cac("vsh");
		defineLintCommand(vsh);
		definePubLintCommand(vsh);
		defineCodeWorkspaceCommand(vsh);
		defineCheckCircularCommand(vsh);
		defineDepcheckCommand(vsh);
		vsh.usage("vsh <command> [options]");
		vsh.help();
		vsh.version(version);
		vsh.parse(void 0, { run: false });
		if (!vsh.matchedCommand && vsh.args.length > 0) {
			const unknownCmd = String(vsh.args[0]);
			consola.error(colors.red(`Invalid command: ${unknownCmd}`), "\n", colors.yellow("Available commands:"), "\n", Object.entries(COMMAND_DESCRIPTIONS).map(([name, desc]) => `  ${colors.cyan(name)} - ${desc}`).join("\n"));
			process.exit(1);
		}
		await vsh.runMatchedCommand();
	} catch (error) {
		consola.error(colors.red("An unexpected error occurred:"), "\n", error instanceof Error ? error.message : error);
		process.exit(1);
	}
}
main().catch((error) => {
	consola.error(colors.red("Failed to start CLI:"), "\n", error instanceof Error ? error.message : error);
	process.exit(1);
});
//#endregion
export {};
