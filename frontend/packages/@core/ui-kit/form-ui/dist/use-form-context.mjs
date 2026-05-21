import { computed, unref, useSlots } from "vue";
import { createContext } from "@vben-core/shadcn-ui";
import { useForm } from "vee-validate";
import { isString, mergeWithArrayOverride, set } from "@vben-core/shared/utils";
import { ZodIntersection, ZodNumber, ZodObject, ZodString, object } from "zod";
import { getDefaultsForSchema } from "zod-defaults";
//#region src/use-form-context.ts
const [injectFormProps, provideFormProps] = createContext("VbenFormProps");
const [injectComponentRefMap, provideComponentRefMap] = createContext("ComponentRefMap");
function useFormInitial(props) {
	const slots = useSlots();
	const initialValues = generateInitialValues();
	const form = useForm({ ...Object.keys(initialValues)?.length ? { initialValues } : {} });
	const delegatedSlots = computed(() => {
		const resultSlots = [];
		for (const key of Object.keys(slots)) if (key !== "default") resultSlots.push(key);
		return resultSlots;
	});
	function generateInitialValues() {
		const initialValues = {};
		const zodObject = {};
		(unref(props).schema || []).forEach((item) => {
			if (Reflect.has(item, "defaultValue")) set(initialValues, item.fieldName, item.defaultValue);
			else if (item.rules && !isString(item.rules)) {
				const customDefaultValue = getCustomDefaultValue(item.rules);
				zodObject[item.fieldName] = item.rules;
				if (customDefaultValue !== void 0) initialValues[item.fieldName] = customDefaultValue;
			}
		});
		const schemaInitialValues = getDefaultsForSchema(object(zodObject));
		const zodDefaults = {};
		for (const key in schemaInitialValues) set(zodDefaults, key, schemaInitialValues[key]);
		return mergeWithArrayOverride(initialValues, zodDefaults);
	}
	function getCustomDefaultValue(rule) {
		if (rule instanceof ZodString) return "";
		else if (rule instanceof ZodNumber) return null;
		else if (rule instanceof ZodObject) {
			const defaultValues = {};
			for (const [key, valueSchema] of Object.entries(rule.shape)) defaultValues[key] = getCustomDefaultValue(valueSchema);
			return defaultValues;
		} else if (rule instanceof ZodIntersection) {
			const leftDefaultValue = getCustomDefaultValue(rule._def.left);
			const rightDefaultValue = getCustomDefaultValue(rule._def.right);
			if (typeof leftDefaultValue === "object" && typeof rightDefaultValue === "object") return {
				...leftDefaultValue,
				...rightDefaultValue
			};
			return leftDefaultValue ?? rightDefaultValue;
		} else return;
	}
	return {
		delegatedSlots,
		form
	};
}
//#endregion
export { injectComponentRefMap, injectFormProps, provideComponentRefMap, provideFormProps, useFormInitial };
