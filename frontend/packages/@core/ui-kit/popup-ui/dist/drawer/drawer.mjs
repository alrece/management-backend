import { Separator, Sheet, SheetClose, SheetContent, SheetDescription, SheetFooter, SheetHeader, SheetTitle, VbenButton, VbenHelpTooltip, VbenIconButton, VbenLoading, VisuallyHidden } from "@vben-core/shadcn-ui";
import { computed, createBlock, createCommentVNode, createElementVNode, createTextVNode, createVNode, defineComponent, normalizeClass, onDeactivated, openBlock, provide, ref, renderSlot, resolveDynamicComponent, toDisplayString, unref, useId, watch, withCtx } from "vue";
import { useIsMobile, usePriorityValues, useSimpleLocale } from "@vben-core/composables";
import { X } from "@vben-core/icons";
import { globalShareState } from "@vben-core/shared/global-state";
import { cn } from "@vben-core/shared/utils";
import { ELEMENT_ID_MAIN_CONTENT } from "@vben-core/shared/constants";
//#region src/drawer/drawer.vue
const _hoisted_1 = { class: "flex items-center" };
const _hoisted_2 = { class: "flex-center" };
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "drawer",
	props: {
		drawerApi: { default: void 0 },
		appendToMain: {
			type: Boolean,
			default: false
		},
		cancelText: {},
		class: { type: [
			Array,
			Boolean,
			null,
			Object,
			String
		] },
		closable: { type: Boolean },
		closeIconPlacement: { default: "right" },
		closeOnClickModal: { type: Boolean },
		closeOnPressEscape: { type: Boolean },
		confirmLoading: { type: Boolean },
		confirmText: {},
		contentClass: {},
		description: {},
		destroyOnClose: {
			type: Boolean,
			default: false
		},
		footer: { type: Boolean },
		footerClass: { type: [
			Array,
			Boolean,
			null,
			Object,
			String
		] },
		header: { type: Boolean },
		headerClass: { type: [
			Array,
			Boolean,
			null,
			Object,
			String
		] },
		loading: { type: Boolean },
		modal: { type: Boolean },
		openAutoFocus: { type: Boolean },
		overlayBlur: {},
		placement: {},
		showCancelButton: { type: Boolean },
		showConfirmButton: { type: Boolean },
		submitting: {
			type: Boolean,
			default: false
		},
		title: {},
		titleTooltip: {},
		zIndex: { default: 1e3 }
	},
	setup(__props) {
		const props = __props;
		const components = globalShareState.getComponents();
		const id = useId();
		provide("DISMISSABLE_DRAWER_ID", id);
		const wrapperRef = ref();
		const { $t } = useSimpleLocale();
		const { isMobile } = useIsMobile();
		const state = props.drawerApi?.useStore?.();
		const { appendToMain, cancelText, class: drawerClass, closable, closeIconPlacement, closeOnClickModal, closeOnPressEscape, confirmLoading, confirmText, contentClass, description, destroyOnClose, footer: showFooter, footerClass, header: showHeader, headerClass, loading: showLoading, modal, openAutoFocus, overlayBlur, placement, showCancelButton, showConfirmButton, submitting, title, titleTooltip, zIndex } = usePriorityValues(props, state);
		/**
		* 在开启keepAlive情况下 直接通过浏览器按钮/手势等返回 不会关闭弹窗
		*/
		onDeactivated(() => {
			if (!appendToMain.value) props.drawerApi?.close();
		});
		function interactOutside(e) {
			if (!closeOnClickModal.value || submitting.value) e.preventDefault();
		}
		function escapeKeyDown(e) {
			if (!closeOnPressEscape.value || submitting.value) e.preventDefault();
		}
		function pointerDownOutside(e) {
			const dismissableDrawer = e.target?.dataset.dismissableDrawer;
			if (submitting.value || !closeOnClickModal.value || dismissableDrawer !== id) e.preventDefault();
		}
		function handerOpenAutoFocus(e) {
			if (!openAutoFocus.value) e?.preventDefault();
		}
		function handleFocusOutside(e) {
			e.preventDefault();
			e.stopPropagation();
		}
		const getAppendTo = computed(() => {
			return appendToMain.value ? `#${ELEMENT_ID_MAIN_CONTENT}>div:not(.absolute)>div` : void 0;
		});
		/**
		* destroyOnClose功能完善
		*/
		const hasOpened = ref(false);
		const isClosed = ref(true);
		watch(() => state?.value?.isOpen, (value) => {
			isClosed.value = false;
			if (value && !unref(hasOpened)) hasOpened.value = true;
		});
		function handleClosed() {
			isClosed.value = true;
			props.drawerApi?.onClosed();
		}
		const getForceMount = computed(() => {
			return !unref(destroyOnClose) && unref(hasOpened);
		});
		return (_ctx, _cache) => {
			return openBlock(), createBlock(unref(Sheet), {
				modal: false,
				open: unref(state)?.isOpen,
				"onUpdate:open": _cache[3] || (_cache[3] = () => __props.drawerApi?.close())
			}, {
				default: withCtx(() => [createVNode(unref(SheetContent), {
					"append-to": getAppendTo.value,
					class: normalizeClass(unref(cn)("flex w-130 flex-col", unref(drawerClass), {
						"w-full!": unref(isMobile) || unref(placement) === "bottom" || unref(placement) === "top",
						"max-h-screen": unref(placement) === "bottom" || unref(placement) === "top",
						hidden: isClosed.value
					})),
					modal: unref(modal),
					open: unref(state)?.isOpen,
					side: unref(placement),
					"z-index": unref(zIndex),
					"force-mount": getForceMount.value,
					"overlay-blur": unref(overlayBlur),
					onCloseAutoFocus: handleFocusOutside,
					onClosed: handleClosed,
					onEscapeKeyDown: escapeKeyDown,
					onFocusOutside: handleFocusOutside,
					onInteractOutside: interactOutside,
					onOpenAutoFocus: handerOpenAutoFocus,
					onOpened: _cache[2] || (_cache[2] = () => __props.drawerApi?.onOpened()),
					onPointerDownOutside: pointerDownOutside
				}, {
					default: withCtx(() => [
						unref(showHeader) ? (openBlock(), createBlock(unref(SheetHeader), {
							key: 0,
							class: normalizeClass(unref(cn)("flex! flex-row items-center justify-between border-b px-6 py-5", unref(headerClass), {
								"px-4 py-3": unref(closable),
								"pl-2": unref(closable) && unref(closeIconPlacement) === "left"
							}))
						}, {
							default: withCtx(() => [
								createElementVNode("div", _hoisted_1, [
									unref(closable) && unref(closeIconPlacement) === "left" ? (openBlock(), createBlock(unref(SheetClose), {
										key: 0,
										"as-child": "",
										disabled: unref(submitting),
										class: "ml-0.5 cursor-pointer rounded-full opacity-80 transition-opacity hover:opacity-100 focus:outline-hidden disabled:pointer-events-none data-[state=open]:bg-secondary"
									}, {
										default: withCtx(() => [renderSlot(_ctx.$slots, "close-icon", {}, () => [createVNode(unref(VbenIconButton), null, {
											default: withCtx(() => [createVNode(unref(X), { class: "size-4" })]),
											_: 1
										})])]),
										_: 3
									}, 8, ["disabled"])) : createCommentVNode("v-if", true),
									unref(closable) && unref(closeIconPlacement) === "left" ? (openBlock(), createBlock(unref(Separator), {
										key: 1,
										class: "mr-2 ml-1 h-8",
										decorative: "",
										orientation: "vertical"
									})) : createCommentVNode("v-if", true),
									unref(title) ? (openBlock(), createBlock(unref(SheetTitle), {
										key: 2,
										class: "text-left"
									}, {
										default: withCtx(() => [renderSlot(_ctx.$slots, "title", {}, () => [createTextVNode(toDisplayString(unref(title)) + " ", 1), unref(titleTooltip) ? (openBlock(), createBlock(unref(VbenHelpTooltip), {
											key: 0,
											"trigger-class": "pb-1"
										}, {
											default: withCtx(() => [createTextVNode(toDisplayString(unref(titleTooltip)), 1)]),
											_: 1
										})) : createCommentVNode("v-if", true)])]),
										_: 3
									})) : createCommentVNode("v-if", true),
									unref(description) ? (openBlock(), createBlock(unref(SheetDescription), {
										key: 3,
										class: "mt-1 text-xs"
									}, {
										default: withCtx(() => [renderSlot(_ctx.$slots, "description", {}, () => [createTextVNode(toDisplayString(unref(description)), 1)])]),
										_: 3
									})) : createCommentVNode("v-if", true)
								]),
								!unref(title) || !unref(description) ? (openBlock(), createBlock(unref(VisuallyHidden), { key: 0 }, {
									default: withCtx(() => [!unref(title) ? (openBlock(), createBlock(unref(SheetTitle), { key: 0 })) : createCommentVNode("v-if", true), !unref(description) ? (openBlock(), createBlock(unref(SheetDescription), { key: 1 })) : createCommentVNode("v-if", true)]),
									_: 1
								})) : createCommentVNode("v-if", true),
								createElementVNode("div", _hoisted_2, [renderSlot(_ctx.$slots, "extra"), unref(closable) && unref(closeIconPlacement) === "right" ? (openBlock(), createBlock(unref(SheetClose), {
									key: 0,
									"as-child": "",
									disabled: unref(submitting),
									class: "ml-0.5 cursor-pointer rounded-full opacity-80 transition-opacity hover:opacity-100 focus:outline-hidden disabled:pointer-events-none data-[state=open]:bg-secondary"
								}, {
									default: withCtx(() => [renderSlot(_ctx.$slots, "close-icon", {}, () => [createVNode(unref(VbenIconButton), null, {
										default: withCtx(() => [createVNode(unref(X), { class: "size-4" })]),
										_: 1
									})])]),
									_: 3
								}, 8, ["disabled"])) : createCommentVNode("v-if", true)])
							]),
							_: 3
						}, 8, ["class"])) : (openBlock(), createBlock(unref(VisuallyHidden), { key: 1 }, {
							default: withCtx(() => [createVNode(unref(SheetTitle)), createVNode(unref(SheetDescription))]),
							_: 1
						})),
						createElementVNode("div", {
							ref_key: "wrapperRef",
							ref: wrapperRef,
							class: normalizeClass(unref(cn)("relative flex-1 overflow-y-auto p-3", unref(contentClass), { "pointer-events-none": unref(showLoading) || unref(submitting) }))
						}, [renderSlot(_ctx.$slots, "default")], 2),
						unref(showLoading) || unref(submitting) ? (openBlock(), createBlock(unref(VbenLoading), {
							key: 2,
							spinning: ""
						})) : createCommentVNode("v-if", true),
						unref(showFooter) ? (openBlock(), createBlock(unref(SheetFooter), {
							key: 3,
							class: normalizeClass(unref(cn)("w-full flex-row items-center justify-end border-t p-2 px-3", unref(footerClass)))
						}, {
							default: withCtx(() => [
								renderSlot(_ctx.$slots, "prepend-footer"),
								renderSlot(_ctx.$slots, "footer", {}, () => [
									unref(showCancelButton) ? (openBlock(), createBlock(resolveDynamicComponent(unref(components).DefaultButton || unref(VbenButton)), {
										key: 0,
										variant: "ghost",
										disabled: unref(submitting),
										onClick: _cache[0] || (_cache[0] = () => __props.drawerApi?.onCancel())
									}, {
										default: withCtx(() => [renderSlot(_ctx.$slots, "cancelText", {}, () => [createTextVNode(toDisplayString(unref(cancelText) || unref($t)("cancel")), 1)])]),
										_: 3
									}, 8, ["disabled"])) : createCommentVNode("v-if", true),
									renderSlot(_ctx.$slots, "center-footer"),
									unref(showConfirmButton) ? (openBlock(), createBlock(resolveDynamicComponent(unref(components).PrimaryButton || unref(VbenButton)), {
										key: 1,
										loading: unref(confirmLoading) || unref(submitting),
										onClick: _cache[1] || (_cache[1] = () => __props.drawerApi?.onConfirm())
									}, {
										default: withCtx(() => [renderSlot(_ctx.$slots, "confirmText", {}, () => [createTextVNode(toDisplayString(unref(confirmText) || unref($t)("confirm")), 1)])]),
										_: 3
									}, 8, ["loading"])) : createCommentVNode("v-if", true)
								]),
								renderSlot(_ctx.$slots, "append-footer")
							]),
							_: 3
						}, 8, ["class"])) : createCommentVNode("v-if", true)
					]),
					_: 3
				}, 8, [
					"append-to",
					"class",
					"modal",
					"open",
					"side",
					"z-index",
					"force-mount",
					"overlay-blur"
				])]),
				_: 3
			}, 8, ["open"]);
		};
	}
});
//#endregion
export { _sfc_main as default };
