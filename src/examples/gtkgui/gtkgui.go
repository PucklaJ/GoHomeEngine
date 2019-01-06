package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/GTK"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/renderers/OpenGL"
)

func main() {
	gohome.MainLop.Run(&framework.GTKFramework{UseWholeWindowAsGLArea: false, UseExternalWindow: true},
		&renderer.OpenGLRenderer{}, 640, 480, "GTKGUI", &GTKGUIScene{})
}
