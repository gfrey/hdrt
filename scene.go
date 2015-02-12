package hdrt

import (
	"encoding/json"
	"fmt"
	"github.com/dynport/dgtk/log"
	"image/color"
	"math"
)

type Scene struct {
	Objects []Object `json:"objects"`
	Lights  []*Light `json:"lights"`
}

func (sc *Scene) UnmarshalJSON(data []byte) error {
	rsc := &struct {
		Objects []*rawObject
		Lights  []*Light
	}{}

	err := json.Unmarshal(data, &rsc)
	if err != nil {
		return err
	}

	for i := range rsc.Objects {
		sc.Objects = append(sc.Objects, rsc.Objects[i].obj)
	}

	sc.Lights = rsc.Lights
	return nil
}

func (sc *Scene) Render(pos, dir *Vector) *color.RGBA {
	for i := range sc.Objects {
		o := sc.Objects[i]
		ipos := o.Intersect(pos, dir)
		if ipos != nil {
			return o.GetColor()
		}
	}
	return &color.RGBA{255, 0, 0, 255}
}

type rawObject struct {
	Type       string
	Position   *Vector
	Properties json.RawMessage
	obj        Object
}

func (robj *rawObject) UnmarshalJSON(data []byte) error {
	tobj := &struct {
		Type       string
		Position   *Vector
		Color      *color.RGBA
		Properties json.RawMessage
	}{}
	err := json.Unmarshal(data, &tobj)
	if err != nil {
		return err
	}

	switch tobj.Type {
	case "sphere":
		s := new(objSphere)
		s.BaseObject = &BaseObject{Position: tobj.Position, Color: tobj.Color}
		robj.obj = s
	case "box":
		b := new(objBox)
		b.BaseObject = &BaseObject{Position: tobj.Position, Color: tobj.Color}
		robj.obj = b
	default:
		return fmt.Errorf("type %q not supported", robj.Type)
	}

	return json.Unmarshal(tobj.Properties, &robj.obj)
}

type Object interface {
	Intersect(pos *Vector, dir *Vector) (intersection *Vector) // returns nil on no intersection
	GetColor() *color.RGBA
}

type BaseObject struct {
	Position *Vector
	Color    *color.RGBA
}

func (o *BaseObject) GetColor() *color.RGBA {
	if o.Color == nil {
		return &color.RGBA{0, 0, 255, 255}
	}
	return o.Color
}

type objSphere struct {
	*BaseObject
	Radius float64
}

func (o *objSphere) Intersect(p, d *Vector) *Vector {
	log.Info(">> Intersect(%s, %s)", p, d)

	c := o.Position
	vpc := VectorSub(c, p)
	vpcd := VectorDot(d, vpc)

	if vpcd < 0.0 {
		log.Log("SPHERE", "ddc <= 0 | %.2f", vpcd)
		// sphere is behind the viewplane
	} else {
		log.Info("ddc > 0 | %.2f", vpcd)
		puv := VectorProject(d, c)
		pc := VectorAdd(p, puv) // center of the sphere projected onto the ray
		log.Info("pc: %s", pc)
		log.Info("pc.DistanceTo(c)=%.2f o.Radius=%.2f", pc.DistanceTo(c), o.Radius)

		if pc.DistanceTo(c) > o.Radius {
			// no intersection
		} else {
			// pc is intersection in the middle
			return o.findFirstIntersectionPoint(vpc, pc, p, d)
		}
	}
	return nil
}

func (o *objSphere) findFirstIntersectionPoint(vpc *Vector, pc *Vector, p *Vector, d *Vector) *Vector {
	pcmcl := VectorSub(pc, o.Position).Length()
	dist := math.Sqrt(o.Radius*o.Radius - pcmcl*pcmcl)

	var di1 float64
	if vpc.Length() > o.Radius {
		// ray origin is outside sphere
		di1 = VectorSub(pc, p).Length() - dist
	} else {
		// ray origin is inside sphere
		di1 = VectorSub(pc, p).Length() + dist
	}

	return VectorScalarMultiply(VectorAdd(p, d), di1)
}

type objBox struct {
	*BaseObject
	Width, Height, Depth float64
}

func (o *objBox) Intersect(pos, dir *Vector) *Vector {
	w, h, d := o.Width/2.0, o.Height/2.0, o.Depth/2.0

	p0 := NewVector(o.Position[0]+w, o.Position[1]+h, o.Position[2]+d)
	p1 := NewVector(o.Position[0]-w, o.Position[1]+h, o.Position[2]+d)
	p2 := NewVector(o.Position[0]+w, o.Position[1]-h, o.Position[2]+d)
	p3 := NewVector(o.Position[0]-w, o.Position[1]-h, o.Position[2]+d)
	p4 := NewVector(o.Position[0]+w, o.Position[1]+h, o.Position[2]-d)
	p5 := NewVector(o.Position[0]-w, o.Position[1]+h, o.Position[2]-d)
	p6 := NewVector(o.Position[0]+w, o.Position[1]-h, o.Position[2]-d)
	p7 := NewVector(o.Position[0]-w, o.Position[1]-h, o.Position[2]-d)

	var cand *Vector
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

func intersectSquare(l0, l, p0, p1, p2 *Vector) *Vector {
	a := VectorSub(p1, p0)
	b := VectorSub(p2, p0)
	normal := VectorCross(a, b)
	normal.Normalize()
	cand := intersectPlane(l0, l, p0, normal)
	if cand != nil && pointInPlane(a, b, VectorSub(cand, p0)) {
		return cand
	}
	return nil
}

// point w in plane opened by u and v (from origin) (can be reused for triangles with r+t < 1.0)
func pointInPlane(u, v, w *Vector) bool {
	vCrossW := VectorCross(v, w)
	vCrossU := VectorCross(v, u)

	// Test sign of r
	if VectorDot(vCrossW, vCrossU) < 0.0 {
		return false
	}

	uCrossW := VectorCross(u, w)
	uCrossV := VectorCross(u, v)

	// Test sign of t
	if VectorDot(uCrossW, uCrossV) < 0.0 {
		return false
	}

	// At this point, we know that r and t and both > 0.
	// Therefore, as long as their sum is <= 1, each must be less <= 1
	denom := uCrossV.Length()
	r := vCrossW.Length() / denom
	t := uCrossW.Length() / denom

	return FloatLessThan(r, 1.0, epsilon) && FloatLessThan(t, 1.0, epsilon)
}

// l0 point on the ray
// l  direction of the ray
// p0 point on the plane
// n  normal of the plane
//
// d := (p0 - l0) * n / l * n
// if divisor is zero the line is parallel to the plane
// if divisor and divident are zero line is contained in plane
func intersectPlane(l0, l, p0, n *Vector) *Vector {
	divident := VectorDot(VectorSub(p0, l0), n)
	divisor := VectorDot(l, n)

	switch {
	case FloatEqual(divisor, 0.0, epsilon):
		if FloatEqual(divident, 0.0, epsilon) {
			return l0
		}
		return nil
	case FloatGreaterThan(divisor, 0.0, epsilon):
		return nil
	default:
		return VectorAdd(l0, VectorScalarMultiply(l, divident/divisor))
	}
}

type Light interface{}
