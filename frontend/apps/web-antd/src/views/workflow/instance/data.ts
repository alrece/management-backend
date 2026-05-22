import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Tag } from 'ant-design-vue';

/** 实例状态枚举 */
export const INSTANCE_STATUS_OPTIONS = [
  { label: '运行中', value: 'RUNNING' },
  { label: '成功', value: 'SUCCESS' },
  { label: '失败', value: 'FAILED' },
  { label: '已取消', value: 'CANCELLED' },
] as const;

const INSTANCE_STATUS_MAP: Record<string, { color: string; text: string }> = {
  RUNNING: { color: 'processing', text: '运行中' },
  SUCCESS: { color: 'success', text: '成功' },
  FAILED: { color: 'error', text: '失败' },
  CANCELLED: { color: 'default', text: '已取消' },
};

export function renderInstanceStatusTag(status: string) {
  const item = INSTANCE_STATUS_MAP[status] ?? { color: 'default', text: '未知' };
  return h(Tag, { color: item.color }, () => item.text);
}

/** 列表的搜索表单 */
export function useGridFormSchema(): VbenFormSchema[] {
  return [
    {
      fieldName: 'workflowId',
      label: '工作流ID',
      component: 'Input',
      componentProps: { placeholder: '请输入工作流ID', allowClear: true },
    },
    {
      fieldName: 'status',
      label: '状态',
      component: 'Select',
      componentProps: {
        options: INSTANCE_STATUS_OPTIONS.map((o) => ({ label: o.label, value: o.value })),
        placeholder: '请选择状态',
        allowClear: true,
      },
    },
  ];
}

/** 列表的字段 */
export function useGridColumns(): VxeTableGridOptions['columns'] {
  return [
    { field: 'id', title: '实例编号', minWidth: 100 },
    { field: 'workflowId', title: '工作流ID', minWidth: 90 },
    {
      field: 'status',
      title: '状态',
      minWidth: 80,
      slots: { default: 'status' },
    },
    { field: 'durationMs', title: '耗时(ms)', minWidth: 90 },
    { field: 'variables', title: '入参', minWidth: 160, showOverflow: true },
    { field: 'result', title: '执行结果', minWidth: 200, showOverflow: true },
    { field: 'errorMsg', title: '错误信息', minWidth: 200, showOverflow: true },
    {
      field: 'startedAt',
      title: '开始时间',
      minWidth: 170,
      formatter: 'formatDateTime',
    },
    {
      field: 'finishedAt',
      title: '结束时间',
      minWidth: 170,
      formatter: 'formatDateTime',
    },
  ];
}
