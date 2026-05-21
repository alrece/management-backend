import { Transition, createBlock, defineComponent, mergeProps, openBlock, renderSlot, toHandlers, withCtx } from "vue";
//#region src/components/collapse-transition.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "CollapseTransition",
	__name: "collapse-transition",
	setup(__props) {
		const reset = (el) => {
			el.style.maxHeight = "";
			el.style.overflow = el.dataset.oldOverflow;
			el.style.paddingTop = el.dataset.oldPaddingTop;
			el.style.paddingBottom = el.dataset.oldPaddingBottom;
		};
		const on = {
			afterEnter(el) {
				el.style.maxHeight = "";
				el.style.overflow = el.dataset.oldOverflow;
			},
			afterLeave(el) {
				reset(el);
			},
			beforeEnter(el) {
				if (!el.dataset) el.dataset = {};
				el.dataset.oldPaddingTop = el.style.paddingTop;
				el.dataset.oldMarginTop = el.style.marginTop;
				el.dataset.oldPaddingBottom = el.style.paddingBottom;
				el.dataset.oldMarginBottom = el.style.marginBottom;
				if (el.style.height) el.dataset.elExistsHeight = el.style.height;
				el.style.maxHeight = 0;
				el.style.paddingTop = 0;
				el.style.marginTop = 0;
				el.style.paddingBottom = 0;
				el.style.marginBottom = 0;
			},
			beforeLeave(el) {
				if (!el.dataset) el.dataset = {};
				el.dataset.oldPaddingTop = el.style.paddingTop;
				el.dataset.oldMarginTop = el.style.marginTop;
				el.dataset.oldPaddingBottom = el.style.paddingBottom;
				el.dataset.oldMarginBottom = el.style.marginBottom;
				el.dataset.oldOverflow = el.style.overflow;
				el.style.maxHeight = `${el.scrollHeight}px`;
				el.style.overflow = "hidden";
			},
			enter(el) {
				requestAnimationFrame(() => {
					el.dataset.oldOverflow = el.style.overflow;
					if (el.dataset.elExistsHeight) el.style.maxHeight = el.dataset.elExistsHeight;
					else if (el.scrollHeight === 0) el.style.maxHeight = 0;
					else el.style.maxHeight = `${el.scrollHeight}px`;
					el.style.paddingTop = el.dataset.oldPaddingTop;
					el.style.paddingBottom = el.dataset.oldPaddingBottom;
					el.style.marginTop = el.dataset.oldMarginTop;
					el.style.marginBottom = el.dataset.oldMarginBottom;
					el.style.overflow = "hidden";
				});
			},
			enterCancelled(el) {
				reset(el);
			},
			leave(el) {
				if (el.scrollHeight !== 0) {
					el.style.maxHeight = 0;
					el.style.paddingTop = 0;
					el.style.paddingBottom = 0;
					el.style.marginTop = 0;
					el.style.marginBottom = 0;
				}
			},
			leaveCancelled(el) {
				reset(el);
			}
		};
		return (_ctx, _cache) => {
			return openBlock(), createBlock(Transition, mergeProps({ name: "collapse-transition" }, toHandlers(on)), {
				default: withCtx(() => [renderSlot(_ctx.$slots, "default")]),
				_: 3
			}, 16);
		};
	}
});
//#endregion
export { _sfc_main as default };
