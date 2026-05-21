import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import dayjs from "dayjs";
import timezone from "dayjs/plugin/timezone.js";
import utc from "dayjs/plugin/utc.js";
import CryptoJS from "crypto-js";
import { JSEncrypt } from "jsencrypt";
import { isFunction, isObject, isString } from "@vue/shared";
import { createDefu, createDefu as createMerge, defu as merge } from "defu";
import { get, isEqual, set } from "es-toolkit/compat";
import cloneDeep from "lodash.clonedeep";
//#region src/utils/cn.ts
function cn(...inputs) {
	return twMerge(clsx(inputs));
}
//#endregion
//#region src/utils/date.ts
dayjs.extend(utc);
dayjs.extend(timezone);
function formatDate(time, format = "YYYY-MM-DD") {
	if (!time) return "";
	try {
		const date = dayjs.isDayjs(time) ? time : dayjs(time);
		if (!date.isValid()) throw new Error("Invalid date");
		return date.tz().format(format);
	} catch (error) {
		console.error(`Error formatting date: ${error}`);
		return String(time ?? "");
	}
}
function formatDateTime(time) {
	return formatDate(time, "YYYY-MM-DD HH:mm:ss");
}
function formatDate2(date, format) {
	if (!date) return "";
	return date ? dayjs(date).format(format ?? "YYYY-MM-DD HH:mm:ss") : "";
}
function isDate(value) {
	return value instanceof Date;
}
function isDayjsObject(value) {
	return dayjs.isDayjs(value);
}
/**
* element plus 的时间 Formatter 实现，使用 YYYY-MM-DD HH:mm:ss 格式
*
* @param _row
* @param _column
* @param cellValue 字段值
*/
function dateFormatter(_row, _column, cellValue) {
	return cellValue ? formatDate(cellValue)?.toString() || "" : "";
}
/**
* 获取当前时区
* @returns 当前时区
*/
const getSystemTimezone = () => {
	return dayjs.tz.guess();
};
/**
* 自定义设置的时区
*/
let currentTimezone = getSystemTimezone();
/**
* 设置默认时区
* @param timezone
*/
const setCurrentTimezone = (timezone) => {
	currentTimezone = timezone || getSystemTimezone();
	dayjs.tz.setDefault(currentTimezone);
};
/**
* 获取设置的时区
* @returns 设置的时区
*/
const getCurrentTimezone = () => {
	return currentTimezone;
};
//#endregion
//#region src/utils/diff.ts
function arraysEqual(a, b) {
	if (a.length !== b.length) return false;
	const counter = /* @__PURE__ */ new Map();
	for (const value of a) counter.set(value, (counter.get(value) || 0) + 1);
	for (const value of b) {
		const count = counter.get(value);
		if (count === void 0 || count === 0) return false;
		counter.set(value, count - 1);
	}
	return true;
}
function diff(obj1, obj2) {
	function findDifferences(o1, o2) {
		if (Array.isArray(o1) && Array.isArray(o2)) {
			if (!arraysEqual(o1, o2)) return o2;
			return;
		}
		if (typeof o1 === "object" && typeof o2 === "object" && o1 !== null && o2 !== null) {
			const diffResult = {};
			new Set([...Object.keys(o1), ...Object.keys(o2)]).forEach((key) => {
				const valueDiff = findDifferences(o1[key], o2[key]);
				if (valueDiff !== void 0) diffResult[key] = valueDiff;
			});
			return Object.keys(diffResult).length > 0 ? diffResult : void 0;
		}
		return o1 === o2 ? void 0 : o2;
	}
	return findDifferences(obj1, obj2);
}
//#endregion
//#region src/utils/dom.ts
/**
* 获取元素可见信息
* @param element
*/
function getElementVisibleRect(element) {
	if (!element) return {
		bottom: 0,
		height: 0,
		left: 0,
		right: 0,
		top: 0,
		width: 0
	};
	const rect = element.getBoundingClientRect();
	const viewHeight = Math.max(document.documentElement.clientHeight, window.innerHeight);
	const top = Math.max(rect.top, 0);
	const bottom = Math.min(rect.bottom, viewHeight);
	const viewWidth = Math.max(document.documentElement.clientWidth, window.innerWidth);
	const left = Math.max(rect.left, 0);
	const right = Math.min(rect.right, viewWidth);
	if (top >= viewHeight || bottom <= 0 || left >= viewWidth || right <= 0) return {
		bottom: 0,
		height: 0,
		left: 0,
		right: 0,
		top: 0,
		width: 0
	};
	return {
		bottom,
		height: Math.max(0, bottom - top),
		left,
		right,
		top,
		width: Math.max(0, right - left)
	};
}
function getScrollbarWidth() {
	const scrollDiv = document.createElement("div");
	scrollDiv.style.visibility = "hidden";
	scrollDiv.style.overflow = "scroll";
	scrollDiv.style.position = "absolute";
	scrollDiv.style.top = "-9999px";
	document.body.append(scrollDiv);
	const innerDiv = document.createElement("div");
	scrollDiv.append(innerDiv);
	const scrollbarWidth = scrollDiv.offsetWidth - innerDiv.offsetWidth;
	scrollDiv.remove();
	return scrollbarWidth;
}
function needsScrollbar() {
	const doc = document.documentElement;
	const body = document.body;
	const overflowY = window.getComputedStyle(body).overflowY;
	if (overflowY === "scroll" || overflowY === "auto") return doc.scrollHeight > window.innerHeight;
	return doc.scrollHeight > window.innerHeight;
}
function triggerWindowResize() {
	const resizeEvent = new Event("resize");
	window.dispatchEvent(resizeEvent);
}
//#endregion
//#region src/utils/window.ts
/**
* 新窗口打开URL。
*
* @param url - 需要打开的网址。
* @param options - 打开窗口的选项。
*/
function openWindow(url, options = {}) {
	const { noopener = true, noreferrer = true, target = "_blank" } = options;
	const features = [noopener && "noopener=yes", noreferrer && "noreferrer=yes"].filter(Boolean).join(",");
	window.open(url, target, features);
}
/**
* 在新窗口中打开路由。
* @param path
*/
function openRouteInNewWindow(path) {
	const { hash, origin } = location;
	const fullPath = path.startsWith("/") ? path : `/${path}`;
	openWindow(`${origin}${hash && !fullPath.startsWith("/#") ? "/#" : ""}${fullPath}`, { target: "_blank" });
}
//#endregion
//#region src/utils/download.ts
const DEFAULT_FILENAME = "downloaded_file";
/**
* 通过 URL 下载文件，支持跨域
* @throws {Error} - 当下载失败时抛出错误
*/
async function downloadFileFromUrl({ fileName, source, target = "_blank" }) {
	if (!source || typeof source !== "string") throw new Error("Invalid URL.");
	const isChrome = window.navigator.userAgent.toLowerCase().includes("chrome");
	const isSafari = window.navigator.userAgent.toLowerCase().includes("safari");
	if (/iP/.test(window.navigator.userAgent)) {
		console.error("Your browser does not support download!");
		return;
	}
	if (isChrome || isSafari) {
		triggerDownload(source, resolveFileName(source, fileName));
		return;
	}
	if (!source.includes("?")) source += "?download";
	openWindow(source, { target });
}
/**
* 下载图片（允许跨域）
* @param url - 图片 URL
* @param canvasWidth - 画布宽度
* @param canvasHeight - 画布高度
* @param drawWithImageSize - 将图片绘制在画布上时带上图片的宽高值, 默认是要带上的
* @returns
*/
function downloadImageByCanvas({ url, canvasWidth, canvasHeight, drawWithImageSize = true }) {
	const image = new Image();
	image.src = url;
	image.addEventListener("load", () => {
		const canvas = document.createElement("canvas");
		canvas.width = canvasWidth || image.width;
		canvas.height = canvasHeight || image.height;
		const ctx = canvas.getContext("2d");
		ctx?.clearRect(0, 0, canvas.width, canvas.height);
		if (drawWithImageSize) ctx.drawImage(image, 0, 0, image.width, image.height);
		else ctx.drawImage(image, 0, 0);
		downloadFileFromImageUrl({
			source: canvas.toDataURL("image/png"),
			fileName: "image.png"
		});
	});
}
/**
* 通过 Base64 下载文件
*/
function downloadFileFromBase64({ fileName, source }) {
	if (!source || typeof source !== "string") throw new Error("Invalid Base64 data.");
	triggerDownload(source, fileName || DEFAULT_FILENAME);
}
/**
* 通过图片 URL 下载图片文件
*/
async function downloadFileFromImageUrl({ fileName, source }) {
	downloadFileFromBase64({
		fileName,
		source: await urlToBase64(source)
	});
}
/**
* 通过 Blob 下载文件
*/
function downloadFileFromBlob({ fileName = DEFAULT_FILENAME, source }) {
	if (!(source instanceof Blob)) throw new TypeError("Invalid Blob data.");
	triggerDownload(URL.createObjectURL(source), fileName);
}
/**
* 下载文件，支持 Blob、字符串和其他 BlobPart 类型
*/
function downloadFileFromBlobPart({ fileName = DEFAULT_FILENAME, source }) {
	const blob = source instanceof Blob ? source : new Blob([source], { type: "application/octet-stream" });
	triggerDownload(URL.createObjectURL(blob), fileName);
}
/**
* @description: base64 to blob
*/
function dataURLtoBlob(base64Buf) {
	const arr = base64Buf.split(",");
	const mime = arr[0].match(/:(.*?);/)[1];
	const bstr = window.atob(arr[1]);
	let n = bstr.length;
	const u8arr = new Uint8Array(n);
	while (n--) u8arr[n] = bstr.codePointAt(n);
	return new Blob([u8arr], { type: mime });
}
/**
* img url to base64
* @param url
*/
function urlToBase64(url, mineType) {
	return new Promise((resolve, reject) => {
		let canvas = document.createElement("CANVAS");
		const ctx = canvas?.getContext("2d");
		const img = new Image();
		img.crossOrigin = "";
		img.addEventListener("load", () => {
			if (!canvas || !ctx) return reject(/* @__PURE__ */ new Error("Failed to create canvas."));
			canvas.height = img.height;
			canvas.width = img.width;
			ctx.drawImage(img, 0, 0);
			const dataURL = canvas.toDataURL(mineType || "image/png");
			canvas = null;
			resolve(dataURL);
		});
		img.src = url;
	});
}
/**
* 将 Base64 字符串转换为文件对象
* @param base64 - Base64 字符串
* @param fileName - 文件名
* @returns File 对象
*/
function base64ToFile(base64, fileName) {
	if (!base64 || typeof base64 !== "string") throw new Error("base64 参数必须是非空字符串");
	const data = base64.split(",");
	if (data.length !== 2 || !data[0] || !data[1]) throw new Error("无效的 base64 格式");
	const typeMatch = data[0].match(/:(.*?);/);
	if (!typeMatch || !typeMatch[1]) throw new Error("无法解析 base64 类型信息");
	const type = typeMatch[1];
	const typeParts = type.split("/");
	if (typeParts.length !== 2 || !typeParts[1]) throw new Error("无效的 MIME 类型格式");
	const suffix = typeParts[1];
	try {
		const bstr = window.atob(data[1]);
		const n = bstr.length;
		const u8arr = new Uint8Array(n);
		for (let i = 0; i < n; i++) u8arr[i] = bstr.charCodeAt(i);
		return new File([u8arr], `${fileName}.${suffix}`, { type });
	} catch (error) {
		throw new Error(`Base64 解码失败: ${error instanceof Error ? error.message : "未知错误"}`, { cause: error });
	}
}
/**
* 通用下载触发函数
* @param href - 文件下载的 URL
* @param fileName - 下载文件的名称，如果未提供则自动识别
* @param revokeDelay - 清理 URL 的延迟时间 (毫秒)
*/
function triggerDownload(href, fileName, revokeDelay = 100) {
	const finalFileName = fileName || "downloaded_file";
	const link = document.createElement("a");
	link.href = href;
	link.download = finalFileName;
	link.style.display = "none";
	if (link.download === void 0) link.setAttribute("target", "_blank");
	document.body.append(link);
	link.click();
	link.remove();
	setTimeout(() => URL.revokeObjectURL(href), revokeDelay);
}
function resolveFileName(url, fileName) {
	return fileName || url.slice(url.lastIndexOf("/") + 1) || DEFAULT_FILENAME;
}
//#endregion
//#region src/utils/encrypt.ts
/**
* API 加解密工具类
* 支持 AES 和 RSA 加密算法
*/
/**
* AES 加密工具类
*/
const AES = {
	encrypt(data, key) {
		try {
			if (!key) throw new Error("AES 加密密钥不能为空");
			if (key.length !== 32 && key.length !== 16) throw new Error(`AES 加密密钥长度必须为 32 位或 16 位，当前长度: ${key.length}`);
			const keyUtf8 = CryptoJS.enc.Utf8.parse(key);
			return CryptoJS.AES.encrypt(data, keyUtf8, {
				mode: CryptoJS.mode.ECB,
				padding: CryptoJS.pad.Pkcs7
			}).toString();
		} catch (error) {
			console.error("AES 加密失败:", error);
			throw error;
		}
	},
	decrypt(encryptedData, key) {
		try {
			if (!key) throw new Error("AES 解密密钥不能为空");
			if (key.length !== 32) throw new Error(`AES 解密密钥长度必须为 32 位，当前长度: ${key.length}`);
			if (!encryptedData) throw new Error("AES 解密数据不能为空");
			const keyUtf8 = CryptoJS.enc.Utf8.parse(key);
			const result = CryptoJS.AES.decrypt(encryptedData, keyUtf8, {
				mode: CryptoJS.mode.ECB,
				padding: CryptoJS.pad.Pkcs7
			}).toString(CryptoJS.enc.Utf8);
			if (!result) throw new Error("AES 解密结果为空，可能是密钥错误或数据损坏");
			return result;
		} catch (error) {
			console.error("AES 解密失败:", error);
			throw error;
		}
	}
};
/**
* RSA 加密工具类
*/
const RSA = {
	encrypt(data, publicKey) {
		try {
			if (!publicKey) throw new Error("RSA 公钥不能为空");
			const encryptor = new JSEncrypt();
			encryptor.setPublicKey(publicKey);
			const result = encryptor.encrypt(data);
			if (result === false) throw new Error("RSA 加密失败，可能是公钥格式错误或数据过长");
			return result;
		} catch (error) {
			console.error("RSA 加密失败:", error);
			throw error;
		}
	},
	decrypt(encryptedData, privateKey) {
		try {
			if (!privateKey) throw new Error("RSA 私钥不能为空");
			if (!encryptedData) throw new Error("RSA 解密数据不能为空");
			const encryptor = new JSEncrypt();
			encryptor.setPrivateKey(privateKey);
			const result = encryptor.decrypt(encryptedData);
			if (result === false) throw new Error("RSA 解密失败，可能是私钥错误或数据损坏");
			return result;
		} catch (error) {
			console.error("RSA 解密失败:", error);
			throw error;
		}
	}
};
/**
* API 加解密主类
*/
var ApiEncrypt = class {
	config;
	constructor(config) {
		this.config = config;
	}
	/**
	* 解密响应数据
	* @param encryptedData 加密的响应数据
	* @returns 解密后的数据
	*/
	decryptResponse(encryptedData) {
		if (!this.config.enable) return encryptedData;
		try {
			let decryptedData = "";
			if (this.config.algorithm.toUpperCase() === "AES") {
				if (!this.config.responseKey) throw new Error("AES 响应解密密钥未配置");
				decryptedData = AES.decrypt(encryptedData, this.config.responseKey);
			} else if (this.config.algorithm.toUpperCase() === "RSA") {
				if (!this.config.responseKey) throw new Error("RSA 私钥未配置");
				decryptedData = RSA.decrypt(encryptedData, this.config.responseKey);
				if (decryptedData === false) throw new Error("RSA 解密失败");
			} else throw new Error(`不支持的解密算法: ${this.config.algorithm}`);
			if (!decryptedData) throw new Error("解密结果为空");
			try {
				return JSON.parse(decryptedData);
			} catch {
				return decryptedData;
			}
		} catch (error) {
			console.error("响应数据解密失败:", error);
			throw error;
		}
	}
	/**
	* 加密请求数据
	* @param data 要加密的数据
	* @returns 加密后的数据
	*/
	encryptRequest(data) {
		if (!this.config.enable) return data;
		try {
			const jsonData = typeof data === "string" ? data : JSON.stringify(data);
			if (this.config.algorithm.toUpperCase() === "AES") {
				if (!this.config.requestKey) throw new Error("AES 请求加密密钥未配置");
				return AES.encrypt(jsonData, this.config.requestKey);
			} else if (this.config.algorithm.toUpperCase() === "RSA") {
				if (!this.config.requestKey) throw new Error("RSA 公钥未配置");
				const result = RSA.encrypt(jsonData, this.config.requestKey);
				if (result === false) throw new Error("RSA 加密失败");
				return result;
			} else throw new Error(`不支持的加密算法: ${this.config.algorithm}`);
		} catch (error) {
			console.error("请求数据加密失败:", error);
			throw error;
		}
	}
	/**
	* 获取加密头名称
	*/
	getEncryptHeader() {
		return this.config.header;
	}
};
/**
* 创建基于环境变量的 API 加解密实例
* @param env 环境变量对象
* @returns ApiEncrypt 实例
*/
function createApiEncrypt(env) {
	return new ApiEncrypt({
		enable: env.VITE_APP_API_ENCRYPT_ENABLE === "true",
		header: env.VITE_APP_API_ENCRYPT_HEADER || "X-Api-Encrypt",
		algorithm: env.VITE_APP_API_ENCRYPT_ALGORITHM || "AES",
		requestKey: env.VITE_APP_API_ENCRYPT_REQUEST_KEY || "",
		responseKey: env.VITE_APP_API_ENCRYPT_RESPONSE_KEY || ""
	});
}
//#endregion
//#region src/utils/inference.ts
/**
* 检查传入的值是否为undefined。
*
* @param {unknown} value 要检查的值。
* @returns {boolean} 如果值是undefined，返回true，否则返回false。
*/
function isUndefined(value) {
	return value === void 0;
}
/**
* 检查传入的值是否为boolean
* @param value
* @returns 如果值是布尔值，返回true，否则返回false。
*/
function isBoolean(value) {
	return typeof value === "boolean";
}
/**
* 检查传入的值是否为空。
*
* 以下情况将被认为是空：
* - 值为null。
* - 值为undefined。
* - 值为一个空字符串。
* - 值为一个长度为0的数组。
* - 值为一个没有元素的Map或Set。
* - 值为一个没有属性的对象。
*
* @param {T} value 要检查的值。
* @returns {boolean} 如果值为空，返回true，否则返回false。
*/
function isEmpty(value) {
	if (value === null || value === void 0) return true;
	if (Array.isArray(value) || isString(value)) return value.length === 0;
	if (value instanceof Map || value instanceof Set) return value.size === 0;
	if (isObject(value)) return Object.keys(value).length === 0;
	return false;
}
/**
* 检查传入的字符串是否为有效的HTTP或HTTPS URL。
*
* @param {string} url 要检查的字符串。
* @return {boolean} 如果字符串是有效的HTTP或HTTPS URL，返回true，否则返回false。
*/
function isHttpUrl(url) {
	if (!url) return false;
	return /^https?:\/\/.*$/.test(url);
}
/**
* 检查传入的值是否为window对象。
*
* @param {any} value 要检查的值。
* @returns {boolean} 如果值是window对象，返回true，否则返回false。
*/
function isWindow(value) {
	return typeof window !== "undefined" && value !== null && value === value.window;
}
/**
* 检查当前运行环境是否为Mac OS。
*
* 这个函数通过检查navigator.userAgent字符串来判断当前运行环境。
* 如果userAgent字符串中包含"macintosh"或"mac os x"（不区分大小写），则认为当前环境是Mac OS。
*
* @returns {boolean} 如果当前环境是Mac OS，返回true，否则返回false。
*/
function isMacOs() {
	return /macintosh|mac os x/i.test(navigator.userAgent);
}
/**
* 检查当前运行环境是否为Windows OS。
*
* 这个函数通过检查navigator.userAgent字符串来判断当前运行环境。
* 如果userAgent字符串中包含"windows"或"win32"（不区分大小写），则认为当前环境是Windows OS。
*
* @returns {boolean} 如果当前环境是Windows OS，返回true，否则返回false。
*/
function isWindowsOs() {
	return /windows|win32/i.test(navigator.userAgent);
}
/**
* 检查传入的值是否为数字
* @param value
*/
function isNumber(value) {
	return typeof value === "number" && Number.isFinite(value);
}
/**
* Returns the first value in the provided list that is neither `null` nor `undefined`.
*
* This function iterates over the input values and returns the first one that is
* not strictly equal to `null` or `undefined`. If all values are either `null` or
* `undefined`, it returns `undefined`.
*
* @template T - The type of the input values.
* @param {...(T | null | undefined)[]} values - A list of values to evaluate.
* @returns {T | undefined} - The first value that is not `null` or `undefined`, or `undefined` if none are found.
*
* @example
* // Returns 42 because it is the first non-null, non-undefined value.
* getFirstNonNullOrUndefined(undefined, null, 42, 'hello'); // 42
*
* @example
* // Returns 'hello' because it is the first non-null, non-undefined value.
* getFirstNonNullOrUndefined(null, undefined, 'hello', 123); // 'hello'
*
* @example
* // Returns undefined because all values are either null or undefined.
* getFirstNonNullOrUndefined(undefined, null); // undefined
*/
function getFirstNonNullOrUndefined(...values) {
	for (const value of values) if (value !== void 0 && value !== null) return value;
}
//#endregion
//#region src/utils/formatNumber.ts
/**
* 将一个整数转换为分数保留传入的小数
* @param num
* @param digit
*/
function formatToFractionDigit(num, digit = 2) {
	if (isUndefined(num)) return "0.00";
	return ((isString(num) ? Number.parseFloat(num) : num) / 100).toFixed(digit);
}
/**
* 将一个整数转换为分数保留两位小数
* @param num
*/
function formatToFraction(num) {
	return formatToFractionDigit(num, 2);
}
/**
* 将一个数转换为 1.00 这样
* 数据呈现的时候使用
*
* @param num 整数
*/
function floatToFixed2(num) {
	let str = "0.00";
	if (isUndefined(num)) return str;
	const f = formatToFraction(num);
	const decimalPart = f.toString().split(".")[1];
	switch (decimalPart ? decimalPart.length : 0) {
		case 0:
			str = `${f.toString()}.00`;
			break;
		case 1:
			str = `${f.toString()}0`;
			break;
		case 2:
			str = f.toString();
			break;
	}
	return str;
}
/**
* 将一个分数转换为整数
* @param num
*/
function convertToInteger(num) {
	if (isUndefined(num)) return 0;
	const parsedNumber = isString(num) ? Number.parseFloat(num) : num;
	return Math.round(parsedNumber * 100);
}
/**
* 元转分
*/
function yuanToFen(amount) {
	return convertToInteger(amount);
}
/**
* 分转元
*/
function fenToYuan(price) {
	return formatToFraction(price);
}
const fenToYuanFormat = (_, __, cellValue, ___) => {
	return `￥${floatToFixed2(cellValue)}`;
};
/**
* 计算环比
*
* @param value 当前数值
* @param reference 对比数值
*/
function calculateRelativeRate(value, reference) {
	if (!reference || reference === 0) return 0;
	return Number.parseFloat((100 * ((value || 0) - reference) / reference).toFixed(0));
}
const ERP_COUNT_DIGIT = 3;
const ERP_PRICE_DIGIT = 2;
/**
* 【ERP】格式化 Input 数字
*
* 例如说：库存数量
*
* @param num 数量
* @package
* @return 格式化后的数量
*/
function erpNumberFormatter(num, digit) {
	if (num === null || num === void 0) return "";
	if (typeof num === "string") num = Number.parseFloat(num);
	if (Number.isNaN(num)) return "";
	return num.toFixed(digit);
}
/**
* 【ERP】格式化数量，保留三位小数
*
* 例如说：库存数量
*
* @param num 数量
* @return 格式化后的数量
*/
function erpCountInputFormatter(num) {
	return erpNumberFormatter(num, ERP_COUNT_DIGIT);
}
/**
* 【ERP】格式化数量，保留三位小数
*
* @param cellValue 数量
* @return 格式化后的数量
*/
function erpCountTableColumnFormatter(cellValue) {
	return erpNumberFormatter(cellValue, ERP_COUNT_DIGIT);
}
/**
* 【ERP】格式化金额，保留二位小数
*
* 例如说：库存数量
*
* @param num 数量
* @return 格式化后的数量
*/
function erpPriceInputFormatter(num) {
	return erpNumberFormatter(num, ERP_PRICE_DIGIT);
}
/**
* 【ERP】格式化金额，保留二位小数
*
* @param cellValue 数量
* @return 格式化后的数量
*/
function erpPriceTableColumnFormatter(cellValue) {
	return erpNumberFormatter(cellValue, ERP_PRICE_DIGIT);
}
/**
* 【ERP】价格计算，四舍五入保留两位小数
*
* @param price 价格
* @param count 数量
* @return 总价格。如果有任一为空，则返回 undefined
*/
function erpPriceMultiply(price, count) {
	if (isEmpty(price) || isEmpty(count)) return void 0;
	return Number.parseFloat((price * count).toFixed(ERP_PRICE_DIGIT));
}
/**
* 【ERP】百分比计算，四舍五入保留两位小数
*
* 如果 total 为 0，则返回 0
*
* @param value 当前值
* @param total 总值
*/
function erpCalculatePercentage(value, total) {
	if (total === 0) return 0;
	return (value / total * 100).toFixed(2);
}
//#endregion
//#region src/utils/letter.ts
/**
* 将字符串的首字母大写
* @param string
*/
function capitalizeFirstLetter(string) {
	return string.charAt(0).toUpperCase() + string.slice(1);
}
/**
* 将字符串的首字母转换为小写。
*
* @param str 要转换的字符串
* @returns 首字母小写的字符串
*/
function toLowerCaseFirstLetter(str) {
	if (!str) return str;
	return str.charAt(0).toLowerCase() + str.slice(1);
}
/**
*  生成驼峰命名法的键名
* @param key
* @param parentKey
*/
function toCamelCase(key, parentKey) {
	if (!parentKey) return key;
	return parentKey + key.charAt(0).toUpperCase() + key.slice(1);
}
function kebabToCamelCase(str) {
	return str.split("-").filter(Boolean).map((word, index) => index === 0 ? word : word.charAt(0).toUpperCase() + word.slice(1)).join("");
}
//#endregion
//#region src/utils/merge.ts
const mergeWithArrayOverride = createDefu((originObj, key, updates) => {
	if (Array.isArray(originObj[key]) && Array.isArray(updates)) {
		originObj[key] = updates;
		return true;
	}
});
//#endregion
//#region src/utils/nprogress.ts
let nProgressInstance = null;
/**
* 动态加载NProgress库，并进行配置。
* 此函数首先检查是否已经加载过NProgress库，如果已经加载过，则直接返回NProgress实例。
* 否则，动态导入NProgress库，进行配置，然后返回NProgress实例。
*
* @returns  NProgress实例的Promise对象。
*/
async function loadNprogress() {
	if (nProgressInstance) return nProgressInstance;
	nProgressInstance = await import("nprogress");
	nProgressInstance.configure({
		showSpinner: true,
		speed: 300
	});
	return nProgressInstance;
}
/**
* 开始显示进度条。
* 此函数首先加载NProgress库，然后调用NProgress的start方法开始显示进度条。
*/
async function startProgress() {
	(await loadNprogress())?.start();
}
/**
* 停止显示进度条，并隐藏进度条。
* 此函数首先加载NProgress库，然后调用NProgress的done方法停止并隐藏进度条。
*/
async function stopProgress() {
	(await loadNprogress())?.done();
}
//#endregion
//#region src/utils/resources.ts
/**
* 加载js文件
* @param src js文件地址
*/
function loadScript(src) {
	return new Promise((resolve, reject) => {
		if (document.querySelector(`script[src="${src}"]`)) return resolve();
		const script = document.createElement("script");
		script.src = src;
		script.addEventListener("load", () => resolve());
		script.addEventListener("error", () => reject(/* @__PURE__ */ new Error(`Failed to load script: ${src}`)));
		document.head.append(script);
	});
}
//#endregion
//#region src/utils/stack.ts
/**
* @zh_CN 栈数据结构
*/
var Stack = class {
	/**
	* @zh_CN 栈内元素数量
	*/
	get size() {
		return this.items.length;
	}
	/**
	* @zh_CN 是否去重
	*/
	dedup;
	/**
	* @zh_CN 栈内元素
	*/
	items = [];
	/**
	* @zh_CN 栈的最大容量
	*/
	maxSize;
	constructor(dedup = true, maxSize) {
		this.maxSize = maxSize;
		this.dedup = dedup;
	}
	/**
	* @zh_CN 清空栈内元素
	*/
	clear() {
		this.items.length = 0;
	}
	/**
	* @zh_CN 查看栈顶元素
	* @returns 栈顶元素
	*/
	peek() {
		return this.items[this.items.length - 1];
	}
	/**
	* @zh_CN 出栈
	* @returns 栈顶元素
	*/
	pop() {
		return this.items.pop();
	}
	/**
	* @zh_CN 入栈
	* @param items 要入栈的元素
	*/
	push(...items) {
		items.forEach((item) => {
			if (this.dedup) {
				const index = this.items.indexOf(item);
				if (index !== -1) this.items.splice(index, 1);
			}
			this.items.push(item);
			if (this.maxSize && this.items.length > this.maxSize) this.items.splice(0, this.items.length - this.maxSize);
		});
	}
	/**
	* @zh_CN 移除栈内元素
	* @param itemList 要移除的元素列表
	*/
	remove(...itemList) {
		this.items = this.items.filter((i) => !itemList.includes(i));
	}
	/**
	* @zh_CN 保留栈内元素
	* @param itemList 要保留的元素列表
	*/
	retain(itemList) {
		this.items = this.items.filter((i) => itemList.includes(i));
	}
	/**
	* @zh_CN 转换为数组
	* @returns 栈内元素数组
	*/
	toArray() {
		return [...this.items];
	}
};
/**
* @zh_CN 创建一个栈实例
* @param dedup 是否去重
* @param maxSize 栈的最大容量
* @returns 栈实例
*/
const createStack = (dedup = true, maxSize) => new Stack(dedup, maxSize);
//#endregion
//#region src/utils/state-handler.ts
var StateHandler = class {
	condition = false;
	rejectCondition = null;
	resolveCondition = null;
	isConditionTrue() {
		return this.condition;
	}
	reset() {
		this.condition = false;
		this.clearPromises();
	}
	setConditionFalse() {
		this.condition = false;
		if (this.rejectCondition) {
			this.rejectCondition();
			this.clearPromises();
		}
	}
	setConditionTrue() {
		this.condition = true;
		if (this.resolveCondition) {
			this.resolveCondition();
			this.clearPromises();
		}
	}
	waitForCondition() {
		return new Promise((resolve, reject) => {
			if (this.condition) resolve();
			else {
				this.resolveCondition = resolve;
				this.rejectCondition = reject;
			}
		});
	}
	clearPromises() {
		this.resolveCondition = null;
		this.rejectCondition = null;
	}
};
//#endregion
//#region src/utils/time.ts
/**
* @param {Date | number | string} time 需要转换的时间
* @param {string} fmt 需要转换的格式 如 yyyy-MM-dd、yyyy-MM-dd HH:mm:ss
*/
function formatTime(time, fmt) {
	if (time) {
		const date = new Date(time);
		const o = {
			"M+": date.getMonth() + 1,
			"d+": date.getDate(),
			"H+": date.getHours(),
			"m+": date.getMinutes(),
			"s+": date.getSeconds(),
			"q+": Math.floor((date.getMonth() + 3) / 3),
			S: date.getMilliseconds()
		};
		const yearMatch = fmt.match(/y+/);
		if (yearMatch) fmt = fmt.replace(yearMatch[0], `${date.getFullYear()}`.slice(4 - yearMatch[0].length));
		for (const k in o) {
			const match = fmt.match(new RegExp(`(${k})`));
			if (match) fmt = fmt.replace(match[0], match[0].length === 1 ? o[k] : `00${o[k]}`.slice(`${o[k]}`.length));
		}
		return fmt;
	} else return "";
}
/**
* 获取当前日期是第几周
* @param dateTime 当前传入的日期值
* @returns 返回第几周数字值
*/
function getWeek(dateTime) {
	const temptTime = new Date(dateTime);
	const weekday = temptTime.getDay() || 7;
	temptTime.setDate(temptTime.getDate() - weekday + 1 + 5);
	let firstDay = new Date(temptTime.getFullYear(), 0, 1);
	const dayOfWeek = firstDay.getDay();
	let spendDay = 1;
	if (dayOfWeek !== 0) spendDay = 7 - dayOfWeek + 1;
	firstDay = new Date(temptTime.getFullYear(), 0, 1 + spendDay);
	const d = Math.ceil((temptTime.valueOf() - firstDay.valueOf()) / 864e5);
	return Math.ceil(d / 7);
}
/**
* 将时间转换为 `几秒前`、`几分钟前`、`几小时前`、`几天前`
* @param param 当前时间，new Date() 格式或者字符串时间格式
* @param format 需要转换的时间格式字符串
* @description param 10秒：  10 * 1000
* @description param 1分：   60 * 1000
* @description param 1小时： 60 * 60 * 1000
* @description param 24小时：60 * 60 * 24 * 1000
* @description param 3天：   60 * 60* 24 * 1000 * 3
* @returns 返回拼接后的时间字符串
*/
function formatPast(param, format = "YYYY-MM-DD HH:mm:ss") {
	let s, t;
	let time = Date.now();
	typeof param === "string" || typeof param === "object" ? t = new Date(param).getTime() : t = param;
	time = Number.parseInt(`${time - t}`);
	if (time < 1e4) return "刚刚";
	else if (time < 6e4 && time >= 1e4) {
		s = Math.floor(time / 1e3);
		return `${s}秒前`;
	} else if (time < 36e5 && time >= 6e4) {
		s = Math.floor(time / 6e4);
		return `${s}分钟前`;
	} else if (time < 864e5 && time >= 36e5) {
		s = Math.floor(time / 36e5);
		return `${s}小时前`;
	} else if (time < 2592e5 && time >= 864e5) {
		s = Math.floor(time / 864e5);
		return `${s}天前`;
	} else return formatDate(typeof param === "string" || typeof param === "object" ? new Date(param) : param, format);
}
/**
* 时间问候语
* @param param 当前时间，new Date() 格式
* @description param 调用 `formatAxis(new Date())` 输出 `上午好`
* @returns 返回拼接后的时间字符串
*/
function formatAxis(param) {
	const hour = new Date(param).getHours();
	if (hour < 6) return "凌晨好";
	else if (hour < 9) return "早上好";
	else if (hour < 12) return "上午好";
	else if (hour < 14) return "中午好";
	else if (hour < 17) return "下午好";
	else if (hour < 19) return "傍晚好";
	else if (hour < 22) return "晚上好";
	else return "夜里好";
}
/**
* 将毫秒，转换成时间字符串。例如说，xx 分钟
*
* @param ms 毫秒
* @returns {string} 字符串
*/
function formatPast2(ms) {
	const day = Math.floor(ms / (1440 * 60 * 1e3));
	const hour = Math.floor(ms / (3600 * 1e3) - day * 24);
	const minute = Math.floor(ms / (60 * 1e3) - day * 24 * 60 - hour * 60);
	const second = Math.floor(ms / 1e3 - day * 24 * 60 * 60 - hour * 60 * 60 - minute * 60);
	if (day > 0) return `${day} 天${hour} 小时 ${minute} 分钟`;
	if (hour > 0) return `${hour} 小时 ${minute} 分钟`;
	if (minute > 0) return `${minute} 分钟`;
	return second > 0 ? `${second} 秒` : `0 秒`;
}
/**
* 设置起始日期，时间为00:00:00
* @param param 传入日期
* @returns 带时间00:00:00的日期
*/
function beginOfDay(param) {
	return new Date(param.getFullYear(), param.getMonth(), param.getDate(), 0, 0, 0);
}
/**
* 设置结束日期，时间为23:59:59
* @param param 传入日期
* @returns 带时间23:59:59的日期
*/
function endOfDay(param) {
	return new Date(param.getFullYear(), param.getMonth(), param.getDate(), 23, 59, 59);
}
/**
* 计算两个日期间隔天数
* @param param1 日期1
* @param param2 日期2
*/
function betweenDay(param1, param2) {
	param1 = convertDate(param1);
	param2 = convertDate(param2);
	return Math.floor((param2.getTime() - param1.getTime()) / (24 * 3600 * 1e3));
}
/**
* 日期计算
* @param param1 日期
* @param param2 添加的时间
*/
function addTime(param1, param2) {
	param1 = convertDate(param1);
	return new Date(param1.getTime() + param2);
}
/**
* 日期转换
* @param param 日期
*/
function convertDate(param) {
	if (typeof param === "string") return new Date(param);
	return param;
}
/**
* 指定的两个日期, 是否为同一天
* @param a 日期 A
* @param b 日期 B
*/
function isSameDay(a, b) {
	if (!a || !b) return false;
	const aa = dayjs(a);
	const bb = dayjs(b);
	return aa.year() === bb.year() && aa.month() === bb.month() && aa.day() === bb.day();
}
/**
* 获取一天的开始时间、截止时间
* @param date 日期
* @param days 天数
*/
function getDayRange(date, days) {
	const day = dayjs(date).add(days, "d");
	return getDateRange(day, day);
}
/**
* 获取最近7天的开始时间、截止时间
*/
function getLast7Days() {
	return getDateRange(dayjs().subtract(7, "d"), dayjs().subtract(1, "d"));
}
/**
* 获取最近30天的开始时间、截止时间
*/
function getLast30Days() {
	return getDateRange(dayjs().subtract(30, "d"), dayjs().subtract(1, "d"));
}
/**
* 获取最近1年的开始时间、截止时间
*/
function getLast1Year() {
	return getDateRange(dayjs().subtract(1, "y"), dayjs().subtract(1, "d"));
}
/**
* 获取指定日期的开始时间、截止时间
* @param beginDate 开始日期
* @param endDate 截止日期
*/
function getDateRange(beginDate, endDate) {
	return [dayjs(beginDate).startOf("d").format("YYYY-MM-DD HH:mm:ss"), dayjs(endDate).endOf("d").format("YYYY-MM-DD HH:mm:ss")];
}
//#endregion
//#region src/utils/to.ts
/**
* @param { Readonly<Promise> } promise
* @param {object=} errorExt - Additional Information you can pass to the err object
* @return { Promise }
*/
async function to(promise, errorExt) {
	try {
		return [null, await promise];
	} catch (error) {
		if (errorExt) return [Object.assign({}, error, errorExt), void 0];
		return [error, void 0];
	}
}
//#endregion
//#region src/utils/tree.ts
/**
* @zh_CN 遍历树形结构，并返回所有节点中指定的值。
* @param tree 树形结构数组
* @param getValue 获取节点值的函数
* @param options 作为子节点数组的可选属性名称。
* @returns 所有节点中指定的值的数组
*/
function traverseTreeValues(tree, getValue, options) {
	const result = [];
	const { childProps } = options || { childProps: "children" };
	const dfs = (treeNode) => {
		const value = getValue(treeNode);
		result.push(value);
		const children = treeNode?.[childProps];
		if (!children) return;
		if (children.length > 0) for (const child of children) dfs(child);
	};
	for (const treeNode of tree) dfs(treeNode);
	return result.filter(Boolean);
}
/**
* 根据条件过滤给定树结构的节点，并以原有顺序返回所有匹配节点的数组。
* @param tree 要过滤的树结构的根节点数组。
* @param filter 用于匹配每个节点的条件。
* @param options 作为子节点数组的可选属性名称。
* @returns 包含所有匹配节点的数组。
*/
function filterTree(tree, filter, options) {
	const { childProps } = options || { childProps: "children" };
	const _filterTree = (nodes) => {
		return nodes.filter((node) => {
			if (filter(node)) {
				if (node[childProps]) node[childProps] = _filterTree(node[childProps]);
				return true;
			}
			return false;
		});
	};
	return _filterTree(tree);
}
/**
* 根据条件重新映射给定树结构的节
* @param tree 要过滤的树结构的根节点数组。
* @param mapper 用于map每个节点的条件。
* @param options 作为子节点数组的可选属性名称。
*/
function mapTree(tree, mapper, options) {
	const { childProps } = options || { childProps: "children" };
	return tree.map((node) => {
		const mapperNode = mapper(node);
		if (mapperNode[childProps]) mapperNode[childProps] = mapTree(mapperNode[childProps], mapper, options);
		return mapperNode;
	});
}
/**
* 构造树型结构数据
*
* @param {*} data 数据源
* @param {*} id id字段 默认 'id'
* @param {*} parentId 父节点字段 默认 'parentId'
* @param {*} children 孩子节点字段 默认 'children'
*/
function handleTree(data, id = "id", parentId = "parentId", children = "children") {
	if (!Array.isArray(data)) {
		console.warn("data must be an array");
		return [];
	}
	const config = {
		id,
		parentId,
		childrenList: children
	};
	const childrenListMap = {};
	const nodeIds = {};
	const tree = [];
	for (const d of data) {
		const pId = d[config.parentId];
		if (childrenListMap[pId] === void 0) childrenListMap[pId] = [];
		nodeIds[d[config.id]] = d;
		childrenListMap[pId].push(d);
	}
	for (const d of data) if (nodeIds[d[config.parentId]] === void 0) tree.push(d);
	const adaptToChildrenList = (node) => {
		const nodeId = node[config.id];
		if (childrenListMap[nodeId]) {
			node[config.childrenList] = childrenListMap[nodeId];
			for (const child of node[config.childrenList]) adaptToChildrenList(child);
		}
	};
	for (const rootNode of tree) adaptToChildrenList(rootNode);
	return tree;
}
/**
* 获取节点的完整结构
* @param tree 树数据
* @param nodeId 节点 id
*/
function treeToString(tree, nodeId) {
	if (tree === void 0 || !Array.isArray(tree) || tree.length === 0) {
		console.warn("tree must be an array");
		return "";
	}
	const node = tree.find((item) => item.id === nodeId);
	if (node !== void 0) return node.name;
	let str = "";
	function performAThoroughValidation(arr) {
		if (arr === void 0 || !Array.isArray(arr) || arr.length === 0) return false;
		for (const item of arr) if (item.id === nodeId) {
			str += ` / ${item.name}`;
			return true;
		} else if (item.children !== void 0 && item.children.length > 0) {
			str += ` / ${item.name}`;
			if (performAThoroughValidation(item.children)) return true;
		}
		return false;
	}
	for (const item of tree) {
		str = `${item.name}`;
		if (performAThoroughValidation(item.children)) break;
	}
	return str;
}
/**
* 对树形结构数据进行递归排序
* @param treeData - 树形数据数组
* @param sortFunction - 排序函数，用于定义排序规则
* @param options - 配置选项，包括子节点属性名
* @returns 排序后的树形数据
*/
function sortTree(treeData, sortFunction, options) {
	const { childProps } = options || { childProps: "children" };
	return treeData.toSorted(sortFunction).map((item) => {
		const children = item[childProps];
		if (children && Array.isArray(children) && children.length > 0) return {
			...item,
			[childProps]: sortTree(children, sortFunction, options)
		};
		return item;
	});
}
//#endregion
//#region src/utils/unique.ts
/**
* 根据指定字段对对象数组进行去重
* @param arr 要去重的对象数组
* @param key 去重依据的字段名
* @returns 去重后的对象数组
*/
function uniqueByField(arr, key) {
	const seen = /* @__PURE__ */ new Map();
	return arr.filter((item) => {
		const value = item[key];
		return seen.has(value) ? false : (seen.set(value, item), true);
	});
}
//#endregion
//#region src/utils/update-css-variables.ts
/**
* 更新 CSS 变量的函数
* @param variables 要更新的 CSS 变量与其新值的映射
*/
function updateCSSVariables(variables, id = "__vben-styles__") {
	const styleElement = document.querySelector(`#${id}`) || document.createElement("style");
	styleElement.id = id;
	let cssText = ":root {";
	for (const key in variables) if (Object.prototype.hasOwnProperty.call(variables, key)) cssText += `${key}: ${variables[key]};`;
	cssText += "}";
	styleElement.textContent = cssText;
	if (!document.querySelector(`#${id}`)) setTimeout(() => {
		document.head.append(styleElement);
	});
}
//#endregion
//#region src/utils/upload.ts
/**
* 根据支持的文件类型生成 accept 属性值
*
* @param supportedFileTypes 支持的文件类型数组，如 ['PDF', 'DOC', 'DOCX']
* @returns 用于文件上传组件 accept 属性的字符串
*/
function generateAcceptedFileTypes(supportedFileTypes) {
	const allowedExtensions = supportedFileTypes.map((ext) => ext.toLowerCase());
	const mimeTypes = [];
	if (allowedExtensions.includes("txt")) mimeTypes.push("text/plain");
	if (allowedExtensions.includes("pdf")) mimeTypes.push("application/pdf");
	if (allowedExtensions.includes("html") || allowedExtensions.includes("htm")) mimeTypes.push("text/html");
	if (allowedExtensions.includes("csv")) mimeTypes.push("text/csv");
	if (allowedExtensions.includes("xlsx") || allowedExtensions.includes("xls")) mimeTypes.push("application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet");
	if (allowedExtensions.includes("docx") || allowedExtensions.includes("doc")) mimeTypes.push("application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document");
	if (allowedExtensions.includes("pptx") || allowedExtensions.includes("ppt")) mimeTypes.push("application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation");
	if (allowedExtensions.includes("xml")) mimeTypes.push("application/xml", "text/xml");
	if (allowedExtensions.includes("md") || allowedExtensions.includes("markdown")) mimeTypes.push("text/markdown");
	if (allowedExtensions.includes("epub")) mimeTypes.push("application/epub+zip");
	if (allowedExtensions.includes("eml")) mimeTypes.push("message/rfc822");
	if (allowedExtensions.includes("msg")) mimeTypes.push("application/vnd.ms-outlook");
	const extensions = allowedExtensions.map((ext) => `.${ext}`);
	return [...mimeTypes, ...extensions].join(",");
}
/**
* 从 URL 中提取文件名
*
* @param url 文件 URL
* @returns 文件名，如果无法提取则返回 'unknown'
*/
function getFileNameFromUrl(url) {
	if (!url) return "unknown";
	try {
		const fileName = new URL(url).pathname.split("/").pop() || "unknown";
		return decodeURIComponent(fileName);
	} catch {
		const parts = url.split("/");
		return parts[parts.length - 1] || "unknown";
	}
}
/**
* 默认图片类型
*/
const defaultImageAccepts = [
	"bmp",
	"gif",
	"jpeg",
	"jpg",
	"png",
	"svg",
	"webp"
];
/**
* 图片类 MIME 子类型到扩展名的别名映射；未列出的子类型按字面量与扩展名比较
*/
const IMAGE_MIME_SUBTYPE_ALIASES = {
	apng: ["apng", "png"],
	jpeg: ["jpeg", "jpg"],
	pjpeg: ["jpeg", "jpg"],
	"svg+xml": ["svg"],
	tiff: ["tif", "tiff"],
	"x-icon": ["ico"]
};
/**
* 判断 MIME 子类型是否与文件扩展名匹配；image/* 限定为已知图片扩展名集合
*/
function matchMimeSubtype(subtype, ext) {
	if (subtype === "*") return defaultImageAccepts.includes(ext);
	const aliases = IMAGE_MIME_SUBTYPE_ALIASES[subtype];
	if (aliases) return aliases.includes(ext);
	return subtype === ext;
}
/**
* 判断文件是否为图片
*
* @param filename 文件名
* @param accepts 支持的文件类型，兼容 MIME（如 image/png）、.ext（如 .png）与纯后缀（如 png）
* @returns 是否为图片
*/
function isImage(filename, accepts = defaultImageAccepts) {
	if (!filename || accepts.length === 0) return false;
	const ext = filename.split(".").pop()?.toLowerCase() || "";
	if (!ext) return false;
	return accepts.some((accept) => {
		const lower = accept.toLowerCase();
		if (lower.includes("/")) return matchMimeSubtype(lower.split("/").pop() || "", ext);
		if (lower.startsWith(".")) return lower.slice(1) === ext;
		return lower === ext;
	});
}
/**
* 判断文件是否为指定类型
*
* @param file 文件
* @param accepts 支持的文件类型
* @returns 是否为指定类型
*/
function checkFileType(file, accepts) {
	if (!accepts || accepts.length === 0) return true;
	const newTypes = accepts.join("|");
	return new RegExp(`${String.raw`\.(` + newTypes})$`, "i").test(file.name);
}
/**
* 格式化文件大小
*
* @param bytes 文件大小（字节）
* @returns 格式化后的文件大小字符串
*/
function formatFileSize(bytes, digits = 2) {
	if (bytes === 0) return "0 B";
	const k = 1024;
	const unitArr = [
		"B",
		"KB",
		"MB",
		"GB",
		"TB",
		"PB",
		"EB",
		"ZB",
		"YB"
	];
	const index = Math.floor(Math.log(bytes) / Math.log(k));
	return `${Number.parseFloat((bytes / k ** index).toFixed(digits))} ${unitArr[index]}`;
}
/**
* 获取文件图标（Lucide Icons）
*
* @param filename 文件名
* @returns Lucide 图标名称
*/
function getFileIcon(filename) {
	if (!filename) return "lucide:file";
	const ext = filename.split(".").pop()?.toLowerCase() || "";
	if (isImage(ext)) return "lucide:image";
	if (["pdf"].includes(ext)) return "lucide:file-text";
	if (["doc", "docx"].includes(ext)) return "lucide:file-text";
	if (["xls", "xlsx"].includes(ext)) return "lucide:file-spreadsheet";
	if (["ppt", "pptx"].includes(ext)) return "lucide:presentation";
	if ([
		"aac",
		"m4a",
		"mp3",
		"wav"
	].includes(ext)) return "lucide:music";
	if ([
		"avi",
		"mov",
		"mp4",
		"wmv"
	].includes(ext)) return "lucide:video";
	return "lucide:file";
}
/**
* 获取文件类型样式类（Tailwind CSS 渐变色）
*
* @param filename 文件名
* @returns Tailwind CSS 渐变类名
*/
function getFileTypeClass(filename) {
	if (!filename) return "from-gray-500 to-gray-700";
	const ext = filename.split(".").pop()?.toLowerCase() || "";
	if (isImage(ext)) return "from-yellow-400 to-orange-500";
	if (["pdf"].includes(ext)) return "from-red-500 to-red-700";
	if (["doc", "docx"].includes(ext)) return "from-blue-600 to-blue-800";
	if (["xls", "xlsx"].includes(ext)) return "from-green-600 to-green-800";
	if (["ppt", "pptx"].includes(ext)) return "from-orange-600 to-orange-800";
	if ([
		"aac",
		"m4a",
		"mp3",
		"wav"
	].includes(ext)) return "from-purple-500 to-purple-700";
	if ([
		"avi",
		"mov",
		"mp4",
		"wmv"
	].includes(ext)) return "from-red-500 to-red-700";
	return "from-gray-500 to-gray-700";
}
//#endregion
//#region src/utils/url.ts
/**
* URL 验证
* @param path URL 路径
*/
function isUrl(path) {
	try {
		return Boolean(new URL(path));
	} catch {
		return false;
	}
}
//#endregion
//#region src/utils/util.ts
function bindMethods(instance) {
	const prototype = Object.getPrototypeOf(instance);
	Object.getOwnPropertyNames(prototype).forEach((propertyName) => {
		const descriptor = Object.getOwnPropertyDescriptor(prototype, propertyName);
		const propertyValue = instance[propertyName];
		if (typeof propertyValue === "function" && propertyName !== "constructor" && descriptor && !descriptor.get && !descriptor.set) instance[propertyName] = propertyValue.bind(instance);
	});
}
/**
* 获取嵌套对象的字段值
* @param obj - 要查找的对象
* @param path - 用于查找字段的路径，使用小数点分隔
* @returns 字段值，或者未找到时返回 undefined
*/
function getNestedValue(obj, path) {
	if (typeof path !== "string" || path.length === 0) throw new Error("Path must be a non-empty string");
	const keys = path.split(".");
	let current = obj;
	for (const key of keys) {
		if (current === null || current === void 0) return;
		current = current[key];
	}
	return current;
}
/**
* 获取链接的参数值（值类型）
* @param key 参数键名
* @param urlStr 链接地址，默认为当前浏览器的地址
*/
function getUrlNumberValue(key, urlStr = location.href) {
	return Number(getUrlValue(key, urlStr));
}
/**
* 获取链接的参数值
* @param key 参数键名
* @param urlStr 链接地址，默认为当前浏览器的地址
*/
function getUrlValue(key, urlStr = location.href) {
	if (!urlStr || !key) return "";
	return new URL(decodeURIComponent(urlStr)).searchParams.get(key) ?? "";
}
/**
* 将值复制到目标对象，且以目标对象属性为准，例：target: {a:1} source:{a:2,b:3} 结果为：{a:2}
* @param target 目标对象
* @param source 源对象
*/
function copyValueToTarget(target, source) {
	const newObj = Object.assign({}, target, source);
	Object.keys(newObj).forEach((key) => {
		if (!Object.keys(target).includes(key)) delete newObj[key];
	});
	Object.assign(target, newObj);
}
/** 实现 groupBy 功能 */
function groupBy(array, key) {
	const result = {};
	for (const item of array) {
		const groupKey = item[key];
		if (!result[groupKey]) result[groupKey] = [];
		result[groupKey].push(item);
	}
	return result;
}
/**
* 解析 JSON 字符串
*
* @param str
*/
function jsonParse(str) {
	try {
		return JSON.parse(str);
	} catch {
		console.warn(`str[${str}] 不是一个 JSON 字符串`);
		return str;
	}
}
//#endregion
//#region src/utils/uuid.ts
const hexList = [];
for (let i = 0; i <= 15; i++) hexList[i] = i.toString(16);
function buildUUID() {
	let uuid = "";
	for (let i = 1; i <= 36; i++) switch (i) {
		case 9:
		case 14:
		case 19:
		case 24:
			uuid += "-";
			break;
		case 15:
			uuid += 4;
			break;
		case 20:
			uuid += hexList[Math.random() * 4 | 8];
			break;
		default: uuid += hexList[Math.trunc(Math.random() * 16)];
	}
	return uuid.replaceAll("-", "");
}
let unique = 0;
function buildShortUUID(prefix = "") {
	const time = Date.now();
	const random = Math.floor(Math.random() * 1e9);
	unique++;
	return `${prefix}_${random}${unique}${String(time)}`;
}
//#endregion
export { AES, ApiEncrypt, RSA, Stack, StateHandler, addTime, arraysEqual, base64ToFile, beginOfDay, betweenDay, bindMethods, buildShortUUID, buildUUID, calculateRelativeRate, capitalizeFirstLetter, checkFileType, cloneDeep, cn, convertDate, convertToInteger, copyValueToTarget, createApiEncrypt, createMerge, createStack, dataURLtoBlob, dateFormatter, defaultImageAccepts, diff, downloadFileFromBase64, downloadFileFromBlob, downloadFileFromBlobPart, downloadFileFromImageUrl, downloadFileFromUrl, downloadImageByCanvas, endOfDay, erpCalculatePercentage, erpCountInputFormatter, erpCountTableColumnFormatter, erpNumberFormatter, erpPriceInputFormatter, erpPriceMultiply, erpPriceTableColumnFormatter, fenToYuan, fenToYuanFormat, filterTree, floatToFixed2, formatAxis, formatDate, formatDate2, formatDateTime, formatFileSize, formatPast, formatPast2, formatTime, formatToFraction, formatToFractionDigit, generateAcceptedFileTypes, get, getCurrentTimezone, getDateRange, getDayRange, getElementVisibleRect, getFileIcon, getFileNameFromUrl, getFileTypeClass, getFirstNonNullOrUndefined, getLast1Year, getLast30Days, getLast7Days, getNestedValue, getScrollbarWidth, getSystemTimezone, getUrlNumberValue, getUrlValue, getWeek, groupBy, handleTree, isBoolean, isDate, isDayjsObject, isEmpty, isEqual, isFunction, isHttpUrl, isImage, isMacOs, isNumber, isObject, isSameDay, isString, isUndefined, isUrl, isWindow, isWindowsOs, jsonParse, kebabToCamelCase, loadScript, mapTree, merge, mergeWithArrayOverride, needsScrollbar, openRouteInNewWindow, openWindow, set, setCurrentTimezone, sortTree, startProgress, stopProgress, to, toCamelCase, toLowerCaseFirstLetter, traverseTreeValues, treeToString, triggerDownload, triggerWindowResize, uniqueByField, updateCSSVariables, urlToBase64, yuanToFen };
