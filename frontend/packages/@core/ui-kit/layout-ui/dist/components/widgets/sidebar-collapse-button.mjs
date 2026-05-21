import { createBlock, createElementBlock, defineComponent, openBlock, unref, useModel, withModifiers } from "vue";
import { ChevronsLeft, ChevronsRight } from "@vben-core/icons";
//#region src/components/widgets/sidebar-collapse-button.vue
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "sidebar-collapse-button",
	props: {
		"collapsed": { type: Boolean },
		"collapsedModifiers": {}
	},
	emits: ["update:collapsed"],
	setup(__props) {
		const collapsed = useModel(__props, "collapsed");
		function handleCollapsed() {
			collapsed.value = !collapsed.value;
		}
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("div", {
				class: "absolute bottom-2 left-3 z-10 flex-center cursor-pointer rounded-sm bg-accent p-1 text-foreground/60 hover:bg-accent-hover hover:text-foreground",
				onClick: withModifiers(handleCollapsed, ["stop"])
			}, [collapsed.value ? (openBlock(), createBlock(unref(ChevronsRight), {
				key: 0,
				class: "size-4"
			})) : (openBlock(), createBlock(unref(ChevronsLeft), {
				key: 1,
				class: "size-4"
			}))]);
		};
	}
});
//#endregion
export { _sfc_main as default };
