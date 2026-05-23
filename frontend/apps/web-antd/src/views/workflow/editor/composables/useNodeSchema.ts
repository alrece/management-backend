import { ref } from 'vue';

import type { WorkflowApi } from '#/api/workflow/workflow';
import { getNodeSchemas } from '#/api/workflow/workflow';

export function useNodeSchema() {
  const schemas = ref<WorkflowApi.NodeSchema[]>([]);
  const loading = ref(false);

  async function fetchSchemas() {
    loading.value = true;
    try {
      schemas.value = await getNodeSchemas();
    } catch {
      schemas.value = [];
    } finally {
      loading.value = false;
    }
  }

  return { schemas, loading, fetchSchemas };
}
