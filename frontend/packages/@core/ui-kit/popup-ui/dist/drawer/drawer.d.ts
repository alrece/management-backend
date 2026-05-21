import { DrawerApi } from "./drawer-api.js";
import { Component, Ref } from "vue";
import { ClassType, MaybePromise } from "@vben-core/typings";

//#region src/drawer/drawer.d.ts
type DrawerPlacement = 'bottom' | 'left' | 'right' | 'top';
type CloseIconPlacement = 'left' | 'right';
interface DrawerProps {
  appendToMain?: boolean;
  cancelText?: string;
  class?: ClassType;
  closable?: boolean;
  closeIconPlacement?: CloseIconPlacement;
  closeOnClickModal?: boolean;
  closeOnPressEscape?: boolean;
  confirmLoading?: boolean;
  confirmText?: string;
  contentClass?: string;
  description?: string;
  destroyOnClose?: boolean;
  footer?: boolean;
  footerClass?: ClassType;
  header?: boolean;
  headerClass?: ClassType;
  loading?: boolean;
  modal?: boolean;
  openAutoFocus?: boolean;
  overlayBlur?: number;
  placement?: DrawerPlacement;
  showCancelButton?: boolean;
  showConfirmButton?: boolean;
  submitting?: boolean;
  title?: string;
  titleTooltip?: string;
  zIndex?: number;
}
interface DrawerState extends DrawerProps {
  isOpen?: boolean;
  sharedData?: Record<string, any>;
}
type ExtendedDrawerApi = DrawerApi & {
  useStore: <T = NoInfer<DrawerState>>(selector?: (state: NoInfer<DrawerState>) => T) => Readonly<Ref<T>>;
};
interface DrawerApiOptions extends DrawerState {
  connectedComponent?: Component;
  onBeforeClose?: () => MaybePromise<boolean | undefined>;
  onCancel?: () => void;
  onClosed?: () => void;
  onConfirm?: () => void;
  onOpenChange?: (isOpen: boolean) => void;
  onOpened?: () => void;
}
//#endregion
export { CloseIconPlacement, DrawerApiOptions, DrawerPlacement, DrawerProps, DrawerState, ExtendedDrawerApi };