import { useSidebarDrag } from "../hooks/use-sidebar-drag.mjs";
import _sfc_main$1 from "./widgets/sidebar-collapse-button.mjs";
import _sfc_main$2 from "./widgets/sidebar-fixed-button.mjs";
import { Fragment, computed, createBlock, createCommentVNode, createElementBlock, createElementVNode, createVNode, defineComponent, mergeModels, normalizeClass, normalizeStyle, onUnmounted, openBlock, renderSlot, shallowRef, unref, useModel, useSlots, watchEffect, withCtx } from "vue";
import { VbenScrollbar } from "@vben-core/shadcn-ui";
import { useScrollLock } from "@vueuse/core";
//#region src/components/layout-sidebar.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "layout-sidebar",
	props: /* @__PURE__ */ mergeModels({
		collapseHeight: { default: 42 },
		collapseWidth: { default: 48 },
		domVisible: {
			type: Boolean,
			default: true
		},
		extraWidth: {},
		fixedExtra: {
			type: Boolean,
			default: false
		},
		headerHeight: {},
		isSidebarMixed: {
			type: Boolean,
			default: false
		},
		marginTop: { default: 0 },
		mixedWidth: { default: 70 },
		paddingTop: { default: 0 },
		show: {
			type: Boolean,
			default: true
		},
		showCollapseButton: {
			type: Boolean,
			default: true
		},
		showFixedButton: {
			type: Boolean,
			default: true
		},
		theme: {},
		themeSub: {},
		width: {},
		zIndex: { default: 0 }
	}, {
		"draggable": { type: Boolean },
		"draggableModifiers": {},
		"collapse": { type: Boolean },
		"collapseModifiers": {},
		"extraCollapse": { type: Boolean },
		"extraCollapseModifiers": {},
		"expandOnHovering": { type: Boolean },
		"expandOnHoveringModifiers": {},
		"expandOnHover": { type: Boolean },
		"expandOnHoverModifiers": {},
		"extraVisible": { type: Boolean },
		"extraVisibleModifiers": {}
	}),
	emits: /* @__PURE__ */ mergeModels(["leave", "update:width"], [
		"update:draggable",
		"update:collapse",
		"update:extraCollapse",
		"update:expandOnHovering",
		"update:expandOnHover",
		"update:extraVisible"
	]),
	setup(__props, { emit: __emit }) {
		const props = __props;
		const emit = __emit;
		const draggable = useModel(__props, "draggable");
		const collapse = useModel(__props, "collapse");
		const extraCollapse = useModel(__props, "extraCollapse");
		const expandOnHovering = useModel(__props, "expandOnHovering");
		const expandOnHover = useModel(__props, "expandOnHover");
		const extraVisible = useModel(__props, "extraVisible");
		const isLocked = useScrollLock(document.body);
		const slots = useSlots();
		const asideRef = shallowRef(null);
		const dragBarRef = shallowRef(null);
		const hiddenSideStyle = computed(() => calcMenuWidthStyle(true));
		const style = computed(() => {
			const { isSidebarMixed, marginTop, paddingTop, zIndex } = props;
			return {
				"--scroll-shadow": "var(--sidebar)",
				...calcMenuWidthStyle(false),
				height: `calc(100% - ${marginTop}px)`,
				marginTop: `${marginTop}px`,
				paddingTop: `${paddingTop}px`,
				zIndex,
				...isSidebarMixed && extraVisible.value ? { transition: "none" } : {}
			};
		});
		const extraStyle = computed(() => {
			const { extraWidth, show, width, zIndex } = props;
			return {
				left: `${width}px`,
				width: extraVisible.value && show ? `${extraWidth}px` : 0,
				zIndex
			};
		});
		const extraTitleStyle = computed(() => {
			const { headerHeight } = props;
			return { height: `${headerHeight - 1}px` };
		});
		const contentWidthStyle = computed(() => {
			const { fixedExtra, isSidebarMixed, mixedWidth } = props;
			if (isSidebarMixed && fixedExtra) return { width: `${mixedWidth}px` };
			return {};
		});
		const contentStyle = computed(() => {
			const { collapseHeight, headerHeight } = props;
			return {
				height: `calc(100% - ${headerHeight + collapseHeight}px)`,
				paddingTop: "8px",
				...contentWidthStyle.value
			};
		});
		const headerStyle = computed(() => {
			const { headerHeight, isSidebarMixed } = props;
			return {
				...isSidebarMixed ? {
					display: "flex",
					justifyContent: "center"
				} : {},
				height: `${headerHeight - 1}px`,
				...contentWidthStyle.value
			};
		});
		const extraContentStyle = computed(() => {
			const { collapseHeight, headerHeight } = props;
			return { height: `calc(100% - ${headerHeight + collapseHeight}px)` };
		});
		const collapseStyle = computed(() => {
			return { height: `${props.collapseHeight}px` };
		});
		watchEffect(() => {
			extraVisible.value = props.fixedExtra ? true : extraVisible.value;
		});
		function calcMenuWidthStyle(isHiddenDom) {
			const { collapseWidth, extraWidth, mixedWidth, fixedExtra, isSidebarMixed, show, width } = props;
			let widthValue = width === 0 ? "0px" : `${width + (isSidebarMixed && fixedExtra && extraVisible.value ? extraWidth : 0)}px`;
			if (isHiddenDom && expandOnHovering.value && !expandOnHover.value) widthValue = isSidebarMixed ? `${mixedWidth}px` : `${collapseWidth}px`;
			return {
				...widthValue === "0px" ? { overflow: "hidden" } : {},
				flex: `0 0 ${widthValue}`,
				marginLeft: show ? 0 : `-${widthValue}`,
				maxWidth: widthValue,
				minWidth: widthValue,
				width: widthValue
			};
		}
		function handleMouseenter(e) {
			if (e?.offsetX < 10) return;
			if (expandOnHover.value) return;
			if (!expandOnHovering.value) collapse.value = false;
			if (props.isSidebarMixed) isLocked.value = true;
			expandOnHovering.value = true;
		}
		function handleMouseleave() {
			emit("leave");
			if (props.isSidebarMixed) isLocked.value = false;
			if (expandOnHover.value) return;
			expandOnHovering.value = false;
			collapse.value = true;
			extraVisible.value = false;
		}
		const { startDrag, endDrag } = useSidebarDrag();
		const handleDragSidebar = (e) => {
			const { isSidebarMixed, collapseWidth, width } = props;
			startDrag(e, {
				min: isSidebarMixed ? width + collapseWidth : collapseWidth,
				max: isSidebarMixed ? width + 320 : 320
			}, {
				target: asideRef.value,
				dragBar: dragBarRef.value
			}, (newWidth) => {
				if (isSidebarMixed) {
					emit("update:width", newWidth - width);
					extraCollapse.value = collapse.value = newWidth - width <= collapseWidth;
				} else {
					emit("update:width", newWidth);
					collapse.value = extraCollapse.value = newWidth <= collapseWidth;
				}
			});
		};
		onUnmounted(() => {
			endDrag();
		});
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock(Fragment, null, [__props.domVisible ? (openBlock(), createElementBlock("div", {
				key: 0,
				class: normalizeClass([__props.theme, "h-full transition-all duration-150"]),
				style: normalizeStyle(hiddenSideStyle.value)
			}, null, 6)) : createCommentVNode("v-if", true), createElementVNode("aside", {
				ref_key: "asideRef",
				ref: asideRef,
				style: normalizeStyle(style.value),
				class: normalizeClass(["fixed left-0 top-0 h-full transition-all duration-150", __props.theme]),
				onMouseenter: handleMouseenter,
				onMouseleave: handleMouseleave
			}, [
				createElementVNode("div", {
					class: normalizeClass(["h-full", [{
						"bg-sidebar-deep": __props.isSidebarMixed,
						"border-r border-border bg-sidebar": !__props.isSidebarMixed
					}]]),
					style: normalizeStyle({ width: `${__props.width}px` })
				}, [
					!collapse.value && !__props.isSidebarMixed && __props.showFixedButton ? (openBlock(), createBlock(unref(_sfc_main$2), {
						key: 0,
						"expand-on-hover": expandOnHover.value,
						"onUpdate:expandOnHover": _cache[0] || (_cache[0] = ($event) => expandOnHover.value = $event)
					}, null, 8, ["expand-on-hover"])) : createCommentVNode("v-if", true),
					unref(slots).logo ? (openBlock(), createElementBlock("div", {
						key: 1,
						style: normalizeStyle(headerStyle.value)
					}, [renderSlot(_ctx.$slots, "logo")], 4)) : createCommentVNode("v-if", true),
					createVNode(unref(VbenScrollbar), {
						style: normalizeStyle(contentStyle.value),
						shadow: "",
						"shadow-border": ""
					}, {
						default: withCtx(() => [renderSlot(_ctx.$slots, "default")]),
						_: 3
					}, 8, ["style"]),
					createElementVNode("div", { style: normalizeStyle(collapseStyle.value) }, null, 4),
					__props.showCollapseButton && !__props.isSidebarMixed ? (openBlock(), createBlock(unref(_sfc_main$1), {
						key: 2,
						collapsed: collapse.value,
						"onUpdate:collapsed": _cache[1] || (_cache[1] = ($event) => collapse.value = $event)
					}, null, 8, ["collapsed"])) : createCommentVNode("v-if", true)
				], 6),
				__props.isSidebarMixed ? (openBlock(), createElementBlock("div", {
					key: 0,
					class: normalizeClass([[__props.themeSub, { "border-l": extraVisible.value }], "fixed top-0 h-full overflow-hidden border-r border-border bg-sidebar transition-all duration-200"]),
					style: normalizeStyle(extraStyle.value)
				}, [
					__props.isSidebarMixed && expandOnHover.value ? (openBlock(), createBlock(unref(_sfc_main$1), {
						key: 0,
						collapsed: extraCollapse.value,
						"onUpdate:collapsed": _cache[2] || (_cache[2] = ($event) => extraCollapse.value = $event)
					}, null, 8, ["collapsed"])) : createCommentVNode("v-if", true),
					!extraCollapse.value ? (openBlock(), createBlock(unref(_sfc_main$2), {
						key: 1,
						"expand-on-hover": expandOnHover.value,
						"onUpdate:expandOnHover": _cache[3] || (_cache[3] = ($event) => expandOnHover.value = $event)
					}, null, 8, ["expand-on-hover"])) : createCommentVNode("v-if", true),
					!extraCollapse.value ? (openBlock(), createElementBlock("div", {
						key: 2,
						style: normalizeStyle(extraTitleStyle.value),
						class: "pl-2"
					}, [renderSlot(_ctx.$slots, "extra-title")], 4)) : createCommentVNode("v-if", true),
					createVNode(unref(VbenScrollbar), {
						style: normalizeStyle(extraContentStyle.value),
						class: "border-border py-2",
						shadow: "",
						"shadow-border": ""
					}, {
						default: withCtx(() => [renderSlot(_ctx.$slots, "extra")]),
						_: 3
					}, 8, ["style"])
				], 6)) : createCommentVNode("v-if", true),
				draggable.value ? (openBlock(), createElementBlock("div", {
					key: 1,
					ref_key: "dragBarRef",
					ref: dragBarRef,
					class: "absolute inset-y-0 -right-px z-1000 w-0.5 cursor-col-resize hover:bg-primary",
					onMousedown: handleDragSidebar
				}, null, 544)) : createCommentVNode("v-if", true)
			], 38)], 64);
		};
	}
});
//#endregion
export { _sfc_main as default };
