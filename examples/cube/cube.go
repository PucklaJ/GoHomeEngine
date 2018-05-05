package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GLFW"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGL"
)

func main() {
	gohome.MainLoop{}.Run(&framework.GLFWFramework{}, &renderer.OpenGLRenderer{}, 640, 480, "Cube", &CubeScene{})
}
