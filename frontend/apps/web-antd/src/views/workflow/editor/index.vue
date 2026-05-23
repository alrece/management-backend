<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import { VueFlow } from '@vue-flow/core';

import EditorToolbar from './components/EditorToolbar.vue';
import NodePalette from './components/NodePalette.vue';
import FlowCanvas from './components/FlowCanvas.vue';
import PropertyPanel from './components/PropertyPanel.vue';
import { useDragAndDrop } from './composables/useDragAndDrop';
import { useWorkflow } from './composables/useWorkflow';
import { useNodeSchema } from './composables/useNodeSchema';

const route = useRoute();
const workflowId = ref<number | null>(null);
const workflowName = ref('');
const selectedNodeId = ref<string | null>(null);

if (route.params.id && route.params.id !== 'new') {
  workflowId.value = Number(route.params.id);
}

const { onDragStart, onDragOver, onDrop } = useDragAndDrop();
const { saving, loading, load, save, execute } = useWorkflow(workflowId);
const { schemas, fetchSchemas } = useNodeSchema();

onMounted(async () => {
  await fetchSchemas();
  if (workflowId.value) {
    await load();
  }
});

async function handleSave() {
  await save(workflowName.value || '未命名工作流');
}

const selectedNode = computed(() => {
  if (!selectedNodeId.value) return null;
  return null;
});

function handleNodeUpdate(nodeId: string, data: Record<string, any>) {
  // 节点属性更新通过 Vue Flow 内部状态
}
</script>

<template>
  <div class="workflow-editor">
    <EditorToolbar
      :saving="saving"
      :workflow-name="workflowName"
      @save="handleSave"
      @execute="execute"
      @undo="() => {}"
      @redo="() => {}"
    />
    <div class="editor-body">
      <NodePalette
        :schemas="schemas"
        @dragstart="onDragStart"
      />
      <div class="canvas-wrapper" @dragover="onDragOver" @drop="onDrop">
        <VueFlow
          :default-viewport="{ zoom: 1, x: 0, y: 0 }"
          fit-view-on-init
          @node-click="(_event: any, node: any) => { selectedNodeId = node.id }"
          @pane-click="selectedNodeId = null"
        >
          <template #node-triggerNode="nodeProps">
            <component :is="($options.components?.TriggerNode)" v-bind="nodeProps" />
          </template>
          <template #node-actionNode="nodeProps">
            <component :is="($options.components?.ActionNode)" v-bind="nodeProps" />
          </template>
          <template #node-controlNode="nodeProps">
            <component :is="($options.components?.ControlNode)" v-bind="nodeProps" />
          </template>
          <template #node-transformNode="nodeProps">
            <component :is="($options.components?.TransformNode)" v-bind="nodeProps" />
          </template>
        </VueFlow>
      </div>
      <PropertyPanel
        :node="null"
        @update="handleNodeUpdate"
      />
    </div>
  </div>
</template>

<script lang="ts">
import TriggerNode from './components/nodes/TriggerNode.vue';
import ActionNode from './components/nodes/ActionNode.vue';
import ControlNode from './components/nodes/ControlNode.vue';
import TransformNode from './components/nodes/TransformNode.vue';

export default {
  components: { TriggerNode, ActionNode, ControlNode, TransformNode },
};
</script>

<style scoped>
.workflow-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
}
.editor-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}
.canvas-wrapper {
  flex: 1;
  position: relative;
}
</style>
