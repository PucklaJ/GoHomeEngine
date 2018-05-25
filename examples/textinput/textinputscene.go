package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type TextInputScene struct {
	text gohome.Text2D
}

func (this *TextInputScene) Init() {
	gohome.ErrorMgr.ErrorLevel = gohome.ERROR_LEVEL_WARNING
	gohome.Init2DShaders()

	gohome.ResourceMgr.LoadFont("Font", "FreeMonoBold.ttf")

	this.text.Init("Font", 30, "Write:")
	gohome.RenderMgr.AddObject(&this.text)

	gohome.RenderMgr.EnableBackBuffer = false
	gohome.LightMgr.CurrentLightCollection = -1

	gohome.Framew.StartTextInput()
}

func (this *TextInputScene) Update(delta_time float32) {
	if input := gohome.Framew.GetTextInput(); len(input) != 0 {
		this.text.Text += input
	}

	if gohome.InputMgr.JustPressed(gohome.KeyEnter) {
		this.text.Text += "\n"
	} else if gohome.InputMgr.JustPressed(gohome.KeyBackspace) {
		if len(this.text.Text) > 0 {
			this.text.Text = this.text.Text[:len(this.text.Text)-1]
		}
	}
}

func (this *TextInputScene) Terminate() {
	this.text.Terminate()
}
