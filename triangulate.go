package delaunay

import (
	"math/rand"
	"time"
  "sort"
)

func Triangulate(points []*Point) []*Edge {
  pointsCopy := make([]*Point, len(points))
  copy(pointsCopy, points)
  sort.Slice(pointsCopy, func(i, j int) bool { return pointsCopy[i].LexicallyLessThan(*pointsCopy[j]) })

  graph := newAdjacencyGraph(pointsCopy)
  dacTriangulate(pointsCopy, &graph)
	edges := graph.toEdgeSlice()
	return edges
}

// Divide and Conquer Triangulation Method
// points must be sorted lexicographically and graph must contain all of the points, with no edges to any of them
func dacTriangulate(points []*Point, graph *adjacencyGraph) {
  if len(points) <= 3 {
    for i := 0; i < len(points); i++ {
      for j := i+1; j < len(points); j++ {
        graph.addEdge(points[i],points[j])
      }
    }
    return
  }

  left := points[:len(points)/2]
  right := points[len(points)/2:]

  dacTriangulate(left,graph)
  dacTriangulate(right,graph)
  dacTriangulateMerge(left,right,graph)
}

func dacTriangulateMerge(left, right []*Point, graph *adjacencyGraph) {
  utL, utR := uct(left, right, graph)
  L, R := lct(left, right, graph)
  graph.addEdge(L,R)
  for L!=utL || R!=utR {
    A := false
    B := false

    R1 := graph.pred(R,L)
    if (crossProduct(L,R,R1) > 0) {
      R2 := graph.pred(R,R1)
      for !qTest(R1,L,R,R2) {
        graph.deleteEdge(R,R1)
        R1 = R2
        R2 = graph.pred(R,R1)
      }
    } else {
      A = true
    }

    L1 := graph.succ(L,R)
    if (crossProduct(R,L,L1) < 0) {
      L2 := graph.succ(L,L1)
      for !qTest(L,R,L1,L2) {
        graph.deleteEdge(L,L1)
        L1 = L2
        L2 = graph.succ(L,L1)
      }
    } else {
      B = true
    }

    if A {
      L = L1
    } else if B {
      R = R1
    } else if qTest(L,R,R1,L1){
      R = R1
    } else {
      L = L1
    }

    graph.addEdge(L,R)
  }
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
func lct(left, right []*Point, graph *adjacencyGraph) (*Point, *Point) {
  // rightmost point in left set
  x := left[len(left)-1]
  // leftmost point in right set
  y := right[0]

  nextX := graph.hullPrevious(x)
  nextY := graph.hullNext(y)

  for lctFound := false; !lctFound; {
    if crossProduct(x,y,nextY) < 0 {
      y = nextY
      nextY = graph.hullNext(y)
    } else if crossProduct(x,y,nextX) < 0 {
      x = nextX
      nextX = graph.hullPrevious(x)
    } else {
      lctFound = true
    }
  }
  return x,y
}

func uct(left, right []*Point, graph *adjacencyGraph) (*Point, *Point) {
  // rightmost point in left set
  x := left[len(left)-1]
  // leftmost point in right set
  y := right[0]

  nextX := graph.hullNext(x)
  nextY := graph.hullPrevious(y)

  for lctFound := false; !lctFound; {
    if crossProduct(x,y,nextY) > 0 {
      y = nextY
      nextY = graph.hullPrevious(y)
    } else if crossProduct(x,y,nextX) > 0 {
      x = nextX
      nextX = graph.hullNext(x)
    } else {
      lctFound = true
    }
  }
  return x,y}

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

func (graph adjacencyGraph) succ(u, v *Point) *Point {
  return graph[u].successor(v)
}
func (graph adjacencyGraph) pred(u, v *Point) *Point {
  return graph[u].predecessor(v)
}
func (graph adjacencyGraph) hullNext(u *Point) *Point {
  return graph[u].firstNeighbor()
}
func (graph adjacencyGraph) hullPrevious(u *Point) *Point {
  return graph[u].lastNeighbor()
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

func (al adjacencyList) find(point *Point) *adjacencyListNode {
  for thisNeighbor := al.root.next; thisNeighbor != al.root; thisNeighbor = thisNeighbor.next {
    if thisNeighbor.point == point {
      return thisNeighbor
    }
  }
  return nil
}

func (al *adjacencyList) remove(point *Point) {
  node := al.find(point)
  if (node == nil) {
    panic("point not found in adjacency list")
  }
  node.previous.next = node.next
  node.next.previous = node.previous
}

func (al adjacencyList) neighbors() []*Point {
	neighbors := make([]*Point, 0)
	for thisNeighbor := al.root.next; thisNeighbor != al.root; thisNeighbor = thisNeighbor.next {
		neighbors = append(neighbors, thisNeighbor.point)
	}
	return neighbors
}

func (al adjacencyList) successor(point *Point) *Point{
  node := al.find(point)
  if (node == nil) {
    panic("point not found in adjacency list")
  }
  if (node.next != al.root) {
    return node.next.point
  } else {
    return node.next.next.point
  }
}

func (al adjacencyList) predecessor(point *Point) *Point{
  node := al.find(point)
  if (node == nil) {
    panic("point not found in adjacency list")
  }
  if (node.previous != al.root) {
    return node.previous.point
  } else {
    return node.previous.previous.point
  }
}

func (al adjacencyList) firstNeighbor() *Point{
  neighbor := al.root.next.point
  if neighbor != nil {
    return neighbor
  } else {
    return al.center
  }

}
func (al adjacencyList) lastNeighbor() *Point{
  neighbor := al.root.previous.point
  if neighbor != nil {
    return neighbor
  } else {
    return al.center
  }
}
