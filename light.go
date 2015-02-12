package hdrt

import "math"

type Light struct {
	Position  *Vector
	Direction *Vector
	Angle     float64
	Distance  float64
}

func (l *Light) InCone(pos *Vector) (float64, float64, *Vector) {
	v := VectorSub(pos, l.Position)
	d := v.Length()
	if d > l.Distance { // outside of light cone
		return 0.0, 0.0, nil
	}

	cosDelta := VectorDot(v, l.Direction) / (d * l.Direction.Length())
	delta := math.Acos(cosDelta)
	if delta > (deg2rad(l.Angle) / 2.0) {
		return 0.0, 0.0, nil
	}

	return d, delta, v
}
