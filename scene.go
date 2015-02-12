package hdrt

import (
	"encoding/json"
	"fmt"
	"github.com/dynport/dgtk/log"
	"image/color"
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
		ipos := sc.Objects[i].Intersect(pos, dir)
		if ipos != nil {
			return &color.RGBA{0, 0, 255, 255}
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
		Properties json.RawMessage
	}{}
	err := json.Unmarshal(data, &tobj)
	if err != nil {
		return err
	}

	switch tobj.Type {
	case "sphere":
		s := new(objSphere)
		s.BaseObject = &BaseObject{Position: tobj.Position}
		robj.obj = s
	case "box":
		b := new(objBox)
		b.BaseObject = &BaseObject{Position: tobj.Position}
		robj.obj = b
	default:
		return fmt.Errorf("type %q not supported", robj.Type)
	}

	return json.Unmarshal(tobj.Properties, &robj.obj)
}

type Object interface {
	Intersect(pos *Vector, dir *Vector) (intersction *Vector) // returns nil on no intersection
}

type BaseObject struct {
	Position *Vector
}

type objSphere struct {
	*BaseObject
	Radius float64
}

func (o *objSphere) Intersect(p, d *Vector) *Vector {
	c := o.Position
	dc := VectorSub(c, p)
	ddc := VectorDot(d, dc)

	if ddc > 0 {
		log.Debug("ddc > 0 | %.2f", ddc)
		puv := VectorProject(d, c)
		pc := VectorAdd(p, puv) // center of the sphere projected onto the ray
		log.Debug("pc: %s", pc)

		log.Debug("pc.DistanceTo(c)=%.2f o.Radius=%.2f", pc.DistanceTo(c), o.Radius)
		if pc.DistanceTo(c) <= o.Radius {
			// pc is intersection in the middle
			// TODO return o.findFirstIntersectionPoint()
			return pc
		}
	} else {
		log.Log("SPHERE", "ddc <= 0 | %.2f", ddc)
		// sphere is behind the viewplane
	}

	return nil
}

type objBox struct {
	*BaseObject
	Width, Height, Depth float64
}

func (o *objBox) Intersect(pos, dir *Vector) *Vector {
	return nil
}

type Light interface{}
