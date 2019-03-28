package pf

// The uint32 is on the form ARGB

// Invert the colors
func Invert(v uint32) uint32 {
	// Invert the colors, but set the alpha value to 0xff
	return (0xffffffff - v) | 0xff000000
}

// Invert the colors, including the alpha value
func InvertEverything(v uint32) uint32 {
	// Invert everything
	return 0xffffffff - v
}

// Keep the red component
func OnlyRed(v uint32) uint32 {
	// Keep alpha and the red value
	return v & 0xffff0000
}

// Keep the green component
func OnlyGreen(v uint32) uint32 {
	// Keep alpha and the green value
	return v & 0xff00ff00
}

// Keep the blue component
func OnlyBlue(v uint32) uint32 {
	// Keep alpha and the blue value
	return v & 0xff0000ff
}

// Keep the alpha component
func OnlyAlpha(v uint32) uint32 {
	// Keep only the alpha value
	return v & 0xff000000
}

// Remove the red component
func RemoveRed(v uint32) uint32 {
	// Keep everything but red
	return v & 0xff00ffff
}

// Remove the green component
func RemoveGreen(v uint32) uint32 {
	// Keep everything but green
	return v & 0xffff00ff
}

// Remove the blue component
func RemoveBlue(v uint32) uint32 {
	// Keep everything but blue
	return v & 0xffffff00
}

// Remove the alpha component, making the pixels transparent
func RemoveAlpha(v uint32) uint32 {
	// Keep everything but the alpha value
	return v & 0x00ffffff
}

// Make the red component of every pixel 0xff
func SetRedBits(v uint32) uint32 {
	return v | 0x00ff0000
}

// Make the green component of every pixel 0xff
func SetGreenBits(v uint32) uint32 {
	return v | 0x0000ff00
}

// Make the blue component of every pixel 0xff
func SetBlueBits(v uint32) uint32 {
	return v | 0x000000ff
}

// Make the alpha component of every pixel 0xff
func OrAlpha(v uint32) uint32 {
	return v | 0xff000000
}
