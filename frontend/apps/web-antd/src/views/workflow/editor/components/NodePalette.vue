<script setup lang="ts">
import type { WorkflowApi } from '#/api/workflow/workflow';
import type { NodeCategory } from '../types';

import { CATEGORY_LABELS, CATEGORY_COLORS } from '../types';

interface Props {
  schemas: WorkflowApi.NodeSchema[];
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'dragstart', event: DragEvent, type: string, category: NodeCategory): void;
}>();

const categories: NodeCategory[] = ['trigger', 'action', 'control', 'transform'];
</script>

<template>
  <div class="node-palette">
    <h4 style="margin: 0 0 12px; font-size: 14px; color: #333">节点面板</h4>
    <div v-for="cat in categories" :key="cat" class="category-group">
      <div class="category-title" :style="{ color: CATEGORY_COLORS[cat] }">
        {{ CATEGORY_LABELS[cat] }}
      </div>
      <div
        v-for="schema in schemas.filter((s) => s.category === cat)"
        :key="schema.type"
        class="palette-item"
        :style="{ borderLeftColor: CATEGORY_COLORS[cat] }"
        draggable="true"
        @dragstart="(e) => emit('dragstart', e as DragEvent, schema.type, cat)"
      >
        <span class="item-icon">{{ schema.icon || '&#9679;' }}</span>
        <span class="item-label">{{ schema.label }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.node-palette {
  width: 200px;
  padding: 12px;
  background: #fafafa;
  border-right: 1px solid #e8e8e8;
  overflow-y: auto;
  height: 100%;
}
.category-group {
  margin-bottom: 12px;
}
.category-title {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 6px;
  text-transform: uppercase;
}
.palette-item {
  padding: 6px 10px;
  margin-bottom: 4px;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-left: 3px solid;
  border-radius: 4px;
  cursor: grab;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.palette-item:hover {
  background: #f0f7ff;
}
.palette-item:active {
  cursor: grabbing;
}
</style>
