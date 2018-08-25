package geom

import (
	"fmt"
	"math"
)

const tolerance = 0.000001

// Float64 compares two floats to see if they are within the given tolerance.
func cmpFloat(f1, f2 float64) bool {
	if math.IsInf(f1, 1) {
		return math.IsInf(f2, 1)
	}
	if math.IsInf(f2, 1) {
		return math.IsInf(f1, 1)
	}
	if math.IsInf(f1, -1) {
		return math.IsInf(f2, -1)
	}
	if math.IsInf(f2, -1) {
		return math.IsInf(f1, -1)
	}
	return math.Abs(f1-f2) < tolerance
}
func pointEqual(p1, p2 [2]float64) bool { return cmpFloat(p1[0], p2[0]) && cmpFloat(p1[1], p2[1]) }

type Triangle [3][2]float64

func (t Triangle) Center() (pt [2]float64) {
	pt[0] = (t[0][0] + t[1][0] + t[2][0]) / 3
	pt[1] = (t[0][1] + t[1][1] + t[2][1]) / 3
	return pt
}

func (t Triangle) ThirdPoint(p1, p2 [2]float64) [2]float64 {
	switch {
	case (pointEqual(t[0], p1) && pointEqual(t[1], p2)) ||
		(pointEqual(t[1], p1) && pointEqual(t[0], p2)):
		return t[2]
	case (pointEqual(t[0], p1) && pointEqual(t[2], p2)) ||
		(pointEqual(t[2], p1) && pointEqual(t[0], p2)):
		return t[1]
	default:
		return t[0]
	}
}

func NewTriangleFromPolygon(py [][][2]float64) Triangle {
	// Assume we are getting triangles from the function.
	if debug && len(py) != 1 {
		panic(fmt.Sprintf("Step   3 : assumption invalid for triangle. %v", py))
	}
	if debug && len(py[0]) < 3 {
		panic(fmt.Sprintf("Step   3 : assumption invalid for triangle. %v", py))
	}
	t := Triangle{py[0][0], py[0][1], py[0][2]}
	return t
}
