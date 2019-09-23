package gtk

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

func (this Window) ConfigureParametersAdv(width, height int, title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	C.configureWindowParameters(this.Handle, C.int(width), C.int(height), ctitle)
}

func (this Window) ConfigureParameters() {
	this.ConfigureParametersAdv(0, 0, "")
}

func (this Window) ConnectSignals() {
	C.connectWindowSignals(this.Handle)
}

func (this Window) SetAttachedTo(widget Widget) {
	C.gtk_window_set_attached_to(this.Handle, widget.Handle)
}

func (this Container) Add(widget Widget) {
	C.gtk_container_add(this.Handle, widget.Handle)
	C.gtk_widget_show(widget.Handle)
}

func (this Container) Remove(widget Widget) {
	C.gtk_container_remove(this.Handle, widget.Handle)
}

func (this Container) GetChildren() GList {
	return GList{C.gtk_container_get_children(this.Handle)}
}

func (this GLArea) Configure() {
	C.configureGLArea(this.Handle)
}

func (this Widget) ShowAll() {
	C.gtk_widget_show_all(this.Handle)
}

func (this Widget) Show() {
	C.gtk_widget_show(this.Handle)
}

func (this Widget) Unrealize() {
	C.gtk_widget_unrealize(this.Handle)
}

func (this Widget) SetSizeRequest(width, height int) {
	C.gtk_widget_set_size_request(this.Handle, C.gint(width), C.gint(height))
}

func (this Widget) GetSizeRequest() (int32, int32) {
	var width, height C.gint

	C.gtk_widget_get_size_request(this.Handle, &width, &height)

	return int32(width), int32(height)
}

func (this Widget) GetSize() (int32, int32) {
	var width, height C.gint
	C.widgetGetSize(this.Handle, &width, &height)
	return int32(width), int32(height)
}

func (this Widget) GetParent() Widget {
	return Widget{C.gtk_widget_get_parent(this.Handle)}
}

func (this Widget) IsNULL() bool {
	return this.Handle == nil
}

func (this Widget) Realize() {
	C.gtk_widget_realize(this.Handle)
}

func (this Widget) GetName() string {
	name := C.gtk_widget_get_name(this.Handle)
	return C.GoString(name)
}

func (this Widget) Destroy() {
	C.gtk_widget_destroy(this.Handle)
}

func (this Widget) HasFocus() bool {
	return C.gtk_widget_has_focus(this.Handle) == C.TRUE
}

func (this Widget) SetCanFocus(value bool) {
	C.gtk_widget_set_can_focus(this.Handle, boolTogboolean(value))
}

func (this Widget) GrabFocus() {
	C.gtk_widget_grab_focus(this.Handle)
}

func (this Widget) SetName(name string) {
	namec := C.CString(name)
	defer C.free(unsafe.Pointer(namec))
	C.gtk_widget_set_name(this.Handle, namec)
}

func (this Builder) GetObject(name string) GObject {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))

	handle := C.gtk_builder_get_object(this.Handle, cstr)

	return GObject{handle}
}

func (this Builder) AddFromFile(file string) error {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	if err := C.gtk_builder_add_from_file(this.Handle, cfile, nil); err == 0 {
		return errors.New("Error while loading file")
	}

	return nil
}

func (this GList) Next() GList {
	return GList{this.Handle.next}
}

func (this GList) Prev() GList {
	return GList{this.Handle.prev}
}

func (this GList) Data() GPointer {
	return GPointer{this.Handle.data}
}

func (this GList) Equals(other GList) bool {
	return this.Handle == other.Handle
}

func (this Widget) Equals(other Widget) bool {
	return this.Handle == other.Handle
}

func (this Grid) Attach(child Widget, left, top, width, height int32) {
	C.gtk_grid_attach(this.Handle, child.Handle, C.gint(left), C.gint(top), C.gint(width), C.gint(height))
}

func (this ListBox) Insert(widget Widget, position int32) {
	C.gtk_list_box_insert(this.Handle, widget.Handle, C.gint(position))
}

func (this ListBox) GetSelectedRow() ListBoxRow {
	return ListBoxRow{C.gtk_list_box_get_selected_row(this.Handle)}
}

func (this ListBoxRow) IsNULL() bool {
	return this.Handle == nil
}

func (this Label) SetText(text string) {
	textcs := C.CString(text)
	defer C.free(unsafe.Pointer(textcs))
	C.gtk_label_set_text(this.Handle, textcs)
}

func (this Label) GetText() string {
	textcs := C.gtk_label_get_text(this.Handle)
	return C.GoString(textcs)
}

func (this Dialog) Run() int32 {
	switch C.gtk_dialog_run(this.Handle) {
	case C.GTK_RESPONSE_ACCEPT:
		return RESPONSE_ACCEPT
	default:
		return RESPONSE_NONE
	}
}

func (this FileChooser) SetSelectMultiple(select_multiple bool) {
	var b C.gboolean
	if select_multiple {
		b = C.TRUE
	} else {
		b = C.FALSE
	}
	C.gtk_file_chooser_set_select_multiple(this.Handle, b)
}

func (this FileChooser) GetFilenames() (files []string) {
	list := C.gtk_file_chooser_get_filenames(this.Handle)
	for file := list; file != nil; file = file.next {
		files = append(files, C.GoString(C.gpointerToGChar(file.data)))
	}
	return
}

func (this FileChooser) GetFilename() string {
	filencs := C.gtk_file_chooser_get_filename(this.Handle)
	defer C.free(unsafe.Pointer(filencs))
	filen := C.GoString(filencs)
	return filen
}

func (this FileChooser) SetFilter(filter FileFilter) {
	C.gtk_file_chooser_set_filter(this.Handle, filter.Handle)
}

func (this GObject) SetData(key string, data string) {
	datacs := C.CString(data)
	keycs := C.CString(key)
	C.g_object_set_data(this.Handle, keycs, C.gcharToGPointer(datacs))
	C.free(unsafe.Pointer(keycs))
}

func (this GObject) GetData(key string) string {
	keycs := C.CString(key)
	datacs := C.g_object_get_data(this.Handle, keycs)
	C.free(unsafe.Pointer(keycs))
	return C.GoString(C.gpointerToGChar(datacs))
}

func (this FileFilter) AddPattern(pattern string) {
	cp := C.CString(pattern)
	C.gtk_file_filter_add_pattern(this.Handle, cp)
	C.free(unsafe.Pointer(cp))
}

func (this MenuShell) Append(item MenuItem) {
	C.gtk_menu_shell_append(this.Handle, item.ToWidget().Handle)
}

func (this MenuItem) SetSubmenu(menu Menu) {
	C.gtk_menu_item_set_submenu(this.Handle, menu.ToWidget().Handle)
}

func (this Entry) GetText() string {
	return C.GoString(C.gtk_entry_get_text(this.Handle))
}

func (this Entry) SetText(text string) {
	textc := C.CString(text)
	defer C.free(unsafe.Pointer(textc))
	C.gtk_entry_set_text(this.Handle, textc)
}

func (this Editable) SetEditable(editable bool) {
	C.gtk_editable_set_editable(this.Handle, boolTogboolean(editable))
}

func boolTogboolean(value bool) C.gboolean {
	if value {
		return C.TRUE
	}
	return C.FALSE
}

func (this EventKey) Keyval() gohome.Key {
	keyval := this.Handle.keyval
	return gdkkeysymTogohomekey(keyval)
}

func (this Switch) SetActive(active bool) {
	C.gtk_switch_set_active(this.Handle, boolTogboolean(active))
}

func (this Switch) GetActive() bool {
	return C.gtk_switch_get_active(this.Handle) == C.TRUE
}

func (this Switch) SetState(state bool) {
	C.gtk_switch_set_state(this.Handle, boolTogboolean(state))
}

func (this Switch) GetState() bool {
	return C.gtk_switch_get_state(this.Handle) == C.TRUE
}

func (this SpinButton) SetValue(val float64) {
	C.gtk_spin_button_set_value(this.Handle, C.gdouble(val))
}

func (this SpinButton) GetValue() float64 {
	return float64(C.gtk_spin_button_get_value(this.Handle))
}
