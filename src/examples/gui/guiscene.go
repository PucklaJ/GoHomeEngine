package main

import (
	"fmt"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

type GUIScene struct {
	btn  gohome.Button
	btn1 gohome.Button

	slider gohome.Slider
}

func (this *GUIScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadFont("Button", "/usr/share/fonts/truetype/ubuntu/Ubuntu-R.ttf")
	gohome.ButtonFont = "Button"
	this.btn.Text = "0"
	this.btn.Init(gohome.Framew.WindowGetSize().Mul(0.25), "")
	this.btn.Transform.Origin = [2]float32{0.5, 0.5}
	this.btn.PressCallback = func(btn *gohome.Button) {
		fmt.Println("Pressed 0")
	}
	this.btn1.Text = "1"
	this.btn1.Init(gohome.Framew.WindowGetSize().Mul(0.75), "")
	this.btn1.Transform.Origin = [2]float32{0.5, 0.5}
	this.btn1.PressCallback = func(btn *gohome.Button) {
		fmt.Println("Pressed 1")
	}

	this.slider.Init([2]float32{300.0, 100.0}, "", "")
}

func (this *GUIScene) Update(delta_time float32) {
}

func (this *GUIScene) Terminate() {
	this.btn.Terminate()
	this.btn.Sprite2D.Terminate()
	this.btn1.Terminate()
	this.btn1.Sprite2D.Terminate()
	this.slider.Terminate()
	this.slider.Long.Terminate()
	this.slider.Circle.Terminate()
}
