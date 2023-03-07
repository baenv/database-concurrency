package graph

import "math"

type GraphSpec interface {
	Vertices() []Vertex
	Neighbors(v Vertex) []Vertex
	Weight(u, v Vertex) float64
}

var Infinity = math.Inf(0)

type Vertex string

type Graph struct {
	vertexes []Vertex
	edges    map[Vertex]map[Vertex]float64
}

func NewGraph(vertexes []Vertex, edges map[Vertex]map[Vertex]float64) Graph {
	return Graph{vertexes, edges}
}

func (g Graph) AddEdge(from Vertex, to Vertex, weight float64) {
	if _, ok := g.edges[from]; !ok {
		g.edges[from] = make(map[Vertex]float64)
	}
	g.edges[from][to] = weight
}

func (g Graph) Vertices() []Vertex {
	return g.vertexes
}

func (g Graph) Neighbors(v Vertex) (vs []Vertex) {
	for k := range g.edges[v] {
		vs = append(vs, k)
	}
	return vs
}

func (g Graph) Weight(u, v Vertex) float64 {
	return g.edges[u][v]
}

func GetPaths(u Vertex, v Vertex, next map[Vertex]map[Vertex]*Vertex) (paths []Vertex) {
	if next[u][v] == nil {
		return
	}
	paths = []Vertex{u}
	for u != v {
		u = *next[u][v]
		paths = append(paths, u)
	}
	return paths
}

func PrintPaths(vv []Vertex) (s string) {
	if len(vv) == 0 {
		return ""
	}
	s = string(vv[0])
	for _, v := range vv[1:] {
		s += " -> " + string(v)
	}
	return s
}
