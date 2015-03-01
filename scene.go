package hdrt

import (
	"encoding/json"
	"image/color"
	"math"

	"github.com/gfrey/hdrt/mat"
	"github.com/gfrey/hdrt/obj"
	"github.com/gfrey/hdrt/vec"
)

const MAX_STEPS = 4

type Scene struct {
	AmbientLight *obj.Material
	Objects      []obj.Object `json:"objects"`
	Lights       []*Light     `json:"lights"`
}

func (sc *Scene) UnmarshalJSON(data []byte) error {
	rsc := &struct {
		AmbientLight *obj.Material
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
	return sc.render(pos, dir, 0)
}

func (sc *Scene) render(pos, dir *vec.Vector, step uint) *color.RGBA {
	if step > MAX_STEPS {
		return &color.RGBA{}
	}
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
		cDirect := sc.DirectLight(cand, ipos, pos)
		cReflect := sc.Reflect(cand, ipos, pos, step)
		cRefract := sc.Refract(cand, ipos, pos)
		return &color.RGBA{
			sumColorChans(cDirect.R, cReflect.R, cRefract.R),
			sumColorChans(cDirect.G, cReflect.G, cRefract.G),
			sumColorChans(cDirect.B, cReflect.B, cRefract.B),
			255,
		}

	}
	return &color.RGBA{}
}

func sumColorChans(a, b, c uint8) uint8 {
	r := uint16(a) + uint16(b) + uint16(c)
	if r > 255 {
		return 255
	}
	return uint8(r)
}

func (sc *Scene) Reflect(o obj.Object, ipos, spos *vec.Vector, step uint) *color.RGBA {
	n := o.Normal(ipos)
	i := vec.Sub(ipos, spos)
	r := vec.Sub(i, vec.ScalarMultiply(n, 2*vec.Dot(i, n)))

	col := sc.render(ipos, r, step+1)
	delta := 255 / o.Reflection()
	col.R /= delta
	col.G /= delta
	col.B /= delta
	return col
}

func (sc *Scene) Refract(o obj.Object, ipos, spos *vec.Vector) *color.RGBA {
	return &color.RGBA{}
}

func (sc *Scene) DirectLight(o obj.Object, ipos, spos *vec.Vector) *color.RGBA {
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
			ln := vec.Dot(dir, normal)
			// is obj "visible" from light source?
			for j := range sc.Objects {
				if sc.Objects[j] == o {
					continue
				}

				tmpPos := sc.Objects[j].Intersect(lPos, vec.ScalarMultiply(dir, -1))
				if tmpPos == nil {
					continue
				}

				d := vec.Sub(tmpPos, lPos).Length()
				if mat.FloatLessThan(d, dist) { // in shadow
					continue LIGHTSOURCES
				}
			}

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
				rn := math.Pow(vec.Dot(rm, v), float64(o.Reflection())) * lcoeff
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
