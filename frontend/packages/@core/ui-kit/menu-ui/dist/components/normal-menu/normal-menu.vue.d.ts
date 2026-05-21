import { NormalMenuProps } from "./normal-menu.js";
import * as _$vue from "vue";
import { MenuRecordRaw } from "@vben-core/typings";

//#region src/components/normal-menu/normal-menu.vue.d.ts
interface Props extends NormalMenuProps {}
declare const __VLS_export: _$vue.DefineComponent<Props, {}, {}, {}, {}, _$vue.ComponentOptionsMixin, _$vue.ComponentOptionsMixin, {
  enter: (args_0: MenuRecordRaw) => any;
  select: (args_0: MenuRecordRaw) => any;
}, string, _$vue.PublicProps, Readonly<Props> & Readonly<{
  onEnter?: ((args_0: MenuRecordRaw) => any) | undefined;
  onSelect?: ((args_0: MenuRecordRaw) => any) | undefined;
}>, {
  menus: MenuRecordRaw[];
  collapse: boolean;
  theme: "dark" | "light";
  activePath: string;
}, {}, {}, {}, string, _$vue.ComponentProvideOptions, false, {}, any>;
declare const _default: typeof __VLS_export;
//#endregion
export { _default };