package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK"
	// "github.com/PucklaMotzer09/gohomeengine/src/frameworks/Android"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGL"
	// "github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"
)

func main() {
	gohome.MainLop.Run(&framework.GTKFramework{}, &renderer.OpenGLRenderer{}, 1280, 720, "TextRendering", &TextRenderingScene{})
}
