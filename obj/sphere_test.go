package obj

import (
	"testing"

	"github.com/gfrey/hdrt/vec"
)

func TestObjSphereIntersect(t *testing.T) {
	tt := []struct {
		Pos    *vec.Vector
		Radius float64
		RPos   *vec.Vector
		RDir   *vec.Vector
		Exp    *vec.Vector
	}{
		{
			Pos: vec.New(2, 0, 0), Radius: 1,
			RPos: vec.New(0, 0, 0), RDir: vec.New(1, 0, 0),
			Exp: vec.New(1, 0, 0),
		},
		{
			Pos: vec.New(2, 0, 0), Radius: 1,
			RPos: vec.New(4, 0, 0), RDir: vec.New(-1, 0, 0),
			Exp: vec.New(3, 0, 0),
		},
		{
			Pos: vec.New(2, 0, 0), Radius: 1,
			RPos: vec.New(2, -3, 0), RDir: vec.New(0, 1, 0),
			Exp: vec.New(2, -1, 0),
		},
		{
			Pos: vec.New(2, 0, 0), Radius: 1,
			RPos: vec.New(2, -3, 0), RDir: vec.New(0, -1, 0),
			Exp: nil,
		},
		{
			Pos: vec.New(2, 0, 0), Radius: 0.75,
			RPos: vec.New(2, 3, 0), RDir: vec.New(0, -1, 0),
			Exp: vec.New(2, 0.75, 0),
		},
	}

	for i := range tt {
		o := &objSphere{
			BaseObject: &BaseObject{Position: tt[i].Pos},
			Radius:     tt[i].Radius,
		}

		got := o.Intersect(tt[i].RPos, tt[i].RDir)

		if !vec.Equal(got, tt[i].Exp) {
			t.Errorf("in test %d %s was expected, got %s", i+1, tt[i].Exp, got)
		}

	}
}
