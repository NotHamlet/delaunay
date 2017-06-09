package delaunay

import "fmt"

import (
  "math"
)

type Point struct {
  X, Y float64
}

func (this Point) LexicallyLessThan(other Point) bool {
  if (this.X != other.X) {
    return (this.X < other.X)
  }
  return (this.Y < other.Y)
}


type Edge struct {
  a, b *Point
}

func (e Edge) length() float64 {
  return segmentLength(e.a, e.b)
}

// computes the z component of the cross product of the vectors pA->pB and pA->pC
func crossProduct(pA, pB, pC *Point) float64 {
  x1 := pB.X - pA.X
  y1 := pB.Y - pA.Y
  x2 := pC.X - pA.X
  y2 := pC.Y - pA.Y
  return (x1*y2 - x2*y1)
}

func dotProduct(pA, pB, pC *Point) float64 {
  x1 := pB.X - pA.X
  y1 := pB.Y - pA.Y
  x2 := pC.X - pA.X
  y2 := pC.Y - pA.Y
  return x1*x2 + y1*y2
}

func segmentLength(pA, pB *Point) float64 {
  dsquared := math.Pow(pA.X-pB.X, 2) + math.Pow(pA.Y-pB.Y, 2)
  return math.Sqrt(dsquared)
}

// checks if the line segment p1->p2 intersects p3->p4
func doLineSegmentsIntersect(p1, p2, p3, p4 *Point) bool {
  // if p1==p2 or p3==p4, panic
  if (*p1 == *p2 || *p3 == *p4) {
    panic("no line segment between equivalent points")
  }
  a := p1.X - p3.X
  b := p1.Y - p3.Y
  m1 := p1.X - p2.X
  m2 := p4.X - p3.X
  m3 := p1.Y - p2.Y
  m4 := p4.Y - p3.Y
  det := m1*m4 - m2*m3
  if (det == 0) {
    // the segments are parallel
    // FIXME
    return false
  }
  k1 := (m4*a - m2*b)/det
  k2 := (m1*b - m3*a)/det
  fmt.Sprintf("%f %f\n", k1, k2)
  return (k1>=0 && k1<=1 && k2>=0 && k2<=1)
}

func qTest(pA, pB, pC, pD *Point) bool {
  var m11, m12, m13 float64 = pA.X-pD.X, pA.Y-pD.Y, (pA.X*pA.X-pD.X*pD.X)+(pA.Y*pA.Y-pD.Y*pD.Y)
  var m21, m22, m23 float64 = pB.X-pD.X, pB.Y-pD.Y, (pB.X*pB.X-pD.X*pD.X)+(pB.Y*pB.Y-pD.Y*pD.Y)
  var m31, m32, m33 float64 = pC.X-pD.X, pC.Y-pD.Y, (pC.X*pC.X-pD.X*pD.X)+(pC.Y*pC.Y-pD.Y*pD.Y)

  det := m11*m22*m33 + m12*m23*m31 + m13*m21*m32 - m11*m23*m32 - m12*m21*m33 - m13*m22*m31
  return det <= 0
}
