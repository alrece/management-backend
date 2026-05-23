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
  <div class="control-node" :class="{ selected }">
    <Handle type="target" :position="Position.Top" />
    <Handle type="source" :position="Position.Bottom" id="true" />
    <Handle type="source" :position="Position.Right" id="false" />
    <div class="node-header" style="background: #E6A23C">
      <span class="node-icon">&#9874;</span>
      <span class="node-label">{{ data.label || '控制' }}</span>
    </div>
  </div>
</template>

<style scoped>
.control-node {
  background: #fff;
  border: 2px solid #E6A23C;
  border-radius: 8px;
  min-width: 160px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
.control-node.selected {
  box-shadow: 0 0 0 2px #E6A23C;
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
