package gtk

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"
import (
	"github.com/go-gl/mathgl/mgl32"
	"os"
	"unsafe"
)

type Orientation int
type RenderCallback func()
type MotionNotifyCallback func(x, y int16)
type UseWholeWindowCallback func() bool

var OnMotion MotionNotifyCallback
var OnRender RenderCallback
var OnUseWholeScreen UseWholeWindowCallback

const (
	ORIENTATION_HORIZONTAL Orientation = iota
	ORIENTATION_VERTICAL   Orientation = iota
)

type Box struct {
	Handle *C.GtkBox
}

type Window struct {
	Handle *C.GtkWindow
}

type Container struct {
	Handle *C.GtkContainer
}

type Widget struct {
	Handle *C.GtkWidget
}

type GLArea struct {
	Handle *C.GtkGLArea
}

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

func CreateWindow(windowWidth, windowHeight uint32, title string) error {
	titleCS := C.CString(title)
	defer C.free(unsafe.Pointer(titleCS))
	if int(C.createWindow(C.uint(windowWidth), C.uint(windowHeight), titleCS)) == 0 {
		return nil
	}
	return nil
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

func BoxNew(orient Orientation, spacing int) Box {
	var corient C.GtkOrientation
	switch orient {
	case ORIENTATION_HORIZONTAL:
		corient = C.GTK_ORIENTATION_HORIZONTAL
	default:
		corient = C.GTK_ORIENTATION_VERTICAL
	}

	gtkWidget := C.gtk_box_new(corient, C.gint(spacing))
	var this Box
	this.Handle = C.widgetToBox(gtkWidget)
	return this
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

func (this Container) Add(widget Widget) {
	C.gtk_container_add(this.Handle, widget.Handle)
	C.gtk_widget_show(widget.Handle)
}

func (this Box) ToContainer() Container {
	return Container{C.boxToContainer(this.Handle)}
}

func (this Box) ToWidget() Widget {
	return Widget{C.boxToWidget(this.Handle)}
}

func (this GLArea) ToWidget() Widget {
	return Widget{C.glareaToWidget(this.Handle)}
}

func (this Widget) ToBox() Box {
	return Box{C.widgetToBox(this.Handle)}
}

func (this Window) ToContainer() Container {
	return Container{C.windowToContainer(this.Handle)}
}

func GetGLArea() GLArea {
	return GLArea{C.GLarea}
}

func (this Widget) SetSizeRequest(width, height int) {
	C.gtk_widget_set_size_request(this.Handle, C.gint(width), C.gint(height))
}
