package quat

import (
	"math"
	"testing"

	"github.com/gfrey/hdrt/vec"
)

func TestQuaternionAdd(t *testing.T) {
	a := New(1.0, 1.0, 1.0, 1.0)
	b := New(1.0, 2.0, 3.0, 4.0)

	c := Add(a, b)

	d := []struct {
		exp float64
		got float64
		msg string
	}{
		{2.0, c.Data[0], "X"},
		{3.0, c.Data[1], "Y"},
		{4.0, c.Data[2], "Z"},
		{5.0, c.Data[3], "W"},
	}

	epsilon := math.Nextafter(1, 2) - 1

	for i := range d {
		if d[i].got < d[i].exp-epsilon || d[i].exp+epsilon < d[i].got {
			t.Errorf("expected %q to be %.6f, got %.6f", d[i].msg, d[i].exp, d[i].got)
		}
	}
}

var (
	unitX = New(1.0, 0.0, 0.0, 0.0)
	unitY = New(0.0, 1.0, 0.0, 0.0)
	unitZ = New(0.0, 0.0, 1.0, 0.0)

	a = New(1.0, 2.0, 3.0, 4.0)
	b = New(5.0, 6.0, 7.0, 8.0)

	epsilon = math.Nextafter(1, 2) - 1
)

func TestQuaternionMultiply(t *testing.T) {
	d := []struct {
		exp *Quaternion
		got *Quaternion
		msg string
	}{
		{New(0.0, 0.0, 1.0, 0.0), Multiply(unitX, unitY), "unit Z"},
		{New(0.0, 1.0, 0.0, 0.0), Multiply(unitZ, unitX), "unit Y"},
		{New(1.0, 0.0, 0.0, 0.0), Multiply(unitY, unitZ), "unit X"},

		{New(0.0, 0.0, -1.0, 0.0), Multiply(unitY, unitX), "-unit Z"},
		{New(0.0, -1.0, 0.0, 0.0), Multiply(unitX, unitZ), "-unit Y"},
		{New(-1.0, 0.0, 0.0, 0.0), Multiply(unitZ, unitY), "-unit X"},

		{New(0.0, 0.0, 0.0, -1.0), Multiply(unitX, unitX), "uX*uX"},
		{New(0.0, 0.0, 0.0, -1.0), Multiply(unitY, unitY), "uY*uY"},
		{New(0.0, 0.0, 0.0, -1.0), Multiply(unitZ, unitZ), "uZ*uZ"},

		{New(24.0, 48.0, 48.0, -6.0), Multiply(a, b), "a*b"},
	}

	for i := range d {
		if !Equal(d[i].got, d[i].exp) {
			t.Errorf("in %s expected %s, got %s", d[i].msg, d[i].exp, d[i].got)
		}
	}
}

func TestQuaternionS(t *testing.T) {
	d := []struct {
		exp *Quaternion
		got *Quaternion
		msg string
	}{
		{New(0.0, 0.0, 0.0, -1.0), Square(unitX), "uX^2"},
		{New(0.0, 0.0, 0.0, -1.0), Square(unitY), "uY^2"},
		{New(0.0, 0.0, 0.0, -1.0), Square(unitZ), "uZ^2"},

		{New(8.0, 16.0, 24.0, 2.0), Square(a), "a^2"},
	}

	for i := range d {
		if !Equal(d[i].got, d[i].exp) {
			t.Errorf("in %s expected %s, got %s", d[i].msg, d[i].exp, d[i].got)
		}
	}
}

func TestQuaternionRotation(t *testing.T) {
	tt := []struct {
		q   *Quaternion
		v   *vec.Vector
		exp *vec.Vector
	}{
		{NewRotation(vec.New(1, 0, 0), 90), vec.New(0, 1, 0), vec.New(0, 0, 1)},
		{NewRotation(vec.New(1, 0, 0), 180), vec.New(0, 1, 0), vec.New(0, -1, 0)},
		{NewRotation(vec.New(1, 0, 0), 270), vec.New(0, 1, 0), vec.New(0, 0, -1)},
		{NewRotation(vec.New(1, 0, 0), 270), vec.New(0, 2, 0), vec.New(0, 0, -2)},
		{NewRotation(vec.New(2, 0, 0), 270), vec.New(0, 2, 0), vec.New(0, 0, -2)},
		{NewRotation(vec.New(0, 0, 1), 90), vec.New(1, 0, 0), vec.New(0, 1, 0)},
	}

	for i := range tt {
		got := tt[i].q.Rotate(tt[i].v)

		if !vec.Equal(got, tt[i].exp) {
			t.Errorf("in test %d expected rotated vector to be %s, got %s", i, tt[i].exp, got)
		}
	}

}