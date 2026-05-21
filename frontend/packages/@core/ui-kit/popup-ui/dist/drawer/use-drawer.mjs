import _sfc_main from "./drawer.mjs";
import { DrawerApi } from "./drawer-api.mjs";
import { defineComponent, h, inject, nextTick, provide, reactive, ref } from "vue";
import { useStore } from "@vben-core/shared/store";
//#region src/drawer/use-drawer.ts
const USER_DRAWER_INJECT_KEY = Symbol("VBEN_DRAWER_INJECT");
const DEFAULT_DRAWER_PROPS = {};
function setDefaultDrawerProps(props) {
	Object.assign(DEFAULT_DRAWER_PROPS, props);
}
function useVbenDrawer(options = {}) {
	const { connectedComponent } = options;
	if (connectedComponent) {
		const extendedApi = reactive({});
		const isDrawerReady = ref(true);
		return [defineComponent((props, { attrs, slots }) => {
			provide(USER_DRAWER_INJECT_KEY, {
				extendApi(api) {
					Object.setPrototypeOf(extendedApi, api);
				},
				options,
				async reCreateDrawer() {
					isDrawerReady.value = false;
					await nextTick();
					isDrawerReady.value = true;
				}
			});
			checkProps(extendedApi, {
				...props,
				...attrs,
				...slots
			});
			return () => h(isDrawerReady.value ? connectedComponent : "div", {
				...props,
				...attrs
			}, slots);
		}, {
			name: "VbenParentDrawer",
			inheritAttrs: false
		}), extendedApi];
	}
	const injectData = inject(USER_DRAWER_INJECT_KEY, {});
	const mergedOptions = {
		...DEFAULT_DRAWER_PROPS,
		...injectData.options,
		...options
	};
	mergedOptions.onOpenChange = (isOpen) => {
		options.onOpenChange?.(isOpen);
		injectData.options?.onOpenChange?.(isOpen);
	};
	const onClosed = mergedOptions.onClosed;
	mergedOptions.onClosed = () => {
		onClosed?.();
		if (mergedOptions.destroyOnClose) injectData.reCreateDrawer?.();
	};
	const api = new DrawerApi(mergedOptions);
	const extendedApi = api;
	extendedApi.useStore = (selector) => {
		return useStore(api.store, selector);
	};
	const Drawer = defineComponent((props, { attrs, slots }) => {
		return () => h(_sfc_main, {
			...props,
			...attrs,
			drawerApi: extendedApi
		}, slots);
	}, {
		name: "VbenDrawer",
		inheritAttrs: false
	});
	injectData.extendApi?.(extendedApi);
	return [Drawer, extendedApi];
}
async function checkProps(api, attrs) {
	if (!attrs || Object.keys(attrs).length === 0) return;
	await nextTick();
	const state = api?.store?.state;
	if (!state) return;
	const stateKeys = new Set(Object.keys(state));
	for (const attr of Object.keys(attrs)) if (stateKeys.has(attr) && !["class"].includes(attr)) console.warn(`[Vben Drawer]: When 'connectedComponent' exists, do not set props or slots '${attr}', which will increase complexity. If you need to modify the props of Drawer, please use useVbenDrawer or api.`);
}
//#endregion
export { setDefaultDrawerProps, useVbenDrawer };
