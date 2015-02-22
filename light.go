package hdrt

import (
	"math"

	"github.com/gfrey/hdrt/obj"
	"github.com/gfrey/hdrt/vec"
)

type Light struct {
	Position  *vec.Vector
	Direction *vec.Vector
	Angle     float64
	Diffuse   *obj.Material
	Specular  *obj.Material
}

// Returns the emitted diffuse and specular light as color value for the given direction.
func (l *Light) InCone(pos, normal *vec.Vector) (float64, *obj.Material, *obj.Material) {
	v := vec.Sub(pos, l.Position)
	if vec.Dot(v, normal) > 0 {
		return 0.0, nil, nil
	}

	ang := (deg2rad(l.Angle) / 2.0)
	cosDelta := vec.Dot(v, l.Direction) / (v.Length() * l.Direction.Length())
	delta := math.Acos(cosDelta) // outside of cone?
	if delta > ang {
		return 0.0, nil, nil
	}

	// TODO: currently the light source will emit one uniform color; that
	//       should degrade from center to peripherie.
	return 1.0 - (delta / ang), l.Diffuse, l.Specular
}
