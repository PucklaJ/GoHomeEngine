package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/SDL2"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/renderers/OpenGLES2"
	"golang.org/x/image/colornames"
)

type ExampleScene struct {
}

func (*ExampleScene) Init() {
	gohome.Render.SetBackgroundColor(colornames.Lime)
}

func (*ExampleScene) Update(delta_time float32) {

}

func (*ExampleScene) Terminate() {

}

func main() {
	gohome.MainLop.Run(&framework.SDL2Framework{}, &renderer.OpenGLES2Renderer{}, 640, 480, "Example", &ExampleScene{})
}
