package main

import "C"

import (
	"github.com/veandco/go-sdl2/sdl"
	"runtime"
)

const (
	winTitle  = "GoHomeApp"
	winWidth  = 480
	winHeight = 800
)

var Window *sdl.Window
var Running bool = true

func Init() (err error) {
	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return
	}

	Window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return
	}

	return
}

// Destroy destroys SDL and releases the memory.
func Destroy() {
	Window.Destroy()
	sdl.Quit()
}

// Quit exits main loop.
func Quit() {
	Running = false
}

//export SDL_main
func SDL_main() {
	runtime.LockOSThread()

	err := Init()
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Init: %s\n", err)
	}
	defer Destroy()

	for Running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				Quit()

			case *sdl.KeyboardEvent:
				if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
					Quit()
				}
			}
		}
	}
}

func main() {
	SDL_main()
}
