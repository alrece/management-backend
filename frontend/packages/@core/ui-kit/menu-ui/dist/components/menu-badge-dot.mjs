import { createElementBlock, createElementVNode, defineComponent, normalizeClass, normalizeStyle, openBlock } from "vue";
//#region src/components/menu-badge-dot.vue
const _hoisted_1 = { class: "relative mr-1 flex size-1.5" };
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "menu-badge-dot",
	props: {
		dotClass: { default: "" },
		dotStyle: { default: () => ({}) }
	},
	setup(__props) {
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("span", _hoisted_1, [createElementVNode("span", {
				class: normalizeClass([__props.dotClass, "absolute inline-flex size-full animate-ping rounded-full opacity-75"]),
				style: normalizeStyle(__props.dotStyle)
			}, null, 6), createElementVNode("span", {
				class: normalizeClass([__props.dotClass, "relative inline-flex size-1.5 rounded-full"]),
				style: normalizeStyle(__props.dotStyle)
			}, null, 6)]);
		};
	}
});
//#endregion
export { _sfc_main as default };
