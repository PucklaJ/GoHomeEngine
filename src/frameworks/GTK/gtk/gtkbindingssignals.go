package gtk

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"
import (
	"unsafe"
)

type ButtonSignalCallback func(button Button)
type MenuItemSignalCallback func(menuItem MenuItem)

var OnMotion MotionNotifyCallback
var OnRender RenderCallback
var OnUseWholeScreen UseWholeWindowCallback

var buttonSignalCallbacks map[int]map[string]ButtonSignalCallback
var widgetSignalCallbacks map[string]map[string]func(widget Widget)
var widgetEventSignalCallbacks map[string]map[string]func(widget Widget, event Event)
var menuItemSignalCallbacks map[string]map[string]MenuItemSignalCallback
var listBoxRowSelectedSignalCallbacks map[string]map[string]func(listBox ListBox, listBoxRow ListBoxRow)
var toolButtonSignalCallbacks map[string]func(toolButton ToolButton)

func (this Button) SignalConnect(signal string, callback ButtonSignalCallback) {

	if buttonSignalCallbacks == nil {
		buttonSignalCallbacks = make(map[int]map[string]ButtonSignalCallback)
	}
	if buttonSignalCallbacks[this.ID] == nil {
		buttonSignalCallbacks[this.ID] = make(map[string]ButtonSignalCallback)
	}
	var alreadyConnected = false
	if _, ok := buttonSignalCallbacks[this.ID]; ok {
		if _, ok1 := buttonSignalCallbacks[this.ID][signal]; ok1 {
			alreadyConnected = true
		}
	}
	if !alreadyConnected {
		signalcs := C.CString(signal)
		C.signalConnectButton(this.Handle, signalcs, C.int(this.ID))
		C.free(unsafe.Pointer(signalcs))
	}

	buttonSignalCallbacks[this.ID][signal] = callback
}

func (this Widget) signalConnect(signal string, callback func(widget Widget)) {
	name := this.GetName()
	if widgetSignalCallbacks == nil {
		widgetSignalCallbacks = make(map[string]map[string]func(widget Widget))
	}

	if widgetSignalCallbacks[name] == nil {
		widgetSignalCallbacks[name] = make(map[string]func(widget Widget))
	}

	var alreadyConnected = false
	if _, ok := widgetSignalCallbacks[name]; ok {
		if _, ok1 := widgetSignalCallbacks[name][signal]; ok1 {
			alreadyConnected = true
		}
	}
	if !alreadyConnected {
		signalcs := C.CString(signal)
		namecs := C.CString(name)

		if signal == "size-allocate" {
			C.sizeAllocateSignalConnectWidget(this.Handle, signalcs, namecs)
		}

		C.free(unsafe.Pointer(signalcs))
		C.free(unsafe.Pointer(namecs))
	}
	widgetSignalCallbacks[name][signal] = callback
}

func (this Widget) eventSignalConnect(signal string, callback func(widget Widget, event Event)) {
	name := this.GetName()
	if widgetEventSignalCallbacks == nil {
		widgetEventSignalCallbacks = make(map[string]map[string]func(widget Widget, event Event))
	}
	if widgetEventSignalCallbacks[name] == nil {
		widgetEventSignalCallbacks[name] = make(map[string]func(widget Widget, event Event))
	}
	var alreadyConnected1 = false
	if _, ok := widgetEventSignalCallbacks[name]; ok {
		if _, ok1 := widgetEventSignalCallbacks[name][signal]; ok1 {
			alreadyConnected1 = true
		}
	}
	if !alreadyConnected1 {
		signalcs := C.CString(signal)
		namecs := C.CString(name)

		if signal == "" {
		} else {
			C.eventSignalConnectWidget(this.Handle, signalcs, namecs)
		}

		C.free(unsafe.Pointer(signalcs))
		C.free(unsafe.Pointer(namecs))
	}

	widgetEventSignalCallbacks[name][signal] = callback
}

func (this Widget) SignalConnect(signal string, callback interface{}) {
	switch callback.(type) {
	case func(widget Widget):
		this.signalConnect(signal, callback.(func(widget Widget)))
	case func(widget Widget, event Event):
		this.eventSignalConnect(signal, callback.(func(widget Widget, event Event)))
	}
}

func (this MenuItem) SignalConnect(signal string, callback MenuItemSignalCallback) {
	if menuItemSignalCallbacks == nil {
		menuItemSignalCallbacks = make(map[string]map[string]MenuItemSignalCallback)
	}
	name := this.ToWidget().GetName()
	if menuItemSignalCallbacks[name] == nil {
		menuItemSignalCallbacks[name] = make(map[string]MenuItemSignalCallback)
	}
	var alreadyConnected = false
	if _, ok := menuItemSignalCallbacks[name]; ok {
		if _, ok1 := menuItemSignalCallbacks[name][signal]; ok1 {
			alreadyConnected = true
		}
	}
	if !alreadyConnected {
		signalcs := C.CString(signal)
		namecs := C.CString(name)

		if signal == "" {

		} else {
			C.signalConnectMenuItem(this.Handle, signalcs, namecs)
		}

		C.free(unsafe.Pointer(signalcs))
		C.free(unsafe.Pointer(namecs))
	}

	menuItemSignalCallbacks[name][signal] = callback
}

func (this ListBox) rowSelectedSignalConnect(signal string, callback func(listBox ListBox, listBoxRow ListBoxRow)) {
	if listBoxRowSelectedSignalCallbacks == nil {
		listBoxRowSelectedSignalCallbacks = make(map[string]map[string]func(listBox ListBox, listBoxRow ListBoxRow))
	}
	name := this.ToWidget().GetName()
	if listBoxRowSelectedSignalCallbacks[name] == nil {
		listBoxRowSelectedSignalCallbacks[name] = make(map[string]func(listBox ListBox, listBoxRow ListBoxRow))
	}
	var alreadyConnected = false
	if _, ok := listBoxRowSelectedSignalCallbacks[name]; ok {
		if _, ok1 := listBoxRowSelectedSignalCallbacks[name][signal]; ok1 {
			alreadyConnected = true
		}
	}
	if !alreadyConnected {
		signalcs := C.CString(signal)
		namecs := C.CString(name)

		if signal == "" {

		} else {
			C.rowSelectedSignalConnectListBox(this.Handle, signalcs, namecs)
		}

		C.free(unsafe.Pointer(signalcs))
		C.free(unsafe.Pointer(namecs))
	}

	listBoxRowSelectedSignalCallbacks[name][signal] = callback
}

func (this ListBox) SignalConnect(signal string, callback interface{}) {
	switch callback.(type) {
	case func(listBox ListBox, listBoxRow ListBoxRow):
		this.rowSelectedSignalConnect(signal, callback.(func(listBox ListBox, listBoxRow ListBoxRow)))
	}
}

func (this ToolButton) SignalConnect(callback func(toolButton ToolButton)) {
	if toolButtonSignalCallbacks == nil {
		toolButtonSignalCallbacks = make(map[string]func(toolButton ToolButton))
	}
	name := this.ToWidget().GetName()
	var alreadyConnected = false
	if _, ok := toolButtonSignalCallbacks[name]; ok {
		alreadyConnected = true
	}
	if !alreadyConnected {
		namecs := C.CString(name)

		C.signalConnectToolButton(this.Handle, namecs)

		C.free(unsafe.Pointer(namecs))
	}

	toolButtonSignalCallbacks[name] = callback
}
