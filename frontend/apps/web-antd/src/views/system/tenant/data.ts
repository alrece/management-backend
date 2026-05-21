import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

import { z } from '#/adapter/form';

/** 新增/修改的表单 */
export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      fieldName: 'id',
      component: 'Input',
      dependencies: {
        triggerFields: [''],
        show: () => false,
      },
    },
    {
      fieldName: 'tenantName',
      label: '租户名称',
      component: 'Input',
      componentProps: { placeholder: '请输入租户名称' },
      rules: 'required',
    },
    {
      fieldName: 'contactName',
      label: '联系人',
      component: 'Input',
      componentProps: { placeholder: '请输入联系人' },
    },
    {
      fieldName: 'contactMobile',
      label: '联系手机',
      component: 'Input',
      componentProps: { placeholder: '请输入联系手机' },
    },
    {
      fieldName: 'status',
      label: '状态',
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '开启', value: 0 },
          { label: '关闭', value: 1 },
        ],
        buttonStyle: 'solid',
        optionType: 'button',
      },
      rules: z.number().default(0),
    },
  ];
}

/** 列表的搜索表单 */
export function useGridFormSchema(): VbenFormSchema[] {
  return [
    {
      fieldName: 'tenantName',
      label: '租户名',
      component: 'Input',
      componentProps: { placeholder: '请输入租户名', allowClear: true },
    },
  ];
}

/** 列表的字段 */
export function useGridColumns(): VxeTableGridOptions['columns'] {
  return [
    { type: 'checkbox', width: 40 },
    { field: 'id', title: '租户编号', minWidth: 100 },
    { field: 'tenantName', title: '租户名', minWidth: 180 },
    { field: 'status', title: '状态', minWidth: 100 },
    { field: 'createTime', title: '创建时间', minWidth: 180, formatter: 'formatDateTime' },
    { title: '操作', width: 130, fixed: 'right', slots: { default: 'actions' } },
  ];
}
