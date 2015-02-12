package hdrt

import "math"

type Camera struct {
	Position  *Vector
	Direction *Vector
	Up        *Vector
	FOV       float64 // angle in degree
}

func (c *Camera) Init() error {
	c.Direction.Normalize()
	c.Up.Normalize()
	return nil
}

type Viewplane struct {
	Distance   float64
	ResX, ResY int

	pos, a, b *Vector
}

func (v *Viewplane) Init(c *Camera) error {
	v.span(c)
	return nil
}

func (v *Viewplane) span(c *Camera) {
	vpCenter := VectorAdd(c.Position, VectorScalarMultiply(c.Direction, v.Distance))

	aspectRatio := float64(v.ResX) / float64(v.ResY)
	alpha := deg2rad(c.FOV)
	beta := alpha / aspectRatio

	a := v.Distance * math.Tan(alpha/2.0)
	b := v.Distance * math.Tan(beta/2.0)

	vpTop := c.Up
	vpSide := VectorCross(c.Direction, c.Up)

	v.a = VectorScalarMultiply(vpSide, a)
	v.b = VectorScalarMultiply(vpTop, -b)
	v.pos = VectorAdd(VectorAdd(vpCenter, VectorScalarMultiply(v.a, -0.5)), VectorScalarMultiply(v.b, -0.5))

}
