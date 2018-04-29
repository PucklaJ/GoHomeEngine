package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type TestScene3 struct {
	spr     gohome.Sprite2D
	sprTobj gohome.TransformableObject2D
}

func (this *TestScene3) Init() {
	gohome.InitDefaultValues()
	gohome.FPSLimit.MaxFPS = 1000

	gohome.LightMgr.CurrentLightCollection = -1
	gohome.RenderMgr.EnableBackBuffer = true
	gohome.ResourceMgr.LoadTexture("Image", "textures/image.tga")

	this.spr.Init("Image", &this.sprTobj)
	nw, nh := gohome.Render.GetNativeResolution()
	this.sprTobj.Size = [2]float32{float32(nw), float32(nh)}
	gohome.RenderMgr.AddObject(&this.spr, &this.sprTobj)
}

func (this *TestScene3) Update(delta_time float32) {
}

func (this *TestScene3) Terminate() {
	this.spr.Terminate()
}
