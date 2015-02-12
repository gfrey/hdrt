package hdrt

import "testing"

func TestBoxIntersection(t *testing.T) {
	b := &objBox{
		BaseObject: &BaseObject{Position: NewVector(2.0, 0.0, 0.0)},
		Width:      1.0,
		Height:     1.0,
		Depth:      1.0,
	}

	i1 := b.Intersect(NewVector(0.0, 0.0, 0.0), NewVector(1.0, 0.0, 0.0))
	i1Exp := NewVector(1.5, 0.0, 0.0)
	if i1 == nil || !VectorEqual(i1, i1Exp, epsilon) {
		t.Errorf("expected i1 to be %s, got %s", i1Exp, i1)
	}

	i2 := b.Intersect(NewVector(0.0, 5.0, 0.0), NewVector(1.0, 0.0, 0.0))
	if i2 != nil {
		t.Errorf("expected i2 to be nil, got %s", i2)
	}
}

func TestIntersectPlane(t *testing.T) {
	r_pos, r_dir := NewVector(0.0, 0.0, 0.0), NewVector(1.0, 0.0, 0.0)

	tt := []struct {
		p0, n *Vector
		exp   *Vector
	}{
		{
			p0:  NewVector(2.0, 0.0, 0.0),
			n:   NewVector(-1.0, 0.0, 0.0),
			exp: NewVector(2.0, 0.0, 0.0),
		},
		{
			p0:  NewVector(2.0, 0.0, 0.0),
			n:   NewVector(1.0, 0.0, 0.0),
			exp: nil,
		},
		{
			p0:  NewVector(2.0, 0.0, 0.0),
			n:   NewVector(0.0, 1.0, 0.0),
			exp: r_pos,
		},
	}

	for i := range tt {
		cand := intersectPlane(r_pos, r_dir, tt[i].p0, tt[i].n)
		if !VectorEqual(cand, tt[i].exp, epsilon) {
			t.Errorf("expected test %d to have %s, got %s", i, tt[i].exp, cand)
		}
	}
}

func TestPointInPlane(t *testing.T) {
	av := NewVector(1.0, 0.0, 0.0)
	bv := NewVector(0.0, 1.0, 0.0)
	tt := []struct {
		a, b, c *Vector
		exp     bool
	}{
		{
			a: av, b: bv, c: NewVector(-1.0, -1.0, 0.0),
			exp: false,
		},
		{
			a: av, b: bv, c: NewVector(0.0, 0.0, 0.0),
			exp: true,
		},
		{
			a: av, b: bv, c: NewVector(0.5, 0.5, 0.0),
			exp: true,
		},
		{
			a: av, b: bv, c: NewVector(1.0, 1.0, 0.0),
			exp: true,
		},
		{
			a: av, b: bv, c: NewVector(2.0, 2.0, 0.0),
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
