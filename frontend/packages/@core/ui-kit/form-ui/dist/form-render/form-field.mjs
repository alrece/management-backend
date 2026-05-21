import { injectComponentRefMap } from "../use-form-context.mjs";
import { injectRenderFormProps, useFormContext } from "./context.mjs";
import useDependencies from "./dependencies.mjs";
import _sfc_main$1 from "./form-label.mjs";
import { isEventObjectLike } from "./helper.mjs";
import { Transition, computed, createBlock, createCommentVNode, createElementBlock, createElementVNode, createSlots, createVNode, defineComponent, guardReactiveProps, mergeProps, nextTick, normalizeClass, normalizeProps, normalizeStyle, onUnmounted, openBlock, renderList, renderSlot, resolveDynamicComponent, unref, useTemplateRef, vShow, watch, withCtx, withDirectives } from "vue";
import { FormControl, FormDescription, FormField, FormItem, FormMessage, VbenRenderContent, VbenTooltip } from "@vben-core/shadcn-ui";
import { useFieldError, useFormValues } from "vee-validate";
import { cn, isFunction, isObject, isString } from "@vben-core/shared/utils";
import { CircleAlert } from "@vben-core/icons";
import { toTypedSchema } from "@vee-validate/zod";
//#region src/form-render/form-field.vue
const _hoisted_1 = { class: "flex-auto overflow-hidden p-px" };
const _hoisted_2 = {
	key: 0,
	class: "ml-1"
};
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "form-field",
	props: {
		component: {},
		componentProps: { type: Function },
		defaultValue: {},
		dependencies: {},
		description: { type: [Function, String] },
		fieldName: {},
		help: { type: [Function, String] },
		hide: { type: Boolean },
		label: { type: [Function, String] },
		renderComponentContent: { type: Function },
		rules: {},
		suffix: { type: [Function, String] },
		valueFormat: { type: Function },
		colon: { type: Boolean },
		controlClass: {},
		disabled: { type: Boolean },
		disabledOnChangeListener: { type: Boolean },
		disabledOnInputListener: { type: Boolean },
		emptyStateValue: {},
		formFieldProps: {},
		formItemClass: { type: [Function, String] },
		hideLabel: { type: Boolean },
		hideRequiredMark: { type: Boolean },
		labelClass: {},
		labelWidth: {},
		modelPropName: {},
		wrapperClass: {},
		commonComponentProps: {}
	},
	setup(__props) {
		const { componentBindEventMap, componentMap, isVertical } = useFormContext();
		const formRenderProps = injectRenderFormProps();
		const values = useFormValues();
		const errors = useFieldError(__props.fieldName);
		const fieldComponentRef = useTemplateRef("fieldComponentRef");
		const formApi = formRenderProps.form;
		const compact = computed(() => formRenderProps.compact);
		const isInValid = computed(() => errors.value?.length > 0);
		function getFormApi() {
			if (!formApi) throw new Error("Form api is required in <FormField />");
			return formApi;
		}
		const FieldComponent = computed(() => {
			const finalComponent = isString(__props.component) ? componentMap.value[__props.component] : __props.component;
			if (!finalComponent) console.warn(`Component ${__props.component} is not registered`);
			return finalComponent;
		});
		const { dynamicComponentProps, dynamicRules, isDisabled, isIf, isRequired, isShow } = useDependencies(() => __props.dependencies);
		const labelStyle = computed(() => {
			return __props.labelClass?.includes("w-") || isVertical.value ? {} : { width: `${__props.labelWidth}px` };
		});
		const currentRules = computed(() => {
			return dynamicRules.value || __props.rules;
		});
		const visible = computed(() => {
			return !__props.hide && isIf.value && isShow.value;
		});
		const shouldRequired = computed(() => {
			if (!visible.value) return false;
			if (!currentRules.value) return isRequired.value;
			if (isRequired.value) return true;
			if (isString(currentRules.value)) return [
				"mobileRequired",
				"required",
				"selectRequired"
			].includes(currentRules.value);
			let isOptional = currentRules?.value?.isOptional?.();
			if (currentRules?.value?._def?.typeName === "ZodDefault") {
				const innerType = currentRules?.value?._def.innerType;
				if (innerType) isOptional = innerType.isOptional?.();
			}
			return !isOptional;
		});
		const fieldRules = computed(() => {
			if (!visible.value) return null;
			let rules = currentRules.value;
			if (!rules) return isRequired.value ? "required" : null;
			if (isString(rules)) return rules;
			if (!!shouldRequired.value) {
				const unwrappedRules = rules?.unwrap?.();
				if (unwrappedRules) rules = unwrappedRules;
			}
			return toTypedSchema(rules);
		});
		const computedProps = computed(() => {
			const finalComponentProps = isFunction(__props.componentProps) ? __props.componentProps(values.value, getFormApi()) : __props.componentProps;
			return {
				...__props.commonComponentProps,
				...finalComponentProps,
				...dynamicComponentProps.value
			};
		});
		const computedHelp = computed(() => {
			const helpContent = __props.help;
			if (!helpContent) return;
			return () => isFunction(helpContent) ? helpContent(values.value, getFormApi()) : helpContent;
		});
		watch(() => computedProps.value?.autofocus, (value) => {
			if (value === true) nextTick(() => {
				autofocus();
			});
		}, { immediate: true });
		const shouldDisabled = computed(() => {
			return isDisabled.value || __props.disabled || computedProps.value?.disabled;
		});
		const customContentRender = computed(() => {
			if (!isFunction(__props.renderComponentContent)) return {};
			return __props.renderComponentContent(values.value, getFormApi());
		});
		const renderContentKey = computed(() => {
			return Object.keys(customContentRender.value);
		});
		const fieldProps = computed(() => {
			const rules = fieldRules.value;
			return {
				keepValue: true,
				label: isString(__props.label) ? __props.label : "",
				...rules ? { rules } : {},
				...__props.formFieldProps
			};
		});
		function fieldBindEvent(slotProps) {
			const modelValue = slotProps.componentField.modelValue;
			const handler = slotProps.componentField["onUpdate:modelValue"];
			const bindEventField = __props.modelPropName || (isString(__props.component) ? componentBindEventMap.value?.[__props.component] : null);
			let value = modelValue;
			if (modelValue && isObject(modelValue) && bindEventField) value = isEventObjectLike(modelValue) ? modelValue?.target?.[bindEventField] : modelValue?.[bindEventField] ?? modelValue;
			if (bindEventField) return {
				[`onUpdate:${bindEventField}`]: handler,
				[bindEventField]: value === void 0 ? __props.emptyStateValue : value,
				onChange: __props.disabledOnChangeListener ? void 0 : (e) => {
					const shouldUnwrap = isEventObjectLike(e);
					const onChange = slotProps?.componentField?.onChange;
					if (!shouldUnwrap) return onChange?.(e);
					return onChange?.(e?.target?.[bindEventField] ?? e);
				},
				...__props.disabledOnInputListener ? { onInput: void 0 } : {}
			};
			return {
				...__props.disabledOnInputListener ? { onInput: void 0 } : {},
				...__props.disabledOnChangeListener ? { onChange: void 0 } : {}
			};
		}
		function createComponentProps(slotProps) {
			const bindEvents = fieldBindEvent(slotProps);
			return {
				...slotProps.componentField,
				...computedProps.value,
				...bindEvents,
				...Reflect.has(computedProps.value, "onChange") ? { onChange: computedProps.value.onChange } : {},
				...Reflect.has(computedProps.value, "onInput") ? { onInput: computedProps.value.onInput } : {}
			};
		}
		function autofocus() {
			if (fieldComponentRef.value && isFunction(fieldComponentRef.value.focus) && document.activeElement !== fieldComponentRef.value) fieldComponentRef.value?.focus?.();
		}
		const componentRefMap = injectComponentRefMap();
		watch(fieldComponentRef, (componentRef) => {
			componentRefMap?.set(__props.fieldName, componentRef);
		});
		onUnmounted(() => {
			if (componentRefMap?.has(__props.fieldName)) componentRefMap.delete(__props.fieldName);
		});
		return (_ctx, _cache) => {
			return !__props.hide && unref(isIf) ? (openBlock(), createBlock(unref(FormField), mergeProps({ key: 0 }, fieldProps.value, { name: __props.fieldName }), {
				default: withCtx((slotProps) => [withDirectives(createVNode(unref(FormItem), mergeProps({ class: [{
					"form-valid-error": isInValid.value,
					"form-is-required": shouldRequired.value,
					"flex-col": unref(isVertical),
					"flex-row items-center": !unref(isVertical),
					"pb-4": !compact.value,
					"pb-2": compact.value
				}, "relative flex"] }, _ctx.$attrs), {
					default: withCtx(() => [!__props.hideLabel ? (openBlock(), createBlock(_sfc_main$1, {
						key: 0,
						class: normalizeClass(unref(cn)("flex leading-6", {
							"mr-2 shrink-0 justify-end": !unref(isVertical),
							"mb-1 flex-row": unref(isVertical)
						}, __props.labelClass)),
						help: computedHelp.value,
						colon: __props.colon,
						label: __props.label,
						required: shouldRequired.value && !__props.hideRequiredMark,
						style: normalizeStyle(labelStyle.value)
					}, {
						default: withCtx(() => [__props.label ? (openBlock(), createBlock(unref(VbenRenderContent), {
							key: 0,
							content: __props.label
						}, null, 8, ["content"])) : createCommentVNode("v-if", true)]),
						_: 1
					}, 8, [
						"class",
						"help",
						"colon",
						"label",
						"required",
						"style"
					])) : createCommentVNode("v-if", true), createElementVNode("div", _hoisted_1, [
						createElementVNode("div", { class: normalizeClass(unref(cn)("relative flex w-full items-center", __props.wrapperClass)) }, [
							createVNode(unref(FormControl), { class: normalizeClass(unref(cn)(__props.controlClass)) }, {
								default: withCtx(() => [renderSlot(_ctx.$slots, "default", normalizeProps(guardReactiveProps({
									...slotProps,
									...createComponentProps(slotProps),
									disabled: shouldDisabled.value,
									isInValid: isInValid.value
								})), () => [(openBlock(), createBlock(resolveDynamicComponent(FieldComponent.value), mergeProps({
									ref_key: "fieldComponentRef",
									ref: fieldComponentRef,
									class: { "border-destructive hover:border-destructive/80 focus:border-destructive focus:shadow-[0_0_0_2px_rgba(255,38,5,0.06)]": isInValid.value }
								}, createComponentProps(slotProps), { disabled: shouldDisabled.value }), createSlots({ _: 2 }, [renderList(renderContentKey.value, (name) => {
									return {
										name,
										fn: withCtx((renderSlotProps) => [createVNode(unref(VbenRenderContent), mergeProps({ content: customContentRender.value[name] }, {
											...renderSlotProps,
											formContext: slotProps
										}), null, 16, ["content"])])
									};
								})]), 1040, ["class", "disabled"])), compact.value && isInValid.value ? (openBlock(), createBlock(unref(VbenTooltip), {
									key: 0,
									"delay-duration": 300,
									side: "left"
								}, {
									trigger: withCtx(() => [renderSlot(_ctx.$slots, "trigger", {}, () => [createVNode(unref(CircleAlert), { class: normalizeClass(unref(cn)("inline-flex size-5 cursor-pointer text-foreground/80 hover:text-foreground")) }, null, 8, ["class"])])]),
									default: withCtx(() => [createVNode(unref(FormMessage))]),
									_: 3
								})) : createCommentVNode("v-if", true)])]),
								_: 2
							}, 1032, ["class"]),
							createCommentVNode(" 自定义后缀 "),
							__props.suffix ? (openBlock(), createElementBlock("div", _hoisted_2, [createVNode(unref(VbenRenderContent), { content: __props.suffix }, null, 8, ["content"])])) : createCommentVNode("v-if", true)
						], 2),
						__props.description ? (openBlock(), createBlock(unref(FormDescription), {
							key: 0,
							class: "text-xs"
						}, {
							default: withCtx(() => [createVNode(unref(VbenRenderContent), { content: __props.description }, null, 8, ["content"])]),
							_: 1
						})) : createCommentVNode("v-if", true),
						!compact.value ? (openBlock(), createBlock(Transition, {
							key: 1,
							name: "slide-up"
						}, {
							default: withCtx(() => [createVNode(unref(FormMessage), { class: "absolute" })]),
							_: 1
						})) : createCommentVNode("v-if", true)
					])]),
					_: 2
				}, 1040, ["class"]), [[vShow, unref(isShow)]])]),
				_: 3
			}, 16, ["name"])) : createCommentVNode("v-if", true);
		};
	}
});
//#endregion
export { _sfc_main as default };
