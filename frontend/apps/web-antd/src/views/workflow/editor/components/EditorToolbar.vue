<script setup lang="ts">
interface Props {
  saving?: boolean;
  workflowName?: string;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'save'): void;
  (e: 'execute'): void;
  (e: 'undo'): void;
  (e: 'redo'): void;
}>();
</script>

<template>
  <div class="editor-toolbar">
    <div class="toolbar-left">
      <span class="workflow-name">{{ workflowName || '未命名工作流' }}</span>
    </div>
    <div class="toolbar-right">
      <button class="toolbar-btn" @click="emit('undo')" title="撤销">&#8617;</button>
      <button class="toolbar-btn" @click="emit('redo')" title="重做">&#8618;</button>
      <button class="toolbar-btn primary" :disabled="saving" @click="emit('save')">
        {{ saving ? '保存中...' : '保存' }}
      </button>
      <button class="toolbar-btn success" @click="emit('execute')">执行</button>
    </div>
  </div>
</template>

<style scoped>
.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  height: 48px;
}
.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.workflow-name {
  font-size: 15px;
  font-weight: 500;
  color: #333;
}
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}
.toolbar-btn {
  padding: 4px 14px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  font-size: 13px;
}
.toolbar-btn:hover {
  border-color: #409EFF;
  color: #409EFF;
}
.toolbar-btn.primary {
  background: #409EFF;
  color: #fff;
  border-color: #409EFF;
}
.toolbar-btn.success {
  background: #67C23A;
  color: #fff;
  border-color: #67C23A;
}
.toolbar-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
