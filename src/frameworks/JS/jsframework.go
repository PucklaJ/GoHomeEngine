package framework

import (
	"errors"
	"fmt"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"github.com/gopherjs/gopherjs/js"
	"io"
	"log"
	"net/http"
	"strconv"
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
	addtimestart = time.Now()
	framew = this
	addEventListeners()
	if !ml.InitWindow() {
		return errors.New("Failed to create Canvas")
	}
	ml.InitRenderer()
	ml.InitManagers()
	gohome.Render.AfterInit()
	ml.SetupStartScene()

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
	this.Canvas.Set("id", "gohome_canvas")
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
func (*JSFramework) LoadLevel(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	extension := getFileExtension(path)
	if equalIgnoreCase(extension, "obj") {
		return loadLevelOBJ(rsmgr, name, path, preloaded, loadToGPU)
	}
	gohome.ErrorMgr.Error("Level", name, "The extension "+extension+" is not supported")
	return nil
}
func (*JSFramework) LoadLevelString(rsmgr *gohome.ResourceManager, name, contents, fileName string, preloaded, loadToGPU bool) *gohome.Level {
	return loadLevelOBJString(rsmgr, name, contents, fileName, preloaded, loadToGPU)
}

func (*JSFramework) Log(a ...interface{}) {
	var str = log.Prefix() + " "
	for _, val := range a {
		str += fmt.Sprint(val) + " "
	}
	println(str)
}

var framew *JSFramework
