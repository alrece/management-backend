import { ContentCompactType, LayoutHeaderModeType, LayoutType, ThemeModeType } from "@vben-core/typings";

//#region src/vben-layout.d.ts
interface VbenLayoutProps {
  contentCompact?: ContentCompactType;
  contentCompactWidth?: number;
  contentPadding?: number;
  contentPaddingBottom?: number;
  contentPaddingLeft?: number;
  contentPaddingRight?: number;
  contentPaddingTop?: number;
  footerEnable?: boolean;
  footerFixed?: boolean;
  footerHeight?: number;
  headerHeight?: number;
  headerHidden?: boolean;
  headerMode?: LayoutHeaderModeType;
  headerTheme?: ThemeModeType;
  headerToggleSidebarButton?: boolean;
  headerVisible?: boolean;
  isMobile?: boolean;
  layout?: LayoutType;
  sidebarCollapse?: boolean;
  sidebarCollapsedButton?: boolean;
  sidebarCollapseShowTitle?: boolean;
  sidebarEnable?: boolean;
  sidebarExtraCollapsedWidth?: number;
  sidebarFixedButton?: boolean;
  sidebarHidden?: boolean;
  sidebarMixedWidth?: number;
  sidebarTheme?: ThemeModeType;
  sidebarThemeSub?: ThemeModeType;
  sidebarWidth?: number;
  sideCollapseWidth?: number;
  tabbarEnable?: boolean;
  tabbarHeight?: number;
  zIndex?: number;
}
//#endregion
export { type VbenLayoutProps };