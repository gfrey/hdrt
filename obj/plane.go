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

func (o *objPlane) Intersect(p, r *vec.Vector) (float64, *vec.Vector) {
	if o.d == nil {
		o.Up.Normalize()
		o.d = new(float64)
		*o.d = vec.Dot(o.Up, o.Position)
	}
	nr := vec.Dot(o.Up, r)
	if mat.FloatGreaterThanEqual(nr, 0.0) { // ray parallel to plane
		return 0,nil
	}
	lambda := (*o.d - vec.Dot(o.Up, p)) / nr
	if mat.FloatLessThan(lambda, 0) {
		return 0, nil
	}
	return lambda, vec.Add(p, vec.ScalarMultiply(r, lambda))
}