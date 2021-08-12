package multirender

import (
	"errors"
	"fmt"
	"image/color"
	"math"
)

type Vec3 struct {
	x float32
	y float32
	z float32
}

// Vertex has a position and a color value
type Vertex struct {
	pos        *Vec3
	colorValue uint32
}

func NewVertex(x, y, z float32, r, g, b, a uint8) *Vertex {
	return &Vertex{&Vec3{x, y, z}, RGBAToColorValue(r, g, b, a)}
}

func (v *Vertex) X() float32 {
	return v.pos.x
}

func (v *Vertex) Y() float32 {
	return v.pos.y
}

func (v *Vertex) Z() float32 {
	return v.pos.z
}

func (v *Vertex) Normalize() error {
	l := float32(math.Sqrt(float64(v.pos.x*v.pos.x + v.pos.y*v.pos.y + v.pos.z*v.pos.z)))
	if l == 0 {
		return errors.New("normalizing a Vertex with length 0")
	}
	v.pos.x /= l
	v.pos.y /= l
	v.pos.z /= l
	return nil
}

func (v *Vertex) Set(x, y, z float32) {
	v.pos.x = x
	v.pos.y = y
	v.pos.z = z
}

func (v *Vertex) Get() (float32, float32, float32) {
	return v.pos.x, v.pos.y, v.pos.z
}

func (v *Vertex) GetVec3() *Vec3 {
	return v.pos
}

func (v *Vertex) SetRGBA(r, g, b, a uint8) {
	v.colorValue = RGBAToColorValue(r, g, b, a)
}

func (v *Vertex) SetColor(c color.RGBA) {
	v.colorValue = RGBAToColorValue(c.R, c.G, c.B, c.A)
}

func (v *Vertex) GetRGBA() (uint8, uint8, uint8, uint8) {
	return ColorValueToRGBA(v.colorValue)
}

func (v *Vertex) GetColor() color.RGBA {
	r, g, b, a := ColorValueToRGBA(v.colorValue)
	return color.RGBA{r, g, b, a}
}

func (v *Vertex) R() uint8 {
	r, _, _, _ := v.GetRGBA()
	return r
}

func (v *Vertex) G() uint8 {
	_, g, _, _ := v.GetRGBA()
	return g
}

func (v *Vertex) B() uint8 {
	_, _, b, _ := v.GetRGBA()
	return b
}

func (v *Vertex) A() uint8 {
	_, _, _, a := v.GetRGBA()
	return a
}

func (v *Vertex) String() string {
	r, g, b, a := v.GetRGBA()
	return fmt.Sprintf("v(%v, %v, %v) color(%v, %v, %v, %v)", v.pos.x, v.pos.y, v.pos.z, r, g, b, a)
}
