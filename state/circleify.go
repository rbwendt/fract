package state

import "math"

func circleify(x float64) float64 {
	return math.Sqrt(1 - math.Pow(x, 2))
}
