package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/Android"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"


)

func main() {
	gohome.MainLop.Run(&framework.AndroidFramework{}, &renderer.OpenGLESRenderer{}, 640, 480, "Basic", &BasicScene{})
}
