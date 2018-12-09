package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type Lines3DScene struct {
	lines gohome.Lines3D
}

func (this *Lines3DScene) Init() {
	this.lines.Init()
	this.lines.AddLines([]gohome.Line3D{
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
	this.lines.Transform.Position[2] = -3.0
	this.lines.Load()

	gohome.RenderMgr.AddObject(&this.lines)
}

func (this *Lines3DScene) Update(delta_time float32) {
	this.lines.Transform.Rotation = this.lines.Transform.Rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{0.0, 1.0, 0.0}))
}

func (this *Lines3DScene) Terminate() {

}
