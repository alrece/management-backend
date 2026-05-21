import { useModalDraggable } from "./use-modal-draggable.mjs";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, VbenButton, VbenHelpTooltip, VbenIconButton, VbenLoading, VisuallyHidden } from "@vben-core/shadcn-ui";
import { computed, createBlock, createCommentVNode, createElementVNode, createTextVNode, createVNode, defineComponent, nextTick, normalizeClass, onDeactivated, openBlock, provide, ref, renderSlot, resolveDynamicComponent, toDisplayString, unref, useId, watch, withCtx } from "vue";
import { useIsMobile, usePriorityValues, useSimpleLocale } from "@vben-core/composables";
import { Expand, Shrink } from "@vben-core/icons";
import { globalShareState } from "@vben-core/shared/global-state";
import { cn } from "@vben-core/shared/utils";
import { ELEMENT_ID_MAIN_CONTENT } from "@vben-core/shared/constants";
//#region src/modal/modal.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "modal",
	props: {
		modalApi: { default: void 0 },
		animationType: {},
		appendToMain: {
			type: Boolean,
			default: false
		},
		bordered: { type: Boolean },
		cancelText: {},
		centered: { type: Boolean },
		class: { type: [
			Array,
			Boolean,
			null,
			Object,
			String
		] },
		closable: { type: Boolean },
		closeOnClickModal: { type: Boolean },
		closeOnPressEscape: { type: Boolean },
		confirmDisabled: { type: Boolean },
		confirmLoading: { type: Boolean },
		confirmText: {},
		contentClass: { type: [
			Array,
			Boolean,
			null,
			Object,
			String
		] },
		description: {},
		destroyOnClose: {
			type: Boolean,
			default: false
		},
		draggable: { type: Boolean },
		footer: { type: Boolean },
		footerClass: { type: [
			Array,
			Boolean,
			null,
			Object,
			String
		] },
		fullscreen: { type: Boolean },
		fullscreenButton: { type: Boolean },
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
		overflow: { type: Boolean },
		overlayBlur: {},
		showCancelButton: { type: Boolean },
		showConfirmButton: { type: Boolean },
		submitting: { type: Boolean },
		title: {},
		titleTooltip: {},
		zIndex: {}
	},
	setup(__props) {
		const props = __props;
		const components = globalShareState.getComponents();
		const contentRef = ref();
		const wrapperRef = ref();
		const dialogRef = ref();
		const headerRef = ref();
		const footerRef = ref();
		const id = useId();
		provide("DISMISSABLE_MODAL_ID", id);
		const { $t } = useSimpleLocale();
		const { isMobile } = useIsMobile();
		const state = props.modalApi?.useStore?.();
		const { appendToMain, bordered, cancelText, centered, class: modalClass, closable, closeOnClickModal, closeOnPressEscape, confirmDisabled, confirmLoading, confirmText, contentClass, description, destroyOnClose, draggable, overflow, footer: showFooter, footerClass, fullscreen, fullscreenButton, header, headerClass, loading: showLoading, modal, openAutoFocus, overlayBlur, showCancelButton, showConfirmButton, submitting, title, titleTooltip, animationType, zIndex } = usePriorityValues(props, state);
		const shouldFullscreen = computed(() => fullscreen.value || isMobile.value);
		const shouldDraggable = computed(() => draggable.value && !shouldFullscreen.value && header.value);
		const shouldCentered = computed(() => centered.value && !shouldFullscreen.value);
		const getAppendTo = computed(() => {
			return appendToMain.value ? `#${ELEMENT_ID_MAIN_CONTENT}>div:not(.absolute)>div` : void 0;
		});
		const { dragging, transform } = useModalDraggable(dialogRef, headerRef, shouldDraggable, getAppendTo, shouldCentered, overflow);
		const firstOpened = ref(false);
		const isClosed = ref(true);
		watch(() => state?.value?.isOpen, async (v) => {
			if (v) {
				isClosed.value = false;
				if (!firstOpened.value) firstOpened.value = true;
				await nextTick();
				if (!contentRef.value) return;
				dialogRef.value = contentRef.value.getContentRef().$el;
				const { offsetX, offsetY } = transform;
				dialogRef.value.style.transform = shouldCentered.value ? `translate(${offsetX}px, calc(-50% + ${offsetY}px))` : `translate(${offsetX}px, ${offsetY}px)`;
			}
		}, { immediate: true });
		/**
		* 在开启keepAlive情况下 直接通过浏览器按钮/手势等返回 不会关闭弹窗
		*/
		onDeactivated(() => {
			if (!appendToMain.value) props.modalApi?.close();
		});
		function handleFullscreen() {
			props.modalApi?.setState((prev) => {
				return {
					...prev,
					fullscreen: !fullscreen.value
				};
			});
		}
		function interactOutside(e) {
			if (!closeOnClickModal.value || submitting.value) {
				e.preventDefault();
				e.stopPropagation();
			}
		}
		function escapeKeyDown(e) {
			if (!closeOnPressEscape.value || submitting.value) e.preventDefault();
		}
		function handleOpenAutoFocus(e) {
			if (!openAutoFocus.value) e?.preventDefault();
		}
		function pointerDownOutside(e) {
			const isDismissableModal = e.target?.dataset.dismissableModal;
			if (!closeOnClickModal.value || isDismissableModal !== id || submitting.value) {
				e.preventDefault();
				e.stopPropagation();
			}
		}
		function handleFocusOutside(e) {
			e.preventDefault();
			e.stopPropagation();
		}
		const getForceMount = computed(() => {
			return !unref(destroyOnClose) && unref(firstOpened);
		});
		const handleOpened = () => {
			requestAnimationFrame(() => {
				props.modalApi?.onOpened();
			});
		};
		function handleClosed() {
			isClosed.value = true;
			props.modalApi?.onClosed();
		}
		return (_ctx, _cache) => {
			return openBlock(), createBlock(unref(Dialog), {
				modal: false,
				open: unref(state)?.isOpen,
				"onUpdate:open": _cache[2] || (_cache[2] = () => !unref(submitting) ? __props.modalApi?.close() : void 0)
			}, {
				default: withCtx(() => [createVNode(unref(DialogContent), {
					ref_key: "contentRef",
					ref: contentRef,
					"append-to": getAppendTo.value,
					class: normalizeClass(unref(cn)("inset-x-0 top-[10vh] mx-auto flex max-h-[80%] w-130 flex-col p-0", shouldFullscreen.value ? "sm:rounded-none" : "sm:rounded-(--radius)", unref(modalClass), {
						"border border-border": unref(bordered),
						"shadow-3xl": !unref(bordered),
						"top-0 left-0 size-full max-h-full transform-[translate(0,0)]!": shouldFullscreen.value,
						"top-1/2": unref(centered) && !shouldFullscreen.value,
						"duration-300": !unref(dragging),
						hidden: isClosed.value
					})),
					"force-mount": getForceMount.value,
					modal: unref(modal),
					open: unref(state)?.isOpen,
					"show-close": unref(closable),
					"animation-type": unref(animationType),
					"z-index": unref(zIndex),
					"overlay-blur": unref(overlayBlur),
					"close-class": "top-3",
					onCloseAutoFocus: handleFocusOutside,
					onClosed: handleClosed,
					"close-disabled": unref(submitting),
					onEscapeKeyDown: escapeKeyDown,
					onFocusOutside: handleFocusOutside,
					onInteractOutside: interactOutside,
					onOpenAutoFocus: handleOpenAutoFocus,
					onOpened: handleOpened,
					onPointerDownOutside: pointerDownOutside
				}, {
					default: withCtx(() => [
						createVNode(unref(DialogHeader), {
							ref_key: "headerRef",
							ref: headerRef,
							class: normalizeClass(unref(cn)("px-5 py-4", {
								"border-b": unref(bordered),
								hidden: !unref(header),
								"cursor-move select-none": shouldDraggable.value
							}, unref(headerClass)))
						}, {
							default: withCtx(() => [
								unref(title) ? (openBlock(), createBlock(unref(DialogTitle), {
									key: 0,
									class: "text-left"
								}, {
									default: withCtx(() => [renderSlot(_ctx.$slots, "title", {}, () => [createTextVNode(toDisplayString(unref(title)) + " ", 1), unref(titleTooltip) ? renderSlot(_ctx.$slots, "titleTooltip", { key: 0 }, () => [createVNode(unref(VbenHelpTooltip), { "trigger-class": "pb-1" }, {
										default: withCtx(() => [createTextVNode(toDisplayString(unref(titleTooltip)), 1)]),
										_: 1
									})]) : createCommentVNode("v-if", true)])]),
									_: 3
								})) : createCommentVNode("v-if", true),
								unref(description) ? (openBlock(), createBlock(unref(DialogDescription), { key: 1 }, {
									default: withCtx(() => [renderSlot(_ctx.$slots, "description", {}, () => [createTextVNode(toDisplayString(unref(description)), 1)])]),
									_: 3
								})) : createCommentVNode("v-if", true),
								!unref(title) || !unref(description) ? (openBlock(), createBlock(unref(VisuallyHidden), { key: 2 }, {
									default: withCtx(() => [!unref(title) ? (openBlock(), createBlock(unref(DialogTitle), { key: 0 })) : createCommentVNode("v-if", true), !unref(description) ? (openBlock(), createBlock(unref(DialogDescription), { key: 1 })) : createCommentVNode("v-if", true)]),
									_: 1
								})) : createCommentVNode("v-if", true)
							]),
							_: 3
						}, 8, ["class"]),
						createElementVNode("div", {
							ref_key: "wrapperRef",
							ref: wrapperRef,
							class: normalizeClass(unref(cn)("relative min-h-40 flex-1 overflow-y-auto p-3", unref(contentClass), { "pointer-events-none": unref(showLoading) || unref(submitting) }))
						}, [renderSlot(_ctx.$slots, "default")], 2),
						unref(showLoading) || unref(submitting) ? (openBlock(), createBlock(unref(VbenLoading), {
							key: 0,
							spinning: ""
						})) : createCommentVNode("v-if", true),
						unref(fullscreenButton) ? (openBlock(), createBlock(unref(VbenIconButton), {
							key: 1,
							class: "absolute top-3 right-10 flex-center hidden size-6 rounded-full px-1 text-lg text-foreground/80 opacity-70 transition-opacity hover:bg-accent hover:text-accent-foreground hover:opacity-100 focus:outline-hidden disabled:pointer-events-none sm:block",
							onClick: handleFullscreen
						}, {
							default: withCtx(() => [unref(fullscreen) ? (openBlock(), createBlock(unref(Shrink), {
								key: 0,
								class: "size-3.5"
							})) : (openBlock(), createBlock(unref(Expand), {
								key: 1,
								class: "size-3.5"
							}))]),
							_: 1
						})) : createCommentVNode("v-if", true),
						unref(showFooter) ? (openBlock(), createBlock(unref(DialogFooter), {
							key: 2,
							ref_key: "footerRef",
							ref: footerRef,
							class: normalizeClass(unref(cn)("flex-row items-center justify-end p-2", { "border-t": unref(bordered) }, unref(footerClass)))
						}, {
							default: withCtx(() => [
								renderSlot(_ctx.$slots, "prepend-footer"),
								renderSlot(_ctx.$slots, "footer", {}, () => [
									unref(showCancelButton) ? (openBlock(), createBlock(resolveDynamicComponent(unref(components).DefaultButton || unref(VbenButton)), {
										key: 0,
										variant: "ghost",
										disabled: unref(submitting),
										onClick: _cache[0] || (_cache[0] = () => __props.modalApi?.onCancel())
									}, {
										default: withCtx(() => [renderSlot(_ctx.$slots, "cancelText", {}, () => [createTextVNode(toDisplayString(unref(cancelText) || unref($t)("cancel")), 1)])]),
										_: 3
									}, 8, ["disabled"])) : createCommentVNode("v-if", true),
									renderSlot(_ctx.$slots, "center-footer"),
									unref(showConfirmButton) ? (openBlock(), createBlock(resolveDynamicComponent(unref(components).PrimaryButton || unref(VbenButton)), {
										key: 1,
										disabled: unref(confirmDisabled),
										loading: unref(confirmLoading) || unref(submitting),
										onClick: _cache[1] || (_cache[1] = () => __props.modalApi?.onConfirm())
									}, {
										default: withCtx(() => [renderSlot(_ctx.$slots, "confirmText", {}, () => [createTextVNode(toDisplayString(unref(confirmText) || unref($t)("confirm")), 1)])]),
										_: 3
									}, 8, ["disabled", "loading"])) : createCommentVNode("v-if", true)
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
					"force-mount",
					"modal",
					"open",
					"show-close",
					"animation-type",
					"z-index",
					"overlay-blur",
					"close-disabled"
				])]),
				_: 3
			}, 8, ["open"]);
		};
	}
});
//#endregion
export { _sfc_main as default };
