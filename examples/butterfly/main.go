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

func Trippy(pixels []uint32, width, height uint32, pitch int32) {
	for y := int32(1); y < int32(height-1); y++ {
		for x := int32(1); x < int32(width-1); x++ {
			left := pixels[y*pitch+x-1]
			right := pixels[y*pitch+x+1]
			this := pixels[y*pitch+x]
			above := pixels[(y+1)*pitch+x]
			// Dividing the raw uint32 color value!
			average := (left + right + this + above) / 8
			pixels[y*pitch+x] = average
		}
	}
}

// TODO: Find out why this is trippy and not the flame effect :D
func Flame(time float32, pixels []uint32, width, height uint32, pitch int32) {

	// Make the effect increase and decrease in intensity instead of increasing and then dropping down to 0 again
	stime := float32(math.Sin(float64(time) * math.Pi))

	for y := int32(2); y < int32(height-2); y++ {
		for x := int32(2); x < int32(width-2); x++ {

			// "snow patterns"
			//left := pixels[y*pitch+x-1]
			//right := pixels[y*pitch+x+1]
			//this := pixels[y*pitch+x]
			//above := pixels[(y+1)*pitch+x]

			// "highway"
			//left := pixels[(y-1)*pitch+x-1]
			//right := pixels[(y-1)*pitch+x+1]
			//this := pixels[y*pitch+x]
			//above := pixels[(y-1)*pitch+x]

			// "dither highway"
			//left := pixels[(y-1)*pitch+x-1]
			//right := pixels[(y-1)*pitch+x+1]
			//this := pixels[(y-1)*pitch+(x-1)]
			//above := pixels[(y+1)*pitch+(x+1)]

			// "speeding"
			//two1 := int32(2.0 - stime * 4.0)
			//two2 := int32(2.0 - time * 4.0)
			//left := pixels[y*pitch+(x-two2)]
			//right := pixels[y*pitch+(x+two1)]
			//this := pixels[y*pitch+int32(float32(x)*stime)]
			//above := pixels[(y-two2)*pitch+int32(float32(x)*stime)]

			// "castle"
			//one1 := int32(1.0 - stime * 2.0)
			//one2 := int32(1.0 - time * 2.0)
			//two1 := int32(2.0 - stime * 4.0)
			//left := pixels[y*pitch+(x-one1)]
			//right := pixels[y*pitch+(x+one1)]
			//this := pixels[y*pitch+x*two1]
			//above := pixels[(y-one2)*pitch+x*two1]

			// "butterfly"
			two1 := int32(2.0 - stime*4.0)
			two2 := int32(2.0 - time*4.0)
			left := pixels[y*pitch+(x-two1)]
			right := pixels[y*pitch+(x+two1)]
			this := pixels[y*pitch+x*two2]
			above := pixels[(y-two1)*pitch+x*two2]

			lr, lg, lb, _ := multirender.ColorValueToRGBA(left)
			rr, rg, rb, _ := multirender.ColorValueToRGBA(right)
			tr, tg, tb, _ := multirender.ColorValueToRGBA(this)
			ar, ag, ab, _ := multirender.ColorValueToRGBA(above)

			averageR := uint8(float32(lr+rr+tr+ar) / float32(4.55-stime))
			averageG := uint8(float32(lg+rg+tg+ag) / float32(4.55-stime))
			averageB := uint8(float32(lb+rb+tb+ab) / float32(4.55-stime))

			pixels[y*pitch+x] = multirender.RGBAToColorValue(averageR, averageG, averageB, 0xff)
		}
	}
	// Top row
	for y := int32(0); y < int32(2); y++ {
		for x := int32(0); x < int32(width); x++ {
			pixels[y*pitch+x] = 0
		}
	}
	// Bottom row
	for y := int32(height - 2); y < int32(height-1); y++ {
		for x := int32(0); x < int32(width); x++ {
			pixels[y*pitch+x] = 0
		}
	}
	// Left col
	for x := int32(0); x < int32(2); x++ {
		for y := int32(0); y < int32(height); y++ {
			pixels[y*pitch+x] = 0
		}
	}
	// Right col
	for x := int32(width - 2); x < int32(width-1); x++ {
		for y := int32(0); y < int32(height); y++ {
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

// Every pixel that is not really light, decrease the intensity
func Darken(pixels []uint32) {
	for i := range pixels {
		r, g, b, _ := multirender.ColorValueToRGBA(pixels[i])
		if r > 240 && g > 240 && b > 240 {
			continue
		}
		pixels[i] = multirender.RGBAToColorValue(clamp(float32(r)*1.1, 255), clamp(float32(g)*1.1, 255), clamp(float32(b)*1.1, 255), 255)
	}
}

// TriangleDance draws a dancing triangle, as time goes from 0.0 to 1.0.
// The returned value signals to wich degree the graphics should be transitioned out.
func TriangleDance(time float32, pixels []uint32, width, height uint32, pitch int32, cores int, xdirection, ydirection int) (transition float32) {

	size := uint32(70)

	//var bgColorValue uint32 = 0x4e7f9eff

	// The function is responsible for clearing the pixels,
	// it might want to reuse the pixels from the last time (flame effect)
	//multirender.FastClear(pixels, bgColorValue)

	// Find a suitable placement and color
	x := int32(0)
	if xdirection > 0 {
		x = int32(multirender.Clamp(uint32(float32(width)*time), size, width-size))
	} else if xdirection == 0 {
		x = int32(width / 2)
	} else {
		x = int32(multirender.Clamp(uint32(float32(width)*(1.0-time)), size, width-size))
	}
	y := int32(0)
	if ydirection > 0 {
		y = int32(multirender.Clamp(uint32(float32(height)*time), size, height-size))
	} else if ydirection == 0 {
		y = int32(height / 2)
	} else {
		y = int32(multirender.Clamp(uint32(float32(height)*(1.0-time)), size, height-size))
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

	cycleTime := float32(0.0)
	flameStart := float32(0.75)
	flameTime := flameStart
	flameTimeAdd := float32(0.0001)

	var loopCounter int64 = 0

	// Innerloop
	for !quit {

		if !pause {
			if loopCounter%4 == 0 {
				// Draw to the pixel buffer
				TriangleDance(cycleTime, pixels, width, height, pitch, cores, 1, 0)
				TriangleDance(cycleTime, pixels, width, height, pitch, cores, -1, 0)
				TriangleDance(cycleTime, pixels, width, height, pitch, cores, 0, 1)
				TriangleDance(cycleTime, pixels, width, height, pitch, cores, 0, -1)
				Flame(flameTime, pixels, width, height, pitch)
			}
			Flame(flameTime, pixels, width, height, pitch)

			// Keep track of the time given to TriangleDance
			cycleTime += 0.0002
			if cycleTime >= 1.0 {
				cycleTime = 0.0
			}

			// Keep track of the time given to Flame
			flameTime += flameTimeAdd
			if flameTime >= 0.81 {
				flameTime = flameStart
				flameTimeAdd = -flameTimeAdd
			} else if flameTime <= flameStart {
				flameTime = flameStart
				flameTimeAdd = -flameTimeAdd
			}

			//Darken(pixels)

			// Draw pixel buffer to screen, but inverted
			Invert(pixels)

			// Draw the center red triangle, in flameTime
			//TriangleDance(flameTime, pixels, width, height, pitch, cores, 0, 0)
			texture.UpdateRGBA(nil, pixels, width)
			Invert(pixels)

			renderer.Copy(texture, nil, nil)
			renderer.Present()

			// Clear the render buffer between each frame
			//renderer.Clear()

		} else {
			fmt.Println("FLAME TIME", flameTime)
			fmt.Println("CYCLE TIME", cycleTime)
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
