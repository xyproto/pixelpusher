package pf

import (
	"sync"
)

// Perform an operation on a single ARGB pixel
type PixelFunction func(v uint32) uint32

// Combine two functions to a single PixelFunction.
// The functions are applied in the same order as the arguments.
func Combine(a, b PixelFunction) PixelFunction {
	return func(v uint32) uint32 {
		return b(a(v))
	}
}

// Combine three functions to a single PixelFunction.
// The functions are applied in the same order as the arguments.
func Combine3(a, b, c PixelFunction) PixelFunction {
	return Combine(a, Combine(b, c))
}

// Divide a slice of pixels into several slices
func Divide(pixels []uint32, n int) [][]uint32 {
	length := len(pixels)

	sliceLen := length / n
	leftover := length % n

	var sliceOfSlices [][]uint32
	for i := 0; i < (length - leftover); i += sliceLen {
		sliceOfSlices = append(sliceOfSlices, pixels[i:i+sliceLen])
	}
	if leftover > 0 {
		sliceOfSlices = append(sliceOfSlices, pixels[length-leftover:length])
	}
	return sliceOfSlices
}

func Map(cores int, f PixelFunction, pixels []uint32) {
	wg := &sync.WaitGroup{}

	// First copy the pixels into several separate slices
	sliceOfSlices := Divide(pixels, cores)

	// Then process the slices individually
	wg.Add(len(sliceOfSlices))
	for _, subPixels := range sliceOfSlices {
		// subPixels is a slice of pixels
		go func(wg *sync.WaitGroup, subPixels []uint32) {
			for i := range subPixels {
				subPixels[i] = f(subPixels[i])
			}
			wg.Done()
		}(wg, subPixels)
	}
	wg.Wait()

	// Then combine the slices into a new and better slice
	newPixels := make([]uint32, len(pixels))
	for _, subPixels := range sliceOfSlices {
		newPixels = append(newPixels, subPixels...)
	}

	// Finally, replace the pixels with the processed pixels
	pixels = newPixels
}
