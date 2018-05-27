package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/Android"
)

func main() {
	gohome.MainLop.Run(&framework.AndroidFramework{},&renderer.OpenGLESRenderer{},640,480,"MultipleTouches",&MultipleTouchesScene{})
}
