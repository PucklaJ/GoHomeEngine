package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

type BasicScene struct {
	gopher gohome.Sprite2D
}

func (this *BasicScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadTexture("Gopher", "gopher.png")

	this.gopher.Init("Gopher")

	this.gopher.Transform.Origin = [2]float32{0.5, 0.5}
	this.gopher.Transform.Position = gohome.Render.GetNativeResolution().Div(2.0)

	gohome.RenderMgr.AddObject(&this.gopher)
}

func (this *BasicScene) Update(delta_time float32) {

}

func (this *BasicScene) Terminate() {
	gohome.RenderMgr.RemoveObject(&this.gopher)
	this.gopher.Terminate()
}
