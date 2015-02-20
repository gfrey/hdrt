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

	if len(wrld.Scene.Objects) != 3 {
		t.Errorf("expected %d objects, got %d", 3, len(wrld.Scene.Objects))
	}
}
