import { computed, createElementBlock, defineComponent, normalizeStyle, openBlock, renderSlot } from "vue";
//#region src/components/layout-tabbar.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "layout-tabbar",
	props: { height: {} },
	setup(__props) {
		const props = __props;
		const style = computed(() => {
			const { height } = props;
			return { height: `${height}px` };
		});
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("section", {
				style: normalizeStyle(style.value),
				class: "flex w-full border-b border-border bg-background transition-all"
			}, [renderSlot(_ctx.$slots, "default")], 4);
		};
	}
});
//#endregion
export { _sfc_main as default };
