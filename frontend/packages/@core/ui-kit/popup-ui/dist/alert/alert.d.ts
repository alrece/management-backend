import { Component, VNode, VNodeArrayChildren } from "vue";
import { Recordable } from "@vben-core/typings";

//#region src/alert/alert.d.ts
type IconType = 'error' | 'info' | 'question' | 'success' | 'warning';
type BeforeCloseScope = {
  isConfirm: boolean;
};
type AlertProps = {
  beforeClose?: (scope: BeforeCloseScope) => boolean | Promise<boolean | undefined> | undefined;
  bordered?: boolean;
  buttonAlign?: 'center' | 'end' | 'start';
  cancelText?: string;
  centered?: boolean;
  confirmText?: string;
  containerClass?: string;
  content: Component | string;
  contentClass?: string;
  contentMasking?: boolean;
  footer?: Component | string;
  icon?: Component | IconType;
  overlayBlur?: number;
  showCancel?: boolean;
  title?: string;
};
type PromptProps<T = any> = {
  beforeClose?: (scope: {
    isConfirm: boolean;
    value: T | undefined;
  }) => boolean | Promise<boolean | undefined> | undefined;
  component?: Component;
  componentProps?: Recordable<any>;
  componentSlots?: (() => any) | Recordable<unknown> | VNode | VNodeArrayChildren;
  defaultValue?: T;
  modelPropName?: string;
} & Omit<AlertProps, 'beforeClose'>;
type AlertContext = {
  doCancel: () => void;
  doConfirm: () => void;
};
declare function useAlertContext(): AlertContext;
//#endregion
export { AlertProps, BeforeCloseScope, IconType, PromptProps, useAlertContext };