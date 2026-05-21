import { BaseFormComponentType, VbenFormAdapterOptions } from "./types.js";
import { Component } from "vue";

//#region src/config.d.ts
declare function setupVbenForm<T extends BaseFormComponentType = BaseFormComponentType>(options: VbenFormAdapterOptions<T>): void;
//#endregion
export { setupVbenForm };