package multirender

import (
	"encoding/binary"
	"image/color"
	"sync"
)

// pointInTriangle tries to decide if the given x and y are within the triangle defined by p0, p1 and p2
// area is 1.0 divided on (the area of the triangle, times 2)
func pointInTriangle(x, y int32, p0, p1, p2 *Pos, areaMod float32) bool {
	s := areaMod * float32(p0.y*p2.x-p0.x*p2.y+(p2.y-p0.y)*x+(p0.x-p2.x)*y)
	t := areaMod * float32(p0.x*p1.y-p0.y*p1.x+(p0.y-p1.y)*x+(p1.x-p0.x)*y)
	return s > 0 && t > 0 && (1-s-t) > 0
}

// Area returns the area of a triangle
func area(p0, p1, p2 *Pos) float32 {
	return 0.5 * float32(-p1.y*p2.x+p0.y*(-p1.x+p2.x)+p0.x*(p1.y-p2.y)+p1.x*p2.y)
}

// drawPartialTriangle draws a part of a triangle
// areaMod is 1.0 divided on (the area of the triangle, times 2)
func drawPartialTriangle(wg *sync.WaitGroup, pixels []uint32, p1, p2, p3 *Pos, minX, maxX, minY, maxY int32, areaMod float32, colorValue uint32, pitch int32) {
	for y := minY; y < maxY; y++ {
		offset := y * pitch
		for x := minX; x < maxX; x++ {
			if pointInTriangle(x, y, p1, p2, p3, areaMod) {
				pixels[offset+x] = colorValue
			}
		}
	}
	wg.Done()
}

// Triangle draws a triangle, concurrently.
// Core is the number of goroutines that will be used.
// pitch is the "width" of the pixel buffer.
func Triangle(cores int, pixels []uint32, x1, y1, x2, y2, x3, y3 int32, c color.RGBA, pitch int32) {
	var wg sync.WaitGroup

	p1 := &Pos{x1, y1}
	p2 := &Pos{x2, y2}
	p3 := &Pos{x3, y3}

	minY, maxY := MinMax3(y1, y2, y3)
	minX, maxX := MinMax3(x1, x2, x3)

	// Triangle area, with modifications, for performance
	areaMod := 1.0 / (area(p1, p2, p3) * 2.0)

	colorValue := binary.BigEndian.Uint32([]uint8{c.A, c.R, c.G, c.B})

	ylength := (maxY - minY)

	if ylength == 0 {
		// This is not a triangle, but a horizontal line, since ylength == 0
		HorizontalLineFast(pixels, minY, minX, maxX, c, pitch)
		return
	}
	// cores-1 because of the final part after the for loop
	ystep := ylength / int32(cores-1)
	if ystep == 0 {
		return
	}
	var minYCore int32
	var maxYCore int32
	for y := minY; y < (maxY - ystep); y += ystep {
		minYCore = y
		maxYCore = y + ystep
		wg.Add(1)
		go drawPartialTriangle(&wg, pixels, p1, p2, p3, minX, maxX, minYCore, maxYCore, areaMod, colorValue, pitch)
	}
	// Draw the final part, if there are a few pixels missing at the end
	if maxYCore < maxY {
		wg.Add(1)
		go drawPartialTriangle(&wg, pixels, p1, p2, p3, minX, maxX, maxYCore, maxY, areaMod, colorValue, pitch)
	}

	wg.Wait()
}
