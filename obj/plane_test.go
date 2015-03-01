package obj

import (
	"testing"

	"github.com/gfrey/hdrt/vec"
)

func TestPlaneIntersection(t *testing.T) {

	tt := []struct {
		Pos  *vec.Vector
		Up   *vec.Vector
		RPos *vec.Vector
		RDir *vec.Vector
		Exp  *vec.Vector
	}{
		{
			Pos: vec.New(1, 0, 0), Up: vec.New(-1, 0, 0),
			RPos: vec.New(0, 1, 0), RDir: vec.New(1, 0, 0),
			Exp: vec.New(1, 1, 0),
		},
		{
			Pos: vec.New(1, 0, 0), Up: vec.New(-1, 0, 0),
			RPos: vec.New(-5, 1, 0), RDir: vec.New(2, 0, 0),
			Exp: vec.New(1, 1, 0),
		},
		{
			Pos: vec.New(1, 0, 0), Up: vec.New(-1, 0, 0),
			RPos: vec.New(5, 1, 0), RDir: vec.New(1, 0, 0),
			Exp: nil,
		},
		{
			Pos: vec.New(1, 0, 0), Up: vec.New(-1, 0, 0),
			RPos: vec.New(0, 0, 0), RDir: vec.New(0, 1, 0),
			Exp: nil,
		},
		{
			Pos: vec.New(0, 0, 0), Up: vec.New(0, 1, 0),
			RPos: vec.New(0, 2, 0), RDir: vec.New(0, -1, 0),
			Exp: vec.New(0, 0, 0),
		},
		{
			Pos: vec.New(3, 0, 0), Up: vec.New(0, 1, 0),
			RPos: vec.New(1, 2, 0), RDir: vec.New(0, -1, 0),
			Exp: vec.New(1, 0, 0),
		},
	}

	for i := range tt {
		p := &objPlane{
			BaseObject: &BaseObject{Position: tt[i].Pos},
			Up:         tt[i].Up,
		}

		got := p.Intersect(tt[i].RPos, tt[i].RDir)

		if !vec.Equal(got, tt[i].Exp) {
			t.Errorf("in test %d %s was expected, got %s", i+1, tt[i].Exp, got)
		}
	}

}