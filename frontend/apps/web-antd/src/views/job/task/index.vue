<script lang="ts" setup>
import type { VxeTableGridOptions } from '#/adapter/vxe-table';
import type { JobTaskApi } from '#/api/job/task';

import { ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message } from 'ant-design-vue';

import { ACTION_ICON, TableAction, useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteJobTask,
  getJobTaskPage,
  triggerJobTask,
  updateJobTask,
} from '#/api/job/task';

import { useGridColumns, useGridFormSchema } from './data';
import { renderStatusTag } from './data';
import Form from './modules/form.vue';

const [FormModal, formModalApi] = useVbenModal({
  connectedComponent: Form,
  destroyOnClose: true,
});

const checkedIds = ref<number[]>([]);

function handleRefresh() {
  gridApi.query();
}

function handleCreate() {
  formModalApi.setData(null).open();
}

function handleEdit(row: JobTaskApi.Task) {
  formModalApi.setData(row).open();
}

async function handleDelete(row: JobTaskApi.Task) {
  const hideLoading = message.loading({ content: '删除中...', duration: 0 });
  try {
    await deleteJobTask(row.id!);
    message.success('删除成功');
    handleRefresh();
  } finally {
    hideLoading();
  }
}

async function handleToggleStatus(row: JobTaskApi.Task) {
  const newStatus = row.status === 0 ? 1 : 0;
  const label = newStatus === 0 ? '恢复' : '暂停';
  const hideLoading = message.loading({ content: `${label}中...`, duration: 0 });
  try {
    await updateJobTask(row.id!, { status: newStatus });
    message.success(`${label}成功`);
    handleRefresh();
  } finally {
    hideLoading();
  }
}

async function handleTrigger(row: JobTaskApi.Task) {
  const hideLoading = message.loading({ content: '触发中...', duration: 0 });
  try {
    await triggerJobTask(row.id!);
    message.success('触发成功');
  } finally {
    hideLoading();
  }
}

function handleRowCheckboxChange({ records }: { records: JobTaskApi.Task[] }) {
  checkedIds.value = records.map((item) => item.id!);
}

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: useGridFormSchema(),
  },
  gridOptions: {
    columns: useGridColumns(),
    height: 'auto',
    keepSource: true,
    proxyConfig: {
      ajax: {
        query: async ({ page }, formValues) => {
          return await getJobTaskPage({
            pageNo: page.currentPage,
            pageSize: page.pageSize,
            ...formValues,
          });
        },
      },
    },
    rowConfig: { keyField: 'id', isHover: true },
    toolbarConfig: { refresh: true, search: true },
  } as VxeTableGridOptions<JobTaskApi.Task>,
  gridEvents: {
    checkboxAll: handleRowCheckboxChange,
    checkboxChange: handleRowCheckboxChange,
  },
});
</script>

<template>
  <Page auto-content-height>
    <FormModal @success="handleRefresh" />
    <Grid table-title="定时任务列表">
      <template #status="{ row }">
        <component :is="renderStatusTag(row.status)" />
      </template>
      <template #toolbar-tools>
        <TableAction
          :actions="[
            {
              label: '新建任务',
              type: 'primary',
              icon: ACTION_ICON.ADD,
              auth: ['job:task:create'],
              onClick: handleCreate,
            },
          ]"
        />
      </template>
      <template #actions="{ row }">
        <TableAction
          :actions="[
            {
              label: '编辑',
              type: 'link',
              icon: ACTION_ICON.EDIT,
              auth: ['job:task:update'],
              onClick: handleEdit.bind(null, row),
            },
            {
              label: row.status === 0 ? '暂停' : '恢复',
              type: 'link',
              icon: ACTION_ICON.EDIT,
              auth: ['job:task:update'],
              popConfirm: {
                title: `确认${row.status === 0 ? '暂停' : '恢复'}吗？`,
                confirm: handleToggleStatus.bind(null, row),
              },
            },
            {
              label: '执行',
              type: 'link',
              icon: ACTION_ICON.EDIT,
              auth: ['job:task:trigger'],
              popConfirm: {
                title: '确认立即执行吗？',
                confirm: handleTrigger.bind(null, row),
              },
            },
            {
              label: '删除',
              type: 'link',
              danger: true,
              icon: ACTION_ICON.DELETE,
              auth: ['job:task:delete'],
              popConfirm: {
                title: `确认删除「${row.name}」吗？`,
                confirm: handleDelete.bind(null, row),
              },
            },
          ]"
        />
      </template>
    </Grid>
  </Page>
</template>
