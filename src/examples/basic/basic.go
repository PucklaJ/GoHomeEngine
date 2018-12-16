package main

import "C"

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/SDL2"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/renderers/OpenGLES31"
)

func main() {
	gohome.MainLop.Run(&framework.SDL2Framework{}, &renderer.OpenGLES31Renderer{}, 640, 480, "Basic", &BasicScene{})
}

//export SDL_main
func SDL_main() {
	main()
}
