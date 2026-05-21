import _sfc_main$1 from "./components/layout-content.mjs";
import _sfc_main$2 from "./components/layout-footer.mjs";
import _sfc_main$3 from "./components/layout-header.mjs";
import _sfc_main$4 from "./components/layout-sidebar.mjs";
import _sfc_main$5 from "./components/layout-tabbar.mjs";
import { useLayout } from "./hooks/use-layout.mjs";
import { computed, createBlock, createCommentVNode, createElementBlock, createElementVNode, createSlots, createVNode, defineComponent, mergeModels, normalizeClass, normalizeStyle, openBlock, ref, renderSlot, unref, useModel, watch, withCtx } from "vue";
import { SCROLL_FIXED_CLASS, useLayoutFooterStyle, useLayoutHeaderStyle } from "@vben-core/composables";
import { IconifyIcon } from "@vben-core/icons";
import { VbenIconButton } from "@vben-core/shadcn-ui";
import { ELEMENT_ID_MAIN_CONTENT } from "@vben-core/shared/constants";
import { useMouse, useScroll, useThrottleFn } from "@vueuse/core";
//#region src/vben-layout.vue
const _hoisted_1 = { class: "relative flex min-h-full w-full" };
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "VbenLayout",
	__name: "vben-layout",
	props: /* @__PURE__ */ mergeModels({
		contentCompact: { default: "wide" },
		contentCompactWidth: { default: 1200 },
		contentPadding: { default: 0 },
		contentPaddingBottom: { default: 0 },
		contentPaddingLeft: { default: 0 },
		contentPaddingRight: { default: 0 },
		contentPaddingTop: { default: 0 },
		footerEnable: {
			type: Boolean,
			default: false
		},
		footerFixed: {
			type: Boolean,
			default: true
		},
		footerHeight: { default: 32 },
		headerHeight: { default: 50 },
		headerHidden: {
			type: Boolean,
			default: false
		},
		headerMode: { default: "fixed" },
		headerTheme: {},
		headerToggleSidebarButton: {
			type: Boolean,
			default: true
		},
		headerVisible: {
			type: Boolean,
			default: true
		},
		isMobile: {
			type: Boolean,
			default: false
		},
		layout: { default: "sidebar-nav" },
		sidebarCollapse: { type: Boolean },
		sidebarCollapsedButton: {
			type: Boolean,
			default: true
		},
		sidebarCollapseShowTitle: {
			type: Boolean,
			default: false
		},
		sidebarEnable: { type: Boolean },
		sidebarExtraCollapsedWidth: { default: 60 },
		sidebarFixedButton: {
			type: Boolean,
			default: true
		},
		sidebarHidden: {
			type: Boolean,
			default: false
		},
		sidebarMixedWidth: { default: 80 },
		sidebarTheme: { default: "dark" },
		sidebarThemeSub: { default: "dark" },
		sidebarWidth: { default: 180 },
		sideCollapseWidth: { default: 60 },
		tabbarEnable: {
			type: Boolean,
			default: true
		},
		tabbarHeight: { default: 40 },
		zIndex: { default: 200 }
	}, {
		"sidebarDraggable": {
			type: Boolean,
			default: true
		},
		"sidebarDraggableModifiers": {},
		"sidebarCollapse": {
			type: Boolean,
			default: false
		},
		"sidebarCollapseModifiers": {},
		"sidebarExtraVisible": { type: Boolean },
		"sidebarExtraVisibleModifiers": {},
		"sidebarExtraCollapse": {
			type: Boolean,
			default: false
		},
		"sidebarExtraCollapseModifiers": {},
		"sidebarExpandOnHover": {
			type: Boolean,
			default: false
		},
		"sidebarExpandOnHoverModifiers": {},
		"sidebarEnable": {
			type: Boolean,
			default: true
		},
		"sidebarEnableModifiers": {}
	}),
	emits: /* @__PURE__ */ mergeModels([
		"sideMouseLeave",
		"toggleSidebar",
		"update:sidebar-width"
	], [
		"update:sidebarDraggable",
		"update:sidebarCollapse",
		"update:sidebarExtraVisible",
		"update:sidebarExtraCollapse",
		"update:sidebarExpandOnHover",
		"update:sidebarEnable"
	]),
	setup(__props, { emit: __emit }) {
		const props = __props;
		const emit = __emit;
		const sidebarDraggable = useModel(__props, "sidebarDraggable");
		const sidebarCollapse = useModel(__props, "sidebarCollapse");
		const sidebarExtraVisible = useModel(__props, "sidebarExtraVisible");
		const sidebarExtraCollapse = useModel(__props, "sidebarExtraCollapse");
		const sidebarExpandOnHover = useModel(__props, "sidebarExpandOnHover");
		const sidebarEnable = useModel(__props, "sidebarEnable");
		const sidebarExpandOnHovering = ref(false);
		const headerIsHidden = ref(false);
		const contentRef = ref();
		const { arrivedState, directions, isScrolling, y: scrollY } = useScroll(document);
		const { setLayoutHeaderHeight } = useLayoutHeaderStyle();
		const { setLayoutFooterHeight } = useLayoutFooterStyle();
		const { y: mouseY } = useMouse({
			target: contentRef,
			type: "client"
		});
		const { currentLayout, isFullContent, isHeaderMixedNav, isHeaderNav, isMixedNav, isSidebarMixedNav } = useLayout(props);
		/**
		* 顶栏是否自动隐藏
		*/
		const isHeaderAutoMode = computed(() => props.headerMode === "auto");
		const headerWrapperHeight = computed(() => {
			let height = 0;
			if (props.headerVisible && !props.headerHidden) height += props.headerHeight;
			if (props.tabbarEnable) height += props.tabbarHeight;
			return height;
		});
		const getSideCollapseWidth = computed(() => {
			const { sidebarCollapseShowTitle, sidebarExtraCollapsedWidth, sideCollapseWidth } = props;
			return sidebarCollapseShowTitle || isSidebarMixedNav.value || isHeaderMixedNav.value ? sidebarExtraCollapsedWidth : sideCollapseWidth;
		});
		/**
		* 动态获取侧边区域是否可见
		*/
		const sidebarEnableState = computed(() => {
			return !isHeaderNav.value && sidebarEnable.value;
		});
		/**
		* 侧边区域离顶部高度
		*/
		const sidebarMarginTop = computed(() => {
			const { headerHeight, isMobile } = props;
			return isMixedNav.value && !isMobile ? headerHeight : 0;
		});
		/**
		* 动态获取侧边宽度
		*/
		const getSidebarWidth = computed(() => {
			const { isMobile, sidebarHidden, sidebarMixedWidth, sidebarWidth } = props;
			let width = 0;
			if (sidebarHidden) return width;
			if (!sidebarEnableState.value || sidebarHidden && !isSidebarMixedNav.value && !isMixedNav.value && !isHeaderMixedNav.value) return width;
			if ((isHeaderMixedNav.value || isSidebarMixedNav.value) && !isMobile) width = sidebarMixedWidth;
			else if (sidebarCollapse.value) width = isMobile ? 0 : getSideCollapseWidth.value;
			else width = sidebarWidth;
			return width;
		});
		/**
		* 获取扩展区域宽度
		*/
		const sidebarExtraWidth = computed(() => {
			const { sidebarExtraCollapsedWidth, sidebarWidth } = props;
			return sidebarExtraCollapse.value ? sidebarExtraCollapsedWidth : sidebarWidth;
		});
		/**
		* 是否侧边栏模式，包含混合侧边
		*/
		const isSideMode = computed(() => currentLayout.value === "mixed-nav" || currentLayout.value === "sidebar-mixed-nav" || currentLayout.value === "sidebar-nav" || currentLayout.value === "header-mixed-nav" || currentLayout.value === "header-sidebar-nav");
		/**
		* header fixed值
		*/
		const headerFixed = computed(() => {
			const { headerMode } = props;
			return isMixedNav.value || headerMode === "fixed" || headerMode === "auto-scroll" || headerMode === "auto";
		});
		const showSidebar = computed(() => {
			return isSideMode.value && sidebarEnable.value && !props.sidebarHidden;
		});
		/**
		* 遮罩可见性
		*/
		const maskVisible = computed(() => !sidebarCollapse.value && props.isMobile);
		const mainStyle = computed(() => {
			let width = "100%";
			let sidebarAndExtraWidth = "unset";
			if (headerFixed.value && currentLayout.value !== "header-nav" && currentLayout.value !== "mixed-nav" && currentLayout.value !== "header-sidebar-nav" && showSidebar.value && !props.isMobile) if ((isSidebarMixedNav.value || isHeaderMixedNav.value) && sidebarExpandOnHover.value && sidebarExtraVisible.value) {
				sidebarAndExtraWidth = `${props.sidebarMixedWidth + (sidebarExtraCollapse.value ? props.sidebarExtraCollapsedWidth : props.sidebarWidth)}px`;
				width = `calc(100% - ${sidebarAndExtraWidth})`;
			} else {
				let sidebarWidth = getSidebarWidth.value;
				if (sidebarExpandOnHovering.value && !sidebarExpandOnHover.value) sidebarWidth = isSidebarMixedNav.value || isHeaderMixedNav.value ? props.sidebarMixedWidth : getSideCollapseWidth.value;
				sidebarAndExtraWidth = `${sidebarWidth}px`;
				width = `calc(100% - ${sidebarAndExtraWidth})`;
			}
			return {
				sidebarAndExtraWidth,
				width
			};
		});
		const tabbarStyle = computed(() => {
			let width;
			let marginLeft = 0;
			if (!isMixedNav.value || props.sidebarHidden) width = "100%";
			else if (sidebarEnable.value) {
				const onHoveringWidth = sidebarExpandOnHover.value ? props.sidebarWidth : getSideCollapseWidth.value;
				marginLeft = sidebarCollapse.value ? getSideCollapseWidth.value : onHoveringWidth;
				width = `calc(100% - ${sidebarCollapse.value ? getSidebarWidth.value : onHoveringWidth}px)`;
			} else width = "100%";
			return {
				marginLeft: `${marginLeft}px`,
				width
			};
		});
		const contentStyle = computed(() => {
			const fixed = headerFixed.value;
			const { footerEnable, footerFixed, footerHeight } = props;
			return {
				marginTop: fixed && !isFullContent.value && !headerIsHidden.value && (!isHeaderAutoMode.value || scrollY.value < headerWrapperHeight.value) ? `${headerWrapperHeight.value}px` : 0,
				paddingBottom: `${footerEnable && footerFixed ? footerHeight : 0}px`
			};
		});
		const headerZIndex = computed(() => {
			const { zIndex } = props;
			return zIndex + (isMixedNav.value ? 1 : 0);
		});
		const headerWrapperStyle = computed(() => {
			const fixed = headerFixed.value;
			return {
				height: isFullContent.value ? "0" : `${headerWrapperHeight.value}px`,
				left: isMixedNav.value ? 0 : mainStyle.value.sidebarAndExtraWidth,
				position: fixed ? "fixed" : "static",
				top: headerIsHidden.value || isFullContent.value ? `-${headerWrapperHeight.value}px` : 0,
				width: mainStyle.value.width,
				"z-index": headerZIndex.value
			};
		});
		/**
		* 侧边栏z-index
		*/
		const sidebarZIndex = computed(() => {
			const { isMobile, zIndex } = props;
			let offset = isMobile || isSideMode.value ? 1 : -1;
			if (isMixedNav.value) offset += 1;
			return zIndex + offset;
		});
		const footerWidth = computed(() => {
			if (!props.footerFixed) return "100%";
			return mainStyle.value.width;
		});
		const maskStyle = computed(() => {
			return { zIndex: props.zIndex };
		});
		const showHeaderToggleButton = computed(() => {
			return props.isMobile || props.headerToggleSidebarButton && isSideMode.value && !isSidebarMixedNav.value && !isMixedNav.value && !props.isMobile;
		});
		const showHeaderLogo = computed(() => {
			return !isSideMode.value || isMixedNav.value || props.isMobile;
		});
		watch(() => props.isMobile, (val) => {
			if (val) sidebarCollapse.value = true;
		}, { immediate: true });
		watch([() => headerWrapperHeight.value, () => isFullContent.value], ([height]) => {
			setLayoutHeaderHeight(isFullContent.value ? 0 : height);
		}, { immediate: true });
		watch(() => props.footerHeight, (height) => {
			setLayoutFooterHeight(height);
		}, { immediate: true });
		{
			const HEADER_TRIGGER_DISTANCE = 12;
			watch([
				() => props.headerMode,
				() => mouseY.value,
				() => headerIsHidden.value
			], () => {
				if (!isHeaderAutoMode.value || isMixedNav.value || isFullContent.value) {
					if (props.headerMode !== "auto-scroll") headerIsHidden.value = false;
					return;
				}
				const isInTriggerZone = mouseY.value <= HEADER_TRIGGER_DISTANCE;
				const isInHeaderZone = !headerIsHidden.value && mouseY.value <= headerWrapperHeight.value;
				headerIsHidden.value = !(isInTriggerZone || isInHeaderZone);
			}, { immediate: true });
		}
		{
			const checkHeaderIsHidden = useThrottleFn((top, bottom, topArrived) => {
				if (scrollY.value < headerWrapperHeight.value) {
					headerIsHidden.value = false;
					return;
				}
				if (topArrived) {
					headerIsHidden.value = false;
					return;
				}
				if (top) headerIsHidden.value = false;
				else if (bottom) headerIsHidden.value = true;
			}, 300);
			watch(() => scrollY.value, () => {
				if (props.headerMode !== "auto-scroll" || isMixedNav.value || isFullContent.value) return;
				if (isScrolling.value) checkHeaderIsHidden(directions.top, directions.bottom, arrivedState.top);
			});
		}
		function handleClickMask() {
			sidebarCollapse.value = true;
		}
		function handleHeaderToggle() {
			if (props.isMobile) sidebarCollapse.value = false;
			else emit("toggleSidebar");
		}
		const idMainContent = ELEMENT_ID_MAIN_CONTENT;
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("div", _hoisted_1, [
				sidebarEnableState.value ? (openBlock(), createBlock(unref(_sfc_main$4), {
					key: 0,
					draggable: sidebarDraggable.value,
					"onUpdate:draggable": _cache[0] || (_cache[0] = ($event) => sidebarDraggable.value = $event),
					collapse: sidebarCollapse.value,
					"onUpdate:collapse": _cache[1] || (_cache[1] = ($event) => sidebarCollapse.value = $event),
					"expand-on-hover": sidebarExpandOnHover.value,
					"onUpdate:expandOnHover": _cache[2] || (_cache[2] = ($event) => sidebarExpandOnHover.value = $event),
					"expand-on-hovering": sidebarExpandOnHovering.value,
					"onUpdate:expandOnHovering": _cache[3] || (_cache[3] = ($event) => sidebarExpandOnHovering.value = $event),
					"extra-collapse": sidebarExtraCollapse.value,
					"onUpdate:extraCollapse": _cache[4] || (_cache[4] = ($event) => sidebarExtraCollapse.value = $event),
					"extra-visible": sidebarExtraVisible.value,
					"onUpdate:extraVisible": _cache[5] || (_cache[5] = ($event) => sidebarExtraVisible.value = $event),
					"show-collapse-button": __props.sidebarCollapsedButton,
					"show-fixed-button": __props.sidebarFixedButton,
					"collapse-width": getSideCollapseWidth.value,
					"dom-visible": !__props.isMobile,
					"extra-width": sidebarExtraWidth.value,
					"fixed-extra": sidebarExpandOnHover.value,
					"header-height": unref(isMixedNav) ? 0 : __props.headerHeight,
					"is-sidebar-mixed": unref(isSidebarMixedNav) || unref(isHeaderMixedNav),
					"margin-top": sidebarMarginTop.value,
					"mixed-width": __props.sidebarMixedWidth,
					show: showSidebar.value,
					theme: __props.sidebarTheme,
					"theme-sub": __props.sidebarThemeSub,
					width: getSidebarWidth.value,
					"z-index": sidebarZIndex.value,
					onLeave: _cache[6] || (_cache[6] = () => emit("sideMouseLeave")),
					"onUpdate:width": _cache[7] || (_cache[7] = (val) => emit("update:sidebar-width", val))
				}, createSlots({
					extra: withCtx(() => [renderSlot(_ctx.$slots, "side-extra")]),
					"extra-title": withCtx(() => [renderSlot(_ctx.$slots, "side-extra-title")]),
					default: withCtx(() => [unref(isSidebarMixedNav) || unref(isHeaderMixedNav) ? renderSlot(_ctx.$slots, "mixed-menu", { key: 0 }) : renderSlot(_ctx.$slots, "menu", { key: 1 })]),
					_: 2
				}, [isSideMode.value && !unref(isMixedNav) ? {
					name: "logo",
					fn: withCtx(() => [renderSlot(_ctx.$slots, "logo")]),
					key: "0"
				} : void 0]), 1032, [
					"draggable",
					"collapse",
					"expand-on-hover",
					"expand-on-hovering",
					"extra-collapse",
					"extra-visible",
					"show-collapse-button",
					"show-fixed-button",
					"collapse-width",
					"dom-visible",
					"extra-width",
					"fixed-extra",
					"header-height",
					"is-sidebar-mixed",
					"margin-top",
					"mixed-width",
					"show",
					"theme",
					"theme-sub",
					"width",
					"z-index"
				])) : createCommentVNode("v-if", true),
				createElementVNode("div", {
					ref_key: "contentRef",
					ref: contentRef,
					class: "flex flex-1 flex-col overflow-hidden transition-all duration-300 ease-in"
				}, [
					createElementVNode("div", {
						class: normalizeClass([[{ "shadow-[0_16px_24px_hsl(var(--background))]": unref(scrollY) > 20 }, unref(SCROLL_FIXED_CLASS)], "overflow-hidden transition-all duration-200"]),
						style: normalizeStyle(headerWrapperStyle.value)
					}, [__props.headerVisible ? (openBlock(), createBlock(unref(_sfc_main$3), {
						key: 0,
						"full-width": !isSideMode.value,
						height: __props.headerHeight,
						"is-mobile": __props.isMobile,
						show: !unref(isFullContent) && !__props.headerHidden,
						"sidebar-width": __props.sidebarWidth,
						theme: __props.headerTheme,
						width: mainStyle.value.width,
						"z-index": headerZIndex.value
					}, createSlots({
						"toggle-button": withCtx(() => [showHeaderToggleButton.value ? (openBlock(), createBlock(unref(VbenIconButton), {
							key: 0,
							class: "my-0 mr-1 rounded-md",
							onClick: handleHeaderToggle
						}, {
							default: withCtx(() => [showSidebar.value ? (openBlock(), createBlock(unref(IconifyIcon), {
								key: 0,
								icon: "ep:fold"
							})) : (openBlock(), createBlock(unref(IconifyIcon), {
								key: 1,
								icon: "ep:expand"
							}))]),
							_: 1
						})) : createCommentVNode("v-if", true)]),
						default: withCtx(() => [renderSlot(_ctx.$slots, "header")]),
						_: 2
					}, [showHeaderLogo.value ? {
						name: "logo",
						fn: withCtx(() => [renderSlot(_ctx.$slots, "logo")]),
						key: "0"
					} : void 0]), 1032, [
						"full-width",
						"height",
						"is-mobile",
						"show",
						"sidebar-width",
						"theme",
						"width",
						"z-index"
					])) : createCommentVNode("v-if", true), __props.tabbarEnable ? (openBlock(), createBlock(unref(_sfc_main$5), {
						key: 1,
						height: __props.tabbarHeight,
						style: normalizeStyle(tabbarStyle.value)
					}, {
						default: withCtx(() => [renderSlot(_ctx.$slots, "tabbar")]),
						_: 3
					}, 8, ["height", "style"])) : createCommentVNode("v-if", true)], 6),
					createCommentVNode(" </div> "),
					createVNode(unref(_sfc_main$1), {
						id: unref(idMainContent),
						"content-compact": __props.contentCompact,
						"content-compact-width": __props.contentCompactWidth,
						padding: __props.contentPadding,
						"padding-bottom": __props.contentPaddingBottom,
						"padding-left": __props.contentPaddingLeft,
						"padding-right": __props.contentPaddingRight,
						"padding-top": __props.contentPaddingTop,
						style: normalizeStyle(contentStyle.value),
						class: "transition-[margin-top] duration-200"
					}, {
						overlay: withCtx(() => [renderSlot(_ctx.$slots, "content-overlay")]),
						default: withCtx(() => [renderSlot(_ctx.$slots, "content")]),
						_: 3
					}, 8, [
						"id",
						"content-compact",
						"content-compact-width",
						"padding",
						"padding-bottom",
						"padding-left",
						"padding-right",
						"padding-top",
						"style"
					]),
					__props.footerEnable ? (openBlock(), createBlock(unref(_sfc_main$2), {
						key: 0,
						fixed: __props.footerFixed,
						height: __props.footerHeight,
						show: !unref(isFullContent),
						width: footerWidth.value,
						"z-index": __props.zIndex
					}, {
						default: withCtx(() => [renderSlot(_ctx.$slots, "footer")]),
						_: 3
					}, 8, [
						"fixed",
						"height",
						"show",
						"width",
						"z-index"
					])) : createCommentVNode("v-if", true)
				], 512),
				renderSlot(_ctx.$slots, "extra"),
				maskVisible.value ? (openBlock(), createElementBlock("div", {
					key: 1,
					style: normalizeStyle(maskStyle.value),
					class: "fixed top-0 left-0 size-full bg-overlay transition-[background-color] duration-200",
					onClick: handleClickMask
				}, null, 4)) : createCommentVNode("v-if", true)
			]);
		};
	}
});
//#endregion
export { _sfc_main as default };
