package framework

import (
	"fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
	"io"
	"time"
)

type AndroidFramework struct {
	appl                app.App
	mainLoop            *gohome.MainLoop
	renderer            *renderer.OpenGLESRenderer
	shouldClose         bool
	prevMousePos        [2]int16
	startOtherThanPaint time.Time
	endOtherThanPaint   time.Time
}

func (this *AndroidFramework) Init(ml *gohome.MainLoop) error {
	this.mainLoop = ml
	this.startOtherThanPaint = time.Now()
	app.Main(androidFrameworkmain)
	return nil
}
func androidFrameworkmain(a app.App) {
	fmt.Println("Starting App ...")

	var androidFramework *AndroidFramework
	androidFramework = gohome.Framew.(*AndroidFramework)
	androidFramework.appl = a
	androidFramework.renderer, _ = gohome.Render.(*renderer.OpenGLESRenderer)
	androidFramework.shouldClose = false
	androidFramework.prevMousePos = [2]int16{0, 0}

	for e := range a.Events() {
		switch e := a.Filter(e).(type) {
		case lifecycle.Event:
			androidFramework.onLifecycle(e)
			break
		case paint.Event:
			androidFramework.onPaint(e)
			break
		case touch.Event:
			androidFramework.onTouch(e)
			break
		}
	}
}
func (this *AndroidFramework) onLifecycle(e lifecycle.Event) {
	fmt.Println("Lifecycle: From:", e.From, "To:", e.To, "CrossVisible:", e.Crosses(lifecycle.StageVisible).String(), "CrossFocused:", e.Crosses(lifecycle.StageFocused).String())
	if e.From == lifecycle.StageDead {
		this.initStuff(e)
	} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
		this.renderer.SetOpenGLESContex(e.DrawContext.(gl.Context))
	}
	if e.To == lifecycle.StageDead {
		this.mainLoop.Quit()
	}
}
func (this *AndroidFramework) initStuff(e lifecycle.Event) {
	context, _ := e.DrawContext.(gl.Context)
	this.renderer.SetOpenGLESContex(context)
	this.mainLoop.InitWindowAndRenderer()
	this.mainLoop.InitManagers()
	this.mainLoop.SetupStartScene()
	fmt.Println("Windowsize:", this.WindowGetSize())
	this.appl.Send(paint.Event{})
}
func (this *AndroidFramework) onTouch(e touch.Event) {
	gohome.InputMgr.Mouse.Pos[0] = int16(e.X)
	gohome.InputMgr.Mouse.Pos[1] = int16(e.Y)
	gohome.InputMgr.Mouse.DPos[0] = this.prevMousePos[0] - int16(e.X)
	gohome.InputMgr.Mouse.DPos[1] = this.prevMousePos[1] - int16(e.Y)
	this.prevMousePos[0] = gohome.InputMgr.Mouse.Pos[0]
	this.prevMousePos[1] = gohome.InputMgr.Mouse.Pos[1]
	if e.Type == touch.TypeBegin {
		gohome.InputMgr.PressKey(gohome.MouseButtonLeft)
	} else if e.Type == touch.TypeEnd {
		gohome.InputMgr.ReleaseKey(gohome.MouseButtonLeft)
	}
}
func (this *AndroidFramework) onPaint(e paint.Event) {
	this.endOtherThanPaint = time.Now()
	gohome.FPSLimit.AddTime(float32(this.endOtherThanPaint.Sub(this.startOtherThanPaint).Seconds()))
	gohome.FPSLimit.StartMeasurement()

	this.mainLoop.InnerLoop()
	this.appl.Send(paint.Event{})

	gohome.FPSLimit.EndMeasurement()
	gohome.FPSLimit.LimitFPS()
	this.startOtherThanPaint = time.Now()
}

func (this *AndroidFramework) Update() {
	gohome.InputMgr.Mouse.DPos[0] = 0
	gohome.InputMgr.Mouse.DPos[1] = 1
}
func (this *AndroidFramework) Terminate() {
	this.appl.Send(lifecycle.Event{
		From: lifecycle.StageFocused,
		To:   lifecycle.StageDead,
	})
}
func (this *AndroidFramework) PollEvents() {

}
func (this *AndroidFramework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	return nil
}
func (this *AndroidFramework) WindowClosed() bool {
	return this.shouldClose
}
func (this *AndroidFramework) WindowSwap() {
	this.appl.Publish()
}
func (this *AndroidFramework) WindowGetSize() mgl32.Vec2 {
	viewport := gohome.Render.GetViewport()
	return mgl32.Vec2{float32(viewport.Width), float32(viewport.Height)}
}
func (this *AndroidFramework) WindowSetFullscreen(b bool) {

}
func (this *AndroidFramework) WindowIsFullscreen() bool {
	return false
}
func (this *AndroidFramework) CurserShow() {

}
func (this *AndroidFramework) CursorHide() {

}
func (this *AndroidFramework) CursorDisable() {

}
func (this *AndroidFramework) CursorShown() bool {
	return true
}
func (this *AndroidFramework) CursorHidden() bool {
	return false
}
func (this *AndroidFramework) CursorDisabled() bool {
	return false
}

func (this *AndroidFramework) OpenFile(file string) (io.ReadCloser, error) {
	return asset.Open(file)
}

func (this *AndroidFramework) LoadLevel(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	return nil
}
