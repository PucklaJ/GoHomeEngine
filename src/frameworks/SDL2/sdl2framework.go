package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

type SDL2Framework struct {
	gohome.NilFramework

	window  *sdl.Window
	context sdl.GLContext
	running bool
}

func (this *SDL2Framework) Init(ml *gohome.MainLoop) error {
	this.window = nil
	this.running = true
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}
	ml.DoStuff()

	return nil
}

func (this *SDL2Framework) Terminate() {
	defer sdl.Quit()
	defer this.window.Destroy()
	defer sdl.GLDeleteContext(this.context)
}

func (this *SDL2Framework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	var err error
	if this.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(windowWidth), int32(windowWidth), sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL); err != nil {
		return err
	}

	if err1 := sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1); err1 != nil {
		return err1
	}
	if err1 := sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4); err1 != nil {
		return err1
	}

	if this.context, err = this.window.GLCreateContext(); err != nil {
		return err
	}

	return nil
}

func (this *SDL2Framework) PollEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			this.running = false
		}
	}
}

func (this *SDL2Framework) WindowClosed() bool {
	return !this.running
}

func (this *SDL2Framework) WindowSwap() {
	this.window.GLSwap()
}

func (this *SDL2Framework) OpenFile(file string) (*gohome.File, error) {
	gFile := &gohome.File{}
	osFile, err := os.Open(file)
	gFile.ReadSeeker = osFile
	gFile.Closer = osFile
	return gFile, err
}
