<script lang="ts" setup>
import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { createCategory, updateCategory } from '#/api/workflow/category';

import { useFormSchema } from '../data';

interface FormData {
  id?: number;
  name: string;
  parentId: number;
  sort: number;
  treeData?: any[];
}

const emit = defineEmits(['success']);
const formData = ref<FormData>();
const getTitle = computed(() => (formData.value?.id ? '编辑分类' : '新建分类'));

const [Form, formApi] = useVbenForm({
  commonConfig: {
    componentProps: { class: 'w-full' },
    formItemClass: 'col-span-2',
    labelWidth: 100,
  },
  layout: 'horizontal',
  schema: useFormSchema(),
  showDefaultActions: false,
});

const [Modal, modalApi] = useVbenModal({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) return;
    modalApi.lock();
    const data = (await formApi.getValues()) as FormData;
    try {
      await (formData.value?.id
        ? updateCategory(formData.value.id, { name: data.name, parentId: data.parentId, sort: data.sort })
        : createCategory({ name: data.name, parentId: data.parentId, sort: data.sort }));
      await modalApi.close();
      emit('success');
      message.success('操作成功');
    } finally {
      modalApi.unlock();
    }
  },
  async onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      formData.value = undefined;
      return;
    }
    const data = modalApi.getData<FormData>();
    // 更新 treeData 到表单 schema
    const schema = useFormSchema(data?.treeData ?? []);
    await formApi.updateSchema(schema);

    if (data?.id) {
      formData.value = data;
      await formApi.setValues(data);
    } else {
      await formApi.resetForm();
      await formApi.setValues({ parentId: data?.parentId ?? 0, sort: 0 });
    }
  },
});
</script>

<template>
  <Modal class="w-[500px]" :title="getTitle">
    <Form class="mx-4" />
  </Modal>
</template>
