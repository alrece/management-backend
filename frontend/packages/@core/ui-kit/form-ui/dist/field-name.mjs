//#region src/field-name.ts
function resolveFieldNamePath(fieldName) {
	if (fieldName.startsWith("[") && fieldName.endsWith("]")) {
		const rawKey = fieldName.slice(1, -1);
		return {
			pathSegments: [rawKey],
			rawKey
		};
	}
	return {
		pathSegments: fieldName.match(/[^.[\]]+/g) ?? [],
		rawKey: void 0
	};
}
//#endregion
export { resolveFieldNamePath };
