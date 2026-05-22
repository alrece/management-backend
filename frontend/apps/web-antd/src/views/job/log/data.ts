import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Tag } from 'ant-design-vue';

/** 执行状态枚举 */
export const EXEC_STATUS_OPTIONS = [
  { label: '执行中', value: 0 },
  { label: '成功', value: 1 },
  { label: '失败', value: 2 },
  { label: '超时', value: 3 },
] as const;

/** 触发类型枚举 */
export const TRIGGER_TYPE_OPTIONS = [
  { label: '定时', value: 1 },
  { label: '手动', value: 2 },
] as const;

const EXEC_STATUS_MAP: Record<number, { color: string; text: string }> = {
  0: { color: 'processing', text: '执行中' },
  1: { color: 'success', text: '成功' },
  2: { color: 'error', text: '失败' },
  3: { color: 'warning', text: '超时' },
};

export function renderExecStatusTag(status: number) {
  const item = EXEC_STATUS_MAP[status] ?? { color: 'default', text: '未知' };
  return h(Tag, { color: item.color }, () => item.text);
}

const TRIGGER_TYPE_MAP: Record<number, string> = { 1: '定时', 2: '手动' };

export function renderTriggerType(type: number) {
  return TRIGGER_TYPE_MAP[type] ?? '未知';
}

/** 列表的搜索表单 */
export function useGridFormSchema(): VbenFormSchema[] {
  return [
    {
      fieldName: 'taskId',
      label: '任务ID',
      component: 'Input',
      componentProps: { placeholder: '请输入任务ID', allowClear: true },
    },
    {
      fieldName: 'status',
      label: '执行状态',
      component: 'Select',
      componentProps: {
        options: EXEC_STATUS_OPTIONS.map((o) => ({ label: o.label, value: o.value })),
        placeholder: '请选择状态',
        allowClear: true,
      },
    },
  ];
}

/** 列表的字段 */
export function useGridColumns(): VxeTableGridOptions['columns'] {
  return [
    { field: 'id', title: '日志编号', minWidth: 100 },
    { field: 'taskId', title: '任务ID', minWidth: 80 },
    { field: 'taskName', title: '任务名称', minWidth: 150 },
    {
      field: 'triggerType',
      title: '触发类型',
      minWidth: 80,
      slots: { default: 'triggerType' },
    },
    {
      field: 'status',
      title: '执行状态',
      minWidth: 80,
      slots: { default: 'status' },
    },
    { field: 'durationMs', title: '耗时(ms)', minWidth: 90 },
    { field: 'result', title: '执行结果', minWidth: 200, showOverflow: true },
    { field: 'error', title: '错误信息', minWidth: 200, showOverflow: true },
    {
      field: 'startTime',
      title: '开始时间',
      minWidth: 170,
      formatter: 'formatDateTime',
    },
  ];
}
