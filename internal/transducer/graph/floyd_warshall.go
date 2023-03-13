package graph

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
