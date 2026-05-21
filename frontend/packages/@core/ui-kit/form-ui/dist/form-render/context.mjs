import { computed } from "vue";
import { createContext } from "@vben-core/shadcn-ui";
//#region src/form-render/context.ts
const [injectRenderFormProps, provideFormRenderProps] = createContext("FormRenderProps");
const useFormContext = () => {
	const formRenderProps = injectRenderFormProps();
	const isVertical = computed(() => formRenderProps.layout === "vertical");
	const componentMap = computed(() => formRenderProps.componentMap);
	return {
		componentBindEventMap: computed(() => formRenderProps.componentBindEventMap),
		componentMap,
		isVertical
	};
};
//#endregion
export { injectRenderFormProps, provideFormRenderProps, useFormContext };
