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

	v.Span(c)

	if v.pos[0] != 1.0 {
		t.Errorf("expected v.pos[0] = %.2f, got %.2f", 1.0, v.pos[0])
	}
	if v.pos[1] != 0.5 {
		t.Errorf("expected v.pos[1] = %.2f, got %.2f", 0.5, v.pos[1])
	}
	if v.pos[2] != -0.5 {
		t.Errorf("expected v.pos[2] = %.2f, got %.2f", -0.5, v.pos[2])
	}

	if v.a[0] != 0.0 {
		t.Errorf("expected v.a[0] = %.2f, got %.2f", 0.0, v.a[0])
	}
	if v.a[1] != 0.0 {
		t.Errorf("expected v.a[1] = %.2f, got %.2f", 0.0, v.a[1])
	}
	if v.a[2] != 1.0 {
		t.Errorf("expected v.a[2] = %.2f, got %.2f", 1.0, v.a[2])
	}

	if v.b[0] != 0.0 {
		t.Errorf("expected v.b[0] = %.2f, got %.2f", 0.0, v.b[0])
	}
	if v.b[1] != -1.0 {
		t.Errorf("expected v.b[1] = %.2f, got %.2f", -1.0, v.b[1])
	}
	if v.b[2] != 0.0 {
		t.Errorf("expected v.b[2] = %.2f, got %.2f", 0.0, v.b[2])
	}
}
