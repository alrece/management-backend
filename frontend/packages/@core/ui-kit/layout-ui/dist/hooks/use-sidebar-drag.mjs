import { ref } from "vue";
//#region src/hooks/use-sidebar-drag.ts
function useSidebarDrag() {
	const isDragging = ref(false);
	let cleanup = null;
	let dragOverlay = null;
	const startDrag = (e, options, elements, onDrag) => {
		const { min, max } = options;
		const { dragBar, target } = elements;
		if (isDragging.value || !dragBar || !target) return;
		e.preventDefault();
		e.stopPropagation();
		isDragging.value = true;
		const startX = e.clientX;
		const startWidth = target.getBoundingClientRect().width;
		const startLeft = dragBar.offsetLeft;
		dragBar.classList.add("bg-primary");
		dragBar.classList.remove("bg-primary/30");
		const dragBarTransition = dragBar.style.transition;
		const targetTransition = target.style.transition;
		dragBar.style.transition = "none";
		target.style.transition = "none";
		dragOverlay = document.createElement("div");
		dragOverlay.style.position = "fixed";
		dragOverlay.style.inset = "0";
		dragOverlay.style.zIndex = "9999";
		dragOverlay.style.cursor = "col-resize";
		dragOverlay.style.userSelect = "none";
		dragOverlay.style.outline = "none";
		dragOverlay.tabIndex = -1;
		dragOverlay.style.background = "rgba(0,0,0,0)";
		document.body.append(dragOverlay);
		const onMouseMove = (moveEvent) => {
			if (!isDragging.value || !dragBar || !target) {
				endDrag();
				return;
			}
			let currentWidth = startWidth + (moveEvent.clientX - startX);
			const isOutOfMin = currentWidth < min;
			const isOutOfMax = currentWidth > max;
			const isOutOfBounds = isOutOfMin || isOutOfMax;
			if (isOutOfMin) currentWidth = min;
			if (isOutOfMax) currentWidth = max;
			const newLeft = startLeft + (currentWidth - startWidth);
			if (dragOverlay) dragOverlay.style.cursor = isOutOfBounds ? "not-allowed" : "col-resize";
			dragBar.style.left = `${newLeft}px`;
			if (isOutOfBounds) {
				dragBar.classList.add("bg-primary/30");
				dragBar.classList.remove("bg-primary");
			} else {
				dragBar.classList.add("bg-primary");
				dragBar.classList.remove("bg-primary/30");
			}
		};
		const onMouseUp = (upEvent) => {
			if (!isDragging.value || !dragBar || !target) {
				endDrag();
				return;
			}
			let newWidth = startWidth + (upEvent.clientX - startX);
			newWidth = Math.min(max, Math.max(min, newWidth));
			dragBar.classList.remove("bg-primary", "bg-primary/30");
			try {
				onDrag?.(Math.round(newWidth));
			} finally {
				endDrag();
			}
		};
		document.addEventListener("mousemove", onMouseMove);
		document.addEventListener("mouseup", onMouseUp);
		cleanup = () => {
			if (!cleanup) return;
			document.removeEventListener("mousemove", onMouseMove);
			document.removeEventListener("mouseup", onMouseUp);
			if (dragBar) {
				dragBar.style.transition = dragBarTransition;
				dragBar.style.left = "";
				dragBar.classList.remove("bg-primary", "bg-primary/30");
			}
			if (target) target.style.transition = targetTransition;
			if (dragOverlay) {
				dragOverlay.remove();
				dragOverlay = null;
			}
			isDragging.value = false;
			cleanup = null;
		};
	};
	const endDrag = () => {
		cleanup?.();
	};
	return {
		startDrag,
		endDrag,
		get isDragging() {
			return isDragging.value;
		}
	};
}
//#endregion
export { useSidebarDrag };
