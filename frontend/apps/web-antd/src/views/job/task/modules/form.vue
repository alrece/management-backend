<script lang="ts" setup>
import type { JobTaskApi } from '#/api/job/task';

import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { createJobTask, updateJobTask } from '#/api/job/task';

import { useFormSchema } from '../data';

const emit = defineEmits(['success']);
const formData = ref<JobTaskApi.Task>();
const getTitle = computed(() => {
  return formData.value?.id ? '编辑任务' : '新建任务';
});

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
    const data = (await formApi.getValues()) as JobTaskApi.Task;
    try {
      await (formData.value?.id
        ? updateJobTask(formData.value.id, data)
        : createJobTask(data));
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
    const data = modalApi.getData<JobTaskApi.Task>();
    if (data?.id) {
      formData.value = data;
      await formApi.setValues(data);
    } else {
      await formApi.resetForm();
      await formApi.setValues({ status: 0 });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]" :title="getTitle">
    <Form class="mx-4" />
  </Modal>
</template>
