package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type ModelScene struct {
	gopher gohome.Entity3D
}

func (this *ModelScene) Init() {
	gohome.Init3DShaders()
	gohome.LightMgr.DisableLighting()
	gohome.ResourceMgr.LoadLevel("Gopher", "assets/models/gopher.obj", true)

	this.gopher.InitModel(gohome.ResourceMgr.GetModel("Gopher"))

	this.gopher.Transform.Position = [3]float32{0.0, 0, -5.0}
	this.gopher.Transform.Scale = [3]float32{0.75, 0.75, 0.75}

	gohome.RenderMgr.AddObject(&this.gopher)

	gohome.RenderMgr.EnableBackBuffer = false
}

func (this *ModelScene) Update(delta_time float32) {
	this.gopher.Transform.Rotation = this.gopher.Transform.Rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{0.0, 1.0, 0.0}))
}

func (this *ModelScene) Terminate() {
	this.gopher.Terminate()
}
