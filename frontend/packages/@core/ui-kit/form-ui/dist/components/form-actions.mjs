import { COMPONENT_MAP } from "../config.mjs";
import { injectFormProps } from "../use-form-context.mjs";
import { Fragment, computed, createBlock, createCommentVNode, createElementBlock, createElementVNode, createTextVNode, defineComponent, mergeProps, normalizeClass, openBlock, renderSlot, resolveDynamicComponent, toDisplayString, toRaw, unref, useModel, watch, withCtx } from "vue";
import { VbenExpandableArrow } from "@vben-core/shadcn-ui";
import { cn, isFunction, triggerWindowResize } from "@vben-core/shared/utils";
import { useSimpleLocale } from "@vben-core/composables";
//#region src/components/form-actions.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "form-actions",
	props: {
		"modelValue": { default: false },
		"modelModifiers": {}
	},
	emits: ["update:modelValue"],
	setup(__props, { expose: __expose }) {
		const { $t } = useSimpleLocale();
		const [rootProps, form] = injectFormProps();
		const collapsed = useModel(__props, "modelValue");
		const resetButtonOptions = computed(() => {
			return {
				content: `${$t.value("reset")}`,
				show: true,
				...unref(rootProps).resetButtonOptions
			};
		});
		const submitButtonOptions = computed(() => {
			return {
				content: `${$t.value("submit")}`,
				show: true,
				...unref(rootProps).submitButtonOptions
			};
		});
		async function handleSubmit(e) {
			e?.preventDefault();
			e?.stopPropagation();
			const props = unref(rootProps);
			if (!props.formApi) return;
			const { valid } = await props.formApi.validate();
			if (!valid) return;
			const values = toRaw(await props.formApi.getValues()) ?? {};
			await props.handleSubmit?.(values);
		}
		async function handleReset(e) {
			e?.preventDefault();
			e?.stopPropagation();
			const props = unref(rootProps);
			const values = toRaw(await props.formApi?.getValues()) ?? {};
			if (isFunction(props.handleReset)) await props.handleReset?.(values);
			else form.resetForm();
		}
		watch(() => collapsed.value, () => {
			if (unref(rootProps).collapseTriggerResize) triggerWindowResize();
		});
		const actionWrapperClass = computed(() => {
			const props = unref(rootProps);
			const actionLayout = props.actionLayout || "rowEnd";
			const actionPosition = props.actionPosition || "right";
			const cls = [
				"flex",
				"items-center",
				"gap-3",
				props.compact ? "pb-2" : "pb-4",
				props.layout === "vertical" ? "self-end" : "self-center",
				props.layout === "inline" ? "" : "w-full",
				props.actionWrapperClass
			];
			switch (actionLayout) {
				case "newLine":
					cls.push("col-span-full");
					break;
				case "rowEnd":
					cls.push("col-[-2/-1]");
					break;
			}
			switch (actionPosition) {
				case "center":
					cls.push("justify-center");
					break;
				case "left":
					cls.push("justify-start");
					break;
				default:
					cls.push("justify-end");
					break;
			}
			return cls.join(" ");
		});
		__expose({
			handleReset,
			handleSubmit
		});
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("div", { class: normalizeClass(unref(cn)(actionWrapperClass.value)) }, [
				unref(rootProps).actionButtonsReverse ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
					createCommentVNode(" 提交按钮前 "),
					renderSlot(_ctx.$slots, "submit-before"),
					submitButtonOptions.value.show ? (openBlock(), createBlock(resolveDynamicComponent(unref(COMPONENT_MAP).PrimaryButton), mergeProps({
						key: 0,
						type: "button",
						onClick: handleSubmit
					}, submitButtonOptions.value), {
						default: withCtx(() => [createTextVNode(toDisplayString(submitButtonOptions.value.content), 1)]),
						_: 1
					}, 16)) : createCommentVNode("v-if", true)
				], 64)) : createCommentVNode("v-if", true),
				createCommentVNode(" 重置按钮前 "),
				renderSlot(_ctx.$slots, "reset-before"),
				resetButtonOptions.value.show ? (openBlock(), createBlock(resolveDynamicComponent(unref(COMPONENT_MAP).DefaultButton), mergeProps({
					key: 1,
					type: "button",
					onClick: handleReset
				}, resetButtonOptions.value), {
					default: withCtx(() => [createTextVNode(toDisplayString(resetButtonOptions.value.content), 1)]),
					_: 1
				}, 16)) : createCommentVNode("v-if", true),
				!unref(rootProps).actionButtonsReverse ? (openBlock(), createElementBlock(Fragment, { key: 2 }, [
					createCommentVNode(" 提交按钮前 "),
					renderSlot(_ctx.$slots, "submit-before"),
					submitButtonOptions.value.show ? (openBlock(), createBlock(resolveDynamicComponent(unref(COMPONENT_MAP).PrimaryButton), mergeProps({
						key: 0,
						type: "button",
						onClick: handleSubmit
					}, submitButtonOptions.value), {
						default: withCtx(() => [createTextVNode(toDisplayString(submitButtonOptions.value.content), 1)]),
						_: 1
					}, 16)) : createCommentVNode("v-if", true)
				], 64)) : createCommentVNode("v-if", true),
				createCommentVNode(" 展开按钮前 "),
				renderSlot(_ctx.$slots, "expand-before"),
				unref(rootProps).showCollapseButton ? (openBlock(), createBlock(unref(VbenExpandableArrow), {
					key: 3,
					class: "ml-[-0.3em]",
					"model-value": collapsed.value,
					"onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => collapsed.value = $event)
				}, {
					default: withCtx(() => [createElementVNode("span", null, toDisplayString(collapsed.value ? unref($t)("expand") : unref($t)("collapse")), 1)]),
					_: 1
				}, 8, ["model-value"])) : createCommentVNode("v-if", true),
				createCommentVNode(" 展开按钮后 "),
				renderSlot(_ctx.$slots, "expand-after")
			], 2);
		};
	}
});
//#endregion
export { _sfc_main as default };
