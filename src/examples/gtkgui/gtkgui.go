package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGL"
)

func main() {

	gohome.MainLop.Run(&framework.GTKFramework{UseWholeWindowAsGLArea: false},
		&renderer.OpenGLRenderer{}, 640, 480, "GTKGUI", &CubeScene{})
}
