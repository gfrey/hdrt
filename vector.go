package hdrt

import (
	"encoding/json"
	"fmt"
	"math"
)

type Vector [3]float64

func NewVector(x, y, z float64) *Vector {
	v := new(Vector)
	v[0], v[1], v[2] = x, y, z
	return v
}

func (v *Vector) String() string {
	return fmt.Sprintf("[%.2f, %.2f, %.2f]", v[0], v[1], v[2])
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

func VectorMultiply(v1 *Vector, v2 *Vector) *Vector {
	return NewVector(v1[0]*v2[0], v1[1]*v2[1], v1[2]*v2[2])
}

func VectorScalarDivide(v *Vector, a float64) *Vector {
	return NewVector(v[0]/a, v[1]/a, v[2]/a)
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

func VectorSub(a, b *Vector) *Vector {
	return NewVector(a[0]-b[0], a[1]-b[1], a[2]-b[2])
}

func (v *Vector) Sub(other *Vector) {
	v[0] -= other[0]
	v[1] -= other[1]
	v[2] -= other[2]
}

func VectorCross(a, b *Vector) *Vector {
	v := new(Vector)

	v[0] = a[1]*b[2] - a[2]*b[1]
	v[1] = a[0]*b[2] - a[2]*b[0]
	v[2] = a[0]*b[1] - a[1]*b[0]

	return v
}

func VectorDot(a, b *Vector) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func VectorNormalize(a *Vector) *Vector {
	b := new(Vector)
	b[0], b[1], b[2] = a[0], a[1], a[2]
	b.Normalize()
	return b
}

func (v *Vector) Normalize() {
	l := v.Length()
	v[0] /= l
	v[1] /= l
	v[2] /= l
}

func VectorProject(v *Vector, u *Vector) *Vector {
	vu := VectorMultiply(v, u)
	vl := v.Length()

	puv := VectorMultiply(VectorScalarDivide(vu, vl), v)

	return puv
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func (v *Vector) DistanceTo(v2 *Vector) float64 {
	return VectorSub(v2, v).Length()
}

var epsilon = 0.0001

func VectorEqual(a, b *Vector, ε float64) bool {
	return FloatEqual(a[0], b[0], ε) && FloatEqual(a[1], b[1], ε) && FloatEqual(a[2], b[2], ε)
}

func FloatEqual(a, b, ε float64) bool {
	return a-ε < b && b < a+ε
}
