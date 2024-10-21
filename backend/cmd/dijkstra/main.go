package main

import (
	"container/heap"
	"fmt"
)

// Estrutura que representa uma aresta do grafo
type Edge struct {
	to     int     // vértice destino
	weight float64 // peso da aresta
}

// Estrutura que representa o grafo como uma lista de adjacências
type Graph struct {
	nodes map[int][]Edge
}

// Adiciona uma aresta ao grafo
func (g *Graph) addEdge(from, to int, weight float64) {
	if g.nodes == nil {
		g.nodes = make(map[int][]Edge)
	}
	g.nodes[from] = append(g.nodes[from], Edge{to, weight})
}

// Item é uma estrutura usada na fila de prioridade
type Item struct {
	node     int     // vértice
	priority float64 // prioridade (menor caminho)
	index    int     // índice do heap
	path     []int   // caminho percorrido até o nó
}

// PriorityQueue implementa uma fila de prioridade (menor prioridade tem precedência)
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
	old[n-1] = nil  // evitar memory leak
	item.index = -1 // segurança
	*pq = old[0 : n-1]
	return item
}

// Função para calcular os 3 melhores caminhos com Dijkstra entre dois pontos
func dijkstraKBest(g *Graph, start, target int, k int) []Path {
	var results []Path
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &Item{node: start, priority: 0, path: []int{start}})

	for pq.Len() > 0 && len(results) < k {
		// Pega o nó com menor distância
		current := heap.Pop(&pq).(*Item)
		currentNode := current.node
		currentPath := current.path
		currentWeight := current.priority

		// Se chegamos ao destino, adiciona o caminho e o peso aos resultados
		if currentNode == target {
			// Monta o resultado no formato desejado
			pathMap := Path{
				Path:   formatPath(currentPath, g),
				Weight: currentWeight,
			}
			results = append(results, pathMap)
			if len(results) == k {
				break
			}
		}

		// Atualiza as distâncias dos vizinhos e continua a busca
		for _, edge := range g.nodes[currentNode] {
			newPath := append([]int(nil), currentPath...)
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

type Route struct {
	Source int
	Target int
}

// Função auxiliar para formatar o caminho no formato [{1:2}, {2:3}, ...]
func formatPath(path []int, g *Graph) []Route {
	var formattedPath []Route
	for i := 0; i < len(path)-1; i++ {
		formattedPath = append(formattedPath, Route{
			Source: path[i],
			Target: path[i+1],
		})
	}
	return formattedPath
}

func main() {
	// Exemplo de uso
	g := &Graph{}
	g.addEdge(1, 2, 1)
	g.addEdge(1, 3, 4)
	g.addEdge(2, 3, 2)
	g.addEdge(2, 4, 6)
	g.addEdge(3, 4, 3)

	// Calcula os 3 melhores caminhos entre o nó 1 e o nó 4
	k := 3
	bestPaths := dijkstraKBest(g, 1, 10, k)
	if len(bestPaths) == 0 {
		fmt.Println("Nenhum caminho encontrado")
		return
	}

	// Exibe os caminhos e os pesos
	for i, result := range bestPaths[0].Path {
		fmt.Printf("%d: Caminho: %v, Peso: %f\n", i+1, result, bestPaths[0].Weight)
	}
}
