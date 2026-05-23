import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

import { h } from 'vue';
import { Tag } from 'ant-design-vue';
import { useRouter } from 'vue-router';

import { getCategoryTree } from '#/api/workflow/category';
import { WORKFLOW_STATUS_OPTIONS } from '#/views/workflow/editor/types';

const WORKFLOW_STATUS_MAP: Record<number, { color: string; text: string }> = {
  0: { color: 'green', text: '正常' },
  1: { color: 'orange', text: '停用' },
};

export function renderWorkflowStatusTag(status: number) {
  const item = WORKFLOW_STATUS_MAP[status] ?? { color: 'default', text: '未知' };
  return h(Tag, { color: item.color }, () => item.text);
}

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
      component: 'Select',
      fieldName: 'triggerType',
      label: '触发类型',
      componentProps: {
        options: [
          { label: '手动', value: 'manual' },
          { label: '定时', value: 'cron' },
          { label: 'Webhook', value: 'webhook' },
        ],
        placeholder: '请选择触发类型',
      },
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注' },
    },
  ];
}

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

export function useGridColumns(): VxeTableGridOptions['columns'] {
  return [
    { type: 'checkbox', width: 40 },
    { field: 'id', title: '编号', minWidth: 80 },
    { field: 'name', title: '流程名称', minWidth: 180 },
    { field: 'triggerType', title: '触发类型', minWidth: 100 },
    { field: 'version', title: '版本', minWidth: 60 },
    { field: 'status', title: '状态', minWidth: 90, slots: { default: 'status' } },
    { field: 'createTime', title: '创建时间', minWidth: 170, formatter: 'formatDateTime' },
    { title: '操作', width: 300, fixed: 'right', slots: { default: 'actions' } },
  ];
}
