import tabs_default from "./components/tabs-chrome/tabs.mjs";
import _sfc_main$1 from "./components/tabs/tabs.mjs";
import { useTabsDrag } from "./use-tabs-drag.mjs";
import { useTabsViewScroll } from "./use-tabs-view-scroll.mjs";
import { createBlock, createCommentVNode, createElementBlock, createElementVNode, createVNode, defineComponent, mergeProps, normalizeClass, normalizeProps, openBlock, unref, vShow, withCtx, withDirectives } from "vue";
import { ChevronsLeft, ChevronsRight } from "@vben-core/icons";
import { VbenScrollbar } from "@vben-core/shadcn-ui";
import { useForwardPropsEmits } from "@vben-core/composables";
//#region src/tabs-view.vue
const _hoisted_1 = { class: "flex h-full flex-1 overflow-hidden" };
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "TabsView",
	__name: "tabs-view",
	props: {
		active: {},
		contentClass: { default: "vben-tabs-content" },
		contextMenus: {},
		draggable: {
			type: Boolean,
			default: true
		},
		gap: {},
		maxWidth: {},
		middleClickToClose: { type: Boolean },
		minWidth: {},
		showIcon: { type: Boolean },
		styleType: { default: "chrome" },
		tabs: {},
		wheelable: {
			type: Boolean,
			default: true
		}
	},
	emits: [
		"close",
		"sortTabs",
		"unpin"
	],
	setup(__props, { emit: __emit }) {
		const props = __props;
		const emit = __emit;
		const forward = useForwardPropsEmits(props, emit);
		const { handleScrollAt, handleWheel, scrollbarRef, scrollDirection, scrollIsAtLeft, scrollIsAtRight, showScrollButton } = useTabsViewScroll(props);
		function onWheel(e) {
			if (props.wheelable) {
				handleWheel(e);
				e.stopPropagation();
				e.preventDefault();
			}
		}
		useTabsDrag(props, emit);
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("div", _hoisted_1, [
				createCommentVNode(" 左侧滚动按钮 "),
				withDirectives(createElementVNode("span", {
					class: normalizeClass([{
						"cursor-pointer text-muted-foreground hover:bg-muted": !unref(scrollIsAtLeft),
						"pointer-events-none opacity-30": unref(scrollIsAtLeft)
					}, "border-r px-2"]),
					onClick: _cache[0] || (_cache[0] = ($event) => unref(scrollDirection)("left"))
				}, [createVNode(unref(ChevronsLeft), { class: "size-4 h-full" })], 2), [[vShow, unref(showScrollButton)]]),
				createElementVNode("div", { class: normalizeClass([{ "pt-0.75": __props.styleType === "chrome" }, "size-full flex-1 overflow-hidden"]) }, [createVNode(unref(VbenScrollbar), {
					ref_key: "scrollbarRef",
					ref: scrollbarRef,
					"shadow-bottom": false,
					"shadow-top": false,
					class: "h-full",
					horizontal: "",
					"scroll-bar-class": "z-10 hidden ",
					shadow: "",
					"shadow-left": "",
					"shadow-right": "",
					onScrollAt: unref(handleScrollAt),
					onWheel
				}, {
					default: withCtx(() => [__props.styleType === "chrome" ? (openBlock(), createBlock(unref(tabs_default), normalizeProps(mergeProps({ key: 0 }, {
						...unref(forward),
						..._ctx.$attrs,
						..._ctx.$props
					})), null, 16)) : (openBlock(), createBlock(unref(_sfc_main$1), normalizeProps(mergeProps({ key: 1 }, {
						...unref(forward),
						..._ctx.$attrs,
						..._ctx.$props
					})), null, 16))]),
					_: 1
				}, 8, ["onScrollAt"])], 2),
				createCommentVNode(" 右侧滚动按钮 "),
				withDirectives(createElementVNode("span", {
					class: normalizeClass([{
						"cursor-pointer text-muted-foreground hover:bg-muted": !unref(scrollIsAtRight),
						"pointer-events-none opacity-30": unref(scrollIsAtRight)
					}, "cursor-pointer border-l px-2 text-muted-foreground hover:bg-muted"]),
					onClick: _cache[1] || (_cache[1] = ($event) => unref(scrollDirection)("right"))
				}, [createVNode(unref(ChevronsRight), { class: "size-4 h-full" })], 2), [[vShow, unref(showScrollButton)]])
			]);
		};
	}
});
//#endregion
export { _sfc_main as default };
