package framework

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/loaders/obj"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
	"io"
	"strings"
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

	onResizeCallbacks []func(newWidth, newHeight uint32)
	onMoveCallbacks   []func(newPosX, newPosY uint32)
	onCloseCallbacks  []func()
	onFocusCallbacks  []func(focused bool)
}

func (this *AndroidFramework) Init(ml *gohome.MainLoop) error {
	this.mainLoop = ml
	this.startOtherThanPaint = time.Now()
	app.Main(androidFrameworkmain)
	return nil
}
func androidFrameworkmain(a app.App) {
	gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_LOG, "Framework", "Android\t", "Starting App ...")

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
		case size.Event:
			androidFramework.onSize(e)
			break
		}
	}

}
func (this *AndroidFramework) onLifecycle(e lifecycle.Event) {
	msg := "Lifecycle: From: " + e.From.String() + " To: " + e.To.String() + " CrossVisible: " + e.Crosses(lifecycle.StageVisible).String() + " CrossFocused: " + e.Crosses(lifecycle.StageFocused).String()
	gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_LOG, "Framework", "Android\t", msg)
	if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
		for i:=0;i<len(this.onFocusCallbacks);i++ {
			this.onFocusCallbacks[i](true)
		}
	} else {
		for i:=0;i<len(this.onFocusCallbacks);i++ {
			this.onFocusCallbacks[i](false)
		}
	}
	if e.Crosses(lifecycle.StageVisible) == lifecycle.CrossOn {
		this.initStuff(e)
	} else if e.Crosses(lifecycle.StageVisible) == lifecycle.CrossOff {
		for i:=0;i<len(this.onCloseCallbacks);i++ {
			this.onCloseCallbacks[i]()
		}
		this.mainLoop.Quit()
		this.renderer.SetOpenGLESContex(nil)
	}
}
func (this *AndroidFramework) initStuff(e lifecycle.Event) {
	context, _ := e.DrawContext.(gl.Context)
	this.renderer.SetOpenGLESContex(context)
	this.mainLoop.InitWindowAndRenderer()
	this.mainLoop.InitManagers()
	gohome.RenderMgr.UpdateProjectionWithViewport = true
	gohome.Render.AfterInit()
	this.mainLoop.SetupStartScene()
	this.appl.Send(paint.Event{})
}
func (this *AndroidFramework) onTouch(e touch.Event) {
	if e.Sequence == 0 {
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

	inputTouch := gohome.InputMgr.Touches[uint8(e.Sequence)]
	inputTouch.ID = uint8(e.Sequence)
	inputTouch.Pos[0] = int16(e.X)
	inputTouch.Pos[1] = int16(e.Y)
	inputTouch.DPos[0] = inputTouch.Pos[0] - inputTouch.PPos[0]
	inputTouch.DPos[1] = inputTouch.Pos[1] - inputTouch.PPos[1]
	inputTouch.PPos = inputTouch.Pos

	if e.Type == touch.TypeBegin {
		gohome.InputMgr.Touch(uint8(e.Sequence))
		inputTouch.DPos = [2]int16{0,0}
	} else if e.Type == touch.TypeEnd {
		inputTouch.DPos = [2]int16{0,0}
		gohome.InputMgr.ReleaseTouch(uint8(e.Sequence))
	}

	gohome.InputMgr.Touches[uint8(e.Sequence)] = inputTouch
}
func (this *AndroidFramework) onPaint(e paint.Event) {
	if (*this.renderer.GetContext()) == nil {
		return
	}
	this.endOtherThanPaint = time.Now()
	gohome.FPSLimit.AddTime(float32(this.endOtherThanPaint.Sub(this.startOtherThanPaint).Seconds()))
	gohome.FPSLimit.StartMeasurement()

	this.mainLoop.InnerLoop()
	this.appl.Send(paint.Event{})

	gohome.FPSLimit.EndMeasurement()
	gohome.FPSLimit.LimitFPS()
	this.startOtherThanPaint = time.Now()
}

func (this *AndroidFramework) onSize(e size.Event) {
	gohome.Render.OnResize(uint32(e.WidthPx), uint32(e.HeightPx))

	for i:=0;i<len(this.onResizeCallbacks);i++ {
		this.onResizeCallbacks[i](uint32(e.WidthPx),uint32(e.HeightPx))
	}
}

func (this *AndroidFramework) Update() {
	gohome.InputMgr.Mouse.DPos[0] = 0
	gohome.InputMgr.Mouse.DPos[1] = 1
}
func (this *AndroidFramework) Terminate() {

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
func (this *AndroidFramework) WindowSetSize(size mgl32.Vec2) {

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

func getFileExtension(file string) string {
	index := strings.LastIndex(file, ".")
	if index == -1 {
		return ""
	}
	return file[index+1:]
}

func equalIgnoreCase(str1, str string) bool {
	if len(str1) != len(str) {
		return false
	}
	for i := 0; i < len(str1); i++ {
		if str1[i] != str[i] {
			if str1[i] >= 65 && str1[i] <= 90 {
				if str[i] >= 97 && str[i] <= 122 {
					if str1[i]+32 != str[i] {
						return false
					}
				} else {
					return false
				}
			} else if str1[i] >= 97 && str1[i] <= 122 {
				if str[i] >= 65 && str[i] <= 90 {
					if str1[i]-32 != str[i] {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

func (this *AndroidFramework) LoadLevel(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	extension := getFileExtension(path)
	if !equalIgnoreCase(extension, "obj") {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Level", name, "Couldn't load file: The file format "+extension+" is not supported")
		return nil
	}
	return loader.LoadLevelOBJ(rsmgr, name, path, preloaded, loadToGPU)
}

func Quit() {
	framew, _ := gohome.Framew.(*AndroidFramework)
	framew.mainLoop.Quit()
}

func (this *AndroidFramework) ShowYesNoDialog(title, message string) uint8 {
	return gohome.DIALOG_CANCELLED
}

func (this *AndroidFramework) OnResize(callback func(newWidth, newHeight uint32)) {
	this.onResizeCallbacks = append(this.onResizeCallbacks,callback)
}
func (this *AndroidFramework) OnMove(callback func(newPosX, newPosY uint32)) {
	this.onMoveCallbacks = append(this.onMoveCallbacks,callback)
}
func (this *AndroidFramework) OnClose(callback func()) {
	this.onCloseCallbacks = append(this.onCloseCallbacks,callback)
}
func (this *AndroidFramework) OnFocus(callback func(focused bool)) {
	this.onFocusCallbacks = append(this.onFocusCallbacks,callback)
}

func (this *AndroidFramework) StartTextInput() {

}

func (this *AndroidFramework) EndTextInput() {

}

func (this *AndroidFramework) GetTextInput() string {
	return ""
}
