package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"io"
	"log"
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

	MonitorGetSize() mgl32.Vec2

	CurserShow()
	CursorHide()
	CursorDisable()
	CursorShown() bool
	CursorHidden() bool
	CursorDisabled() bool

	OpenFile(file string) (*File, error)
	LoadLevel(rsmgr *ResourceManager, name, path string, preloaded, loadToGPU bool) *Level
	LoadLevelString(rsmgr *ResourceManager, name, contents, fileName string, preloaded, loadToGPU bool) *Level
	LoadSound(name, path string) Sound
	LoadMusic(name, path string) Music

	ShowYesNoDialog(title, message string) uint8
	Log(message string)

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

type NilFramework struct {
}

func (*NilFramework) Init(ml *MainLoop) error {
	return nil
}
func (*NilFramework) Update() {

}
func (*NilFramework) Terminate() {

}
func (*NilFramework) PollEvents() {

}
func (*NilFramework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	return nil
}
func (*NilFramework) WindowClosed() bool {
	return false
}
func (*NilFramework) WindowSwap() {

}
func (*NilFramework) WindowSetSize(size mgl32.Vec2) {

}
func (*NilFramework) WindowGetSize() mgl32.Vec2 {
	return [2]float32{0.0, 0.0}
}
func (*NilFramework) WindowSetFullscreen(b bool) {

}
func (*NilFramework) WindowIsFullscreen() bool {
	return false
}
func (*NilFramework) MonitorGetSize() mgl32.Vec2 {
	return [2]float32{0.0, 0.0}
}
func (*NilFramework) CurserShow() {

}
func (*NilFramework) CursorHide() {

}
func (*NilFramework) CursorDisable() {

}
func (*NilFramework) CursorShown() bool {
	return true
}
func (*NilFramework) CursorHidden() bool {
	return false
}
func (*NilFramework) CursorDisabled() bool {
	return false
}
func (*NilFramework) OpenFile(file string) (*File, error) {
	return nil, nil
}
func (*NilFramework) LoadLevel(rsmgr *ResourceManager, name, path string, preloaded, loadToGPU bool) *Level {
	return nil
}
func (*NilFramework) LoadLevelString(rsmgr *ResourceManager, name, contents, fileName string, preloaded, loadToGPU bool) *Level {
	return nil
}
func (*NilFramework) LoadSound(name, path string) Sound {
	return &NilSound{}
}
func (*NilFramework) LoadMusic(name, path string) Music {
	return &NilMusic{}
}
func (*NilFramework) ShowYesNoDialog(title, message string) uint8 {
	return DIALOG_ERROR
}

func (*NilFramework) Log(message string) {
	log.Println(message)
}

func (*NilFramework) OnResize(callback func(newWidth, newHeight uint32)) {

}
func (*NilFramework) OnMove(callback func(newPosX, newPosY uint32)) {

}
func (*NilFramework) OnClose(callback func()) {

}
func (*NilFramework) OnFocus(callback func(focused bool)) {

}
func (*NilFramework) StartTextInput() {

}
func (*NilFramework) GetTextInput() string {
	return ""
}
func (*NilFramework) EndTextInput() {

}

func (*NilFramework) GetAudioManager() AudioManager {
	return &NilAudioManager{}
}
