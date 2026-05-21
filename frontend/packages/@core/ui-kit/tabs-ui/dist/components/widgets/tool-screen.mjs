import { createBlock, createElementBlock, defineComponent, openBlock, unref, useModel } from "vue";
import { Fullscreen, Minimize2 } from "@vben-core/icons";
//#region src/components/widgets/tool-screen.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "tool-screen",
	props: {
		"screen": { type: Boolean },
		"screenModifiers": {}
	},
	emits: ["update:screen"],
	setup(__props) {
		const screen = useModel(__props, "screen");
		function toggleScreen() {
			screen.value = !screen.value;
		}
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("div", {
				class: "flex-center h-full cursor-pointer border-l border-border px-2 text-lg font-semibold text-muted-foreground hover:bg-muted hover:text-foreground",
				onClick: toggleScreen
			}, [screen.value ? (openBlock(), createBlock(unref(Minimize2), {
				key: 0,
				class: "size-4"
			})) : (openBlock(), createBlock(unref(Fullscreen), {
				key: 1,
				class: "size-4"
			}))]);
		};
	}
});
//#endregion
export { _sfc_main as default };
