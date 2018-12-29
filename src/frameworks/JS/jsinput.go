package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
)

type buffered_event interface {
	ApplyValues()
}

type keyEvent struct {
	keyCode int
	pressed bool
}

type mouseButtonEvent struct {
	button  int
	pressed bool
}

type mouseMoveEvent struct {
	x int
	y int
}

var buffered_events []buffered_event
var lock_events bool

func addEvent(e buffered_event) {
	for lock_events {

	}
	lock_events = true
	buffered_events = append(buffered_events, e)
	lock_events = false
}

func onKeyDown(event *js.Object) {
	addEvent(
		&keyEvent{
			keyCode: event.Get("keyCode").Int(),
			pressed: true,
		},
	)

	if event.Get("keyCode").Int() == 120 { // F9
		framew.WindowSetFullscreen(!framew.WindowIsFullscreen())
	}
}

func onKeyUp(event *js.Object) {
	addEvent(
		&keyEvent{
			keyCode: event.Get("keyCode").Int(),
			pressed: false,
		},
	)
}

func onMouseButtonDown(event *js.Object) {
	addEvent(
		&mouseButtonEvent{
			button:  event.Get("button").Int(),
			pressed: true,
		},
	)
}

func onMouseButtonUp(event *js.Object) {
	addEvent(
		&mouseButtonEvent{
			button:  event.Get("button").Int(),
			pressed: false,
		},
	)
}

func onMouseMove(event *js.Object) {
	cx := framew.Canvas.Get("left").Int()
	cy := framew.Canvas.Get("top").Int()
	addEvent(
		&mouseMoveEvent{
			x: event.Get("x").Int() - cx,
			y: event.Get("y").Int() - cy,
		},
	)
}

func (this *keyEvent) ApplyValues() {
	if this.pressed {
		gohome.InputMgr.PressKey(jskeyCodeTogohomeKey(this.keyCode))
	} else {
		gohome.InputMgr.ReleaseKey(jskeyCodeTogohomeKey(this.keyCode))
	}
}

func (this *mouseButtonEvent) ApplyValues() {
	if this.pressed {
		gohome.InputMgr.PressKey(jsmouseButtonTogohomeKey(this.button))
	} else {
		gohome.InputMgr.ReleaseKey(jsmouseButtonTogohomeKey(this.button))
	}
}

func (this *mouseMoveEvent) ApplyValues() {
	gohome.InputMgr.Mouse.Pos[0] = int16(this.x)
	gohome.InputMgr.Mouse.Pos[1] = int16(this.y)
}

func addEventListeners() {
	document := js.Global.Get("document")
	document.Call("addEventListener", "keydown", onKeyDown, false)
	document.Call("addEventListener", "keyup", onKeyUp, false)
	document.Call("addEventListener", "mousedown", onMouseButtonDown, false)
	document.Call("addEventListener", "mouseup", onMouseButtonUp, false)
	document.Call("addEventListener", "mousemove", onMouseMove, false)
	js.Global.Call("addEventListener", "beforeunload", onBeforeUnload, false)
}

func onBeforeUnload(event *js.Object) {
	running = false
}

func onResize() {
	nw, nh := uint32(framew.Canvas.Get("width").Int()), uint32(framew.Canvas.Get("height").Int())
	gohome.Render.OnResize(nw, nh)
}
