import { DrawerApiOptions, DrawerProps, ExtendedDrawerApi } from "./drawer.js";
import * as _$vue from "vue";

//#region src/drawer/use-drawer.d.ts
declare function setDefaultDrawerProps(props: Partial<DrawerProps>): void;
declare function useVbenDrawer<TParentDrawerProps extends DrawerProps = DrawerProps>(options?: DrawerApiOptions): readonly [_$vue.DefineSetupFnComponent<TParentDrawerProps, {}, {}, TParentDrawerProps & {}, _$vue.PublicProps>, ExtendedDrawerApi] | readonly [_$vue.DefineSetupFnComponent<DrawerProps, {}, {}, DrawerProps & {}, _$vue.PublicProps>, ExtendedDrawerApi];
//#endregion
export { setDefaultDrawerProps, useVbenDrawer };