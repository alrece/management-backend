import { Component, ComputedRef, MaybeRef } from "vue";
import { RouteLocationNormalized, RouteMeta as RouteMeta$1, RouteRecordRaw, RouteRecordRaw as RouteRecordRaw$1, Router } from "vue-router";

//#region src/app.d.ts
type LayoutType = 'full-content' | 'header-mixed-nav' | 'header-nav' | 'header-sidebar-nav' | 'mixed-nav' | 'sidebar-mixed-nav' | 'sidebar-nav';
type ThemeModeType = 'auto' | 'dark' | 'light';
/**
 * 偏好设置按钮位置
 * fixed 固定在右侧
 * header 顶栏
 * auto 自动
 */
type PreferencesButtonPositionType = 'auto' | 'fixed' | 'header';
type BuiltinThemeType = 'custom' | 'deep-blue' | 'deep-green' | 'default' | 'gray' | 'green' | 'neutral' | 'orange' | 'pink' | 'red' | 'rose' | 'sky-blue' | 'slate' | 'stone' | 'violet' | 'yellow' | 'zinc' | (Record<never, never> & string);
type ContentCompactType = 'compact' | 'wide';
type LayoutHeaderModeType = 'auto' | 'auto-scroll' | 'fixed' | 'static';
type LayoutHeaderMenuAlignType = 'center' | 'end' | 'start';
/**
 * 登录过期模式
 * modal 弹窗模式
 * page 页面模式
 */
type LoginExpiredModeType = 'modal' | 'page';
/**
 * 面包屑样式
 * background 背景
 * normal 默认
 */
type BreadcrumbStyleType = 'background' | 'normal';
/**
 * 权限模式
 * backend 后端权限模式
 * frontend 前端权限模式
 * mixed 混合权限模式
 */
type AccessModeType = 'backend' | 'frontend' | 'mixed';
/**
 * 导航风格
 * plain 朴素
 * rounded 圆润
 */
type NavigationStyleType = 'plain' | 'rounded';
/**
 * 标签栏风格
 * brisk 轻快
 * card 卡片
 * chrome 谷歌
 * plain 朴素
 */
type TabsStyleType = 'brisk' | 'card' | 'chrome' | 'plain';
/**
 * 页面切换动画
 */
type PageTransitionType = 'fade' | 'fade-down' | 'fade-slide' | 'fade-up';
/**
 * 页面切换动画
 * panel-center 居中布局
 * panel-left 居左布局
 * panel-right 居右布局
 */
type AuthPageLayoutType = 'panel-center' | 'panel-left' | 'panel-right';
/**
 * 时区选项
 */
interface TimezoneOption {
  label: string;
  offset: number;
  timezone: string;
}
//#endregion
//#region src/basic.d.ts
interface BasicOption {
  label: string;
  value: string;
}
type SelectOption = BasicOption;
type TabOption = BasicOption;
interface BasicUserInfo {
  /**
   * 头像
   */
  avatar: string;
  /**
   * 用户昵称
   */
  nickname: string;
  /**
   * 用户角色
   */
  roles?: string[];
  /**
   * 用户id
   */
  userId: string;
  /**
   * 用户名
   */
  username: string;
}
type ClassType = Array<ClassType> | boolean | null | object | string | undefined;
//#endregion
//#region src/helper.d.ts
/**
 * 类型级递归中增加深度计数
 */
type Increment<A extends unknown[]> = [...A, unknown];
/**
 * 深层递归所有属性为可选
 */
type DeepPartial<T, D extends number = 10, C extends unknown[] = []> = C['length'] extends D ? T : T extends object ? { [P in keyof T]?: DeepPartial<T[P], D, Increment<C>> } : T;
/**
 * 深层递归所有属性为只读
 */
type DeepReadonly<T, D extends number = 10, C extends unknown[] = []> = C['length'] extends D ? T : T extends object ? { readonly [P in keyof T]: DeepReadonly<T[P], D, Increment<C>> } : T;
/**
 * 任意类型的异步函数
 */
type AnyPromiseFunction<T extends any[] = any[], R = void> = (...arg: T) => PromiseLike<R>;
/**
 * 任意类型的普通函数
 */
type AnyNormalFunction<T extends any[] = any[], R = void> = (...arg: T) => R;
/**
 * 任意类型的函数
 */
type AnyFunction<T extends any[] = any[], R = void> = AnyNormalFunction<T, R> | AnyPromiseFunction<T, R>;
/**
 *  T | null 包装
 */
type Nullable<T> = null | T;
/**
 * T | Not null 包装
 */
type NonNullable<T> = T extends null | undefined ? never : T;
/**
 * 字符串类型对象
 */
type Recordable$1<T> = Record<string, T>;
/**
 * 字符串类型对象（只读）
 */
interface ReadonlyRecordable<T = any> {
  readonly [key: string]: T;
}
/**
 * setTimeout 返回值类型
 */
type TimeoutHandle = ReturnType<typeof setTimeout>;
/**
 * setInterval 返回值类型
 */
type IntervalHandle = ReturnType<typeof setInterval>;
/**
 * 也许它是一个计算的 ref，或者一个 getter 函数
 *
 */
type MaybeReadonlyRef<T> = (() => T) | ComputedRef<T>;
/**
 * 也许它是一个 ref，或者一个普通值，或者一个 getter 函数
 *
 */
type MaybeComputedRef<T> = MaybeReadonlyRef<T> | MaybeRef<T>;
type Merge<O extends object, T extends object> = { [K in keyof O | keyof T]: K extends keyof T ? T[K] : K extends keyof O ? O[K] : never };
/**
 * T = [
 *  { name: string; age: number; },
 *  { sex: 'male' | 'female'; age: string }
 * ]
 * =>
 * MergeAll<T> = {
 *  name: string;
 *  sex: 'male' | 'female';
 *  age: string
 * }
 */
type MergeAll<T extends object[], R extends object = Record<string, any>> = T extends [infer F extends object, ...infer Rest extends object[]] ? MergeAll<Rest, Merge<R, F>> : R;
type EmitType = (name: Name, ...args: any[]) => void;
type MaybePromise<T> = Promise<T> | T;
//#endregion
//#region src/menu-record.d.ts
interface AppRouteRecordRaw extends Omit<RouteRecordRaw$1, 'meta'> {
  children?: AppRouteRecordRaw[];
  component?: any;
  componentName?: string;
  components?: any;
  fullPath?: string;
  icon?: string;
  id?: any;
  keepAlive?: boolean;
  meta: RouteMeta$1;
  name: string;
  parentId?: number;
  props?: any;
  sort?: number;
  visible?: boolean;
}
type ExRouteRecordRaw = RouteRecordRaw$1 & {
  parent?: string;
  parents?: string[];
  path?: any;
};
interface MenuRecordBadgeRaw {
  badge?: string;
  badgeType?: 'dot' | 'normal';
  badgeVariants?: 'destructive' | 'primary' | string;
}
interface MenuRecordRaw extends MenuRecordBadgeRaw {
  activeIcon?: string;
  children?: MenuRecordRaw[];
  disabled?: boolean;
  icon?: Component | string;
  name: string;
  order?: number;
  parent?: string;
  parents?: string[];
  path: string;
  query?: Recordable$1<any>;
  show?: boolean;
}
//#endregion
//#region src/tabs.d.ts
interface TabDefinition extends RouteLocationNormalized {
  key?: string;
}
//#endregion
//#region src/vue-router.d.ts
interface RouteMeta {
  /**
   * 激活图标（菜单/tab）
   */
  activeIcon?: string;
  /**
   * 当前激活的菜单，有时候不想激活现有菜单，需要激活父级菜单时使用
   */
  activePath?: string;
  /**
   * 是否固定标签页
   * @default false
   */
  affixTab?: boolean;
  /**
   * 固定标签页的顺序
   * @default 0
   */
  affixTabOrder?: number;
  /**
   * 需要特定的角色标识才可以访问
   * @default []
   */
  authority?: string[];
  /**
   * 徽标
   */
  badge?: string;
  /**
   * 徽标类型
   */
  badgeType?: 'dot' | 'normal';
  /**
   * 徽标颜色
   */
  badgeVariants?: 'default' | 'destructive' | 'primary' | 'success' | 'warning' | string;
  /**
   * 路由对应dom是否缓存起来
   */
  domCached?: boolean;
  /**
   * 路由的完整路径作为key（默认true）
   */
  fullPathKey?: boolean;
  /**
   * 当前路由的子级在菜单中不展现
   * @default false
   */
  hideChildrenInMenu?: boolean;
  /**
   * 当前路由在面包屑中不展现
   * @default false
   */
  hideInBreadcrumb?: boolean;
  /**
   * 当前路由在菜单中不展现
   * @default false
   */
  hideInMenu?: boolean;
  /**
   * 当前路由在标签页不展现
   * @default false
   */
  hideInTab?: boolean;
  /**
   * 图标（菜单/tab）
   */
  icon?: Component | string;
  /**
   * iframe 地址
   */
  iframeSrc?: string;
  /**
   * 忽略权限，直接可以访问
   * @default false
   */
  ignoreAccess?: boolean;
  /**
   * 开启KeepAlive缓存
   */
  keepAlive?: boolean;
  /**
   * 外链-跳转路径
   */
  link?: string;
  /**
   * 路由是否已经加载过
   */
  loaded?: boolean;
  /**
   * 标签页最大打开数量
   * @default -1
   */
  maxNumOfOpenTab?: number;
  /**
   * 菜单可以看到，但是访问会被重定向到403
   */
  menuVisibleWithForbidden?: boolean;
  /**
   * 不使用基础布局（仅在顶级生效）
   */
  noBasicLayout?: boolean;
  /**
   * 在新窗口打开
   */
  openInNewWindow?: boolean;
  /**
   * 用于路由->菜单排序
   */
  order?: number;
  /**
   * 菜单所携带的参数
   */
  query?: Recordable;
  /**
   * 标题名称
   */
  title: string;
} // 定义递归类型以将 RouteRecordRaw 的 component 属性更改为 string
type RouteRecordStringComponent<T = string> = Omit<RouteRecordRaw, 'children' | 'component'> & {
  children?: RouteRecordStringComponent<T>[];
  component: T;
};
type ComponentRecordType = Record<string, () => Promise<Component>>;
interface GenerateMenuAndRoutesOptions {
  fetchMenuListAsync?: () => Promise<RouteRecordStringComponent[]>;
  forbiddenComponent?: RouteRecordRaw['component'];
  layoutMap?: ComponentRecordType;
  pageMap?: ComponentRecordType;
  roles?: string[];
  router: Router;
  routes: RouteRecordRaw[];
}
//#endregion
export { type AccessModeType, type AnyFunction, type AnyNormalFunction, type AnyPromiseFunction, type AppRouteRecordRaw, type AuthPageLayoutType, type BasicOption, type BasicUserInfo, type BreadcrumbStyleType, type BuiltinThemeType, type ClassType, type ComponentRecordType, type ContentCompactType, type DeepPartial, type DeepReadonly, type EmitType, type ExRouteRecordRaw, type GenerateMenuAndRoutesOptions, type IntervalHandle, type LayoutHeaderMenuAlignType, type LayoutHeaderModeType, type LayoutType, type LoginExpiredModeType, type MaybeComputedRef, type MaybePromise, type MaybeReadonlyRef, type MenuRecordBadgeRaw, type MenuRecordRaw, type Merge, type MergeAll, type NavigationStyleType, type NonNullable, type Nullable, type PageTransitionType, type PreferencesButtonPositionType, type ReadonlyRecordable, type Recordable$1 as Recordable, type RouteMeta, type RouteRecordRaw, type RouteRecordStringComponent, type SelectOption, TabDefinition, type TabOption, type TabsStyleType, type ThemeModeType, type TimeoutHandle, type TimezoneOption };