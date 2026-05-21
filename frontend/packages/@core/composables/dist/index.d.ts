import * as _$vue from "vue";
import { CSSProperties, ComputedRef, Ref } from "vue";
import { VisibleDomRect } from "@vben-core/shared/utils";
import { useEmitAsProps, useForwardExpose, useForwardProps, useForwardPropsEmits } from "reka-ui";
import Sortable, { SortableOptions } from "sortablejs";

//#region src/use-is-mobile.d.ts
declare function useIsMobile(): {
  isMobile: _$vue.ComputedRef<boolean>;
};
//#endregion
//#region src/use-layout-style.d.ts
declare function useLayoutContentStyle(): {
  contentElement: _$vue.Ref<HTMLDivElement | null, HTMLDivElement | null>;
  overlayStyle: _$vue.ComputedRef<CSSProperties>;
  visibleDomRect: _$vue.Ref<{
    bottom: number;
    height: number;
    left: number;
    right: number;
    top: number;
    width: number;
  } | null, VisibleDomRect | {
    bottom: number;
    height: number;
    left: number;
    right: number;
    top: number;
    width: number;
  } | null>;
};
declare function useLayoutHeaderStyle(): {
  getLayoutHeaderHeight: () => number;
  setLayoutHeaderHeight: (height: number) => void;
};
declare function useLayoutFooterStyle(): {
  getLayoutFooterHeight: () => number;
  setLayoutFooterHeight: (height: number) => void;
};
//#endregion
//#region src/use-namespace.d.ts
declare const useNamespace: (block: string) => {
  b: (blockSuffix?: string) => string;
  be: (blockSuffix?: string, element?: string) => string;
  bem: (blockSuffix?: string, element?: string, modifier?: string) => string;
  bm: (blockSuffix?: string, modifier?: string) => string;
  cssVar: (object: Record<string, string>) => Record<string, string>;
  cssVarBlock: (object: Record<string, string>) => Record<string, string>;
  cssVarBlockName: (name: string) => string;
  cssVarName: (name: string) => string;
  e: (element?: string) => string;
  em: (element?: string, modifier?: string) => string;
  is: {
    (name: string): string;
    (name: string, state: boolean | undefined): string;
  };
  m: (modifier?: string) => string;
  namespace: string;
};
type UseNamespaceReturn = ReturnType<typeof useNamespace>;
//#endregion
//#region src/use-priority-value.d.ts
declare function usePriorityValue<T extends Record<string, any>, S extends Record<string, any>, K extends keyof T = keyof T>(key: K, props: T, state: Readonly<Ref<NoInfer<S>>> | undefined): ComputedRef<T[K]>;
declare function usePriorityValues<T extends Record<string, any>, S extends Ref<Record<string, any>> = Readonly<Ref<NoInfer<T>, NoInfer<T>>>>(props: T, state: S | undefined): { [K in keyof T]: ComputedRef<T[K]> };
declare function useForwardPriorityValues<T extends Record<string, any>, S extends Ref<Record<string, any>> = Readonly<Ref<NoInfer<T>, NoInfer<T>>>>(props: T, state: S | undefined): ComputedRef<{ [K in keyof T]: T[K] }>;
//#endregion
//#region src/use-scroll-lock.d.ts
declare const SCROLL_FIXED_CLASS = "_scroll__fixed_";
declare function useScrollLock(): void;
//#endregion
//#region src/use-simple-locale/messages.d.ts
type Locale = 'en-US' | 'zh-CN';
//#endregion
//#region src/use-simple-locale/index.d.ts
declare const useSimpleLocale: () => {
  $t: _$vue.ComputedRef<(key: string) => string>;
  currentLocale: _$vue.Ref<Locale, Locale>;
  setSimpleLocale: (locale: Locale) => void;
};
//#endregion
//#region src/use-sortable.d.ts
declare function useSortable<T extends HTMLElement>(sortableContainer: T, options?: SortableOptions): {
  initializeSortable: () => Promise<Sortable>;
};
//#endregion
export { SCROLL_FIXED_CLASS, type Sortable, type UseNamespaceReturn, useEmitAsProps, useForwardExpose, useForwardPriorityValues, useForwardProps, useForwardPropsEmits, useIsMobile, useLayoutContentStyle, useLayoutFooterStyle, useLayoutHeaderStyle, useNamespace, usePriorityValue, usePriorityValues, useScrollLock, useSimpleLocale, useSortable };