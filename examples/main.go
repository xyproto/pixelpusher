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

	pitch = width

	// Target framerate
	frameRate = 60

	// Alpha value for opaque colors
	opaque = 255
)

// ranb returns a random byte
func ranb() uint8 {
	return uint8(rand.Intn(255))
}

// DrawAll fills the pixel buffer with pixels
func DrawAll(cores int, pixels []uint32) {

	// First draw triangles, concurrently
	multirender.Triangle(cores, pixels, rand.Int31n(width), rand.Int31n(height), rand.Int31n(width), rand.Int31n(height), rand.Int31n(width), rand.Int31n(height), color.RGBA{ranb(), ranb(), ranb(), opaque}, pitch)

	// Then draw lines and pixels, without caring about which order they appear in
	go multirender.Line(pixels, rand.Int31n(width), rand.Int31n(height), rand.Int31n(width), rand.Int31n(height), color.RGBA{ranb(), ranb(), ranb(), opaque}, pitch)
	go multirender.Pixel(pixels, rand.Int31n(width), rand.Int31n(height), color.RGBA{255, 0, 0, ranb()}, pitch)
}

// isFullscreen checks if the current window has the WINDOW_FULLSCREEN_DESKTOP flag set
func isFullscreen(window *sdl.Window) bool {
	currentFlags := window.GetFlags()
	fullscreen1 := (currentFlags & sdl.WINDOW_FULLSCREEN_DESKTOP) != 0
	fullscreen2 := (currentFlags & sdl.WINDOW_FULLSCREEN) != 0
	return fullscreen1 || fullscreen2
}

// toggleFullscreen switches to fullscreen and back
// returns true if the mode has been switched to fullscreen
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

	renderer.SetDrawColor(0, 0, 0, opaque)
	renderer.Clear()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	rand.Seed(time.Now().UnixNano())

	cores := runtime.NumCPU()
	pixels := make([]uint32, width*height)

	var event sdl.Event
	var running = true

	// Innerloop
	for running {

		// Draw to pixel buffer
		DrawAll(cores, pixels)

		// Draw pixel buffer to screen
		texture.UpdateRGBA(nil, pixels, width)
		renderer.Clear()
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		// Check for events
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				ke := event.(*sdl.KeyboardEvent)
				if ke.Type == sdl.KEYDOWN {
					ks := ke.Keysym
					switch ks.Sym {
					case sdl.K_ESCAPE:
						running = false
					case sdl.K_q:
						running = false
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
