package hdrt

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/gfrey/hdrt/vec"
)

type Scene struct {
	AmbientLight float64
	Objects      []Object `json:"objects"`
	Lights       []*Light `json:"lights"`
}

func (sc *Scene) UnmarshalJSON(data []byte) error {
	rsc := &struct {
		AmbientLight float64
		Objects      []*rawObject
		Lights       []*Light
	}{}

	err := json.Unmarshal(data, &rsc)
	if err != nil {
		return err
	}

	for i := range rsc.Objects {
		sc.Objects = append(sc.Objects, rsc.Objects[i].obj)
	}

	sc.AmbientLight = rsc.AmbientLight
	sc.Lights = rsc.Lights
	return nil
}

type intersection struct {
	d   float64
	obj Object
}

func (sc *Scene) Render(pos, dir *vec.Vector) *color.RGBA {
	var (
		cand     Object
		distance float64
		ipos     *vec.Vector
	)
	for i := range sc.Objects {
		o := sc.Objects[i]
		p := o.Intersect(pos, dir)
		if p != nil {
			d := (p.Data[0] - pos.Data[0]) / dir.Data[0]
			if cand == nil || d < distance {
				ipos = p
				cand = sc.Objects[i]
				distance = d
			}
		}
	}
	if cand != nil {
		return sc.ColorWithLights(cand, ipos)
	}
	return &color.RGBA{0, 0, 0, 0}
}

func (sc *Scene) ColorWithLights(obj Object, pos *vec.Vector) *color.RGBA {
	baseLight := sc.AmbientLight
	normal := obj.Normal(pos)
LIGHTSOURCES:
	for i := range sc.Lights {
		dist, delta, dir := sc.Lights[i].InCone(pos, normal)
		if dir != nil {
			lPos := sc.Lights[i].Position
			for j := range sc.Objects {
				if sc.Objects[j] == obj {
					continue
				}

				tmpPos := sc.Objects[j].Intersect(lPos, dir)
				if tmpPos == nil {
					continue
				}
				d := (tmpPos.Data[0] - lPos.Data[0]) / dir.Data[0]
				if d < dist { // in shadow
					continue LIGHTSOURCES
				}
			}
			// not in shadow
			maxAngle := deg2rad(sc.Lights[i].Angle) / 2.0
			baseLight += (1.0 - (dist / sc.Lights[i].Distance)) * (maxAngle - delta) / maxAngle
		}
	}

	if baseLight > 1.0 {
		baseLight = 1.0
	}

	c := obj.GetColor()
	return &color.RGBA{
		uint8(float64(c.R) * baseLight),
		uint8(float64(c.G) * baseLight),
		uint8(float64(c.B) * baseLight),
		c.A}
}

type rawObject struct {
	Type       string
	Position   *vec.Vector
	Properties json.RawMessage
	obj        Object
}

func (robj *rawObject) UnmarshalJSON(data []byte) error {
	tobj := &struct {
		Type       string
		Position   *vec.Vector
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
	Intersect(pos *vec.Vector, dir *vec.Vector) (intersection *vec.Vector) // returns nil on no intersection
	Normal(pos *vec.Vector) *vec.Vector
	GetColor() *color.RGBA
}

type BaseObject struct {
	Position *vec.Vector
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

func (o *objSphere) Normal(pos *vec.Vector) *vec.Vector {
	n := vec.Sub(pos, o.Position)
	d := n.Length()
	if !vec.FloatEqual(d, o.Radius, epsilon) {
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

type objBox struct {
	*BaseObject
	Width, Height, Depth float64
}

func (o *objBox) Normal(pos *vec.Vector) *vec.Vector {
	w, h, d := o.Width/2.0, o.Height/2.0, o.Depth/2.0

	switch {
	case vec.FloatEqual(o.Position.Data[0]+w, pos.Data[0], epsilon):
		return vec.New(1.0, 0.0, 0.0)
	case vec.FloatEqual(o.Position.Data[0]-w, pos.Data[0], epsilon):
		return vec.New(-1.0, 0.0, 0.0)
	case vec.FloatEqual(o.Position.Data[1]+h, pos.Data[1], epsilon):
		return vec.New(0.0, 1.0, 0.0)
	case vec.FloatEqual(o.Position.Data[1]-h, pos.Data[1], epsilon):
		return vec.New(0.0, -1.0, 0.0)
	case vec.FloatEqual(o.Position.Data[2]+d, pos.Data[2], epsilon):
		return vec.New(0.0, 0.0, 1.0)
	case vec.FloatEqual(o.Position.Data[2]-d, pos.Data[2], epsilon):
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

	return vec.FloatLessThan(r, 1.0, epsilon) && vec.FloatLessThan(t, 1.0, epsilon)
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
	case vec.FloatEqual(divisor, 0.0, epsilon):
		if vec.FloatEqual(divident, 0.0, epsilon) {
			return l0
		}
		return nil
	case vec.FloatGreaterThan(divisor, 0.0, epsilon):
		return nil
	default:
		return vec.Add(l0, vec.ScalarMultiply(l, divident/divisor))
	}
}
