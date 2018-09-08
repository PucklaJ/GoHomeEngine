package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type PhysicsScene struct {
	PhysicsMgr physics2d.PhysicsManager2D
}

func (this *PhysicsScene) Init() {
	this.PhysicsMgr.Init(mgl32.Vec2{0.0, 100.0})
	debugDraw := this.PhysicsMgr.GetDebugDraw()
	gohome.RenderMgr.AddObject(&debugDraw)
	gohome.UpdateMgr.AddObject(&this.PhysicsMgr)

	this.AddBox(gohome.Framew.WindowGetSize().Mul(0.5))

	this.PhysicsMgr.CreateStaticBox(mgl32.Vec2{640.0 / 2.0, 480}, mgl32.Vec2{720, 20})

	gohome.RenderMgr.EnableBackBuffer = false
}

func (this *PhysicsScene) AddBox(pos mgl32.Vec2) {
	var size mgl32.Vec2
	size[0] = 20.0
	size[1] = 20.0

	body := this.PhysicsMgr.CreateDynamicBox(pos, size)
	body.SetAngularVelocity(physics2d.ToBox2DAngle(90))
}

func (this *PhysicsScene) AddCircle(pos mgl32.Vec2) {
	this.PhysicsMgr.CreateStaticCircle(pos, 10.0)
}

func (this *PhysicsScene) AddCar(pos mgl32.Vec2) {
	body1 := this.PhysicsMgr.CreateDynamicCircle(pos.Sub(mgl32.Vec2{20.0, 0.0}), 10.0)
	body2 := this.PhysicsMgr.CreateDynamicCircle(pos.Add(mgl32.Vec2{20.0, 0.0}), 10.0)
	jointDef := box2d.MakeB2RopeJointDef()
	jointDef.SetBodyA(body1)
	jointDef.SetBodyB(body2)
	jointDef.SetCollideConnected(true)
	jointDef.MaxLength = physics2d.ScalarToBox2D(10.0)
	this.PhysicsMgr.World.CreateJoint(&jointDef)
}

func (this *PhysicsScene) Update(delta_time float32) {
	if gohome.InputMgr.IsPressed(gohome.MouseButtonLeft) {
		pos := gohome.InputMgr.Mouse.ToWorldPosition2D()
		this.AddBox(pos)
	} else if gohome.InputMgr.JustPressed(gohome.MouseButtonRight) {
		pos := gohome.InputMgr.Mouse.ToWorldPosition2D()
		this.AddCircle(pos)
	} else if gohome.InputMgr.JustPressed(gohome.MouseButtonMiddle) {
		pos := gohome.InputMgr.Mouse.ToWorldPosition2D()
		this.AddCar(pos)
	}
}

func (this *PhysicsScene) Terminate() {
}
