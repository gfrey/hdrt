package vec

import "testing"

func TestVectorScalarMultiply(t *testing.T) {
	v := New(1.0, 2.0, 0.0)

	vP := ScalarMultiply(v, 2.0)
	if vP.Data[0] != 2.0 {
		t.Errorf("expected %.2f, got %.2f", 2.0, vP.Data[0])
	}
	if vP.Data[1] != 4.0 {
		t.Errorf("expected %.2f, got %.2f", 4.0, vP.Data[1])
	}
	if vP.Data[2] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, vP.Data[2])
	}

	v.ScalarMultiply(2.0)
	if v.Data[0] != 2.0 {
		t.Errorf("expected %.2f, got %.2f", 2.0, v.Data[0])
	}
	if v.Data[1] != 4.0 {
		t.Errorf("expected %.2f, got %.2f", 4.0, v.Data[1])
	}
	if v.Data[2] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, v.Data[2])
	}
}

func TestVectorAdd(t *testing.T) {
	a := New(1.0, 2.0, 3.0)
	b := New(3.0, 4.0, -5.0)

	c := Add(a, b)
	if c.Data[0] != 4.0 {
		t.Errorf("expected %.2f, got %.2f", 4.0, c.Data[0])
	}
	if c.Data[1] != 6.0 {
		t.Errorf("expected %.2f, got %.2f", 6.0, c.Data[1])
	}
	if c.Data[2] != -2.0 {
		t.Errorf("expected %.2f, got %.2f", -2.0, c.Data[2])
	}

	a.Add(b)
	if a.Data[0] != 4.0 {
		t.Errorf("expected %.2f, got %.2f", 4.0, a.Data[0])
	}
	if a.Data[1] != 6.0 {
		t.Errorf("expected %.2f, got %.2f", 6.0, a.Data[1])
	}
	if a.Data[2] != -2.0 {
		t.Errorf("expected %.2f, got %.2f", -2.0, a.Data[2])
	}
}

func TestVectorCross(t *testing.T) {
	a := New(1.0, 0.0, 0.0)
	b := New(0.0, 1.0, 0.0)

	c := Cross(a, b)
	if c.Data[0] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, c.Data[0])
	}
	if c.Data[1] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, c.Data[1])
	}
	if c.Data[2] != 1.0 {
		t.Errorf("expected %.2f, got %.2f", 1.0, c.Data[2])
	}
}

func TestNormalize(t *testing.T) {
	a := New(5.0, 0.0, 0.0)

	b := Normalize(a)
	if b.Data[0] != 1.0 {
		t.Errorf("expected %.2f, got %.2f", 1.0, b.Data[0])
	}
	if b.Data[1] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, b.Data[1])
	}
	if b.Data[2] != 0.0 {
		t.Errorf("expected %.2f, got %.2f", 0.0, b.Data[2])
	}
}
