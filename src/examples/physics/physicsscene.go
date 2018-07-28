package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/go-gl/mathgl/mgl32"
)

type PhysicsScene struct {
	PhysicsMgr physics2d.PhysicsManager2D
	boxes      []gohome.Sprite2D
}

func (this *PhysicsScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadTexture("Box", "rect1.png")

	this.PhysicsMgr.Init(mgl32.Vec2{0.0, 100.0})
	gohome.UpdateMgr.AddObject(&this.PhysicsMgr)

	this.AddBox(gohome.Framew.WindowGetSize().Mul(0.5))

	this.PhysicsMgr.CreateStaticBox(mgl32.Vec2{640.0 / 2.0, 480}, mgl32.Vec2{720, 20})
}

func (this *PhysicsScene) AddBox(pos mgl32.Vec2) {
	var box gohome.Sprite2D
	box.Init("Box")
	box.Transform.Scale = [2]float32{0.1, 0.1}
	gohome.RenderMgr.AddObject(&box)
	var size mgl32.Vec2
	size[0] = box.Transform.Size[0] * box.Transform.Scale[0]
	size[1] = box.Transform.Size[1] * box.Transform.Scale[1]

	body := this.PhysicsMgr.CreateDynamicBox(pos, size)
	body.SetAngularVelocity(physics2d.ToBox2DAngle(90))
	var con physics2d.PhysicsConnector2D
	con.Init(box.Transform, body)
	gohome.UpdateMgr.AddObject(&con)

	this.boxes = append(this.boxes, box)

	gohome.RenderMgr.EnableBackBuffer = true
}

func (this *PhysicsScene) Update(delta_time float32) {
	if gohome.InputMgr.IsPressed(gohome.MouseButtonLeft) {
		pos := gohome.InputMgr.Mouse.ToWorldPosition2D()
		this.AddBox(pos)
	}
}

func (this *PhysicsScene) Terminate() {
	for i := 0; i < len(this.boxes); i++ {
		this.boxes[i].Terminate()
	}
}
