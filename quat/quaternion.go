package quat

import (
	"math"

	"github.com/gfrey/hdrt/mat"
	"github.com/gfrey/hdrt/vec"
)

// Xi+Yj+Zk+W
type Quaternion struct {
	X, Y, Z, W float64
	conj       *Quaternion
}

func New(x, y, z, w float64) *Quaternion {
	q := new(Quaternion)
	q.X, q.Y, q.Z, q.W = x, y, z, w
	return q
}

func NewRotation(axis *vec.Vector, angle float64) *Quaternion {
	a := *axis
	a.Normalize()
	q := new(Quaternion)
	α := mat.Deg2Rad(angle) / 2.0
	sin := math.Sin(α)
	a.ScalarMultiply(sin)
	q.X, q.Y, q.Z = a.Data[0], a.Data[1], a.Data[2]
	q.W = math.Cos(α)
	return q
}

func Add(a, b *Quaternion) *Quaternion {
	q := new(Quaternion)
	*q = *a
	q.X += b.X
	q.Y += b.Y
	q.Z += b.Z
	q.W += b.W
	return q
}

func (q *Quaternion) Multiply(b *Quaternion) *Quaternion {
	x, y, z, w := q.X, q.Y, q.Z, q.W
	q.W = w*b.W - x*b.X - y*b.Y - z*b.Z
	q.X = w*b.X + x*b.W + y*b.Z - z*b.Y
	q.Y = w*b.Y - x*b.Z + y*b.W + z*b.X
	q.Z = w*b.Z + x*b.Y - y*b.X + z*b.W
	return q
}

func Multiply(a, b *Quaternion) *Quaternion {
	q := new(Quaternion)
	q.W = a.W*b.W - a.X*b.X - a.Y*b.Y - a.Z*b.Z
	q.X = a.W*b.X + a.X*b.W + a.Y*b.Z - a.Z*b.Y
	q.Y = a.W*b.Y - a.X*b.Z + a.Y*b.W + a.Z*b.X
	q.Z = a.W*b.Z + a.X*b.Y - a.Y*b.X + a.Z*b.W
	return q
}

func Square(a *Quaternion) *Quaternion {
	q := new(Quaternion)
	q.W = a.W*a.W - a.X*a.X - a.Y*a.Y - a.Z*a.Z
	q.X = 2.0 * a.W * a.X
	q.Y = 2.0 * a.W * a.Y
	q.Z = 2.0 * a.W * a.Z
	return q
}

func Conjugate(a *Quaternion) *Quaternion {
	q := new(Quaternion)
	q.W = a.W
	q.X = -a.X
	q.Y = -a.Y
	q.Z = -a.Z
	return q
}

func (q *Quaternion) Rotate(v *vec.Vector) *vec.Vector {
	if q.conj == nil {
		q.conj = Conjugate(q)
	}

	vp := new(Quaternion)
	vp.X, vp.Y, vp.Z = v.Data[0], v.Data[1], v.Data[2]

	res := Multiply(q, vp).Multiply(q.conj)

	return vec.New(res.X, res.Y, res.Z)
}