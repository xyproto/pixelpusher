package pixelpusher

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
)

// x, y and pitch should be int32, since they are likely to be used together with other int32 types.
// go-sdl2 uses int32 for most things.

// PixelsToImage converts a pixel buffer to an image.RGBA image
func PixelsToImage(pixels []uint32, pitch int32) *image.RGBA {
	width := pitch
	height := int32(len(pixels)) / pitch

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	bs := make([]uint8, 4)
	for y := int32(0); y < height; y++ {
		for x := int32(0); x < width; x++ {
			binary.LittleEndian.PutUint32(bs, pixels[y*pitch+x])
			c := color.RGBA{bs[2], bs[1], bs[0], bs[3]}
			img.Set(int(x), int(y), c)
		}
	}

	return img
}

// BlitImage blits an image on top of a pixel buffer, while blending
func BlitImage(pixels []uint32, pitch int32, img *image.RGBA) error {
	width := pitch
	height := int32(len(pixels)) / pitch

	rectWidth := int32(img.Rect.Size().X)
	rectHeight := int32(img.Rect.Size().Y)

	if rectWidth < width || rectHeight < height {
		return fmt.Errorf("invalid size (%d, %d) for blitting on pixel buffer of size (%d, %d)", rectWidth, rectHeight, width, height)
	}

	// Loop through target coordinates
	for y := int32(0); y < height; y++ {
		offset := y * pitch
		for x := int32(0); x < width; x++ {
			cv := ColorToColorValue(img.At(int(x), int(y)).(color.RGBA))
			pixels[offset+x] = Blend(pixels[offset+x], cv)
		}
	}

	return nil
}

// BlitImageOnTop blits an image on top of a pixel buffer, disregarding any previous pixels.
// The resulting pixels are opaque.
func BlitImageOnTop(pixels []uint32, pitch int32, img *image.RGBA) error {
	width := pitch
	height := int32(len(pixels)) / pitch

	rectWidth := int32(img.Rect.Size().X)
	rectHeight := int32(img.Rect.Size().Y)

	if rectWidth < width || rectHeight < height {
		return fmt.Errorf("invalid size (%d, %d) for blitting on pixel buffer of size (%d, %d)", rectWidth, rectHeight, width, height)
	}

	// Loop through target coordinates
	for y := int32(0); y < height; y++ {
		offset := y * pitch
		for x := int32(0); x < width; x++ {
			cv := ColorToColorValue(img.At(int(x), int(y)).(color.RGBA))
			pixels[offset+x] = Add(pixels[offset+x], cv)
		}
	}

	return nil
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
