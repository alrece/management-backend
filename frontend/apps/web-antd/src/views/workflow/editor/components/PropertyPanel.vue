<script setup lang="ts">
import type { GraphNode } from '@vue-flow/core';

import { computed } from 'vue';
import { Form, Input, InputNumber, Select, Switch } from 'ant-design-vue';

const props = defineProps<{
  node: GraphNode | null;
}>();

const emit = defineEmits<{
  (e: 'update', nodeId: string, data: Record<string, any>): void;
}>();

const formData = computed(() => {
  if (!props.node?.data) return {} as Record<string, any>;
  return props.node.data as Record<string, any>;
});

const errorStrategies = [
  { label: '重试', value: 'retry' },
  { label: '跳过', value: 'skip' },
  { label: '终止', value: 'abort' },
];

function handleUpdate(field: string, value: any) {
  if (!props.node) return;
  emit('update', props.node.id, { ...formData.value, [field]: value });
}
</script>

<template>
  <div class="property-panel">
    <template v-if="node">
      <h4 style="margin: 0 0 12px; font-size: 14px; color: #333">属性面板</h4>
      <Form layout="vertical" size="small">
        <Form.Item label="节点名称">
          <Input
            :value="formData.label"
            @change="(e: any) => handleUpdate('label', e.target.value)"
          />
        </Form.Item>
        <Form.Item label="错误策略">
          <Select
            :value="formData.errorStrategy"
            :options="errorStrategies"
            @change="(v: string) => handleUpdate('errorStrategy', v)"
          />
        </Form.Item>
        <Form.Item label="重试次数">
          <InputNumber
            :value="formData.retryCount"
            :min="0"
            :max="10"
            @change="(v: number) => handleUpdate('retryCount', v)"
            style="width: 100%"
          />
        </Form.Item>
        <Form.Item label="超时(秒)">
          <InputNumber
            :value="formData.timeout"
            :min="1"
            :max="300"
            @change="(v: number) => handleUpdate('timeout', v)"
            style="width: 100%"
          />
        </Form.Item>
      </Form>
    </template>
    <div v-else class="empty-hint">选择节点查看属性</div>
  </div>
</template>

<style scoped>
.property-panel {
  width: 260px;
  padding: 12px;
  background: #fafafa;
  border-left: 1px solid #e8e8e8;
  overflow-y: auto;
  height: 100%;
}
.empty-hint {
  color: #999;
  text-align: center;
  padding-top: 40px;
  font-size: 13px;
}
</style>
