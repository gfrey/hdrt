package hdrt

import "testing"

func TestVectorScalarMultiply(t *testing.T) {
	v := NewVector(1.0, 2.0, 0.0)

	vP := VectorScalarMultiply(v, 2.0)
	if vP[0] != 2.0 {
		t.Errorf("expected %.2f, got %.2f", 2.0, vP[0])
	}
	if vP[1] != 4.0 {
		t.Errorf("expected %.2f, got %.2f", 4.0, vP[1])
	}
	if vP[2] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, vP[2])
	}

	v.ScalarMultiply(2.0)
	if v[0] != 2.0 {
		t.Errorf("expected %.2f, got %.2f", 2.0, v[0])
	}
	if v[1] != 4.0 {
		t.Errorf("expected %.2f, got %.2f", 4.0, v[1])
	}
	if v[2] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, v[2])
	}
}

func TestVectorAdd(t *testing.T) {
	a := NewVector(1.0, 2.0, 3.0)
	b := NewVector(3.0, 4.0, -5.0)

	c := VectorAdd(a, b)
	if c[0] != 4.0 {
		t.Errorf("expected %.2f, got %.2f", 4.0, c[0])
	}
	if c[1] != 6.0 {
		t.Errorf("expected %.2f, got %.2f", 6.0, c[1])
	}
	if c[2] != -2.0 {
		t.Errorf("expected %.2f, got %.2f", -2.0, c[2])
	}

	a.Add(b)
	if a[0] != 4.0 {
		t.Errorf("expected %.2f, got %.2f", 4.0, a[0])
	}
	if a[1] != 6.0 {
		t.Errorf("expected %.2f, got %.2f", 6.0, a[1])
	}
	if a[2] != -2.0 {
		t.Errorf("expected %.2f, got %.2f", -2.0, a[2])
	}
}

func TestVectorCross(t *testing.T) {
	a := NewVector(1.0, 0.0, 0.0)
	b := NewVector(0.0, 1.0, 0.0)

	c := VectorCross(a, b)
	if c[0] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, c[0])
	}
	if c[1] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, c[1])
	}
	if c[2] != 1.0 {
		t.Errorf("expected %.2f, got %.2f", 1.0, c[2])
	}
}

func TestVectorNormalize(t *testing.T) {
	a := NewVector(5.0, 0.0, 0.0)

	b := VectorNormalize(a)
	if b[0] != 1.0 {
		t.Errorf("expected %.2f, got %.2f", 1.0, b[0])
	}
	if b[1] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, b[1])
	}
	if b[2] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, b[2])
	}
}
