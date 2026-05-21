import { computed } from "vue";
//#region src/hooks/use-layout.ts
function useLayout(props) {
	const currentLayout = computed(() => props.isMobile ? "sidebar-nav" : props.layout);
	/**
	* 是否全屏显示content，不需要侧边、底部、顶部、tab区域
	*/
	const isFullContent = computed(() => currentLayout.value === "full-content");
	/**
	* 是否侧边混合模式
	*/
	const isSidebarMixedNav = computed(() => currentLayout.value === "sidebar-mixed-nav");
	/**
	* 是否为头部导航模式
	*/
	const isHeaderNav = computed(() => currentLayout.value === "header-nav");
	/**
	* 是否为混合导航模式
	*/
	const isMixedNav = computed(() => currentLayout.value === "mixed-nav" || currentLayout.value === "header-sidebar-nav");
	return {
		currentLayout,
		isFullContent,
		isHeaderMixedNav: computed(() => currentLayout.value === "header-mixed-nav"),
		isHeaderNav,
		isMixedNav,
		isSidebarMixedNav
	};
}
//#endregion
export { useLayout };
