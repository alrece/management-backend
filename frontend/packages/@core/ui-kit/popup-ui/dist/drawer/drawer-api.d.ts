import { DrawerApiOptions, DrawerState } from "./drawer.js";
import { Store } from "@vben-core/shared/store";

//#region src/drawer/drawer-api.d.ts
declare class DrawerApi {
  sharedData: Record<'payload', any>;
  store: Store<DrawerState>;
  private api;
  private state;
  constructor(options?: DrawerApiOptions);
  close(): Promise<void>;
  getData<T extends object = Record<string, any>>(): T;
  lock(isLocked?: boolean): this;
  onCancel(): void;
  onClosed(): void;
  onConfirm(): void;
  onOpened(): void;
  open(): void;
  setData<T>(payload: T): this;
  setState(stateOrFn: ((prev: DrawerState) => Partial<DrawerState>) | Partial<DrawerState>): this;
  unlock(): this;
}
//#endregion
export { DrawerApi };