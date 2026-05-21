import { IContextMenuItem } from "@vben-core/shadcn-ui";
import { TabDefinition, TabsStyleType } from "@vben-core/typings";

//#region src/types.d.ts
interface TabsProps {
  active?: string;
  contentClass?: string;
  contextMenus?: (data: any) => IContextMenuItem[];
  draggable?: boolean;
  gap?: number;
  maxWidth?: number;
  middleClickToClose?: boolean;
  minWidth?: number;
  showIcon?: boolean;
  styleType?: TabsStyleType;
  tabs?: TabDefinition[];
  wheelable?: boolean;
}
//#endregion
export { TabsProps };