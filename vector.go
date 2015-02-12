package hdrt

import (
	"encoding/json"
	"math"
)

type Vector [3]float64

func NewVector(x, y, z float64) *Vector {
	v := new(Vector)
	v[0], v[1], v[2] = x, y, z
	return v
}

func (v *Vector) UnmarshalJSON(data []byte) error {
	tv := &struct {
		X float64
		Y float64
		Z float64
	}{}
	err := json.Unmarshal(data, &tv)
	if err != nil {
		return err
	}

	v[0] = tv.X
	v[1] = tv.Y
	v[2] = tv.Z

	return nil
}

func VectorScalarMultiply(v *Vector, a float64) *Vector {
	return NewVector(v[0]*a, v[1]*a, v[2]*a)
}

func (v *Vector) ScalarMultiply(a float64) {
	v[0] *= a
	v[1] *= a
	v[2] *= a
}

func VectorAdd(a, b *Vector) *Vector {
	return NewVector(a[0]+b[0], a[1]+b[1], a[2]+b[2])
}

func (v *Vector) Add(other *Vector) {
	v[0] += other[0]
	v[1] += other[1]
	v[2] += other[2]
}

func VectorCross(a, b *Vector) *Vector {
	v := new(Vector)

	v[0] = a[1]*b[2] - a[2]*b[1]
	v[1] = a[0]*b[2] - a[2]*b[0]
	v[2] = a[0]*b[1] - a[1]*b[0]

	return v
}

func VectorNormalize(a *Vector) *Vector {
	b := new(Vector)
	b[0], b[1], b[2] = a[0], a[1], a[2]
	b.Normalize()
	return b
}

func (v *Vector) Normalize() {
	l := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	v[0] /= l
	v[1] /= l
	v[2] /= l
}
