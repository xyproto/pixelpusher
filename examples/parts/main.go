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
	"github.com/xyproto/sdl2utils"
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

// TriangleDance draws a dancing triangle, as time goes from 0.0 to 1.0.
// The returned value signals to wich degree the graphics should be transitioned out.
func TriangleDance(time float32, pixels []uint32, width, height uint32, pitch int32, cores int) (transition float32) {

	var bgColorValue uint32 = 0x4e7f9eff

	// The function is responsible for clearing the pixels,
	// it might want to reuse the pixels from the last time (flame effect)
	multirender.FastClear(pixels, bgColorValue)

	// Find a suitable placement and color
	x := int32(multirender.Clamp(uint32(float32(width)*time), 80, width-80))
	y := int32(height / 2)
	c := color.RGBA{rb(), rb(), rb(), 0xff}

	x1 := x
	y1 := y
	x2 := x + rand.Int31n(40) - 20
	y2 := y + rand.Int31n(40) - 20
	x3 := x + rand.Int31n(40) - 20
	y3 := y + rand.Int31n(40) - 20

	multirender.Triangle(cores, pixels, x1, y1, x2, y2, x3, y3, c, pitch)

	return 0.0
}

func run() int {

	sdl.Init(sdl.INIT_VIDEO)

	var (
		window   *sdl.Window
		renderer *sdl.Renderer
		err      error
	)

	window, err = sdl.CreateWindow("Butterfly", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(width*pixelscale), int32(height*pixelscale), sdl.WINDOW_SHOWN)
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

	// Fill the render buffer with black
	renderer.SetDrawColor(0, 0, 0, opaque)
	renderer.Clear()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	//texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	rand.Seed(time.Now().UnixNano())

	var (
		pixels = make([]uint32, width*height)
		cores  = runtime.NumCPU()
		event  sdl.Event
		quit   bool
		pause  bool
	)

	partTime := float32(0.0)

	var loopCounter int64 = 0

	// Innerloop
	for !quit {

		if !pause {
			if loopCounter%5 == 0 {
				// Draw to the pixel buffer
				TriangleDance(partTime, pixels, width, height, pitch, cores)
			}

			// Keep track of the time given to TriangleDance
			partTime += 0.002
			if partTime >= 1.0 {
				partTime = 0.0
			}

			// Draw pixel buffer to screen
			texture.UpdateRGBA(nil, pixels, width)

			renderer.Copy(texture, nil, nil)
			renderer.Present()

			// Clear the render buffer between each frame
			//renderer.Clear()

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
						sdl2utils.ToggleFullscreen(window)
					case sdl.K_p, sdl.K_SPACE:
						// pause toggle
						pause = !pause
					case sdl.K_s:
						ctrlHeldDown := ks.Mod == sdl.KMOD_LCTRL || ks.Mod == sdl.KMOD_RCTRL
						if !ctrlHeldDown {
							// ctrl+s is not pressed
							break
						}
						// ctrl+s is pressed
						fallthrough
					case sdl.K_F12:
						// screenshot
						sdl2utils.Screenshot(renderer, "screenshot.png", true)
					}
				}
			}
		}
		sdl.Delay(1000 / frameRate)
		loopCounter++
	}
	return 0
}

func main() {
	// This is to allow the deferred functions in run() to kick in at exit
	os.Exit(run())
}
