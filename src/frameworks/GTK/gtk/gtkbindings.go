package gtk

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"
import (
	"os"
	"unsafe"

	"github.com/PucklaJ/mathgl/mgl32"
)

func Init() {
	argv := os.Args
	args := len(argv)
	var argvCString []*C.char

	for i := 0; i < args; i++ {
		argvCString = append(argvCString, C.CString(argv[i]))
		defer C.free(unsafe.Pointer(argvCString[i]))
	}
	C.initialise(C.int(args), &argvCString[0])
}

func Main() {
	C.loop()
}

func CreateWindow(windowWidth, windowHeight int, title string) error {
	titleCS := C.CString(title)
	defer C.free(unsafe.Pointer(titleCS))
	if int(C.createWindow(C.int(windowWidth), C.int(windowHeight), titleCS)) == 0 {
		return nil
	}
	return nil
}

func CreateWindowObject() Window {
	return Window{C.createWindowObject()}
}

func CreateGLArea() {
	C.createGLArea()
}

func WindowSetSize(size mgl32.Vec2) {
	C.windowSetSize(C.float(size.X()), C.float(size.Y()))
}

func WindowGetSize() mgl32.Vec2 {
	var v [2]C.float
	var v1 mgl32.Vec2

	C.windowGetSize(&v[0], &v[1])

	v1[0] = float32(v[0])
	v1[1] = float32(v[1])

	return v1
}

func WindowSetFullscreen(b bool) {
	if b {
		C.gtk_window_fullscreen(C.Window)
	} else {
		C.gtk_window_unfullscreen(C.Window)
	}
}

func CursorShow() {
	C.windowShowCursor()
}
func CursorHide() {
	C.windowHideCursor()
}
func CursorDisable() {
	C.windowDisableCursor()
}
func CursorShown() bool {
	return int(C.windowCursorShown()) == 1
}
func CursorHidden() bool {
	return int(C.windowCursorHidden()) == 1
}
func CursorDisabled() bool {
	return int(C.windowCursorDisabled()) == 1
}

func CreateGLAreaAndAddToWindow() {
	C.createGLArea()
	C.addGLAreaToWindow()
}

func CreateGLAreaAndAddToContainer(container Container) {
	C.createGLArea()
	C.addGLAreaToContainer(container.Handle)
}

func GetWindow() Window {
	return Window{C.Window}
}

func SetWindow(window Window) {
	C.Window = window.Handle
}

func GetGLArea() GLArea {
	return GLArea{C.GLarea}
}

func SetGLArea(area GLArea) {
	C.GLarea = area.Handle
}

func MainQuit() {
	C.gtk_main_quit()
}
