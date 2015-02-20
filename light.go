package hdrt

import (
	"math"

	"github.com/gfrey/hdrt/vec"
)

type Light struct {
	Position  *vec.Vector
	Direction *vec.Vector
	Angle     float64
	Distance  float64
}

func (l *Light) InCone(pos, normal *vec.Vector) (float64, float64, *vec.Vector) {
	v := vec.Sub(pos, l.Position)
	if vec.Dot(v, normal) > 0 {
		return 0.0, 0.0, nil
	}
	d := v.Length()
	if d > l.Distance { // outside of light cone
		return 0.0, 0.0, nil
	}

	cosDelta := vec.Dot(v, l.Direction) / (d * l.Direction.Length())
	delta := math.Acos(cosDelta)
	if delta > (deg2rad(l.Angle) / 2.0) {
		return 0.0, 0.0, nil
	}

	return d, delta, v
}
