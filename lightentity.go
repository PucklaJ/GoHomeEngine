package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type LightEntity struct {
	entity        gohome.Entity3D
	tobj          gohome.TransformableObject3D
	rotationPoint mgl32.Vec3
	light         gohome.PointLight
	elapsed_time  float32
	radius        float32
}

func (this *LightEntity) Init(rotate mgl32.Vec3, offset, radius float32) {
	this.rotationPoint = rotate
	this.radius = radius
	this.entity.InitMesh(gohome.Box("Box", [3]float32{1.0, 1.0, 1.0}), &this.tobj)
	this.entity.Model3D.GetMesh("Box").GetMaterial().SetColors(gohome.Color{255, 255, 255, 255}, gohome.Color{255, 255, 255, 255})
	this.tobj.Scale = [3]float32{0.2, 0.2, 0.2}
	this.light = gohome.PointLight{
		DiffuseColor:  gohome.Color{255, 255, 255, 255},
		SpecularColor: gohome.Color{255, 255, 255, 255},
		Attentuation: gohome.Attentuation{
			Linear: 3.0,
		},
		CastsShadows: 1,
		FarPlane:     25.0,
	}
	this.light.InitShadowmap(1024, 1024)
	// gohome.RenderMgr.AddObject(&this.entity, &this.tobj)
	gohome.UpdateMgr.AddObject(this)
	gohome.LightMgr.AddPointLight(&this.light, 0)
	this.elapsed_time = offset
}

func (this *LightEntity) Update(delta_time float32) {
	this.elapsed_time += delta_time
	this.tobj.Position[0] = float32(math.Cos(float64(this.elapsed_time)))*this.radius + this.rotationPoint[0]
	this.tobj.Position[2] = float32(math.Sin(float64(this.elapsed_time)))*this.radius + this.rotationPoint[2]
	this.tobj.Position[1] = this.rotationPoint[1]
	this.light.Position = this.tobj.Position.Add(mgl32.Vec3{0.0, -2.0 * this.tobj.Scale[1], 0.0})
}

func (this *LightEntity) Terminate() {
	this.entity.Terminate()
}
