package gohome

import (
	// "fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/gl"
	"math"
)

type Framework interface {
	Init(ml *MainLoop) error
	Update()
	Terminate()
	PollEvents()

	CreateWindow(windowWidth, windowHeight uint32, title string) error
	WindowClosed() bool
	WindowSwap()
	WindowGetSize() mgl32.Vec2
	WindowSetFullscreen(b bool)
	WindowIsFullscreen() bool

	CurserShow()
	CursorHide()
	CursorDisable()
	CursorShown() bool
	CursorHidden() bool
	CursorDisabled() bool
}

type GLFWFramework struct {
	window           *glfw.Window
	prevMousePos     [2]int16
	prevWindowWidth  int
	prevWindowHeight int
	prevWindowX      int
	prevWindowY      int
}

func (gfw *GLFWFramework) Init(ml *MainLoop) error {
	gfw.window = nil
	if err := glfw.Init(); err != nil {
		return err
	}
	ml.doStuff()

	return nil
}
func (GLFWFramework) Update() {
	InputMgr.Mouse.Wheel[0] = 0
	InputMgr.Mouse.Wheel[1] = 0
	InputMgr.Mouse.DPos[0] = 0
	InputMgr.Mouse.DPos[1] = 0
}

func (gfw *GLFWFramework) Terminate() {
	defer glfw.Terminate()
	defer gfw.window.Destroy()
}

func (gfw *GLFWFramework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 8)
	var err error
	gfw.window, err = glfw.CreateWindow(int(windowWidth), int(windowHeight), title, nil, nil)
	if err != nil {
		return err
	}
	gfw.window.MakeContextCurrent()
	gfw.window.SetKeyCallback(onKey)
	gfw.window.SetCursorPosCallback(onMousePosChanged)
	gfw.window.SetScrollCallback(onMouseWheelChanged)
	gfw.window.SetMouseButtonCallback(onMouseButton)
	gfw.window.SetFramebufferSizeCallback(onResize)

	glfw.SwapInterval(1)
	return nil
}

func (GLFWFramework) PollEvents() {
	glfw.PollEvents()
}

func (gfw *GLFWFramework) WindowClosed() bool {
	return gfw.window.ShouldClose()
}

func (gfw *GLFWFramework) WindowSwap() {
	gfw.window.SwapBuffers()
}

func (gfw *GLFWFramework) WindowGetSize() mgl32.Vec2 {
	var size mgl32.Vec2
	x, y := gfw.window.GetSize()
	size[0] = float32(x)
	size[1] = float32(y)

	return size
}

func glfwKeysTogohomeKeys(key glfw.Key) Key {
	switch key {
	case glfw.KeyUnknown:
		return KeyUnknown
	case glfw.KeySpace:
		return KeySpace
	case glfw.KeyApostrophe:
		return KeyApostrophe
	case glfw.KeyComma:
		return KeyComma
	case glfw.KeyMinus:
		return KeyMinus
	case glfw.KeyPeriod:
		return KeyPeriod
	case glfw.KeySlash:
		return KeySlash
	case glfw.Key0:
		return Key0
	case glfw.Key1:
		return Key1
	case glfw.Key2:
		return Key2
	case glfw.Key3:
		return Key3
	case glfw.Key4:
		return Key4
	case glfw.Key5:
		return Key5
	case glfw.Key6:
		return Key6
	case glfw.Key7:
		return Key7
	case glfw.Key8:
		return Key8
	case glfw.Key9:
		return Key9
	case glfw.KeySemicolon:
		return KeySemicolon
	case glfw.KeyEqual:
		return KeyEqual
	case glfw.KeyA:
		return KeyA
	case glfw.KeyB:
		return KeyB
	case glfw.KeyC:
		return KeyC
	case glfw.KeyD:
		return KeyD
	case glfw.KeyE:
		return KeyE
	case glfw.KeyF:
		return KeyF
	case glfw.KeyG:
		return KeyG
	case glfw.KeyH:
		return KeyH
	case glfw.KeyI:
		return KeyI
	case glfw.KeyJ:
		return KeyJ
	case glfw.KeyK:
		return KeyK
	case glfw.KeyL:
		return KeyL
	case glfw.KeyM:
		return KeyM
	case glfw.KeyN:
		return KeyN
	case glfw.KeyO:
		return KeyO
	case glfw.KeyP:
		return KeyP
	case glfw.KeyQ:
		return KeyQ
	case glfw.KeyR:
		return KeyR
	case glfw.KeyS:
		return KeyS
	case glfw.KeyT:
		return KeyT
	case glfw.KeyU:
		return KeyU
	case glfw.KeyV:
		return KeyV
	case glfw.KeyW:
		return KeyW
	case glfw.KeyX:
		return KeyX
	case glfw.KeyY:
		return KeyY
	case glfw.KeyZ:
		return KeyZ
	case glfw.KeyLeftBracket:
		return KeyLeftBracket
	case glfw.KeyBackslash:
		return KeyBackslash
	case glfw.KeyRightBracket:
		return KeyRightBracket
	case glfw.KeyGraveAccent:
		return KeyGraveAccent
	case glfw.KeyWorld1:
		return KeyWorld1
	case glfw.KeyWorld2:
		return KeyWorld2
	case glfw.KeyEscape:
		return KeyEscape
	case glfw.KeyEnter:
		return KeyEnter
	case glfw.KeyTab:
		return KeyTab
	case glfw.KeyBackspace:
		return KeyBackspace
	case glfw.KeyInsert:
		return KeyInsert
	case glfw.KeyDelete:
		return KeyDelete
	case glfw.KeyRight:
		return KeyRight
	case glfw.KeyLeft:
		return KeyLeft
	case glfw.KeyDown:
		return KeyDown
	case glfw.KeyUp:
		return KeyUp
	case glfw.KeyPageUp:
		return KeyPageUp
	case glfw.KeyPageDown:
		return KeyPageDown
	case glfw.KeyHome:
		return KeyHome
	case glfw.KeyEnd:
		return KeyEnd
	case glfw.KeyCapsLock:
		return KeyCapsLock
	case glfw.KeyScrollLock:
		return KeyScrollLock
	case glfw.KeyNumLock:
		return KeyNumLock
	case glfw.KeyPrintScreen:
		return KeyPrintScreen
	case glfw.KeyPause:
		return KeyPause
	case glfw.KeyF1:
		return KeyF1
	case glfw.KeyF2:
		return KeyF2
	case glfw.KeyF3:
		return KeyF3
	case glfw.KeyF4:
		return KeyF4
	case glfw.KeyF5:
		return KeyF5
	case glfw.KeyF6:
		return KeyF6
	case glfw.KeyF7:
		return KeyF7
	case glfw.KeyF8:
		return KeyF8
	case glfw.KeyF9:
		return KeyF9
	case glfw.KeyF10:
		return KeyF10
	case glfw.KeyF11:
		return KeyF11
	case glfw.KeyF12:
		return KeyF12
	case glfw.KeyF13:
		return KeyF13
	case glfw.KeyF14:
		return KeyF14
	case glfw.KeyF15:
		return KeyF15
	case glfw.KeyF16:
		return KeyF16
	case glfw.KeyF17:
		return KeyF17
	case glfw.KeyF18:
		return KeyF18
	case glfw.KeyF19:
		return KeyF19
	case glfw.KeyF20:
		return KeyF20
	case glfw.KeyF21:
		return KeyF21
	case glfw.KeyF22:
		return KeyF22
	case glfw.KeyF23:
		return KeyF23
	case glfw.KeyF24:
		return KeyF24
	case glfw.KeyF25:
		return KeyF25
	case glfw.KeyKP0:
		return KeyKP0
	case glfw.KeyKP1:
		return KeyKP1
	case glfw.KeyKP2:
		return KeyKP2
	case glfw.KeyKP3:
		return KeyKP3
	case glfw.KeyKP4:
		return KeyKP4
	case glfw.KeyKP5:
		return KeyKP5
	case glfw.KeyKP6:
		return KeyKP6
	case glfw.KeyKP7:
		return KeyKP7
	case glfw.KeyKP8:
		return KeyKP8
	case glfw.KeyKP9:
		return KeyKP9
	case glfw.KeyKPDecimal:
		return KeyKPDecimal
	case glfw.KeyKPDivide:
		return KeyKPDivide
	case glfw.KeyKPMultiply:
		return KeyKPMultiply
	case glfw.KeyKPSubtract:
		return KeyKPSubtract
	case glfw.KeyKPAdd:
		return KeyKPAdd
	case glfw.KeyKPEnter:
		return KeyKPEnter
	case glfw.KeyKPEqual:
		return KeyKPEqual
	case glfw.KeyLeftShift:
		return KeyLeftShift
	case glfw.KeyLeftControl:
		return KeyLeftControl
	case glfw.KeyLeftAlt:
		return KeyLeftAlt
	case glfw.KeyLeftSuper:
		return KeyLeftSuper
	case glfw.KeyRightShift:
		return KeyRightShift
	case glfw.KeyRightControl:
		return KeyRightControl
	case glfw.KeyRightAlt:
		return KeyRightAlt
	case glfw.KeyRightSuper:
		return KeyRightSuper
	case glfw.KeyMenu:
		return KeyMenu
	}

	return KeyUnknown
}

func glfwMouseButtonTogohomeKeys(mb glfw.MouseButton) Key {
	switch mb {
	// case glfw.MouseButton1:
	// 	return MouseButton1
	// case glfw.MouseButton2:
	// 	return MouseButton2
	// case glfw.MouseButton3:
	// 	return MouseButton3
	// case glfw.MouseButton4:
	// 	return MouseButton4
	// case glfw.MouseButton5:
	// 	return MouseButton5
	// case glfw.MouseButton6:
	// 	return MouseButton6
	// case glfw.MouseButton7:
	// 	return MouseButton7
	// case glfw.MouseButton8:
	// 	return MouseButton8
	case glfw.MouseButtonLast:
		return MouseButtonLast
	case glfw.MouseButtonLeft:
		return MouseButtonLeft
	case glfw.MouseButtonRight:
		return MouseButtonRight
	case glfw.MouseButtonMiddle:
		return MouseButtonMiddle
	}

	return MouseButtonLast
}

func (gfw *GLFWFramework) CurserShow() {
	gfw.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}
func (gfw *GLFWFramework) CursorHide() {
	gfw.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}
func (gfw *GLFWFramework) CursorDisable() {
	gfw.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}
func (gfw *GLFWFramework) CursorShown() bool {
	return gfw.window.GetInputMode(glfw.CursorMode) == glfw.CursorNormal
}
func (gfw *GLFWFramework) CursorHidden() bool {
	return gfw.window.GetInputMode(glfw.CursorMode) == glfw.CursorHidden
}
func (gfw *GLFWFramework) CursorDisabled() bool {
	return gfw.window.GetInputMode(glfw.CursorMode) == glfw.CursorDisabled
}
func (gfw *GLFWFramework) WindowSetFullscreen(b bool) {
	var monitor *glfw.Monitor
	var refreshRate int
	var x, y int
	var width, height int
	if b {
		monitor = getFocusedMonitor(gfw.window)
		if monitor == nil {
			monitor = glfw.GetPrimaryMonitor()
		}
		refreshRate = monitor.GetVideoMode().RefreshRate
		x = 0
		y = 0
		gfw.prevWindowWidth, gfw.prevWindowHeight = gfw.window.GetSize()
		gfw.prevWindowX, gfw.prevWindowY = gfw.window.GetPos()
	} else {
		monitor = nil
		temp := getFocusedMonitor(gfw.window)
		if temp == nil {
			temp = glfw.GetPrimaryMonitor()
		}
		refreshRate = temp.GetVideoMode().RefreshRate
		x, y = gfw.prevWindowX, gfw.prevWindowY
	}
	width, height = gfw.prevWindowWidth, gfw.prevWindowHeight

	gfw.window.SetMonitor(monitor, x, y, width, height, refreshRate)
}

func getFocusedMonitor(window *glfw.Window) *glfw.Monitor {
	wx, wy := window.GetPos()
	ww, wh := window.GetSize()

	var mx, my, mw, mh int
	var monitor *glfw.Monitor
	var max, maxIndex int = -1, -1
	var axmin, axmax, aymin, aymax int
	var bxmin, bxmax, bymin, bymax int

	for i := 0; i < len(glfw.GetMonitors()); i++ {
		monitor = glfw.GetMonitors()[i]
		mx, my = monitor.GetPos()
		mw, mh = monitor.GetVideoMode().Width, monitor.GetVideoMode().Height

		if wx+ww > mx && wx < mx+mw &&
			wy+wh > my && wy < my+mh {

			axmin, axmax = wx, wx+ww
			aymin, aymax = wy, wy+wh
			bxmin, bxmax = mx, mx+mw
			bymin, bymax = my, my+mh

			dx := int(math.Min(float64(axmax), float64(bxmax)) - math.Max(float64(axmin), float64(bxmin)))
			dy := int(math.Min(float64(aymax), float64(bymax)) - math.Max(float64(aymin), float64(bymin)))

			mean := dx * dy
			if mean > max {
				max = mean
				maxIndex = i
			}
		}
	}

	if maxIndex == -1 {
		return nil
	} else {
		return glfw.GetMonitors()[maxIndex]
	}

}

func (gfw *GLFWFramework) WindowIsFullscreen() bool {
	return gfw.window.GetMonitor() != nil
}

type AndroidFramework struct {
	appl   app.App
	glcotx gl.Context
}

func (this *AndroidFramework) Init(ml *MainLoop) error {
	app.Main(androidFrameworkmain)
	return nil
}
func androidFrameworkmain(a app.App) {
	var androidFramework *AndroidFramework
	androidFramework = Framew.(*AndroidFramework)
	androidFramework.appl = a

	for e := range a.Events() {
		switch e := a.Filter(e).(type) {
		case lifecycle.Event:
			switch e.Crosses(lifecycle.StageVisible) {
			case lifecycle.CrossOn:
				androidFramework.glcotx, _ = e.DrawContext.(gl.Context)
				a.Send(paint.Event{})
				break
			}
			break
		case paint.Event:
			androidFramework.glcotx.ClearColor(1.0, 0.0, 0.0, 1.0)
			androidFramework.glcotx.Clear(gl.COLOR_BUFFER_BIT)
			a.Publish()
			break
		}
	}
}
func (this *AndroidFramework) Update() {

}
func (this *AndroidFramework) Terminate() {

}
func (this *AndroidFramework) PollEvents() {

}
func (this *AndroidFramework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	return nil
}
func (this *AndroidFramework) WindowClosed() bool {
	return false
}
func (this *AndroidFramework) WindowSwap() {

}
func (this *AndroidFramework) WindowGetSize() mgl32.Vec2 {
	return mgl32.Vec2{0.0, 0.0}
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

var Framew Framework
