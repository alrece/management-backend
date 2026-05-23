<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core';

interface Props {
  data: { label: string; config: Record<string, any> };
  id: string;
  selected?: boolean;
}

const props = defineProps<Props>();
</script>

<template>
  <div class="trigger-node" :class="{ selected }">
    <Handle type="source" :position="Position.Bottom" />
    <div class="node-header" style="background: #67C23A">
      <span class="node-icon">&#9654;</span>
      <span class="node-label">{{ data.label || '触发器' }}</span>
    </div>
  </div>
</template>

<style scoped>
.trigger-node {
  background: #fff;
  border: 2px solid #67C23A;
  border-radius: 8px;
  min-width: 160px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
.trigger-node.selected {
  box-shadow: 0 0 0 2px #67C23A;
}
.node-header {
  padding: 8px 12px;
  border-radius: 6px 6px 0 0;
  color: #fff;
  font-size: 13px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
}
.node-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
