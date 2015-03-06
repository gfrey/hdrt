package quat

import (
	"fmt"
	"math"

	"github.com/gfrey/hdrt/mat"
	"github.com/gfrey/hdrt/vec"
)

// Xi+Yj+Zk+W
type Quaternion struct {
	Data [4]float64
	conj *Quaternion
}

func (q *Quaternion) String() string {
	return fmt.Sprintf("[%.3f,%.3f,%.3f|%.3f]", q.Data[0], q.Data[1], q.Data[2], q.Data[3])
}

func New(x, y, z, w float64) *Quaternion {
	return &Quaternion{Data: [4]float64{x, y, z, w}}
}

func Copy(a *Quaternion) *Quaternion {
	q := new(Quaternion)
	q.Data[0] = a.Data[0]
	q.Data[1] = a.Data[1]
	q.Data[2] = a.Data[2]
	q.Data[3] = a.Data[3]
	return q
}

func NewRotation(axis *vec.Vector, angle float64) *Quaternion {
	a := vec.Copy(axis).Normalize()
	α := mat.Deg2Rad(angle) * 0.5
	sin := math.Sin(α)
	a.ScalarMultiply(sin)
	return New(a.Data[0], a.Data[1], a.Data[2], math.Cos(α))
}

func Add(a, b *Quaternion) *Quaternion {
	q := Copy(a)
	q.Data[0] += b.Data[0]
	q.Data[1] += b.Data[1]
	q.Data[2] += b.Data[2]
	q.Data[3] += b.Data[3]
	return q
}

func (q *Quaternion) Normalize() *Quaternion {
	sum := 0.0
	for i := 0; i < 4; i++ {
		sum += q.Data[i] * q.Data[i]
	}
	sum = 1 / sum
	for i := 0; i < 4; i++ {
		q.Data[i] *= sum
	}
	return q
}

func (q *Quaternion) Multiply(b *Quaternion) *Quaternion {
	var tmp [4]float64 = q.Data
	q.Data[3] = tmp[3]*b.Data[3] - tmp[0]*b.Data[0] - tmp[1]*b.Data[1] - tmp[2]*b.Data[2]
	q.Data[0] = tmp[3]*b.Data[0] + tmp[0]*b.Data[3] + tmp[1]*b.Data[2] - tmp[2]*b.Data[1]
	q.Data[1] = tmp[3]*b.Data[1] - tmp[0]*b.Data[2] + tmp[1]*b.Data[3] + tmp[2]*b.Data[0]
	q.Data[2] = tmp[3]*b.Data[2] + tmp[0]*b.Data[1] - tmp[1]*b.Data[0] + tmp[2]*b.Data[3]
	return q
}

func Multiply(a, b *Quaternion) *Quaternion {
	q := New(a.Data[0], a.Data[1], a.Data[2], a.Data[3])
	return q.Multiply(b)
}

func Square(a *Quaternion) *Quaternion {
	q := new(Quaternion)
	q.Data[3] = a.Data[3]*a.Data[3] - a.Data[0]*a.Data[0] - a.Data[1]*a.Data[1] - a.Data[2]*a.Data[2]
	t := 2.0 * a.Data[3]
	q.Data[0] = t * a.Data[0]
	q.Data[1] = t * a.Data[1]
	q.Data[2] = t * a.Data[2]
	return q
}

func Conjugate(a *Quaternion) *Quaternion {
	return New(-a.Data[0], -a.Data[1], -a.Data[2], a.Data[3])
}

func (q *Quaternion) Rotate(v *vec.Vector) *vec.Vector {
	if q.conj == nil {
		q.conj = Conjugate(q)
	}

	vp := New(v.Data[0], v.Data[1], v.Data[2], 0.0)
	res := Multiply(q, vp).Multiply(q.conj)
	return vec.New(res.Data[0], res.Data[1], res.Data[2])
}

func Equal(a, b *Quaternion) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil, b == nil:
		return false
	default:
		for i := 0; i < 4; i++ {
			if !mat.FloatEqual(a.Data[i], b.Data[i]) {
				return false
			}
		}
		return true
	}
}