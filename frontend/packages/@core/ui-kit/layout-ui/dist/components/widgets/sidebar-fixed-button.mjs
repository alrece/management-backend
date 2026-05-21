import { createBlock, createElementBlock, defineComponent, openBlock, unref, useModel } from "vue";
import { Pin, PinOff } from "@vben-core/icons";
//#region src/components/widgets/sidebar-fixed-button.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "sidebar-fixed-button",
	props: {
		"expandOnHover": { type: Boolean },
		"expandOnHoverModifiers": {}
	},
	emits: ["update:expandOnHover"],
	setup(__props) {
		const expandOnHover = useModel(__props, "expandOnHover");
		function toggleFixed() {
			expandOnHover.value = !expandOnHover.value;
		}
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("div", {
				class: "absolute right-3 bottom-2 z-10 flex-center cursor-pointer rounded-sm bg-accent p-1.25 text-foreground/60 transition-all duration-300 hover:bg-accent-hover hover:text-foreground",
				onClick: toggleFixed
			}, [!expandOnHover.value ? (openBlock(), createBlock(unref(PinOff), {
				key: 0,
				class: "size-3.5"
			})) : (openBlock(), createBlock(unref(Pin), {
				key: 1,
				class: "size-3.5"
			}))]);
		};
	}
});
//#endregion
export { _sfc_main as default };
