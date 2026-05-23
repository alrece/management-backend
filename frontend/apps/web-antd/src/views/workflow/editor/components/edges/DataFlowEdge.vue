<script setup lang="ts">
import { getBezierPath } from '@vue-flow/core';
import type { ConnectionLineProps } from '@vue-flow/core';

const props = defineProps<ConnectionLineProps>();
const path = getBezierPath({
  sourceX: props.sourceX,
  sourceY: props.sourceY,
  sourcePosition: props.sourcePosition,
  targetX: props.targetX,
  targetY: props.targetY,
  targetPosition: props.targetPosition,
});
</script>

<template>
  <g>
    <path :d="path" class="vue-flow__edge-path" style="stroke: #409EFF; stroke-width: 2" />
    <circle :cx="targetX" :cy="targetY" r="4" fill="#409EFF" />
  </g>
</template>
