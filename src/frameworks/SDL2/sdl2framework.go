package framework

import (
	"fmt"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/loaders/defaultlevel"
	"github.com/PucklaMotzer09/go-sdl2/sdl"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"os"
	"runtime"
)

type SDL2Framework struct {
	defaultlevel.Loader

	window          *sdl.Window
	context         sdl.GLContext
	running         bool
	textInputBuffer string

	onResizeCallbacks []func(newWidth, newHeight int)
	onMoveCallbacks   []func(newPosX, newPosY int)
	onCloseCallbacks  []func()
	onFocusCallbacks  []func(focused bool)
}

func (this *SDL2Framework) Init(ml *gohome.MainLoop) error {
	this.window = nil
	this.running = true
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		return err
	}
	ml.DoStuff()

	return nil
}

func (this *SDL2Framework) Update() {
	gohome.InputMgr.Mouse.Wheel[0] = 0
	gohome.InputMgr.Mouse.Wheel[1] = 0
	gohome.InputMgr.Mouse.DPos[0] = 0
	gohome.InputMgr.Mouse.DPos[1] = 0
}

func (this *SDL2Framework) Terminate() {
	defer sdl.Quit()
	defer this.window.Destroy()
	defer sdl.GLDeleteContext(this.context)
}

func setGLAttributesNormal() error {
	if err1 := setGLAttributesCompatible(); err1 != nil {
		return err1
	}

	return nil
}

func setGLAttributesCompatible() error {
	if err1 := sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1); err1 != nil {
		return err1
	}
	return nil
}

func setGLAttributesProfile() error {
	if runtime.GOOS != "android" {
		if err1 := sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE); err1 != nil {
			return err1
		}
		if err1 := sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4); err1 != nil {
			return err1
		}
		if err1 := sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 5); err1 != nil {
			return err1
		}
	}

	return nil
}

func getWindowCreationFlags() uint32 {
	flags := sdl.WINDOW_SHOWN | sdl.WINDOW_OPENGL
	if runtime.GOOS == "android" {
		flags |= sdl.WINDOW_FULLSCREEN
	} else {
		flags |= sdl.WINDOW_RESIZABLE
	}

	return uint32(flags)
}

func (this *SDL2Framework) createWindowLight(windowWidth, windowHeight int, title string) error {
	sdl.GLResetAttributes()
	if err := setGLAttributesCompatible(); err != nil {
		return err
	}
	var err error
	if this.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(windowWidth), int32(windowHeight), getWindowCreationFlags()); err != nil {
		return err
	}
	return nil
}

func (this *SDL2Framework) CreateWindow(windowWidth, windowHeight int, title string) error {

	if err := setGLAttributesNormal(); err != nil {
		return err
	}
	if err := setGLAttributesProfile(); err != nil {
		return err
	}

	var err error
	if this.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(windowWidth), int32(windowHeight), getWindowCreationFlags()); err != nil {
		if err = this.createWindowLight(windowWidth, windowHeight, title); err != nil {
			return err
		}
	}

	if this.context, err = this.window.GLCreateContext(); err != nil {
		this.window.Destroy()
		if err = this.createWindowLight(windowWidth, windowHeight, title); err != nil {
			return err
		}
		if this.context, err = this.window.GLCreateContext(); err != nil {
			this.window.Destroy()
			sdl.Quit()
			if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
				return err
			}
			if err = this.createWindowLight(windowWidth, windowHeight, title); err != nil {
				return err
			}
			if this.context, err = this.window.GLCreateContext(); err != nil {
				return err
			}
		}
	}

	if err1 := sdl.GLSetSwapInterval(1); err1 != nil {
		return err1
	}

	return nil
}

func (this *SDL2Framework) PollEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			this.running = false
		case *sdl.MouseMotionEvent:
			this.onMouseMotion(t)
		case *sdl.MouseButtonEvent:
			this.onMouseButton(t)
		case *sdl.MouseWheelEvent:
			this.onMouseWheel(t)
		case *sdl.KeyboardEvent:
			this.onKeyEvent(t)
		case *sdl.TextInputEvent:
			this.onTextInput(t)
		case *sdl.WindowEvent:
			this.onWindowEvent(t)
		case *sdl.TouchFingerEvent:
			this.onTouch(t)
		}
	}
}

func (this *SDL2Framework) WindowClosed() bool {
	return !this.running
}

func (this *SDL2Framework) WindowSwap() {
	this.window.GLSwap()
}

func (this *SDL2Framework) WindowSetSize(size mgl32.Vec2) {
	this.window.SetSize(int32(size.X()), int32(size.Y()))
}

func (this *SDL2Framework) WindowGetSize() mgl32.Vec2 {
	var size mgl32.Vec2
	x, y := this.window.GetSize()
	size[0] = float32(x)
	size[1] = float32(y)

	return size
}

func (this *SDL2Framework) CurserShow() {
	sdl.ShowCursor(sdl.ENABLE)
	sdl.SetRelativeMouseMode(false)
}
func (this *SDL2Framework) CursorHide() {
	sdl.ShowCursor(sdl.DISABLE)
	sdl.SetRelativeMouseMode(false)
}
func (this *SDL2Framework) CursorDisable() {
	sdl.SetRelativeMouseMode(true)
}
func (this *SDL2Framework) CursorShown() bool {
	val, err := sdl.ShowCursor(sdl.QUERY)
	return err != nil && val == sdl.DISABLE
}
func (this *SDL2Framework) CursorHidden() bool {
	return !this.CursorShown()
}
func (this *SDL2Framework) CursorDisabled() bool {
	return sdl.GetRelativeMouseMode()
}

func (this *SDL2Framework) WindowSetFullscreen(b bool) {
	var flags uint32
	if b {
		flags = sdl.WINDOW_FULLSCREEN_DESKTOP
	} else {
		flags = 0
	}
	if err := this.window.SetFullscreen(flags); err != nil {
		gohome.ErrorMgr.Error("Framework", "SDL2", err.Error())
	}
}

func (this *SDL2Framework) WindowIsFullscreen() bool {
	flags := this.window.GetFlags()
	return (flags&sdl.WINDOW_FULLSCREEN | sdl.WINDOW_FULLSCREEN_DESKTOP) != 0
}

func (this *SDL2Framework) OpenFile(file string) (gohome.File, error) {
	if runtime.GOOS == "android" {
		rw := sdl.RWFromFile(file, "rb")
		var err error
		if rw == nil {
			err = sdl.GetError()
		}
		return rw, err
	} else {
		return os.Open(file)
	}
}

func (this *SDL2Framework) ShowYesNoDialog(title, message string) uint8 {
	var data sdl.MessageBoxData
	data.Title = title
	data.Message = message
	data.Buttons = []sdl.MessageBoxButtonData{
		sdl.MessageBoxButtonData{
			ButtonID: 2,
			Text:     "Yes",
		},
		sdl.MessageBoxButtonData{
			ButtonID: 3,
			Text:     "No",
		},
	}
	data.NumButtons = 2
	data.Window = this.window
	data.Flags = sdl.MESSAGEBOX_INFORMATION
	id, err := sdl.ShowMessageBox(&data)
	if err == nil {
		switch id {
		case 2:
			return gohome.DIALOG_YES
		case 3:
			return gohome.DIALOG_NO
		default:
			return gohome.DIALOG_CANCELLED
		}
	} else {
		return gohome.DIALOG_ERROR
	}
}

func (*SDL2Framework) Log(a ...interface{}) {
	var str = ""
	for _, val := range a {
		str += fmt.Sprint(val) + " "
	}
	sdl.Log(str[:len(str)-1])
}

func (this *SDL2Framework) OnResize(callback func(newWidth, newHeight int)) {
	this.onResizeCallbacks = append(this.onResizeCallbacks, callback)
}
func (this *SDL2Framework) OnMove(callback func(newPosX, newPosY int)) {
	this.onMoveCallbacks = append(this.onMoveCallbacks, callback)
}
func (this *SDL2Framework) OnClose(callback func()) {
	this.onCloseCallbacks = append(this.onCloseCallbacks, callback)
}

func (this *SDL2Framework) OnFocus(callback func(focused bool)) {
	this.onFocusCallbacks = append(this.onFocusCallbacks, callback)
}

func (this *SDL2Framework) StartTextInput() {
	sdl.StartTextInput()
}

func (this *SDL2Framework) GetTextInput() string {
	text := this.textInputBuffer
	this.textInputBuffer = ""
	return text
}

func (this *SDL2Framework) EndTextInput() {
	sdl.StopTextInput()
}

func (this *SDL2Framework) MonitorGetSize() mgl32.Vec2 {
	return this.WindowGetSize()
}
