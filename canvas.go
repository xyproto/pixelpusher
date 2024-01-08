package pixelpusher

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
)

// Canvas is a window title + pixels + additional info
type Canvas struct {
	Title      string
	PixelScale int
	Width      int
	Height     int
	Pitch      int32
	FrameRate  int
	Opaque     uint8
	Pixels     []uint32
}

// DrawFunction can be used to draw pixels to canvas.Pixels
type DrawFunction func(*Canvas) error

// ActionFunction is called when keys are pressed or released. Order: left, right, up, down, space, return, esc
type ActionFunction func(bool, bool, bool, bool, bool, bool, bool) error

// TickFunction is called at every loop
type TickFunction func() error

var errQuit = errors.New("quit")

// New creates a new Canvas
func New(title string) *Canvas {
	return &Canvas{
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

// Plot can plot a pixel to a canvas. It is a bit slow because it contains additional checks. Modify canvas.Pixels directly for better performance.
func Plot(c *Canvas, x, y, r, g, b int) error {
	if x < 0 || x >= c.Width {
		return fmt.Errorf("x is out of range: %d", x)
	}
	if y < 0 || y >= c.Height {
		return fmt.Errorf("y is out of range: %d", y)
	}
	Pixel(c.Pixels, int32(x), int32(y), color.RGBA{uint8(r), uint8(g), uint8(b), c.Opaque}, c.Pitch)
	return nil
}

// Run takes an optional function draw drawing pixels, an optional function for when an action is pressed, an optional function for when an action is released and a function for each loop
func (c *Canvas) Run(drawFunc DrawFunction, pressFunc ActionFunction, releaseFunc ActionFunction, tickFunc TickFunction) error {
	var (
		window                    *sdl.Window
		renderer                  *sdl.Renderer
		event                     sdl.Event
		joystick                  *sdl.Joystick
		pause, recording          bool
		loopCounter, frameCounter uint64
		err                       error
	)

	// Initialize SDL (video + joystick)
	sdl.Init(uint32(sdl.INIT_VIDEO) | uint32(sdl.INIT_JOYSTICK))
	defer sdl.Quit()

	// Create a window
	window, err = sdl.CreateWindow(c.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(c.Width*c.PixelScale), int32(c.Height*c.PixelScale), sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("failed to create window: %s", err)
	}
	defer window.Destroy()

	// Create a renderer
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %s", err)
	}
	defer renderer.Destroy()

	// Fill the render buffer with black
	renderer.SetDrawColor(0, 0, 0, c.Opaque)
	renderer.Clear()

	// Create a texture to draw to
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, int32(c.Width), int32(c.Height))
	if err != nil {
		return fmt.Errorf("failed to create texture: %s", err)
	}

	// Initialize joystick
	if sdl.NumJoysticks() > 0 {
		joystick = sdl.JoystickOpen(0)
		defer joystick.Close()
	}

	// Innerloop
	for {
		if !pause {
			if drawFunc != nil {
				if err := drawFunc(c); err != nil {
					return err
				}
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
				return errQuit
			case *sdl.JoyAxisEvent:
				ae := event.(*sdl.JoyAxisEvent)
				if ae.Axis == 0 {
					if ae.Value > 0 {
						// right
						if pressFunc != nil {
							if err := pressFunc(false, true, false, false, false, false, false); err != nil {
								return err
							}
						}
					} else if ae.Value < 0 {
						// left
						if pressFunc != nil {
							if err := pressFunc(true, false, false, false, false, false, false); err != nil {
								return err
							}
						}
					}
				}
				if ae.Axis == 1 {
					if ae.Value > 0 {
						// down
						if pressFunc != nil {
							if err := pressFunc(false, false, false, true, false, false, false); err != nil {
								return err
							}
						}
					} else if ae.Value < 0 {
						// up
						if pressFunc != nil {
							if err := pressFunc(false, false, true, false, false, false, false); err != nil {
								return err
							}
						}
					}
				}
			case *sdl.JoyButtonEvent:
				be := event.(*sdl.JoyButtonEvent)
				if be.Button == 0 {
					if be.State == 1 { // pressed A
						// fire pressed
						if pressFunc != nil {
							if err := pressFunc(false, false, false, false, true, false, false); err != nil {
								return err
							}
						}
					} else if be.State == 0 { // released A
						// fire released
						if releaseFunc != nil {
							if err := releaseFunc(false, false, false, false, true, false, false); err != nil {
								return err
							}
						}
					}
				}
				if be.Button == 1 {
					if be.State == 1 { // pressed B
						// secondary pressed
						if pressFunc != nil {
							if err := pressFunc(false, false, false, false, false, true, false); err != nil {
								return err
							}
						}
					} else if be.State == 0 { // released B
						// secondary released
						if releaseFunc != nil {
							if err := releaseFunc(false, false, false, false, false, true, false); err != nil {
								return err
							}
						}
					}
				}
				if be.Button == 2 {
					if be.State == 1 { // pressed C
						// tertiary pressed
						if pressFunc != nil {
							if err := pressFunc(false, false, false, false, false, false, true); err != nil {
								return err
							}
						}
					} else if be.State == 0 { // released C
						// tertiary released
						if releaseFunc != nil {
							if err := releaseFunc(false, false, false, false, false, false, true); err != nil {
								return err
							}
						}
					}
				}
			case *sdl.KeyboardEvent:
				ke := event.(*sdl.KeyboardEvent)
				if ke.Type == sdl.KEYDOWN {
					ks := ke.Keysym
					switch ks.Sym {
					case sdl.K_ESCAPE, sdl.K_q:
						// quit
						if pressFunc == nil {
							return errQuit
						} else if err := pressFunc(false, false, false, false, false, false, true); err != nil {
							return err
						}
					case sdl.K_SPACE:
						// fire
						if pressFunc != nil {
							if err := pressFunc(false, false, false, false, true, false, false); err != nil {
								return err
							}
						}
					case sdl.K_LEFT, sdl.K_a:
						// left
						if pressFunc != nil {
							if err := pressFunc(true, false, false, false, false, false, false); err != nil {
								return err
							}
						}
					case sdl.K_RIGHT, sdl.K_d:
						// right
						if pressFunc != nil {
							if err := pressFunc(false, true, false, false, false, false, false); err != nil {
								return err
							}
						}
					case sdl.K_UP, sdl.K_w:
						// up
						if pressFunc != nil {
							if err := pressFunc(false, false, true, false, false, false, false); err != nil {
								return err
							}
						}
					case sdl.K_DOWN:
						// down
						if pressFunc != nil {
							if err := pressFunc(false, false, false, true, false, false, false); err != nil {
								return err
							}
						}
					case sdl.K_RETURN:
						altHeldDown := ks.Mod == sdl.KMOD_LALT || ks.Mod == sdl.KMOD_RALT
						if !altHeldDown {
							// alt+enter is not pressed
							// enter is pressed
							if pressFunc != nil {
								if err := pressFunc(false, false, false, false, false, true, false); err != nil {
									return err
								}
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
							if pressFunc != nil {
								if err := pressFunc(false, false, false, true, false, false, false); err != nil {
									return err
								}
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
						if releaseFunc == nil {
							return errQuit
						} else if err := releaseFunc(false, false, false, false, false, false, true); err != nil {
							return err
						}
					case sdl.K_SPACE:
						// fire
						if releaseFunc != nil {
							if err := releaseFunc(false, false, false, false, true, false, false); err != nil {
								return err
							}
						}
					case sdl.K_LEFT, sdl.K_a:
						// left
						if releaseFunc != nil {
							if err := releaseFunc(true, false, false, false, false, false, false); err != nil {
								return err
							}
						}
					case sdl.K_RIGHT, sdl.K_d:
						// right
						if releaseFunc != nil {
							if err := releaseFunc(false, true, false, false, false, false, false); err != nil {
								return err
							}
						}
					case sdl.K_UP, sdl.K_w:
						// up
						if releaseFunc != nil {
							if err := releaseFunc(false, false, true, false, false, false, false); err != nil {
								return err
							}
						}
					case sdl.K_DOWN, sdl.K_s:
						// down
						if releaseFunc != nil {
							if err := releaseFunc(false, false, false, true, false, false, false); err != nil {
								return err
							}
						}
					case sdl.K_RETURN:
						if releaseFunc != nil {
							if err := releaseFunc(false, false, false, false, false, true, false); err != nil {
								return err
							}
						}
					}
				}
			}
		}
		if tickFunc != nil {
			if err := tickFunc(); err != nil {
				return err
			}
		}
		sdl.Delay(uint32(1000 / c.FrameRate))
		loopCounter++
	}
	return nil
}
