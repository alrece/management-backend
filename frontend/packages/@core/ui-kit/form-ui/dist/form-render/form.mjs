import { provideFormRenderProps } from "./context.mjs";
import { getBaseRules, getDefaultValueInZodStack } from "./helper.mjs";
import _sfc_main$1 from "./form-field.mjs";
import { useExpandable } from "./expandable.mjs";
import { Fragment, computed, createBlock, createCommentVNode, createElementBlock, createElementVNode, createVNode, defineComponent, guardReactiveProps, mergeProps, normalizeClass, normalizeProps, openBlock, renderList, renderSlot, resolveDynamicComponent, withCtx } from "vue";
import { Form } from "@vben-core/shadcn-ui";
import { cn, isFunction, isString, mergeWithArrayOverride } from "@vben-core/shared/utils";
//#region src/form-render/form.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "form",
	props: {
		arrayToStringFields: {},
		collapsed: { type: Boolean },
		collapsedRows: { default: 1 },
		collapseTriggerResize: { type: Boolean },
		commonConfig: { default: () => ({}) },
		compact: { type: Boolean },
		componentBindEventMap: {},
		componentMap: {},
		fieldMappingTime: {},
		form: {},
		layout: {},
		schema: {},
		showCollapseButton: {
			type: Boolean,
			default: false
		},
		wrapperClass: { default: "grid-cols-1 sm:grid-cols-2 md:grid-cols-3" },
		globalCommonConfig: { default: () => ({}) }
	},
	emits: ["submit"],
	setup(__props, { emit: __emit }) {
		const props = __props;
		const emits = __emit;
		const wrapperClass = computed(() => {
			const cls = ["flex"];
			if (props.layout === "inline") cls.push("flex-wrap gap-x-2");
			else cls.push(props.compact ? "gap-x-2" : "gap-x-4", "flex-col grid");
			return cn(...cls, props.wrapperClass);
		});
		provideFormRenderProps(props);
		const { isCalculated, keepFormItemIndex, wrapperRef } = useExpandable(props);
		const shapes = computed(() => {
			const resultShapes = [];
			props.schema?.forEach((schema) => {
				const { fieldName } = schema;
				const rules = schema.rules;
				let typeName = "";
				if (rules && !isString(rules)) typeName = rules._def.typeName;
				const baseRules = getBaseRules(rules);
				resultShapes.push({
					default: getDefaultValueInZodStack(rules),
					fieldName,
					required: !["ZodNullable", "ZodOptional"].includes(typeName),
					rules: baseRules
				});
			});
			return resultShapes;
		});
		const formComponent = computed(() => props.form ? "form" : Form);
		const formComponentProps = computed(() => {
			return props.form ? { onSubmit: props.form.handleSubmit((val) => emits("submit", val)) } : { onSubmit: (val) => emits("submit", val) };
		});
		const formCollapsed = computed(() => {
			return props.collapsed && isCalculated.value;
		});
		const computedSchema = computed(() => {
			const { colon = false, componentProps = {}, controlClass = "", disabled, disabledOnChangeListener = true, disabledOnInputListener = true, emptyStateValue = void 0, formFieldProps = {}, formItemClass = "", hideLabel = false, hideRequiredMark = false, labelClass = "", labelWidth = 100, modelPropName = "", wrapperClass = "" } = mergeWithArrayOverride(props.commonConfig, props.globalCommonConfig);
			return (props.schema || []).map((schema, index) => {
				const keepIndex = keepFormItemIndex.value;
				const hidden = props.showCollapseButton && !!formCollapsed.value && keepIndex ? keepIndex <= index : false;
				let resolvedSchemaFormItemClass = schema.formItemClass;
				if (isFunction(schema.formItemClass)) try {
					resolvedSchemaFormItemClass = schema.formItemClass();
				} catch (error) {
					console.error("Error calling formItemClass function:", error);
					resolvedSchemaFormItemClass = "";
				}
				return {
					colon,
					disabled,
					disabledOnChangeListener,
					disabledOnInputListener,
					emptyStateValue,
					hideLabel,
					hideRequiredMark,
					labelWidth,
					modelPropName,
					wrapperClass,
					...schema,
					commonComponentProps: componentProps,
					componentProps: schema.componentProps,
					controlClass: cn(controlClass, schema.controlClass),
					formFieldProps: {
						...formFieldProps,
						...schema.formFieldProps
					},
					formItemClass: cn("shrink-0", { hidden }, formItemClass, resolvedSchemaFormItemClass),
					labelClass: cn(labelClass, schema.labelClass)
				};
			});
		});
		return (_ctx, _cache) => {
			return openBlock(), createBlock(resolveDynamicComponent(formComponent.value), normalizeProps(guardReactiveProps(formComponentProps.value)), {
				default: withCtx(() => [createElementVNode("div", {
					ref_key: "wrapperRef",
					ref: wrapperRef,
					class: normalizeClass(wrapperClass.value)
				}, [(openBlock(true), createElementBlock(Fragment, null, renderList(computedSchema.value, (cSchema) => {
					return openBlock(), createElementBlock(Fragment, { key: cSchema.fieldName }, [createCommentVNode(" <div v-if=\"$slots[cSchema.fieldName]\" :class=\"cSchema.formItemClass\">\n          <slot :definition=\"cSchema\" :name=\"cSchema.fieldName\"> </slot>\n        </div> "), createVNode(_sfc_main$1, mergeProps({ ref_for: true }, cSchema, {
						class: cSchema.formItemClass,
						rules: cSchema.rules
					}), {
						default: withCtx((slotProps) => [renderSlot(_ctx.$slots, cSchema.fieldName, mergeProps({ ref_for: true }, slotProps))]),
						_: 2
					}, 1040, ["class", "rules"])], 64);
				}), 128)), renderSlot(_ctx.$slots, "default", { shapes: shapes.value })], 2)]),
				_: 3
			}, 16);
		};
	}
});
//#endregion
export { _sfc_main as default };
