import { useVueFlow } from '@vue-flow/core';
import { ref } from 'vue';
import { message } from 'ant-design-vue';

import type { WorkflowDefinition } from '../types';

import { getWorkflow, updateWorkflow, createWorkflow, executeWorkflow } from '#/api/workflow/workflow';

export function useWorkflow(workflowId: ref<number | null>) {
  const { toObject, fromJSON } = useVueFlow();
  const saving = ref(false);
  const loading = ref(false);
  const currentVersion = ref(1);

  async function load() {
    if (!workflowId.value) return;
    loading.value = true;
    try {
      const wf = await getWorkflow(workflowId.value);
      if (wf.definition) {
        const def: WorkflowDefinition = typeof wf.definition === 'string'
          ? JSON.parse(wf.definition)
          : wf.definition;
        fromJSON({
          nodes: def.nodes.map((n) => ({
            id: n.id,
            type: getNodeType(n.type),
            position: n.position,
            data: n.data,
            label: n.data?.label || n.type,
          })),
          edges: def.edges.map((e) => ({
            id: e.id,
            source: e.source,
            target: e.target,
            sourceHandle: e.sourceHandle,
            targetHandle: e.targetHandle,
            animated: e.animated,
            label: e.label,
          })),
        } as any);
      }
      currentVersion.value = wf.version || 1;
    } catch {
      message.error('加载工作流失败');
    } finally {
      loading.value = false;
    }
  }

  async function save(name: string, categoryId?: number) {
    saving.value = true;
    try {
      const flow = toObject();
      const definition: WorkflowDefinition = {
        nodes: flow.nodes.map((n) => ({
          id: n.id,
          type: n.data?.type || n.type,
          position: n.position,
          data: n.data || {},
        })),
        edges: flow.edges.map((e) => ({
          id: e.id,
          source: e.source,
          target: e.target,
          sourceHandle: e.sourceHandle || undefined,
          targetHandle: e.targetHandle || undefined,
          animated: e.animated || false,
          label: e.label || undefined,
        })),
        settings: {},
      };

      const payload = {
        name,
        categoryId,
        definition: JSON.stringify(definition),
        triggerType: 'manual',
        version: currentVersion.value + 1,
      };

      if (workflowId.value) {
        await updateWorkflow({ id: workflowId.value, ...payload });
        currentVersion.value++;
        message.success('保存成功');
      } else {
        const id = await createWorkflow(payload);
        workflowId.value = id;
        message.success('创建成功');
      }
    } catch (e: any) {
      message.error(e?.message || '保存失败');
    } finally {
      saving.value = false;
    }
  }

  async function execute(input?: any) {
    if (!workflowId.value) {
      message.warning('请先保存工作流');
      return;
    }
    try {
      await executeWorkflow(workflowId.value, input);
      message.success('执行成功');
    } catch (e: any) {
      message.error(e?.message || '执行失败');
    }
  }

  return { saving, loading, currentVersion, load, save, execute };
}

function getNodeType(type: string): string {
  if (type === 'trigger' || type === 'cron' || type === 'webhook') return 'triggerNode';
  if (type === 'condition' || type === 'delay') return 'controlNode';
  if (type === 'transform') return 'transformNode';
  return 'actionNode';
}
