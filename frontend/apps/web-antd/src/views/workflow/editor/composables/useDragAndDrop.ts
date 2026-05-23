import type { GraphNode } from '@vue-flow/core';

import { useVueFlow } from '@vue-flow/core';
import { ref } from 'vue';

import type { NodeCategory } from '../types';

export function useDragAndDrop() {
  const { addNodes, screenToFlowCoordinate, onConnect, addEdges } = useVueFlow();
  const dragNodeType = ref<string | null>(null);
  const dragNodeCategory = ref<NodeCategory | null>(null);

  function onDragStart(event: DragEvent, type: string, category: NodeCategory) {
    if (event.dataTransfer) {
      event.dataTransfer.setData('application/vueflow', type);
      event.dataTransfer.effectAllowed = 'move';
    }
    dragNodeType.value = type;
    dragNodeCategory.value = category;
  }

  function onDragOver(event: DragEvent) {
    event.preventDefault();
    if (event.dataTransfer) {
      event.dataTransfer.dropEffect = 'move';
    }
  }

  function onDrop(event: DragEvent) {
    const type = event.dataTransfer?.getData('application/vueflow') || dragNodeType.value;
    if (!type) return;

    const { left, top } = (
      event.target as HTMLElement
    ).getBoundingClientRect();
    const position = screenToFlowCoordinate({ x: event.clientX - left, y: event.clientY - top });

    const newNode: GraphNode = {
      id: `${type}_${Date.now()}`,
      type: getNodeComponentType(dragNodeCategory.value || 'action'),
      position,
      data: { label: type, config: {}, errorStrategy: 'retry', retryCount: 3, timeout: 60 },
    };

    addNodes([newNode] as any);
  }

  function getNodeComponentType(category: string): string {
    switch (category) {
      case 'trigger': return 'triggerNode';
      case 'control': return 'controlNode';
      case 'transform': return 'transformNode';
      default: return 'actionNode';
    }
  }

  return { onDragStart, onDragOver, onDrop, dragNodeType, dragNodeCategory };
}
