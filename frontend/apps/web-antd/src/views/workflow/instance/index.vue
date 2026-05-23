<script lang="ts" setup>
import type { VxeTableGridOptions } from '#/adapter/vxe-table';
import type { InstanceApi } from '#/api/workflow/instance';

import { ref, h } from 'vue';
import { Page } from '@vben/common-ui';
import { message, Tag, Modal } from 'ant-design-vue';
import { useVbenVxeGrid, ACTION_ICON, TableAction } from '#/adapter/vxe-table';
import { getInstancePage, getInstanceLogs } from '#/api/workflow/instance';

const logModalVisible = ref(false);
const currentLogs = ref<InstanceApi.NodeLog[]>([]);

const STATUS_MAP: Record<string, { color: string; text: string }> = {
  running: { color: 'processing', text: '运行中' },
  success: { color: 'green', text: '成功' },
  failed: { color: 'red', text: '失败' },
  canceled: { color: 'default', text: '已取消' },
};

function renderStatus(status: string) {
  const item = STATUS_MAP[status] ?? { color: 'default', text: '未知' };
  return h(Tag, { color: item.color }, () => item.text);
}

async function handleViewLogs(row: InstanceApi.Instance) {
  const hide = message.loading({ content: '加载日志...', duration: 0 });
  try {
    currentLogs.value = await getInstanceLogs(row.id);
    logModalVisible.value = true;
  } catch {
    message.error('加载日志失败');
  } finally {
    hide();
  }
}

const [Grid] = useVbenVxeGrid({
  formOptions: {
    schema: [
      {
        fieldName: 'status',
        label: '状态',
        component: 'Select',
        componentProps: {
          options: [
            { label: '运行中', value: 'running' },
            { label: '成功', value: 'success' },
            { label: '失败', value: 'failed' },
          ],
          placeholder: '请选择状态',
          allowClear: true,
        },
      },
    ],
  },
  gridOptions: {
    columns: [
      { field: 'id', title: '执行ID', minWidth: 80 },
      { field: 'workflowId', title: '工作流ID', minWidth: 80 },
      { field: 'status', title: '状态', minWidth: 90, slots: { default: 'status' } },
      { field: 'triggerType', title: '触发类型', minWidth: 90 },
      { field: 'startTime', title: '开始时间', minWidth: 170 },
      { field: 'endTime', title: '结束时间', minWidth: 170 },
      { field: 'errorMsg', title: '错误信息', minWidth: 200 },
      { title: '操作', width: 120, fixed: 'right', slots: { default: 'actions' } },
    ],
    height: 'auto',
    proxyConfig: {
      ajax: {
        query: async ({ page }, formValues) => {
          return await getInstancePage({
            pageNo: page.currentPage,
            pageSize: page.pageSize,
            ...formValues,
          });
        },
      },
    },
    rowConfig: { keyField: 'id', isHover: true },
  } as VxeTableGridOptions<InstanceApi.Instance>,
});
</script>

<template>
  <Page auto-content-height>
    <Grid table-title="执行实例列表">
      <template #status="{ row }">
        <component :is="renderStatus(row.status)" />
      </template>
      <template #actions="{ row }">
        <TableAction
          :actions="[
            {
              label: '查看日志',
              type: 'link',
              icon: ACTION_ICON.VIEW,
              onClick: () => handleViewLogs(row),
            },
          ]"
        />
      </template>
    </Grid>
    <Modal
      v-model:open="logModalVisible"
      title="节点执行日志"
      width="800px"
      :footer="null"
    >
      <div v-if="currentLogs.length === 0" style="text-align: center; color: #999; padding: 20px">
        暂无日志
      </div>
      <div v-else>
        <div
          v-for="log in currentLogs"
          :key="log.id"
          style="padding: 8px 0; border-bottom: 1px solid #f0f0f0"
        >
          <div style="display: flex; justify-content: space-between; margin-bottom: 4px">
            <span style="font-weight: 500">{{ log.nodeId }} ({{ log.nodeType }})</span>
            <component :is="renderStatus(log.status)" />
          </div>
          <div v-if="log.errorMsg" style="color: #f56c6c; font-size: 12px">{{ log.errorMsg }}</div>
          <div style="color: #999; font-size: 12px">
            {{ log.startTime }} → {{ log.endTime || '...' }}
            <span v-if="log.retryCount > 0"> | 重试 {{ log.retryCount }} 次</span>
          </div>
        </div>
      </div>
    </Modal>
  </Page>
</template>
