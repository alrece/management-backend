import { findComponentUpward } from "../utils/index.mjs";
import { computed, getCurrentInstance } from "vue";
//#region src/hooks/use-menu.ts
function useMenu() {
	const instance = getCurrentInstance();
	if (!instance) throw new Error("instance is required");
	/**
	* @zh_CN 获取所有父级菜单链路
	*/
	const parentPaths = computed(() => {
		let parent = instance.parent;
		const paths = [instance.props.path];
		while (parent?.type.name !== "Menu") {
			if (parent?.props.path) paths.unshift(parent.props.path);
			parent = parent?.parent ?? null;
		}
		return paths;
	});
	return {
		parentMenu: computed(() => {
			return findComponentUpward(instance, ["Menu", "SubMenu"]);
		}),
		parentPaths
	};
}
function useMenuStyle(menu) {
	return computed(() => {
		return { "--menu-level": menu ? menu?.level ?? 1 : 0 };
	});
}
//#endregion
export { useMenu, useMenuStyle };
