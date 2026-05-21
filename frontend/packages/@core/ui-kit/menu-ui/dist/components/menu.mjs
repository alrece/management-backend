import { flattedChildren } from "../utils/index.mjs";
import { useMenuStyle } from "../hooks/use-menu.mjs";
import { createMenuContext, createSubMenuContext } from "../hooks/use-menu-context.mjs";
import { useMenuScroll } from "../hooks/use-menu-scroll.mjs";
import _sfc_main$1 from "./sub-menu.mjs";
/* empty css                                      */
import { Fragment, computed, createBlock, createElementBlock, createVNode, defineComponent, nextTick, normalizeClass, normalizeStyle, openBlock, reactive, ref, renderList, renderSlot, resolveDynamicComponent, toRef, unref, useSlots, watch, watchEffect, withCtx } from "vue";
import { useNamespace } from "@vben-core/composables";
import { Ellipsis } from "@vben-core/icons";
import { useResizeObserver } from "@vueuse/core";
//#region src/components/menu.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "Menu",
	__name: "menu",
	props: {
		accordion: {
			type: Boolean,
			default: true
		},
		collapse: {
			type: Boolean,
			default: false
		},
		collapseShowTitle: { type: Boolean },
		defaultActive: {},
		defaultOpeneds: {},
		mode: { default: "vertical" },
		rounded: {
			type: Boolean,
			default: true
		},
		scrollToActive: {
			type: Boolean,
			default: false
		},
		theme: { default: "dark" }
	},
	emits: [
		"close",
		"open",
		"select"
	],
	setup(__props, { emit: __emit }) {
		const props = __props;
		const emit = __emit;
		const { b, is } = useNamespace("menu");
		const menuStyle = useMenuStyle();
		const slots = useSlots();
		const menu = ref();
		const sliceIndex = ref(-1);
		const openedMenus = ref(props.defaultOpeneds && !props.collapse ? [...props.defaultOpeneds] : []);
		const activePath = ref(props.defaultActive);
		const items = ref({});
		const subMenus = ref({});
		const mouseInChild = ref(false);
		const isMenuPopup = computed(() => {
			return props.mode === "horizontal" || props.mode === "vertical" && props.collapse;
		});
		const getSlot = computed(() => {
			const originalSlot = flattedChildren(slots.default?.() ?? []);
			const slotDefault = sliceIndex.value === -1 ? originalSlot : originalSlot.slice(0, sliceIndex.value);
			const slotMore = sliceIndex.value === -1 ? [] : originalSlot.slice(sliceIndex.value);
			return {
				showSlotMore: slotMore.length > 0,
				slotDefault,
				slotMore
			};
		});
		watch(() => props.collapse, (value) => {
			if (value) openedMenus.value = [];
		});
		watch(items.value, initMenu);
		watch(() => props.defaultActive, (currentActive = "") => {
			if (!items.value[currentActive]) activePath.value = "";
			updateActiveName(currentActive);
		});
		let resizeStopper;
		watchEffect(() => {
			if (props.mode === "horizontal") resizeStopper = useResizeObserver(menu, handleResize).stop;
			else resizeStopper?.();
		});
		createMenuContext(reactive({
			activePath,
			addMenuItem,
			addSubMenu,
			closeMenu,
			handleMenuItemClick,
			handleSubMenuClick,
			isMenuPopup,
			openedMenus,
			openMenu,
			props,
			removeMenuItem,
			removeSubMenu,
			subMenus,
			theme: toRef(props, "theme"),
			items
		}));
		createSubMenuContext({
			addSubMenu,
			level: 1,
			mouseInChild,
			removeSubMenu
		});
		function calcMenuItemWidth(menuItem) {
			const computedStyle = getComputedStyle(menuItem);
			const marginLeft = Number.parseInt(computedStyle.marginLeft, 10);
			const marginRight = Number.parseInt(computedStyle.marginRight, 10);
			return menuItem.offsetWidth + marginLeft + marginRight || 0;
		}
		function calcSliceIndex() {
			if (!menu.value) return -1;
			const items = [...menu.value?.childNodes ?? []].filter((item) => item.nodeName !== "#comment" && (item.nodeName !== "#text" || item.nodeValue));
			const moreItemWidth = 46;
			const computedMenuStyle = getComputedStyle(menu?.value);
			const paddingLeft = Number.parseInt(computedMenuStyle.paddingLeft, 10);
			const paddingRight = Number.parseInt(computedMenuStyle.paddingRight, 10);
			const menuWidth = menu.value?.clientWidth - paddingLeft - paddingRight;
			let calcWidth = 0;
			let sliceIndex = 0;
			items.forEach((item, index) => {
				calcWidth += calcMenuItemWidth(item);
				if (calcWidth <= menuWidth - moreItemWidth) sliceIndex = index + 1;
			});
			return sliceIndex === items.length ? -1 : sliceIndex;
		}
		function debounce(fn, wait = 33.34) {
			let timer;
			return () => {
				timer && clearTimeout(timer);
				timer = setTimeout(() => {
					fn();
				}, wait);
			};
		}
		let isFirstTimeRender = true;
		function handleResize() {
			if (sliceIndex.value === calcSliceIndex()) return;
			const callback = () => {
				sliceIndex.value = -1;
				nextTick(() => {
					sliceIndex.value = calcSliceIndex();
				});
			};
			callback();
			isFirstTimeRender ? callback() : debounce(callback)();
			isFirstTimeRender = false;
		}
		useMenuScroll(activePath, {
			enable: computed(() => props.scrollToActive && props.mode === "vertical" && !props.collapse),
			delay: 320
		});
		function initMenu() {
			getActivePaths().forEach((path) => {
				const subMenu = subMenus.value[path];
				subMenu && openMenu(path, subMenu.parentPaths);
			});
		}
		function updateActiveName(val) {
			const itemsInData = items.value;
			const item = itemsInData[val] || activePath.value && itemsInData[activePath.value] || itemsInData[props.defaultActive || ""];
			activePath.value = item ? item.path : val;
		}
		function handleMenuItemClick(data) {
			const { collapse, mode } = props;
			if (mode === "horizontal" || collapse) openedMenus.value = [];
			const { parentPaths, path } = data;
			if (!path || !parentPaths) return;
			emit("select", path, parentPaths);
		}
		function handleSubMenuClick({ parentPaths, path }) {
			if (openedMenus.value.includes(path)) closeMenu(path, parentPaths);
			else openMenu(path, parentPaths);
		}
		function close(path) {
			const i = openedMenus.value.indexOf(path);
			if (i !== -1) openedMenus.value.splice(i, 1);
		}
		/**
		* 关闭、折叠菜单
		*/
		function closeMenu(path, parentPaths) {
			if (props.accordion) openedMenus.value = subMenus.value[path]?.parentPaths ?? [];
			close(path);
			emit("close", path, parentPaths);
		}
		/**
		* 点击展开菜单
		*/
		function openMenu(path, parentPaths) {
			if (openedMenus.value.includes(path)) return;
			if (props.accordion) {
				const activeParentPaths = getActivePaths();
				if (activeParentPaths.includes(path)) parentPaths = activeParentPaths;
				openedMenus.value = openedMenus.value.filter((path) => parentPaths.includes(path));
			}
			openedMenus.value.push(path);
			emit("open", path, parentPaths);
		}
		function addMenuItem(item) {
			items.value[item.path] = item;
		}
		function addSubMenu(subMenu) {
			subMenus.value[subMenu.path] = subMenu;
		}
		function removeSubMenu(subMenu) {
			Reflect.deleteProperty(subMenus.value, subMenu.path);
		}
		function removeMenuItem(item) {
			Reflect.deleteProperty(items.value, item.path);
		}
		function getActivePaths() {
			const activeItem = activePath.value && items.value[activePath.value];
			if (!activeItem || props.mode === "horizontal" || props.collapse) return [];
			return activeItem.parentPaths;
		}
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("ul", {
				ref_key: "menu",
				ref: menu,
				class: normalizeClass([
					__props.theme,
					unref(b)(),
					unref(is)(__props.mode, true),
					unref(is)(__props.theme, true),
					unref(is)("rounded", __props.rounded),
					unref(is)("collapse", __props.collapse),
					unref(is)("menu-align", __props.mode === "horizontal")
				]),
				style: normalizeStyle(unref(menuStyle)),
				role: "menu"
			}, [__props.mode === "horizontal" && getSlot.value.showSlotMore ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [(openBlock(true), createElementBlock(Fragment, null, renderList(getSlot.value.slotDefault, (item, index) => {
				return openBlock(), createBlock(resolveDynamicComponent(item), { key: index });
			}), 128)), createVNode(_sfc_main$1, {
				"is-sub-menu-more": "",
				path: "sub-menu-more"
			}, {
				title: withCtx(() => [createVNode(unref(Ellipsis), { class: "size-4" })]),
				default: withCtx(() => [(openBlock(true), createElementBlock(Fragment, null, renderList(getSlot.value.slotMore, (item, index) => {
					return openBlock(), createBlock(resolveDynamicComponent(item), { key: index });
				}), 128))]),
				_: 1
			})], 64)) : renderSlot(_ctx.$slots, "default", { key: 1 })], 6);
		};
	}
});
//#endregion
export { _sfc_main as default };
