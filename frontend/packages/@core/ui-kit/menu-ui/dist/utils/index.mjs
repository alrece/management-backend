import { isVNode } from "vue";
//#region src/utils/index.ts
/**
* @zh_CN Find the parent component upward
* @param instance
* @param parentNames
*/
function findComponentUpward(instance, parentNames) {
	let parent = instance.parent;
	while (parent && !parentNames.includes(parent?.type?.name ?? "")) parent = parent.parent;
	return parent;
}
const flattedChildren = (children) => {
	const vNodes = Array.isArray(children) ? children : [children];
	const result = [];
	vNodes.forEach((child) => {
		if (Array.isArray(child)) result.push(...flattedChildren(child));
		else if (isVNode(child) && Array.isArray(child.children)) result.push(...flattedChildren(child.children));
		else {
			result.push(child);
			if (isVNode(child) && child.component?.subTree) result.push(...flattedChildren(child.component.subTree));
		}
	});
	return result;
};
//#endregion
export { findComponentUpward, flattedChildren };
