package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"golang.org/x/image/colornames"
)

type BasicScene struct {
	gopher gohome.Sprite2D
	rtex   gohome.RenderTexture
}

func (this *BasicScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadTexture("Gopher", "gopher.png")

	this.gopher.Init("Gopher")

	this.gopher.Transform.Origin = [2]float32{0.5, 0.5}
	this.gopher.Transform.Position = gohome.Render.GetNativeResolution().Div(2.0)

	gohome.RenderMgr.AddObject(&this.gopher)

	gohome.Render.SetBackgroundColor(colornames.Lime)

	gohome.Framew.OnResize(func(nw, nh uint32) {
		gohome.Framew.Log("Resize:", nw, nh)
		gohome.Render.SetNativeResolution(nw, nh)
	})
	gohome.RenderMgr.UpdateProjectionWithViewport = true

	this.rtex = gohome.Render.CreateRenderTexture("TestRTex", 200, 200, 1, false, false, false, false)
	this.rtex.SetAsTarget()
	gohome.Render.ClearScreen(colornames.Midnightblue)
	gohome.DrawCircle2D([2]float32{100, 100}, 50)
	this.rtex.UnsetAsTarget()

	var spr gohome.Sprite2D
	spr.InitTexture(this.rtex)
	spr.Flip = gohome.FLIP_VERTICAL
	gohome.RenderMgr.AddObject(&spr)
}

func (this *BasicScene) Update(delta_time float32) {
	this.gopher.Transform.Position = gohome.Render.GetNativeResolution().Div(2.0)
}

func (this *BasicScene) Terminate() {
	gohome.RenderMgr.RemoveObject(&this.gopher)
	this.gopher.Terminate()
}
