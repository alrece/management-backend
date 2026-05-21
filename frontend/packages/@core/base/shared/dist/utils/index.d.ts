import { ClassValue } from "clsx";
import dayjs from "dayjs";
import { isFunction, isObject, isString } from "@vue/shared";
import * as _$defu from "defu";
import { createDefu as createMerge, defu as merge } from "defu";
import { get, isEqual, set } from "es-toolkit/compat";
import cloneDeep from "lodash.clonedeep";

//#region src/utils/cn.d.ts
declare function cn(...inputs: ClassValue[]): string;
//#endregion
//#region src/utils/date.d.ts
type FormatDate = Date | dayjs.Dayjs | number | string;
type Format = 'HH' | 'HH:mm' | 'HH:mm:ss' | 'YYYY' | 'YYYY-MM' | 'YYYY-MM-DD' | 'YYYY-MM-DD HH' | 'YYYY-MM-DD HH:mm' | 'YYYY-MM-DD HH:mm:ss' | (string & {});
declare function formatDate(time?: FormatDate, format?: Format): string;
declare function formatDateTime(time?: FormatDate): string;
declare function formatDate2(date: Date, format?: string): string;
declare function isDate(value: any): value is Date;
declare function isDayjsObject(value: any): value is dayjs.Dayjs;
declare function dateFormatter(_row: any, _column: any, cellValue: any): string;
declare const getSystemTimezone: () => string;
declare const setCurrentTimezone: (timezone?: string) => void;
declare const getCurrentTimezone: () => string;
//#endregion
//#region src/utils/diff.d.ts
declare function arraysEqual<T>(a: T[], b: T[]): boolean;
type DiffResult<T> = Partial<{ [K in keyof T]: T[K] extends object ? DiffResult<T[K]> : T[K] }>;
declare function diff<T extends Record<string, any>>(obj1: T, obj2: T): DiffResult<T>;
//#endregion
//#region src/utils/dom.d.ts
interface VisibleDomRect {
  bottom: number;
  height: number;
  left: number;
  right: number;
  top: number;
  width: number;
}
declare function getElementVisibleRect(element?: HTMLElement | null | undefined): VisibleDomRect;
declare function getScrollbarWidth(): number;
declare function needsScrollbar(): boolean;
declare function triggerWindowResize(): void;
//#endregion
//#region src/utils/download.d.ts
interface DownloadOptions<T = string> {
  fileName?: string;
  source: T;
  target?: string;
}
declare function downloadFileFromUrl({
  fileName,
  source,
  target
}: DownloadOptions): Promise<void>;
declare function downloadImageByCanvas({
  url,
  canvasWidth,
  canvasHeight,
  drawWithImageSize
}: {
  canvasHeight?: number;
  canvasWidth?: number;
  drawWithImageSize?: boolean;
  url: string;
}): void;
declare function downloadFileFromBase64({
  fileName,
  source
}: DownloadOptions): void;
declare function downloadFileFromImageUrl({
  fileName,
  source
}: DownloadOptions): Promise<void>;
declare function downloadFileFromBlob({
  fileName,
  source
}: DownloadOptions<Blob>): void;
declare function downloadFileFromBlobPart({
  fileName,
  source
}: DownloadOptions<BlobPart>): void;
declare function dataURLtoBlob(base64Buf: string): Blob;
declare function urlToBase64(url: string, mineType?: string): Promise<string>;
declare function base64ToFile(base64: string, fileName: string): File;
declare function triggerDownload(href: string, fileName: string | undefined, revokeDelay?: number): void;
//#endregion
//#region src/utils/encrypt.d.ts
declare const AES: {
  encrypt(data: string, key: string): string;
  decrypt(encryptedData: string, key: string): string;
};
declare const RSA: {
  encrypt(data: string, publicKey: string): false | string;
  decrypt(encryptedData: string, privateKey: string): false | string;
};
interface ApiEncryptConfig {
  algorithm: 'AES' | 'RSA';
  enable: boolean;
  header: string;
  requestKey: string;
  responseKey: string;
}
declare class ApiEncrypt {
  private config;
  constructor(config: ApiEncryptConfig);
  decryptResponse(encryptedData: string): any;
  encryptRequest(data: any): string;
  getEncryptHeader(): string;
}
declare function createApiEncrypt(env: Record<string, any>): ApiEncrypt;
//#endregion
//#region src/utils/formatNumber.d.ts
declare function formatToFractionDigit(num: number | string | undefined, digit?: number): string;
declare function formatToFraction(num: number | string | undefined): string;
declare function floatToFixed2(num: number | string | undefined): string;
declare function convertToInteger(num: number | string | undefined): number;
declare function yuanToFen(amount: number | string): number;
declare function fenToYuan(price: number | string): string;
declare const fenToYuanFormat: (_: any, __: any, cellValue: any, ___: any) => string;
declare function calculateRelativeRate(value?: number, reference?: number): number;
declare function erpNumberFormatter(num: number | string | undefined, digit: number): string;
declare function erpCountInputFormatter(num: number | string | undefined): string;
declare function erpCountTableColumnFormatter(cellValue: any): string;
declare function erpPriceInputFormatter(num: number | string | undefined): string;
declare function erpPriceTableColumnFormatter(cellValue: any): string;
declare function erpPriceMultiply(price: number, count: number): number | undefined;
declare function erpCalculatePercentage(value: number, total: number): string | 0;
//#endregion
//#region src/utils/inference.d.ts
declare function isUndefined(value?: unknown): value is undefined;
declare function isBoolean(value: unknown): value is boolean;
declare function isEmpty<T = unknown>(value?: T): value is T;
declare function isHttpUrl(url?: string): boolean;
declare function isWindow(value: any): value is Window;
declare function isMacOs(): boolean;
declare function isWindowsOs(): boolean;
declare function isNumber(value: any): value is number;
declare function getFirstNonNullOrUndefined<T>(...values: (null | T | undefined)[]): T | undefined;
//#endregion
//#region src/utils/letter.d.ts
declare function capitalizeFirstLetter(string: string): string;
declare function toLowerCaseFirstLetter(str: string): string;
declare function toCamelCase(key: string, parentKey: string): string;
declare function kebabToCamelCase(str: string): string;
//#endregion
//#region src/utils/merge.d.ts
declare const mergeWithArrayOverride: _$defu.DefuFn;
//#endregion
//#region src/utils/nprogress.d.ts
declare function startProgress(): Promise<void>;
declare function stopProgress(): Promise<void>;
//#endregion
//#region src/utils/resources.d.ts
declare function loadScript(src: string): Promise<void>;
//#endregion
//#region src/utils/stack.d.ts
declare class Stack<T> {
  get size(): number;
  private readonly dedup;
  private items;
  private readonly maxSize?;
  constructor(dedup?: boolean, maxSize?: number);
  clear(): void;
  peek(): T | undefined;
  pop(): T | undefined;
  push(...items: T[]): void;
  remove(...itemList: T[]): void;
  retain(itemList: T[]): void;
  toArray(): T[];
}
declare const createStack: <T>(dedup?: boolean, maxSize?: number) => Stack<T>;
//#endregion
//#region src/utils/state-handler.d.ts
declare class StateHandler {
  private condition;
  private rejectCondition;
  private resolveCondition;
  isConditionTrue(): boolean;
  reset(): void;
  setConditionFalse(): void;
  setConditionTrue(): void;
  waitForCondition(): Promise<void>;
  private clearPromises;
}
//#endregion
//#region src/utils/time.d.ts
declare function formatTime(time: Date | number | string, fmt: string): string;
declare function getWeek(dateTime: Date): number;
declare function formatPast(param: Date | string, format?: string): string;
declare function formatAxis(param: Date): string;
declare function formatPast2(ms: number): string;
declare function beginOfDay(param: Date): Date;
declare function endOfDay(param: Date): Date;
declare function betweenDay(param1: Date, param2: Date): number;
declare function addTime(param1: Date, param2: number): Date;
declare function convertDate(param: Date | string): Date;
declare function isSameDay(a: dayjs.ConfigType, b: dayjs.ConfigType): boolean;
declare function getDayRange(date: dayjs.ConfigType, days: number): [dayjs.ConfigType, dayjs.ConfigType];
declare function getLast7Days(): [dayjs.ConfigType, dayjs.ConfigType];
declare function getLast30Days(): [dayjs.ConfigType, dayjs.ConfigType];
declare function getLast1Year(): [dayjs.ConfigType, dayjs.ConfigType];
declare function getDateRange(beginDate: dayjs.ConfigType, endDate: dayjs.ConfigType): [string, string];
//#endregion
//#region src/utils/to.d.ts
declare function to<T, U = Error>(promise: Readonly<Promise<T>>, errorExt?: object): Promise<[null, T] | [U, undefined]>;
//#endregion
//#region src/utils/tree.d.ts
interface TreeConfigOptions {
  childProps: string;
}
interface TreeNode {
  [key: string]: any;
  children?: TreeNode[];
}
declare function traverseTreeValues<T, V>(tree: T[], getValue: (node: T) => V, options?: TreeConfigOptions): V[];
declare function filterTree<T extends Record<string, any>>(tree: T[], filter: (node: T) => boolean, options?: TreeConfigOptions): T[];
declare function mapTree<T, V extends Record<string, any>>(tree: T[], mapper: (node: T) => V, options?: TreeConfigOptions): V[];
declare function handleTree(data: TreeNode[], id?: string, parentId?: string, children?: string): TreeNode[];
declare function treeToString(tree: any[], nodeId: number | string): any;
declare function sortTree<T extends Record<string, any>>(treeData: T[], sortFunction: (a: T, b: T) => number, options?: TreeConfigOptions): T[];
//#endregion
//#region src/utils/unique.d.ts
declare function uniqueByField<T>(arr: T[], key: keyof T): T[];
//#endregion
//#region src/utils/update-css-variables.d.ts
declare function updateCSSVariables(variables: {
  [key: string]: string;
}, id?: string): void;
//#endregion
//#region src/utils/upload.d.ts
declare function generateAcceptedFileTypes(supportedFileTypes: string[]): string;
declare function getFileNameFromUrl(url: null | string | undefined): string;
declare const defaultImageAccepts: string[];
declare function isImage(filename: null | string | undefined, accepts?: string[]): boolean;
declare function checkFileType(file: File, accepts: string[]): boolean;
declare function formatFileSize(bytes: number, digits?: number): string;
declare function getFileIcon(filename: null | string | undefined): string;
declare function getFileTypeClass(filename: null | string | undefined): string;
//#endregion
//#region src/utils/url.d.ts
declare function isUrl(path: string): boolean;
//#endregion
//#region src/utils/util.d.ts
declare function bindMethods<T extends object>(instance: T): void;
declare function getNestedValue<T>(obj: T, path: string): any;
declare function getUrlNumberValue(key: string, urlStr?: string): number;
declare function getUrlValue(key: string, urlStr?: string): string;
declare function copyValueToTarget(target: any, source: any): void;
declare function groupBy(array: any[], key: string): Record<string, any[]>;
declare function jsonParse(str: string): any;
//#endregion
//#region src/utils/uuid.d.ts
declare function buildUUID(): string;
declare function buildShortUUID(prefix?: string): string;
//#endregion
//#region src/utils/window.d.ts
interface OpenWindowOptions {
  noopener?: boolean;
  noreferrer?: boolean;
  target?: '_blank' | '_parent' | '_self' | '_top' | string;
}
declare function openWindow(url: string, options?: OpenWindowOptions): void;
declare function openRouteInNewWindow(path: string): void;
//#endregion
export { AES, ApiEncrypt, ApiEncryptConfig, RSA, Stack, StateHandler, VisibleDomRect, addTime, arraysEqual, base64ToFile, beginOfDay, betweenDay, bindMethods, buildShortUUID, buildUUID, calculateRelativeRate, capitalizeFirstLetter, checkFileType, cloneDeep, cn, convertDate, convertToInteger, copyValueToTarget, createApiEncrypt, createMerge, createStack, dataURLtoBlob, dateFormatter, defaultImageAccepts, diff, downloadFileFromBase64, downloadFileFromBlob, downloadFileFromBlobPart, downloadFileFromImageUrl, downloadFileFromUrl, downloadImageByCanvas, endOfDay, erpCalculatePercentage, erpCountInputFormatter, erpCountTableColumnFormatter, erpNumberFormatter, erpPriceInputFormatter, erpPriceMultiply, erpPriceTableColumnFormatter, fenToYuan, fenToYuanFormat, filterTree, floatToFixed2, formatAxis, formatDate, formatDate2, formatDateTime, formatFileSize, formatPast, formatPast2, formatTime, formatToFraction, formatToFractionDigit, generateAcceptedFileTypes, get, getCurrentTimezone, getDateRange, getDayRange, getElementVisibleRect, getFileIcon, getFileNameFromUrl, getFileTypeClass, getFirstNonNullOrUndefined, getLast1Year, getLast30Days, getLast7Days, getNestedValue, getScrollbarWidth, getSystemTimezone, getUrlNumberValue, getUrlValue, getWeek, groupBy, handleTree, isBoolean, isDate, isDayjsObject, isEmpty, isEqual, isFunction, isHttpUrl, isImage, isMacOs, isNumber, isObject, isSameDay, isString, isUndefined, isUrl, isWindow, isWindowsOs, jsonParse, kebabToCamelCase, loadScript, mapTree, merge, mergeWithArrayOverride, needsScrollbar, openRouteInNewWindow, openWindow, set, setCurrentTimezone, sortTree, startProgress, stopProgress, to, toCamelCase, toLowerCaseFirstLetter, traverseTreeValues, treeToString, triggerDownload, triggerWindowResize, uniqueByField, updateCSSVariables, urlToBase64, yuanToFen };