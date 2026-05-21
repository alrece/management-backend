import { useMenu, useMenuStyle } from "../hooks/use-menu.mjs";
import { createSubMenuContext, useMenuContext, useSubMenuContext } from "../hooks/use-menu-context.mjs";
import _sfc_main$1 from "./collapse-transition.mjs";
import _sfc_main$2 from "./sub-menu-content.mjs";
import { Fragment, computed, createBlock, createElementBlock, createElementVNode, createVNode, defineComponent, normalizeClass, normalizeStyle, onBeforeUnmount, onMounted, openBlock, reactive, ref, renderSlot, unref, vShow, withCtx, withDirectives, withModifiers } from "vue";
import { useNamespace } from "@vben-core/composables";
import { VbenHoverCard } from "@vben-core/shadcn-ui";
//#region src/components/sub-menu.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "SubMenu",
	__name: "sub-menu",
	props: {
		isSubMenuMore: {
			type: Boolean,
			default: false
		},
		activeIcon: {},
		disabled: {
			type: Boolean,
			default: false
		},
		icon: {},
		path: {},
		badge: {},
		badgeType: {},
		badgeVariants: {}
	},
	setup(__props) {
		const props = __props;
		const { parentMenu, parentPaths } = useMenu();
		const { b, is } = useNamespace("sub-menu");
		const nsMenu = useNamespace("menu");
		const rootMenu = useMenuContext();
		const subMenu = useSubMenuContext();
		const subMenuStyle = useMenuStyle(subMenu);
		const mouseInChild = ref(false);
		const items = ref({});
		const subMenus = ref({});
		const timer = ref(null);
		createSubMenuContext({
			addSubMenu,
			handleMouseleave,
			level: (subMenu?.level ?? 0) + 1,
			mouseInChild,
			removeSubMenu
		});
		const opened = computed(() => {
			return rootMenu?.openedMenus.includes(props.path);
		});
		const isTopLevelMenuSubmenu = computed(() => parentMenu.value?.type.name === "Menu");
		const mode = computed(() => rootMenu?.props.mode ?? "vertical");
		const rounded = computed(() => rootMenu?.props.rounded);
		const currentLevel = computed(() => subMenu?.level ?? 0);
		const isFirstLevel = computed(() => {
			return currentLevel.value === 1;
		});
		const contentProps = computed(() => {
			const isHorizontal = mode.value === "horizontal";
			return {
				collisionPadding: { top: 20 },
				side: isHorizontal && isFirstLevel.value ? "bottom" : "right",
				sideOffset: isHorizontal ? 5 : 10
			};
		});
		const active = computed(() => {
			let isActive = false;
			Object.values(items.value).forEach((item) => {
				if (item.active) isActive = true;
			});
			Object.values(subMenus.value).forEach((subItem) => {
				if (subItem.active) isActive = true;
			});
			return isActive;
		});
		function addSubMenu(subMenu) {
			subMenus.value[subMenu.path] = subMenu;
		}
		function removeSubMenu(subMenu) {
			Reflect.deleteProperty(subMenus.value, subMenu.path);
		}
		/**
		* 点击submenu展开/关闭
		*/
		function handleClick() {
			const mode = rootMenu?.props.mode;
			if (props.disabled || rootMenu?.props.collapse && mode === "vertical" || mode === "horizontal") return;
			rootMenu?.handleSubMenuClick({
				active: active.value,
				parentPaths: parentPaths.value,
				path: props.path
			});
		}
		function handleMouseenter(event, showTimeout = 300) {
			if (event.type === "focus") return;
			if (!rootMenu?.props.collapse && rootMenu?.props.mode === "vertical" || props.disabled) {
				if (subMenu) subMenu.mouseInChild.value = true;
				return;
			}
			if (subMenu) subMenu.mouseInChild.value = true;
			timer.value && window.clearTimeout(timer.value);
			timer.value = setTimeout(() => {
				rootMenu?.openMenu(props.path, parentPaths.value);
			}, showTimeout);
			parentMenu.value?.vnode.el?.dispatchEvent(new MouseEvent("mouseenter"));
		}
		function handleMouseleave(deepDispatch = false) {
			if (!rootMenu?.props.collapse && rootMenu?.props.mode === "vertical" && subMenu) {
				subMenu.mouseInChild.value = false;
				return;
			}
			timer.value && window.clearTimeout(timer.value);
			if (subMenu) subMenu.mouseInChild.value = false;
			timer.value = setTimeout(() => {
				!mouseInChild.value && rootMenu?.closeMenu(props.path, parentPaths.value);
			}, 300);
			if (deepDispatch) subMenu?.handleMouseleave?.(true);
		}
		const menuIcon = computed(() => active.value ? props.activeIcon || props.icon : props.icon);
		const item = reactive({
			active,
			parentPaths,
			path: props.path
		});
		onMounted(() => {
			subMenu?.addSubMenu?.(item);
			rootMenu?.addSubMenu?.(item);
		});
		onBeforeUnmount(() => {
			subMenu?.removeSubMenu?.(item);
			rootMenu?.removeSubMenu?.(item);
		});
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("li", {
				class: normalizeClass([
					unref(b)(),
					unref(is)("opened", opened.value),
					unref(is)("active", active.value),
					unref(is)("disabled", __props.disabled)
				]),
				onFocus: handleMouseenter,
				onMouseenter: handleMouseenter,
				onMouseleave: _cache[3] || (_cache[3] = () => handleMouseleave())
			}, [unref(rootMenu).isMenuPopup ? (openBlock(), createBlock(unref(VbenHoverCard), {
				key: 0,
				"content-class": [
					unref(rootMenu).theme,
					unref(nsMenu).e("popup-container"),
					unref(is)(unref(rootMenu).theme, true),
					opened.value ? "" : "hidden",
					"overflow-auto",
					"max-h-[calc(var(--reka-hover-card-content-available-height)-20px)]",
					mode.value === "horizontal" ? "is-horizontal" : ""
				],
				"content-props": contentProps.value,
				open: true,
				"open-delay": 0
			}, {
				trigger: withCtx(() => [createVNode(_sfc_main$2, {
					class: normalizeClass(unref(is)("active", active.value)),
					icon: menuIcon.value,
					"is-menu-more": __props.isSubMenuMore,
					"is-top-level-menu-submenu": isTopLevelMenuSubmenu.value,
					level: currentLevel.value,
					path: __props.path,
					onClick: withModifiers(handleClick, ["stop"])
				}, {
					title: withCtx(() => [renderSlot(_ctx.$slots, "title")]),
					_: 3
				}, 8, [
					"class",
					"icon",
					"is-menu-more",
					"is-top-level-menu-submenu",
					"level",
					"path"
				])]),
				default: withCtx(() => [createElementVNode("div", {
					class: normalizeClass([unref(nsMenu).is(mode.value, true), unref(nsMenu).e("popup")]),
					onFocus: _cache[0] || (_cache[0] = (e) => handleMouseenter(e, 100)),
					onMouseenter: _cache[1] || (_cache[1] = (e) => handleMouseenter(e, 100)),
					onMouseleave: _cache[2] || (_cache[2] = () => handleMouseleave(true))
				}, [createElementVNode("ul", {
					class: normalizeClass([unref(nsMenu).b(), unref(is)("rounded", rounded.value)]),
					style: normalizeStyle(unref(subMenuStyle))
				}, [renderSlot(_ctx.$slots, "default")], 6)], 34)]),
				_: 3
			}, 8, ["content-class", "content-props"])) : (openBlock(), createElementBlock(Fragment, { key: 1 }, [createVNode(_sfc_main$2, {
				class: normalizeClass(unref(is)("active", active.value)),
				icon: menuIcon.value,
				"is-menu-more": __props.isSubMenuMore,
				"is-top-level-menu-submenu": isTopLevelMenuSubmenu.value,
				level: currentLevel.value,
				path: __props.path,
				onClick: withModifiers(handleClick, ["stop"])
			}, {
				title: withCtx(() => [renderSlot(_ctx.$slots, "title")]),
				default: withCtx(() => [renderSlot(_ctx.$slots, "content")]),
				_: 3
			}, 8, [
				"class",
				"icon",
				"is-menu-more",
				"is-top-level-menu-submenu",
				"level",
				"path"
			]), createVNode(_sfc_main$1, null, {
				default: withCtx(() => [withDirectives(createElementVNode("ul", {
					class: normalizeClass([unref(nsMenu).b(), unref(is)("rounded", rounded.value)]),
					style: normalizeStyle(unref(subMenuStyle))
				}, [renderSlot(_ctx.$slots, "default")], 6), [[vShow, opened.value]])]),
				_: 3
			})], 64))], 34);
		};
	}
});
//#endregion
export { _sfc_main as default };
