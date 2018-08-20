# multirender [![Build Status](https://travis-ci.org/xyproto/multirender.svg?branch=master)](https://travis-ci.org/xyproto/multirender) [![GoDoc](https://godoc.org/github.com/xyproto/multirender?status.svg)](http://godoc.org/github.com/xyproto/multirender) [![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/xyproto/multirender/master/LICENSE) [![Report Card](https://img.shields.io/badge/go_report-A+-brightgreen.svg?style=flat)](http://goreportcard.com/report/xyproto/multirender)

Concurrent software rendering and triangle rasterization.

![screencap](img/screencap.gif)

## Features and limitations

* Can draw software-rendered triangles concurrently, using goroutines. The work of drawing the triangles are divided on the available CPU cores.
* Provides flat-shaded triangles.
* Everything is drawn to a `[]uint32` pixel buffer (containing "red", "green", "blue" and "alpha").
* Tested together with SDL2, but can be used with any graphics library that can output pixels from a pixel buffer.
* Everything you need for creating an oldschool game or demoscene demo that will run on Linux, macOS and Windows, while using all of the cores.
* Does not support palette cycling, yet.

## Example, using multirender and [SDL2](https://github.com/veandco/go-sdl2):

```go
package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/xyproto/multirender"
)

const (
	// Size of "worldspace pixels", measured in "screenspace pixels"
	pixelscale = 4

	// The resolution (worldspace)
	width  = 320
	height = 200

	// The width of the pixel buffer, used when calculating where to place pixels (y*pitch+x)
	pitch = width

	// Target framerate
	frameRate = 60

	// Alpha value for opaque colors
	opaque = 255
)

// rb returns a random byte
func rb() uint8 {
	return uint8(rand.Intn(255))
}

// rw returns a random int32 in the range [0,width)
func rw() int32 {
	return rand.Int31n(width)
}

// rh returns a random int32 in the range [0,height)
func rh() int32 {
	return rand.Int31n(height)
}

// DrawAll fills the pixel buffer with pixels.
// "cores" is how many CPU cores should be targeted when drawing triangles,
// by launching the same number of goroutines.
func DrawAll(pixels []uint32, cores int) {

	// Draw a triangle, concurrently
	multirender.Triangle(cores, pixels, rw(), rh(), rw(), rh(), rw(), rh(), color.RGBA{rb(), rb(), rb(), opaque}, pitch)

	// Draw a line and a red pixel, without caring about which order they appear in, or if they will complete before the next frame is drawn
	go multirender.Line(pixels, rw(), rh(), rw(), rh(), color.RGBA{0xff, 0xff, 0, opaque}, pitch)
	go multirender.Pixel(pixels, rw(), rh(), color.RGBA{0xff, 0xff, 0xff, opaque}, pitch)
}

// isFullscreen checks if the current window has the WINDOW_FULLSCREEN
// or WINDOW_FULLSCREEN_DESKTOP flag set.
func isFullscreen(window *sdl.Window) bool {
	flags := window.GetFlags()
	window_fullscreen := (flags & sdl.WINDOW_FULLSCREEN) != 0
	window_fullscreen_desktop := (flags & sdl.WINDOW_FULLSCREEN_DESKTOP) != 0
	return window_fullscreen || window_fullscreen_desktop
}

// toggleFullscreen switches to fullscreen and back.
// Returns true if the mode has been switched to fullscreen
func toggleFullscreen(window *sdl.Window) bool {
	if !isFullscreen(window) {
		// Switch to fullscreen mode
		window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
		return true
	}
	// Switch to windowed mode
	window.SetFullscreen(sdl.WINDOW_SHOWN)
	return false
}

func run() int {

	sdl.Init(sdl.INIT_VIDEO)

	var (
		window   *sdl.Window
		renderer *sdl.Renderer
		err      error
	)

	window, err = sdl.CreateWindow("Pixels!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(width*pixelscale), int32(height*pixelscale), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 1
	}
	defer renderer.Destroy()

	// Fill the render buffer with color #FF00CC
	renderer.SetDrawColor(0xff, 0, 0xcc, opaque)
	renderer.Clear()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	texture.SetBlendMode(sdl.BLENDMODE_BLEND) // sdl.BLENDMODE_ADD is also possible

	rand.Seed(time.Now().UnixNano())

	var (
		pixels = make([]uint32, width*height)
		cores  = runtime.NumCPU()
		event  sdl.Event
		quit   bool
		pause  bool
	)

	// Innerloop
	for !quit {

		if !pause {
			// Draw to pixel buffer
			DrawAll(pixels, cores)

			// Draw pixel buffer to screen
			texture.UpdateRGBA(nil, pixels, width)

			// Clear the render buffer between each frame
			renderer.Clear()

			renderer.Copy(texture, nil, nil)
			renderer.Present()
		}

		// Check for events
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyboardEvent:
				ke := event.(*sdl.KeyboardEvent)
				if ke.Type == sdl.KEYDOWN {
					ks := ke.Keysym
					switch ks.Sym {
					case sdl.K_ESCAPE:
						quit = true
					case sdl.K_q:
						quit = true
					case sdl.K_RETURN:
						altHeldDown := ks.Mod == sdl.KMOD_LALT || ks.Mod == sdl.KMOD_RALT
						if !altHeldDown {
							// alt+enter is not pressed
							break
						}
						// alt+enter is pressed
						fallthrough
					case sdl.K_f, sdl.K_F11:
						if toggleFullscreen(window) {
							sdl.ShowCursor(0)
						} else {
							sdl.ShowCursor(1)
						}
					case sdl.K_SPACE, sdl.K_p:
						pause = !pause
					}
				}
			}
		}
		sdl.Delay(1000 / frameRate)
	}
	return 0
}

func main() {
	// This is to allow the deferred functions in run() to kick in at exit
	os.Exit(run())
}
```
