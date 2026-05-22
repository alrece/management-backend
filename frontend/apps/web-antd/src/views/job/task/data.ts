import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Tag } from 'ant-design-vue';

import { z } from '#/adapter/form';

/** 任务状态枚举 */
export const TASK_STATUS_OPTIONS = [
  { label: '正常', value: 0 },
  { label: '暂停', value: 1 },
] as const;

/** 状态渲染 Tag */
export function renderStatusTag(status: number) {
  const map: Record<number, { color: string; text: string }> = {
    0: { color: 'green', text: '正常' },
    1: { color: 'orange', text: '暂停' },
  };
  const item = map[status] ?? { color: 'default', text: '未知' };
  return h(Tag, { color: item.color }, () => item.text);
}

/** 新增/修改的表单 */
export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'id',
      dependencies: {
        triggerFields: [''],
        show: () => false,
      },
    },
    {
      component: 'Input',
      fieldName: 'name',
      label: '任务名称',
      componentProps: { placeholder: '请输入任务名称' },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'handler',
      label: '处理器',
      componentProps: { placeholder: '请输入处理器名称' },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'cronExpr',
      label: 'Cron 表达式',
      componentProps: { placeholder: '如: 0 0/5 * * * ?（每5分钟）' },
      rules: 'required',
    },
    {
      component: 'Textarea',
      fieldName: 'params',
      label: '参数',
      componentProps: { placeholder: 'JSON 格式参数，可选' },
    },
    {
      component: 'RadioGroup',
      fieldName: 'status',
      label: '状态',
      componentProps: {
        options: TASK_STATUS_OPTIONS.map((o) => ({ label: o.label, value: o.value })),
        buttonStyle: 'solid',
        optionType: 'button',
      },
      rules: z.number().default(0),
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注' },
    },
  ];
}

/** 列表的搜索表单 */
export function useGridFormSchema(): VbenFormSchema[] {
  return [
    {
      fieldName: 'name',
      label: '任务名称',
      component: 'Input',
      componentProps: { placeholder: '请输入任务名称', allowClear: true },
    },
    {
      fieldName: 'status',
      label: '状态',
      component: 'Select',
      componentProps: {
        options: TASK_STATUS_OPTIONS.map((o) => ({ label: o.label, value: o.value })),
        placeholder: '请选择状态',
        allowClear: true,
      },
    },
  ];
}

/** 列表的字段 */
export function useGridColumns(): VxeTableGridOptions['columns'] {
  return [
    { type: 'checkbox', width: 40 },
    { field: 'id', title: '任务编号', minWidth: 100 },
    { field: 'name', title: '任务名称', minWidth: 150 },
    { field: 'handler', title: '处理器', minWidth: 150 },
    { field: 'cronExpr', title: 'Cron 表达式', minWidth: 140 },
    {
      field: 'status',
      title: '状态',
      minWidth: 80,
      slots: { default: 'status' },
    },
    { field: 'remark', title: '备注', minWidth: 150 },
    { field: 'creator', title: '创建者', minWidth: 80 },
    {
      title: '操作',
      width: 200,
      fixed: 'right',
      slots: { default: 'actions' },
    },
  ];
}
