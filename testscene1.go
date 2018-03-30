package main

import "github.com/PucklaMotzer09/gohomeengine/src/gohome"

type TestScene1 struct {
	boxEnt     gohome.Entity3D
	boxEntTobj gohome.TransformableObject3D
}

func (this *TestScene1) Init() {
	gohome.InitDefaultValues()

	this.boxEnt.InitMesh(gohome.Box("Box", [3]float32{1.0, 1.0, 1.0}), &this.boxEntTobj)
	this.boxEntTobj.Position = [3]float32{0.0, 0.0, -3.0}
	gohome.LightMgr.CurrentLightCollection = -1

	gohome.RenderMgr.AddObject(&this.boxEnt, &this.boxEntTobj)
}

func (this *TestScene1) Update(delta_time float32) {
}

func (this *TestScene1) Terminate() {

}
