package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type BasicScene struct {
	gopher       gohome.Sprite2D
	littleGopher gohome.Sprite2D
	cam          gohome.Camera2D
}

func (this *BasicScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadTexture("Gopher", "gopher.png")

	this.gopher.Init("Gopher")
	this.littleGopher.Init("Gopher")

	this.gopher.Transform.Origin = [2]float32{0.5, 0.5}
	this.littleGopher.Transform.Origin = [2]float32{0.5, 0.5}
	nr := gohome.Render.GetNativeResolution()
	this.gopher.Transform.Position = nr.Div(2.0)
	this.littleGopher.Transform.Scale = [2]float32{0.3, 0.3}
	this.littleGopher.Depth = 1
	this.gopher.Depth = 0

	gohome.RenderMgr.AddObject(&this.littleGopher)
	gohome.RenderMgr.AddObject(&this.gopher)

	this.cam.Zoom = 1.0
	this.cam.Origin = mgl32.Vec2{0.5, 0.5}
	gohome.RenderMgr.SetCamera2D(&this.cam, 0)

	gohome.RenderMgr.EnableBackBuffer = false
	gohome.ErrorMgr.ErrorLevel = gohome.ERROR_LEVEL_WARNING
}

func (this *BasicScene) Update(delta_time float32) {
	if gohome.InputMgr.IsPressed(gohome.KeyW) {
		this.cam.AddPositionRotated(mgl32.Vec2{0.0, -100.0}.Mul(delta_time))
	} else if gohome.InputMgr.IsPressed(gohome.KeyS) {
		this.cam.AddPositionRotated(mgl32.Vec2{0.0, 100.0}.Mul(delta_time))
	} else if gohome.InputMgr.IsPressed(gohome.KeyA) {
		this.cam.AddPositionRotated(mgl32.Vec2{-100.0, 0.0}.Mul(delta_time))
	} else if gohome.InputMgr.IsPressed(gohome.KeyD) {
		this.cam.AddPositionRotated(mgl32.Vec2{100.0, 0.0}.Mul(delta_time))
	}

	this.cam.Zoom += float32(gohome.InputMgr.Mouse.Wheel[1]) * 0.1

	if gohome.InputMgr.IsPressed(gohome.KeyUp) {
		this.cam.Rotation += 30.0 * delta_time
	} else if gohome.InputMgr.IsPressed(gohome.KeyDown) {
		this.cam.Rotation -= 30.0 * delta_time
	}

	worldPos := gohome.InputMgr.Mouse.ToWorldPosition2D()
	this.littleGopher.Transform.Position = worldPos
}

func (this *BasicScene) Terminate() {
	gohome.RenderMgr.RemoveObject(&this.gopher)
	gohome.RenderMgr.RemoveObject(&this.littleGopher)
	this.gopher.Terminate()
	this.littleGopher.Terminate()
}
