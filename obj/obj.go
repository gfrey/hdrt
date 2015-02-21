package obj

import (
	"fmt"

	"github.com/gfrey/hdrt/vec"
)

type Object interface {
	Intersect(pos *vec.Vector, dir *vec.Vector) (intersection *vec.Vector) // returns nil on no intersection
	Normal(pos *vec.Vector) *vec.Vector
	Material(MaterialType) *MaterialC
	Reflection() float64
}

type BaseObject struct {
	Position *vec.Vector
	mat      *Material `json:"material"`
}

func (bo *BaseObject) Reflection() float64 {
	if bo.mat.Reflection == 0 {
		return 1.5
	}
	return float64(bo.mat.Reflection)
}

func (bo *BaseObject) Material(typ MaterialType) *MaterialC {
	switch typ {
	case MATERIAL_AMBIENT:
		return bo.mat.Ambient
	case MATERIAL_DIFFUSE:
		return bo.mat.Diffuse
	case MATERIAL_SPECULAR:
		return bo.mat.Specular
	default:
		panic(fmt.Sprintf("material type %d not supported", typ))
	}
}

type MaterialC [3]uint

type MaterialType int

const (
	MATERIAL_AMBIENT MaterialType = iota
	MATERIAL_DIFFUSE
	MATERIAL_SPECULAR
)

// I need to specify which wavelength of light are reflected and which are held back by the material.
type Material struct {
	Reflection uint8
	Ambient    *MaterialC
	Diffuse    *MaterialC
	Specular   *MaterialC
}
