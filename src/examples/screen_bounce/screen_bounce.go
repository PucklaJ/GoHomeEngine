package main

import (
	framework "github.com/PucklaJ/GoHomeEngine/src/frameworks/SAMURE"
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	renderer "github.com/PucklaJ/GoHomeEngine/src/renderers/OpenGL"
	"golang.org/x/image/colornames"
)

type ScreenBounceScene struct {
}

func (s *ScreenBounceScene) Init() {
	gohome.Render.SetBackgroundColor(colornames.Lime)
}

func (s *ScreenBounceScene) Update(delta_time float32) {
}

func (s *ScreenBounceScene) Terminate() {
}

func main() {
	gohome.MainLop.Run(&framework.SAMUREFramework{}, &renderer.OpenGLRenderer{}, 1920, 1080, "ScreenBounce", &ScreenBounceScene{})
}
