<script lang="ts" setup>
import type { VxeTableGridOptions } from '#/adapter/vxe-table';
import type { WorkflowApi } from '#/api/workflow/workflow';

import { ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message } from 'ant-design-vue';

import { ACTION_ICON, TableAction, useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  activateWorkflow,
  deactivateWorkflow,
  deleteWorkflow,
  executeWorkflow,
  getWorkflowPage,
} from '#/api/workflow/workflow';

import { useGridColumns, useGridFormSchema, renderWorkflowStatusTag } from './data';
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

function handleEdit(row: WorkflowApi.Workflow) {
  formModalApi.setData(row).open();
}

async function handleDelete(row: WorkflowApi.Workflow) {
  const hideLoading = message.loading({ content: '删除中...', duration: 0 });
  try {
    await deleteWorkflow(row.id!);
    message.success('删除成功');
    handleRefresh();
  } finally {
    hideLoading();
  }
}

async function handleActivate(row: WorkflowApi.Workflow) {
  const hideLoading = message.loading({ content: '激活中...', duration: 0 });
  try {
    await activateWorkflow(row.id!);
    message.success('激活成功');
    handleRefresh();
  } finally {
    hideLoading();
  }
}

async function handleDeactivate(row: WorkflowApi.Workflow) {
  const hideLoading = message.loading({ content: '停用中...', duration: 0 });
  try {
    await deactivateWorkflow(row.id!);
    message.success('停用成功');
    handleRefresh();
  } finally {
    hideLoading();
  }
}

async function handleExecute(row: WorkflowApi.Workflow) {
  const hideLoading = message.loading({ content: '执行中...', duration: 0 });
  try {
    await executeWorkflow(row.id!);
    message.success('执行成功');
  } finally {
    hideLoading();
  }
}

function handleRowCheckboxChange({ records }: { records: WorkflowApi.Workflow[] }) {
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
          return await getWorkflowPage({
            pageNo: page.currentPage,
            pageSize: page.pageSize,
            ...formValues,
          });
        },
      },
    },
    rowConfig: { keyField: 'id', isHover: true },
    toolbarConfig: { refresh: true, search: true },
  } as VxeTableGridOptions<WorkflowApi.Workflow>,
  gridEvents: {
    checkboxAll: handleRowCheckboxChange,
    checkboxChange: handleRowCheckboxChange,
  },
});
</script>

<template>
  <Page auto-content-height>
    <FormModal @success="handleRefresh" />
    <Grid table-title="工作流列表">
      <template #status="{ row }">
        <component :is="renderWorkflowStatusTag(row.status)" />
      </template>
      <template #toolbar-tools>
        <TableAction
          :actions="[
            {
              label: '新建流程',
              type: 'primary',
              icon: ACTION_ICON.ADD,
              auth: ['workflow:workflow:create'],
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
              auth: ['workflow:workflow:update'],
              onClick: handleEdit.bind(null, row),
            },
            ...(row.status === 'ACTIVE'
              ? [
                  {
                    label: '停用',
                    type: 'link' as const,
                    icon: ACTION_ICON.EDIT,
                    auth: ['workflow:workflow:activate'],
                    popConfirm: {
                      title: '确认停用吗？',
                      confirm: handleDeactivate.bind(null, row),
                    },
                  },
                ]
              : []),
            ...(row.status !== 'ACTIVE'
              ? [
                  {
                    label: '激活',
                    type: 'link' as const,
                    icon: ACTION_ICON.EDIT,
                    auth: ['workflow:workflow:activate'],
                    popConfirm: {
                      title: '确认激活吗？',
                      confirm: handleActivate.bind(null, row),
                    },
                  },
                ]
              : []),
            {
              label: '执行',
              type: 'link',
              icon: ACTION_ICON.EDIT,
              auth: ['workflow:workflow:execute'],
              popConfirm: {
                title: '确认执行吗？',
                confirm: handleExecute.bind(null, row),
              },
            },
            {
              label: '删除',
              type: 'link',
              danger: true,
              icon: ACTION_ICON.DELETE,
              auth: ['workflow:workflow:delete'],
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
