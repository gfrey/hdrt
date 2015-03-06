package obj

import (
	"fmt"

	"github.com/gfrey/hdrt/vec"
)

type Object interface {
	Intersect(pos, dir *vec.Vector) (float64, *vec.Vector) // returns nil on no intersection
	Normal(pos *vec.Vector) *vec.Vector
	Material(MaterialType) *Material
	Reflection() uint8
}

type BaseObject struct {
	Position *vec.Vector
	Mat      *material `json:"material"`
}

func (bo *BaseObject) Reflection() uint8 {
	if bo.Mat.Reflection < 2 {
		return 2
	}
	return bo.Mat.Reflection
}

func (bo *BaseObject) Material(typ MaterialType) *Material {
	switch typ {
	case MATERIAL_AMBIENT:
		return bo.Mat.Ambient
	case MATERIAL_DIFFUSE:
		return bo.Mat.Diffuse
	case MATERIAL_SPECULAR:
		return bo.Mat.Specular
	default:
		panic(fmt.Sprintf("material type %d not supported", typ))
	}
}
