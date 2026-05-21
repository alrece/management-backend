import { createContext } from "@vben-core/shadcn-ui";
//#region src/alert/alert.ts
const [injectAlertContext, provideAlertContext] = createContext("VbenAlertContext");
/**
* 获取Alert上下文
* @returns AlertContext
*/
function useAlertContext() {
	const context = injectAlertContext();
	if (!context) throw new Error("useAlertContext must be used within an AlertProvider");
	return context;
}
//#endregion
export { provideAlertContext, useAlertContext };
