import { BaseFormComponentType, FormActions, FormSchema, VbenFormProps } from "./types.js";
import { ComponentPublicInstance } from "vue";
import * as _$vee_validate0 from "vee-validate";
import { FormState, GenericObject, ResetFormOpts, ValidationOptions } from "vee-validate";
import { Store } from "@vben-core/shared/store";
import { StateHandler } from "@vben-core/shared/utils";
import { Recordable } from "@vben-core/typings";

//#region src/form-api.d.ts
declare class FormApi {
  form: FormActions;
  isMounted: boolean;
  state: null | VbenFormProps;
  stateHandler: StateHandler;
  store: Store<VbenFormProps>;
  private componentRefMap;
  private latestSubmissionValues;
  private prevState;
  constructor(options?: VbenFormProps);
  getFieldComponentRef<T = ComponentPublicInstance>(fieldName: string): T | undefined;
  getFocusedField(): string | undefined;
  getLatestSubmissionValues(): Recordable<any>;
  getState(): VbenFormProps<BaseFormComponentType, Record<never, never>> | null;
  getValues<T = Recordable<any>>(): Promise<T>;
  isFieldValid(fieldName: string): Promise<boolean>;
  merge(formApi: FormApi): any;
  mount(formActions: FormActions, componentRefMap?: Map<string, unknown>): void;
  removeSchemaByFields(fields: string[]): Promise<void>;
  resetForm(state?: Partial<FormState<GenericObject>> | undefined, opts?: Partial<ResetFormOpts>): Promise<void>;
  resetValidate(): Promise<void>;
  scrollToFirstError(errors: Record<string, any> | string): void;
  setDisabled(disabled: boolean): void;
  setFieldValue(field: string, value: any, shouldValidate?: boolean): Promise<void>;
  setLatestSubmissionValues(values: null | Recordable<any>): void;
  setLoading(loading: boolean): void;
  setState(stateOrFn: ((prev: VbenFormProps) => Partial<VbenFormProps>) | Partial<VbenFormProps>): void;
  setValues(fields: Record<string, any>, filterFields?: boolean, shouldValidate?: boolean): Promise<void>;
  submitForm(e?: Event): Promise<Recordable<any>>;
  unmount(): void;
  updateSchema(schema: Partial<FormSchema>[]): void;
  validate(opts?: Partial<ValidationOptions>): Promise<_$vee_validate0.FormValidationResult<GenericObject, GenericObject>>;
  validateAndSubmitForm(): Promise<Recordable<any> | undefined>;
  validateField(fieldName: string, opts?: Partial<ValidationOptions>): Promise<_$vee_validate0.ValidationResult<any>>;
  private deleteValueByFieldName;
  private getForm;
  private handleMultiFields;
  private handleRangeTimeValue;
  private handleValueFormat;
  private processFields;
  private resolveValueByFieldName;
  private setValueByFieldName;
  private updateState;
}
//#endregion
export { FormApi };