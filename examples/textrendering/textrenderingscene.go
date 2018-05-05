package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type TextRenderingScene struct {
}

func (this *TextRenderingScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadFont("Font", "/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf")

	var text gohome.Text2D
	var textTobj gohome.TransformableObject2D

	text.Init("Font", 30, "Hello World! I can render text too!", &textTobj)

	textTobj.Origin = [2]float32{0.5, 0.5}
	textTobj.Position = gohome.Framew.WindowGetSize().Mul(0.5)

	gohome.RenderMgr.AddObject(&text, &textTobj)
}

func (this *TextRenderingScene) Update(delta_time float32) {

}

func (this *TextRenderingScene) Terminate() {

}
