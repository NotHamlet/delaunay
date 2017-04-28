package delauney

import (
  "math"
)

type Point struct {
  X, Y float64
}

type Edge struct {
  a, b *Point
}

func (e Edge) length() float64 {
  dsquared := math.Pow(e.a.X-e.b.X, 2) + math.Pow(e.a.Y-e.b.Y, 2)
  return math.Sqrt(dsquared)
}
