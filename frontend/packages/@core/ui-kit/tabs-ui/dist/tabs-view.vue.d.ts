import { TabsProps } from "./types.js";
import * as _$vue from "vue";
import * as _$_vben_core_typings0 from "@vben-core/typings";

//#region src/tabs-view.vue.d.ts
interface Props extends TabsProps {}
declare const __VLS_export: _$vue.DefineComponent<Props, {}, {}, {}, {}, _$vue.ComponentOptionsMixin, _$vue.ComponentOptionsMixin, {
  close: (args_0: string) => any;
  sortTabs: (args_0: number, args_1: number) => any;
  unpin: (args_0: _$_vben_core_typings0.TabDefinition) => any;
}, string, _$vue.PublicProps, Readonly<Props> & Readonly<{
  onClose?: ((args_0: string) => any) | undefined;
  onSortTabs?: ((args_0: number, args_1: number) => any) | undefined;
  onUnpin?: ((args_0: _$_vben_core_typings0.TabDefinition) => any) | undefined;
}>, {
  contentClass: string;
  draggable: boolean;
  styleType: _$_vben_core_typings0.TabsStyleType;
  wheelable: boolean;
}, {}, {}, {}, string, _$vue.ComponentProvideOptions, false, {}, any>;
declare const _default: typeof __VLS_export;
//#endregion
export { _default };