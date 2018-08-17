package multirender

import (
	"encoding/binary"
	"image/color"
)

// Line draws a line to the pixel buffer.
// pixels are the pixels, pitch is the width of the pixel buffer.
func Line(pixels []uint32, x1, y1, x2, y2 int32, c color.RGBA, pitch int32) {
	if y1 == y2 {
		HorizontalLine(pixels, y1, x1, x2, c, pitch)
		return
	}
	if x1 == x2 {
		VerticalLine(pixels, x1, y1, y2, c, pitch)
		return
	}

	startx, stopx := MinMax(x1, x2)
	starty, stopy := MinMax(y1, y2)

	xdiff := stopx - startx
	ydiff := stopy - starty

	colorValue := binary.BigEndian.Uint32([]uint8{c.A, c.R, c.G, c.B})

	if xdiff > ydiff {
		// We're going along X
		y := float32(starty)
		ystep := float32(ydiff) / float32(xdiff)
		if y1 != starty {
			// Move in the other direction along Y
			ystep = -ystep
			y = float32(stopy)
		}
		// Draw the line
		for x := startx; x < stopx; x++ {
			pixels[int32(y)*pitch+int32(x)] = colorValue
			y += ystep
		}
	} else {
		// We're going along Y
		x := float32(startx)
		xstep := float32(xdiff) / float32(ydiff)
		if x1 != startx {
			// Move in the other direction along X
			xstep = -xstep
			x = float32(stopx)
		}
		// Draw the line
		starty *= pitch
		stopy *= pitch
		for y := starty; y < stopy; y += pitch {
			pixels[y+int32(x)] = colorValue
			x += xstep
		}
	}
}

// HorizontalLineFast draws a line from (x1, y) to (x2, y), but x1 must be smaller than x2!
func HorizontalLineFast(pixels []uint32, y, x1, x2 int32, c color.RGBA, pitch int32) {
	colorValue := binary.BigEndian.Uint32([]uint8{c.A, c.R, c.G, c.B})
	xstart, xstop := x1, x2
	offset := y * pitch
	xstart += offset
	xstop += offset
	for x := xstart; x < xstop; x++ {
		pixels[x] = colorValue
	}
}

// HorizontalLine draws a line from (x1, y) to (x2, y)
func HorizontalLine(pixels []uint32, y, x1, x2 int32, c color.RGBA, pitch int32) {
	if x1 < x2 {
		HorizontalLineFast(pixels, y, x1, x2, c, pitch)
	} else {
		HorizontalLineFast(pixels, y, x2, x1, c, pitch)
	}
}

// VerticalLineFast draws a line from (x, y1) to (x, y2), but y1 must be smaller than y2!
func VerticalLineFast(pixels []uint32, x, y1, y2 int32, c color.RGBA, pitch int32) {
	colorValue := binary.BigEndian.Uint32([]uint8{c.A, c.R, c.G, c.B})
	for y := y1; y < y2; y += pitch {
		pixels[y+x] = colorValue
	}
}

// VerticalLine draws a line from (x, y1) to (x, y2)
func VerticalLine(pixels []uint32, x, y1, y2 int32, c color.RGBA, pitch int32) {
	if y1 < y2 {
		VerticalLineFast(pixels, x, y1, y2, c, pitch)
	} else {
		VerticalLineFast(pixels, x, y2, y1, c, pitch)
	}
}
