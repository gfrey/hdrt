package hdrt

import (
	"encoding/json"
	"math"

	"github.com/gfrey/hdrt/mat"
	"github.com/gfrey/hdrt/quat"
	"github.com/gfrey/hdrt/vec"
)

type Camera struct {
	Orientation *quat.Quaternion
	Position    *vec.Vector
	FOV         float64 // angle in degree

	direction *vec.Vector
	up        *vec.Vector
}

func (c *Camera) UnmarshalJSON(data []byte) error {
	d := &struct {
		Orientation [][4]float64
		Position    *vec.Vector
		FOV         float64
	}{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	c.Orientation = quat.NewRotation(vec.New(1, 0, 0), 0)
	for _, o := range d.Orientation {
		r := quat.NewRotation(vec.New(o[0], o[1], o[2]), o[3])
		c.Orientation.Multiply(r)
	}
	c.Position = d.Position
	c.FOV = d.FOV
	return nil
}

func (c *Camera) Init() error {
	c.direction = vec.New(1, 0, 0)
	c.up = vec.New(0, 1, 0)
	if c.Orientation != nil {
		c.direction = c.Orientation.Rotate(c.direction)
		c.up = c.Orientation.Rotate(c.up)
	}
	return nil
}

type Viewplane struct {
	Distance   float64
	ResX, ResY int

	pos, a, b *vec.Vector
}

func (v *Viewplane) Init(c *Camera) error {
	v.span(c)
	return nil
}

func (v *Viewplane) span(c *Camera) {
	vpCenter := vec.Add(c.Position, vec.ScalarMultiply(c.direction, v.Distance))

	aspectRatio := float64(v.ResX) / float64(v.ResY)
	alpha := mat.Deg2Rad(c.FOV)
	beta := alpha / aspectRatio

	a := v.Distance * math.Tan(alpha/2.0)
	b := v.Distance * math.Tan(beta/2.0)

	vpTop := c.up
	vpSide := vec.Cross(c.direction, c.up)

	v.a = vec.ScalarMultiply(vpSide, a)
	v.b = vec.ScalarMultiply(vpTop, -b)
	v.pos = vec.Add(vec.Add(vpCenter, vec.ScalarMultiply(v.a, -0.5)), vec.ScalarMultiply(v.b, -0.5))
}

func (v *Viewplane) PixelPosition(x, y int) *vec.Vector {
	posX := vec.Add(v.pos, vec.ScalarMultiply(v.a, float64(x)/float64(v.ResX)))
	return vec.Add(posX, vec.ScalarMultiply(v.b, float64(y)/float64(v.ResY)))
}
