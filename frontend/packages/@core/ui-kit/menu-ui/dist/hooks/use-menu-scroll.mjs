import { watch } from "vue";
import { useDebounceFn } from "@vueuse/core";
//#region src/hooks/use-menu-scroll.ts
function useMenuScroll(activePath, options = {}) {
	const { enable = true, delay = 320 } = options;
	function scrollToActiveItem() {
		if (!(typeof enable === "boolean" ? enable : enable.value)) return;
		const activeElement = document.querySelector(`aside li[role=menuitem].is-active`);
		if (activeElement) activeElement.scrollIntoView({
			behavior: "smooth",
			block: "center",
			inline: "center"
		});
	}
	const debouncedScroll = useDebounceFn(scrollToActiveItem, delay);
	watch(activePath, () => {
		if (!(typeof enable === "boolean" ? enable : enable.value)) return;
		debouncedScroll();
	});
	return { scrollToActiveItem };
}
//#endregion
export { useMenuScroll };
