package sdl2utils

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/xyproto/pixelpusher"
	"image"
	"unsafe"
)

// RendererToImage converts a *sdl.Renderer to an *image.RGBA image
func RendererToImage(renderer *sdl.Renderer) (*image.RGBA, error) {
	w, h, err := renderer.GetOutputSize()
	if err != nil {
		return nil, err
	}

	rect := &sdl.Rect{0, 0, w, h}
	format := uint32(sdl.PIXELFORMAT_ARGB8888)
	pitch := int32(w)
	pixelBuffer := make([]uint32, w*h)
	pixelBufferPointer := unsafe.Pointer(&pixelBuffer[0])

	// pitch * size of uint32, in bytes
	if err := renderer.ReadPixels(rect, format, pixelBufferPointer, int(pitch*4)); err != nil {
		return nil, err
	}

	// Convert the extracted pixel buffer to an image
	return pixelpusher.PixelsToImage(pixelBuffer, pitch), nil
}

// Screenshot saves the contents of the given *sdl.Renderer to a PNG file.
// Set overwrite to true for overwriting any existing files.
func Screenshot(renderer *sdl.Renderer, filename string, overwrite bool) error {
	img, err := RendererToImage(renderer)
	if err != nil {
		return err
	}
	return pixelpusher.SaveImageToPNG(img, filename, overwrite)
}
