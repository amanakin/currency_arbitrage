package graph

import (
	"fmt"
	"math"
)

type Vertex struct {
	SliceIndex int
	Name       string
}

type Edge struct {
	From int
	To   int
	Cost float64
}

type Graph struct {
	vertexes []Vertex
	edges    []Edge

	nameToIndex map[string]int
}

func NewGraph() *Graph {
	return &Graph{
		nameToIndex: make(map[string]int),
	}
}

func (g *Graph) InsertVertex(Name string) {
	v := Vertex{
		SliceIndex: len(g.vertexes),
		Name:       Name,
	}

	g.vertexes = append(g.vertexes, v)
	g.nameToIndex[Name] = v.SliceIndex
}

func (g *Graph) InsertEdge(from, to string, cost float64) {
	fromIndex, ok := g.nameToIndex[from]
	if !ok {
		return
	}
	toIndex, ok := g.nameToIndex[to]
	if !ok {
		return
	}

	g.edges = append(g.edges, Edge{
		From: fromIndex,
		To:   toIndex,
		Cost: cost,
	})
}

func (g *Graph) FindNegativeCycleFromStart(name string) []Vertex {
	start := g.nameToIndex[name]

	dists := make([]float64, len(g.vertexes)) // distance from start
	prev := make([]int, len(g.vertexes))
	for i := range dists {
		dists[i] = math.MaxFloat64
		prev[i] = -1
	}
	dists[start] = 0

	for i := 0; i < len(g.vertexes)-1; i++ {
		for _, edge := range g.edges {
			if dists[edge.From]+edge.Cost >= dists[edge.To] {
				continue
			}

			dists[edge.To] = dists[edge.From] + edge.Cost
			prev[edge.To] = edge.From
		}
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("graph=%+v\n", g)
			fmt.Printf("prev= %v\n", prev)
			fmt.Printf("dists=%v\n", dists)
			panic(r)
		}
	}()

	for _, edge := range g.edges {
		if dists[edge.From]+edge.Cost >= dists[edge.To] {
			continue
		}
		nodes := make(map[int]bool)
		curr := edge.To

		for {
			// Cycle not in start vertex
			if curr < 0 {
				nodes = nil
				break
			}

			if !nodes[curr] {
				nodes[curr] = true
				curr = prev[curr]
				continue
			}

			break
		}

		if !nodes[start] {
			continue
		}

		result := []Vertex{g.vertexes[start]}
		curr = start

		for {
			curr = prev[curr]
			result = append(result, g.vertexes[curr])
			if curr == start {
				break
			}
		}

		// Reverse
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}
		return result
	}

	return nil
}
