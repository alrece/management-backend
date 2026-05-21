import { createBlock, createCommentVNode, createElementBlock, createVNode, defineComponent, normalizeClass, openBlock, renderSlot, unref, withCtx } from "vue";
import { FormLabel, VbenHelpTooltip, VbenRenderContent } from "@vben-core/shadcn-ui";
import { cn } from "@vben-core/shared/utils";
//#region src/form-render/form-label.vue
const _hoisted_1 = {
	key: 0,
	class: "mr-0.5 text-destructive"
};
const _hoisted_2 = {
	key: 2,
	class: "ml-0.5"
};
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "form-label",
	props: {
		class: {},
		colon: { type: Boolean },
		help: { type: [Function, String] },
		label: { type: [Function, String] },
		required: { type: Boolean }
	},
	setup(__props) {
		const props = __props;
		return (_ctx, _cache) => {
			return openBlock(), createBlock(unref(FormLabel), { class: normalizeClass(unref(cn)("flex items-center", props.class)) }, {
				default: withCtx(() => [
					__props.required ? (openBlock(), createElementBlock("span", _hoisted_1, "*")) : createCommentVNode("v-if", true),
					renderSlot(_ctx.$slots, "default"),
					__props.help ? (openBlock(), createBlock(unref(VbenHelpTooltip), {
						key: 1,
						"trigger-class": "size-3.5 ml-1"
					}, {
						default: withCtx(() => [createVNode(unref(VbenRenderContent), { content: __props.help }, null, 8, ["content"])]),
						_: 1
					})) : createCommentVNode("v-if", true),
					__props.colon && __props.label ? (openBlock(), createElementBlock("span", _hoisted_2, ":")) : createCommentVNode("v-if", true)
				]),
				_: 3
			}, 8, ["class"]);
		};
	}
});
//#endregion
export { _sfc_main as default };
