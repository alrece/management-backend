import { provideAlertContext } from "./alert.mjs";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogTitle, VbenButton, VbenLoading, VbenRenderContent } from "@vben-core/shadcn-ui";
import { computed, createBlock, createCommentVNode, createElementVNode, createTextVNode, createVNode, defineComponent, h, mergeModels, nextTick, normalizeClass, openBlock, ref, resolveDynamicComponent, toDisplayString, unref, useModel, withCtx } from "vue";
import { useSimpleLocale } from "@vben-core/composables";
import { CircleAlert, CircleCheckBig, CircleHelp, CircleX, Info, X } from "@vben-core/icons";
import { globalShareState } from "@vben-core/shared/global-state";
import { cn } from "@vben-core/shared/utils";
//#region src/alert/alert.vue
const _hoisted_1 = { class: "flex items-center" };
const _hoisted_2 = { class: "flex-auto" };
const _hoisted_3 = { class: "m-4 min-h-7.5" };
const _sfc_main = /* @__PURE__ */ defineComponent({
	__name: "alert",
	props: /* @__PURE__ */ mergeModels({
		beforeClose: {},
		bordered: {
			type: Boolean,
			default: true
		},
		buttonAlign: { default: "end" },
		cancelText: {},
		centered: {
			type: Boolean,
			default: true
		},
		confirmText: {},
		containerClass: {},
		content: {},
		contentClass: {},
		contentMasking: { type: Boolean },
		footer: {},
		icon: {},
		overlayBlur: {},
		showCancel: { type: Boolean },
		title: {}
	}, {
		"open": {
			type: Boolean,
			default: false
		},
		"openModifiers": {}
	}),
	emits: /* @__PURE__ */ mergeModels([
		"closed",
		"confirm",
		"opened"
	], ["update:open"]),
	setup(__props, { emit: __emit }) {
		const props = __props;
		const emits = __emit;
		const open = useModel(__props, "open");
		const { $t } = useSimpleLocale();
		const components = globalShareState.getComponents();
		const isConfirm = ref(false);
		function onAlertClosed() {
			emits("closed", isConfirm.value);
			isConfirm.value = false;
		}
		function onEscapeKeyDown() {
			isConfirm.value = false;
		}
		const getIconRender = computed(() => {
			let iconRender = null;
			if (props.icon) {
				if (typeof props.icon === "string") switch (props.icon) {
					case "error":
						iconRender = h(CircleX, { style: { color: "hsl(var(--destructive))" } });
						break;
					case "info":
						iconRender = h(Info, { style: { color: "hsl(var(--info))" } });
						break;
					case "question":
						iconRender = CircleHelp;
						break;
					case "success":
						iconRender = h(CircleCheckBig, { style: { color: "hsl(var(--success))" } });
						break;
					case "warning":
						iconRender = h(CircleAlert, { style: { color: "hsl(var(--warning))" } });
						break;
					default:
						iconRender = null;
						break;
				}
			} else iconRender = props.icon ?? null;
			return iconRender;
		});
		function doCancel() {
			handleCancel();
			handleOpenChange(false);
		}
		function doConfirm() {
			handleConfirm();
			handleOpenChange(false);
		}
		provideAlertContext({
			doCancel,
			doConfirm
		});
		function handleConfirm() {
			isConfirm.value = true;
			emits("confirm");
		}
		function handleCancel() {
			isConfirm.value = false;
		}
		const loading = ref(false);
		async function handleOpenChange(val) {
			await nextTick();
			if (!val && props.beforeClose) {
				loading.value = true;
				try {
					if (await props.beforeClose({ isConfirm: isConfirm.value }) !== false) open.value = false;
				} finally {
					loading.value = false;
				}
			} else open.value = val;
		}
		return (_ctx, _cache) => {
			return openBlock(), createBlock(unref(AlertDialog), {
				open: open.value,
				"onUpdate:open": handleOpenChange
			}, {
				default: withCtx(() => [createVNode(unref(AlertDialogContent), {
					open: open.value,
					centered: __props.centered,
					"overlay-blur": __props.overlayBlur,
					onOpened: _cache[0] || (_cache[0] = ($event) => emits("opened")),
					onClosed: onAlertClosed,
					onEscapeKeyDown,
					class: normalizeClass(unref(cn)(__props.containerClass, "inset-x-0 mx-auto flex max-h-[80%] flex-col p-0 duration-300 sm:w-130 sm:max-w-[80%] sm:rounded-(--radius)", {
						"border border-border": __props.bordered,
						"shadow-3xl": !__props.bordered
					}))
				}, {
					default: withCtx(() => [createElementVNode("div", { class: normalizeClass(unref(cn)("relative flex-1 overflow-y-auto p-3", __props.contentClass)) }, [
						__props.title ? (openBlock(), createBlock(unref(AlertDialogTitle), { key: 0 }, {
							default: withCtx(() => [createElementVNode("div", _hoisted_1, [
								(openBlock(), createBlock(resolveDynamicComponent(getIconRender.value), { class: "mr-2" })),
								createElementVNode("span", _hoisted_2, toDisplayString(unref($t)(__props.title)), 1),
								__props.showCancel ? (openBlock(), createBlock(unref(AlertDialogCancel), {
									key: 0,
									"as-child": ""
								}, {
									default: withCtx(() => [createVNode(unref(VbenButton), {
										variant: "ghost",
										size: "icon",
										class: "rounded-full",
										disabled: loading.value,
										onClick: handleCancel
									}, {
										default: withCtx(() => [createVNode(unref(X), { class: "size-4 text-muted-foreground" })]),
										_: 1
									}, 8, ["disabled"])]),
									_: 1
								})) : createCommentVNode("v-if", true)
							])]),
							_: 1
						})) : createCommentVNode("v-if", true),
						createVNode(unref(AlertDialogDescription), null, {
							default: withCtx(() => [createElementVNode("div", _hoisted_3, [createVNode(unref(VbenRenderContent), {
								content: __props.content,
								"render-br": ""
							}, null, 8, ["content"])]), loading.value && __props.contentMasking ? (openBlock(), createBlock(unref(VbenLoading), {
								key: 0,
								spinning: loading.value
							}, null, 8, ["spinning"])) : createCommentVNode("v-if", true)]),
							_: 1
						}),
						createElementVNode("div", { class: normalizeClass(["flex items-center justify-end gap-x-2", `justify-${__props.buttonAlign}`]) }, [
							createVNode(unref(VbenRenderContent), { content: __props.footer }, null, 8, ["content"]),
							__props.showCancel ? (openBlock(), createBlock(unref(AlertDialogCancel), {
								key: 0,
								"as-child": ""
							}, {
								default: withCtx(() => [(openBlock(), createBlock(resolveDynamicComponent(unref(components).DefaultButton || unref(VbenButton)), {
									disabled: loading.value,
									variant: "ghost",
									onClick: handleCancel
								}, {
									default: withCtx(() => [createTextVNode(toDisplayString(__props.cancelText || unref($t)("cancel")), 1)]),
									_: 1
								}, 8, ["disabled"]))]),
								_: 1
							})) : createCommentVNode("v-if", true),
							createVNode(unref(AlertDialogAction), { "as-child": "" }, {
								default: withCtx(() => [(openBlock(), createBlock(resolveDynamicComponent(unref(components).PrimaryButton || unref(VbenButton)), {
									loading: loading.value,
									onClick: handleConfirm
								}, {
									default: withCtx(() => [createTextVNode(toDisplayString(__props.confirmText || unref($t)("confirm")), 1)]),
									_: 1
								}, 8, ["loading"]))]),
								_: 1
							})
						], 2)
					], 2)]),
					_: 1
				}, 8, [
					"open",
					"centered",
					"overlay-blur",
					"class"
				])]),
				_: 1
			}, 8, ["open"]);
		};
	}
});
//#endregion
export { _sfc_main as default };
