package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type BasicScene struct {
}

func (this *BasicScene) Init() {
	gohome.ResourceMgr.LoadTexture("Gopher", "gopher.png")
	gohome.ResourceMgr.LoadShader(gohome.SPRITE2D_SHADER_NAME, "vertex1.glsl", "fragment.glsl", "", "", "", "")

	var gopher gohome.Sprite2D
	var gopherTobj gohome.TransformableObject2D
	gopher.Init("Gopher", &gopherTobj)

	gopherTobj.Origin = [2]float32{0.5, 0.5}
	nw, nh := gohome.Render.GetNativeResolution()
	gopherTobj.Position = [2]float32{float32(nw) / 2.0, float32(nh) / 2.0}

	gohome.RenderMgr.AddObject(&gopher, &gopherTobj)
}

func (this *BasicScene) Update(delta_time float32) {

}

func (this *BasicScene) Terminate() {

}
