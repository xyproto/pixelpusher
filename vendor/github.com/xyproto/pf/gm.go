package pf

import (
	"sync"
)

// GlitchyMap can map a PixelFunction over every pixel (uint32 ARGB value).
// This function has data race issues and should not be
// used for anything but creating glitch effects on purpose.
func GlitchyMap(cores int, f PixelFunction, pixels []uint32) {
	// Map a pixel function over every pixel, concurrently
	var (
		wg      sync.WaitGroup
		iLength = int32(len(pixels))
		iStep   = iLength / int32(cores)

		// iConcurrentlyDone keeps track of how much work have been done by launching goroutines
		iConcurrentlyDone = int32(cores) * iStep

		// iDone keeps track of how much work have been done in total
		iDone int32
	)

	// Apply partialMap for each of the partitions
	if iStep < iLength {
		for i := int32(0); i < iConcurrentlyDone; i += iStep {
			// run a PixelFunction on parts of the pixel buffer
			wg.Add(1)
			go func(wg *sync.WaitGroup, f PixelFunction, pixels []uint32, iStart, iStop int32) {
				for i := iStart; i < iStop; i++ {
					pixels[i] = f(pixels[i])
				}
				wg.Done()
			}(&wg, f, pixels, i, i+iStep)
		}
		iDone = iConcurrentlyDone
	}

	if iDone == iLength {
		// No leftover pixels
		return
	}

	// Apply partialMap to the final leftover pixels
	wg.Add(1)
	go func(wg *sync.WaitGroup, f PixelFunction, pixels []uint32, iStart, iStop int32) {
		for i := iStart; i < iStop; i++ {
			pixels[i] = f(pixels[i])
		}
		wg.Done()
	}(&wg, f, pixels, iDone, iLength)

	wg.Wait()
}
