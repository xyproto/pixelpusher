package main

import (
	pp "github.com/xyproto/pixelpusher"
)

func onDraw(canvas *pp.Canvas) error {
	// x=0, y=0, red=255, green=0, blue=0
	return pp.Plot(canvas, 0, 0, 255, 0, 0)
}

func main() {
	// The window title is "Red Pixel"
	canvas := pp.New("Red Pixel")
	// onDraw will be called whenever it is time to draw a frame
	canvas.Run(onDraw, nil, nil, nil)
}
