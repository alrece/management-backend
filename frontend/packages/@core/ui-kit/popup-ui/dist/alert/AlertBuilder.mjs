import _sfc_main from "./alert2.mjs";
import { Input, VbenRenderContent } from "@vben-core/shadcn-ui";
import { h, nextTick, ref, render } from "vue";
import { useSimpleLocale } from "@vben-core/composables";
import { isFunction, isString } from "@vben-core/shared/utils";
//#region src/alert/AlertBuilder.ts
const alerts = ref([]);
const { $t } = useSimpleLocale();
function vbenAlert(arg0, arg1, arg2) {
	return new Promise((resolve, reject) => {
		const options = isString(arg0) ? { content: arg0 } : { ...arg0 };
		if (arg1) {
			if (isString(arg1)) options.title = arg1;
			else if (!isString(arg1)) Object.assign(options, arg1);
		}
		if (arg2 && !isString(arg2)) Object.assign(options, arg2);
		const container = document.createElement("div");
		document.body.append(container);
		const alertRef = {
			container,
			instance: null
		};
		const vnode = h(_sfc_main, {
			onClosed: (isConfirm) => {
				alerts.value = alerts.value.filter((item) => item !== alertRef);
				render(null, container);
				if (container.parentNode) container.remove();
				if (isConfirm) resolve();
				else reject(/* @__PURE__ */ new Error("dialog cancelled"));
			},
			...options,
			open: true,
			title: options.title ?? $t.value("prompt")
		});
		render(vnode, container);
		alertRef.instance = vnode.component?.proxy;
		alerts.value.push(alertRef);
	});
}
function vbenConfirm(arg0, arg1, arg2) {
	const defaultProps = { showCancel: true };
	if (!arg1) return isString(arg0) ? vbenAlert(arg0, defaultProps) : vbenAlert({
		...defaultProps,
		...arg0
	});
	else if (!arg2) return isString(arg1) ? vbenAlert(arg0, arg1, defaultProps) : vbenAlert(arg0, {
		...defaultProps,
		...arg1
	});
	return vbenAlert(arg0, arg1, {
		...defaultProps,
		...arg2
	});
}
async function vbenPrompt(options) {
	const { component: _component, componentProps: _componentProps, componentSlots, content, defaultValue, modelPropName: _modelPropName, ...delegated } = options;
	const modelValue = ref(defaultValue);
	const inputComponentRef = ref(null);
	const staticContents = [h(VbenRenderContent, {
		content,
		renderBr: true
	})];
	const modelPropName = _modelPropName || "modelValue";
	const componentProps = { ..._componentProps };
	const contentRenderer = () => {
		const currentProps = {
			...componentProps,
			[modelPropName]: modelValue.value,
			[`onUpdate:${modelPropName}`]: (val) => {
				modelValue.value = val;
			}
		};
		inputComponentRef.value = h(_component || Input, currentProps, componentSlots);
		return h("div", { class: "flex flex-col gap-2" }, { default: () => [...staticContents, inputComponentRef.value] });
	};
	await vbenConfirm({
		...delegated,
		async beforeClose(scope) {
			if (delegated.beforeClose) return await delegated.beforeClose({
				...scope,
				value: modelValue.value
			});
		},
		content: contentRenderer,
		contentMasking: true,
		async onOpened() {
			await nextTick();
			const componentRef = inputComponentRef.value;
			if (componentRef) {
				if (componentRef.component?.exposed && isFunction(componentRef.component.exposed.focus)) componentRef.component.exposed.focus();
				else if (componentRef.el) {
					if (isFunction(componentRef.el.focus) && [
						"BUTTON",
						"INPUT",
						"SELECT",
						"TEXTAREA"
					].includes(componentRef.el.tagName)) componentRef.el.focus();
					else if (isFunction(componentRef.el.querySelector)) {
						const focusableElement = componentRef.el.querySelector("input, select, textarea, button");
						if (focusableElement && isFunction(focusableElement.focus)) focusableElement.focus();
					} else if (componentRef.el.nextElementSibling && isFunction(componentRef.el.nextElementSibling.focus)) componentRef.el.nextElementSibling.focus();
				}
			}
		}
	});
	return modelValue.value;
}
function clearAllAlerts() {
	alerts.value.forEach((alert) => {
		render(null, alert.container);
		if (alert.container.parentNode) alert.container.remove();
	});
	alerts.value = [];
}
//#endregion
export { clearAllAlerts, vbenAlert, vbenConfirm, vbenPrompt };
