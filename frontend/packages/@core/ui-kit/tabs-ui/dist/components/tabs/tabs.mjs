import { Fragment, TransitionGroup, computed, createBlock, createCommentVNode, createElementBlock, createElementVNode, createVNode, defineComponent, mergeModels, normalizeClass, openBlock, renderList, toDisplayString, unref, useModel, vShow, withCtx, withDirectives, withModifiers } from "vue";
import { Pin, X } from "@vben-core/icons";
import { VbenContextMenu, VbenIcon } from "@vben-core/shadcn-ui";
//#region src/components/tabs/tabs.vue
const _hoisted_1 = [
	"data-index",
	"onClick",
	"onMousedown"
];
const _hoisted_2 = { class: "relative flex size-full items-center" };
const _hoisted_3 = { class: "absolute top-1/2 right-1.5 z-3 translate-y-[-50%] overflow-hidden" };
const _hoisted_4 = { class: "mx-3 mr-4 flex h-full items-center overflow-hidden rounded-tl-[5px] rounded-tr-[5px] pr-3 text-accent-foreground transition-all duration-300 group-[.is-active]:text-primary group-[.is-active]:dark:text-accent-foreground" };
const _hoisted_5 = { class: "flex-1 overflow-hidden text-sm whitespace-nowrap" };
const _sfc_main = /* @__PURE__ */ defineComponent({
	name: "VbenTabs",
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
		gap: {},
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
		const typeWithClass = computed(() => {
			return {
				brisk: { content: `h-full after:content-['']  after:absolute after:bottom-0 after:left-0 after:w-full after:h-[1.5px] after:bg-primary after:scale-x-0 after:transition-[transform] after:ease-out after:duration-300 hover:after:scale-x-100 after:origin-left [&.is-active]:after:scale-x-100 [&:not(:first-child)]:border-l last:border-r last:border-r border-border` },
				card: { content: "h-[calc(100%-6px)] rounded-md ml-2 border border-border  transition-all" },
				plain: { content: "h-full [&:not(:first-child)]:border-l last:border-r border-border" }
			}[props.styleType || "plain"] || { content: "" };
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
			return openBlock(), createElementBlock("div", { class: normalizeClass([__props.contentClass, "relative flex! h-full w-max items-center overflow-hidden pr-6"]) }, [createVNode(TransitionGroup, { name: "slide-left" }, {
				default: withCtx(() => [(openBlock(true), createElementBlock(Fragment, null, renderList(tabsView.value, (tab, i) => {
					return openBlock(), createElementBlock("div", {
						key: tab.key,
						class: normalizeClass([[{
							"is-active bg-primary/15 dark:bg-accent": tab.key === active.value,
							draggable: !tab.affixTab,
							"affix-tab": tab.affixTab
						}, typeWithClass.value.content], "group tab-item translate-all relative flex cursor-pointer select-none [&:not(.is-active)]:hover:bg-accent"]),
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
							createCommentVNode(" extra "),
							createElementVNode("div", _hoisted_3, [
								createCommentVNode(" close-icon "),
								withDirectives(createVNode(unref(X), {
									class: "size-3 cursor-pointer rounded-full stroke-accent-foreground/80 transition-all group-[.is-active]:text-primary hover:bg-accent hover:stroke-accent-foreground group-[.is-active]:dark:text-accent-foreground",
									onClick: withModifiers(() => emit("close", tab.key), ["stop"])
								}, null, 8, ["onClick"]), [[vShow, !tab.affixTab && tabsView.value.length > 1 && tab.closable]]),
								withDirectives(createVNode(unref(Pin), {
									class: "mt-px size-3.5 cursor-pointer rounded-full transition-all group-[.is-active]:text-primary hover:bg-accent hover:stroke-accent-foreground group-[.is-active]:dark:text-accent-foreground",
									onClick: withModifiers(() => emit("unpin", tab), ["stop"])
								}, null, 8, ["onClick"]), [[vShow, tab.affixTab && tabsView.value.length > 1 && tab.closable]])
							]),
							createCommentVNode(" tab-item-main "),
							createElementVNode("div", _hoisted_4, [__props.showIcon ? (openBlock(), createBlock(unref(VbenIcon), {
								key: 0,
								icon: tab.icon,
								class: "mr-2 flex size-4 items-center overflow-hidden group-hover:animate-[shrink_0.3s_ease-in-out]",
								fallback: ""
							}, null, 8, ["icon"])) : createCommentVNode("v-if", true), createElementVNode("span", _hoisted_5, toDisplayString(tab.title), 1)])
						])]),
						_: 2
					}, 1032, ["handler-data", "menus"])], 42, _hoisted_1);
				}), 128))]),
				_: 1
			})], 2);
		};
	}
});
//#endregion
export { _sfc_main as default };
