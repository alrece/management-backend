/* empty css                                                      */
import export_helper_default from "../../_virtual/_/plugin-vue/export-helper.mjs";
import { Fragment, TransitionGroup, computed, createBlock, createCommentVNode, createElementBlock, createElementVNode, createVNode, defineComponent, mergeModels, normalizeClass, normalizeStyle, openBlock, ref, renderList, toDisplayString, unref, useModel, vShow, withCtx, withDirectives, withModifiers } from "vue";
import { Pin, X } from "@vben-core/icons";
import { VbenContextMenu, VbenIcon } from "@vben-core/shadcn-ui";
//#region src/components/tabs-chrome/tabs.vue
const _hoisted_1 = [
	"data-active-tab",
	"data-index",
	"onClick",
	"onMousedown"
];
const _hoisted_2 = { class: "relative size-full px-1" };
const _hoisted_3 = {
	key: 0,
	class: "tabs-chrome__divider absolute top-1/2 left-(--gap) z-0 h-4 w-px translate-y-[-50%] bg-border transition-all"
};
const _hoisted_4 = { class: "tabs-chrome__extra absolute top-1/2 right-(--gap) z-3 size-4 translate-y-[-50%]" };
const _hoisted_5 = { class: "tabs-chrome__item-main z-2 mx-[calc(var(--gap)*2)] my-0 flex h-full items-center overflow-hidden rounded-tl-[5px] rounded-tr-[5px] pr-4 pl-2 text-accent-foreground duration-150 group-[.is-active]:text-primary group-[.is-active]:dark:text-accent-foreground" };
const _hoisted_6 = { class: "flex-1 overflow-hidden text-sm whitespace-nowrap" };
var tabs_default = /* @__PURE__ */ export_helper_default(/* @__PURE__ */ defineComponent({
	name: "VbenTabsChrome",
	inheritAttrs: false,
	__name: "tabs",
	props: /* @__PURE__ */ mergeModels({
		active: {},
		contentClass: { default: "vben-tabs-content" },
		contextMenus: {
			type: Function,
			default: () => []
		},
		draggable: { type: Boolean },
		gap: { default: 7 },
		maxWidth: {},
		middleClickToClose: { type: Boolean },
		minWidth: {},
		showIcon: { type: Boolean },
		styleType: {},
		tabs: { default: () => [] },
		wheelable: { type: Boolean }
	}, {
		"active": {},
		"activeModifiers": {}
	}),
	emits: /* @__PURE__ */ mergeModels(["close", "unpin"], ["update:active"]),
	setup(__props, { emit: __emit }) {
		const props = __props;
		const emit = __emit;
		const active = useModel(__props, "active");
		const contentRef = ref();
		const tabRef = ref();
		const style = computed(() => {
			const { gap } = props;
			return { "--gap": `${gap}px` };
		});
		const tabsView = computed(() => {
			return props.tabs.map((tab) => {
				const { fullPath, meta, name, path, key } = tab || {};
				const { affixTab, icon, newTabTitle, tabClosable, title } = meta || {};
				return {
					affixTab: !!affixTab,
					closable: Reflect.has(meta, "tabClosable") ? !!tabClosable : true,
					fullPath,
					icon,
					key,
					meta,
					name,
					path,
					title: newTabTitle || title || name
				};
			});
		});
		function onMouseDown(e, tab) {
			if (e.button === 1 && tab.closable && !tab.affixTab && tabsView.value.length > 1 && props.middleClickToClose) {
				e.preventDefault();
				e.stopPropagation();
				emit("close", tab.key);
			}
		}
		return (_ctx, _cache) => {
			return openBlock(), createElementBlock("div", {
				ref_key: "contentRef",
				ref: contentRef,
				class: normalizeClass([__props.contentClass, "tabs-chrome flex! h-full w-max overflow-y-hidden pr-6"]),
				style: normalizeStyle(style.value)
			}, [createVNode(TransitionGroup, { name: "slide-left" }, {
				default: withCtx(() => [(openBlock(true), createElementBlock(Fragment, null, renderList(tabsView.value, (tab, i) => {
					return openBlock(), createElementBlock("div", {
						key: tab.key,
						ref_for: true,
						ref_key: "tabRef",
						ref: tabRef,
						class: normalizeClass([[{
							"is-active": tab.key === active.value,
							draggable: !tab.affixTab,
							"affix-tab": tab.affixTab
						}], "draggable group tabs-chrome__item translate-all relative -mr-3 flex h-full items-center select-none"]),
						"data-active-tab": active.value,
						"data-index": i,
						"data-tab-item": "true",
						onClick: ($event) => active.value = tab.key,
						onMousedown: ($event) => onMouseDown($event, tab)
					}, [createVNode(unref(VbenContextMenu), {
						"handler-data": tab,
						menus: __props.contextMenus,
						modal: false,
						"item-class": "pr-6"
					}, {
						default: withCtx(() => [createElementVNode("div", _hoisted_2, [
							createCommentVNode(" divider "),
							i !== 0 && tab.key !== active.value ? (openBlock(), createElementBlock("div", _hoisted_3)) : createCommentVNode("v-if", true),
							createCommentVNode(" background "),
							_cache[0] || (_cache[0] = createElementVNode("div", { class: "tabs-chrome__background absolute z-[-1] size-full px-[calc(var(--gap)-1px)] py-0 transition-opacity duration-150" }, [
								createElementVNode("div", { class: "tabs-chrome__background-content h-full rounded-tl-(--gap) rounded-tr-(--gap) duration-150 group-[.is-active]:bg-primary/15 group-[.is-active]:dark:bg-accent" }),
								createElementVNode("svg", {
									class: "tabs-chrome__background-before absolute bottom-0 -left-px fill-transparent transition-all duration-150 group-[.is-active]:fill-primary/15 group-[.is-active]:dark:fill-accent",
									height: "7",
									width: "7"
								}, [createElementVNode("path", { d: "M 0 7 A 7 7 0 0 0 7 0 L 7 7 Z" })]),
								createElementVNode("svg", {
									class: "tabs-chrome__background-after absolute -right-px bottom-0 fill-transparent transition-all duration-150 group-[.is-active]:fill-primary/15 group-[.is-active]:dark:fill-accent",
									height: "7",
									width: "7"
								}, [createElementVNode("path", { d: "M 0 0 A 7 7 0 0 0 7 7 L 0 7 Z" })])
							], -1)),
							createCommentVNode(" extra "),
							createElementVNode("div", _hoisted_4, [
								createCommentVNode(" close-icon "),
								withDirectives(createVNode(unref(X), {
									class: "mt-0.5 size-3 cursor-pointer rounded-full stroke-accent-foreground/80 text-accent-foreground/80 transition-all group-[.is-active]:text-accent-foreground hover:bg-accent hover:stroke-accent-foreground",
									onClick: withModifiers(() => emit("close", tab.key), ["stop"])
								}, null, 8, ["onClick"]), [[vShow, !tab.affixTab && tabsView.value.length > 1 && tab.closable]]),
								withDirectives(createVNode(unref(Pin), {
									class: "mt-px size-3.5 cursor-pointer rounded-full text-accent-foreground/80 transition-all group-[.is-active]:text-accent-foreground hover:text-accent-foreground",
									onClick: withModifiers(() => emit("unpin", tab), ["stop"])
								}, null, 8, ["onClick"]), [[vShow, tab.affixTab && tabsView.value.length > 1 && tab.closable]])
							]),
							createCommentVNode(" tab-item-main "),
							createElementVNode("div", _hoisted_5, [__props.showIcon ? (openBlock(), createBlock(unref(VbenIcon), {
								key: 0,
								icon: tab.icon,
								class: "mr-1 flex size-4 items-center overflow-hidden group-hover:animate-[shrink_0.3s_ease-in-out]"
							}, null, 8, ["icon"])) : createCommentVNode("v-if", true), createElementVNode("span", _hoisted_6, toDisplayString(tab.title), 1)])
						])]),
						_: 2
					}, 1032, ["handler-data", "menus"])], 42, _hoisted_1);
				}), 128))]),
				_: 1
			})], 6);
		};
	}
}), [["__scopeId", "data-v-62fdb4ca"]]);
//#endregion
export { tabs_default as default };
