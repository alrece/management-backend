import type { VbenFormSchema } from '#/adapter/form';

import { z } from '#/adapter/form';

/** 分类表单 */
export function useFormSchema(treeData: any[] = []): VbenFormSchema[] {
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
      label: '分类名称',
      componentProps: { placeholder: '请输入分类名称' },
      rules: 'required',
    },
    {
      component: 'TreeSelect',
      fieldName: 'parentId',
      label: '上级分类',
      componentProps: {
        treeData,
        placeholder: '请选择上级分类（留空为顶级）',
        allowClear: true,
        treeDefaultExpandAll: true,
        fieldNames: { label: 'name', value: 'id', children: 'children' },
      },
    },
    {
      component: 'InputNumber',
      fieldName: 'sort',
      label: '排序',
      componentProps: { min: 0, placeholder: '请输入排序' },
      rules: z.number().default(0),
    },
  ];
}

/** 树列定义 */
export const treeColumns = [
  { title: '分类名称', dataIndex: 'name', key: 'name' },
  { title: '排序', dataIndex: 'sort', key: 'sort', width: 80 },
  { title: '操作', key: 'action', width: 200 },
];
