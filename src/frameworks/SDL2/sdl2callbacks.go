package framework

import (
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/PucklaJ/go-sdl2/sdl"
)

func setMousePosition(x, y, xrel, yrel int32) {
	gohome.InputMgr.Mouse.Pos[0] = int16(x)
	gohome.InputMgr.Mouse.Pos[1] = int16(y)
	gohome.InputMgr.Mouse.DPos[0] = int16(xrel)
	gohome.InputMgr.Mouse.DPos[1] = int16(yrel)
}

func setTouchPosition(x, y, xrel, yrel int32, touchID sdl.FingerID) {
	inputTouch := gohome.InputMgr.Touches[uint8(touchID)]
	inputTouch.Pos = [2]int16{int16(x), int16(y)}
	inputTouch.DPos = [2]int16{int16(xrel), int16(yrel)}
	inputTouch.PPos[0] = inputTouch.Pos[0] - inputTouch.DPos[0]
	inputTouch.PPos[1] = inputTouch.Pos[1] - inputTouch.DPos[1]

	inputTouch.ID = uint8(touchID)
	gohome.InputMgr.Touches[uint8(touchID)] = inputTouch
}

func (this *SDL2Framework) onMouseMotion(event *sdl.MouseMotionEvent) {
	setMousePosition(event.X, event.Y, event.XRel, event.YRel)
	setTouchPosition(event.X, event.Y, event.XRel, event.YRel, 0)
}

func (this *SDL2Framework) onMouseWheel(event *sdl.MouseWheelEvent) {
	gohome.InputMgr.Mouse.Wheel[0] = int8(event.X)
	gohome.InputMgr.Mouse.Wheel[1] = int8(event.Y)
}

func (this *SDL2Framework) onMouseButton(event *sdl.MouseButtonEvent) {
	if event.GetType() == sdl.MOUSEBUTTONDOWN {
		gohome.InputMgr.PressKey(sdlMouseButtonTogohomeKeys(event.Button))
	} else { // MOUSEBUTTONUP
		gohome.InputMgr.ReleaseKey(sdlMouseButtonTogohomeKeys(event.Button))
	}
}

func (this *SDL2Framework) onKeyEvent(event *sdl.KeyboardEvent) {
	if event.GetType() == sdl.KEYDOWN {
		gohome.InputMgr.PressKey(sdlKeysTogohomeKeys(event.Keysym.Sym))
	} else { // KEYUP
		gohome.InputMgr.ReleaseKey(sdlKeysTogohomeKeys(event.Keysym.Sym))
	}
}

func (this *SDL2Framework) onTextInput(event *sdl.TextInputEvent) {
	for i := 0; i < len(event.Text); i++ {
		if event.Text[i] == '\x00' {
			break
		} else {
			this.textInputBuffer += string(event.Text[i])
		}
	}
}

func (this *SDL2Framework) onWindowEvent(event *sdl.WindowEvent) {
	switch event.Event {
	case sdl.WINDOWEVENT_MOVED:
		for _, c := range this.onMoveCallbacks {
			c(int(event.Data1), int(event.Data2))
		}
	case sdl.WINDOWEVENT_SIZE_CHANGED:
		gohome.Render.OnResize(int(event.Data1), int(event.Data2))
		for _, c := range this.onResizeCallbacks {
			c(int(event.Data1), int(event.Data2))
		}
	case sdl.WINDOWEVENT_FOCUS_GAINED:
		for _, c := range this.onFocusCallbacks {
			c(true)
		}
	case sdl.WINDOWEVENT_FOCUS_LOST:
		for _, c := range this.onFocusCallbacks {
			c(false)
		}
	case sdl.WINDOWEVENT_CLOSE:
		for _, c := range this.onCloseCallbacks {
			c()
		}
	}
}

func (this *SDL2Framework) onTouch(event *sdl.TouchFingerEvent) {
	windowSize := this.WindowGetSize()
	x := int32(event.X * windowSize[0])
	y := int32(event.Y * windowSize[1])
	xrel := int32(event.DX * windowSize[0])
	yrel := int32(event.DY * windowSize[1])
	setTouchPosition(x, y, xrel, yrel, event.FingerID)
	if event.FingerID == 0 {
		setMousePosition(x, y, xrel, yrel)
	}

	inputTouch := gohome.InputMgr.Touches[uint8(event.FingerID)]
	inputTouch.ID = uint8(event.FingerID)

	switch event.Type {
	case sdl.FINGERDOWN:
		gohome.InputMgr.Touch(uint8(event.FingerID))
	case sdl.FINGERUP:
		gohome.InputMgr.ReleaseTouch(uint8(event.FingerID))
	}

	gohome.InputMgr.Touches[uint8(event.FingerID)] = inputTouch
}
