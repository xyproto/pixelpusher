package multirender

import (
	"errors"
	"image"
	"image/png"
	"os"
)

// exists checks if a file already exists
func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
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
func SavePixelsToPNG(pixels []uint32, pitch int32, filename string, overwrite bool) error {
	return SaveImageToPNG(PixelsToImage(pixels, pitch), filename, overwrite)
}
