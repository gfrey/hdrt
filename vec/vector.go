package vec

import (
	"encoding/json"
	"fmt"
	"math"
)

type Vector struct {
	Data   [3]float64
	length *float64
}

func NewVector(x, y, z float64) *Vector {
	v := new(Vector)
	v.Data[0], v.Data[1], v.Data[2] = x, y, z
	return v
}

func (v *Vector) String() string {
	return fmt.Sprintf("[%.2f, %.2f, %.2f]", v.Data[0], v.Data[1], v.Data[2])
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

	v.Data[0] = tv.X
	v.Data[1] = tv.Y
	v.Data[2] = tv.Z

	return nil
}

func VectorScalarMultiply(v *Vector, a float64) *Vector {
	return NewVector(v.Data[0]*a, v.Data[1]*a, v.Data[2]*a)
}

func VectorScalarDivide(v *Vector, a float64) *Vector {
	return NewVector(v.Data[0]/a, v.Data[1]/a, v.Data[2]/a)
}

func (v *Vector) ScalarMultiply(a float64) {
	v.Data[0] *= a
	v.Data[1] *= a
	v.Data[2] *= a
	if v.length != nil {
		*v.length *= a
	}
}

func VectorAdd(a, b *Vector) *Vector {
	return NewVector(a.Data[0]+b.Data[0], a.Data[1]+b.Data[1], a.Data[2]+b.Data[2])
}

func (v *Vector) Add(other *Vector) {
	v.length = nil
	v.Data[0] += other.Data[0]
	v.Data[1] += other.Data[1]
	v.Data[2] += other.Data[2]
}

func VectorSub(a, b *Vector) *Vector {
	return NewVector(a.Data[0]-b.Data[0], a.Data[1]-b.Data[1], a.Data[2]-b.Data[2])
}

func (v *Vector) Sub(other *Vector) {
	v.length = nil
	v.Data[0] -= other.Data[0]
	v.Data[1] -= other.Data[1]
	v.Data[2] -= other.Data[2]
}

func VectorCross(a, b *Vector) *Vector {
	v := new(Vector)

	v.Data[0] = a.Data[1]*b.Data[2] - a.Data[2]*b.Data[1]
	v.Data[1] = a.Data[2]*b.Data[0] - a.Data[0]*b.Data[2]
	v.Data[2] = a.Data[0]*b.Data[1] - a.Data[1]*b.Data[0]

	return v
}

func VectorDot(a, b *Vector) float64 {
	return a.Data[0]*b.Data[0] + a.Data[1]*b.Data[1] + a.Data[2]*b.Data[2]
}

func VectorNormalize(a *Vector) *Vector {
	b := new(Vector)
	b.Data[0], b.Data[1], b.Data[2] = a.Data[0], a.Data[1], a.Data[2]
	b.Normalize()
	return b
}

func (v *Vector) Normalize() {
	l := v.Length()
	v.Data[0] /= l
	v.Data[1] /= l
	v.Data[2] /= l
	// we normalized didn't we? The length computation at the beginning will have set the length value
	*v.length = 1.0
}

func VectorProject(u *Vector, v *Vector) *Vector {
	vu := VectorDot(v, u)
	vl := v.Length()
	return VectorScalarMultiply(v, vu/vl)
}

func (v *Vector) Length() float64 {
	if v.length == nil {
		v.length = new(float64)
		*v.length = math.Sqrt(v.Data[0]*v.Data[0] + v.Data[1]*v.Data[1] + v.Data[2]*v.Data[2])
	}
	return *v.length
}

func (v *Vector) DistanceTo(v2 *Vector) float64 {
	return VectorSub(v2, v).Length()
}

var epsilon = 0.0001

func VectorEqual(a, b *Vector, ε float64) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil, b == nil:
		return false
	default:
		return FloatEqual(a.Data[0], b.Data[0], ε) &&
			FloatEqual(a.Data[1], b.Data[1], ε) &&
			FloatEqual(a.Data[2], b.Data[2], ε)
	}
}

func FloatEqual(a, b, ε float64) bool {
	return a-ε < b && b < a+ε
}

func FloatLessThan(a, b, ε float64) bool {
	return a < b+ε
}

func FloatLessThanEqual(a, b, ε float64) bool {
	return a <= b+ε
}

func FloatGreaterThan(a, b, ε float64) bool {
	return a > b-ε
}
