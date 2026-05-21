import { FormApi } from "./form-api.js";
import { Component, HtmlHTMLAttributes, Ref } from "vue";
import { VbenButtonProps } from "@vben-core/shadcn-ui";
import { FieldOptions, FormContext, GenericObject } from "vee-validate";
import { ZodTypeAny } from "zod";
import { ClassType, MaybeComputedRef } from "@vben-core/typings";

//#region src/types.d.ts
type FormLayout = 'horizontal' | 'inline' | 'vertical';
type BaseFormComponentType = 'DefaultButton' | 'PrimaryButton' | 'VbenCheckbox' | 'VbenInput' | 'VbenInputPassword' | 'VbenPinInput' | 'VbenSelect' | (Record<never, never> & string);
type Breakpoints = '2xl:' | '3xl:' | '' | 'lg:' | 'md:' | 'sm:' | 'xl:';
type GridCols = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13;
type WrapperClassType = `${Breakpoints}grid-cols-${GridCols}` | (Record<never, never> & string);
type FormFieldOptions = Partial<FieldOptions & {
  validateOnBlur?: boolean;
  validateOnChange?: boolean;
  validateOnInput?: boolean;
  validateOnModelUpdate?: boolean;
}>;
type MaybeComponentPropKey = 'options' | 'placeholder' | 'title' | keyof HtmlHTMLAttributes | (Record<never, never> & string);
type MaybeComponentProps = { [K in MaybeComponentPropKey]?: any };
type FormActions = FormContext<GenericObject>;
type CustomRenderType = (() => Component | string) | string;
type CustomParamsRenderType = ((value: Partial<Record<string, any>>, actions: FormActions) => Component | string) | string;
type FormSchemaRuleType = 'mobile' | 'mobileRequired' | 'required' | 'selectRequired' | null | (Record<never, never> & string) | ZodTypeAny;
type FormItemDependenciesCondition<T = boolean | PromiseLike<boolean>> = (value: Partial<Record<string, any>>, actions: FormActions) => T;
type FormItemDependenciesConditionWithRules = (value: Partial<Record<string, any>>, actions: FormActions) => FormSchemaRuleType | PromiseLike<FormSchemaRuleType>;
type FormItemDependenciesConditionWithProps = (value: Partial<Record<string, any>>, actions: FormActions) => MaybeComponentProps | PromiseLike<MaybeComponentProps>;
interface FormItemDependencies {
  componentProps?: FormItemDependenciesConditionWithProps;
  disabled?: boolean | FormItemDependenciesCondition;
  if?: boolean | FormItemDependenciesCondition;
  required?: FormItemDependenciesCondition;
  rules?: FormItemDependenciesConditionWithRules;
  show?: boolean | FormItemDependenciesCondition;
  trigger?: FormItemDependenciesCondition<void>;
  triggerFields: string[];
}
type ComponentProps = ((value: Partial<Record<string, any>>, actions: FormActions) => MaybeComponentProps) | MaybeComponentProps;
interface FormCommonConfig {
  colon?: boolean;
  componentProps?: ComponentProps;
  controlClass?: string;
  disabled?: boolean;
  disabledOnChangeListener?: boolean;
  disabledOnInputListener?: boolean;
  emptyStateValue?: null | undefined;
  formFieldProps?: FormFieldOptions;
  formItemClass?: (() => string) | string;
  hideLabel?: boolean;
  hideRequiredMark?: boolean;
  labelClass?: string;
  labelWidth?: number;
  modelPropName?: string;
  wrapperClass?: string;
}
type RenderComponentContentType = (value: Partial<Record<string, any>>, api: FormActions) => Record<string, any>;
type MappedComponentProps<P> = ((value: Partial<Record<string, any>>, actions: FormActions) => P & Record<string, any>) | (P & Record<string, any>);
type FormValueFormat = (value: any, setValue: (fieldName: string, value: any) => void, values: Record<string, any>) => any;
interface FormSchemaBody extends Omit<FormCommonConfig, 'componentProps'> {
  defaultValue?: any;
  dependencies?: FormItemDependencies;
  description?: CustomRenderType;
  fieldName: string;
  help?: CustomParamsRenderType;
  hide?: boolean;
  label?: CustomRenderType;
  renderComponentContent?: RenderComponentContentType;
  rules?: FormSchemaRuleType;
  suffix?: CustomRenderType;
  valueFormat?: FormValueFormat;
}
type FormSchemaDiscriminated<T extends BaseFormComponentType, P extends Record<string, any>> = { [K in Extract<keyof P, T>]: {
  component: K;
  componentProps?: MappedComponentProps<P[K]>;
} & FormSchemaBody }[Extract<keyof P, T>];
type FormSchemaFallback<T extends BaseFormComponentType> = {
  component: Component | T;
  componentProps?: ComponentProps;
} & FormSchemaBody;
type FormSchema<T extends BaseFormComponentType = BaseFormComponentType, P extends Record<string, any> = Record<never, never>> = FormSchemaDiscriminated<T, P> | FormSchemaFallback<T>;
type HandleSubmitFn = (values: Record<string, any>) => Promise<void> | void;
type HandleResetFn = (values: Record<string, any>) => Promise<void> | void;
type FieldMappingTime = [string, [string, string], (((value: any, fieldName: string) => any) | [string, string] | null | string)?][];
type ArrayToStringFields = Array<[string[], string?] | string | string[]>;
interface FormRenderProps<T extends BaseFormComponentType = BaseFormComponentType, P extends Record<string, any> = Record<never, never>> {
  arrayToStringFields?: ArrayToStringFields;
  collapsed?: boolean;
  collapsedRows?: number;
  collapseTriggerResize?: boolean;
  commonConfig?: FormCommonConfig;
  compact?: boolean;
  componentBindEventMap?: Partial<Record<BaseFormComponentType, string>>;
  componentMap: Record<BaseFormComponentType, Component>;
  fieldMappingTime?: FieldMappingTime;
  form?: FormContext<GenericObject>;
  layout?: FormLayout;
  schema?: FormSchema<T, P>[];
  showCollapseButton?: boolean;
  wrapperClass?: WrapperClassType;
}
interface ActionButtonOptions extends VbenButtonProps {
  [key: string]: any;
  content?: MaybeComputedRef<string>;
  show?: boolean;
}
interface VbenFormProps<T extends BaseFormComponentType = BaseFormComponentType, P extends Record<string, any> = Record<never, never>> extends Omit<FormRenderProps<T, P>, 'componentBindEventMap' | 'componentMap' | 'form'> {
  actionButtonsReverse?: boolean;
  actionLayout?: 'inline' | 'newLine' | 'rowEnd';
  actionPosition?: 'center' | 'left' | 'right';
  actionWrapperClass?: ClassType;
  arrayToStringFields?: ArrayToStringFields;
  fieldMappingTime?: FieldMappingTime;
  handleCollapsedChange?: (collapsed: boolean) => void;
  handleReset?: HandleResetFn;
  handleSubmit?: HandleSubmitFn;
  handleValuesChange?: (values: Record<string, any>, fieldsChanged: string[]) => void;
  resetButtonOptions?: ActionButtonOptions;
  scrollToFirstError?: boolean;
  showDefaultActions?: boolean;
  submitButtonOptions?: ActionButtonOptions;
  submitOnChange?: boolean;
  submitOnEnter?: boolean;
}
type ExtendedFormApi = FormApi & {
  useStore: <T = NoInfer<VbenFormProps>>(selector?: (state: NoInfer<VbenFormProps>) => T) => Readonly<Ref<T>>;
};
interface VbenFormAdapterOptions<T extends BaseFormComponentType = BaseFormComponentType> {
  config?: {
    baseModelPropName?: string;
    disabledOnChangeListener?: boolean;
    disabledOnInputListener?: boolean;
    emptyStateValue?: null | undefined;
    modelPropNameMap?: Partial<Record<T, string>>;
  };
  defineRules?: {
    mobile?: (value: any, params: any, ctx: Record<string, any>) => boolean | string;
    mobileRequired?: (value: any, params: any, ctx: Record<string, any>) => boolean | string;
    required?: (value: any, params: any, ctx: Record<string, any>) => boolean | string;
    selectRequired?: (value: any, params: any, ctx: Record<string, any>) => boolean | string;
  };
}
//#endregion
export { BaseFormComponentType, ExtendedFormApi, FormActions, FormSchema, VbenFormAdapterOptions, VbenFormProps };