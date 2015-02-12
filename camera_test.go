package hdrt

import "testing"

func TestViewplaneInit(t *testing.T) {
	c := new(Camera)
	c.Position = NewVector(0.0, 0.0, 0.0)
	c.Direction = NewVector(1.0, 0.0, 0.0)
	c.Up = NewVector(0.0, 1.0, 0.0)
	c.FOV = 90.0

	v := new(Viewplane)
	v.Distance = 1.0
	v.ResX, v.ResY = 100, 100

	v.span(c)

	posExp := NewVector(1.0, 0.5, -0.5)
	if !VectorEqual(v.pos, posExp, epsilon) {
		t.Errorf("expected v.pos = %s, got %s", posExp, v.pos)
	}

	aExp := NewVector(0.0, 0.0, 1.0)
	if !VectorEqual(v.a, aExp, epsilon) {
		t.Errorf("expected v.a = %s, got %s", aExp, v.a)
	}

	bExp := NewVector(0.0, -1.0, 0.0)
	if !VectorEqual(v.b, bExp, epsilon) {
		t.Errorf("expected v.b = %s, got %s", bExp, v.b)
	}

	pixelPos1 := v.PixelPosition(0, 0)
	if !VectorEqual(pixelPos1, posExp, epsilon) {
		t.Errorf("expected pixelPos1 = %s, got %s", posExp, pixelPos1)
	}

	pixelPos2 := v.PixelPosition(v.ResX/2, v.ResY/2)
	pixelPos2Exp := NewVector(1.0, 0.0, 0.0)
	if !VectorEqual(pixelPos2, pixelPos2Exp, epsilon) {
		t.Errorf("expected pixelPos2 = %s, got %s", pixelPos2Exp, pixelPos2)
	}
}
