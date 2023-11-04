package main

import (
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
)

type BasicScene struct {
	gopher gohome.Sprite2D
}

func (this *BasicScene) Init() {
	gohome.ResourceMgr.LoadTexture("Gopher", "gopher.png")

	this.gopher.Init("Gopher")

	this.gopher.Transform.Position = gohome.Render.GetNativeResolution().Div(2)
	this.gopher.Transform.Origin = [2]float32{0.5, 0.5}

	gohome.RenderMgr.AddObject(&this.gopher)
}

func (this *BasicScene) Update(delta_time float32) {
}

func (this *BasicScene) Terminate() {
}
