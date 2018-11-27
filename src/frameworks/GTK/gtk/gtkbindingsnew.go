package gtk

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"
import (
	"unsafe"
)

func WindowNew() Window {
	return Window{C.widgetToWindow(C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL))}
}

func GLAreaNew() GLArea {
	return GLArea{C.widgetToGLArea(C.gtk_gl_area_new())}
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

var buttonID int = 0

func ButtonNew() Button {
	defer func() {
		buttonID++
	}()
	return Button{C.widgetToButton(C.gtk_button_new()), buttonID}
}

func ButtonNewWithLabel(label string) Button {
	defer func() {
		buttonID++
	}()
	cs := C.CString(label)
	defer C.free(unsafe.Pointer(cs))
	return Button{C.widgetToButton(C.gtk_button_new_with_label(cs)), buttonID}
}

func BuilderNew() Builder {
	return Builder{C.gtk_builder_new()}
}

func GridNew() Grid {
	return Grid{C.widgetToGrid(C.gtk_grid_new())}
}

func ListBoxNew() ListBox {
	return ListBox{C.widgetToListBox(C.gtk_list_box_new())}
}

func LabelNew(text string) Label {
	textcs := C.CString(text)
	defer C.free(unsafe.Pointer(textcs))
	return Label{C.widgetToLabel(C.gtk_label_new(textcs))}
}

func FileChooserDialogNew(title string, parent Window, action FileChooserAction) FileChooserDialog {
	titlecs := C.CString(title)
	defer C.free(unsafe.Pointer(titlecs))

	var faction C.GtkFileChooserAction
	switch action {
	case FILE_CHOOSER_ACTION_OPEN:
		faction = C.GTK_FILE_CHOOSER_ACTION_OPEN
	case FILE_CHOOSER_ACTION_SAVE:
		faction = C.GTK_FILE_CHOOSER_ACTION_SAVE
	case FILE_CHOOSER_ACTION_SELECT_FOLDER:
		faction = C.GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER
	case FILE_CHOOSER_ACTION_CREATE_FOLDER:
		faction = C.GTK_FILE_CHOOSER_ACTION_CREATE_FOLDER
	default:
		faction = C.GTK_FILE_CHOOSER_ACTION_OPEN
	}

	return FileChooserDialog{C.widgetToFileChooserDialog(C.gohome_file_chooser_dialog_new(
		titlecs, parent.Handle, faction,
	))}
}

func FileFilterNew() FileFilter {
	return FileFilter{C.gtk_file_filter_new()}
}
