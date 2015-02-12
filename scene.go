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
	return nil
}

type Light interface{}
