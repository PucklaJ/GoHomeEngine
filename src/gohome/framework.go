package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"io"
)

type Framework interface {
	Init(ml *MainLoop) error
	Update()
	Terminate()
	PollEvents()

	CreateWindow(windowWidth, windowHeight uint32, title string) error
	WindowClosed() bool
	WindowSwap()
	WindowGetSize() mgl32.Vec2
	WindowSetFullscreen(b bool)
	WindowIsFullscreen() bool

	CurserShow()
	CursorHide()
	CursorDisable()
	CursorShown() bool
	CursorHidden() bool
	CursorDisabled() bool

	OpenFile(file string) (io.ReadCloser, error)
	LoadLevel(rsmgr *ResourceManager, name, path string, preloaded, loadToGPU bool) *Level
}

var Framew Framework
