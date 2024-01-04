package pixelpusher

import (
	"fmt"
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
)

type Config struct {
	Title      string
	PixelScale int
	Width      int
	Height     int
	Pitch      int32
	FrameRate  int
	Opaque     uint8
	Pixels     []uint32
}

// DrawFunction can be used to draw pixels to Config.Pixels
type DrawFunction func(*Config) error

// ActionFunction is called when keys are pressed or released. Order: left, right, up, down, space, return, esc
type ActionFunction func(bool, bool, bool, bool, bool, bool, bool) error

type TickFunction func() error

func New(title string) *Config {
	return &Config{
		Title:      title, // window title
		PixelScale: 4,     // size of "worldspace pixels" measured in "screenspace pixels"
		Width:      320,   // width, worldspace
		Height:     200,   // height, worldspace
		Pitch:      320,   // Same as width, used when calculating where to place pixels (y*pitch+x)
		FrameRate:  60,    // Target framerate
		Opaque:     255,   // Alpha value for opaque colors
		Pixels:     make([]uint32, 320*200),
	}
}

func Plot(c *Config, x, y, r, g, b int) error {
	if x < 0 || x >= c.Width {
		return fmt.Errorf("x is out of range: %d", x)
	}
	if y < 0 || y >= c.Height {
		return fmt.Errorf("y is out of range: %d", y)
	}
	Pixel(c.Pixels, int32(x), int32(y), color.RGBA{uint8(r), uint8(g), uint8(b), c.Opaque}, c.Pitch)
	return nil
}

func (c *Config) Run(drawFunc DrawFunction, pressFunc ActionFunction, releaseFunc ActionFunction, tickFunc TickFunction) error {

	sdl.Init(uint32(sdl.INIT_VIDEO))
	defer sdl.Quit()

	var (
		window   *sdl.Window
		renderer *sdl.Renderer
		err      error
	)

	window, err = sdl.CreateWindow(c.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(c.Width*c.PixelScale), int32(c.Height*c.PixelScale), sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("failed to create window: %s", err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %s", err)
	}
	defer renderer.Destroy()

	// Fill the render buffer with black
	renderer.SetDrawColor(0, 0, 0, c.Opaque)
	renderer.Clear()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, int32(c.Width), int32(c.Height))
	if err != nil {
		return fmt.Errorf("failed to create texture: %s", err)
	}

	var (
		event                     sdl.Event
		pause, recording, quit    bool
		loopCounter, frameCounter uint64
	)

	// Innerloop
	for !quit {

		if !pause {

			if drawFunc != nil {
				drawFunc(c)
			}

			texture.UpdateRGBA(nil, c.Pixels, int(c.Pitch))

			renderer.Copy(texture, nil, nil)
			renderer.Present()

			if recording {
				filename := fmt.Sprintf("frame%05d.png", frameCounter)
				SavePixelsToPNG(c.Pixels, c.Pitch, filename, true)
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
					case sdl.K_ESCAPE, sdl.K_q:
						// quit
						if pressFunc != nil && pressFunc(false, false, false, false, false, false, true) != nil {
							quit = true
						}
					case sdl.K_SPACE:
						// fire
						if pressFunc != nil && pressFunc(false, false, false, false, true, false, false) != nil {
							quit = true
						}
					case sdl.K_LEFT, sdl.K_a:
						// left
						if pressFunc != nil && pressFunc(true, false, false, false, false, false, false) != nil {
							quit = true
						}
					case sdl.K_RIGHT, sdl.K_d:
						// right
						if pressFunc != nil && pressFunc(false, true, false, false, false, false, false) != nil {
							quit = true
						}
					case sdl.K_UP, sdl.K_w:
						// up
						if pressFunc != nil && pressFunc(false, false, true, false, false, false, false) != nil {
							quit = true
						}
					case sdl.K_DOWN:
						// down
						if pressFunc != nil && pressFunc(false, false, false, true, false, false, false) != nil {
							quit = true
						}

					case sdl.K_RETURN:
						altHeldDown := ks.Mod == sdl.KMOD_LALT || ks.Mod == sdl.KMOD_RALT
						if !altHeldDown {
							// alt+enter is not pressed
							// enter is pressed
							if pressFunc != nil && pressFunc(false, false, false, false, false, true, false) != nil {
								quit = true
							}
							break
						}
						// alt+enter is pressed
						fallthrough
					case sdl.K_f, sdl.K_F11:
						ToggleFullscreen(window)
					case sdl.K_p:
						// pause toggle
						pause = !pause
					case sdl.K_s:
						ctrlHeldDown := ks.Mod == sdl.KMOD_LCTRL || ks.Mod == sdl.KMOD_RCTRL
						if !ctrlHeldDown {
							// ctrl+s is not pressed
							// s is pressed
							if pressFunc != nil && pressFunc(false, false, false, true, false, false, false) != nil {
								quit = true
							}
							break
						}
						// ctrl+s is pressed
						fallthrough
					case sdl.K_F12:
						// screenshot
						Screenshot(renderer, "screenshot.png", true)
					case sdl.K_r:
						// recording
						recording = !recording
						frameCounter = 0
					}
				} else if ke.Type == sdl.KEYUP {
					ks := ke.Keysym
					switch ks.Sym {
					case sdl.K_ESCAPE, sdl.K_q:
						// quit
						if releaseFunc != nil && releaseFunc(false, false, false, false, false, false, true) != nil {
							quit = true
						}
					case sdl.K_SPACE:
						// fire
						if releaseFunc != nil && releaseFunc(false, false, false, false, true, false, false) != nil {
							quit = true
						}
					case sdl.K_LEFT, sdl.K_a:
						// left
						if releaseFunc != nil && releaseFunc(true, false, false, false, false, false, false) != nil {
							quit = true
						}
					case sdl.K_RIGHT, sdl.K_d:
						// right
						if releaseFunc != nil && releaseFunc(false, true, false, false, false, false, false) != nil {
							quit = true
						}
					case sdl.K_UP, sdl.K_w:
						// up
						if releaseFunc != nil && releaseFunc(false, false, true, false, false, false, false) != nil {
							quit = true
						}
					case sdl.K_DOWN, sdl.K_s:
						// down
						if releaseFunc != nil && releaseFunc(false, false, false, true, false, false, false) != nil {
							quit = true
						}
					case sdl.K_RETURN:
						if releaseFunc != nil && releaseFunc(false, false, false, false, false, true, false) != nil {
							quit = true
						}
					}
				}
			}
		}
		sdl.Delay(uint32(1000 / c.FrameRate))
		loopCounter++
		if tickFunc != nil && tickFunc() != nil {
			quit = true
		}
	}
	return nil
}
