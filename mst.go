package delaunay

import (
	"sort"
)

func mst(edges []*Edge) []*Edge {
	var mst []*Edge
	indexMap := make(map[*Point]int)
	i := 0
	for _, edge := range edges {
		if _, ok := indexMap[edge.a]; !ok {
			indexMap[edge.a] = i
			i++
		}
		if _, ok := indexMap[edge.b]; !ok {
			indexMap[edge.b] = i
			i++
		}
	}

	uf := createUF(i)
	sort.Slice(edges, func(i, j int) bool { return edges[i].length() < edges[j].length() })
	for _, edge := range edges {
		a := indexMap[edge.a]
		b := indexMap[edge.b]
		if uf.find(a) != uf.find(b) {
			uf.union(a, b)
			mst = append(mst, edge)
		}
	}

	return mst
}

func EuclideanMST(points []*Point) []*Edge {
	edges := completeEdgeSet(points)
	mst := mst(edges)
	return mst
}

func completeEdgeSet(points []*Point) []*Edge {
	var edges []*Edge
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			edges = append(edges, &Edge{points[i], points[j]})
		}
	}
	return edges
}

type unionFind struct {
	parent []int
}

func createUF(n int) unionFind {
	var uf unionFind
	uf.parent = make([]int, n)
	for i := 0; i < n; i++ {
		uf.parent[i] = i
	}
	return uf
}
func (uf unionFind) find(a int) int {
	if uf.parent[a] != a {
		uf.parent[a] = uf.find(uf.parent[a])
	}
	return uf.parent[a]
}
func (uf unionFind) union(a, b int) {
	aRoot := uf.find(a)
	bRoot := uf.find(b)
	if aRoot == bRoot {
		return
	}
	uf.parent[aRoot] = bRoot
}
