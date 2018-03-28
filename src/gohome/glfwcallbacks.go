package gohome

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		InputMgr.PressKey(glfwKeysTogohomeKeys(key))
	} else if action == glfw.Release {
		InputMgr.ReleaseKey(glfwKeysTogohomeKeys(key))
	}
}

func onMousePosChanged(w *glfw.Window, xpos, ypos float64) {
	InputMgr.Mouse.Pos[0] = int16(xpos)
	InputMgr.Mouse.Pos[1] = int16(ypos)
	framework, ok := Framew.(*GLFWFramework)
	if ok {
		InputMgr.Mouse.DPos[0] = InputMgr.Mouse.Pos[0] - framework.prevMousePos[0]
		InputMgr.Mouse.DPos[1] = InputMgr.Mouse.Pos[1] - framework.prevMousePos[1]
		framework.prevMousePos[0] = InputMgr.Mouse.Pos[0]
		framework.prevMousePos[1] = InputMgr.Mouse.Pos[1]
	}

}

func onMouseWheelChanged(w *glfw.Window, xoff, yoff float64) {
	InputMgr.Mouse.Wheel[0] = int8(xoff)
	InputMgr.Mouse.Wheel[1] = int8(yoff)
}

func onMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		InputMgr.PressKey(glfwMouseButtonTogohomeKeys(button))
	} else if action == glfw.Release {
		InputMgr.ReleaseKey(glfwMouseButtonTogohomeKeys(button))
	}
}

func onResize(window *glfw.Window, width, height int) {
	Render.OnResize(uint32(width), uint32(height))
}
