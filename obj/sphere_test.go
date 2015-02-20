package obj

import (
	"testing"

	"github.com/gfrey/hdrt/vec"
)

func TestObjSphereIntersect(t *testing.T) {
	o := &objSphere{
		BaseObject: &BaseObject{},
		Radius:     0.75,
	}

	rayPos := vec.New(0, 0, 0)
	dir := vec.New(1, 0, 0)

	_, _ = rayPos, dir

	o.Position = vec.New(0.1, 0, 0)
	if o.Intersect(rayPos, dir) == nil {
		t.Errorf("did expect to intersect when sphere pos is %s, but did not", o.Position)
	}

	o.Position = vec.New(0.5, 0, 0)
	if o.Intersect(rayPos, dir) == nil {
		t.Errorf("did expect to intersect when sphere pos is %s, but did not", o.Position)
	}

	o.Position = vec.New(5, 0.5, 0)
	if o.Intersect(rayPos, dir) == nil {
		t.Errorf("did expect to intersect when sphere pos is %s, but did not", o.Position)
	}

	o.Position = vec.New(5, 0.74, 0)
	if o.Intersect(rayPos, dir) == nil {
		t.Errorf("did expect to intersect when sphere pos is %s, but did not", o.Position)
	}

	o.Position = vec.New(5, 1.1, 0)
	if o.Intersect(rayPos, dir) != nil {
		t.Errorf("did expect to NOT intersect when sphere pos is %s, but DID", o.Position)
	}
}
