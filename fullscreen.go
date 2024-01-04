package pixelpusher

import (
	"github.com/veandco/go-sdl2/sdl"
)

// ShowCursor shows the mouse cursor
func ShowCursor() {
	sdl.ShowCursor(1)
}

// HideCursor hides the mouse cursor
func HideCursor() {
	sdl.ShowCursor(0)
}

// Fullscreen checks if the current window has the WINDOW_FULLSCREEN
// or WINDOW_FULLSCREEN_DESKTOP flag set.
func IsFullscreen(window *sdl.Window) bool {
	flags := window.GetFlags()
	window_fullscreen := (flags & sdl.WINDOW_FULLSCREEN) != 0
	window_fullscreen_desktop := (flags & sdl.WINDOW_FULLSCREEN_DESKTOP) != 0
	return window_fullscreen || window_fullscreen_desktop
}

// ToggleFullscreen switches to fullscreen, or back.
// Returns true if the mode has been switched to fullscreen.
// Also toggles the visibility of the mouse cursor.
func ToggleFullscreen(window *sdl.Window) bool {
	if !IsFullscreen(window) {
		// Switch to fullscreen mode
		window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
		HideCursor()
		// Return the new fullscreen status
		return true
		// Returning IsFullscreen(window) is also possible, but unreliable, since it takes a while for the window status to change
	}
	// Switch to windowed mode
	window.SetFullscreen(sdl.WINDOW_SHOWN)
	ShowCursor()
	// Return the new fullscreen status
	return false
	// Returning IsFullscreen(window) is also possible, but unreliable, since it takes a while for the window status to change
}
