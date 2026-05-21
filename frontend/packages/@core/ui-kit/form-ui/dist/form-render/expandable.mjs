import { computed, nextTick, onMounted, ref, useTemplateRef, watch } from "vue";
import { breakpointsTailwind, useBreakpoints, useElementVisibility } from "@vueuse/core";
//#region src/form-render/expandable.ts
/**
* 动态计算行数
*/
function useExpandable(props) {
	const wrapperRef = useTemplateRef("wrapperRef");
	const isVisible = useElementVisibility(wrapperRef);
	const rowMapping = ref({});
	const isCalculated = ref(false);
	const breakpoints = useBreakpoints(breakpointsTailwind);
	const keepFormItemIndex = computed(() => {
		const rows = props.collapsedRows ?? 1;
		const mapping = rowMapping.value;
		let maxItem = 0;
		for (let index = 1; index <= rows; index++) maxItem += mapping?.[index] ?? 0;
		return maxItem - 1 || 1;
	});
	watch([
		() => props.showCollapseButton,
		() => breakpoints.active().value,
		() => props.schema?.length,
		() => isVisible.value
	], async ([val]) => {
		if (val) {
			await nextTick();
			rowMapping.value = {};
			isCalculated.value = false;
			await calculateRowMapping();
		}
	});
	async function calculateRowMapping() {
		if (!props.showCollapseButton) return;
		await nextTick();
		if (!wrapperRef.value) return;
		const formItems = [...wrapperRef.value.children];
		const container = wrapperRef.value;
		const rowHeights = window.getComputedStyle(container).getPropertyValue("grid-template-rows").split(" ");
		const containerRect = container?.getBoundingClientRect();
		formItems.forEach((el) => {
			const itemTop = el.getBoundingClientRect().top - containerRect.top;
			let rowStart = 0;
			let cumulativeHeight = 0;
			for (const [i, rowHeight] of rowHeights.entries()) {
				cumulativeHeight += Number.parseFloat(rowHeight);
				if (itemTop < cumulativeHeight) {
					rowStart = i + 1;
					break;
				}
			}
			if (rowStart > (props?.collapsedRows ?? 1)) return;
			rowMapping.value[rowStart] = (rowMapping.value[rowStart] ?? 0) + 1;
			isCalculated.value = true;
		});
	}
	onMounted(() => {
		calculateRowMapping();
	});
	return {
		isCalculated,
		keepFormItemIndex,
		wrapperRef
	};
}
//#endregion
export { useExpandable };
