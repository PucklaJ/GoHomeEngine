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

func PixbufNewFromBytes(data []byte, colorspace Colorspace, has_alpha bool, bits_per_sample, width, height, rowstride int) Pixbuf {
	var csp C.GdkColorspace
	csp = C.GDK_COLORSPACE_RGB
	var alpha C.gboolean
	if has_alpha {
		alpha = C.TRUE
	} else {
		alpha = C.FALSE
	}

	return Pixbuf{
		C.gdk_pixbuf_new_from_bytes(
			C.voidpToGbytes(unsafe.Pointer(&data[0])),
			csp,
			alpha,
			C.int(bits_per_sample),
			C.int(width),
			C.int(height),
			C.int(rowstride),
		),
	}
}

func ImageNewFromPixbuf(pixbuf Pixbuf) Image {
	return Image{C.widgetToImage(C.gtk_image_new_from_pixbuf(pixbuf.Handle))}
}

func MenuItemNewWithLabel(label string) MenuItem {
	labelc := C.CString(label)
	defer C.free(unsafe.Pointer(labelc))
	return MenuItem{C.widgetToMenuItem(C.gtk_menu_item_new_with_label(labelc))}
}

func MenuNew() Menu {
	return Menu{C.widgetToMenu(C.gtk_menu_new())}
}

func MenuBarNew() MenuBar {
	return MenuBar{C.widgetToMenuBar(C.gtk_menu_bar_new())}
}

func EntryNew() Entry {
	return Entry{C.widgetToEntry(C.gtk_entry_new())}
}

func SwitchNew() Switch {
	return Switch{C.widgetToSwitch(C.gtk_switch_new())}
}

func SpinButtonNew(adjustment *Adjustment, climbRate float64, digits uint) SpinButton {
	var cadjustment *C.GtkAdjustment
	if adjustment == nil {
		cadjustment = nil
	} else {
		cadjustment = adjustment.Handle
	}

	return SpinButton{C.widgetToSpinButton(C.gtk_spin_button_new(cadjustment, C.gdouble(climbRate), C.guint(digits)))}
}

func SpinButtonNewWithRange(min, max, step float64) SpinButton {
	return SpinButton{C.widgetToSpinButton(C.gtk_spin_button_new_with_range(C.gdouble(min), C.gdouble(max), C.gdouble(step)))}
}
