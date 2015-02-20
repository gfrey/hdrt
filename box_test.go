package hdrt

import (
	"testing"

	"github.com/gfrey/hdrt/vec"
)

func TestBoxIntersection(t *testing.T) {
	b := &objBox{
		BaseObject: &BaseObject{Position: vec.New(2.0, 0.0, 0.0)},
		Width:      1.0,
		Height:     1.0,
		Depth:      1.0,
	}

	i1 := b.Intersect(vec.New(0.0, 0.0, 0.0), vec.New(1.0, 0.0, 0.0))
	i1Exp := vec.New(1.5, 0.0, 0.0)
	if i1 == nil || !vec.Equal(i1, i1Exp) {
		t.Errorf("expected i1 to be %s, got %s", i1Exp, i1)
	}

	i2 := b.Intersect(vec.New(0.0, 5.0, 0.0), vec.New(1.0, 0.0, 0.0))
	if i2 != nil {
		t.Errorf("expected i2 to be nil, got %s", i2)
	}
}

func TestIntersectPlane(t *testing.T) {
	r_pos, r_dir := vec.New(0.0, 0.0, 0.0), vec.New(1.0, 0.0, 0.0)

	tt := []struct {
		p0, n *vec.Vector
		exp   *vec.Vector
	}{
		{
			p0:  vec.New(2.0, 0.0, 0.0),
			n:   vec.New(-1.0, 0.0, 0.0),
			exp: vec.New(2.0, 0.0, 0.0),
		},
		{
			p0:  vec.New(2.0, 0.0, 0.0),
			n:   vec.New(1.0, 0.0, 0.0),
			exp: nil,
		},
		{
			p0:  vec.New(2.0, 0.0, 0.0),
			n:   vec.New(0.0, 1.0, 0.0),
			exp: r_pos,
		},
	}

	for i := range tt {
		cand := intersectPlane(r_pos, r_dir, tt[i].p0, tt[i].n)
		if !vec.Equal(cand, tt[i].exp) {
			t.Errorf("expected test %d to have %s, got %s", i, tt[i].exp, cand)
		}
	}
}

func TestPointInPlane(t *testing.T) {
	av := vec.New(1.0, 0.0, 0.0)
	bv := vec.New(0.0, 1.0, 0.0)
	tt := []struct {
		a, b, c *vec.Vector
		exp     bool
	}{
		{
			a: av, b: bv, c: vec.New(-1.0, -1.0, 0.0),
			exp: false,
		},
		{
			a: av, b: bv, c: vec.New(0.0, 0.0, 0.0),
			exp: true,
		},
		{
			a: av, b: bv, c: vec.New(0.5, 0.5, 0.0),
			exp: true,
		},
		{
			a: av, b: bv, c: vec.New(1.0, 1.0, 0.0),
			exp: true,
		},
		{
			a: av, b: bv, c: vec.New(2.0, 2.0, 0.0),
			exp: false,
		},
	}

	for i := range tt {
		got := pointInPlane(tt[i].a, tt[i].b, tt[i].c)
		if got != tt[i].exp {
			if tt[i].exp {
				t.Errorf("expected point %s to be in plane, wasn't", tt[i].c)
			} else {
				t.Errorf("expected point %s to not be in plane, was", tt[i].c)
			}
		}
	}
}
