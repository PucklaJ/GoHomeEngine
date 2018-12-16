package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"golang.org/x/image/colornames"
	"math"
)

type ModelScene struct {
	gopher gohome.InstancedEntity3D
	light  gohome.SpotLight
	cube   gohome.Entity3D
	cam    gohome.Camera3D
	track  bool
	rotate bool
}

const SIZE = 2

func (this *ModelScene) Init() {
	this.light = gohome.SpotLight{
		Position:      [3]float32{0, 0, 0},
		Direction:     [3]float32{0, 0, -1},
		DiffuseColor:  colornames.White,
		SpecularColor: colornames.White,
		CastsShadows:  0,
		Attentuation: gohome.Attentuation{
			Constant:  2,
			Linear:    0,
			Quadratic: 0,
		},
		OuterCutOff: 20.0,
		InnerCutOff: 19.0,
		FarPlane:    1000.0,
		NearPlane:   0.1,
	}
	dir := gohome.DirectionalLight{
		Direction:      mgl32.Vec3{1.0, -1.0, -1.0}.Normalize(),
		DiffuseColor:   colornames.Khaki,
		SpecularColor:  colornames.Khaki,
		CastsShadows:   0,
		ShadowDistance: 15.0,
	}
	// this.light.InitShadowmap(1280, 720)
	// dir.InitShadowmap(3379, 3379)
	gohome.LightMgr.AddSpotLight(&this.light, 0)
	gohome.LightMgr.AddDirectionalLight(&dir, 0)
	gohome.LightMgr.SetAmbientLight(gohome.Color{10, 10, 10, 255}, 0)

	gohome.ResourceMgr.LoadLevel("Meat", "Meat.obj", false)
	planeCol := gohome.ResourceMgr.LoadTexture("PlaneColor", "Brick_Wall_015_COLOR.jpg")
	planeNorm := gohome.ResourceMgr.LoadTexture("PlaneNorm", "Brick_Wall_015_NORM.jpg")
	planeRough := gohome.ResourceMgr.LoadTexture("PlaneSpec", "Brick_Wall_015_ROUGH.jpg")

	this.gopher.InitModel(gohome.InstancedModel3DFromModel3D(gohome.ResourceMgr.GetModel("Meat")), SIZE*SIZE)

	this.gopher.Transforms[0].Position = [3]float32{0.0, -1.75, 0.0}
	this.gopher.Transforms[0].Scale = [3]float32{0.75, 0.75, 0.75}
	pos := this.gopher.Transforms[0].Position
	for x := 0; x < SIZE; x++ {
		for y := 0; y < SIZE; y++ {
			trans := this.gopher.Transforms[x+SIZE*y]
			trans.Position = pos.Add([3]float32{float32(x) * 4.0, 0.0, float32(y) * 4.0})
			trans.Scale = [3]float32{0.75, 0.75, 0.75}
		}
	}

	plane := gohome.Plane("Plane", [2]float32{100.0, 100.0}, 5.0, true)
	plane.GetMaterial().DiffuseTexture = planeCol
	plane.GetMaterial().NormalMap = planeNorm
	plane.GetMaterial().SpecularTexture = planeRough
	var planeEnt gohome.Entity3D
	planeEnt.InitMesh(plane)
	planeEnt.Transform.Position[1] = -2.0

	this.cube.InitMesh(gohome.Box("Box", [3]float32{0.5, 0.5, 0.5}, true))
	this.cube.SetShader(gohome.LoadGeneratedShader3D(gohome.SHADER_TYPE_3D, gohome.SHADER_FLAG_NOUV|gohome.SHADER_FLAG_NO_LIGHTING))
	this.cube.SetType(gohome.TYPE_3D_NORMAL)

	gohome.RenderMgr.AddObject(&this.gopher)
	gohome.RenderMgr.AddObject(&planeEnt)
	gohome.RenderMgr.AddObject(&this.cube)

	this.cam.Init()

	gohome.RenderMgr.SetCamera3D(&this.cam, 0)
	// var move gohome.TestCameraMovement3D
	// move.Init(&this.cam)
	// gohome.UpdateMgr.AddObject(&move)

	this.track = false
	this.rotate = true

	this.cam.LookDirection = mgl32.Vec3{1.0, -1.0, -1.0}.Normalize()
	this.cam.Position = mgl32.Vec3{-1.0, 1.0, 1.0}.Mul(10.0)
}

var val float32 = 0.0

const DIST float32 = 10.0

func (this *ModelScene) Update(delta_time float32) {
	// this.gopher.Transform.Rotation = this.gopher.Transform.Rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{0.0, 1.0, 0.0}))

	if this.rotate {
		pos := this.gopher.Transforms[SIZE*SIZE*0.25].Position
		x := float32(math.Sin(float64(val)) * float64(DIST))
		z := float32(math.Cos(float64(val)) * float64(DIST))

		this.light.Position = [3]float32{pos[0] + x, 1.0, pos[2] + z}
		this.light.Direction = pos.Sub(this.light.Position).Normalize()

		val += float32(math.Pi/2.0) * delta_time
	}

	if gohome.InputMgr.JustPressed(gohome.KeyT) {
		this.track = !this.track
	} else if gohome.InputMgr.JustPressed(gohome.MouseButtonLeft) {
		this.rotate = !this.rotate
	}

	if this.track {
		this.light.Position = this.cam.Position
		this.light.Direction = this.cam.LookDirection
	}

	this.cube.Transform.Position = this.light.Position
}

func (this *ModelScene) Terminate() {
	this.gopher.Terminate()
}
