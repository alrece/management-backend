import { ModalApi } from "./modal-api.js";
import { Component, Ref } from "vue";
import { ClassType, MaybePromise } from "@vben-core/typings";

//#region src/modal/modal.d.ts
interface ModalProps {
  animationType?: 'scale' | 'slide';
  appendToMain?: boolean;
  bordered?: boolean;
  cancelText?: string;
  centered?: boolean;
  class?: ClassType;
  closable?: boolean;
  closeOnClickModal?: boolean;
  closeOnPressEscape?: boolean;
  confirmDisabled?: boolean;
  confirmLoading?: boolean;
  confirmText?: string;
  contentClass?: ClassType;
  description?: string;
  destroyOnClose?: boolean;
  draggable?: boolean;
  footer?: boolean;
  footerClass?: ClassType;
  fullscreen?: boolean;
  fullscreenButton?: boolean;
  header?: boolean;
  headerClass?: ClassType;
  loading?: boolean;
  modal?: boolean;
  openAutoFocus?: boolean;
  overflow?: boolean;
  overlayBlur?: number;
  showCancelButton?: boolean;
  showConfirmButton?: boolean;
  submitting?: boolean;
  title?: string;
  titleTooltip?: string;
  zIndex?: number;
}
interface ModalState extends ModalProps {
  isOpen?: boolean;
  sharedData?: Record<string, any>;
}
type ExtendedModalApi = ModalApi & {
  useStore: <T = NoInfer<ModalState>>(selector?: (state: NoInfer<ModalState>) => T) => Readonly<Ref<T>>;
};
interface ModalApiOptions extends ModalState {
  connectedComponent?: Component;
  onBeforeClose?: () => MaybePromise<boolean | undefined>;
  onCancel?: () => void;
  onClosed?: () => void;
  onConfirm?: () => void;
  onOpenChange?: (isOpen: boolean) => void;
  onOpened?: () => void;
}
//#endregion
export { ExtendedModalApi, ModalApiOptions, ModalProps, ModalState };