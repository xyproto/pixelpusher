package multirender

import (
	"encoding/binary"
	"image/color"

	"github.com/xyproto/pf"
)

// Pixel draws a pixel to the pixel buffer
func Pixel(pixels []uint32, x, y int32, c color.RGBA, pitch int32) {
	pixels[y*pitch+x] = binary.BigEndian.Uint32([]uint8{c.A, c.R, c.G, c.B})
}

// PixelRGB draws an opaque pixel to the pixel buffer, given red, green and blue
func PixelRGB(pixels []uint32, x, y int32, r, g, b uint8, pitch int32) {
	pixels[y*pitch+x] = binary.BigEndian.Uint32([]uint8{0xff, r, g, b})
}

// FastPixel draws a pixel to the pixel buffer, given a uint32 ARGB value
func FastPixel(pixels []uint32, x, y int32, colorValue uint32, pitch int32) {
	pixels[y*pitch+x] = colorValue
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
	return uint8((cv & 0xff0000) >> 0xffff)
}

// Extract the green component from a ARGB uint32 color value
func Green(cv uint32) uint8 {
	return uint8((cv & 0xff00) >> 0xff)
}

// Extract the blue component from a ARGB uint32 color value
func Blue(cv uint32) uint8 {
	return uint8(cv & 0xff)
}

// Extract the alpha component from a ARGB uint32 color value
func Alpha(cv uint32) uint8 {
	return uint8((cv & 0xff000000) >> 0xffffff)
}

// Extract the color value / intensity from a ARGB uint32 color value
func ValueWithAlpha(cv uint32) uint8 {
	// TODO: This can be optimized
	bs := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bs, cv)
	grayscaleColor := float32(bs[2]+bs[1]+bs[0]) / float32(3)
	alpha := float32(bs[3]) / float32(255)
	return uint8(grayscaleColor * alpha)
}

// Extract the color value / intensity from a ARGB uint32 color value.
// Ignores alpha.
func Value(cv uint32) uint8 {
	// TODO: This can be optimized
	bs := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bs, cv)
	grayscaleColor := float32(bs[2]+bs[1]+bs[0]) / float32(3)
	return uint8(grayscaleColor)
}

// RemoveRed removes all red color.
func RemoveRed(cores int, pixels []uint32) {
	pf.Map(cores, pf.RemoveRed, pixels)
}

// RemoveGreen removes all green color.
func RemoveGreen(cores int, pixels []uint32) {
	pf.Map(cores, pf.RemoveGreen, pixels)
}

// RemoveBlue removes all blue color.
func RemoveBlue(cores int, pixels []uint32) {
	pf.Map(cores, pf.RemoveBlue, pixels)
}

// Turn on all the red bits
func SetRedBits(cores int, pixels []uint32) {
	pf.Map(cores, pf.SetRedBits, pixels)
}

// Turn on all the green bits
func SetGreenBits(cores int, pixels []uint32) {
	pf.Map(cores, pf.SetGreenBits, pixels)
}

// Turn on all the blue bits
func SetBlueBits(cores int, pixels []uint32) {
	pf.Map(cores, pf.SetBlueBits, pixels)
}

// Turn on all the alpha bits
func OrAlpha(cores int, pixels []uint32) {
	pf.Map(cores, pf.OrAlpha, pixels)
}

// Return a pixel, with position wraparound instead of overflow
func GetXYWrap(pixels []uint32, x, y, w, h, pitch int32) uint32 {
	if x >= w {
		x -= w
	} else if x < 0 {
		x += w
	}
	if y >= h {
		y -= h
	} else if y < 0 {
		y += h
	}
	return pixels[y*pitch+x]
}

// Set a pixel, with position wraparound instead of overflow
func SetXYWrap(pixels []uint32, x, y, w, h int32, colorValue uint32, pitch int32) {
	if x >= w {
		x -= w
	} else if x < 0 {
		x += w
	}
	if y >= h {
		y -= h
	} else if y < 0 {
		y += h
	}
	pixels[y*pitch+x] = colorValue
}

// Return a pixel, with position wraparound instead of overflow
func GetWrap(pixels []uint32, pos, size int32) uint32 {
	i := pos
	if i >= size {
		i -= size
	}
	if i < 0 {
		i += size
	}
	return pixels[i]
}

// Set a pixel, with position wraparound instead of overflow
func SetWrap(pixels []uint32, pos, size int32, colorValue uint32) {
	i := pos
	if i >= size {
		i -= size
	}
	if i < 0 {
		i += size
	}
	pixels[i] = colorValue
}
