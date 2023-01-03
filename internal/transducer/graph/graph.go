package graph

import "math"

type GraphSpec interface {
	Vertices() []Vertex
	Neighbors(v Vertex) []Vertex
	Weight(u, v Vertex) float64
}

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

var Infinity = math.Inf(0)

func FloydWarshall(g Graph) (dist map[Vertex]map[Vertex]float64, next map[Vertex]map[Vertex]*Vertex) {
	vertexes := g.Vertices()
	dist = make(map[Vertex]map[Vertex]float64)
	next = make(map[Vertex]map[Vertex]*Vertex)

	for _, u := range vertexes {
		dist[u] = make(map[Vertex]float64)
		next[u] = make(map[Vertex]*Vertex)
		for _, v := range vertexes {
			dist[u][v] = Infinity
		}
		dist[u][u] = 0
		for _, v := range g.Neighbors(u) {
			v := v
			dist[u][v] = g.Weight(u, v)
			next[u][v] = &v
		}
	}

	for _, k := range vertexes {
		for _, i := range vertexes {
			for _, j := range vertexes {
				if dist[i][k] < Infinity && dist[k][j] < Infinity {
					if dist[i][j] > dist[i][k]+dist[k][j] {
						dist[i][j] = dist[i][k] + dist[k][j]
						next[i][j] = next[i][k]
					}
				}
			}
		}
	}
	return dist, next
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
