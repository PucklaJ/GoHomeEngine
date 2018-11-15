package gtk

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"

type Orientation int
type RenderCallback func()
type MotionNotifyCallback func(x, y int16)
type UseWholeWindowCallback func() bool

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

type Button struct {
	Handle *C.GtkButton
	ID     int
}

type GObject struct {
	Handle *C.GObject
}

type Builder struct {
	Handle *C.GtkBuilder
}

type GList struct {
	Handle *C.GList
}

type GPointer struct {
	Handle C.gpointer
}

type Grid struct {
	Handle *C.GtkGrid
}

type ListBox struct {
	Handle *C.GtkListBox
}

type Label struct {
	Handle *C.GtkLabel
}

type MenuItem struct {
	Handle *C.GtkMenuItem
}

type Event struct {
	Handle *C.GdkEvent
}

type ListBoxRow struct {
	Handle *C.GtkListBoxRow
}

type ToolButton struct {
	Handle *C.GtkToolButton
}
