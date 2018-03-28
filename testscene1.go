package main

import "github.com/PucklaMotzer09/gohomeengine/src/gohome"

type TestScene1 struct {
	cam3d gohome.Camera3D
	tst   gohome.TestCameraMovement3D
}

func (this *TestScene1) Init() {
	gohome.ResourceMgr.PreloadShader("2D", "vertex1.glsl", "fragment.glsl", "", "", "", "")
	gohome.ResourceMgr.PreloadShader("3D", "vertex3d.glsl", "fragment3d.glsl", "", "", "", "")
	gohome.ResourceMgr.PreloadShader("Normal", "vertex3d.glsl", "normalFrag.glsl", "", "", "", "")
	gohome.ResourceMgr.PreloadLevel("Kratos", "Kratos.obj")
	gohome.ResourceMgr.LoadPreloadedResources()

	nW, nH := gohome.Render.GetNativeResolution()
	lvl := gohome.ResourceMgr.GetLevel("Kratos")

	gohome.RenderMgr.SetProjection2D(&gohome.Ortho2DProjection{
		Left:   0.0,
		Right:  float32(nW),
		Top:    0.0,
		Bottom: float32(nH),
	})
	gohome.RenderMgr.SetProjection3D(&gohome.PerspectiveProjection{
		Width:     float32(nW),
		Height:    float32(nH),
		FOV:       70.0,
		NearPlane: 0.1,
		FarPlane:  1000.0,
	})

	this.tst.Init(&this.cam3d)
	this.cam3d.Position = [3]float32{0.0, 0.5, 2.0}

	gohome.RenderMgr.SetCamera3D(&this.cam3d, 0)
	gohome.LightMgr.CurrentLightCollection = -1
	gohome.UpdateMgr.AddObject(&this.tst)

	lvl.AddToScene()
}

func (this *TestScene1) Update(delta_time float32) {
	if gohome.InputMgr.JustPressed(gohome.KeyB) {
		if gohome.RenderMgr.ForceShader3D == nil {
			gohome.RenderMgr.ForceShader3D = gohome.ResourceMgr.GetShader("Normal")
		} else {
			gohome.RenderMgr.ForceShader3D = nil
		}
	}
}

func (this *TestScene1) Terminate() {

}
