package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type BasicScene struct {
}

func (this *BasicScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadTexture("Gopher", "gopher.png")

	var gopher gohome.Sprite2D
	gopher.Init("Gopher")

	gopher.Transform.Origin = [2]float32{0.5, 0.5}
	nw, nh := gohome.Render.GetNativeResolution()
	gopher.Transform.Position = [2]float32{float32(nw) / 2.0, float32(nh) / 2.0}

	gohome.RenderMgr.AddObject(&gopher)

	gohome.RenderMgr.EnableBackBuffer = false
}

func (this *BasicScene) Update(delta_time float32) {

}

func (this *BasicScene) Terminate() {

}
