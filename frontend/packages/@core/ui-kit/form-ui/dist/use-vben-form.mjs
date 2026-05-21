import { FormApi } from "./form-api.mjs";
import _sfc_main from "./vben-use-form.mjs";
import { defineComponent, h, isReactive, onBeforeUnmount, watch } from "vue";
import { useStore } from "@vben-core/shared/store";
//#region src/use-vben-form.ts
function useVbenForm(options) {
	const IS_REACTIVE = isReactive(options);
	const api = new FormApi(options);
	const extendedApi = api;
	extendedApi.useStore = (selector) => {
		return useStore(api.store, selector);
	};
	const Form = defineComponent((props, { attrs, slots }) => {
		onBeforeUnmount(() => {
			api.unmount();
		});
		api.setState({
			...props,
			...attrs
		});
		return () => h(_sfc_main, {
			...props,
			...attrs,
			formApi: extendedApi
		}, slots);
	}, {
		name: "VbenUseForm",
		inheritAttrs: false
	});
	if (IS_REACTIVE) watch(() => options.schema, () => {
		api.setState({ schema: options.schema });
	}, { immediate: true });
	return [Form, extendedApi];
}
//#endregion
export { useVbenForm };
