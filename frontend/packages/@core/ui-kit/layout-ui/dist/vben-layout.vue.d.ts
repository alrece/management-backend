import { VbenLayoutProps } from "./vben-layout.js";
import * as _$vue from "vue";
import * as _$_vben_core_typings0 from "@vben-core/typings";

//#region src/vben-layout.vue.d.ts
interface Props extends VbenLayoutProps {}
type __VLS_Props = Props;
type __VLS_ModelProps = {
  'sidebarDraggable'?: boolean;
  'sidebarCollapse'?: boolean;
  'sidebarExtraVisible'?: boolean;
  'sidebarExtraCollapse'?: boolean;
  'sidebarExpandOnHover'?: boolean;
  'sidebarEnable'?: boolean;
};
type __VLS_PublicProps = __VLS_Props & __VLS_ModelProps;
declare var __VLS_11: {}, __VLS_13: {}, __VLS_15: {}, __VLS_18: {}, __VLS_21: {}, __VLS_30: {}, __VLS_51: {}, __VLS_59: {}, __VLS_67: {}, __VLS_70: {}, __VLS_78: {}, __VLS_80: {};
type __VLS_Slots = {} & {
  logo?: (props: typeof __VLS_11) => any;
} & {
  'mixed-menu'?: (props: typeof __VLS_13) => any;
} & {
  menu?: (props: typeof __VLS_15) => any;
} & {
  'side-extra'?: (props: typeof __VLS_18) => any;
} & {
  'side-extra-title'?: (props: typeof __VLS_21) => any;
} & {
  logo?: (props: typeof __VLS_30) => any;
} & {
  header?: (props: typeof __VLS_51) => any;
} & {
  tabbar?: (props: typeof __VLS_59) => any;
} & {
  content?: (props: typeof __VLS_67) => any;
} & {
  'content-overlay'?: (props: typeof __VLS_70) => any;
} & {
  footer?: (props: typeof __VLS_78) => any;
} & {
  extra?: (props: typeof __VLS_80) => any;
};
declare const __VLS_base: _$vue.DefineComponent<__VLS_PublicProps, {}, {}, {}, {}, _$vue.ComponentOptionsMixin, _$vue.ComponentOptionsMixin, {
  "update:sidebarDraggable": (value: boolean) => any;
  "update:sidebarCollapse": (value: boolean) => any;
  "update:sidebarExtraVisible": (value: boolean | undefined) => any;
  "update:sidebarExtraCollapse": (value: boolean) => any;
  "update:sidebarExpandOnHover": (value: boolean) => any;
  "update:sidebarEnable": (value: boolean) => any;
  sideMouseLeave: () => any;
  toggleSidebar: () => any;
  "update:sidebar-width": (value: number) => any;
}, string, _$vue.PublicProps, Readonly<__VLS_PublicProps> & Readonly<{
  "onUpdate:sidebarDraggable"?: ((value: boolean) => any) | undefined;
  "onUpdate:sidebarCollapse"?: ((value: boolean) => any) | undefined;
  "onUpdate:sidebarExtraVisible"?: ((value: boolean | undefined) => any) | undefined;
  "onUpdate:sidebarExtraCollapse"?: ((value: boolean) => any) | undefined;
  "onUpdate:sidebarExpandOnHover"?: ((value: boolean) => any) | undefined;
  "onUpdate:sidebarEnable"?: ((value: boolean) => any) | undefined;
  onSideMouseLeave?: (() => any) | undefined;
  onToggleSidebar?: (() => any) | undefined;
  "onUpdate:sidebar-width"?: ((value: number) => any) | undefined;
}>, {
  contentCompact: _$_vben_core_typings0.ContentCompactType;
  contentCompactWidth: number;
  contentPadding: number;
  contentPaddingBottom: number;
  contentPaddingLeft: number;
  contentPaddingRight: number;
  contentPaddingTop: number;
  footerEnable: boolean;
  footerFixed: boolean;
  footerHeight: number;
  headerHeight: number;
  headerHidden: boolean;
  headerMode: _$_vben_core_typings0.LayoutHeaderModeType;
  headerToggleSidebarButton: boolean;
  headerVisible: boolean;
  isMobile: boolean;
  layout: _$_vben_core_typings0.LayoutType;
  sidebarCollapsedButton: boolean;
  sidebarCollapseShowTitle: boolean;
  sidebarExtraCollapsedWidth: number;
  sidebarFixedButton: boolean;
  sidebarHidden: boolean;
  sidebarMixedWidth: number;
  sidebarTheme: _$_vben_core_typings0.ThemeModeType;
  sidebarThemeSub: _$_vben_core_typings0.ThemeModeType;
  sidebarWidth: number;
  sideCollapseWidth: number;
  tabbarEnable: boolean;
  tabbarHeight: number;
  zIndex: number;
}, {}, {}, {}, string, _$vue.ComponentProvideOptions, false, {}, any>;
declare const __VLS_export: __VLS_WithSlots<typeof __VLS_base, __VLS_Slots>;
declare const _default: typeof __VLS_export;
type __VLS_WithSlots<T, S> = T & {
  new (): {
    $slots: S;
  };
};
//#endregion
export { _default };