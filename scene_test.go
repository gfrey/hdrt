package hdrt

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSceneDescription(t *testing.T) {
	fh, err := os.Open("default_scene.json")
	if err != nil {
		t.Fatalf("failed to read default scene: %s", err)
	}

	wrld := new(World)
	err = json.NewDecoder(fh).Decode(&wrld)
	if err != nil {
		t.Fatalf("failed to decode world: %s", err)
	}

	if wrld.Camera == nil {
		t.Errorf("camera not set")
	}
	if wrld.Viewplane == nil {
		t.Errorf("view plane not set")
	}
	if wrld.Scene == nil {
		t.Errorf("scene not set")
	}

	if len(wrld.Scene.Objects) != 2 {
		t.Errorf("expected %d objects, got %d", 2, len(wrld.Scene.Objects))
	}

	if v, ok := wrld.Scene.Objects[0].(*objSphere); !ok {
		t.Errorf("expected first object to be a sphere, got %T", wrld.Scene.Objects[0])
	} else {
		if v.Radius != 0.75 {
			t.Errorf("expected sphere's radius to be 0.75, got %.6f", v.Radius)
		}
	}

	if v, ok := wrld.Scene.Objects[1].(*objBox); !ok {
		t.Errorf("expected first object to be a sphere, got %T", wrld.Scene.Objects[1])
	} else {
		if v.Width != 0.8 {
			t.Errorf("expected box's width to be 0.8, got %.6f", v.Width)
		}
		if v.Height != 1.0 {
			t.Errorf("expected box's height to be 1.0, got %.6f", v.Height)
		}
		if v.Depth != 0.8 {
			t.Errorf("expected box's depth to be 0.8, got %.6f", v.Depth)
		}
	}
}
