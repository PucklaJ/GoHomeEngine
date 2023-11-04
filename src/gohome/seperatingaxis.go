package gohome

import (
	"github.com/PucklaJ/mathgl/mgl32"
)

// Returns wether two Polygons are intersecting
func AreIntersecting(p1, p2 *PolygonMath2D) bool {
	vertices1 := *p1
	vertices2 := *p2

	relativeVecs1, relativeVecs2 := make([]mgl32.Vec2, len(vertices1)), make([]mgl32.Vec2, len(vertices2))

	var j int
	for i := 0; i < len(vertices1); i++ {
		j = i + 1
		if j == len(vertices1) {
			j = 0
		}
		relativeVecs1[i] = vertices1[j].Sub(vertices1[i])
	}
	for i := 0; i < len(vertices2); i++ {
		j = i + 1
		if j == len(vertices2) {
			j = 0
		}
		relativeVecs2[i] = vertices2[j].Sub(vertices2[i])
	}

	edges1, edges2 := make([]mgl32.Vec2, len(vertices1)), make([]mgl32.Vec2, len(vertices2))
	for i := 0; i < len(vertices1); i++ {
		edges1[i][0], edges1[i][1] = -relativeVecs1[i][1], relativeVecs1[i][0]
	}
	for i := 0; i < len(vertices2); i++ {
		edges2[i][0], edges2[i][1] = -relativeVecs2[i][1], relativeVecs2[i][0]
	}

	for i := 0; i < len(vertices1); i++ {
		if !areOverlapping(edges1[i], vertices1, vertices2) {
			return false
		}
	}
	for i := 0; i < len(vertices2); i++ {
		if !areOverlapping(edges2[i], vertices1, vertices2) {
			return false
		}
	}

	return true
}

// Returns wether a polygon is intersecting with a point
func AreIntersectingPoint(p1 *PolygonMath2D, point mgl32.Vec2) bool {
	vertices1 := *p1
	vertices2 := []mgl32.Vec2{point}

	relativeVecs1 := make([]mgl32.Vec2, len(vertices1))
	var j int
	for i := 0; i < len(vertices1); i++ {
		j = i + 1
		if j == len(vertices1) {
			j = 0
		}
		relativeVecs1[i] = vertices1[j].Sub(vertices1[i])
	}
	edges1 := make([]mgl32.Vec2, len(vertices1))
	for i := 0; i < len(vertices1); i++ {
		edges1[i][0], edges1[i][1] = -relativeVecs1[i][1], relativeVecs1[i][0]
	}

	for i := 0; i < len(vertices1); i++ {
		if !areOverlapping(edges1[i], vertices1, vertices2) {
			return false
		}
	}

	return true
}

func areOverlapping(edge mgl32.Vec2, vertices1, vertices2 []mgl32.Vec2) bool {
	var mmDots1, mmDots2 [2]float32

	mmDots1 = getMaxMinDotProducts(edge, vertices1)
	mmDots2 = getMaxMinDotProducts(edge, vertices2)

	return mmDots1[0] < mmDots2[1] && mmDots1[1] > mmDots2[0]
}

func getMaxMinDotProducts(edge mgl32.Vec2, vertices []mgl32.Vec2) (vals [2]float32) {
	products := getDotProducts(edge, vertices)
	for i := 0; i < len(vertices); i++ {
		if i == 0 {
			vals[0] = products[i]
			vals[1] = products[i]
		} else {
			mgl32.SetMax(&vals[1], &products[i])
			mgl32.SetMin(&vals[0], &products[i])
		}
	}

	return
}

func getDotProducts(edge mgl32.Vec2, relativeVecs []mgl32.Vec2) (products []float32) {
	products = make([]float32, len(relativeVecs))
	for i := 0; i < len(relativeVecs); i++ {
		products[i] = edge.Dot(relativeVecs[i])
	}

	return
}
