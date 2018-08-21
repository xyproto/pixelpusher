package multirender

import (
	"encoding/binary"
	"image/color"
)

// Pixel draws a pixel to the pixel buffer
func Pixel(pixels []uint32, x, y int32, c color.RGBA, pitch int32) {
	pixels[y*pitch+x] = binary.BigEndian.Uint32([]uint8{c.A, c.R, c.G, c.B})
}

// Clear changes all pixels to the given color
func Clear(pixels []uint32, c color.RGBA) {
	colorValue := binary.BigEndian.Uint32([]uint8{c.A, c.R, c.G, c.B})
	for i := range pixels {
		pixels[i] = colorValue
	}
}

// FastClear changes all pixels to the given uint32 color value,
// like 0xff0000ff for: 0xff red, 0x00 green, 0x00 blue and 0xff alpha.
func FastClear(pixels []uint32, colorValue uint32) {
	for i := range pixels {
		pixels[i] = colorValue
	}
}

// RGBAToColorValue converts from four bytes to an ARGB uint32 color value
func RGBAToColorValue(r, g, b, a uint8) uint32 {
	return binary.BigEndian.Uint32([]uint8{a, r, g, b})
}

// ColorValueToRGBA converts from an ARGB uint32 color value to four bytes
func ColorValueToRGBA(cv uint32) (uint8, uint8, uint8, uint8) {
	bs := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bs, cv)
	// r, g, b, a
	return bs[2], bs[1], bs[0], bs[3]
}

// Extract the red component from a ARGB uint32 color value
func Red(cv uint32) uint8 {
	// TODO: This can be optimized
	bs := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bs, cv)
	return bs[2]
}

// Extract the green component from a ARGB uint32 color value
func Green(cv uint32) uint8 {
	// TODO: This can be optimized
	bs := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bs, cv)
	return bs[1]
}

// Extract the blue component from a ARGB uint32 color value
func Blue(cv uint32) uint8 {
	// TODO: This can be optimized
	bs := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bs, cv)
	return bs[0]
}

// Extract the alpha component from a ARGB uint32 color value
func Alpha(cv uint32) uint8 {
	// TODO: This can be optimized
	bs := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bs, cv)
	return bs[3]
}

// TODO: Also create a Value function that disregards the alpha value

// Extract the color value / intensity from a ARGB uint32 color value
func Value(cv uint32) uint8 {
	// TODO: This can be optimized
	bs := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bs, cv)
	grayscaleColor := float32(bs[2]+bs[1]+bs[0]) / float32(3)
	alpha := float32(bs[3]) / float32(255)
	return uint8(grayscaleColor * alpha)
}
