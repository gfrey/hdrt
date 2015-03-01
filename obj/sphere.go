package obj

import (
	"log"
	"math"

	"github.com/gfrey/hdrt/mat"
	"github.com/gfrey/hdrt/vec"
)

type objSphere struct {
	*BaseObject
	Radius float64
}

func (o *objSphere) Normal(pos *vec.Vector) *vec.Vector {
	n := vec.Sub(pos, o.Position)
	d := n.Length()
	if !mat.FloatEqual(d, o.Radius) {
		log.Printf("radius %.2f should equal  %.2f", o.Radius, d)
	}
	n.Normalize()
	return n
}

func (o *objSphere) Intersect(p, d *vec.Vector) *vec.Vector {
	deltaP := vec.Sub(o.Position, p)
	dDeltaP := vec.Dot(d, deltaP)

	if mat.FloatLessThan(dDeltaP, 0) {
		return nil
	}

	lambda := vec.Sub(deltaP, vec.ScalarMultiply(d, dDeltaP)).Length()
	if mat.FloatGreaterThan(lambda, o.Radius) {
		return nil
	}

	sqrt := math.Sqrt(o.Radius*o.Radius - lambda*lambda)
	s1 := dDeltaP - sqrt
	s2 := dDeltaP + sqrt

	switch {
	case mat.FloatLessThan(s1, 0.0), mat.FloatLessThan(s2, 0.0):
		return nil // don't show anything if inside the sphere
	case mat.FloatLessThan(s1, s2):
		return vec.Add(p, vec.ScalarMultiply(d, s1))
	default:
		return vec.Add(p, vec.ScalarMultiply(d, s2))
	}
}
