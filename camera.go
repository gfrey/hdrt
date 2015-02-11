package hdrt

type Camera struct {
	Position  Vector
	Direction Vector
	Up        Vector
}

type Viewplane struct {
	Distance float64
	Width    int
	Height   int
}
