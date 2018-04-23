package main

import (
	// "github.com/PucklaMotzer09/gohomeengine/src/frameworks/Android"
	// "fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GLFW"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	// "github.com/PucklaMotzer09/gohomeengine/src/loaders/assimp"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGL"
	// "github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"
	// "github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	gohome.MainLoop{}.Run(&framework.GLFWFramework{}, &renderer.OpenGLRenderer{}, 1280, 720, "GoHomeEngine", &TestScene{})
}
