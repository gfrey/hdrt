package obj

import (
	"image/color"

	"github.com/gfrey/hdrt/vec"
)

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
