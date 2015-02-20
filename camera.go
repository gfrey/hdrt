package hdrt

import (
	"math"

	"github.com/gfrey/hdrt/vec"
)

type Camera struct {
	Position  *vec.Vector
	Direction *vec.Vector
	Up        *vec.Vector
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

	pos, a, b *vec.Vector
}

func (v *Viewplane) Init(c *Camera) error {
	v.span(c)
	return nil
}

func (v *Viewplane) span(c *Camera) {
	vpCenter := vec.VectorAdd(c.Position, vec.VectorScalarMultiply(c.Direction, v.Distance))

	aspectRatio := float64(v.ResX) / float64(v.ResY)
	alpha := deg2rad(c.FOV)
	beta := alpha / aspectRatio

	a := v.Distance * math.Tan(alpha/2.0)
	b := v.Distance * math.Tan(beta/2.0)

	vpTop := c.Up
	vpSide := vec.VectorCross(c.Direction, c.Up)

	v.a = vec.VectorScalarMultiply(vpSide, a)
	v.b = vec.VectorScalarMultiply(vpTop, -b)
	v.pos = vec.VectorAdd(vec.VectorAdd(vpCenter, vec.VectorScalarMultiply(v.a, -0.5)), vec.VectorScalarMultiply(v.b, -0.5))
}

func (v *Viewplane) PixelPosition(x, y int) *vec.Vector {
	return vec.VectorAdd(vec.VectorAdd(v.pos, vec.VectorScalarMultiply(v.a, float64(x)/float64(v.ResX))), vec.VectorScalarMultiply(v.b, float64(y)/float64(v.ResY)))
}
