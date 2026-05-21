import _sfc_main$1 from "./components/menu.mjs";
import _sfc_main$2 from "./sub-menu.mjs";
import { Fragment, createBlock, createElementBlock, defineComponent, guardReactiveProps, normalizeProps, openBlock, renderList, unref, withCtx } from "vue";
import { useForwardProps } from "@vben-core/composables";
//#region src/menu.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "MenuView",
	__name: "menu",
	props: {
		menus: {},
		accordion: { type: Boolean },
		collapse: {
			type: Boolean,
			default: false
		},
		collapseShowTitle: { type: Boolean },
		defaultActive: {},
		defaultOpeneds: {},
		mode: {},
		rounded: { type: Boolean },
		scrollToActive: { type: Boolean },
		theme: {}
	},
	setup(__props) {
		const forward = useForwardProps(__props);
		return (_ctx, _cache) => {
			return openBlock(), createBlock(unref(_sfc_main$1), normalizeProps(guardReactiveProps(unref(forward))), {
				default: withCtx(() => [(openBlock(true), createElementBlock(Fragment, null, renderList(__props.menus, (menu) => {
					return openBlock(), createBlock(_sfc_main$2, {
						key: menu.path,
						menu
					}, null, 8, ["menu"]);
				}), 128))]),
				_: 1
			}, 16);
		};
	}
});
//#endregion
export { _sfc_main as default };
