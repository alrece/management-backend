<script setup lang="ts">
import { VueFlow, useVueFlow } from '@vue-flow/core';
import { Background } from '@vue-flow/background';
import { MiniMap } from '@vue-flow/minimap';
import { Controls } from '@vue-flow/controls';
import '@vue-flow/core/dist/style.css';
import '@vue-flow/core/dist/theme-default.css';
import '@vue-flow/background/dist/style.css';
import '@vue-flow/controls/dist/style.css';
import '@vue-flow/minimap/dist/style.css';

import TriggerNode from './nodes/TriggerNode.vue';
import ActionNode from './nodes/ActionNode.vue';
import ControlNode from './nodes/ControlNode.vue';
import TransformNode from './nodes/TransformNode.vue';

interface Props {
  dragOver: (event: DragEvent) => void;
  drop: (event: DragEvent) => void;
}

defineProps<Props>();

const { onPaneClick } = useVueFlow();
</script>

<template>
  <div class="flow-canvas" @dragover="dragOver" @drop="drop">
    <VueFlow
      :default-viewport="{ zoom: 1, x: 0, y: 0 }"
      :min-zoom="0.2"
      :max-zoom="4"
      fit-view-on-init
      :node-types="{ triggerNode: TriggerNode, actionNode: ActionNode, controlNode: ControlNode, transformNode: TransformNode }"
    >
      <Background pattern-color="#e8e8e8" :gap="20" />
      <MiniMap />
      <Controls />
    </VueFlow>
  </div>
</template>

<style scoped>
.flow-canvas {
  flex: 1;
  height: 100%;
}
</style>
