package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGL"
)

func main() {
	gohome.MainLop.Run(&framework.GTKFramework{UseWholeWindowAsGLArea: false, UseExternalWindow: true}, &renderer.OpenGLRenderer{}, 1280, 720, "GTKBuilder", &GTKBuilderScene{})
}
