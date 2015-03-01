package obj

import (
	"github.com/gfrey/hdrt/mat"
	"github.com/gfrey/hdrt/vec"
)

type objPlane struct {
	*BaseObject
	Up *vec.Vector
	d  *float64
}

func (o *objPlane) Normal(pos *vec.Vector) *vec.Vector {
	return o.Up
}

func (o *objPlane) Intersect(p, r *vec.Vector) *vec.Vector {
	if o.d == nil {
		o.Up.Normalize()
		o.d = new(float64)
		*o.d = vec.Dot(o.Up, o.Position)
	}
	nr := vec.Dot(o.Up, r)
	if mat.FloatGreaterThanEqual(nr, 0.0) { // ray parallel to plane
		return nil
	}
	lambda := (*o.d - vec.Dot(o.Up, p)) / nr
	if mat.FloatLessThan(lambda, 0) {
		return nil
	}
	return vec.Add(p, vec.ScalarMultiply(r, lambda))
}