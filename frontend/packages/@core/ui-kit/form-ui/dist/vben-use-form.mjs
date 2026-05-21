import { COMPONENT_BIND_EVENT_MAP, COMPONENT_MAP, DEFAULT_FORM_COMMON_CONFIG } from "./config.mjs";
import { provideComponentRefMap, provideFormProps, useFormInitial } from "./use-form-context.mjs";
import _sfc_main$1 from "./components/form-actions.mjs";
import _sfc_main$2 from "./form-render/form.mjs";
import { createBlock, createCommentVNode, createSlots, defineComponent, guardReactiveProps, mergeProps, nextTick, normalizeProps, onMounted, openBlock, renderList, renderSlot, unref, watch, withCtx, withKeys } from "vue";
import { cloneDeep, get, isEqual, set } from "@vben-core/shared/utils";
import { useForwardPriorityValues } from "@vben-core/composables";
import { useDebounceFn } from "@vueuse/core";
//#region src/vben-use-form.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "vben-use-form",
	props: {
		formApi: {},
		actionButtonsReverse: { type: Boolean },
		actionLayout: {},
		actionPosition: {},
		actionWrapperClass: { type: [
			Array,
			Boolean,
			null,
			Object,
			String
		] },
		arrayToStringFields: {},
		fieldMappingTime: {},
		handleCollapsedChange: { type: Function },
		handleReset: { type: Function },
		handleSubmit: { type: Function },
		handleValuesChange: { type: Function },
		resetButtonOptions: {},
		scrollToFirstError: { type: Boolean },
		showDefaultActions: { type: Boolean },
		submitButtonOptions: {},
		submitOnChange: { type: Boolean },
		submitOnEnter: { type: Boolean },
		collapsed: { type: Boolean },
		collapsedRows: {},
		collapseTriggerResize: { type: Boolean },
		commonConfig: {},
		compact: { type: Boolean },
		layout: {},
		schema: {},
		showCollapseButton: { type: Boolean },
		wrapperClass: {}
	},
	setup(__props) {
		const props = __props;
		const state = props.formApi?.useStore?.();
		const forward = useForwardPriorityValues(props, state);
		const componentRefMap = /* @__PURE__ */ new Map();
		const { delegatedSlots, form } = useFormInitial(forward);
		provideFormProps([forward, form]);
		provideComponentRefMap(componentRefMap);
		props.formApi?.mount?.(form, componentRefMap);
		const handleUpdateCollapsed = (value) => {
			props.formApi?.setState({ collapsed: value });
			forward.value.handleCollapsedChange?.(value);
		};
		function handleKeyDownEnter(event) {
			if (!state?.value.submitOnEnter || !forward.value.formApi?.isMounted) return;
			if (event.target instanceof HTMLTextAreaElement) return;
			event.preventDefault();
			forward.value.formApi?.validateAndSubmitForm();
		}
		const handleValuesChangeDebounced = useDebounceFn(async () => {
			state?.value.submitOnChange && forward.value.formApi?.validateAndSubmitForm();
		}, 300);
		const valuesCache = {};
		onMounted(async () => {
			await nextTick();
			watch(() => form.values, async (newVal) => {
				if (forward.value.handleValuesChange) {
					const fields = state?.value.schema?.map((item) => {
						return item.fieldName;
					});
					if (fields && fields.length > 0) {
						const changedFields = [];
						fields.forEach((field) => {
							const newFieldValue = get(newVal, field);
							if (!isEqual(newFieldValue, get(valuesCache, field))) {
								changedFields.push(field);
								set(valuesCache, field, newFieldValue);
							}
						});
						if (changedFields.length > 0) {
							const values = await forward.value.formApi?.getValues();
							forward.value.handleValuesChange(cloneDeep(values ?? {}), changedFields);
						}
					}
				}
				handleValuesChangeDebounced();
			}, { deep: true });
		});
		return (_ctx, _cache) => {
			return openBlock(), createBlock(unref(_sfc_main$2), mergeProps({ onKeydown: withKeys(handleKeyDownEnter, ["enter"]) }, unref(forward), {
				collapsed: unref(state)?.collapsed,
				"component-bind-event-map": unref(COMPONENT_BIND_EVENT_MAP),
				"component-map": unref(COMPONENT_MAP),
				form: unref(form),
				"global-common-config": unref(DEFAULT_FORM_COMMON_CONFIG)
			}), createSlots({
				default: withCtx((slotProps) => [renderSlot(_ctx.$slots, "default", normalizeProps(guardReactiveProps(slotProps)), () => [unref(forward).showDefaultActions ? (openBlock(), createBlock(_sfc_main$1, {
					key: 0,
					"model-value": unref(state)?.collapsed,
					"onUpdate:modelValue": handleUpdateCollapsed
				}, {
					"reset-before": withCtx((resetSlotProps) => [renderSlot(_ctx.$slots, "reset-before", normalizeProps(guardReactiveProps(resetSlotProps)))]),
					"submit-before": withCtx((submitSlotProps) => [renderSlot(_ctx.$slots, "submit-before", normalizeProps(guardReactiveProps(submitSlotProps)))]),
					"expand-before": withCtx((expandBeforeSlotProps) => [renderSlot(_ctx.$slots, "expand-before", normalizeProps(guardReactiveProps(expandBeforeSlotProps)))]),
					"expand-after": withCtx((expandAfterSlotProps) => [renderSlot(_ctx.$slots, "expand-after", normalizeProps(guardReactiveProps(expandAfterSlotProps)))]),
					_: 3
				}, 8, ["model-value"])) : createCommentVNode("v-if", true)])]),
				_: 2
			}, [renderList(unref(delegatedSlots), (slotName) => {
				return {
					name: slotName,
					fn: withCtx((slotProps) => [renderSlot(_ctx.$slots, slotName, normalizeProps(guardReactiveProps(slotProps)))])
				};
			})]), 1040, [
				"collapsed",
				"component-bind-event-map",
				"component-map",
				"form",
				"global-common-config"
			]);
		};
	}
});
//#endregion
export { _sfc_main as default };
