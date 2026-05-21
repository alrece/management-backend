import { computed, createElementBlock, createVNode, defineComponent, normalizeStyle, openBlock, renderSlot, unref, withCtx } from "vue";
import { useLayoutContentStyle } from "@vben-core/composables";
import { Slot } from "@vben-core/shadcn-ui";
//#region src/components/layout-content.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "layout-content",
	props: {
		contentCompact: {},
		contentCompactWidth: {},
		padding: {},
		paddingBottom: {},
		paddingLeft: {},
		paddingRight: {},
		paddingTop: {}
	},
	setup(__props) {
		const props = __props;
		const { contentElement, overlayStyle } = useLayoutContentStyle();
		const style = computed(() => {
			const { contentCompact, padding, paddingBottom, paddingLeft, paddingRight, paddingTop } = props;
			return {
				...contentCompact === "compact" ? {
					margin: "0 auto",
					width: `${props.contentCompactWidth}px`
				} : {},
				flex: 1,
				padding: `${padding}px`,
				paddingBottom: `${paddingBottom}px`,
				paddingLeft: `${paddingLeft}px`,
				paddingRight: `${paddingRight}px`,
				paddingTop: `${paddingTop}px`
			};
		});
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("main", {
				ref_key: "contentElement",
				ref: contentElement,
				style: normalizeStyle(style.value),
				class: "relative bg-background-deep"
			}, [createVNode(unref(Slot), { style: normalizeStyle(unref(overlayStyle)) }, {
				default: withCtx(() => [renderSlot(_ctx.$slots, "overlay")]),
				_: 3
			}, 8, ["style"]), renderSlot(_ctx.$slots, "default")], 4);
		};
	}
});
//#endregion
export { _sfc_main as default };
