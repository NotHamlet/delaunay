package delauney

import (
	"math/rand"
	"testing"
  "fmt"
)

func randomPoint(maxX, maxY float64) *Point {
	return &Point{
		X: maxX * r.Float64(),
		Y: maxY * r.Float64(),
	}
}

var r *rand.Rand

func TestPass(t *testing.T) {
	r = rand.New(rand.NewSource(105))
	N := 100
	randomPoints := make([]*Point, N)
	for n := range randomPoints {
		randomPoints[n] = randomPoint(100, 100)
		// fmt.Println(randomPoints[n])
	}

	newMst := EuclideanMST(randomPoints)
  if (len(newMst) != N-1) {
    t.Fail()
  }
}

func TestMSTEasy(t *testing.T) {
	var points []*Point
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			points = append(points,
        &Point{X: float64(5*i + 5), Y: float64(5*j + 5)})
		}
	}

	newMst := EuclideanMST(points)
  fmt.Println(len(newMst))
}
