# PixelFunctions

Apply functions to each pixel in an image, concurrently.

This module contains functions that fit well together with the [multirender](https://github.com/xyproto/multirender) module.

The `PixelFunction` type has this signature:

    func(v uint32) uint32

If you have a pixel buffer of type `[]uint32`, with colors on the form `ARGB`, then this modules allows you to apply functions of the type `PixelFunction` to that slice, concurrently.

The goal is to avoid looping over all pixels more than once, while applying many different effects, concurrently.

## Combine and Map

* Several `PixelFunction` functions can be combined to a single `PixelFunction` by using the `Combine` function.
* A `PixelFuncion` can be applied to a pixel buffer by using the `Map` function.

Example:

```go
package main

import (
	"fmt"
	"github.com/xyproto/pf"
	"runtime"
)

func main() {
	// Resolution
	const w, h = 320, 200

	pixels := make([]uint32, w*h)

	// Find the number of available CPUs
	n := runtime.NumCPU()

	// Combine two pixel functions
	pfs := pf.Combine(pf.InvertEverything, pf.OnlyBlue)

	// Run the combined pixel functions over all pixels using all available CPUs
	pf.Map(n, pfs, pixels)

	// Retrieve the red, green and blue components of the first pixel
	red := (pixels[0] | 0x00ff0000) >> 0xffff
	green := (pixels[0] | 0x0000ff00) >> 0xff
	blue := (pixels[0] | 0x000000ff)

	// Should output only blue: rgb(0, 0, 255)
	fmt.Printf("rgb(%d, %d, %d)\n", red, green, blue)
}
```

# General info

* License: MIT
* Version: 0.1
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
