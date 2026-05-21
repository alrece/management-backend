import { breakpointsTailwind, createSharedComposable, tryOnBeforeUnmount, tryOnMounted, useBreakpoints, useCssVar, useDebounceFn, useScrollLock as useScrollLock$1 } from "@vueuse/core";
import { computed, getCurrentInstance, onMounted, onUnmounted, ref, unref, useAttrs, useSlots } from "vue";
import { CSS_VARIABLE_LAYOUT_CONTENT_HEIGHT, CSS_VARIABLE_LAYOUT_CONTENT_WIDTH, CSS_VARIABLE_LAYOUT_FOOTER_HEIGHT, CSS_VARIABLE_LAYOUT_HEADER_HEIGHT, DEFAULT_NAMESPACE } from "@vben-core/shared/constants";
import { getElementVisibleRect, getFirstNonNullOrUndefined, getScrollbarWidth, kebabToCamelCase, needsScrollbar } from "@vben-core/shared/utils";
import { useEmitAsProps, useForwardExpose, useForwardProps, useForwardPropsEmits } from "reka-ui";
//#region src/use-is-mobile.ts
function useIsMobile() {
	return { isMobile: useBreakpoints(breakpointsTailwind).smaller("md") };
}
//#endregion
//#region src/use-layout-style.ts
/**
* @zh_CN content style
*/
function useLayoutContentStyle() {
	let resizeObserver = null;
	const contentElement = ref(null);
	const visibleDomRect = ref(null);
	const contentHeight = useCssVar(CSS_VARIABLE_LAYOUT_CONTENT_HEIGHT);
	const contentWidth = useCssVar(CSS_VARIABLE_LAYOUT_CONTENT_WIDTH);
	const overlayStyle = computed(() => {
		const { height, left, top, width } = visibleDomRect.value ?? {};
		return {
			height: `${height}px`,
			left: `${left}px`,
			position: "fixed",
			top: `${top}px`,
			width: `${width}px`,
			zIndex: 150
		};
	});
	const debouncedCalcHeight = useDebounceFn((_entries) => {
		visibleDomRect.value = getElementVisibleRect(contentElement.value);
		contentHeight.value = `${visibleDomRect.value.height}px`;
		contentWidth.value = `${visibleDomRect.value.width}px`;
	}, 16);
	onMounted(() => {
		if (contentElement.value && !resizeObserver) {
			resizeObserver = new ResizeObserver(debouncedCalcHeight);
			resizeObserver.observe(contentElement.value);
		}
	});
	onUnmounted(() => {
		resizeObserver?.disconnect();
		resizeObserver = null;
	});
	return {
		contentElement,
		overlayStyle,
		visibleDomRect
	};
}
function useLayoutHeaderStyle() {
	const headerHeight = useCssVar(CSS_VARIABLE_LAYOUT_HEADER_HEIGHT);
	return {
		getLayoutHeaderHeight: () => {
			return Number.parseInt(`${headerHeight.value}`, 10);
		},
		setLayoutHeaderHeight: (height) => {
			headerHeight.value = `${height}px`;
		}
	};
}
function useLayoutFooterStyle() {
	const footerHeight = useCssVar(CSS_VARIABLE_LAYOUT_FOOTER_HEIGHT);
	return {
		getLayoutFooterHeight: () => {
			return Number.parseInt(`${footerHeight.value}`, 10);
		},
		setLayoutFooterHeight: (height) => {
			footerHeight.value = `${height}px`;
		}
	};
}
//#endregion
//#region src/use-namespace.ts
/**
* @see copy https://github.com/element-plus/element-plus/blob/dev/packages/hooks/use-namespace/index.ts
*/
const statePrefix = "is-";
const _bem = (namespace, block, blockSuffix, element, modifier) => {
	let cls = `${namespace}-${block}`;
	if (blockSuffix) cls += `-${blockSuffix}`;
	if (element) cls += `__${element}`;
	if (modifier) cls += `--${modifier}`;
	return cls;
};
const is = (name, ...args) => {
	const state = args.length > 0 ? args[0] : true;
	return name && state ? `${statePrefix}${name}` : "";
};
const useNamespace = (block) => {
	const namespace = DEFAULT_NAMESPACE;
	const b = (blockSuffix = "") => _bem(namespace, block, blockSuffix, "", "");
	const e = (element) => element ? _bem(namespace, block, "", element, "") : "";
	const m = (modifier) => modifier ? _bem(namespace, block, "", "", modifier) : "";
	const be = (blockSuffix, element) => blockSuffix && element ? _bem(namespace, block, blockSuffix, element, "") : "";
	const em = (element, modifier) => element && modifier ? _bem(namespace, block, "", element, modifier) : "";
	const bm = (blockSuffix, modifier) => blockSuffix && modifier ? _bem(namespace, block, blockSuffix, "", modifier) : "";
	const bem = (blockSuffix, element, modifier) => blockSuffix && element && modifier ? _bem(namespace, block, blockSuffix, element, modifier) : "";
	const cssVar = (object) => {
		const styles = {};
		for (const key in object) if (object[key]) styles[`--${namespace}-${key}`] = object[key];
		return styles;
	};
	const cssVarBlock = (object) => {
		const styles = {};
		for (const key in object) if (object[key]) styles[`--${namespace}-${block}-${key}`] = object[key];
		return styles;
	};
	const cssVarName = (name) => `--${namespace}-${name}`;
	const cssVarBlockName = (name) => `--${namespace}-${block}-${name}`;
	return {
		b,
		be,
		bem,
		bm,
		cssVar,
		cssVarBlock,
		cssVarBlockName,
		cssVarName,
		e,
		em,
		is,
		m,
		namespace
	};
};
//#endregion
//#region src/use-priority-value.ts
/**
* 依次从插槽、attrs、props、state 中获取值
* @param key
* @param props
* @param state
*/
function usePriorityValue(key, props, state) {
	const instance = getCurrentInstance();
	const slots = useSlots();
	const attrs = useAttrs();
	return computed(() => {
		const rawProps = instance?.vnode?.props || {};
		const standardRawProps = {};
		for (const [key, value] of Object.entries(rawProps)) standardRawProps[kebabToCamelCase(key)] = value;
		const propsKey = standardRawProps?.[key] === void 0 ? void 0 : props[key];
		return getFirstNonNullOrUndefined(slots[key], attrs[key], propsKey, state?.value?.[key]);
	});
}
/**
* 批量获取state中的值（每个值都是ref）
* @param props
* @param state
*/
function usePriorityValues(props, state) {
	const result = {};
	Object.keys(props).forEach((key) => {
		result[key] = usePriorityValue(key, props, state);
	});
	return result;
}
/**
* 批量获取state中的值（集中在一个computed，用于透传）
* @param props
* @param state
*/
function useForwardPriorityValues(props, state) {
	const computedResult = {};
	Object.keys(props).forEach((key) => {
		computedResult[key] = usePriorityValue(key, props, state);
	});
	return computed(() => {
		const unwrapResult = {};
		Object.keys(props).forEach((key) => {
			unwrapResult[key] = unref(computedResult[key]);
		});
		return unwrapResult;
	});
}
//#endregion
//#region src/use-scroll-lock.ts
const SCROLL_FIXED_CLASS = `_scroll__fixed_`;
function useScrollLock() {
	const isLocked = useScrollLock$1(document.body);
	const scrollbarWidth = getScrollbarWidth();
	tryOnMounted(() => {
		if (!needsScrollbar()) return;
		document.body.style.paddingRight = `${scrollbarWidth}px`;
		const nodes = [...document.querySelectorAll(`.${SCROLL_FIXED_CLASS}`)];
		if (nodes.length > 0) nodes.forEach((node) => {
			node.dataset.transition = node.style.transition;
			node.style.transition = "none";
			node.style.paddingRight = `${scrollbarWidth}px`;
		});
		isLocked.value = true;
	});
	tryOnBeforeUnmount(() => {
		if (!needsScrollbar()) return;
		isLocked.value = false;
		const nodes = [...document.querySelectorAll(`.${SCROLL_FIXED_CLASS}`)];
		if (nodes.length > 0) nodes.forEach((node) => {
			node.style.paddingRight = "";
			requestAnimationFrame(() => {
				node.style.transition = node.dataset.transition || "";
			});
		});
		document.body.style.paddingRight = "";
	});
}
//#endregion
//#region src/use-simple-locale/messages.ts
const messages = {
	"en-US": {
		cancel: "Cancel",
		collapse: "Collapse",
		confirm: "Confirm",
		expand: "Expand",
		prompt: "Prompt",
		reset: "Reset",
		submit: "Submit"
	},
	"zh-CN": {
		cancel: "取消",
		collapse: "收起",
		confirm: "确认",
		expand: "展开",
		prompt: "提示",
		reset: "重置",
		submit: "提交"
	}
};
const getMessages = (locale) => messages[locale];
//#endregion
//#region src/use-simple-locale/index.ts
const useSimpleLocale = createSharedComposable(() => {
	const currentLocale = ref("zh-CN");
	const setSimpleLocale = (locale) => {
		currentLocale.value = locale;
	};
	return {
		$t: computed(() => {
			const localeMessages = getMessages(currentLocale.value);
			return (key) => {
				return localeMessages[key] || key;
			};
		}),
		currentLocale,
		setSimpleLocale
	};
});
//#endregion
//#region src/use-sortable.ts
function useSortable(sortableContainer, options = {}) {
	const initializeSortable = async () => {
		return (await import("sortablejs/modular/sortable.complete.esm.js"))?.default?.create?.(sortableContainer, {
			animation: 300,
			delay: 400,
			delayOnTouchOnly: true,
			...options
		});
	};
	return { initializeSortable };
}
//#endregion
export { SCROLL_FIXED_CLASS, useEmitAsProps, useForwardExpose, useForwardPriorityValues, useForwardProps, useForwardPropsEmits, useIsMobile, useLayoutContentStyle, useLayoutFooterStyle, useLayoutHeaderStyle, useNamespace, usePriorityValue, usePriorityValues, useScrollLock, useSimpleLocale, useSortable };
