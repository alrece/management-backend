import { resolveFieldNamePath } from "./field-name.mjs";
import { isRef, toRaw } from "vue";
import { Store } from "@vben-core/shared/store";
import { StateHandler, bindMethods, cloneDeep, createMerge, formatDate, get, isDate, isDayjsObject, isFunction, isObject, mergeWithArrayOverride, set } from "@vben-core/shared/utils";
//#region src/form-api.ts
function getDefaultState() {
	return {
		actionWrapperClass: "",
		collapsed: false,
		collapsedRows: 1,
		collapseTriggerResize: false,
		commonConfig: {},
		handleReset: void 0,
		handleSubmit: void 0,
		handleValuesChange: void 0,
		handleCollapsedChange: void 0,
		layout: "horizontal",
		resetButtonOptions: {},
		schema: [],
		scrollToFirstError: false,
		showCollapseButton: false,
		showDefaultActions: true,
		submitButtonOptions: {},
		submitOnChange: false,
		submitOnEnter: false,
		wrapperClass: "grid-cols-1"
	};
}
var FormApi = class {
	form = {};
	isMounted = false;
	state = null;
	stateHandler;
	store;
	/**
	* 组件实例映射
	*/
	componentRefMap = /* @__PURE__ */ new Map();
	latestSubmissionValues = null;
	prevState = null;
	constructor(options = {}) {
		const { ...storeState } = options;
		this.store = new Store({
			...getDefaultState(),
			...storeState
		});
		this.store.subscribe((state) => {
			this.prevState = this.state;
			this.state = state;
			this.updateState();
		});
		this.state = this.store.state;
		this.stateHandler = new StateHandler();
		bindMethods(this);
	}
	/**
	* 获取字段组件实例
	* @param fieldName 字段名
	* @returns 组件实例
	*/
	getFieldComponentRef(fieldName) {
		let target = this.componentRefMap.has(fieldName) ? this.componentRefMap.get(fieldName) : void 0;
		if (target && target.$.type.name === "AsyncComponentWrapper" && target.$.subTree.ref) {
			if (Array.isArray(target.$.subTree.ref)) {
				if (target.$.subTree.ref.length > 0 && isRef(target.$.subTree.ref[0]?.r)) target = target.$.subTree.ref[0]?.r.value;
			} else if (isRef(target.$.subTree.ref.r)) target = target.$.subTree.ref.r.value;
		}
		return target;
	}
	/**
	* 获取当前聚焦的字段，如果没有聚焦的字段则返回undefined
	*/
	getFocusedField() {
		for (const fieldName of this.componentRefMap.keys()) {
			const ref = this.getFieldComponentRef(fieldName);
			if (ref) {
				let el = null;
				if (ref instanceof HTMLElement) el = ref;
				else if (ref.$el instanceof HTMLElement) el = ref.$el;
				if (!el) continue;
				if (el === document.activeElement || el.contains(document.activeElement)) return fieldName;
			}
		}
	}
	getLatestSubmissionValues() {
		return this.latestSubmissionValues || {};
	}
	getState() {
		return this.state;
	}
	async getValues() {
		const form = await this.getForm();
		const values = form.values ? this.handleRangeTimeValue(cloneDeep(toRaw(form.values))) : {};
		return this.handleValueFormat(values);
	}
	async isFieldValid(fieldName) {
		return (await this.getForm()).isFieldValid(fieldName);
	}
	merge(formApi) {
		const chain = [this, formApi];
		const proxy = new Proxy(formApi, { get(target, prop) {
			if (prop === "merge") return (nextFormApi) => {
				chain.push(nextFormApi);
				return proxy;
			};
			if (prop === "submitAllForm") return async (needMerge = true) => {
				try {
					const results = await Promise.all(chain.map(async (api) => {
						if (!(await api.validate()).valid) return;
						return toRaw(await api.getValues() || {});
					}));
					if (needMerge) return Object.assign({}, ...results);
					return results;
				} catch (error) {
					console.error("Validation error:", error);
				}
			};
			return target[prop];
		} });
		return proxy;
	}
	mount(formActions, componentRefMap) {
		if (!this.isMounted) {
			Object.assign(this.form, formActions);
			this.stateHandler.setConditionTrue();
			const initialValues = this.form.values ? this.handleRangeTimeValue(cloneDeep(toRaw(this.form.values))) : {};
			this.setLatestSubmissionValues({ ...this.handleValueFormat(initialValues) });
			this.componentRefMap = componentRefMap ?? this.componentRefMap ?? /* @__PURE__ */ new Map();
			this.isMounted = true;
		}
	}
	/**
	* 根据字段名移除表单项
	* @param fields
	*/
	async removeSchemaByFields(fields) {
		const fieldSet = new Set(fields);
		const filterSchema = (this.state?.schema ?? []).filter((item) => !fieldSet.has(item.fieldName));
		this.setState({ schema: filterSchema });
	}
	/**
	* 重置表单
	*/
	async resetForm(state, opts) {
		return (await this.getForm()).resetForm(state, opts);
	}
	async resetValidate() {
		const form = await this.getForm();
		Object.keys(form.errors.value).forEach((field) => {
			form.setFieldError(field, void 0);
		});
	}
	/**
	* 滚动到第一个错误字段
	* @param errors 验证错误对象
	*/
	scrollToFirstError(errors) {
		const firstErrorFieldName = typeof errors === "string" ? errors : Object.keys(errors)[0];
		if (!firstErrorFieldName) return;
		let el = document.querySelector(`[name="${firstErrorFieldName}"]`);
		if (!el) {
			const componentRef = this.getFieldComponentRef(firstErrorFieldName);
			if (componentRef && componentRef.$el instanceof HTMLElement) el = componentRef.$el;
		}
		if (el) el.scrollIntoView({
			behavior: "smooth",
			block: "center",
			inline: "nearest"
		});
	}
	/**
	* 设置表单禁用状态：用于非 Modal 中使用 Form 时，需要 Form 自己控制禁用状态
	* @author 芋道源码
	* @param disabled 是否禁用
	*/
	setDisabled(disabled) {
		this.setState((prev) => ({
			...prev,
			commonConfig: {
				...prev.commonConfig,
				disabled
			}
		}));
	}
	async setFieldValue(field, value, shouldValidate) {
		(await this.getForm()).setFieldValue(field, value, shouldValidate);
	}
	setLatestSubmissionValues(values) {
		this.latestSubmissionValues = { ...toRaw(values) };
	}
	/**
	* 设置表单提交按钮的加载状态：用于非 Modal 中使用 Form 时，需要 Form 自己控制 loading 状态
	* @author 芋道源码
	* @param loading 是否加载中
	*/
	setLoading(loading) {
		this.setState((prev) => ({
			...prev,
			submitButtonOptions: {
				...prev.submitButtonOptions,
				loading
			}
		}));
	}
	setState(stateOrFn) {
		if (isFunction(stateOrFn)) this.store.setState((prev) => {
			return mergeWithArrayOverride(stateOrFn(prev), prev);
		});
		else this.store.setState((prev) => mergeWithArrayOverride(stateOrFn, prev));
	}
	/**
	* 设置表单值
	* @param fields record
	* @param filterFields 过滤不在schema中定义的字段 默认为true
	* @param shouldValidate
	*/
	async setValues(fields, filterFields = true, shouldValidate = false) {
		const form = await this.getForm();
		if (!filterFields) {
			form.setValues(fields, shouldValidate);
			return;
		}
		/**
		* 合并算法有待改进，目前的算法不支持object类型的值。
		* antd的日期时间相关组件的值类型为dayjs对象
		* element-plus的日期时间相关组件的值类型可能为Date对象
		* 以上两种类型需要排除深度合并
		*/
		const fieldMergeFn = createMerge((obj, key, value) => {
			if (key in obj) obj[key] = !Array.isArray(obj[key]) && isObject(obj[key]) && !isDayjsObject(obj[key]) && !isDate(obj[key]) ? fieldMergeFn(value, obj[key]) : value;
			return true;
		});
		const filteredFields = fieldMergeFn(fields, form.values);
		form.setValues(filteredFields, shouldValidate);
	}
	async submitForm(e) {
		e?.preventDefault();
		e?.stopPropagation();
		await (await this.getForm()).submitForm();
		const rawValues = toRaw(await this.getValues());
		await this.state?.handleSubmit?.(rawValues);
		return rawValues;
	}
	unmount() {
		this.form?.resetForm?.();
		this.componentRefMap = /* @__PURE__ */ new Map();
		this.latestSubmissionValues = null;
		this.isMounted = false;
		this.stateHandler.reset();
	}
	updateSchema(schema) {
		const updated = [...schema];
		if (!updated.every((item) => Reflect.has(item, "fieldName") && item.fieldName)) {
			console.error("All items in the schema array must have a valid `fieldName` property to be updated");
			return;
		}
		const currentSchema = [...this.state?.schema ?? []];
		const updatedMap = {};
		updated.forEach((item) => {
			if (item.fieldName) updatedMap[item.fieldName] = item;
		});
		currentSchema.forEach((schema, index) => {
			const updatedData = updatedMap[schema.fieldName];
			if (updatedData) currentSchema[index] = mergeWithArrayOverride(updatedData, schema);
		});
		this.setState({ schema: currentSchema });
	}
	async validate(opts) {
		const validateResult = await (await this.getForm()).validate(opts);
		if (Object.keys(validateResult?.errors ?? {}).length > 0) {
			console.error("validate error", validateResult?.errors);
			if (this.state?.scrollToFirstError) this.scrollToFirstError(validateResult.errors);
		}
		return validateResult;
	}
	async validateAndSubmitForm() {
		const { valid, errors } = await (await this.getForm()).validate();
		if (!valid) {
			if (this.state?.scrollToFirstError) this.scrollToFirstError(errors);
			return;
		}
		return await this.submitForm();
	}
	async validateField(fieldName, opts) {
		const validateResult = await (await this.getForm()).validateField(fieldName, opts);
		if (Object.keys(validateResult?.errors ?? {}).length > 0) {
			console.error("validate error", validateResult?.errors);
			if (this.state?.scrollToFirstError) this.scrollToFirstError(fieldName);
		}
		return validateResult;
	}
	deleteValueByFieldName(values, fieldName) {
		const { pathSegments, rawKey } = resolveFieldNamePath(fieldName);
		if (rawKey) {
			Reflect.deleteProperty(values, rawKey);
			return;
		}
		if (!pathSegments || pathSegments.length === 0) {
			Reflect.deleteProperty(values, fieldName);
			return;
		}
		let target = values;
		for (const segment of pathSegments.slice(0, -1)) {
			if (!target || !isObject(target)) return;
			target = target[segment];
		}
		if (!target || !isObject(target)) return;
		const lastPathSegment = pathSegments.at(-1);
		if (!lastPathSegment) return;
		Reflect.deleteProperty(target, lastPathSegment);
	}
	async getForm() {
		if (!this.isMounted) await this.stateHandler.waitForCondition();
		if (!this.form?.meta) throw new Error("<VbenForm /> is not mounted");
		return this.form;
	}
	handleMultiFields = (originValues) => {
		const arrayToStringFields = this.state?.arrayToStringFields;
		if (!arrayToStringFields || !Array.isArray(arrayToStringFields)) return;
		const processFields = (fields, separator = ",") => {
			this.processFields(fields, separator, originValues, (value, sep) => {
				if (Array.isArray(value)) return value.join(sep);
				else if (typeof value === "string") {
					if (value === "") return [];
					const escapedSeparator = sep.replaceAll(/[.*+?^${}()|[\]\\]/g, String.raw`\$&`);
					return value.split(new RegExp(escapedSeparator));
				} else return value;
			});
		};
		if (arrayToStringFields.every((item) => typeof item === "string")) {
			const lastItem = arrayToStringFields[arrayToStringFields.length - 1] || "";
			processFields(lastItem.length === 1 ? arrayToStringFields.slice(0, -1) : arrayToStringFields, lastItem.length === 1 ? lastItem : ",");
			return;
		}
		arrayToStringFields.forEach((fieldConfig) => {
			if (Array.isArray(fieldConfig)) {
				const [fields, separator = ","] = fieldConfig;
				if (!Array.isArray(fields)) {
					console.warn(`Invalid field configuration: fields should be an array of strings, got ${typeof fields}`);
					return;
				}
				processFields(fields, separator);
			}
		});
	};
	handleRangeTimeValue = (originValues) => {
		const values = { ...originValues };
		const fieldMappingTime = this.state?.fieldMappingTime;
		this.handleMultiFields(values);
		if (!fieldMappingTime || !Array.isArray(fieldMappingTime)) return values;
		fieldMappingTime.forEach(([field, [startTimeKey, endTimeKey], format = "YYYY-MM-DD"]) => {
			if (startTimeKey && endTimeKey && values[field] === null) {
				Reflect.deleteProperty(values, startTimeKey);
				Reflect.deleteProperty(values, endTimeKey);
			}
			if (!values[field]) {
				Reflect.deleteProperty(values, field);
				return;
			}
			const [startTime, endTime] = values[field];
			if (format === null) {
				values[startTimeKey] = startTime;
				values[endTimeKey] = endTime;
			} else if (isFunction(format)) {
				values[startTimeKey] = format(startTime, startTimeKey);
				values[endTimeKey] = format(endTime, endTimeKey);
			} else {
				const [startTimeFormat, endTimeFormat] = Array.isArray(format) ? format : [format, format];
				values[startTimeKey] = startTime ? formatDate(startTime, startTimeFormat) : void 0;
				values[endTimeKey] = endTime ? formatDate(endTime, endTimeFormat) : void 0;
			}
			Reflect.deleteProperty(values, field);
		});
		return values;
	};
	handleValueFormat = (originValues) => {
		const values = { ...originValues };
		(this.state?.schema ?? []).forEach((schema) => {
			if (!schema.valueFormat) return;
			const fieldName = schema.fieldName;
			const value = this.resolveValueByFieldName(values, fieldName);
			this.deleteValueByFieldName(values, fieldName);
			const formattedValue = schema.valueFormat(value, (key, nextValue) => {
				this.setValueByFieldName(values, key, nextValue);
			}, values);
			if (formattedValue !== void 0) this.setValueByFieldName(values, fieldName, formattedValue);
		});
		return values;
	};
	processFields = (fields, separator, originValues, transformFn) => {
		fields.forEach((field) => {
			const value = originValues[field];
			if (value === void 0 || value === null) return;
			originValues[field] = transformFn(value, separator);
		});
	};
	resolveValueByFieldName(values, fieldName) {
		const { rawKey } = resolveFieldNamePath(fieldName);
		if (rawKey) return values[rawKey];
		return get(values, fieldName);
	}
	setValueByFieldName(values, fieldName, value) {
		const { rawKey } = resolveFieldNamePath(fieldName);
		if (rawKey) {
			values[rawKey] = value;
			return;
		}
		set(values, fieldName, value);
	}
	updateState() {
		const currentSchema = this.state?.schema ?? [];
		const prevSchema = this.prevState?.schema ?? [];
		if (currentSchema.length < prevSchema.length) {
			const currentFields = new Set(currentSchema.map((item) => item.fieldName));
			const deletedSchema = prevSchema.filter((item) => !currentFields.has(item.fieldName));
			for (const schema of deletedSchema) this.form?.setFieldValue?.(schema.fieldName, void 0);
		}
	}
};
//#endregion
export { FormApi };
