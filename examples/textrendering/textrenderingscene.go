package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type TextRenderingScene struct {
}

func (this *TextRenderingScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadFont("Font", "FreeMonoBold.ttf")

	var text gohome.Text2D
	var textTobj gohome.TransformableObject2D

	text.Init("Font", 30, "Lorem ipsum dolor sit amet,\nconsetetur sadipscing elitr,\nsed diam nonumy eirmod tempor invidunt ut\nlabore et dolore magna aliquyam erat,\nsed diam voluptua.\nAt vero eos et accusam et", &textTobj)

	textTobj.Origin = [2]float32{0.5, 0.5}
	textTobj.Position = gohome.Framew.WindowGetSize().Mul(0.5)

	gohome.RenderMgr.AddObject(&text, &textTobj)

	gohome.RenderMgr.EnableBackBuffer = false
}

func (this *TextRenderingScene) Update(delta_time float32) {

}

func (this *TextRenderingScene) Terminate() {

}
