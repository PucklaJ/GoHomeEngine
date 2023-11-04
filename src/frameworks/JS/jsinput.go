package framework

import (
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
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
	x  int
	y  int
	dx int
	dy int
}

type mouseWheelEvent struct {
	dx int
	dy int
	dm int
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
	if event.Get("target") != js.Global.Get("document").Get("body") {
		return
	}
	addEvent(
		&keyEvent{
			keyCode: event.Get("keyCode").Int(),
			pressed: true,
		},
	)

	if event.Get("keyCode").Int() == 120 { // F9
		framew.WindowSetFullscreen(!framew.WindowIsFullscreen())
	} else if event.Get("keyCode").Int() == 76 { // L
		if framew.CursorShown() {
			framew.CursorDisable()
		} else {
			framew.CurserShow()
		}
	}

	if event.Get("altKey").Bool() || event.Get("ctrlKey").Bool() ||
		event.Get("shiftKey").Bool() {
		event.Call("preventDefault")
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

var moveX, moveY int

func getMovementX(event *js.Object) int {
	movementX := event.Get("movementX")
	if movementX != js.Undefined && movementX != nil {
		return movementX.Int()
	}
	mozMovementX := event.Get("mozMovementX")
	if mozMovementX != js.Undefined && mozMovementX != nil {
		return mozMovementX.Int()
	}
	webkitMovementX := event.Get("webkitMovementX")
	if webkitMovementX != js.Undefined && webkitMovementX != nil {
		return webkitMovementX.Int()
	}
	return 0
}

func getMovementY(event *js.Object) int {
	movementY := event.Get("movementY")
	if movementY != js.Undefined && movementY != nil {
		return movementY.Int()
	}
	mozMovementY := event.Get("mozMovementY")
	if mozMovementY != js.Undefined && mozMovementY != nil {
		return mozMovementY.Int()
	}
	webkitMovementY := event.Get("webkitMovementY")
	if webkitMovementY != js.Undefined && webkitMovementY != nil {
		return webkitMovementY.Int()
	}
	return 0
}

func onMouseMove(event *js.Object) {
	crect := framew.Canvas.Call("getClientRects").Index(0)
	cx := crect.Get("x").Int()
	cy := crect.Get("y").Int()
	mx := event.Get("x").Int() - cx
	my := event.Get("y").Int() - cy
	moveX += getMovementX(event)
	moveY += getMovementY(event)
	if framew.WindowIsFullscreen() {
		width := crect.Get("width").Float()
		height := crect.Get("height").Float()
		real_width := framew.Canvas.Get("width").Float()
		real_height := framew.Canvas.Get("height").Float()
		var calc_width, calc_height float64

		xrel1 := real_width / width
		yrel1 := real_height / height

		if xrel1 < yrel1 {
			calc_width = real_width / real_height * height
			calc_height = real_height / real_width * calc_width
		} else {
			calc_height = real_height / real_width * width
			calc_width = real_width / real_height * calc_height
		}

		xrel := real_width / calc_width
		yrel := real_height / calc_height

		mx -= int(width-calc_width) / 2
		my -= int(height-calc_height) / 2

		mx = int(float64(mx) * xrel)
		my = int(float64(my) * yrel)
	}
	addEvent(
		&mouseMoveEvent{
			x:  mx,
			y:  my,
			dx: mx - prevMPos[0],
			dy: my - prevMPos[1],
		},
	)
}

func onWheel(event *js.Object) {
	addEvent(
		&mouseWheelEvent{
			dx: event.Get("deltaX").Int(),
			dy: event.Get("deltaY").Int(),
			dm: event.Get("deltaMode").Int(),
		},
	)

	event.Call("preventDefault")
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
		if this.button == 0 {
			gohome.InputMgr.Touch(0)
		}
	} else {
		gohome.InputMgr.ReleaseKey(jsmouseButtonTogohomeKey(this.button))
		if this.button == 0 {
			gohome.InputMgr.ReleaseTouch(0)
		}
	}
}

func (this *mouseMoveEvent) ApplyValues() {
	if framew.CursorDisabled() {
		if moveX != 0 || moveY != 0 {
			gohome.InputMgr.Mouse.Pos[0] += int16(moveX)
			gohome.InputMgr.Mouse.Pos[1] += int16(moveY)
			gohome.InputMgr.Mouse.DPos[0] = int16(moveX)
			gohome.InputMgr.Mouse.DPos[1] = int16(moveY)
			moveX = 0
			moveY = 0
		}
	} else {
		gohome.InputMgr.Mouse.Pos[0] = int16(this.x)
		gohome.InputMgr.Mouse.Pos[1] = int16(this.y)
		gohome.InputMgr.Mouse.DPos[0] = int16(this.dx)
		gohome.InputMgr.Mouse.DPos[1] = int16(this.dy)
		moveX = 0
		moveY = 0
	}

	inputTouch := gohome.InputMgr.Touches[0]
	inputTouch.Pos = gohome.InputMgr.Mouse.Pos
	inputTouch.DPos = gohome.InputMgr.Mouse.DPos
	inputTouch.PPos[0] = int16(prevMPos[0])
	inputTouch.PPos[1] = int16(prevMPos[1])
	inputTouch.ID = 0
	gohome.InputMgr.Touches[0] = inputTouch
}

func (this *mouseWheelEvent) ApplyValues() {
	if this.dm == 0 || this.dm == 1 {
		if this.dx > 0 {
			gohome.InputMgr.Mouse.Wheel[0] = 1
		} else if this.dx < 0 {
			gohome.InputMgr.Mouse.Wheel[0] = -1
		}
		if this.dy > 0 {
			gohome.InputMgr.Mouse.Wheel[1] = 1
		} else if this.dy < 0 {
			gohome.InputMgr.Mouse.Wheel[1] = -1
		}
	}
}

func addEventListeners() {
	document := js.Global.Get("document")
	document.Call("addEventListener", "keydown", onKeyDown, false)
	document.Call("addEventListener", "keyup", onKeyUp, false)
	framew.Canvas.Call("addEventListener", "mousedown", onMouseButtonDown, false)
	document.Call("addEventListener", "mouseup", onMouseButtonUp, false)
	document.Call("addEventListener", "mousemove", onMouseMove, false)
	framew.Canvas.Call("addEventListener", "wheel", onWheel, false)
	js.Global.Call("addEventListener", "beforeunload", onBeforeUnload, false)
	framew.Canvas.Call("addEventListener", "contextmenu", disableContextMenu, false)
}

func onBeforeUnload(event *js.Object) {
	running = false
}

func disableContextMenu(event *js.Object) {
	event.Call("preventDefault")
}

func onResize() {
	nw, nh := framew.Canvas.Get("width").Int(), framew.Canvas.Get("height").Int()
	gohome.Render.OnResize(nw, nh)
}
