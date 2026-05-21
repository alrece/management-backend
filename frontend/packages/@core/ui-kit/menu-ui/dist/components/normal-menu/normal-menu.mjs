/* empty css                                                             */
import export_helper_default from "../../_virtual/_/plugin-vue/export-helper.mjs";
import { Fragment, createElementBlock, createElementVNode, createVNode, defineComponent, normalizeClass, openBlock, renderList, toDisplayString, unref } from "vue";
import { useNamespace } from "@vben-core/composables";
import { VbenIcon } from "@vben-core/shadcn-ui";
//#region src/components/normal-menu/normal-menu.vue
const _hoisted_1 = ["onClick", "onMouseenter"];
var normal_menu_default = /* @__PURE__ */ export_helper_default(/* @__PURE__ */ defineComponent({
	name: "NormalMenu",
	__name: "normal-menu",
	props: {
		activePath: { default: "" },
		collapse: {
			type: Boolean,
			default: false
		},
		menus: { default: () => [] },
		rounded: { type: Boolean },
		theme: { default: "dark" }
	},
	emits: ["enter", "select"],
	setup(__props, { emit: __emit }) {
		const props = __props;
		const emit = __emit;
		const { b, e, is } = useNamespace("normal-menu");
		function menuIcon(menu) {
			return props.activePath === menu.path ? menu.activeIcon || menu.icon : menu.icon;
		}
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("ul", { class: normalizeClass([[
				__props.theme,
				unref(b)(),
				unref(is)("collapse", __props.collapse),
				unref(is)(__props.theme, true),
				unref(is)("rounded", __props.rounded)
			], "relative"]) }, [(openBlock(true), createElementBlock(Fragment, null, renderList(__props.menus, (menu) => {
				return openBlock(), createElementBlock("li", {
					key: menu.path,
					class: normalizeClass([unref(e)("item"), unref(is)("active", __props.activePath === menu.path)]),
					onClick: () => emit("select", menu),
					onMouseenter: () => emit("enter", menu)
				}, [createVNode(unref(VbenIcon), {
					class: normalizeClass(unref(e)("icon")),
					icon: menuIcon(menu),
					fallback: ""
				}, null, 8, ["class", "icon"]), createElementVNode("span", { class: normalizeClass([unref(e)("name"), "truncate"]) }, toDisplayString(menu.name), 3)], 42, _hoisted_1);
			}), 128))], 2);
		};
	}
}), [["__scopeId", "data-v-424fbe3f"]]);
//#endregion
export { normal_menu_default as default };
