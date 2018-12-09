package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GLFW"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGL"
)

func main() {
	gohome.MainLop.Run(&framework.GLFWFramework{}, &renderer.OpenGLRenderer{}, 800, 480, "TextRendering", &TextRenderingScene{})
}
