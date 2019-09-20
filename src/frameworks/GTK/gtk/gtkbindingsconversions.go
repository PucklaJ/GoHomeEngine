package gtk

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"

func (this Window) ToWidget() Widget {
	return Widget{C.windowToWidget(this.Handle)}
}

func (this Box) ToContainer() Container {
	return Container{C.boxToContainer(this.Handle)}
}

func (this Button) ToContainer() Container {
	return Container{C.buttonToContainer(this.Handle)}
}

func (this Box) ToWidget() Widget {
	return Widget{C.boxToWidget(this.Handle)}
}

func (this GLArea) ToWidget() Widget {
	return Widget{C.glareaToWidget(this.Handle)}
}

func (this Button) ToWidget() Widget {
	return Widget{C.buttonToWidget(this.Handle)}
}

func (this GObject) ToWidget() Widget {
	return Widget{C.gobjectToWidget(this.Handle)}
}

func (this GObject) ToListBox() ListBox {
	return ListBox{C.gobjectToListBox(this.Handle)}
}

func (this GObject) ToGLArea() GLArea {
	return GLArea{C.gobjectToGLArea(this.Handle)}
}

func (this GObject) ToToolButton() ToolButton {
	return ToolButton{C.gobjectToToolButton(this.Handle)}
}

func (this GObject) ToMenuItem() MenuItem {
	return MenuItem{C.gobjectToMenuItem(this.Handle)}
}

func (this GObject) ToBox() Box {
	return Box{C.gobjectToBox(this.Handle)}
}

func (this Widget) ToBox() Box {
	return Box{C.widgetToBox(this.Handle)}
}

func (this Window) ToContainer() Container {
	return Container{C.windowToContainer(this.Handle)}
}

func (this Widget) ToWindow() Window {
	return Window{C.widgetToWindow(this.Handle)}
}

func (this Widget) ToContainer() Container {
	return Container{C.widgetToContainer(this.Handle)}
}

func (this Widget) ToGrid() Grid {
	return Grid{C.widgetToGrid(this.Handle)}
}

func (this Widget) ToListBox() ListBox {
	return ListBox{C.widgetToListBox(this.Handle)}
}

func (this Widget) ToLabel() Label {
	return Label{C.widgetToLabel(this.Handle)}
}

func (this GPointer) ToWidget() Widget {
	return Widget{C.gpointerToWidget(this.Handle)}
}

func (this GPointer) ToLabel() Label {
	return this.ToWidget().ToLabel()
}

func (this ListBox) ToWidget() Widget {
	return Widget{C.listBoxToWidget(this.Handle)}
}

func (this ListBox) ToContainer() Container {
	return Container{C.listBoxToContainer(this.Handle)}
}

func (this ListBoxRow) ToWidget() Widget {
	return Widget{C.listBoxRowToWidget(this.Handle)}
}

func (this ListBoxRow) ToContainer() Container {
	return Container{C.listBoxRowToContainer(this.Handle)}
}

func (this Label) ToWidget() Widget {
	return Widget{C.labelToWidget(this.Handle)}
}

func (this Label) ToGObject() GObject {
	return GObject{C.labelToGObject(this.Handle)}
}

func (this MenuItem) ToWidget() Widget {
	return Widget{C.menuItemToWidget(this.Handle)}
}

func (this ToolButton) ToWidget() Widget {
	return Widget{C.toolButtonToWidget(this.Handle)}
}

func (this FileChooserDialog) ToDialog() Dialog {
	return Dialog{C.fileChooserDialogToDialog(this.Handle)}
}

func (this FileChooserDialog) ToWidget() Widget {
	return Widget{C.fileChooserDialogToWidget(this.Handle)}
}

func (this Dialog) ToWidget() Widget {
	return Widget{C.dialogToWidget(this.Handle)}
}

func (this FileChooserDialog) ToFileChooser() FileChooser {
	return FileChooser{C.fileChooserDialogToFileChooser(this.Handle)}
}

func (this Image) ToWidget() Widget {
	return Widget{C.imageToWidget(this.Handle)}
}

func (this MenuBar) ToContainer() Container {
	return Container{C.menuBarToContainer(this.Handle)}
}

func (this MenuBar) ToWidget() Widget {
	return Widget{C.menuBarToWidget(this.Handle)}
}

func (this MenuBar) ToMenuShell() MenuShell {
	return MenuShell{C.menuBarToMenuShell(this.Handle)}
}

func (this Menu) ToMenuShell() MenuShell {
	return MenuShell{C.menuToMenuShell(this.Handle)}
}

func (this Menu) ToWidget() Widget {
	return Widget{C.menuToWidget(this.Handle)}
}

func (this Entry) ToWidget() Widget {
	return Widget{C.entryToWidget(this.Handle)}
}

func (this Entry) ToEditable() Editable {
	return Editable{C.entryToEditable(this.Handle)}
}

func (this Widget) ToEntry() Entry {
	return Entry{C.widgetToEntry(this.Handle)}
}

func (this Event) ToEventKey() EventKey {
	return EventKey{C.eventToEventKey(this.Handle)}
}
