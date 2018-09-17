package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	as "github.com/vulkan-go/asche"
	"github.com/vulkan-go/demos/vulkancube"
	vk "github.com/vulkan-go/vulkan"
	"github.com/xlab/closer"
)

const (
	// Window title
	windowTitle = "WORK IN PROGRESS Vulkan Cube"

	// Size of "worldspace pixels", measured in "screenspace pixels"
	pixelscale = 4

	// The resolution (worldspace)
	width  = 256
	height = 256

	// The width of the pixel buffer, used when calculating where to place pixels (y*pitch+x)
	pitch = width

	// Target framerate
	frameRate = 60

	// Alpha value for opaque colors
	opaque = 255
)

func init() {
	runtime.LockOSThread()
	log.SetFlags(log.Lshortfile)
}

type Application struct {
	*vulkancube.SpinningCube
	debugEnabled bool
	windowHandle *sdl.Window
}

func (a *Application) VulkanSurface(instance vk.Instance) (surface vk.Surface) {
	surfPtr, err := a.windowHandle.VulkanCreateSurface(instance)
	if err != nil {
		log.Println("vulkan error:", err)
		return vk.NullSurface
	}
	surf := vk.SurfaceFromPointer(surfPtr)
	return surf
}

func (a *Application) VulkanAppName() string {
	return "VulkanCube"
}

func (a *Application) VulkanLayers() []string {
	return []string{
		// "VK_LAYER_GOOGLE_threading",
		// "VK_LAYER_LUNARG_parameter_validation",
		// "VK_LAYER_LUNARG_object_tracker",
		// "VK_LAYER_LUNARG_core_validation",
		// "VK_LAYER_LUNARG_api_dump",
		// "VK_LAYER_LUNARG_swapchain",
		// "VK_LAYER_GOOGLE_unique_objects",
	}
}

func (a *Application) VulkanDebug() bool {
	return false // a.debugEnabled
}

func (a *Application) VulkanDeviceExtensions() []string {
	return []string{
		"VK_KHR_swapchain",
	}
}

func (a *Application) VulkanSwapchainDimensions() *as.SwapchainDimensions {
	return &as.SwapchainDimensions{
		Width: width * pixelscale, Height: height * pixelscale, Format: vk.FormatB8g8r8a8Unorm,
	}
}

func (a *Application) VulkanInstanceExtensions() []string {
	extensions := a.windowHandle.VulkanGetInstanceExtensions()
	if a.debugEnabled {
		extensions = append(extensions, "VK_EXT_debug_report")
	}
	return extensions
}

func NewApplication(debugEnabled bool) *Application {
	return &Application{
		SpinningCube: vulkancube.NewSpinningCube(1.0),
		debugEnabled: debugEnabled,
	}
}

func run() int {
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL2: %s\n", err)
		return 1
	}
	defer sdl.Quit()

	err = sdl.VulkanLoadLibrary("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load Vulkan library: %s\n", err)
		return 1
	}
	defer sdl.VulkanUnloadLibrary()

	vk.SetGetInstanceProcAddr(sdl.VulkanGetVkGetInstanceProcAddr())
	err = vk.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize Vulkan: %s\n", err)
		return 1
	}
	defer closer.Close()

	app := NewApplication(true)
	reqDim := app.VulkanSwapchainDimensions()
	window, err := sdl.CreateWindow(windowTitle,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(reqDim.Width), int32(reqDim.Height),
		sdl.WINDOW_SHOWN | sdl.WINDOW_VULKAN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	app.windowHandle = window

	// creates a new platform, also initializes Vulkan context in the app
	platform, err := as.NewPlatform(app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize platform: %s\n", err)
		return 1
	}

	dim := app.Context().SwapchainDimensions()
	log.Printf("Initialized %s with %+v swapchain", app.VulkanAppName(), dim)

	// some sync logic
	doneC := make(chan struct{}, 2)
	exitC := make(chan struct{}, 2)
	defer closer.Bind(func() {
		exitC <- struct{}{}
		<-doneC
		log.Println("Bye!")
	})

	fpsDelay := time.Second / 60
	fpsTicker := time.NewTicker(fpsDelay)
	start := time.Now()
	frames := 0
_MainLoop:
	for {
		select {
		case <-exitC:
			fmt.Printf("FPS: %.2f\n", float64(frames)/time.Now().Sub(start).Seconds())
			app.Destroy()
			platform.Destroy()
			window.Destroy()
			fpsTicker.Stop()
			doneC <- struct{}{}
			return 0
		case <-fpsTicker.C:
			frames++
			var event sdl.Event
			for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.KeyboardEvent:
					if t.Keysym.Sym == sdl.K_ESCAPE || t.Keysym.Sym == sdl.K_q {
						exitC <- struct{}{}
						continue _MainLoop
					}
				case *sdl.QuitEvent:
					exitC <- struct{}{}
					continue _MainLoop
				}
			}
			app.NextFrame()

			imageIdx, outdated, err := app.Context().AcquireNextImage()
			if err != nil {
				panic(err)
			}
			if outdated {
				imageIdx, _, err = app.Context().AcquireNextImage()
				if err != nil {
					panic(err)
				}
			}
			//fmt.Printf("%v %T\n", platform.Surface(), platform.Surface())
			_, err = app.Context().PresentImage(imageIdx)
			if err != nil {
				panic(err)
			}
		}
	}
}

func orPanic(vr vk.Result) {
	if err := vk.Error(vr); err != nil {
		panic(err)
	}
}

func main() {
	// This is to allow the deferred functions in run() to kick in at exit
	os.Exit(run())
}
