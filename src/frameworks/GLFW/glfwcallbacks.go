package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		gohome.InputMgr.PressKey(glfwKeysTogohomeKeys(key))
	} else if action == glfw.Release {
		gohome.InputMgr.ReleaseKey(glfwKeysTogohomeKeys(key))
	}
}

func onMousePosChanged(w *glfw.Window, xpos, ypos float64) {
	gohome.InputMgr.Mouse.Pos[0] = int16(xpos)
	gohome.InputMgr.Mouse.Pos[1] = int16(ypos)
	framework := gohome.Framew.(*GLFWFramework)

	gohome.InputMgr.Mouse.DPos[0] = gohome.InputMgr.Mouse.Pos[0] - framework.prevMousePos[0]
	gohome.InputMgr.Mouse.DPos[1] = gohome.InputMgr.Mouse.Pos[1] - framework.prevMousePos[1]
	framework.prevMousePos[0] = gohome.InputMgr.Mouse.Pos[0]
	framework.prevMousePos[1] = gohome.InputMgr.Mouse.Pos[1]

	inputTouch := gohome.InputMgr.Touches[0]
	inputTouch.Pos = gohome.InputMgr.Mouse.Pos
	inputTouch.DPos = gohome.InputMgr.Mouse.DPos
	inputTouch.PPos = framework.prevMousePos
	inputTouch.ID = 0
	gohome.InputMgr.Touches[0] = inputTouch
}

func onMouseWheelChanged(w *glfw.Window, xoff, yoff float64) {
	gohome.InputMgr.Mouse.Wheel[0] = int8(xoff)
	gohome.InputMgr.Mouse.Wheel[1] = int8(yoff)
}

func onMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		gohome.InputMgr.PressKey(glfwMouseButtonTogohomeKeys(button))
		gohome.InputMgr.Touch(0)
	} else if action == glfw.Release {
		gohome.InputMgr.ReleaseKey(glfwMouseButtonTogohomeKeys(button))
		gohome.InputMgr.ReleaseTouch(0)
	}
}

func onResize(window *glfw.Window, width, height int) {
	gohome.Render.OnResize(width, height)
	gfw := gohome.Framew.(*GLFWFramework)
	for i := 0; i < len(gfw.onResizeCallbacks); i++ {
		gfw.onResizeCallbacks[i](width, height)
	}
}

func onMove(window *glfw.Window, posx, posy int) {
	gfw := gohome.Framew.(*GLFWFramework)

	for i := 0; i < len(gfw.onMoveCallbacks); i++ {
		gfw.onMoveCallbacks[i](posx, posy)
	}
}

func onClose(window *glfw.Window) {
	gfw := gohome.Framew.(*GLFWFramework)

	for i := 0; i < len(gfw.onCloseCallbacks); i++ {
		gfw.onCloseCallbacks[i]()
	}
}

func onFocus(window *glfw.Window, focused bool) {
	gfw := gohome.Framew.(*GLFWFramework)

	for i := 0; i < len(gfw.onFocusCallbacks); i++ {
		gfw.onFocusCallbacks[i](focused)
	}
}

func onTextInput(window *glfw.Window, char rune) {
	gfw := gohome.Framew.(*GLFWFramework)

	gfw.textInputBuffer += string(char)
}
