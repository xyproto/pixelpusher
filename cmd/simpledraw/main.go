package main

import (
	"errors"

	pp "github.com/xyproto/pixelpusher"
)

var x, y = 160, 100

func onDraw(canvas *pp.Canvas) error {
	return pp.Plot(canvas, x, y, 255, 0, 0)
}

func onPress(left, right, up, down, space, enter, esc bool) error {
	if up {
		y--
	} else if down {
		y++
	}
	if left {
		x--
	} else if right {
		x++
	}
	if esc {
		return errors.New("quit")
	}
	return nil
}

func main() {
	pp.New("Simple Draw").Run(onDraw, onPress, nil, nil)
}
