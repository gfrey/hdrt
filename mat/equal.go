package mat

const Epsilon = 0.0001

func FloatEqual(a, b float64) bool {
	return FloatEqualE(a, b, Epsilon)
}

func FloatLessThan(a, b float64) bool {
	return FloatLessThanE(a, b, Epsilon)
}

func FloatLessThanEqual(a, b float64) bool {
	return FloatLessThanEqualE(a, b, Epsilon)
}

func FloatGreaterThan(a, b float64) bool {
	return FloatGreaterThanE(a, b, Epsilon)
}

func FloatGreaterThanEqual(a, b float64) bool {
	return FloatGreaterThanEqualE(a, b, Epsilon)
}

func FloatEqualE(a, b, ε float64) bool {
	return a-ε < b && b < a+ε
}

func FloatLessThanE(a, b, ε float64) bool {
	return a < b+ε
}

func FloatLessThanEqualE(a, b, ε float64) bool {
	return a <= b+ε
}

func FloatGreaterThanE(a, b, ε float64) bool {
	return a > b-ε
}

func FloatGreaterThanEqualE(a, b, ε float64) bool {
	return a >= b-ε
}
