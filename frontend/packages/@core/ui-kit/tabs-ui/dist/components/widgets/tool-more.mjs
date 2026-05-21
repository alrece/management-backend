import { createBlock, createElementVNode, createVNode, defineComponent, openBlock, unref, withCtx } from "vue";
import { LayoutGrid } from "@vben-core/icons";
import { VbenDropdownMenu } from "@vben-core/shadcn-ui";
//#region src/components/widgets/tool-more.vue
const _hoisted_1 = { class: "flex-center h-full cursor-pointer border-l border-border px-2 text-lg font-semibold text-muted-foreground hover:bg-muted hover:text-foreground" };
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "tool-more",
	props: { menus: {} },
	setup(__props) {
		return (_ctx, _cache) => {
			return openBlock(), createBlock(unref(VbenDropdownMenu), {
				menus: __props.menus,
				modal: false
			}, {
				default: withCtx(() => [createElementVNode("div", _hoisted_1, [createVNode(unref(LayoutGrid), { class: "size-4" })])]),
				_: 1
			}, 8, ["menus"]);
		};
	}
});
//#endregion
export { _sfc_main as default };
