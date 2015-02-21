package hdrt

import (
	"encoding/json"
	"image/color"
	"math"

	"github.com/gfrey/hdrt/obj"
	"github.com/gfrey/hdrt/vec"
)

type Scene struct {
	AmbientLight *obj.MaterialC
	Objects      []obj.Object `json:"objects"`
	Lights       []*Light     `json:"lights"`
}

func (sc *Scene) UnmarshalJSON(data []byte) error {
	rsc := &struct {
		AmbientLight *obj.MaterialC
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
	for i := range sc.Lights {
		sc.Lights[i].Direction.Normalize()
	}
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
			d := vec.Sub(p, pos).Length()
			if cand == nil || d < distance {
				ipos = p
				cand = sc.Objects[i]
				distance = d
			}
		}
	}
	if cand != nil {
		return sc.ColorWithLights(cand, ipos, pos)
	}
	return &color.RGBA{0, 0, 0, 0}
}

func (sc *Scene) ColorWithLights(o obj.Object, ipos, spos *vec.Vector) *color.RGBA {
	normal := o.Normal(ipos)
	var r, g, b uint
	ma := o.Material(obj.MATERIAL_AMBIENT)
	if ma != nil && sc.AmbientLight != nil {
		r = uint(ma[0]) * uint(sc.AmbientLight[0]) / 255
		g = uint(ma[1]) * uint(sc.AmbientLight[1]) / 255
		b = uint(ma[2]) * uint(sc.AmbientLight[2]) / 255
	}

LIGHTSOURCES:
	for i := range sc.Lights {
		lcoeff, ldif, lspec := sc.Lights[i].InCone(ipos, normal)
		if ldif != nil || lspec != nil {
			lPos := sc.Lights[i].Position
			dir := vec.Sub(lPos, ipos)
			dist := dir.Length()
			dir.Normalize()
			// is obj "visible" from light source?
			for j := range sc.Objects {
				if sc.Objects[j] == o {
					continue
				}

				tmpPos := sc.Objects[j].Intersect(lPos, dir)
				if tmpPos == nil {
					continue
				}

				d := vec.Sub(tmpPos, lPos).Length()
				if d < dist { // in shadow
					continue LIGHTSOURCES
				}
			}
			ln := vec.Dot(dir, normal)

			mdif := o.Material(obj.MATERIAL_DIFFUSE)
			if mdif != nil && ldif != nil {
				lnc := ln * lcoeff
				r += uint(float64(mdif[0]*ldif[0])*lnc) / 255
				g += uint(float64(mdif[1]*ldif[1])*lnc) / 255
				b += uint(float64(mdif[2]*ldif[2])*lnc) / 255
			}

			mspec := o.Material(obj.MATERIAL_SPECULAR)
			if mspec != nil && lspec != nil {
				v := vec.Sub(spos, ipos).Normalize()

				rm := vec.ScalarMultiply(normal, 2*ln).Sub(dir).Normalize()
				rn := math.Pow(vec.Dot(rm, v), o.Reflection()) * lcoeff
				r += uint(float64(mspec[0]*lspec[0])*rn) / 255
				g += uint(float64(mspec[1]*lspec[1])*rn) / 255
				b += uint(float64(mspec[2]*lspec[2])*rn) / 255
			}
		}
	}

	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}

	return &color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}
