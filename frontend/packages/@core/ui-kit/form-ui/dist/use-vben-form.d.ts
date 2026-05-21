import { BaseFormComponentType, ExtendedFormApi, VbenFormProps } from "./types.js";
import * as _$vue from "vue";

//#region src/use-vben-form.d.ts
declare function useVbenForm<T extends BaseFormComponentType = BaseFormComponentType, P extends Record<string, any> = Record<never, never>>(options: VbenFormProps<T, P>): readonly [_$vue.DefineSetupFnComponent<VbenFormProps<BaseFormComponentType, Record<never, never>>, {}, {}, VbenFormProps<BaseFormComponentType, Record<never, never>> & {}, _$vue.PublicProps>, ExtendedFormApi];
//#endregion
export { useVbenForm };