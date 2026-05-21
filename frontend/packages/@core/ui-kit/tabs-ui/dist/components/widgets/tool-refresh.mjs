import { createElementBlock, createVNode, defineComponent, openBlock, unref } from "vue";
import { RotateCw } from "@vben-core/icons";
//#region src/components/widgets/tool-refresh.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "tool-refresh",
	emits: ["refresh"],
	setup(__props, { emit: __emit }) {
		const emit = __emit;
		const handleRefresh = () => {
			emit("refresh");
		};
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("div", {
				class: "flex-center h-full cursor-pointer border-l border-border px-2 text-lg font-semibold text-muted-foreground hover:bg-muted hover:text-foreground",
				onClick: handleRefresh
			}, [createVNode(unref(RotateCw), { class: "size-4" })]);
		};
	}
});
//#endregion
export { _sfc_main as default };
