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

type adjacencyGraph map[*Point]*adjacencyList

func newAdjacencyGraph(points []*Point) adjacencyGraph {
	graph := make(adjacencyGraph)
	for _, p := range points {
		newList := newAdjacencyList(p)
		graph[p] = &(newList)
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
	root   *adjacencyListNode
}
type adjacencyListNode struct {
	point          *Point
	previous, next *adjacencyListNode
}

func newAdjacencyList(center *Point) adjacencyList {
	al := adjacencyList{
		center: center,
		root:   &adjacencyListNode{},
	}
  al.root.previous = al.root
  al.root.next = al.root
	return al
}

func (al *adjacencyList) insert(newPoint *Point) {
  //FIXME: does not deal with duplicate points
  //FIXME: does not deal with points lying on a line with a point currently in the set
	newNode := &adjacencyListNode{newPoint, nil, nil}
  var insertionPoint *adjacencyListNode = al.root
  for thisNeighbor := al.root.next; thisNeighbor != al.root; thisNeighbor = thisNeighbor.next {
    cross := crossProduct(al.center, newPoint, thisNeighbor.point)
    if (cross < 0) {
      insertionPoint = thisNeighbor
    } else if (insertionPoint != al.root) {
      break
    }
  }
  newNode.previous = insertionPoint
  newNode.next = insertionPoint.next
  insertionPoint.next.previous = newNode
  insertionPoint.next = newNode
}

func (al adjacencyList) remove(oldPoint *Point) {

}

func (al adjacencyList) neighbors() []*Point {
	neighbors := make([]*Point, 0)
	for thisNeighbor := al.root.next; thisNeighbor != al.root; thisNeighbor = thisNeighbor.next {
		neighbors = append(neighbors, thisNeighbor.point)
	}
	return neighbors
}
