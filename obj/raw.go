package obj

import (
	"encoding/json"
	"fmt"

	"github.com/gfrey/hdrt/vec"
)

type Raw struct {
	Type       string
	Position   *vec.Vector
	Properties json.RawMessage
	obj        Object
}

func (robj *Raw) Object() Object {
	return robj.obj
}

func (robj *Raw) UnmarshalJSON(data []byte) error {
	tobj := &struct {
		Type       string
		Position   *vec.Vector
		Material   *material
		Properties json.RawMessage
	}{}
	err := json.Unmarshal(data, &tobj)
	if err != nil {
		return err
	}

	switch tobj.Type {
	case "sphere":
		s := new(objSphere)
		s.BaseObject = &BaseObject{Position: tobj.Position, Mat: tobj.Material}
		robj.obj = s
	case "box":
		b := new(objBox)
		b.BaseObject = &BaseObject{Position: tobj.Position, Mat: tobj.Material}
		robj.obj = b
	case "plane":
		b := new(objPlane)
		b.BaseObject = &BaseObject{Position: tobj.Position, Mat: tobj.Material}
		robj.obj = b
	default:
		return fmt.Errorf("type %q not supported", robj.Type)
	}

	return json.Unmarshal(tobj.Properties, &robj.obj)
}
