package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type CubeScene struct {
	cubeTobj gohome.TransformableObject3D
}

func (this *CubeScene) Init() {
	gohome.Init3DShaders()
	gohome.ResourceMgr.LoadTexture("CubeImage", "cube.png")

	var cube gohome.Entity3D
	mesh := gohome.Box("Cube", [3]float32{1.0, 1.0, 1.0})
	mesh.GetMaterial().SetTextures("CubeImage", "", "")
	cube.InitMesh(mesh, &this.cubeTobj)
	this.cubeTobj.Position = [3]float32{0.0, 0.0, -3.0}

	gohome.RenderMgr.AddObject(&cube, &this.cubeTobj)
	gohome.LightMgr.CurrentLightCollection = -1
}

func (this *CubeScene) Update(delta_time float32) {
	this.cubeTobj.Rotation[0] += 30.0 * delta_time
	this.cubeTobj.Rotation[1] += 30.0 * delta_time
}

func (this *CubeScene) Terminate() {

}
