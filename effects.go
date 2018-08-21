package multirender

import "fmt"

// TODO: Also create a concurrent version of StretchConstrast
// TODO: Also create a version of StretchContrast that disregards alpha.
// TODO: Also create a version of StretchContrast that uses the colorvalue directly.

func wipStretchContrast(cores int, pixels []uint32, pitch int32, from, to uint8) {

	// Find minimum and maximum color intensity, then stretch the colors out from 0 to 255

	// Assume the image has at least one pixel, for the purpose of initializing the min/max values.
	// This will fail if the image has zero pixels!
	tmpV := Value(pixels[0])
	minV, maxV := tmpV, tmpV

	// Find the minimum and maximum for each color value
	for i := range pixels {
		tmpV = Value(pixels[i])
		minV, maxV = MinMax3Byte(minV, maxV, tmpV)
	}

	if minV == 0 && maxV == 255 {
		// Nothing to do here
		return
	}

	fmt.Println("MINV", minV)
	fmt.Println("MAXV", maxV)

	// 0 to 84

	// Scale the color of each pixel
	var r, g, b, a uint8
	for i := range pixels {
		r, g, b, a = ColorValueToRGBA(pixels[i])

        v := Value(pixels[i])

		// First scale the Value for this pixel to be
		// on the scale from 0 to 255 instead of on
		// the scale from minV to maxV.
        vScaled := Scale(int32(v), int32(minV), int32(maxV), 0, 255)

		fmt.Println("VSCALED", i, vScaled)

		// vScaled is now on the scale from minV to maxV
		// now figure out what v needs to be multiplied with
		// to get vScaled

		// v * x = vScaled
		// x = vScaled / v
        x := float32(vScaled) / float32(v)

		// Now use x to multiply r, g and b and see what happens

		pixels[i] = RGBAToColorValue(
			uint8(Clamp(int32(float32(r) * x), 0, 255)),
			uint8(Clamp(int32(float32(g) * x), 0, 255)),
			uint8(Clamp(int32(float32(b) * x), 0, 255)),
			uint8(Clamp(int32(float32(a) * x), 0, 255)),
		)
	}
}
