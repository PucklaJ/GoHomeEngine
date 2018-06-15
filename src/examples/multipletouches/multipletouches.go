package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	// "github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGL"
	// "github.com/PucklaMotzer09/gohomeengine/src/frameworks/Android"
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GLFW"
)

func main() {
	gohome.MainLop.Run(&framework.GLFWFramework{}, &renderer.OpenGLRenderer{}, 640, 480, "MultipleTouches", &MultipleTouchesScene{})
}
