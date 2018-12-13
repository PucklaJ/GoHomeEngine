package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/veandco/go-sdl2/sdl"
)

func (this *SDL2Framework) onMouseMotion(event *sdl.MouseMotionEvent) {
	gohome.InputMgr.Mouse.Pos[0] = int16(event.X)
	gohome.InputMgr.Mouse.Pos[1] = int16(event.Y)
	gohome.InputMgr.Mouse.DPos[0] = int16(event.XRel)
	gohome.InputMgr.Mouse.DPos[1] = int16(event.YRel)

	inputTouch := gohome.InputMgr.Touches[0]
	inputTouch.Pos = gohome.InputMgr.Mouse.Pos
	inputTouch.DPos = gohome.InputMgr.Mouse.DPos
	inputTouch.PPos[0] = gohome.InputMgr.Mouse.Pos[0] - gohome.InputMgr.Mouse.DPos[0]
	inputTouch.PPos[0] = gohome.InputMgr.Mouse.Pos[1] - gohome.InputMgr.Mouse.DPos[1]
	inputTouch.ID = 0
	gohome.InputMgr.Touches[0] = inputTouch
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
			c(uint32(event.Data1), uint32(event.Data2))
		}
	case sdl.WINDOWEVENT_SIZE_CHANGED:
		gohome.Render.OnResize(uint32(event.Data1), uint32(event.Data2))
		for _, c := range this.onResizeCallbacks {
			c(uint32(event.Data1), uint32(event.Data2))
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

	inputTouch := gohome.InputMgr.Touches[uint8(event.FingerID)]
	inputTouch.Pos[0] = int16(event.X * windowSize[0])
	inputTouch.Pos[1] = int16(event.Y * windowSize[1])

	inputTouch.ID = uint8(event.FingerID)

	switch event.Type {
	case sdl.FINGERMOTION:
		inputTouch.DPos[0] = int16(event.DX * windowSize[0])
		inputTouch.DPos[1] = int16(event.DY * windowSize[1])
		inputTouch.PPos[0] = inputTouch.Pos[0] - inputTouch.DPos[0]
		inputTouch.PPos[1] = inputTouch.Pos[1] - inputTouch.DPos[1]
	case sdl.FINGERDOWN:
		gohome.InputMgr.Touch(uint8(event.FingerID))
	case sdl.FINGERUP:
		gohome.InputMgr.ReleaseTouch(uint8(event.FingerID))
	}

	gohome.InputMgr.Touches[uint8(event.FingerID)] = inputTouch

}
