package framework

import (
	"errors"
	"fmt"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/loaders/defaultlevel"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"github.com/gopherjs/gopherjs/js"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type JSFramework struct {
	defaultlevel.Loader
	Canvas *js.Object
}

var running = false
var prevMPos [2]int

var addtimestart, addtimeend time.Time

func loop(float32) {
	go func() {
		addtimeend = time.Now()
		addtime := float32(addtimeend.Sub(addtimestart).Seconds())
		gohome.FPSLimit.AddTime(addtime)
		gohome.MainLop.LoopOnce()
		addtimestart = time.Now()
		if running {
			js.Global.Call("requestAnimationFrame", loop)
		} else {
			gohome.MainLop.Quit()
		}
	}()
}

func (this *JSFramework) setPointerLockFunctions() {
	if this.Canvas.Get("requestPointerLock") != js.Undefined {
		this.Canvas.Set("requestPointerLock", this.Canvas.Get("requestPointerLock"))
	} else if this.Canvas.Get("mozRequestPointerLock") != js.Undefined {
		this.Canvas.Set("requestPointerLock", this.Canvas.Get("mozRequestPointerLock"))
	} else if this.Canvas.Get("webkitRequestPointerLock") != js.Undefined {
		this.Canvas.Set("requestPointerLock", this.Canvas.Get("webkitRequestPointerLock"))
	}
	document := js.Global.Get("document")
	if document.Get("exitPointerLock") != js.Undefined {
		document.Set("exitPointerLock", document.Get("exitPointerLock"))
	} else if document.Get("mozExitPointerLock") != js.Undefined {
		document.Set("exitPointerLock", document.Get("mozExitPointerLock"))
	} else if document.Get("webkitExitPointerLock") != js.Undefined {
		document.Set("exitPointerLock", document.Get("webkitExitPointerLock"))
	}
}

func (this *JSFramework) Init(ml *gohome.MainLoop) error {
	addtimestart = time.Now()
	framew = this
	if !ml.InitWindow() {
		return errors.New("Failed to create Canvas")
	}
	this.setPointerLockFunctions()
	addEventListeners()
	ml.InitRenderer()
	ml.InitManagers()
	gohome.Render.AfterInit()
	ml.SetupStartScene()
	running = true
	js.Global.Call("requestAnimationFrame", loop)

	return nil
}

func (*JSFramework) PollEvents() {
	lock_events = true

	gohome.InputMgr.Mouse.DPos[0] = 0
	gohome.InputMgr.Mouse.DPos[1] = 0
	gohome.InputMgr.Mouse.Wheel[0] = 0
	gohome.InputMgr.Mouse.Wheel[1] = 0

	for _, event := range buffered_events {
		event.ApplyValues()
	}
	buffered_events = buffered_events[:0]

	prevMPos[0] = int(gohome.InputMgr.Mouse.Pos[0])
	prevMPos[1] = int(gohome.InputMgr.Mouse.Pos[1])

	lock_events = false
}

func (this *JSFramework) CreateWindow(windowWidth, windowHeight int, title string) error {
	document := js.Global.Get("document")
	this.Canvas = document.Call("createElement", "canvas")
	if this.Canvas == nil {
		return errors.New("Failed to create Canvas")
	}
	this.Canvas.Set("width", windowWidth)
	this.Canvas.Set("height", windowHeight)
	this.Canvas.Set("id", "gohome_canvas")
	body := document.Get("body")
	if body == nil {
		return errors.New("Failed to attach Canvas to body")
	}
	body.Call("appendChild", this.Canvas)
	return nil
}

func (this *JSFramework) WindowSetSize(size mgl32.Vec2) {
	this.Canvas.Set("width", size[0])
	this.Canvas.Set("height", size[1])
	canvas_width = int(size[0])
	canvas_height = int(size[1])
	onResize()
}
func (this *JSFramework) WindowGetSize() mgl32.Vec2 {
	width := float32(this.Canvas.Get("width").Float())
	height := float32(this.Canvas.Get("height").Float())

	return [2]float32{width, height}
}

var canvas_width, canvas_height int

func (this *JSFramework) WindowSetFullscreen(b bool) {
	canvas_width = this.Canvas.Get("width").Int()
	canvas_height = this.Canvas.Get("height").Int()
	if b {
		this.Canvas.Call("requestFullscreen")
	} else {
		js.Global.Get("document").Call("exitFullscreen")
	}
	onResize()
}
func (this *JSFramework) WindowIsFullscreen() bool {
	document := js.Global.Get("document")
	fullscreenElement := document.Get("fullscreenElement")
	if fullscreenElement != js.Undefined && fullscreenElement != nil {
		return true
	}
	mozFullscreenElement := document.Get("mozFullScreenElement")
	if mozFullscreenElement != js.Undefined && mozFullscreenElement != nil {
		return true
	}
	webkitFullscreenElement := document.Get("webkitFullscreenElement")
	if webkitFullscreenElement != js.Undefined && webkitFullscreenElement != nil {
		return true
	}

	return false
}

func (this *JSFramework) CurserShow() {
	this.Canvas.Get("style").Set("cursor", "")
	js.Global.Get("document").Call("exitPointerLock")
}
func (this *JSFramework) CursorHide() {
	this.Canvas.Get("style").Set("cursor", "none")
}
func (this *JSFramework) CursorDisable() {
	this.Canvas.Call("requestPointerLock")
}
func (this *JSFramework) CursorShown() bool {
	if this.CursorHidden() {
		return false
	}
	document := js.Global.Get("document")
	pointerLockElement := document.Get("pointerLockElement")
	if pointerLockElement != js.Undefined && pointerLockElement != nil {
		return false
	}
	mozPointerLockElement := document.Get("mozPointerLockElement")
	if mozPointerLockElement != js.Undefined && mozPointerLockElement != nil {
		return false
	}

	return true
}
func (this *JSFramework) CursorHidden() bool {
	cursor := this.Canvas.Get("style").Get("cursor")
	return cursor != js.Undefined && cursor.String() == "none"
}
func (this *JSFramework) CursorDisabled() bool {
	return !this.CursorShown()
}

type JSFile struct {
	io.Reader
	io.Closer
}

func (*JSFramework) OpenFile(file string) (gohome.File, error) {
	resp, err := http.Get(file)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, errors.New("HTTP Request failed: " + strconv.Itoa(resp.StatusCode))
	}

	return resp.Body, nil
}

func (*JSFramework) Log(a ...interface{}) {
	var str = log.Prefix() + " "
	for _, val := range a {
		str += fmt.Sprint(val) + " "
	}
	println(str)
}

var framew *JSFramework
