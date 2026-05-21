import { useMenuContext } from "../hooks/use-menu-context.mjs";
import { computed, createBlock, createCommentVNode, createElementBlock, defineComponent, normalizeClass, normalizeStyle, openBlock, renderSlot, resolveDynamicComponent, unref, vShow, withDirectives } from "vue";
import { useNamespace } from "@vben-core/composables";
import { VbenIcon } from "@vben-core/shadcn-ui";
import { ChevronDown, ChevronRight } from "@vben-core/icons";
//#region src/components/sub-menu-content.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "SubMenuContent",
	__name: "sub-menu-content",
	props: {
		isMenuMore: {
			type: Boolean,
			default: false
		},
		isTopLevelMenuSubmenu: { type: Boolean },
		level: { default: 0 },
		activeIcon: {},
		disabled: { type: Boolean },
		icon: {},
		path: {},
		query: {},
		badge: {},
		badgeType: {},
		badgeVariants: {}
	},
	setup(__props) {
		const props = __props;
		const rootMenu = useMenuContext();
		const { b, e, is } = useNamespace("sub-menu-content");
		const nsMenu = useNamespace("menu");
		const opened = computed(() => {
			return rootMenu?.openedMenus.includes(props.path);
		});
		const collapse = computed(() => {
			return rootMenu.props.collapse;
		});
		const isFirstLevel = computed(() => {
			return props.level === 1;
		});
		const getCollapseShowTitle = computed(() => {
			return rootMenu.props.collapseShowTitle && isFirstLevel.value && collapse.value;
		});
		const mode = computed(() => {
			return rootMenu?.props.mode;
		});
		const showArrowIcon = computed(() => {
			return mode.value === "horizontal" || !(isFirstLevel.value && collapse.value);
		});
		const hiddenTitle = computed(() => {
			return mode.value === "vertical" && isFirstLevel.value && collapse.value && !getCollapseShowTitle.value;
		});
		const iconComp = computed(() => {
			return mode.value === "horizontal" && !isFirstLevel.value || mode.value === "vertical" && collapse.value ? ChevronRight : ChevronDown;
		});
		const iconArrowStyle = computed(() => {
			return opened.value ? { transform: `rotate(180deg)` } : {};
		});
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("div", { class: normalizeClass([
				unref(b)(),
				unref(is)("collapse-show-title", getCollapseShowTitle.value),
				unref(is)("more", __props.isMenuMore)
			]) }, [
				renderSlot(_ctx.$slots, "default"),
				!__props.isMenuMore ? (openBlock(), createBlock(unref(VbenIcon), {
					key: 0,
					class: normalizeClass(unref(nsMenu).e("icon")),
					icon: __props.icon,
					fallback: ""
				}, null, 8, ["class", "icon"])) : createCommentVNode("v-if", true),
				!hiddenTitle.value ? (openBlock(), createElementBlock("div", {
					key: 1,
					class: normalizeClass([unref(e)("title")])
				}, [renderSlot(_ctx.$slots, "title")], 2)) : createCommentVNode("v-if", true),
				!__props.isMenuMore ? withDirectives((openBlock(), createBlock(resolveDynamicComponent(iconComp.value), {
					key: 2,
					class: normalizeClass([[unref(e)("icon-arrow")], "size-4"]),
					style: normalizeStyle(iconArrowStyle.value)
				}, null, 8, ["class", "style"])), [[vShow, showArrowIcon.value]]) : createCommentVNode("v-if", true)
			], 2);
		};
	}
});
//#endregion
export { _sfc_main as default };
