import { ModalApiOptions, ModalState } from "./modal.js";
import { Store } from "@vben-core/shared/store";

//#region src/modal/modal-api.d.ts
declare class ModalApi {
  sharedData: Record<'payload', any>;
  store: Store<ModalState>;
  private api;
  private state;
  constructor(options?: ModalApiOptions);
  close(): Promise<void>;
  getData<T extends object = Record<string, any>>(): T;
  lock(isLocked?: boolean): this;
  onCancel(): void;
  onClosed(): void;
  onConfirm(): void;
  onOpened(): void;
  open(): void;
  setData<T>(payload: T): this;
  setState(stateOrFn: ((prev: ModalState) => Partial<ModalState>) | Partial<ModalState>): this;
  unlock(): this;
}
//#endregion
export { ModalApi };