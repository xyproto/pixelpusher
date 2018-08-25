package multirender

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

// BlitImage blits an image on top of a pixel buffer
func BlitImage(pixels []uint32, pitch int32, img *image.RGBA) error {
	width := pitch
	height := int32(len(pixels)) / pitch

	rectWidth := int32(img.Rect.Size().X)
	rectHeight := int32(img.Rect.Size().Y)

	if rectWidth < width || rectHeight < height {
		return fmt.Errorf("Invalid size (%d, %d) for blitting on pixel buffer of size (%d, %d)", rectWidth, rectHeight, width, height)
	}

	// Loop through target coordinates
	for y := int32(0); y < height; y++ {
		for x := int32(0); x < width; x++ {
			c := img.At(int(x), int(y)).(color.RGBA)
			Pixel(pixels, x, y, c, pitch)
		}
	}

	return nil
}
