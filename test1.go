package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

func main() {
	gohome.MainLoop{}.Run(&gohome.GLFWFramework{}, &gohome.OpenGLRenderer{}, 1280, 720, "GoHomeEngine", &TestScene{})
}
