package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type TextRenderingScene struct {
}

func (this *TextRenderingScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadFont("AbyssinicaSIL-R", "/usr/share/fonts/truetype/abyssinica/AbyssinicaSIL-R.ttf")

	var text gohome.Text2D
	var textTobj gohome.TransformableObject2D
	text.Init("AbyssinicaSIL-R", 24, "Hello World", &textTobj)

	gohome.RenderMgr.AddObject(&text, &textTobj)

}

func (this *TextRenderingScene) Update(delta_time float32) {

}

func (this *TextRenderingScene) Terminate() {

}
