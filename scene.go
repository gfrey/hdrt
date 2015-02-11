package hdrt

import (
	"encoding/json"
	"fmt"
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

type rawObject struct {
	Type       string
	Position   Vector
	Properties json.RawMessage
	obj        Object
}

func (robj *rawObject) UnmarshalJSON(data []byte) error {
	tobj := &struct {
		Type       string
		Position   Vector
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
}

type BaseObject struct {
	Position Vector
}

type objSphere struct {
	*BaseObject
	Radius float64
}

type objBox struct {
	*BaseObject
	Width, Height, Depth float64
}

type Light interface{}
