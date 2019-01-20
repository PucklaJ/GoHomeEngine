package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type Shape3DScene struct {
	shape gohome.Shape3D
}

func (this *Shape3DScene) Init() {
	this.shape.Init()
	this.shape.AddLines([]gohome.Line3D{
		{
			{
				-1.0, -1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
			{
				1.0, -1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
		},
		{
			{
				1.0, -1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
			{
				1.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
		},
		{
			{
				1.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
			{
				-1.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
		},
		{
			{
				-1.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
			{
				-1.0, -1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
		},
		{
			{
				-1.0, -1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
			{
				1.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
		},
		{
			{
				-1.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
			{
				1.0, -1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
		},
		{
			{
				0.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
			{
				0.0, -1.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
		},
		{
			{
				-1.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
			{
				1.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0,
			},
		},
	})
	this.shape.Transform.Position[2] = -3.0
	this.shape.Load()

	gohome.RenderMgr.AddObject(&this.shape)
}

func (this *Shape3DScene) Update(delta_time float32) {
	this.shape.Transform.Rotation = this.shape.Transform.Rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{0.0, 1.0, 0.0}))
}

func (this *Shape3DScene) Terminate() {

}
