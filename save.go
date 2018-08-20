package multirender

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
)

// exists checks if a file already exists
func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// PixelsToImage converts a pixel buffer to an image.RGBA image
func PixelsToImage(pixels []uint32, pitch uint32) *image.RGBA {
	width := pitch
	height := uint32(len(pixels)) / pitch

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	bs := make([]uint8, 4)
	for y := uint32(0); y < height; y++ {
		for x := uint32(0); x < width; x++ {
			binary.LittleEndian.PutUint32(bs, pixels[y*pitch+x])
			c := color.RGBA{bs[2], bs[1], bs[0], bs[3]}
			img.Set(int(x), int(y), c)
			//img.Pix[y*pitch+x*4] = bs[2]
			//img.Pix[y*pitch+x*4+1] = bs[1]
			//img.Pix[y*pitch+x*4+2] = bs[0]
			//img.Pix[y*pitch+x*4+3] = bs[3]
		}
	}

	return img
}

// SaveImageToPNG saves an image.RGBA image to a PNG file.
// Set overwrite to true to allow overwriting files.
func SaveImageToPNG(img *image.RGBA, filename string, overwrite bool) error {
	if !overwrite && exists(filename) {
		return errors.New(filename + " already exists")
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

// Save pixels in uint32 ARGB format to PNG with alpha.
// pitch is the width of the pixel buffer.
// Set overwrite to true to allow overwriting files.
func SavePixelsToPNG(pixels []uint32, pitch uint32, filename string, overwrite bool) error {
	return SaveImageToPNG(PixelsToImage(pixels, pitch), filename, overwrite)
}
