package multirender

import (
	"encoding/binary"
	"image/color"
)

// Pixel draws a pixel to the pixel buffer
func Pixel(pixels []uint32, x, y int32, c color.RGBA, pitch int32) {
	pixels[y*pitch+x] = binary.BigEndian.Uint32([]uint8{c.A, c.R, c.G, c.B})
}
