import _sfc_main$1 from "./components/menu-badge.mjs";
import _sfc_main$2 from "./components/menu-item.mjs";
import _sfc_main$3 from "./components/sub-menu.mjs";
import { Fragment, computed, createBlock, createElementBlock, createElementVNode, createVNode, defineComponent, openBlock, renderList, toDisplayString, unref, withCtx } from "vue";
//#region src/sub-menu.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "SubMenuUi",
	__name: "sub-menu",
	props: { menu: {} },
	setup(__props) {
		const props = __props;
		/**
		* 判断是否有子节点，动态渲染 menu-item/sub-menu-item
		*/
		const hasChildren = computed(() => {
			const { menu } = props;
			return Reflect.has(menu, "children") && !!menu.children && menu.children.length > 0;
		});
		return (_ctx, _cache) => {
			return !hasChildren.value ? (openBlock(), createBlock(unref(_sfc_main$2), {
				key: __props.menu.path,
				"active-icon": __props.menu.activeIcon,
				badge: __props.menu.badge,
				"badge-type": __props.menu.badgeType,
				"badge-variants": __props.menu.badgeVariants,
				icon: __props.menu.icon,
				path: __props.menu.path,
				query: __props.menu.query
			}, {
				title: withCtx(() => [createElementVNode("span", null, toDisplayString(__props.menu.name), 1)]),
				_: 1
			}, 8, [
				"active-icon",
				"badge",
				"badge-type",
				"badge-variants",
				"icon",
				"path",
				"query"
			])) : (openBlock(), createBlock(unref(_sfc_main$3), {
				key: `${__props.menu.path}_sub`,
				"active-icon": __props.menu.activeIcon,
				icon: __props.menu.icon,
				path: __props.menu.path
			}, {
				content: withCtx(() => [createVNode(unref(_sfc_main$1), {
					badge: __props.menu.badge,
					"badge-type": __props.menu.badgeType,
					"badge-variants": __props.menu.badgeVariants,
					class: "right-6"
				}, null, 8, [
					"badge",
					"badge-type",
					"badge-variants"
				])]),
				title: withCtx(() => [createElementVNode("span", null, toDisplayString(__props.menu.name), 1)]),
				default: withCtx(() => [(openBlock(true), createElementBlock(Fragment, null, renderList(__props.menu.children || [], (childItem) => {
					return openBlock(), createBlock(_sfc_main, {
						key: childItem.path,
						menu: childItem
					}, null, 8, ["menu"]);
				}), 128))]),
				_: 1
			}, 8, [
				"active-icon",
				"icon",
				"path"
			]));
		};
	}
});
//#endregion
export { _sfc_main as default };
