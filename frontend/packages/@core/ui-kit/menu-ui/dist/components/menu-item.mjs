import _sfc_main$1 from "./menu-badge.mjs";
import { useMenu } from "../hooks/use-menu.mjs";
import { useMenuContext, useSubMenuContext } from "../hooks/use-menu-context.mjs";
import { computed, createBlock, createCommentVNode, createElementBlock, createElementVNode, createVNode, defineComponent, mergeProps, normalizeClass, onBeforeUnmount, onMounted, openBlock, reactive, renderSlot, resolveComponent, unref, useSlots, vShow, withCtx, withDirectives, withModifiers } from "vue";
import { useNamespace } from "@vben-core/composables";
import { VbenIcon, VbenTooltip } from "@vben-core/shadcn-ui";
import { isHttpUrl } from "@vben-core/shared/utils";
import qs from "qs";
//#region src/components/menu-item.vue
const _hoisted_1 = ["href"];
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "MenuItem",
	__name: "menu-item",
	props: {
		activeIcon: {},
		disabled: {
			type: Boolean,
			default: false
		},
		icon: {},
		path: {},
		query: {},
		badge: {},
		badgeType: {},
		badgeVariants: {}
	},
	emits: ["click"],
	setup(__props, { emit: __emit }) {
		const props = __props;
		const emit = __emit;
		const slots = useSlots();
		const { b, e, is } = useNamespace("menu-item");
		const nsMenu = useNamespace("menu");
		const rootMenu = useMenuContext();
		const subMenu = useSubMenuContext();
		const { parentMenu, parentPaths } = useMenu();
		const active = computed(() => props.path === rootMenu?.activePath);
		const menuIcon = computed(() => active.value ? props.activeIcon || props.icon : props.icon);
		const isHttp = computed(() => isHttpUrl(item.parentPaths.at(-1)));
		const isTopLevelMenuItem = computed(() => parentMenu.value?.type.name === "Menu");
		const collapseShowTitle = computed(() => rootMenu.props?.collapseShowTitle && isTopLevelMenuItem.value && rootMenu.props.collapse);
		const showTooltip = computed(() => rootMenu.props.mode === "vertical" && isTopLevelMenuItem.value && rootMenu.props?.collapse && slots.title);
		const item = reactive({
			active,
			parentPaths: parentPaths.value,
			path: props.path || "",
			query: props.query
		});
		/**
		* 菜单项点击事件
		*/
		function handleClick() {
			if (props.disabled) return;
			rootMenu?.handleMenuItemClick?.({
				parentPaths: parentPaths.value,
				path: props.path
			});
			emit("click", item);
		}
		onMounted(() => {
			subMenu?.addSubMenu?.(item);
			rootMenu?.addMenuItem?.(item);
		});
		onBeforeUnmount(() => {
			subMenu?.removeSubMenu?.(item);
			rootMenu?.removeMenuItem?.(item);
		});
		return (_ctx, _cache) => {
			const _component_router_link = resolveComponent("router-link");
			return openBlock(), createBlock(_component_router_link, {
				custom: "",
				to: (item.parentPaths.at(-1) ?? "") + (item?.query ? `?${unref(qs).stringify(item?.query)}` : "")
			}, {
				default: withCtx(({ href }) => [createElementVNode("a", {
					href: isHttp.value ? item.parentPaths.at(-1) : href,
					class: normalizeClass([
						unref(rootMenu).theme,
						unref(b)(),
						unref(is)("active", active.value),
						unref(is)("disabled", __props.disabled),
						unref(is)("collapse-show-title", collapseShowTitle.value)
					]),
					role: "menuitem",
					onClick: withModifiers(handleClick, ["prevent", "stop"])
				}, [showTooltip.value ? (openBlock(), createBlock(unref(VbenTooltip), {
					key: 0,
					"content-class": [unref(rootMenu).theme],
					side: "right"
				}, {
					trigger: withCtx(() => [createElementVNode("div", { class: normalizeClass([unref(nsMenu).be("tooltip", "trigger")]) }, [
						createVNode(unref(VbenIcon), {
							class: normalizeClass(unref(nsMenu).e("icon")),
							icon: menuIcon.value,
							fallback: ""
						}, null, 8, ["class", "icon"]),
						renderSlot(_ctx.$slots, "default"),
						collapseShowTitle.value ? (openBlock(), createElementBlock("span", {
							key: 0,
							class: normalizeClass(unref(nsMenu).e("name"))
						}, [renderSlot(_ctx.$slots, "title")], 2)) : createCommentVNode("v-if", true)
					], 2)]),
					default: withCtx(() => [renderSlot(_ctx.$slots, "title")]),
					_: 3
				}, 8, ["content-class"])) : createCommentVNode("v-if", true), withDirectives(createElementVNode("div", { class: normalizeClass([unref(e)("content")]) }, [
					unref(rootMenu).props.mode !== "horizontal" ? (openBlock(), createBlock(unref(_sfc_main$1), mergeProps({
						key: 0,
						class: "right-2"
					}, props), null, 16)) : createCommentVNode("v-if", true),
					createVNode(unref(VbenIcon), {
						class: normalizeClass(unref(nsMenu).e("icon")),
						icon: menuIcon.value
					}, null, 8, ["class", "icon"]),
					renderSlot(_ctx.$slots, "default"),
					renderSlot(_ctx.$slots, "title")
				], 2), [[vShow, !showTooltip.value]])], 10, _hoisted_1)]),
				_: 3
			}, 8, ["to"]);
		};
	}
});
//#endregion
export { _sfc_main as default };
