package delaunay

import (
	"bufio"
	"fmt"
  "math"
	"math/rand"
	"os"
	"sort"
	"testing"
)

const folderName = "testFiles"
const EPS = 0.00001

func makePoints(minX, maxX, minY, maxY float64, seed int64, N int) []*Point {
	r := rand.New(rand.NewSource(seed))
	xRange := maxX - minX
	yRange := maxY - minY

	points := make([]*Point, N)

	for i := 0; i < N; i++ {
		points[i] = &Point{
			X: minX + xRange*r.Float64(),
			Y: minY + yRange*r.Float64(),
		}
	}

	return points
}

func TestPointSort(t *testing.T) {
	// fileName := "sortedPoints"
	N := 5000
	var seed int64 = 31
	minX := float64(-100)
	maxX := float64(100)
	minY := float64(-100)
	maxY := float64(100)

	randomPoints := makePoints(minX, maxX, minY, maxY, seed, N)

	sort.Slice(randomPoints, func(i, j int) bool { return randomPoints[i].LexicallyLessThan(*randomPoints[j]) })

	isSorted := sort.SliceIsSorted(randomPoints,
		func(i, j int) bool {
			return randomPoints[i].X < randomPoints[j].X ||
				(randomPoints[i].X == randomPoints[j].X && randomPoints[i].Y < randomPoints[j].Y)
		})

	if !isSorted {
		t.Fail()
	}

	// drawGraph(randomPoints, nil, minX, maxX, minY, maxY, fileName)
}

func TestTriangulate100(t *testing.T) {
	fileName := "tri100"
	N := 100
	var seed int64 = 105060
	minX := float64(0)
	maxX := float64(256)
	minY := float64(0)
	maxY := float64(160)

	randomPoints := makePoints(minX, maxX, minY, maxY, seed, N)
	edges := Triangulate(randomPoints)
	drawGraph(randomPoints, edges, minX, maxX, minY, maxY, fileName)
}

func TestMST500(t *testing.T) {
	fileName := "mst500"
	N := 500
	var seed int64 = 80
	minX := float64(0)
	maxX := float64(300)
	minY := float64(0)
	maxY := float64(300)

	randomPoints := makePoints(minX, maxX, minY, maxY, seed, N)
	mst := EuclideanMST(randomPoints)
  correctLength := 4473.850825
  actualLength := totalWeight(mst)
  if math.Abs(correctLength - actualLength) >= EPS {
    t.Fail()
  }
	drawGraph(randomPoints, mst, minX, maxX, minY, maxY, fileName)
}

func TestAdjacencyListInsert(t *testing.T) {
  center := &Point{-1,-1}
  al := newAdjacencyList(center)
  points := []*Point{
    &Point{4,0},
    &Point{2,2},
    &Point{1,10},
    &Point{-4,200},
    &Point{-3,-1},
  }

  //randomize our list of points before adding them to AL
  pointsCopy := make([]*Point, len(points))
  copy(pointsCopy, points)
  r := rand.New(rand.NewSource(5918))
  for i,_ := range pointsCopy {
    j := r.Intn(len(pointsCopy)-i)+i
    pointsCopy[i],pointsCopy[j] = pointsCopy[j],pointsCopy[i]
  }

  for _,point := range pointsCopy {
    al.insert(point)
  }

  expected := points
  actual := al.neighbors()
  if (len(actual) != len(expected)) {
    t.Fail()
  }
  for i,point := range actual {
    if point != expected[i] {
      t.Fail();
    }
  }

}

func totalWeight(edges []*Edge) float64 {
	var result float64
	for _, edge := range edges {
		result += edge.length()
	}
	return result
}

func drawGraph(points []*Point, edges []*Edge, minX, maxX, minY, maxY float64, fileName string) {
	os.Mkdir(folderName, 0777)
	f, err := os.Create(fmt.Sprintf("%s/%s.ps", folderName, fileName))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	w.WriteString("%!PS-Adobe-2.0 EPSF-1.2\n")
	w.WriteString(fmt.Sprintf("%%%%BoundingBox: %d %d %d %d\n", int(minX), int(minY), int(maxX), int(maxY)))
	w.WriteString("%%EndComments\n")
	w.WriteString("/drawO {newpath 0 360 arc stroke} def\n")
	w.WriteString("/drawLine {newpath moveto lineto stroke} def\n")

	w.WriteString(fmt.Sprintf("%f setgray\n", 0.5))
	for _, edge := range edges {
		w.WriteString(fmt.Sprintf("%f %f %f %f drawLine\n", edge.a.X, edge.a.Y, edge.b.X, edge.b.Y))
	}

  w.WriteString(fmt.Sprintf("%f setgray\n", 0.0))
	for _, point := range points {
		// w.WriteString(fmt.Sprintf("%f setgray\n", float64(n)/float64(len(points))))
		w.WriteString(fmt.Sprintf("%f %f 1 drawO\n", point.X, point.Y))
	}

	w.Flush()
}
