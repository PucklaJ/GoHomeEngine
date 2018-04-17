package framework

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		gohome.InputMgr.PressKey(glfwKeysTogohomeKeys(key))
	} else if action == glfw.Release {
		gohome.InputMgr.ReleaseKey(glfwKeysTogohomeKeys(key))
	}
}

func onMousePosChanged(w *glfw.Window, xpos, ypos float64) {
	gohome.InputMgr.Mouse.Pos[0] = int16(xpos)
	gohome.InputMgr.Mouse.Pos[1] = int16(ypos)
	framework, ok := gohome.Framew.(*GLFWFramework)
	if ok {
		gohome.InputMgr.Mouse.DPos[0] = gohome.InputMgr.Mouse.Pos[0] - framework.prevMousePos[0]
		gohome.InputMgr.Mouse.DPos[1] = gohome.InputMgr.Mouse.Pos[1] - framework.prevMousePos[1]
		framework.prevMousePos[0] = gohome.InputMgr.Mouse.Pos[0]
		framework.prevMousePos[1] = gohome.InputMgr.Mouse.Pos[1]
	}

}

func onMouseWheelChanged(w *glfw.Window, xoff, yoff float64) {
	gohome.InputMgr.Mouse.Wheel[0] = int8(xoff)
	gohome.InputMgr.Mouse.Wheel[1] = int8(yoff)
}

func onMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		gohome.InputMgr.PressKey(glfwMouseButtonTogohomeKeys(button))
	} else if action == glfw.Release {
		gohome.InputMgr.ReleaseKey(glfwMouseButtonTogohomeKeys(button))
	}
}

func onResize(window *glfw.Window, width, height int) {
	gohome.Render.OnResize(uint32(width), uint32(height))
}
