package delaunay

import (
	"math/rand"
	"time"
)

func Triangulate(points []*Point) []*Edge {
	graph := newAdjacencyGraph(points)

	edges := graph.toEdgeSlice()
	return edges
}

func BadTriangulate(points []*Point) []*Edge {
	var edges []*Edge
	randEdges := make([]*Edge, 0)

	for i, p1 := range points {
		for j := i + 1; j < len(points); j++ {
			randEdges = append(randEdges, &Edge{p1, points[j]})
		}
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < len(randEdges); i++ {
		j := i + r.Intn(len(randEdges)-i)
		randEdges[i], randEdges[j] = randEdges[j], randEdges[i]
	}

	for _, e := range randEdges {
		isGood := true
		for _, f := range edges {
			if e.a != f.a && e.a != f.b && e.b != f.a && e.b != f.b &&
				doLineSegmentsIntersect(f.a, f.b, e.a, e.b) {
				isGood = false
			}
		}
		if isGood {
			edges = append(edges, e)
		}
	}

	return edges
}

/*
  input two convex hulls, finds lower common tangent or upper common tangent
  input:  slices of points - the left and right hulls to merge
          adjacency lists for points
  output: pointers totangent endpoints
*/
func lct() (*Point, *Point) {
	return nil, nil
}
func uct() (*Point, *Point) {
	return nil, nil
}

type adjacencyGraph map[*Point]adjacencyList

func newAdjacencyGraph(points []*Point) adjacencyGraph {
	graph := make(adjacencyGraph)
	for _, p := range points {
		graph[p] = adjacencyList{
			center: p,
			first:  nil,
		}
	}
	return graph
}

func (ag adjacencyGraph) toEdgeSlice() []*Edge {
	edges := make([]*Edge, 0)
	seenPoints := make(map[*Point]bool)
	for point, adjList := range ag {
		seenPoints[point] = true
		for _, neighbor := range adjList.neighbors() {
			if seenPoints[neighbor] {
				edges = append(edges, &Edge{point, neighbor})
			}
		}
	}
	return edges
}

func (graph adjacencyGraph) addEdge(p1, p2 *Point) {
	graph[p1].insert(p2)
	graph[p2].insert(p1)
}
func (graph adjacencyGraph) deleteEdge(p1, p2 *Point) {
	graph[p1].remove(p2)
	graph[p2].remove(p1)
}

type adjacencyList struct {
	center *Point
	first  *adjacencyListNode
}
type adjacencyListNode struct {
	point          *Point
	previous, next *adjacencyListNode
}

func (al adjacencyList) insert(newPoint *Point) {

}

func (al adjacencyList) remove(oldPoint *Point) {

}

func (al adjacencyList) neighbors() []*Point {
	neighbors := make([]*Point, 0)
	if al.first == nil {
		return neighbors
	}
	neighbors = append(neighbors, al.first.point)
	for thisNeighbor := al.first.next; thisNeighbor != al.first; thisNeighbor = thisNeighbor.next {
		neighbors = append(neighbors, thisNeighbor.point)
	}

	return neighbors
}
