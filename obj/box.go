package obj

import (
	"github.com/gfrey/hdrt/mat"
	"github.com/gfrey/hdrt/vec"
)

type objBox struct {
	*BaseObject
	Width, Height, Depth float64
}

func (o *objBox) Normal(pos *vec.Vector) *vec.Vector {
	w, h, d := o.Width/2.0, o.Height/2.0, o.Depth/2.0

	switch {
	case mat.FloatEqual(o.Position.Data[0]+w, pos.Data[0]):
		return vec.New(1.0, 0.0, 0.0)
	case mat.FloatEqual(o.Position.Data[0]-w, pos.Data[0]):
		return vec.New(-1.0, 0.0, 0.0)
	case mat.FloatEqual(o.Position.Data[1]+h, pos.Data[1]):
		return vec.New(0.0, 1.0, 0.0)
	case mat.FloatEqual(o.Position.Data[1]-h, pos.Data[1]):
		return vec.New(0.0, -1.0, 0.0)
	case mat.FloatEqual(o.Position.Data[2]+d, pos.Data[2]):
		return vec.New(0.0, 0.0, 1.0)
	case mat.FloatEqual(o.Position.Data[2]-d, pos.Data[2]):
		return vec.New(0.0, 0.0, -1.0)
	}
	panic("don't know how to compute a normal")
}

func (o *objBox) Intersect(pos, dir *vec.Vector) *vec.Vector {
	w, h, d := o.Width/2.0, o.Height/2.0, o.Depth/2.0
	x, y, z := o.Position.Data[0], o.Position.Data[1], o.Position.Data[2]

	p0 := vec.New(x+w, y+h, z+d)
	p1 := vec.New(x-w, y+h, z+d)
	p2 := vec.New(x+w, y-h, z+d)
	p3 := vec.New(x-w, y-h, z+d)
	p4 := vec.New(x+w, y+h, z-d)
	p5 := vec.New(x-w, y+h, z-d)
	p6 := vec.New(x+w, y-h, z-d)
	p7 := vec.New(x-w, y-h, z-d)

	var cand *vec.Vector
	cand = intersectSquare(pos, dir, p0, p4, p1)
	if cand != nil {
		return cand
	}

	cand = intersectSquare(pos, dir, p0, p1, p2)
	if cand != nil {
		return cand
	}

	cand = intersectSquare(pos, dir, p0, p2, p4)
	if cand != nil {
		return cand
	}

	cand = intersectSquare(pos, dir, p7, p5, p6)
	if cand != nil {
		return cand
	}

	cand = intersectSquare(pos, dir, p7, p6, p3)
	if cand != nil {
		return cand
	}

	cand = intersectSquare(pos, dir, p7, p3, p5)
	if cand != nil {
		return cand
	}

	return nil
}

func intersectSquare(l0, l, p0, p1, p2 *vec.Vector) *vec.Vector {
	a := vec.Sub(p1, p0)
	b := vec.Sub(p2, p0)
	normal := vec.Cross(a, b)
	normal.Normalize()
	cand := intersectPlane(l0, l, p0, normal)
	if cand != nil && pointInPlane(a, b, vec.Sub(cand, p0)) {
		return cand
	}
	return nil
}

// point w in plane opened by u and v (from origin) (can be reused for triangles with r+t < 1.0)
func pointInPlane(u, v, w *vec.Vector) bool {
	vCrossW := vec.Cross(v, w)
	vCrossU := vec.Cross(v, u)

	// Test sign of r
	if vec.Dot(vCrossW, vCrossU) < 0.0 {
		return false
	}

	uCrossW := vec.Cross(u, w)
	uCrossV := vec.Cross(u, v)

	// Test sign of t
	if vec.Dot(uCrossW, uCrossV) < 0.0 {
		return false
	}

	// At this point, we know that r and t and both > 0.
	// Therefore, as long as their sum is <= 1, each must be less <= 1
	denom := uCrossV.Length()
	r := vCrossW.Length() / denom
	t := uCrossW.Length() / denom

	return mat.FloatLessThan(r, 1.0) && mat.FloatLessThan(t, 1.0)
}

// l0 point on the ray
// l  direction of the ray
// p0 point on the plane
// n  normal of the plane
//
// d := (p0 - l0) * n / l * n
// if divisor is zero the line is parallel to the plane
// if divisor and divident are zero line is contained in plane
func intersectPlane(l0, l, p0, n *vec.Vector) *vec.Vector {
	divident := vec.Dot(vec.Sub(p0, l0), n)
	divisor := vec.Dot(l, n)

	switch {
	case mat.FloatEqual(divisor, 0.0):
		if mat.FloatEqual(divident, 0.0) {
			return l0
		}
		return nil
	case mat.FloatGreaterThan(divisor, 0.0), mat.FloatGreaterThan(divident, 0.0):
		return nil
	default:
		return vec.Add(l0, vec.ScalarMultiply(l, divident/divisor))
	}
}
