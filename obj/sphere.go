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
	c := o.Position
	vpc := vec.Sub(c, p)

	if vec.Dot(d, vpc) < 0.0 {
		// sphere is behind the viewplane
		return nil
	}

	puv := vec.Project(vpc, d) // vpc on d
	pc := vec.Add(p, puv)      // center of the sphere projected onto the ray

	if pc.DistanceTo(c) > o.Radius {
		return nil
	}

	pcmcl := vec.Sub(pc, c).Length()
	dist := math.Sqrt(o.Radius*o.Radius - pcmcl*pcmcl)

	var di1 float64
	if vpc.Length() > o.Radius {
		// ray origin is outside sphere
		di1 = vec.Sub(pc, p).Length() - dist
	} else {
		// ray origin is inside sphere
		di1 = vec.Sub(pc, p).Length() + dist
	}

	return vec.Add(p, vec.ScalarMultiply(d, di1))
}
