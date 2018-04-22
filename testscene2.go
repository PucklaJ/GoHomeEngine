package main

import (
	// "encoding/binary"
	// "fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	// "github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"
	// "golang.org/x/mobile/exp/f32"
	// "golang.org/x/mobile/gl"
)

type TestScene2 struct {
	box     gohome.Entity3D
	boxTobj gohome.TransformableObject3D
}

func (this *TestScene2) Init() {
	gohome.InitDefaultValues()
	gohome.FPSLimit.MaxFPS = 1000

	gohome.ResourceMgr.LoadTexture("Image", "image.tga")
	// width, height := gohome.Render.GetNativeResolution()

	// var spr gohome.Sprite2D
	// var sprTobj gohome.TransformableObject2D
	// spr.Init("Image", &sprTobj)
	// gohome.RenderMgr.AddObject(&spr, &sprTobj)

	// sprTobj.Size[0] = float32(width)
	// sprTobj.Size[1] = float32(height)

	this.box.InitMesh(gohome.Box("Box", [3]float32{1.0, 1.0, 1.0}), &this.boxTobj)
	this.boxTobj.Position = [3]float32{0.0, 0.0, -3.0}
	this.box.Model3D.GetMeshIndex(0).GetMaterial().DiffuseTexture = gohome.ResourceMgr.GetTexture("Image")
	gohome.RenderMgr.AddObject(&this.box, &this.boxTobj)

	gohome.LightMgr.CurrentLightCollection = -1
}

func (this *TestScene2) Update(delta_time float32) {
	this.boxTobj.Rotation[1] += 30.0 * delta_time
	this.boxTobj.Rotation[0] += 30.0 * delta_time
}

func (this *TestScene2) Terminate() {

}
