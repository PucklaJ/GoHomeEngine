package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/Android"
	// "github.com/PucklaMotzer09/gohomeengine/src/frameworks/GLFW"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	// "github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGL"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"
)

func main() {
	gohome.MainLoop{}.Run(&framework.AndroidFramework{}, &renderer.OpenGLESRenderer{}, 600, 800, "GoHomeEngine", &TestScene2{})
}
