import { MenuRecordRaw } from "@vben-core/typings";

//#region src/components/normal-menu/normal-menu.d.ts
interface NormalMenuProps {
  activePath?: string;
  collapse?: boolean;
  menus?: MenuRecordRaw[];
  rounded?: boolean;
  theme?: 'dark' | 'light';
}
//#endregion
export { type NormalMenuProps };