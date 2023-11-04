package main

import (
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
)

type TextRenderingScene struct {
}

func (this *TextRenderingScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadFont("Font", "FreeMonoBold.ttf")

	var text gohome.Text2D

	text.Init("Font", 30, "Lorem ipsum dolor sit amet,\nconsetetur sadipscing elitr,\nsed diam nonumy eirmod tempor invidunt ut\nlabore et dolore magna aliquyam erat,\nsed diam voluptua.\nAt vero eos et accusam et")

	text.Transform.Origin = [2]float32{0.5, 0.5}
	text.Transform.Position = gohome.Framew.WindowGetSize().Mul(0.5)

	gohome.RenderMgr.AddObject(&text)

	gohome.RenderMgr.EnableBackBuffer = false
	gohome.ErrorMgr.ErrorLevel = gohome.ERROR_LEVEL_WARNING
}

func (this *TextRenderingScene) Update(delta_time float32) {

}

func (this *TextRenderingScene) Terminate() {

}
