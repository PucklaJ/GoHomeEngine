package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type CubeScene struct {
	cube gohome.Entity3D
}

func (this *CubeScene) Init() {
	gohome.Init3DShaders()
	gohome.ResourceMgr.LoadTexture("CubeImage", "cube.png")

	mesh := gohome.Box("Cube", [3]float32{1.0, 1.0, 1.0})
	mesh.GetMaterial().SetTextures("CubeImage", "", "")
	this.cube.InitMesh(mesh)
	this.cube.Transform.Position = [3]float32{0.0, 0.0, -3.0}

	gohome.RenderMgr.AddObject(&this.cube)
	gohome.LightMgr.DisableLighting()
}

func (this *CubeScene) Update(delta_time float32) {
	this.cube.Transform.Rotation[0] += 30.0 * delta_time
	this.cube.Transform.Rotation[1] += 30.0 * delta_time
}

func (this *CubeScene) Terminate() {

}
