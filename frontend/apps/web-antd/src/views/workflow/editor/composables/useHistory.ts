import { ref } from 'vue';
import type { GraphData } from '@vue-flow/core';
import { useVueFlow } from '@vue-flow/core';

interface Command {
  label: string;
  snapshot: GraphData;
}

export function useHistory() {
  const { toObject, fromJSON } = useVueFlow();
  const undoStack = ref<Command[]>([]);
  const redoStack = ref<Command[]>([]);
  const maxSize = 50;

  function push(label: string) {
    const snapshot = toObject() as GraphData;
    undoStack.value.push({ label, snapshot });
    if (undoStack.value.length > maxSize) {
      undoStack.value.shift();
    }
    redoStack.value = [];
  }

  function undo() {
    const cmd = undoStack.value.pop();
    if (!cmd) return;
    redoStack.value.push({ label: cmd.label, snapshot: toObject() as GraphData });
    fromJSON(cmd.snapshot as any);
  }

  function redo() {
    const cmd = redoStack.value.pop();
    if (!cmd) return;
    undoStack.value.push({ label: cmd.label, snapshot: toObject() as GraphData });
    fromJSON(cmd.snapshot as any);
  }

  const canUndo = ref(true);
  const canRedo = ref(true);

  return { undo, redo, push, canUndo, canRedo, undoStack, redoStack };
}
