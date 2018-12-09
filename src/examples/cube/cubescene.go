package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type CubeScene struct {
	cube gohome.Entity3D
}

func (this *CubeScene) Init() {
	gohome.LightMgr.DisableLighting()
	gohome.ResourceMgr.LoadTexture("CubeImage", "cube.png")

	mesh := gohome.Box("Cube", [3]float32{1.0, 1.0, 1.0}, true)
	mesh.GetMaterial().SetTextures("CubeImage", "", "")
	this.cube.InitMesh(mesh)
	this.cube.Transform.Position = [3]float32{0.0, 0.0, -3.0}

	gohome.RenderMgr.AddObject(&this.cube)
}

func (this *CubeScene) Update(delta_time float32) {
	this.cube.Transform.Rotation = this.cube.Transform.Rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{0.0, 1.0, 0.0})).Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{1.0, 0.0, 0.0}))
}

func (this *CubeScene) Terminate() {

}
