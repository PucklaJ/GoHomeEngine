package framework

import (
	"log"
	"math"
	"os"
	"strings"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/PucklaJ/GoHomeEngine/src/loaders/defaultlevel"
	"github.com/PucklaJ/mathgl/mgl32"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type GLFWFramework struct {
	defaultlevel.Loader

	window           *glfw.Window
	prevMousePos     [2]int16
	prevWindowWidth  int
	prevWindowHeight int
	prevWindowX      int
	prevWindowY      int

	onResizeCallbacks []func(newWidth, newHeight int)
	onMoveCallbacks   []func(newPosX, newPosY int)
	onCloseCallbacks  []func()
	onFocusCallbacks  []func(focused bool)

	textInputStarted bool
	textInputBuffer  string
}

func (gfw *GLFWFramework) Init(ml *gohome.MainLoop) error {
	gfw.textInputStarted = false
	gfw.window = nil
	if err := glfw.Init(); err != nil {
		return err
	}
	ml.DoStuff()

	return nil
}
func (GLFWFramework) Update() {
	gohome.InputMgr.Mouse.Wheel[0] = 0
	gohome.InputMgr.Mouse.Wheel[1] = 0
	gohome.InputMgr.Mouse.DPos[0] = 0
	gohome.InputMgr.Mouse.DPos[1] = 0
}

func (gfw *GLFWFramework) Terminate() {
	defer glfw.Terminate()
	defer gfw.window.Destroy()
}

func setProfile() {
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

func (gfw *GLFWFramework) createWindowProfile(windowWidth, windowHeight int, title string, setprofile bool) error {
	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Resizable, glfw.True)
	if setprofile {
		setProfile()
	}

	glfw.WindowHint(glfw.Samples, 4)

	var err error
	gfw.window, err = glfw.CreateWindow(windowWidth, windowHeight, title, nil, nil)
	if err != nil {
		return err
	}

	gfw.window.MakeContextCurrent()
	gfw.window.SetKeyCallback(onKey)
	gfw.window.SetCursorPosCallback(onMousePosChanged)
	gfw.window.SetScrollCallback(onMouseWheelChanged)
	gfw.window.SetMouseButtonCallback(onMouseButton)
	gfw.window.SetFramebufferSizeCallback(onResize)
	gfw.window.SetPosCallback(onMove)
	gfw.window.SetCloseCallback(onClose)
	gfw.window.SetFocusCallback(onFocus)
	gfw.window.SetCharCallback(onTextInput)

	glfw.SwapInterval(1)
	return nil
}

func (gfw *GLFWFramework) CreateWindow(windowWidth, windowHeight int, title string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recover:", r)
			if err := gfw.createWindowProfile(windowWidth, windowHeight, title, false); err != nil {
				gohome.ErrorMgr.MessageError(gohome.ERROR_LEVEL_FATAL, "WindowCreation", "GLFW", err)
			}
		}
	}()
	if err := gfw.createWindowProfile(windowWidth, windowHeight, title, true); err != nil {
		if strings.Contains(err.Error(), "VersionUnavailable") {
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "WindowCreation", "GLFW", "Couldn't set profile: "+err.Error())
			if err1 := gfw.createWindowProfile(windowWidth, windowHeight, title, false); err1 != nil {
				return err1
			}
		}
	}
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

func (gfw *GLFWFramework) WindowSetSize(size mgl32.Vec2) {
	gfw.window.SetSize(int(size.X()), int(size.Y()))
}

func (gfw *GLFWFramework) WindowGetSize() mgl32.Vec2 {
	var size mgl32.Vec2
	x, y := gfw.window.GetSize()
	size[0] = float32(x)
	size[1] = float32(y)

	return size
}

func glfwKeysTogohomeKeys(key glfw.Key) gohome.Key {
	switch key {
	case glfw.KeyUnknown:
		return gohome.KeyUnknown
	case glfw.KeySpace:
		return gohome.KeySpace
	case glfw.KeyApostrophe:
		return gohome.KeyApostrophe
	case glfw.KeyComma:
		return gohome.KeyComma
	case glfw.KeyMinus:
		return gohome.KeyMinus
	case glfw.KeyPeriod:
		return gohome.KeyPeriod
	case glfw.KeySlash:
		return gohome.KeySlash
	case glfw.Key0:
		return gohome.Key0
	case glfw.Key1:
		return gohome.Key1
	case glfw.Key2:
		return gohome.Key2
	case glfw.Key3:
		return gohome.Key3
	case glfw.Key4:
		return gohome.Key4
	case glfw.Key5:
		return gohome.Key5
	case glfw.Key6:
		return gohome.Key6
	case glfw.Key7:
		return gohome.Key7
	case glfw.Key8:
		return gohome.Key8
	case glfw.Key9:
		return gohome.Key9
	case glfw.KeySemicolon:
		return gohome.KeySemicolon
	case glfw.KeyEqual:
		return gohome.KeyEqual
	case glfw.KeyA:
		return gohome.KeyA
	case glfw.KeyB:
		return gohome.KeyB
	case glfw.KeyC:
		return gohome.KeyC
	case glfw.KeyD:
		return gohome.KeyD
	case glfw.KeyE:
		return gohome.KeyE
	case glfw.KeyF:
		return gohome.KeyF
	case glfw.KeyG:
		return gohome.KeyG
	case glfw.KeyH:
		return gohome.KeyH
	case glfw.KeyI:
		return gohome.KeyI
	case glfw.KeyJ:
		return gohome.KeyJ
	case glfw.KeyK:
		return gohome.KeyK
	case glfw.KeyL:
		return gohome.KeyL
	case glfw.KeyM:
		return gohome.KeyM
	case glfw.KeyN:
		return gohome.KeyN
	case glfw.KeyO:
		return gohome.KeyO
	case glfw.KeyP:
		return gohome.KeyP
	case glfw.KeyQ:
		return gohome.KeyQ
	case glfw.KeyR:
		return gohome.KeyR
	case glfw.KeyS:
		return gohome.KeyS
	case glfw.KeyT:
		return gohome.KeyT
	case glfw.KeyU:
		return gohome.KeyU
	case glfw.KeyV:
		return gohome.KeyV
	case glfw.KeyW:
		return gohome.KeyW
	case glfw.KeyX:
		return gohome.KeyX
	case glfw.KeyY:
		return gohome.KeyY
	case glfw.KeyZ:
		return gohome.KeyZ
	case glfw.KeyLeftBracket:
		return gohome.KeyLeftBracket
	case glfw.KeyBackslash:
		return gohome.KeyBackslash
	case glfw.KeyRightBracket:
		return gohome.KeyRightBracket
	case glfw.KeyGraveAccent:
		return gohome.KeyGraveAccent
	case glfw.KeyWorld1:
		return gohome.KeyWorld1
	case glfw.KeyWorld2:
		return gohome.KeyWorld2
	case glfw.KeyEscape:
		return gohome.KeyEscape
	case glfw.KeyEnter:
		return gohome.KeyEnter
	case glfw.KeyTab:
		return gohome.KeyTab
	case glfw.KeyBackspace:
		return gohome.KeyBackspace
	case glfw.KeyInsert:
		return gohome.KeyInsert
	case glfw.KeyDelete:
		return gohome.KeyDelete
	case glfw.KeyRight:
		return gohome.KeyRight
	case glfw.KeyLeft:
		return gohome.KeyLeft
	case glfw.KeyDown:
		return gohome.KeyDown
	case glfw.KeyUp:
		return gohome.KeyUp
	case glfw.KeyPageUp:
		return gohome.KeyPageUp
	case glfw.KeyPageDown:
		return gohome.KeyPageDown
	case glfw.KeyHome:
		return gohome.KeyHome
	case glfw.KeyEnd:
		return gohome.KeyEnd
	case glfw.KeyCapsLock:
		return gohome.KeyCapsLock
	case glfw.KeyScrollLock:
		return gohome.KeyScrollLock
	case glfw.KeyNumLock:
		return gohome.KeyNumLock
	case glfw.KeyPrintScreen:
		return gohome.KeyPrintScreen
	case glfw.KeyPause:
		return gohome.KeyPause
	case glfw.KeyF1:
		return gohome.KeyF1
	case glfw.KeyF2:
		return gohome.KeyF2
	case glfw.KeyF3:
		return gohome.KeyF3
	case glfw.KeyF4:
		return gohome.KeyF4
	case glfw.KeyF5:
		return gohome.KeyF5
	case glfw.KeyF6:
		return gohome.KeyF6
	case glfw.KeyF7:
		return gohome.KeyF7
	case glfw.KeyF8:
		return gohome.KeyF8
	case glfw.KeyF9:
		return gohome.KeyF9
	case glfw.KeyF10:
		return gohome.KeyF10
	case glfw.KeyF11:
		return gohome.KeyF11
	case glfw.KeyF12:
		return gohome.KeyF12
	case glfw.KeyF13:
		return gohome.KeyF13
	case glfw.KeyF14:
		return gohome.KeyF14
	case glfw.KeyF15:
		return gohome.KeyF15
	case glfw.KeyF16:
		return gohome.KeyF16
	case glfw.KeyF17:
		return gohome.KeyF17
	case glfw.KeyF18:
		return gohome.KeyF18
	case glfw.KeyF19:
		return gohome.KeyF19
	case glfw.KeyF20:
		return gohome.KeyF20
	case glfw.KeyF21:
		return gohome.KeyF21
	case glfw.KeyF22:
		return gohome.KeyF22
	case glfw.KeyF23:
		return gohome.KeyF23
	case glfw.KeyF24:
		return gohome.KeyF24
	case glfw.KeyF25:
		return gohome.KeyF25
	case glfw.KeyKP0:
		return gohome.KeyKP0
	case glfw.KeyKP1:
		return gohome.KeyKP1
	case glfw.KeyKP2:
		return gohome.KeyKP2
	case glfw.KeyKP3:
		return gohome.KeyKP3
	case glfw.KeyKP4:
		return gohome.KeyKP4
	case glfw.KeyKP5:
		return gohome.KeyKP5
	case glfw.KeyKP6:
		return gohome.KeyKP6
	case glfw.KeyKP7:
		return gohome.KeyKP7
	case glfw.KeyKP8:
		return gohome.KeyKP8
	case glfw.KeyKP9:
		return gohome.KeyKP9
	case glfw.KeyKPDecimal:
		return gohome.KeyKPDecimal
	case glfw.KeyKPDivide:
		return gohome.KeyKPDivide
	case glfw.KeyKPMultiply:
		return gohome.KeyKPMultiply
	case glfw.KeyKPSubtract:
		return gohome.KeyKPSubtract
	case glfw.KeyKPAdd:
		return gohome.KeyKPAdd
	case glfw.KeyKPEnter:
		return gohome.KeyKPEnter
	case glfw.KeyKPEqual:
		return gohome.KeyKPEqual
	case glfw.KeyLeftShift:
		return gohome.KeyLeftShift
	case glfw.KeyLeftControl:
		return gohome.KeyLeftControl
	case glfw.KeyLeftAlt:
		return gohome.KeyLeftAlt
	case glfw.KeyLeftSuper:
		return gohome.KeyLeftSuper
	case glfw.KeyRightShift:
		return gohome.KeyRightShift
	case glfw.KeyRightControl:
		return gohome.KeyRightControl
	case glfw.KeyRightAlt:
		return gohome.KeyRightAlt
	case glfw.KeyRightSuper:
		return gohome.KeyRightSuper
	case glfw.KeyMenu:
		return gohome.KeyMenu
	}

	return gohome.KeyUnknown
}

func glfwMouseButtonTogohomeKeys(mb glfw.MouseButton) gohome.Key {
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
		return gohome.MouseButtonLast
	case glfw.MouseButtonLeft:
		return gohome.MouseButtonLeft
	case glfw.MouseButtonRight:
		return gohome.MouseButtonRight
	case glfw.MouseButtonMiddle:
		return gohome.MouseButtonMiddle
	}

	return gohome.MouseButtonLast
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
	}

	return glfw.GetMonitors()[maxIndex]
}

func (gfw *GLFWFramework) WindowIsFullscreen() bool {
	return gfw.window.GetMonitor() != nil
}

func (gfw *GLFWFramework) OpenFile(file string) (gohome.File, error) {
	return os.Open(file)
}

func (gfw *GLFWFramework) ShowYesNoDialog(title, message string) uint8 {
	return gohome.DIALOG_CANCELLED
}

func (gfw *GLFWFramework) OnResize(callback func(newWidth, newHeight int)) {
	gfw.onResizeCallbacks = append(gfw.onResizeCallbacks, callback)
}
func (gfw *GLFWFramework) OnMove(callback func(newPosX, newPosY int)) {
	gfw.onMoveCallbacks = append(gfw.onMoveCallbacks, callback)
}
func (gfw *GLFWFramework) OnClose(callback func()) {
	gfw.onCloseCallbacks = append(gfw.onCloseCallbacks, callback)
}

func (gfw *GLFWFramework) OnFocus(callback func(focused bool)) {
	gfw.onFocusCallbacks = append(gfw.onFocusCallbacks, callback)
}

func (gfw *GLFWFramework) StartTextInput() {
	gfw.textInputStarted = true
}
func (gfw *GLFWFramework) GetTextInput() string {
	str := gfw.textInputBuffer
	gfw.textInputBuffer = ""
	return str
}
func (gfw *GLFWFramework) EndTextInput() {
	gfw.textInputStarted = false
	gfw.textInputBuffer = ""
}

func (gfw *GLFWFramework) MonitorGetSize() mgl32.Vec2 {
	m := getFocusedMonitor(gfw.window)
	vm := m.GetVideoMode()
	if vm != nil {
		return [2]float32{
			float32(vm.Width),
			float32(vm.Height),
		}
	} else {
		return gfw.WindowGetSize()
	}
}
