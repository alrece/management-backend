import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Tag } from 'ant-design-vue';

import { z } from '#/adapter/form';

import { getCategoryTree } from '#/api/workflow/category';

/** 工作流状态枚举 */
export const WORKFLOW_STATUS_OPTIONS = [
  { label: '草稿', value: 'DRAFT' },
  { label: '已激活', value: 'ACTIVE' },
  { label: '已停用', value: 'INACTIVE' },
  { label: '已失同步', value: 'DESYNCED' },
] as const;

const WORKFLOW_STATUS_MAP: Record<string, { color: string; text: string }> = {
  DRAFT: { color: 'default', text: '草稿' },
  ACTIVE: { color: 'green', text: '已激活' },
  INACTIVE: { color: 'orange', text: '已停用' },
  DESYNCED: { color: 'red', text: '已失同步' },
};

export function renderWorkflowStatusTag(status: string) {
  const item = WORKFLOW_STATUS_MAP[status] ?? { color: 'default', text: '未知' };
  return h(Tag, { color: item.color }, () => item.text);
}

/** 新增/修改的表单 */
export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'id',
      dependencies: { triggerFields: [''], show: () => false },
    },
    {
      component: 'Input',
      fieldName: 'name',
      label: '流程名称',
      componentProps: { placeholder: '请输入流程名称' },
      rules: 'required',
    },
    {
      component: 'ApiTreeSelect',
      fieldName: 'categoryId',
      label: '所属分类',
      componentProps: {
        api: getCategoryTree,
        placeholder: '请选择分类',
        allowClear: true,
        treeDefaultExpandAll: true,
        fieldNames: { label: 'name', value: 'id', children: 'children' },
      },
    },
    {
      component: 'Textarea',
      fieldName: 'paramsSchema',
      label: '参数 Schema',
      componentProps: { placeholder: 'JSON Schema 格式，定义流程参数结构' },
    },
  ];
}

/** 列表的搜索表单 */
export function useGridFormSchema(): VbenFormSchema[] {
  return [
    {
      fieldName: 'name',
      label: '流程名称',
      component: 'Input',
      componentProps: { placeholder: '请输入流程名称', allowClear: true },
    },
    {
      fieldName: 'status',
      label: '状态',
      component: 'Select',
      componentProps: {
        options: WORKFLOW_STATUS_OPTIONS.map((o) => ({ label: o.label, value: o.value })),
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
    { field: 'id', title: '编号', minWidth: 80 },
    { field: 'name', title: '流程名称', minWidth: 180 },
    { field: 'categoryId', title: '分类', minWidth: 100 },
    {
      field: 'status',
      title: '状态',
      minWidth: 90,
      slots: { default: 'status' },
    },
    {
      field: 'createdAt',
      title: '创建时间',
      minWidth: 170,
      formatter: 'formatDateTime',
    },
    {
      title: '操作',
      width: 260,
      fixed: 'right',
      slots: { default: 'actions' },
    },
  ];
}
