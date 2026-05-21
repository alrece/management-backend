import { MenuProps } from "./types.js";
import * as _$vue from "vue";
import { MenuRecordRaw } from "@vben-core/typings";

//#region src/menu.vue.d.ts
interface Props extends MenuProps {
  menus: MenuRecordRaw[];
}
declare const __VLS_export: _$vue.DefineComponent<Props, {}, {}, {}, {}, _$vue.ComponentOptionsMixin, _$vue.ComponentOptionsMixin, {}, string, _$vue.PublicProps, Readonly<Props> & Readonly<{}>, {
  collapse: boolean;
}, {}, {}, {}, string, _$vue.ComponentProvideOptions, false, {}, any>;
declare const _default: typeof __VLS_export;
//#endregion
export { _default };