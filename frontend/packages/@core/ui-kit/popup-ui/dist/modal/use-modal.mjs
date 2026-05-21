import _sfc_main from "./modal.mjs";
import { ModalApi } from "./modal-api.mjs";
import { defineComponent, h, inject, nextTick, provide, reactive, ref } from "vue";
import { useStore } from "@vben-core/shared/store";
//#region src/modal/use-modal.ts
const USER_MODAL_INJECT_KEY = Symbol("VBEN_MODAL_INJECT");
const DEFAULT_MODAL_PROPS = {};
function setDefaultModalProps(props) {
	Object.assign(DEFAULT_MODAL_PROPS, props);
}
function useVbenModal(options = {}) {
	const { connectedComponent } = options;
	if (connectedComponent) {
		const extendedApi = reactive({});
		const isModalReady = ref(true);
		return [defineComponent((props, { attrs, slots }) => {
			provide(USER_MODAL_INJECT_KEY, {
				extendApi(api) {
					Object.setPrototypeOf(extendedApi, api);
				},
				consumed: false,
				options,
				async reCreateModal() {
					isModalReady.value = false;
					await nextTick();
					isModalReady.value = true;
				}
			});
			checkProps(extendedApi, {
				...props,
				...attrs,
				...slots
			});
			return () => h(isModalReady.value ? connectedComponent : "div", {
				...props,
				...attrs
			}, slots);
		}, {
			name: "VbenParentModal",
			inheritAttrs: false
		}), extendedApi];
	}
	let injectData = inject(USER_MODAL_INJECT_KEY, {});
	if (injectData.consumed) injectData = {};
	else injectData.consumed = true;
	const mergedOptions = {
		...DEFAULT_MODAL_PROPS,
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
		if (mergedOptions.destroyOnClose) {
			injectData.consumed = false;
			injectData.reCreateModal?.();
		}
	};
	const api = new ModalApi(mergedOptions);
	const extendedApi = api;
	extendedApi.useStore = (selector) => {
		return useStore(api.store, selector);
	};
	const Modal = defineComponent((props, { attrs, slots }) => {
		return () => h(_sfc_main, {
			...props,
			...attrs,
			modalApi: extendedApi
		}, slots);
	}, {
		name: "VbenModal",
		inheritAttrs: false
	});
	injectData.extendApi?.(extendedApi);
	return [Modal, extendedApi];
}
async function checkProps(api, attrs) {
	if (!attrs || Object.keys(attrs).length === 0) return;
	await nextTick();
	const state = api?.store?.state;
	if (!state) return;
	const stateKeys = new Set(Object.keys(state));
	for (const attr of Object.keys(attrs)) if (stateKeys.has(attr) && !["class"].includes(attr)) console.warn(`[Vben Modal]: When 'connectedComponent' exists, do not set props or slots '${attr}', which will increase complexity. If you need to modify the props of Modal, please use useVbenModal or api.`);
}
//#endregion
export { setDefaultModalProps, useVbenModal };
