package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/xyproto/multirender"
	"github.com/xyproto/pf"
	"github.com/xyproto/sdl2utils"
)

const (
	// Size of "worldspace pixels", measured in "screenspace pixels"
	pixelscale = 1

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

func Trippy(pixels []uint32, width, height, pitch int32) {
	for y := int32(1); y < int32(height-1); y++ {
		for x := int32(1); x < int32(width-1); x++ {
			left := pixels[y*pitch+x-1]
			right := pixels[y*pitch+x+1]
			this := pixels[y*pitch+x]
			above := pixels[(y+1)*pitch+x]
			// Dividing the raw uint32 color value!
			average := (left + right + this + above) / 4
			pixels[y*pitch+x] = average
		}
	}
}

func Convolution(time float32, pixels []uint32, width, height, pitch int32, enr int) {

	// Make the effect increase and decrease in intensity instead of increasing and then dropping down to 0 again
	stime := float32(math.Sin(float64(time) * math.Pi))
	var left, right, this, above uint32
	two1 := int32(2.0 - stime*4.0)
	two2 := int32(2.0 - time*4.0)
	one1 := int32(1.0 - stime*2.0)
	one2 := int32(1.0 - time*2.0)

	size := width * height

	for y := int32(0); y < height; y++ {
		for x := int32(0); x < width; x++ {
			switch enr {
			case 0:
				// "snow patterns"
				left = multirender.GetWrap(pixels, y*pitch+x-1, size)
				right = multirender.GetWrap(pixels, y*pitch+x+1, size)
				this = multirender.GetWrap(pixels, y*pitch+x, size)
				above = multirender.GetWrap(pixels, (y+1)*pitch+x, size)
			case 1:
				// "highway"
				left = multirender.GetWrap(pixels, (y-1)*pitch+x-1, size)
				right = multirender.GetWrap(pixels, (y-1)*pitch+x+1, size)
				this = multirender.GetWrap(pixels, y*pitch+x, size)
				above = multirender.GetWrap(pixels, (y-1)*pitch+x, size)
			case 2:
				// "dither highway"
				left = multirender.GetWrap(pixels, (y-1)*pitch+x-1, size)
				right = multirender.GetWrap(pixels, (y-1)*pitch+x+1, size)
				this = multirender.GetWrap(pixels, (y-1)*pitch+(x-1), size)
				above = multirender.GetWrap(pixels, (y+1)*pitch+(x+1), size)
			case 3:
				// "butterfly"
				left = multirender.GetWrap(pixels, y*pitch+(x-two1), size)
				right = multirender.GetWrap(pixels, y*pitch+(x+two1), size)
				this = multirender.GetWrap(pixels, y*pitch+x*two2, size)
				above = multirender.GetWrap(pixels, (y-two1)*pitch+x*two2, size)
			case 4:
				// ?
				left = multirender.GetWrap(pixels, y*pitch+(x-two2), size)
				right = multirender.GetWrap(pixels, y*pitch+(x+two1), size)
				this = multirender.GetWrap(pixels, y*pitch+int32(float32(x)*stime), size)
				above = multirender.GetWrap(pixels, (y-two2)*pitch+int32(float32(x)*stime), size)
			case 5:
				// "castle"
				left = multirender.GetWrap(pixels, y*pitch+(x-one1), size)
				right = multirender.GetWrap(pixels, y*pitch+(x+one1), size)
				this = multirender.GetWrap(pixels, y*pitch+x*two1, size)
				above = multirender.GetWrap(pixels, (y-one2)*pitch+x*two1, size)
			}

			lr, lg, lb, _ := multirender.ColorValueToRGBA(left)
			rr, rg, rb, _ := multirender.ColorValueToRGBA(right)
			tr, tg, tb, _ := multirender.ColorValueToRGBA(this)
			ar, ag, ab, _ := multirender.ColorValueToRGBA(above)

			averageR := uint8(float32(lr+rr+tr+ar) / float32(4.8-stime))
			averageG := uint8(float32(lg+rg+tg+ag) / float32(4.8-stime))
			averageB := uint8(float32(lb+rb+tb+ab) / float32(4.8-stime))

			multirender.SetWrap(pixels, y*pitch+x, width*height, multirender.RGBAToColorValue(averageR, averageG, averageB, 0xff))
		}
	}
}

func clamp(v float32, max uint8) uint8 {
	u := uint32(v)
	if u > uint32(max) {
		return 255
	}
	return uint8(u)
}

// TriangleDance draws a dancing triangle, as time goes from 0.0 to 1.0.
// The returned value signals to wich degree the graphics should be transitioned out.
func TriangleDance(cores int, time float32, pixels []uint32, width, height, pitch int32, xdirection, ydirection int) (transition float32) {

	time *= time

	size := int32(120)

	//var bgColorValue uint32 = 0x4e7f9eff

	// The function is responsible for clearing the pixels,
	// it might want to reuse the pixels from the last time (flame effect)
	//multirender.FastClear(pixels, bgColorValue)

	// Find a suitable placement and color
	x := int32(0)
	if xdirection > 0 {
		x = multirender.Clamp(int32(float32(width)*time), size, width-size)
	} else if xdirection == 0 {
		x = int32(width / 2)
	} else {
		x = multirender.Clamp(int32(float32(width)*(1.0-time)), size, width-size)
	}
	y := int32(0)
	if ydirection > 0 {
		y = multirender.Clamp(int32(float32(height)*time), size, height-size)
	} else if ydirection == 0 {
		y = int32(height / 2)
	} else {
		y = multirender.Clamp(int32(float32(height)*(1.0-time)), size, height-size)
	}

	// Make the center triangle red
	var c color.RGBA
	if xdirection == 0 && ydirection == 0 {
		c = color.RGBA{0xff, 0, 0, 0xff}
	} else {
		c = color.RGBA{rb(), rb(), rb(), 0xff}
	}

	x1 := x
	y1 := y
	x2 := x + rand.Int31n(int32(size)) - int32(size/2)
	y2 := y + rand.Int31n(int32(size)) - int32(size/2)
	x3 := x + rand.Int31n(int32(size)) - int32(size/2)
	y3 := y + rand.Int31n(int32(size)) - int32(size/2)

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

	window, err = sdl.CreateWindow("Trippy Triangles", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(width*pixelscale), int32(height*pixelscale), sdl.WINDOW_SHOWN)
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
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return 1
	}

	rand.Seed(time.Now().UnixNano())

	var (
		pixels    = make([]uint32, width*height)
		pixelCopy = make([]uint32, width*height)
		cores     = runtime.NumCPU()
		event     sdl.Event
		quit      bool
		pause     bool
		recording bool
	)

	cycleTime := float32(0.0)
	flameStart := float32(0.75)
	convTime := flameStart
	convTimeAdd := float32(0.0001)

	var loopCounter uint64 = 0
	var frameCounter uint64 = 0

	// effect number
	enr := 3

	// post processing effect number
	lineEffect := 0

	// PixelFunction for inverting the colors
	lineEffectFunction := pf.Combine(pf.Invert, pf.OrBlue)

	// Innerloop
	for !quit {

		if !pause {

			// Invert pixels, and or with Blue, before drawing
			pf.Map(cores, lineEffectFunction, pixels)

			if loopCounter%4 == 0 {
				// Draw to the pixel buffer
				TriangleDance(cores, cycleTime, pixels, width, height, pitch, 1, 0)
				TriangleDance(cores, cycleTime, pixels, width, height, pitch, -1, 0)
				TriangleDance(cores, cycleTime, pixels, width, height, pitch, 0, 1)
				TriangleDance(cores, cycleTime, pixels, width, height, pitch, 0, -1)
				Convolution(convTime, pixels, width, height, pitch, enr)
			}
			Convolution(convTime, pixels, width, height, pitch, enr)

			// Keep track of the time given to TriangleDance
			cycleTime += 0.002
			if cycleTime >= 1.0 {
				cycleTime = 0.0
				enr++
				if enr > 3 {
					enr = 0
				}
			}

			// Keep track of the time given to Convolution
			convTime += convTimeAdd
			if convTime >= 0.81 {
				convTime = flameStart
				convTimeAdd = -convTimeAdd
			} else if convTime <= flameStart {
				convTime = flameStart
				convTimeAdd = -convTimeAdd
			}

			// Take a copy before applying post-processing
			copy(pixelCopy, pixels)

			// Invert the pixels back after adding all the things above
			pf.Map(cores, lineEffectFunction, pixels)

			// Stretch the contrast on a copy of the pixels
			multirender.StretchContrast(cores, pixelCopy, pitch, cycleTime)

			//RemoveBlue(cores, pixelCopy)

			texture.UpdateRGBA(nil, pixelCopy, pitch)

			renderer.Copy(texture, nil, nil)
			renderer.Present()

			if recording {
				filename := fmt.Sprintf("frame%05d.png", frameCounter)
				multirender.SavePixelsToPNG(pixelCopy, pitch, filename, true)
				frameCounter++
			}

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
					case sdl.K_p:
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
					case sdl.K_SPACE:
						// Alternate between PixelFunctions that are applied as a post-processing filter
						lineEffect++
						if lineEffect > 1 {
							lineEffect = 0
						}
						switch lineEffect {
						case 0:
							// The combined pixel functions Invert and OrBlue
							lineEffectFunction = pf.Combine(pf.Invert, pf.OrBlue)
						case 1:
							lineEffectFunction = pf.Invert
						}
					case sdl.K_F12:
						// screenshot
						sdl2utils.Screenshot(renderer, "screenshot.png", true)
					case sdl.K_r:
						// recording
						recording = !recording
						frameCounter = 0
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
