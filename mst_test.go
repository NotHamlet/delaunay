package delauney

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

const folderName = "testFiles"

func randomPoint(minX, maxX, minY, maxY float64) *Point {
	var xRange = maxX - minX
	var yRange = maxY - minY
	return &Point{
		X: minX + xRange*r.Float64(),
		Y: minY + yRange*r.Float64(),
	}
}

var r *rand.Rand

func TestMST99(t *testing.T) {

}

func TestMST1000(t *testing.T) {
	r = rand.New(rand.NewSource(80))

	N := 4000
	minX := float64(0)
	maxX := float64(300)
	minY := float64(0)
	maxY := float64(300)

	randomPoints := make([]*Point, N)
	for n := range randomPoints {
		randomPoints[n] = randomPoint(minX, maxX, minY, maxY)
	}

  mst := EuclideanMST(randomPoints)

	drawGraph(randomPoints, mst, minX, maxX, minY, maxY, "points1000.ps")
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
	f, err := os.Create(fmt.Sprintf("%s/%s", folderName, fileName))
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

  for _, point := range points {
    w.WriteString(fmt.Sprintf("%f %f 1 drawO\n", point.X, point.Y))
  }

  for _,edge := range edges {
    w.WriteString(fmt.Sprintf("%f %f %f %f drawLine\n", edge.a.X, edge.a.Y, edge.b.X, edge.b.Y))
  }

	w.Flush()
}
