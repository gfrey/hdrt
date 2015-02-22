package hdrt

import (
	"testing"

	"github.com/gfrey/hdrt/vec"
)

func TestViewplaneInit(t *testing.T) {
	c := new(Camera)
	c.Position = vec.New(0.0, 0.0, 0.0)
	c.FOV = 90.0
	err := c.Init()
	if err != nil {
		t.Fatalf("failed to initialize the camera: %s", err)
	}

	v := new(Viewplane)
	v.Distance = 1.0
	v.ResX, v.ResY = 100, 100

	v.span(c)

	posExp := vec.New(1.0, 0.5, -0.5)
	if !vec.Equal(v.pos, posExp) {
		t.Errorf("expected v.pos = %s, got %s", posExp, v.pos)
	}

	aExp := vec.New(0.0, 0.0, 1.0)
	if !vec.Equal(v.a, aExp) {
		t.Errorf("expected v.a = %s, got %s", aExp, v.a)
	}

	bExp := vec.New(0.0, -1.0, 0.0)
	if !vec.Equal(v.b, bExp) {
		t.Errorf("expected v.b = %s, got %s", bExp, v.b)
	}

	pixelPos1 := v.PixelPosition(0, 0)
	if !vec.Equal(pixelPos1, posExp) {
		t.Errorf("expected pixelPos1 = %s, got %s", posExp, pixelPos1)
	}

	pixelPos2 := v.PixelPosition(v.ResX/2, v.ResY/2)
	pixelPos2Exp := vec.New(1.0, 0.0, 0.0)
	if !vec.Equal(pixelPos2, pixelPos2Exp) {
		t.Errorf("expected pixelPos2 = %s, got %s", pixelPos2Exp, pixelPos2)
	}
}
