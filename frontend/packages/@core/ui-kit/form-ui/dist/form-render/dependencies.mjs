import { resolveFieldNamePath } from "../field-name.mjs";
import { injectRenderFormProps } from "./context.mjs";
import { computed, ref, watch } from "vue";
import { useFormValues } from "vee-validate";
import { get, isBoolean, isFunction } from "@vben-core/shared/utils";
//#region src/form-render/dependencies.ts
/**
* 解析Nested Objects对应的字段值
* @param values 表单值
* @param fieldName 字段名
*/
function resolveValueByFieldName(values, fieldName) {
	const { rawKey } = resolveFieldNamePath(fieldName);
	if (rawKey) return values[rawKey];
	return get(values, fieldName);
}
function useDependencies(getDependencies) {
	const values = useFormValues();
	const formApi = injectRenderFormProps().form;
	if (!formApi) throw new Error("Form api is required in useDependencies");
	if (!values) throw new Error("useDependencies should be used within <VbenForm>");
	const isIf = ref(true);
	const isDisabled = ref(false);
	const isShow = ref(true);
	const isRequired = ref(false);
	const dynamicComponentProps = ref({});
	const dynamicRules = ref();
	const triggerFieldValues = computed(() => {
		return (getDependencies()?.triggerFields ?? []).map((dep) => {
			return resolveValueByFieldName(values.value, dep);
		});
	});
	const resetConditionState = () => {
		isDisabled.value = false;
		isIf.value = true;
		isShow.value = true;
		isRequired.value = false;
		dynamicRules.value = void 0;
		dynamicComponentProps.value = {};
	};
	watch([triggerFieldValues, getDependencies], async ([_values, dependencies]) => {
		if (!dependencies || !dependencies?.triggerFields?.length) return;
		resetConditionState();
		const { componentProps, disabled, if: whenIf, required, rules, show, trigger } = dependencies;
		const formValues = values.value;
		if (isFunction(whenIf)) {
			isIf.value = !!await whenIf(formValues, formApi);
			if (!isIf.value) return;
		} else if (isBoolean(whenIf)) {
			isIf.value = whenIf;
			if (!isIf.value) return;
		}
		if (isFunction(show)) isShow.value = !!await show(formValues, formApi);
		else if (isBoolean(show)) isShow.value = show;
		if (isFunction(componentProps)) dynamicComponentProps.value = await componentProps(formValues, formApi);
		if (isFunction(rules)) dynamicRules.value = await rules(formValues, formApi);
		if (isFunction(disabled)) isDisabled.value = !!await disabled(formValues, formApi);
		else if (isBoolean(disabled)) isDisabled.value = disabled;
		if (isFunction(required)) isRequired.value = !!await required(formValues, formApi);
		if (isFunction(trigger)) trigger(formValues, formApi);
	}, {
		deep: true,
		immediate: true
	});
	return {
		dynamicComponentProps,
		dynamicRules,
		isDisabled,
		isIf,
		isRequired,
		isShow
	};
}
//#endregion
export { useDependencies as default };
