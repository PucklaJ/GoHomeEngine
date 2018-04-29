package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type TestScene2 struct {
	box     gohome.Entity3D
	boxTobj gohome.TransformableObject3D
}

func (this *TestScene2) Init() {
	gohome.InitDefaultValues()
	gohome.FPSLimit.MaxFPS = 1000

	gohome.ResourceMgr.LoadTexture("Image", "image.tga")
	gohome.ResourceMgr.LoadShader("3D Simple", "vertex3dNoShadows.glsl", "fragment3dSimple.glsl", "", "", "", "")

	this.box.InitMesh(gohome.Box("Box", [3]float32{1.0, 1.0, 1.0}), &this.boxTobj)
	this.boxTobj.Position = [3]float32{0.0, 0.0, -3.0}
	this.box.Model3D.GetMeshIndex(0).GetMaterial().DiffuseTexture = gohome.ResourceMgr.GetTexture("Image")
	gohome.RenderMgr.AddObject(&this.box, &this.boxTobj)

	gohome.LightMgr.CurrentLightCollection = -1
	gohome.RenderMgr.EnableBackBuffer = false

	gohome.ResourceMgr.SetShader("3D", "3D Simple")
}

func (this *TestScene2) Update(delta_time float32) {
	this.boxTobj.Rotation[1] += 30.0 * delta_time
	this.boxTobj.Rotation[0] += 30.0 * delta_time
}

func (this *TestScene2) Terminate() {

}
