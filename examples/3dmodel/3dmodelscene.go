package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type ModelScene struct {
	gopherTobj gohome.TransformableObject3D
	gopher     gohome.Entity3D
}

func (this *ModelScene) Init() {
	gohome.Init3DShaders()
	gohome.ResourceMgr.LoadLevel("Gopher", "gopher.obj", true)

	this.gopher.InitModel(gohome.ResourceMgr.GetModel("Gopher"), &this.gopherTobj)

	this.gopherTobj.Position = [3]float32{0.0, -1.75, -5.0}
	this.gopherTobj.Scale = [3]float32{0.75, 0.75, 0.75}

	gohome.RenderMgr.AddObject(&this.gopher, &this.gopherTobj)

	gohome.LightMgr.CurrentLightCollection = -1
	gohome.RenderMgr.EnableBackBuffer = false
}

func (this *ModelScene) Update(delta_time float32) {
	this.gopherTobj.Rotation[1] += 30.0 * delta_time
}

func (this *ModelScene) Terminate() {
	this.gopher.Terminate()
}
