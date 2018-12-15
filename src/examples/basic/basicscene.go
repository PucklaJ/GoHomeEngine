package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"golang.org/x/image/colornames"
)

type BasicScene struct {
	gopher    gohome.Sprite2D
	wallpaper gohome.Sprite2D
	rtex      gohome.RenderTexture
}

func (this *BasicScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadTexture("Gopher", "gopher.png")
	gohome.ResourceMgr.LoadTexture("Wallpaper", "wallpaper_10.png")

	this.gopher.Init("Gopher")
	this.wallpaper.Init("Wallpaper")

	nr := gohome.Render.GetNativeResolution()

	this.gopher.Transform.Origin = [2]float32{0.5, 0.5}
	this.gopher.Transform.Position = nr.Div(2.0)

	this.wallpaper.Transform.Origin = [2]float32{0.5, 0.5}
	this.wallpaper.Transform.Position = nr.Div(2.0)
	this.wallpaper.Transform.Size = [2]float32{nr[1], nr[0]}
	this.wallpaper.Transform.Rotation = 90.0

	gohome.RenderMgr.AddObject(&this.wallpaper)
	gohome.RenderMgr.AddObject(&this.gopher)

	gohome.Render.SetBackgroundColor(colornames.Lime)

	gohome.Framew.OnResize(func(nw, nh uint32) {
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
	gohome.RenderMgr.EnableBackBuffer = true
	nr := gohome.Render.GetNativeResolution()
	this.gopher.Transform.Position = nr.Div(2.0)
	this.wallpaper.Transform.Position = nr.Div(2.0)
	this.wallpaper.Transform.Size = [2]float32{nr[1], nr[0]}
}

func (this *BasicScene) Terminate() {
	gohome.RenderMgr.RemoveObject(&this.gopher)
	this.gopher.Terminate()
}
