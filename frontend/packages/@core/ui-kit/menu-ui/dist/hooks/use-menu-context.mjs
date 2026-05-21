import { findComponentUpward } from "../utils/index.mjs";
import { getCurrentInstance, inject, provide } from "vue";
//#region src/hooks/use-menu-context.ts
const menuContextKey = Symbol("menuContext");
/**
* @zh_CN Provide menu context
*/
function createMenuContext(injectMenuData) {
	provide(menuContextKey, injectMenuData);
}
/**
* @zh_CN Provide menu context
*/
function createSubMenuContext(injectSubMenuData) {
	provide(`subMenu:${getCurrentInstance()?.uid}`, injectSubMenuData);
}
/**
* @zh_CN Inject menu context
*/
function useMenuContext() {
	if (!getCurrentInstance()) throw new Error("instance is required");
	return inject(menuContextKey);
}
/**
* @zh_CN Inject menu context
*/
function useSubMenuContext() {
	const instance = getCurrentInstance();
	if (!instance) throw new Error("instance is required");
	return inject(`subMenu:${findComponentUpward(instance, ["Menu", "SubMenu"])?.uid}`);
}
//#endregion
export { createMenuContext, createSubMenuContext, useMenuContext, useSubMenuContext };
