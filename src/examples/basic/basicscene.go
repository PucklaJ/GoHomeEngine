package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type BasicScene struct {
	gopher gohome.Sprite2D
}

func (this *BasicScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadTexture("Gopher", "gopher.png")

	this.gopher.Init("Gopher")

	this.gopher.Transform.Origin = [2]float32{0.5, 0.5}
	nw, nh := gohome.Render.GetNativeResolution()
	this.gopher.Transform.Position = [2]float32{float32(nw) / 2.0, float32(nh) / 2.0}

	gohome.RenderMgr.AddObject(&this.gopher)

	gohome.RenderMgr.EnableBackBuffer = true
}

func (this *BasicScene) Update(delta_time float32) {

}

func (this *BasicScene) Terminate() {
	gohome.RenderMgr.RemoveObject(&this.gopher)
	this.gopher.Terminate()
}
