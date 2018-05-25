package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"strconv"
)

type ModelScene struct {
	gopher gohome.Entity3D
}

func (this *ModelScene) Init() {
	gohome.ErrorMgr.DuplicateMessages = false
	gohome.ErrorMgr.ErrorLevel = gohome.ERROR_LEVEL_WARNING

	gohome.Init3DShaders()
	gohome.ResourceMgr.PreloadLevel("Gopher", "gopher.obj", true)
	gohome.ResourceMgr.LoadPreloadedResources()

	this.gopher.InitModel(gohome.ResourceMgr.GetModel("Gopher"))

	this.gopher.Transform.Position = [3]float32{0.0, -1.75, -5.0}
	this.gopher.Transform.Scale = [3]float32{0.75, 0.75, 0.75}

	gohome.RenderMgr.AddObject(&this.gopher)

	gohome.LightMgr.CurrentLightCollection = -1
	gohome.RenderMgr.EnableBackBuffer = false
}

func (this *ModelScene) Update(delta_time float32) {
	this.gopher.Transform.Rotation[1] += 30.0 * delta_time
}

func (this *ModelScene) Terminate() {
	this.gopher.Terminate()
}
