package gohome

import (
	"fmt"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"io"
	"log"
)

// The differnet results of a dialogue
const (
	DIALOG_YES       = iota
	DIALOG_NO        = iota
	DIALOG_CANCELLED = iota
	DIALOG_ERROR     = iota
)

// An interface consisting of a io.Reader and a io.Closer
type File interface {
	io.Reader
	io.Closer
}

// An interface consisting of a io.ReadSeeker
type FileSeeker interface {
	io.ReadSeeker
}

// The interface that handles everything OS related that's not rendering
type Framework interface {
	// Initialises the framework using MainLoop
	Init(ml *MainLoop) error
	// Update everything
	Update()
	// Terminate everything
	Terminate()
	// Get the events on the window
	PollEvents()

	// Creates a window with the given parameters
	CreateWindow(windowWidth, windowHeight int, title string) error
	// Returns wether the window is closed
	WindowClosed() bool
	// Swaps the back buffer (not that of RenderManager) with the front buffer
	// used for double buffering
	WindowSwap()
	// Sets the size of the window
	WindowSetSize(size mgl32.Vec2)
	// Returns the size of the window
	WindowGetSize() mgl32.Vec2
	// Sets the window to be fullscreen or not
	WindowSetFullscreen(b bool)
	// Returns wether the window is in fullscreen
	WindowIsFullscreen() bool

	// Returns the resolution of the monitor
	MonitorGetSize() mgl32.Vec2

	// Shows the mouse cursor
	CurserShow()
	// Hides the mouse cursor (no locking)
	CursorHide()
	// Hides the mouse cursor (with locking)
	CursorDisable()
	// Returns wether the mouse cursor is shown
	CursorShown() bool
	// Returns wether the mouse cursor is hidden
	CursorHidden() bool
	// Returns wether the mouse cursor is disabled
	CursorDisabled() bool

	// Opens file for reading (uses the framework related functionality
	// on desktop the usual methods can be used without problem)
	OpenFile(file string) (File, error)
	// Loads the level file (.obj)
	LoadLevel(name, path string, loadToGPU bool) *Level
	// Loads the level using contents as the file contents
	LoadLevelString(name, contents, fileName string, loadToGPU bool) *Level

	// Pop ups a dialog with yes and no
	// Returns one enum value
	ShowYesNoDialog(title, message string) uint8
	// Uses framework related logging (logcat on android)
	Log(a ...interface{})

	// Add a function that should be called when the window resizes
	OnResize(callback func(newWidth, newHeight int))
	// Add a function that should be called when the window moves
	OnMove(callback func(newPosX, newPosY int))
	// Add a function that should be called when the window closes
	OnClose(callback func())
	// Add a function that shoul be called when the window gets focused
	OnFocus(callback func(focused bool))

	// Starts the input of text
	StartTextInput()
	// Gets the inputted text
	GetTextInput() string
	// Ends the input of text
	EndTextInput()
}

// The Framework that should be used for everything
var Framew Framework

// An implementation of Framework that does nothing
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
func (*NilFramework) CreateWindow(windowWidth, windowHeight int, title string) error {
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
func (*NilFramework) OpenFile(file string) (File, error) {
	return nil, nil
}
func (*NilFramework) LoadLevel(name, path string, loadToGPU bool) *Level {
	return nil
}
func (*NilFramework) LoadLevelString(name, contents, fileName string, loadToGPU bool) *Level {
	return nil
}
func (*NilFramework) ShowYesNoDialog(title, message string) uint8 {
	return DIALOG_ERROR
}

func (*NilFramework) Log(a ...interface{}) {
	var str = ""
	for _, val := range a {
		str += fmt.Sprint(val) + " "
	}
	log.Println(str[:len(str)-1])
}

func (*NilFramework) OnResize(callback func(newWidth, newHeight int)) {

}
func (*NilFramework) OnMove(callback func(newPosX, newPosY int)) {

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
