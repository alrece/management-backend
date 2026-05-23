package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === 拓扑排序测试 ===

func TestTopologicalSort_LinearChain(t *testing.T) {
	nodes := []DAGNode{
		{ID: "a", Type: "trigger"},
		{ID: "b", Type: "action"},
		{ID: "c", Type: "action"},
	}
	edges := []DAGEdge{
		{Source: "a", Target: "b"},
		{Source: "b", Target: "c"},
	}

	result, err := TopologicalSort(nodes, edges)
	assert.NoError(t, err)
	assert.Len(t, result.Layers, 3)
	assert.Equal(t, "a", result.Layers[0][0].ID)
	assert.Equal(t, "b", result.Layers[1][0].ID)
	assert.Equal(t, "c", result.Layers[2][0].ID)
}

func TestTopologicalSort_ParallelNodes(t *testing.T) {
	nodes := []DAGNode{
		{ID: "a", Type: "trigger"},
		{ID: "b", Type: "action"},
		{ID: "c", Type: "action"},
		{ID: "d", Type: "action"},
	}
	edges := []DAGEdge{
		{Source: "a", Target: "b"},
		{Source: "a", Target: "c"},
		{Source: "b", Target: "d"},
		{Source: "c", Target: "d"},
	}

	result, err := TopologicalSort(nodes, edges)
	assert.NoError(t, err)
	assert.Len(t, result.Layers, 3)
	assert.Len(t, result.Layers[0], 1) // a
	assert.Len(t, result.Layers[1], 2) // b, c
	assert.Len(t, result.Layers[2], 1) // d
}

func TestTopologicalSort_DetectsCycle(t *testing.T) {
	nodes := []DAGNode{
		{ID: "a"}, {ID: "b"}, {ID: "c"},
	}
	edges := []DAGEdge{
		{Source: "a", Target: "b"},
		{Source: "b", Target: "c"},
		{Source: "c", Target: "a"},
	}

	_, err := TopologicalSort(nodes, edges)
	assert.Error(t, err)
	var cycleErr *CycleError
	assert.ErrorAs(t, err, &cycleErr)
}

func TestTopologicalSort_EmptyGraph(t *testing.T) {
	result, err := TopologicalSort(nil, nil)
	assert.NoError(t, err)
	assert.Empty(t, result.Layers)
}

func TestTopologicalSort_SingleNode(t *testing.T) {
	nodes := []DAGNode{{ID: "a"}}
	result, err := TopologicalSort(nodes, nil)
	assert.NoError(t, err)
	assert.Len(t, result.Layers, 1)
	assert.Equal(t, "a", result.Layers[0][0].ID)
}

func TestTopologicalSort_DisconnectedGraph(t *testing.T) {
	nodes := []DAGNode{
		{ID: "a"}, {ID: "b"}, {ID: "c"},
	}
	// a→b, c 独立
	edges := []DAGEdge{{Source: "a", Target: "b"}}

	result, err := TopologicalSort(nodes, edges)
	assert.NoError(t, err)
	// 第一层应该包含 a 和 c（都无入度）
	assert.Len(t, result.Layers[0], 2)
}

// === 循环检测测试 ===

func TestDetectCycle_NoCycle(t *testing.T) {
	nodes := []DAGNode{
		{ID: "a"}, {ID: "b"}, {ID: "c"},
	}
	edges := []DAGEdge{
		{Source: "a", Target: "b"},
		{Source: "b", Target: "c"},
	}

	err := DetectCycle(nodes, edges)
	assert.NoError(t, err)
}

func TestDetectCycle_SimpleCycle(t *testing.T) {
	nodes := []DAGNode{
		{ID: "a"}, {ID: "b"}, {ID: "c"},
	}
	edges := []DAGEdge{
		{Source: "a", Target: "b"},
		{Source: "b", Target: "c"},
		{Source: "c", Target: "a"},
	}

	err := DetectCycle(nodes, edges)
	assert.Error(t, err)
	var cycleErr *CycleError
	assert.ErrorAs(t, err, &cycleErr)
	assert.NotEmpty(t, cycleErr.CycleNodes)
}

func TestDetectCycle_SelfLoop(t *testing.T) {
	nodes := []DAGNode{{ID: "a"}}
	edges := []DAGEdge{{Source: "a", Target: "a"}}

	err := DetectCycle(nodes, edges)
	assert.Error(t, err)
}

func TestDetectCycle_TwoNodeCycle(t *testing.T) {
	nodes := []DAGNode{{ID: "a"}, {ID: "b"}}
	edges := []DAGEdge{
		{Source: "a", Target: "b"},
		{Source: "b", Target: "a"},
	}

	err := DetectCycle(nodes, edges)
	assert.Error(t, err)
}

func TestDetectCycle_NoEdges(t *testing.T) {
	nodes := []DAGNode{{ID: "a"}, {ID: "b"}, {ID: "c"}}

	err := DetectCycle(nodes, nil)
	assert.NoError(t, err)
}

// === 节点注册表测试 ===

func TestNodeRegistry_RegisterAndGet(t *testing.T) {
	registry := NewNodeRegistry()

	handler := &mockHandler{
		nodeType:   "test_node",
		category:   "action",
		label:      "测试节点",
	}
	registry.Register(handler)

	got, ok := registry.Get("test_node")
	assert.True(t, ok)
	assert.Equal(t, "test_node", got.Type())
}

func TestNodeRegistry_GetNotExists(t *testing.T) {
	registry := NewNodeRegistry()
	_, ok := registry.Get("nonexistent")
	assert.False(t, ok)
}

func TestNodeRegistry_AllSchemas(t *testing.T) {
	registry := NewNodeRegistry()
	registry.Register(&mockHandler{nodeType: "a", category: "action", label: "A"})
	registry.Register(&mockHandler{nodeType: "b", category: "control", label: "B"})

	schemas := registry.AllSchemas()
	assert.Len(t, schemas, 2)
}

func TestNodeRegistry_SchemasByCategory(t *testing.T) {
	registry := NewNodeRegistry()
	registry.Register(&mockHandler{nodeType: "a", category: "action", label: "A"})
	registry.Register(&mockHandler{nodeType: "b", category: "control", label: "B"})
	registry.Register(&mockHandler{nodeType: "c", category: "action", label: "C"})

	actionSchemas := registry.SchemasByCategory("action")
	assert.Len(t, actionSchemas, 2)
}

func TestNodeRegistry_DuplicateRegister(t *testing.T) {
	registry := NewNodeRegistry()
	registry.Register(&mockHandler{nodeType: "dup", category: "action", label: "Dup"})

	assert.Panics(t, func() {
		registry.Register(&mockHandler{nodeType: "dup", category: "action", label: "Dup2"})
	})
}

// === 下游节点查询测试 ===

func TestGetDownstreamNodes(t *testing.T) {
	edges := []DAGEdge{
		{Source: "a", Target: "b"},
		{Source: "a", Target: "c"},
		{Source: "b", Target: "d"},
		{Source: "c", Target: "d"},
	}

	downstream := GetDownstreamNodes("a", edges)
	assert.Len(t, downstream, 3)
	assert.Contains(t, downstream, "b")
	assert.Contains(t, downstream, "c")
	assert.Contains(t, downstream, "d")
}

func TestGetDownstreamNodes_NoDownstream(t *testing.T) {
	edges := []DAGEdge{{Source: "a", Target: "b"}}
	downstream := GetDownstreamNodes("b", edges)
	assert.Empty(t, downstream)
}

// === mock handler ===

type mockHandler struct {
	nodeType string
	category string
	label    string
}

func (m *mockHandler) Execute(_ context.Context, _ *NodeInput) (*NodeOutput, error) {
	return &NodeOutput{Items: []map[string]any{{"result": "ok"}}}, nil
}

func (m *mockHandler) Type() string     { return m.nodeType }
func (m *mockHandler) Category() string  { return m.category }

func (m *mockHandler) Schema() NodeSchema {
	return NodeSchema{Type: m.nodeType, Label: m.label, Category: m.category}
}
