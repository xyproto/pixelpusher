package main

import (
	"github.com/xyproto/pixelpusher"
)

func onDraw(gfx *pixelpusher.Config) error {
	// x, y, r, g, b
	return pixelpusher.Plot(gfx, 0, 0, 255, 0, 0)
}

func main() {
	gfx := pixelpusher.New("Red Pixel")
	gfx.Run(onDraw, nil, nil, nil)
}
