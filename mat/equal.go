package mat

import "math"

var ε float64 = 1e-12

func FloatEqual(a, b float64) bool {
	absA, absB := math.Abs(a), math.Abs(b)
	diff := math.Abs(a - b)

	switch {
	case a == b:
		return true
	case a == 0, b == 0, diff < math.Nextafter(1.0, 2.0)-1.0:
		return diff < ε
	default:
		return diff/(absA+absB) < ε
	}
}

func FloatLessThan(a, b float64) bool {
	return a < b
}

func FloatLessThanEqual(a, b float64) bool {
	return a < b || FloatEqual(a, b)
}

func FloatGreaterThan(a, b float64) bool {
	return a > b
}

func FloatGreaterThanEqual(a, b float64) bool {
	return a > b || FloatEqual(a, b)
}
