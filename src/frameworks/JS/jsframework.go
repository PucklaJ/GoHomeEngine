package framework

import (
	"errors"
	"fmt"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"github.com/gopherjs/gopherjs/js"
	"log"
	"time"
)

type JSFramework struct {
	gohome.NilFramework
	Canvas *js.Object
}

var addtimestart, addtimeend time.Time

func loop(float32) {
	go func() {
		addtimeend = time.Now()
		addtime := float32(addtimeend.Sub(addtimestart).Seconds())
		gohome.FPSLimit.AddTime(addtime)
		gohome.MainLop.LoopOnce()
		addtimestart = time.Now()
		js.Global.Call("requestAnimationFrame", loop)
	}()
}

func (this *JSFramework) Init(ml *gohome.MainLoop) error {
	framew = this
	addEventListeners()
	if !ml.InitWindow() {
		return errors.New("Failed to create Canvas")
	}
	ml.InitRenderer()
	ml.InitManagers()
	ml.SetupStartScene()

	addtimestart = time.Now()
	js.Global.Call("requestAnimationFrame", loop)

	return nil
}

func (*JSFramework) Update() {

}

func (*JSFramework) Terminate() {

}

func (*JSFramework) PollEvents() {
	lock_events = true

	for _, event := range buffered_events {
		event.ApplyValues()
	}
	buffered_events = buffered_events[:0]

	lock_events = false
}

func (this *JSFramework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	document := js.Global.Get("document")
	this.Canvas = document.Call("createElement", "canvas")
	if this.Canvas == nil {
		return errors.New("Failed to create Canvas")
	}
	this.Canvas.Set("width", windowWidth)
	this.Canvas.Set("height", windowHeight)
	body := document.Get("body")
	if body == nil {
		return errors.New("Failed to attach Canvas to body")
	}
	body.Call("appendChild", this.Canvas)
	return nil
}
func (*JSFramework) WindowClosed() bool {
	return false
}

func (*JSFramework) WindowSetSize(size mgl32.Vec2) {

}
func (this *JSFramework) WindowGetSize() mgl32.Vec2 {
	width := float32(this.Canvas.Get("width").Float())
	height := float32(this.Canvas.Get("height").Float())

	return [2]float32{width, height}
}
func (*JSFramework) WindowSetFullscreen(b bool) {

}
func (*JSFramework) WindowIsFullscreen() bool {
	return false
}
func (*JSFramework) MonitorGetSize() mgl32.Vec2 {
	return [2]float32{0.0, 0.0}
}
func (*JSFramework) CurserShow() {

}
func (*JSFramework) CursorHide() {

}
func (*JSFramework) CursorDisable() {

}
func (*JSFramework) CursorShown() bool {
	return true
}
func (*JSFramework) CursorHidden() bool {
	return false
}
func (*JSFramework) CursorDisabled() bool {
	return false
}
func (*JSFramework) OpenFile(file string) (gohome.File, error) {
	return nil, nil
}
func (*JSFramework) LoadLevel(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	return nil
}
func (*JSFramework) LoadLevelString(rsmgr *gohome.ResourceManager, name, contents, fileName string, preloaded, loadToGPU bool) *gohome.Level {
	return nil
}

func (*JSFramework) Log(a ...interface{}) {
	var str = log.Prefix() + " "
	for _, val := range a {
		str += fmt.Sprint(val) + " "
	}
	println(str)
}

var framew *JSFramework
