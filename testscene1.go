package main

import "github.com/PucklaMotzer09/gohomeengine/src/gohome"
import "fmt"

// import "math"
import "math/rand"

type TestScene1 struct {
	boxEnts     [5]gohome.Entity3D
	boxEntTobjs [5]gohome.TransformableObject3D
	cam3d       gohome.Camera3D
	tst         gohome.TestCameraMovement3D
	light       gohome.PointLight
	light1      gohome.SpotLight
	planes      [6]gohome.Entity3D
	planeTobjs  [6]gohome.TransformableObject3D
}

func (this *TestScene1) Init() {
	fmt.Println("Hello")
	gohome.InitDefaultValues()

	for i := 0; i < 5; i++ {
		this.boxEnts[i].InitMesh(gohome.Box("Box", [3]float32{1.0, 1.0, 1.0}), &this.boxEntTobjs[i])
		this.boxEntTobjs[i].Rotation = [3]float32{rand.Float32() * 180.0, rand.Float32() * 180.0, rand.Float32() * 180.0}
		// if i == 0 {
		gohome.RenderMgr.AddObject(&this.boxEnts[i], &this.boxEntTobjs[i])
		// }
	}
	this.boxEntTobjs[0].Position = [3]float32{2.0, 0.0, 0.0}
	this.boxEntTobjs[1].Position = [3]float32{-2.0, 0.0, 0.0}
	this.boxEntTobjs[2].Position = [3]float32{0.0, 2.0, 0.0}
	this.boxEntTobjs[3].Position = [3]float32{0.0, -2.0, 0.0}
	this.boxEntTobjs[4].Position = [3]float32{0.0, 0.0, -2.0}

	gohome.LightMgr.SetAmbientLight(&gohome.Color{80, 80, 80, 255}, 0)
	gohome.LightMgr.CurrentLightCollection = 0

	this.tst.Init(&this.cam3d)
	gohome.UpdateMgr.AddObject(&this.tst)
	gohome.RenderMgr.SetCamera3D(&this.cam3d, 0)

	this.light = gohome.PointLight{
		Position: [3]float32{0.0, 0.0, 0.0},
		// Direction:     [3]float32{0.0, -1.0, 0.0},
		DiffuseColor:  &gohome.Color{20, 20, 20, 255},
		SpecularColor: &gohome.Color{20, 20, 20, 255},
		Attentuation: gohome.Attentuation{
			Constant: 1.0,
		},
		FarPlane:     25.0,
		CastsShadows: 1,
	}
	this.light1 = gohome.SpotLight{
		Position:      [3]float32{0.0, 1.0, 0.0},
		Direction:     [3]float32{0.0, -1.0, 0.0},
		DiffuseColor:  &gohome.Color{128, 128, 128, 128},
		SpecularColor: &gohome.Color{128, 128, 128, 128},
		Attentuation: gohome.Attentuation{
			Constant: 1.0,
		},
		CastsShadows: 1,
		OuterCutOff:  10.0,
		InnerCutOff:  5.0,
	}
	this.light.InitShadowmap(1024, 1024)
	this.light1.InitShadowmap(1024, 1024)

	gohome.LightMgr.AddPointLight(&this.light, 0)
	gohome.LightMgr.AddSpotLight(&this.light1, 0)

	const PLANE_SIZE float32 = 7.0

	for i := 0; i < 6; i++ {
		this.planes[i].InitMesh(gohome.Plane("Plane", [2]float32{PLANE_SIZE, PLANE_SIZE}, 1), &this.planeTobjs[i])
		gohome.RenderMgr.AddObject(&this.planes[i], &this.planeTobjs[i])
	}

	this.planeTobjs[0].Position = [3]float32{0.0, -PLANE_SIZE / 2.0, 0.0}
	this.planeTobjs[1].Position = [3]float32{0.0, PLANE_SIZE / 2.0, 0.0}
	this.planeTobjs[1].Rotation = [3]float32{180.0, 0.0, 0.0}
	this.planeTobjs[2].Position = [3]float32{0.0, 0.0, PLANE_SIZE / 2.0}
	this.planeTobjs[2].Rotation = [3]float32{-90.0, 0.0, 0.0}
	this.planeTobjs[3].Position = [3]float32{0.0, 0.0, -PLANE_SIZE / 2.0}
	this.planeTobjs[3].Rotation = [3]float32{90.0, 0.0, 0.0}
	this.planeTobjs[4].Position = [3]float32{-PLANE_SIZE / 2.0, 0.0, 0.0}
	this.planeTobjs[4].Rotation = [3]float32{0.0, 0.0, -90.0}
	this.planeTobjs[5].Position = [3]float32{PLANE_SIZE / 2.0, 0.0, 0.0}
	this.planeTobjs[5].Rotation = [3]float32{0.0, 0.0, 90.0}

	var spr gohome.Sprite2D
	var sprTobj gohome.TransformableObject2D
	spr.Init("", &sprTobj)
	spr.Texture = this.light1.ShadowMap
	sprTobj.Size = [2]float32{512.0, 512.0}
	sprTobj.Scale = [2]float32{1.0, 1.0}
	gohome.RenderMgr.AddObject(&spr, &sprTobj)

}

// var elapsed_time float32 = 0.0

func (this *TestScene1) Update(delta_time float32) {
	// elapsed_time += delta_time

	// this.light.Position[1] = float32(math.Sin(float64(elapsed_time)))
}

func (this *TestScene1) Terminate() {

}
