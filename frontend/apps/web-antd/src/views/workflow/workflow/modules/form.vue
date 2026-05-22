<script lang="ts" setup>
import type { WorkflowApi } from '#/api/workflow/workflow';

import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { createWorkflow, updateWorkflow } from '#/api/workflow/workflow';

import { useFormSchema } from '../data';

const emit = defineEmits(['success']);
const formData = ref<WorkflowApi.Workflow>();
const getTitle = computed(() => (formData.value?.id ? '编辑流程' : '新建流程'));

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
    const data = (await formApi.getValues()) as WorkflowApi.WorkflowForm;
    try {
      await (formData.value?.id
        ? updateWorkflow(formData.value.id, data)
        : createWorkflow(data));
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
    const data = modalApi.getData<WorkflowApi.Workflow>();
    if (data?.id) {
      formData.value = data;
      await formApi.setValues(data);
    } else {
      await formApi.resetForm();
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]" :title="getTitle">
    <Form class="mx-4" />
  </Modal>
</template>
