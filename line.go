package multirender

import (
	"encoding/binary"
	"image/color"
)

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

// Line draws a line in a completely wrong way to the pixel buffer.
// pixels are the pixels, pitch is the width of the pixel buffer.
func Line(pixels []uint32, x1, y1, x2, y2 int32, c color.RGBA, pitch int32) {
	//fmt.Printf("Line from (%d, %d) to (%d, %d)\n", x1, y1, x2, y2)
	if y1 == y2 {
		HorizontalLine(pixels, y1, x1, x2, c, pitch)
		return
	}
	if x1 == x2 {
		VerticalLine(pixels, x1, y1, y2, c, pitch)
		return
	}

	// First figure out if it's an X or Y major triangle

	xdiff := Abs(x1 - x2)
	ydiff := Abs(y1 - y2)

	colorValue := binary.BigEndian.Uint32([]uint8{c.A, c.R, c.G, c.B})

	if xdiff > ydiff {
		//fmt.Println("X MAJOR")
		startx := x1
		starty := y1
		stopx := x2
		stopy := y2
		if x2 < x1 {
			startx = x2
			starty = y2
			stopx = x1
			stopy = y1
		}

		// We're going along X
		y := float32(starty)
		ystep := float32(ydiff) / float32(xdiff)
		if stopy < starty {
			// Move in the other direction along Y
			ystep = -ystep
		}
		// Draw the line
		for x := startx; x < stopx; x++ {
			//fmt.Printf("\t(%d, %d)\n", int32(x), int32(y))
			pixels[int32(y)*pitch+int32(x)] = colorValue
			y += ystep
		}
	} else {
		//fmt.Println("Y MAJOR")
		startx := x1
		starty := y1
		stopx := x2
		stopy := y2
		if y2 < y1 {
			startx = x2
			starty = y2
			stopx = x1
			stopy = y1
		}

		// We're going along Y
		x := float32(startx)
		xstep := float32(xdiff) / float32(ydiff)
		if stopx < startx {
			// Move in the other direction along X
			xstep = -xstep
		}
		// Draw the line
		starty *= pitch
		stopy *= pitch
		for y := starty; y < stopy; y += pitch {
			//fmt.Printf("\t(%d, %d)\n", int32(x), int32(y/pitch))
			pixels[y+int32(x)] = colorValue
			x += xstep
		}
	}
}

//func plot(pixels []uint32, x, y int32, brightness float32, c color.RGBA, pitch int32) {
//	b := uint8(brightness * 255.0)
//	// Use the brightness to scale the given alpha value
//	//colorValue := binary.BigEndian.Uint32([]uint8{uint8(float32(c.A) * brightness * 255.0), c.R, c.G, c.B})
//	colorValue := binary.BigEndian.Uint32([]uint8{0xff, b, b, b})
//	fmt.Println("PLOTTING AT", x, ",", y)
//	pixels[y*pitch+x] = colorValue
//}
//
//// integer part of x
//func ipart(x float32) float32 {
//	return float32(math.Floor(float64(x)))
//}
//
//// round x
//func round(x float32) float32 {
//	return ipart(x + 0.5)
//}
//
//// fractional part of x
//func fpart(x float32) float32 {
//	return x - float32(math.Floor(float64(x)))
//}
//
//// 1 - (fractional part of x)
//func rfpart(x float32) float32 {
//	return 1 - fpart(x)
//}
//
//func abs(a float32) float32 {
//	if a >= 0 {
//		return a
//	}
//	return -a
//}
//
//// Line draws an antialiased line using Xiaolin Wu's algorithm
//// https://en.wikipedia.org/wiki/Xiaolin_Wu%27s_line_algorithm
//func ALine(pixels []uint32, ox0, oy0, ox1, oy1 int32, c color.RGBA, pitch int32) {
//
//	x0 := float32(ox0)
//	y0 := float32(ox0)
//	x1 := float32(ox1)
//	y1 := float32(ox1)
//
//	steep := abs(y1-y0) > abs(x1-x0)
//
//	if steep {
//		x0, y0 = y0, x0
//		x1, y1 = y1, x1
//	}
//	if x0 > x1 {
//		x0, x1 = x1, x0
//		y0, y1 = y1, y0
//	}
//
//	dx := x1 - x0
//	dy := y1 - y0
//
//	gradient := dy / dx
//	if dx == 0.0 {
//		gradient = 1.0
//	}
//
//	// handle first endpoint
//	xend := round(x0)
//	yend := y0 + gradient*(xend-x0)
//
//	xgap := rfpart(float32(x0) + 0.5)
//	xpxl1 := xend // this will be used in the main loop
//	ypxl1 := ipart(float32(yend))
//	if steep {
//		plot(pixels, int32(ypxl1), int32(xpxl1), rfpart(yend)*xgap, c, pitch)
//		plot(pixels, int32(ypxl1+1), int32(xpxl1), fpart(yend)*xgap, c, pitch)
//	} else {
//		plot(pixels, int32(xpxl1), int32(ypxl1), rfpart(yend)*xgap, c, pitch)
//		plot(pixels, int32(xpxl1), int32(ypxl1+1), fpart(yend)*xgap, c, pitch)
//	}
//	intery := yend + gradient // first y-intersection for the main loop
//
//	// handle second endpoint
//	xend = round(x1)
//	yend = y1 + gradient*(xend-x1)
//	xgap = fpart(x1 + 0.5)
//	xpxl2 := xend //this will be used in the main loop
//	ypxl2 := ipart(yend)
//	if steep {
//		plot(pixels, int32(ypxl2), int32(xpxl2), rfpart(yend)*xgap, c, pitch)
//		plot(pixels, int32(ypxl2+1), int32(xpxl2), fpart(yend)*xgap, c, pitch)
//	} else {
//		plot(pixels, int32(xpxl2), int32(ypxl2), rfpart(yend)*xgap, c, pitch)
//		plot(pixels, int32(xpxl2), int32(ypxl2+1), fpart(yend)*xgap, c, pitch)
//	}
//
//	// main loop
//	if steep {
//		for x := xpxl1 + 1; x < xpxl2-1; x++ {
//			plot(pixels, int32(ipart(intery)), int32(x), rfpart(intery), c, pitch)
//			plot(pixels, int32(ipart(intery)+1), int32(x), fpart(intery), c, pitch)
//			intery = intery + gradient
//		}
//	} else {
//		for x := xpxl1 + 1; x < xpxl2-1; x++ {
//			plot(pixels, int32(x), int32(ipart(intery)), rfpart(intery), c, pitch)
//			plot(pixels, int32(x), int32(ipart(intery)+1), fpart(intery), c, pitch)
//			intery = intery + gradient
//		}
//	}
//}
