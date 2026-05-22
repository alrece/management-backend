<script lang="ts" setup>
import type { VxeTableGridOptions } from '#/adapter/vxe-table';
import type { InstanceApi } from '#/api/workflow/instance';

import { Page } from '@vben/common-ui';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getInstancePage } from '#/api/workflow/instance';

import { useGridColumns, useGridFormSchema, renderInstanceStatusTag } from './data';

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
          return await getInstancePage({
            pageNo: page.currentPage,
            pageSize: page.pageSize,
            ...formValues,
          });
        },
      },
    },
    rowConfig: { keyField: 'id', isHover: true },
    toolbarConfig: { refresh: true, search: true },
  } as VxeTableGridOptions<InstanceApi.Instance>,
});
</script>

<template>
  <Page auto-content-height>
    <Grid table-title="执行实例列表">
      <template #status="{ row }">
        <component :is="renderInstanceStatusTag(row.status)" />
      </template>
    </Grid>
  </Page>
</template>
