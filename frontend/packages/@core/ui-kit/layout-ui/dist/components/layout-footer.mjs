import { computed, createElementBlock, defineComponent, normalizeStyle, openBlock, renderSlot } from "vue";
//#region src/components/layout-footer.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "layout-footer",
	props: {
		fixed: { type: Boolean },
		height: {},
		show: {
			type: Boolean,
			default: true
		},
		width: {},
		zIndex: {}
	},
	setup(__props) {
		const props = __props;
		const style = computed(() => {
			const { fixed, height, show, width, zIndex } = props;
			return {
				height: `${height}px`,
				marginBottom: show ? "0" : `-${height}px`,
				position: fixed ? "fixed" : "static",
				width,
				zIndex
			};
		});
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("footer", {
				style: normalizeStyle(style.value),
				class: "bottom-0 w-full bg-background-deep transition-all duration-200"
			}, [renderSlot(_ctx.$slots, "default")], 4);
		};
	}
});
//#endregion
export { _sfc_main as default };
