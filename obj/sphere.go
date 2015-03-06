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
		log.Printf("radius %.2f, not %.2f", o.Radius, d)
	}
	n.Normalize()
	return n
}

// s = rayD*deltaP+- sqrt(r^2-|deltaP-rayD*deltaP*rayD|^2)
func (o *objSphere) Intersect(rayP, rayD *vec.Vector) (float64, *vec.Vector) {
	deltaP := vec.Sub(o.Position, rayP)
	dDeltaP := vec.Dot(rayD, deltaP)
	if mat.FloatLessThan(dDeltaP, 0) {
		return 0, nil
	}

	lambda := vec.Sub(deltaP, vec.ScalarMultiply(rayD, dDeltaP)).Length()
	lSqr, rSqr := lambda*lambda, o.Radius*o.Radius
	if mat.FloatLessThanEqual(rSqr, lSqr) {
		return 0, nil
	}

	sqrt := math.Sqrt(rSqr - lSqr)
	s1 := dDeltaP - sqrt
	s2 := dDeltaP + sqrt

	switch {
	case mat.FloatLessThan(s1, 0), mat.FloatLessThan(s2, 0):
		return 0, nil // don't show anything if inside the sphere
	case mat.FloatLessThan(s1, s2):
		return s1, vec.Add(rayP, vec.ScalarMultiply(rayD, s1))
	default:
		return s2, vec.Add(rayP, vec.ScalarMultiply(rayD, s2))
	}
}
