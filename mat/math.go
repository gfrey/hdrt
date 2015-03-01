package mat

import "math"

func Deg2Rad(x float64) float64 {
	return x * math.Pi / 180.0
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
