package engine

import (
	"fmt"
)

// DAGNode DAG 节点（简化版，仅含 ID 和类型）
type DAGNode struct {
	ID   string
	Type string
}

// DAGEdge DAG 边
type DAGEdge struct {
	Source       string
	Target       string
	SourceHandle string // 条件分支用（true/false）
}

// CycleError 循环检测错误（包含参与环的节点 ID）
type CycleError struct {
	CycleNodes []string
}

func (e *CycleError) Error() string {
	return fmt.Sprintf("检测到循环依赖，参与节点: %v", e.CycleNodes)
}

// LayeredResult 拓扑排序分层结果
type LayeredResult struct {
	Layers [][]DAGNode // 按层级组织，同层节点可并行执行
}

// TopologicalSort 拓扑排序（Kahn 算法），返回分层结果
// 同一层内的节点无依赖关系，可并行执行
func TopologicalSort(nodes []DAGNode, edges []DAGEdge) (*LayeredResult, error) {
	if len(nodes) == 0 {
		return &LayeredResult{}, nil
	}

	// 构建 ID → DAGNode 映射
	nodeMap := make(map[string]DAGNode, len(nodes))
	for _, n := range nodes {
		nodeMap[n.ID] = n
	}

	// 计算入度
	inDegree := make(map[string]int, len(nodes))
	adjList := make(map[string][]DAGEdge, len(nodes))

	for _, n := range nodes {
		inDegree[n.ID] = 0
		adjList[n.ID] = nil
	}

	for _, e := range edges {
		inDegree[e.Target]++
		adjList[e.Source] = append(adjList[e.Source], e)
	}

	// BFS 分层处理
	var layers [][]DAGNode
	queue := make([]string, 0)

	// 收集入度为 0 的节点作为第一层
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	visited := 0
	for len(queue) > 0 {
		// 当前层的所有节点
		var layer []DAGNode
		size := len(queue)
		for i := range size {
			id := queue[i]
			layer = append(layer, nodeMap[id])
			visited++

			// 减少下游节点的入度
			for _, edge := range adjList[id] {
				inDegree[edge.Target]--
				if inDegree[edge.Target] == 0 {
					queue = append(queue, edge.Target)
				}
			}
		}

		layers = append(layers, layer)
		queue = queue[size:]
	}

	// 如果还有未访问的节点，说明存在环
	if visited < len(nodes) {
		var cycleNodes []string
		for id, deg := range inDegree {
			if deg > 0 {
				cycleNodes = append(cycleNodes, id)
			}
		}
		return nil, &CycleError{CycleNodes: cycleNodes}
	}

	return &LayeredResult{Layers: layers}, nil
}

// DetectCycle 检测循环依赖（DFS 颜色标记法）
// 返回 nil 表示无环，否则返回参与环的节点列表
func DetectCycle(nodes []DAGNode, edges []DAGEdge) error {
	const (
		white = 0 // 未访问
		gray  = 1 // 访问中
		black = 2 // 已完成
	)

	color := make(map[string]int, len(nodes))
	parent := make(map[string]string, len(nodes))
	adjList := make(map[string][]DAGEdge, len(nodes))

	for _, n := range nodes {
		color[n.ID] = white
	}

	for _, e := range edges {
		adjList[e.Source] = append(adjList[e.Source], e)
	}

	var cycleNodes []string
	var found bool

	var dfs func(nodeID string)
	dfs = func(nodeID string) {
		if found {
			return
		}
		color[nodeID] = gray

		for _, edge := range adjList[nodeID] {
			target := edge.Target
			switch color[target] {
			case gray:
				// 发现环，回溯环路径
				found = true
				cycleNodes = []string{target, nodeID}
				cur := nodeID
				for parent[cur] != "" && parent[cur] != target {
					cur = parent[cur]
					cycleNodes = append(cycleNodes, cur)
				}
				return
			case white:
				parent[target] = nodeID
				dfs(target)
				if found {
					return
				}
			}
		}

		color[nodeID] = black
	}

	for _, n := range nodes {
		if color[n.ID] == white {
			dfs(n.ID)
			if found {
				return &CycleError{CycleNodes: cycleNodes}
			}
		}
	}

	return nil
}

// GetDownstreamNodes 获取指定节点的所有下游节点 ID
func GetDownstreamNodes(nodeID string, edges []DAGEdge) []string {
	seen := make(map[string]bool)
	var result []string
	queue := []string{nodeID}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, e := range edges {
			if e.Source == cur && !seen[e.Target] {
				seen[e.Target] = true
				result = append(result, e.Target)
				queue = append(queue, e.Target)
			}
		}
	}

	return result
}
