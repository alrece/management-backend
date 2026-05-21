import { AlertProps } from "./alert.js";
import * as _$vue from "vue";

//#region src/alert/alert.vue.d.ts
type __VLS_Props = AlertProps;
type __VLS_ModelProps = {
  'open'?: boolean;
};
type __VLS_PublicProps = __VLS_Props & __VLS_ModelProps;
declare const __VLS_export: _$vue.DefineComponent<__VLS_PublicProps, {}, {}, {}, {}, _$vue.ComponentOptionsMixin, _$vue.ComponentOptionsMixin, {
  closed: (...args: any[]) => void;
  confirm: (...args: any[]) => void;
  opened: (...args: any[]) => void;
  "update:open": (value: boolean) => void;
}, string, _$vue.PublicProps, Readonly<__VLS_PublicProps> & Readonly<{
  "onUpdate:open"?: ((value: boolean) => any) | undefined;
  onClosed?: ((...args: any[]) => any) | undefined;
  onConfirm?: ((...args: any[]) => any) | undefined;
  onOpened?: ((...args: any[]) => any) | undefined;
}>, {
  bordered: boolean;
  buttonAlign: "center" | "end" | "start";
  centered: boolean;
}, {}, {}, {}, string, _$vue.ComponentProvideOptions, false, {}, any>;
declare const _default: typeof __VLS_export;
//#endregion
export { _default };