<script lang="ts" setup>
import type { VxeTableGridOptions } from '#/adapter/vxe-table';
import type { ExecLogApi } from '#/api/job/log';

import { Page } from '@vben/common-ui';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getExecLogPage } from '#/api/job/log';

import { useGridColumns, useGridFormSchema, renderExecStatusTag, renderTriggerType } from './data';

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
          return await getExecLogPage({
            pageNo: page.currentPage,
            pageSize: page.pageSize,
            ...formValues,
          });
        },
      },
    },
    rowConfig: { keyField: 'id', isHover: true },
    toolbarConfig: { refresh: true, search: true },
  } as VxeTableGridOptions<ExecLogApi.Log>,
});
</script>

<template>
  <Page auto-content-height>
    <Grid table-title="执行日志列表">
      <template #status="{ row }">
        <component :is="renderExecStatusTag(row.status)" />
      </template>
      <template #triggerType="{ row }">
        {{ renderTriggerType(row.triggerType) }}
      </template>
    </Grid>
  </Page>
</template>
