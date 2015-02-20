package hdrt

import (
	"encoding/json"
	"image/color"

	"github.com/gfrey/hdrt/obj"
	"github.com/gfrey/hdrt/vec"
)

type Scene struct {
	AmbientLight float64
	Objects      []obj.Object `json:"objects"`
	Lights       []*Light     `json:"lights"`
}

func (sc *Scene) UnmarshalJSON(data []byte) error {
	rsc := &struct {
		AmbientLight float64
		Objects      []*obj.Raw
		Lights       []*Light
	}{}

	err := json.Unmarshal(data, &rsc)
	if err != nil {
		return err
	}

	for i := range rsc.Objects {
		sc.Objects = append(sc.Objects, rsc.Objects[i].Object())
	}

	sc.AmbientLight = rsc.AmbientLight
	sc.Lights = rsc.Lights
	return nil
}

func (sc *Scene) Render(pos, dir *vec.Vector) *color.RGBA {
	var (
		cand     obj.Object
		distance float64
		ipos     *vec.Vector
	)
	for i := range sc.Objects {
		o := sc.Objects[i]
		p := o.Intersect(pos, dir)
		if p != nil {
			d := (p.Data[0] - pos.Data[0]) / dir.Data[0]
			if cand == nil || d < distance {
				ipos = p
				cand = sc.Objects[i]
				distance = d
			}
		}
	}
	if cand != nil {
		return sc.ColorWithLights(cand, ipos)
	}
	return &color.RGBA{0, 0, 0, 0}
}

func (sc *Scene) ColorWithLights(obj obj.Object, pos *vec.Vector) *color.RGBA {
	baseLight := sc.AmbientLight
	normal := obj.Normal(pos)
LIGHTSOURCES:
	for i := range sc.Lights {
		dist, delta, dir := sc.Lights[i].InCone(pos, normal)
		if dir != nil {
			lPos := sc.Lights[i].Position
			for j := range sc.Objects {
				if sc.Objects[j] == obj {
					continue
				}

				tmpPos := sc.Objects[j].Intersect(lPos, dir)
				if tmpPos == nil {
					continue
				}
				d := (tmpPos.Data[0] - lPos.Data[0]) / dir.Data[0]
				if d < dist { // in shadow
					continue LIGHTSOURCES
				}
			}
			// not in shadow
			maxAngle := deg2rad(sc.Lights[i].Angle) / 2.0
			baseLight += (1.0 - (dist / sc.Lights[i].Distance)) * (maxAngle - delta) / maxAngle
		}
	}

	if baseLight > 1.0 {
		baseLight = 1.0
	}

	c := obj.GetColor()
	return &color.RGBA{
		uint8(float64(c.R) * baseLight),
		uint8(float64(c.G) * baseLight),
		uint8(float64(c.B) * baseLight),
		c.A}
}
