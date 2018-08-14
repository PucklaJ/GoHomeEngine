package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
)

type TiledMapScene struct {
	tiledMap   gohome.TiledMap
	PhysicsMgr physics2d.PhysicsManager2D
}

func (this *TiledMapScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.LoadTMXMap("outdoor", "orthogonal-outside.tmx")

	this.tiledMap.Init("outdoor")
	gohome.RenderMgr.AddObject(&this.tiledMap)

	gohome.RenderMgr.UpdateProjectionWithViewport = true
	gohome.RenderMgr.EnableBackBuffer = false

	w := this.tiledMap.Texture.GetWidth()
	h := this.tiledMap.Texture.GetHeight()

	gohome.Render.SetNativeResolution(uint32(w), uint32(h))
	gohome.Framew.WindowSetSize([2]float32{float32(w), float32(h)})

	this.PhysicsMgr.Init([2]float32{0.0, 30.0})
	debug := this.PhysicsMgr.GetDebugDraw()
	gohome.RenderMgr.AddObject(&debug)

	this.PhysicsMgr.LayerToCollision(&this.tiledMap, "Objects")
	gohome.UpdateMgr.AddObject(&this.PhysicsMgr)
}

func (this *TiledMapScene) Update(delta_time float32) {
	if gohome.InputMgr.JustPressed(gohome.MouseButtonLeft) {
		pos := gohome.InputMgr.Mouse.ToWorldPosition2D()
		this.PhysicsMgr.CreateDynamicBox(pos, [2]float32{20.0, 20.0})
	}
}

func (this *TiledMapScene) Terminate() {
}
