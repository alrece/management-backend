import { computed, createCommentVNode, createElementBlock, defineComponent, normalizeClass, normalizeStyle, openBlock, renderSlot, unref, useSlots } from "vue";
//#region src/components/layout-header.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "layout-header",
	props: {
		fullWidth: { type: Boolean },
		height: {},
		isMobile: { type: Boolean },
		show: { type: Boolean },
		sidebarWidth: {},
		theme: {},
		width: {},
		zIndex: {}
	},
	setup(__props) {
		const props = __props;
		const slots = useSlots();
		const style = computed(() => {
			const { fullWidth, height, show } = props;
			const right = !show || !fullWidth ? void 0 : 0;
			return {
				height: `${height}px`,
				marginTop: show ? 0 : `-${height}px`,
				right
			};
		});
		const logoStyle = computed(() => {
			return { minWidth: `${props.isMobile ? 40 : props.sidebarWidth}px` };
		});
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("header", {
				class: normalizeClass([__props.theme, "top-0 flex w-full flex-[0_0_auto] items-center border-b border-border bg-header pl-2 transition-[margin-top] duration-200"]),
				style: normalizeStyle(style.value)
			}, [
				unref(slots).logo ? (openBlock(), createElementBlock("div", {
					key: 0,
					style: normalizeStyle(logoStyle.value)
				}, [renderSlot(_ctx.$slots, "logo")], 4)) : createCommentVNode("v-if", true),
				renderSlot(_ctx.$slots, "toggle-button"),
				renderSlot(_ctx.$slots, "default")
			], 6);
		};
	}
});
//#endregion
export { _sfc_main as default };
