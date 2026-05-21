import { AlertProps, PromptProps } from "./alert.js";

//#region src/alert/AlertBuilder.d.ts
declare function vbenAlert(options: AlertProps): Promise<void>;
declare function vbenAlert(message: string, options?: Partial<AlertProps>): Promise<void>;
declare function vbenAlert(message: string, title?: string, options?: Partial<AlertProps>): Promise<void>;
declare function vbenConfirm(options: AlertProps): Promise<void>;
declare function vbenConfirm(message: string, options?: Partial<AlertProps>): Promise<void>;
declare function vbenConfirm(message: string, title?: string, options?: Partial<AlertProps>): Promise<void>;
declare function vbenPrompt<T = any>(options: PromptProps<T>): Promise<T | undefined>;
declare function clearAllAlerts(): void;
//#endregion
export { clearAllAlerts, vbenAlert, vbenConfirm, vbenPrompt };