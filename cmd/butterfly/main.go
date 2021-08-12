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
	"github.com/xyproto/pixelpusher"
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

func convolution(time float32, pixels []uint32, width, height, pitch int32, enr int) {

	// Make the effect increase and decrease in intensity instead of increasing and then dropping down to 0 again
	stime := float32(math.Sin(float64(time) * math.Pi))
	var left, right, this, above uint32
	two1 := int32(2.0 - stime*4.0)
	two2 := int32(2.0 - time*4.0)
	one1 := int32(1.0 - stime*2.0)
	one2 := int32(1.0 - time*2.0)

	for y := int32(2); y < int32(height-2); y++ {
		for x := int32(2); x < int32(width-2); x++ {

			switch enr {
			case 0:
				// "snow patterns"
				left = pixels[y*pitch+x-1]
				right = pixels[y*pitch+x+1]
				this = pixels[y*pitch+x]
				above = pixels[(y+1)*pitch+x]
			case 1:
				// "highway"
				left = pixels[(y-1)*pitch+x-1]
				right = pixels[(y-1)*pitch+x+1]
				this = pixels[y*pitch+x]
				above = pixels[(y-1)*pitch+x]
			case 2:
				// "dither highway"
				left = pixels[(y-1)*pitch+x-1]
				right = pixels[(y-1)*pitch+x+1]
				this = pixels[(y-1)*pitch+(x-1)]
				above = pixels[(y+1)*pitch+(x+1)]
			case 3:
				// "butterfly"
				left = pixels[y*pitch+(x-two1)]
				right = pixels[y*pitch+(x+two1)]
				this = pixels[y*pitch+x*two2]
				above = pixels[(y-two1)*pitch+x*two2]
			case 4:
				// ?
				left = pixels[y*pitch+(x-two2)]
				right = pixels[y*pitch+(x+two1)]
				this = pixels[y*pitch+int32(float32(x)*stime)]
				above = pixels[(y-two2)*pitch+int32(float32(x)*stime)]
			case 5:
				// "castle"
				left = pixels[y*pitch+(x-one1)]
				right = pixels[y*pitch+(x+one1)]
				this = pixels[y*pitch+x*two1]
				above = pixels[(y-one2)*pitch+x*two1]
			}

			lr, lg, lb, _ := pixelpusher.ColorValueToRGBA(left)
			rr, rg, rb, _ := pixelpusher.ColorValueToRGBA(right)
			tr, tg, tb, _ := pixelpusher.ColorValueToRGBA(this)
			ar, ag, ab, _ := pixelpusher.ColorValueToRGBA(above)

			averageR := uint8(float32(lr+rr+tr+ar) / float32(4.55-stime))
			averageG := uint8(float32(lg+rg+tg+ag) / float32(4.55-stime))
			averageB := uint8(float32(lb+rb+tb+ab) / float32(4.55-stime))

			pixels[y*pitch+x] = pixelpusher.RGBAToColorValue(averageR, averageG, averageB, 0xff)
		}
	}
	// Top row
	for y := int32(0); y < int32(2); y++ {
		for x := int32(0); x < width; x++ {
			pixels[y*pitch+x] = 0
		}
	}
	// Bottom row
	for y := height - 2; y < height-1; y++ {
		for x := int32(0); x < width; x++ {
			pixels[y*pitch+x] = 0
		}
	}
	// Left col
	for x := int32(0); x < int32(2); x++ {
		for y := int32(0); y < height; y++ {
			pixels[y*pitch+x] = 0
		}
	}
	// Right col
	for x := width - 2; x < width-1; x++ {
		for y := int32(0); y < height; y++ {
			pixels[y*pitch+x] = 0
		}
	}
}

// Invert the colors, but set the alpha to 255
func Invert(pixels []uint32) {
	for i := range pixels {
		pixels[i] = (0xffffffff - pixels[i]) | 0x000000ff
	}
}

func clamp(v float32, max uint8) uint8 {
	u := uint32(v)
	if u > uint32(max) {
		return 255
	}
	return uint8(u)
}

// Every pixel that is light, decrease the intensity
func darken(pixels []uint32) {
	for i := range pixels {
		r, g, b, _ := pixelpusher.ColorValueToRGBA(pixels[i])
		if r < 20 && g < 20 && b < 20 {
			continue
		}
		pixels[i] = pixelpusher.RGBAToColorValue(clamp(float32(r)*0.99, 255), clamp(float32(g)*0.99, 255), clamp(float32(b)*0.99, 255), 255)
	}
}

// TriangleDance draws a dancing triangle, as time goes from 0.0 to 1.0.
// The returned value signals to which degree the graphics should be transitioned out.
func TriangleDance(time float32, pixels []uint32, width, height, pitch int32, cores int, xdirection, ydirection int) (transition float32) {

	size := int32(70)

	//var bgColorValue uint32 = 0x4e7f9eff

	// The function is responsible for clearing the pixels,
	// it might want to reuse the pixels from the last time (flame effect)
	//pixelpusher.FastClear(pixels, bgColorValue)

	// Find a suitable placement and color
	var x int32
	if xdirection > 0 {
		x = pixelpusher.Clamp(int32(float32(width)*time), size, width-size)
	} else if xdirection == 0 {
		x = int32(width / 2)
	} else {
		x = pixelpusher.Clamp(int32(float32(width)*(1.0-time)), size, width-size)
	}
	var y int32
	if ydirection > 0 {
		y = pixelpusher.Clamp(int32(float32(height)*time), size, height-size)
	} else if ydirection == 0 {
		y = int32(height / 2)
	} else {
		y = pixelpusher.Clamp(int32(float32(height)*(1.0-time)), size, height-size)
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

	pixelpusher.Triangle(cores, pixels, x1, y1, x2, y2, x3, y3, c, pitch)

	return 0.0
}

func run() int {

	sdl.Init(sdl.INIT_VIDEO)
	defer sdl.Quit()

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
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return 1
	}

	//texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	rand.Seed(time.Now().UnixNano())

	var (
		pixels = make([]uint32, width*height)
		cores  = runtime.NumCPU()
		event  sdl.Event

		pause, quit bool

		cycleTime    float32
		flameStart   float32 = 0.75
		flameTime    float32 = 0.75
		flameTimeAdd float32 = 0.0001

		loopCounter int64
		enr         = 3 // effect number
	)

	// Innerloop
	for !quit {

		if !pause {

			// Invert pixels before drawing
			Invert(pixels)

			if loopCounter%4 == 0 {
				// Draw to the pixel buffer
				TriangleDance(cycleTime, pixels, width, height, pitch, cores, 1, 0)
				TriangleDance(cycleTime, pixels, width, height, pitch, cores, -1, 0)
				TriangleDance(cycleTime, pixels, width, height, pitch, cores, 0, 1)
				TriangleDance(cycleTime, pixels, width, height, pitch, cores, 0, -1)
				convolution(flameTime, pixels, width, height, pitch, enr)
			}
			convolution(flameTime, pixels, width, height, pitch, enr)

			// Keep track of the time given to TriangleDance
			cycleTime += 0.002
			if cycleTime >= 1.0 {
				cycleTime = 0.0
				enr++
				if enr > 3 {
					enr = 0
				}
			}

			// Keep track of the time given to convolution
			flameTime += flameTimeAdd
			if flameTime >= 0.81 {
				flameTime = flameStart
				flameTimeAdd = -flameTimeAdd
			} else if flameTime <= flameStart {
				flameTime = flameStart
				flameTimeAdd = -flameTimeAdd
			}

			Invert(pixels)

			// Draw the center red triangle, in flameTime
			//TriangleDance(flameTime, pixels, width, height, pitch, cores, 0, 0)
			texture.UpdateRGBA(nil, pixels, pitch)
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
