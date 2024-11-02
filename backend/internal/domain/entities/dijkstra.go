package entities

import (
	"container/heap"
)

type Edge struct {
	to     string
	weight float64
}

type Graph struct {
	nodes map[string][]Edge
}

func (g *Graph) AddEdge(from, to string, weight float64) {
	if g.nodes == nil {
		g.nodes = make(map[string][]Edge)
	}
	g.nodes[from] = append(g.nodes[from], Edge{to, weight})
	g.nodes[to] = append(g.nodes[to], Edge{from, weight})
}

type Item struct {
	node     string
	priority float64
	index    int
	path     []string
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (g *Graph) DijkstraKBest(start, target string, k int) []Path {
	var results []Path
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &Item{node: start, priority: 0, path: []string{start}})

	for pq.Len() > 0 && len(results) < k {
		current := heap.Pop(&pq).(*Item)
		currentNode := current.node
		currentPath := current.path
		currentWeight := current.priority

		if currentNode == target {
			pathMap := Path{
				Path:   formatPath(currentPath, g),
				Weight: currentWeight,
			}
			results = append(results, pathMap)
			if len(results) == k {
				break
			}
		}

		for _, edge := range g.nodes[currentNode] {
			newPath := append([]string(nil), currentPath...)
			newPath = append(newPath, edge.to)
			heap.Push(&pq, &Item{
				node:     edge.to,
				priority: currentWeight + edge.weight,
				path:     newPath,
			})
		}
	}

	return results
}

type Path struct {
	Path   []Route
	Weight float64
}

func formatPath(path []string, g *Graph) []Route {
	var formattedPath []Route
	for i := 0; i < len(path)-1; i++ {
		formattedPath = append(formattedPath, Route{
			Source: path[i],
			Target: path[i+1],
		})
	}
	return formattedPath
}