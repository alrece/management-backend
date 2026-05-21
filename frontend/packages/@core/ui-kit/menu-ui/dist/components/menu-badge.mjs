import _sfc_main$1 from "./menu-badge-dot.mjs";
import { computed, createBlock, createCommentVNode, createElementBlock, defineComponent, normalizeClass, normalizeStyle, openBlock, toDisplayString } from "vue";
import { isValidColor } from "@vben-core/shared/color";
//#region src/components/menu-badge.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "menu-badge",
	props: {
		hasChildren: { type: Boolean },
		badge: {},
		badgeType: {},
		badgeVariants: {}
	},
	setup(__props) {
		const props = __props;
		const variantsMap = {
			default: "bg-green-500",
			destructive: "bg-destructive",
			primary: "bg-primary",
			success: "bg-green-500",
			warning: "bg-yellow-500"
		};
		const isDot = computed(() => props.badgeType === "dot");
		const badgeClass = computed(() => {
			const { badgeVariants } = props;
			if (!badgeVariants) return variantsMap.default;
			return variantsMap[badgeVariants] || badgeVariants;
		});
		const badgeStyle = computed(() => {
			if (badgeClass.value && isValidColor(badgeClass.value)) return { backgroundColor: badgeClass.value };
			return {};
		});
		return (_ctx, _cache) => {
			return isDot.value || __props.badge ? (openBlock(), createElementBlock("span", {
				key: 0,
				class: normalizeClass([_ctx.$attrs.class, "absolute"])
			}, [isDot.value ? (openBlock(), createBlock(_sfc_main$1, {
				key: 0,
				"dot-class": badgeClass.value,
				"dot-style": badgeStyle.value
			}, null, 8, ["dot-class", "dot-style"])) : (openBlock(), createElementBlock("div", {
				key: 1,
				class: normalizeClass([badgeClass.value, "flex-center rounded-xl px-1.5 py-0.5 text-[10px] text-primary-foreground"]),
				style: normalizeStyle(badgeStyle.value)
			}, toDisplayString(__props.badge), 7))], 2)) : createCommentVNode("v-if", true);
		};
	}
});
//#endregion
export { _sfc_main as default };
