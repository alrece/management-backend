<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core';

interface Props {
  data: { label: string; config: Record<string, any> };
  id: string;
  selected?: boolean;
}

defineProps<Props>();
</script>

<template>
  <div class="action-node" :class="{ selected }">
    <Handle type="target" :position="Position.Top" />
    <Handle type="source" :position="Position.Bottom" />
    <div class="node-header" style="background: #409EFF">
      <span class="node-icon">&#9881;</span>
      <span class="node-label">{{ data.label || '动作' }}</span>
    </div>
  </div>
</template>

<style scoped>
.action-node {
  background: #fff;
  border: 2px solid #409EFF;
  border-radius: 8px;
  min-width: 160px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
.action-node.selected {
  box-shadow: 0 0 0 2px #409EFF;
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
