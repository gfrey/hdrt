package obj

import (
	"fmt"

	"github.com/gfrey/hdrt/vec"
)

type Object interface {
	Intersect(pos *vec.Vector, dir *vec.Vector) (intersection *vec.Vector) // returns nil on no intersection
	Normal(pos *vec.Vector) *vec.Vector
	Material(MaterialType) *Material
	Reflection() float64
}

type BaseObject struct {
	Position *vec.Vector
	Mat      *material `json:"material"`
}

func (bo *BaseObject) Reflection() float64 {
	if bo.Mat.Reflection == 0 {
		return 1.5
	}
	return float64(bo.Mat.Reflection)
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
