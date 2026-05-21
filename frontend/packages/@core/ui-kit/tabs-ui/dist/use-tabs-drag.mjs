import { nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useIsMobile, useSortable } from "@vben-core/composables";
//#region src/use-tabs-drag.ts
function findParentElement(element) {
	const parentCls = "group";
	return element.classList.contains(parentCls) ? element : element.closest(`.${parentCls}`);
}
function useTabsDrag(props, emit) {
	const sortableInstance = ref(null);
	async function initTabsSortable() {
		await nextTick();
		const el = document.querySelectorAll(`.${props.contentClass}`)?.[0];
		if (!el) {
			console.warn("Element not found for sortable initialization");
			return;
		}
		const resetElState = async () => {
			el.style.cursor = "default";
			el.querySelector(".draggable")?.classList.remove("dragging");
		};
		const { initializeSortable } = useSortable(el, {
			filter: (_evt, target) => {
				return !findParentElement(target)?.classList.contains("draggable") || !props.draggable;
			},
			onEnd(evt) {
				const { newIndex, oldIndex } = evt;
				const { srcElement } = evt.originalEvent;
				if (!srcElement) {
					resetElState();
					return;
				}
				const srcParent = findParentElement(srcElement);
				if (!srcParent) {
					resetElState();
					return;
				}
				if (!srcParent.classList.contains("draggable")) {
					resetElState();
					return;
				}
				if (oldIndex !== void 0 && newIndex !== void 0 && !Number.isNaN(oldIndex) && !Number.isNaN(newIndex) && oldIndex !== newIndex) emit("sortTabs", oldIndex, newIndex);
				resetElState();
			},
			onMove(evt) {
				if (findParentElement(evt.related)?.classList.contains("draggable") && props.draggable) return evt.dragged.classList.contains("affix-tab") === evt.related.classList.contains("affix-tab");
				else return false;
			},
			onStart: () => {
				el.style.cursor = "grabbing";
				el.querySelector(".draggable")?.classList.add("dragging");
			}
		});
		sortableInstance.value = await initializeSortable();
	}
	async function init() {
		const { isMobile } = useIsMobile();
		if (isMobile.value) return;
		await nextTick();
		initTabsSortable();
	}
	onMounted(init);
	watch(() => props.styleType, () => {
		sortableInstance.value?.destroy();
		init();
	});
	onUnmounted(() => {
		sortableInstance.value?.destroy();
	});
}
//#endregion
export { useTabsDrag };
