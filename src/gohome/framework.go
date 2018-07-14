package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"io"
)

const (
	DIALOG_YES       = iota
	DIALOG_NO        = iota
	DIALOG_CANCELLED = iota
	DIALOG_ERROR     = iota
)

type File struct {
	io.ReadSeeker
	io.Closer
}

type Framework interface {
	Init(ml *MainLoop) error
	Update()
	Terminate()
	PollEvents()

	CreateWindow(windowWidth, windowHeight uint32, title string) error
	WindowClosed() bool
	WindowSwap()
	WindowSetSize(size mgl32.Vec2)
	WindowGetSize() mgl32.Vec2
	WindowSetFullscreen(b bool)
	WindowIsFullscreen() bool

	CurserShow()
	CursorHide()
	CursorDisable()
	CursorShown() bool
	CursorHidden() bool
	CursorDisabled() bool

	OpenFile(file string) (*File, error)
	LoadLevel(rsmgr *ResourceManager, name, path string, preloaded, loadToGPU bool) *Level
	LoadSound(name,path string) Sound
	LoadMusic(name,path string) Music

	ShowYesNoDialog(title, message string) uint8

	OnResize(callback func(newWidth, newHeight uint32))
	OnMove(callback func(newPosX, newPosY uint32))
	OnClose(callback func())
	OnFocus(callback func(focused bool))

	StartTextInput()
	GetTextInput() string
	EndTextInput()

	GetAudioManager() AudioManager
}

var Framew Framework
