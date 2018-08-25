package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/fogleman/fauxgl"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/xyproto/multirender"
	"github.com/xyproto/sdl2utils"
)

const (
	// Window title
	windowTitle = "Duck!"

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

// DrawAll fills the pixel buffer with pixels.
// "cores" is how many CPU cores should be targeted when drawing triangles,
// by launching the same number of goroutines.
func DrawAll(pixels []uint32, cores int, mesh *fauxgl.Mesh, cameraAngle float32, meshHexColor string, enr int) {
}

func run() int {

	sdl.Init(sdl.INIT_VIDEO)

	var (
		window   *sdl.Window
		renderer *sdl.Renderer
		err      error
	)

	window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(width*pixelscale), int32(height*pixelscale), sdl.WINDOW_SHOWN)
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

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	//texture.SetBlendMode(sdl.BLENDMODE_BLEND) // sdl.BLENDMODE_ADD is also possible

	rand.Seed(time.Now().UnixNano())

	mesh, err := LoadMeshOBJ("duck.obj")
	if err != nil {
		panic(err)
	}
	meshHexColor := "#ffff00"

	var (
		pixels      = make([]uint32, width*height)
		pixelCopy   = make([]uint32, width*height)
		event       sdl.Event
		quit        bool
		pause       bool
		cameraAngle float32
		enr         int  = 5
		cores            = runtime.NumCPU()
		effect      bool = true
	)

	// Innerloop
	for !quit {

		if !pause {
			// Clear the pixel buffer
			multirender.FastClear(pixels, 0xffffffff)

			// Draw a triangle, concurrently
			if effect {
				yellow := rb()
				multirender.Triangle(cores, pixels, rw(), rh(), rw(), rh(), rw(), rh(), color.RGBA{yellow, yellow, 0, 0x80}, pitch)
			}

			// Draw a 3D object
			DrawMesh(pixels, pitch, mesh, cameraAngle, meshHexColor)

			if effect {
				copy(pixelCopy, pixels)

				// Filter
				Convolution(0.9, pixelCopy, width, height, pitch, enr)

				// Stretch contrast
				multirender.StretchContrast(cores, pixelCopy, pitch, 0.09)

				// Draw pixel buffer to screen
				texture.UpdateRGBA(nil, pixelCopy, width)
			} else {
				// Draw pixel buffer to screen
				texture.UpdateRGBA(nil, pixels, width)
			}

			//if int32(cameraAngle) % 4 == 0 {
			//	enr++
			//	if enr > 3 {
			//		enr = 0
			//	}
			//}

			// Clear the pixel buffer
			//multirender.FastClear(pixels, 0)

			renderer.Copy(texture, nil, nil)
			renderer.Present()

			cameraAngle += 0.1
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
					case sdl.K_SPACE:
						effect = !effect
					case sdl.K_p:
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
						// save the image
						multirender.SavePixelsToPNG(pixels, pitch, "screenshot.png", true)
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
