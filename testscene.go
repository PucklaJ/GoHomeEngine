package main

import (
	"fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	// "github.com/go-gl/gl/v4.1-core/gl"
	"golang.org/x/image/colornames"
	// "github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
)

const (
	NUM_PLANES uint32  = 4
	PLANE_SIZE float32 = 10.0
)

type TestScene struct {
	lightEnts    [5]LightEntity
	cam3d        gohome.Camera3D
	cam3d1       gohome.Camera3D
	tst          gohome.TestCameraMovement3D
	planeEnt     gohome.Entity3D
	planeEntTobj gohome.TransformableObject3D
	direct       gohome.DirectionalLight
	direct1      gohome.DirectionalLight
	spot         gohome.SpotLight
	m4           gohome.Entity3D
	m4Tobj       gohome.TransformableObject3D
	m41          gohome.Entity3D
	m41Tobj      gohome.TransformableObject3D
	kratos       gohome.Entity3D
	kratosTobj   gohome.TransformableObject3D
	oakTrees     [NUM_PLANES * NUM_PLANES]gohome.Entity3D
	oakTreeTobjs [NUM_PLANES * NUM_PLANES]gohome.TransformableObject3D
}

func (this *TestScene) Init() {
	gohome.Render.SetNativeResolution(1920, 1080)
	start := time.Now()
	gohome.ResourceMgr.PreloadLevel("M4", "Arma_M4.obj")
	gohome.ResourceMgr.PreloadLevel("Kratos", "Kratos.obj")
	gohome.ResourceMgr.PreloadLevel("Barrel", "barrel.obj")
	gohome.ResourceMgr.PreloadLevel("Crate", "crate.obj")
	gohome.ResourceMgr.PreloadLevel("Pine", "pine.obj")
	gohome.ResourceMgr.PreloadShader("Normal", "vertex3d.glsl", "normalFrag.glsl", "", "", "", "")
	gohome.ResourceMgr.PreloadShader("ShadowMapRender", "vertex1.glsl", "shadowMapRenderFrag.glsl", "", "", "", "")
	gohome.ResourceMgr.PreloadTexture("Kratos_torso_n.tga", "Kratos_torso_n.tga")
	gohome.ResourceMgr.PreloadTexture("Kratos_legs_n.tga", "Kratos_legs_n.tga")
	gohome.ResourceMgr.PreloadTexture("Kratos_head_n.tga", "Kratos_head_n.tga")
	gohome.ResourceMgr.PreloadTexture("PlaneTexture", "159.JPG")
	gohome.ResourceMgr.PreloadTexture("PlaneNormalMap", "159_norm.JPG")
	gohome.ResourceMgr.PreloadTexture("Pine", "pine.png")
	gohome.ResourceMgr.LoadPreloadedResources()
	end := time.Now()
	fmt.Println("Needed Time:", end.Sub(start).Seconds()*1000.0, "ms")
	gohome.InitDefaultValues()

	gohome.LightMgr.SetAmbientLight(gohome.Color{20, 20, 20, 255}, 0)
	// this.lightEnts[0].Init([3]float32{0.0, 2.0, 0.0}, 0.0, 2.0)
	// this.lightEnts[1].Init([3]float32{0.0, 5.0, 0.0}, 1.0, 2.0)
	// this.lightEnts[2].Init([3]float32{3.0, 3.0, 0.0}, 2.0, 2.0)
	// this.lightEnts[3].Init([3]float32{-3.0, 2.0, 0.0}, 3.0, 2.0)
	this.lightEnts[4].Init([3]float32{10.0, 4.0, 10.0}, 4.0, 2.0)
	this.planeEnt.InitMesh(gohome.Plane("Plane", [2]float32{PLANE_SIZE * float32(NUM_PLANES), PLANE_SIZE * float32(NUM_PLANES)}, NUM_PLANES*1), &this.planeEntTobj)
	this.planeEnt.Model3D.GetMesh("Plane").GetMaterial().Shinyness = 0.5
	this.planeEnt.Model3D.GetMesh("Plane").GetMaterial().SetTextures("PlaneTexture", "", "PlaneNormalMap")
	this.planeEnt.Model3D.GetMesh("Plane").GetMaterial().SetColors(gohome.Color{255, 255, 255, 255}, gohome.Color{255, 255, 255, 255})
	gohome.ResourceMgr.GetTexture("PlaneTexture").SetWrapping(gohome.WRAPPING_MIRRORED_REPEAT)
	gohome.ResourceMgr.GetTexture("PlaneNormalMap").SetWrapping(gohome.WRAPPING_MIRRORED_REPEAT)
	this.planeEntTobj.Position = [3]float32{10.0*float32(NUM_PLANES)/2.0 - PLANE_SIZE/2.0, 0.0, PLANE_SIZE*float32(NUM_PLANES)/2.0 - PLANE_SIZE/2.0}

	this.m4.InitName("M4", &this.m4Tobj)
	this.m41.InitName("M4", &this.m41Tobj)
	this.kratos.InitName("Kratos", &this.kratosTobj)
	this.kratos.Model3D.GetMeshIndex(0).GetMaterial().SetTextures("", "", "Kratos_legs_n.tga")
	this.kratos.Model3D.GetMeshIndex(1).GetMaterial().SetTextures("", "", "Kratos_torso_n.tga")
	this.kratos.Model3D.GetMeshIndex(2).GetMaterial().SetTextures("", "", "Kratos_head_n.tga")
	for i := 0; i < int(NUM_PLANES*NUM_PLANES); i++ {
		this.oakTrees[i].InitName("Pine", &this.oakTreeTobjs[i])
		this.oakTrees[i].Model3D.GetMesh("Pine").GetMaterial().DiffuseTexture = gohome.ResourceMgr.GetTexture("Pine")
		this.oakTreeTobjs[i].Scale = [3]float32{0.5, 0.5, 0.5}
	}

	this.direct = gohome.DirectionalLight{
		DiffuseColor:   colornames.Khaki,
		SpecularColor:  colornames.Yellow,
		Direction:      [3]float32{1.0, -1.0, 0.0},
		CastsShadows:   1,
		ShadowDistance: 50.0,
	}
	this.direct1 = gohome.DirectionalLight{
		DiffuseColor:   colornames.Lime,
		SpecularColor:  colornames.Azure,
		Direction:      [3]float32{1.0, -100000000.0, 0.0},
		CastsShadows:   1,
		ShadowDistance: 50.0,
	}
	this.spot = gohome.SpotLight{
		DiffuseColor:  colornames.Blue,
		SpecularColor: colornames.White,
		Position:      [3]float32{0.0, 5.0, 0.0},
		Direction:     [3]float32{0.0, -1.0, 0.0},
		InnerCutOff:   30.0,
		OuterCutOff:   35.0,
		Attentuation: gohome.Attentuation{
			Constant: 1.0,
		},
		CastsShadows: 1,
	}

	this.tst.Init(&this.cam3d)
	this.cam3d.Position = [3]float32{0.0, 2.0, 2.0}
	this.cam3d1.Init()
	this.cam3d1.Position = this.cam3d.Position
	this.m4.NotRelativeToCamera = 0
	this.m41.NotRelativeToCamera = 0
	this.m4Tobj.Position = [3]float32{0.2, -0.2, -0.2}
	this.m4Tobj.Rotation[1] = 180.0
	this.m41Tobj.Position = [3]float32{-0.2, -0.2, -0.2}
	this.m41Tobj.Rotation[1] = 180.0
	for x := 0; x < int(NUM_PLANES); x++ {
		for y := 0; y < int(NUM_PLANES); y++ {
			this.oakTreeTobjs[x+y*int(NUM_PLANES)].Position = [3]float32{PLANE_SIZE * (float32(x) + 0.5), -0.5, PLANE_SIZE * (float32(y) + 0.5)}
			gohome.RenderMgr.AddObject(&this.oakTrees[x+y*int(NUM_PLANES)], &this.oakTreeTobjs[x+y*int(NUM_PLANES)])
		}
	}
	gohome.UpdateMgr.AddObject(&this.tst)
	gohome.RenderMgr.SetCamera3D(&this.cam3d, 0)
	gohome.RenderMgr.SetCamera3D(&this.cam3d1, 1)
	gohome.RenderMgr.AddObject(&this.m4, &this.m4Tobj)
	gohome.RenderMgr.AddObject(&this.m41, &this.m41Tobj)
	gohome.RenderMgr.AddObject(&this.kratos, &this.kratosTobj)
	gohome.RenderMgr.AddObject(&this.planeEnt, &this.planeEntTobj)
	// gohome.LightMgr.AddDirectionalLight(&this.direct, 0)
	// gohome.LightMgr.AddDirectionalLight(&this.direct1, 0)
	gohome.LightMgr.AddSpotLight(&this.spot, 0)

	nWidth, nHeight := gohome.Render.GetNativeResolution()

	gohome.RenderMgr.SetViewport3D(&gohome.Viewport{
		0,
		0, 0,
		int(nWidth) / 2, int(nHeight),
	}, 0)
	gohome.RenderMgr.SetViewport3D(&gohome.Viewport{
		1,
		int(nWidth) - int(nWidth)/4, int(nHeight) - int(nHeight)/4,
		int(nWidth) / 4, int(nHeight) / 4,
	}, 1)

	gohome.RenderMgr.SetViewport2D(&gohome.Viewport{
		0,
		0, 0,
		int(nWidth), int(nHeight),
	}, 0)

	gohome.RenderMgr.ForceShader2D = gohome.ResourceMgr.GetShader("ShadowMapRender")
}

var elapsed_time float32 = 0.0

func (this *TestScene) Update(delta_time float32) {
	this.cam3d1.Position = this.spot.Position

	if gohome.InputMgr.JustPressed(gohome.KeyB) {
		if gohome.RenderMgr.ForceShader3D != nil {
			gohome.RenderMgr.ForceShader3D = nil
		} else {
			gohome.RenderMgr.ForceShader3D = gohome.ResourceMgr.GetShader("Normal")
		}
	}
	if gohome.InputMgr.JustPressed(gohome.KeyF) {
		gohome.Framew.WindowSetFullscreen(!gohome.Framew.WindowIsFullscreen())
	}
	if gohome.InputMgr.JustPressed(gohome.KeyT) {
		if gohome.LightMgr.CurrentLightCollection == 0 {
			gohome.LightMgr.CurrentLightCollection = -1
		} else {
			gohome.LightMgr.CurrentLightCollection = 0
		}
	}
	if gohome.InputMgr.JustPressed(gohome.KeyC) {
		this.cam3d1.LookDirection = this.spot.Direction.Add([3]float32{1e-19, 1e-19, 1e-19})
	}
	if gohome.InputMgr.IsPressed(gohome.KeyX) {
		this.cam3d.Position = this.spot.Position
		this.cam3d.LookDirection = this.spot.Direction.Add([3]float32{1e-19, 1e-19, 1e-19})
	}

	elapsed_time += delta_time

	var x, y float32
	x = float32(math.Cos(float64(elapsed_time) * 0.25))
	y = float32(math.Sin(float64(elapsed_time) * 0.25))

	this.direct1.Direction = [3]float32{x, y, 0.0}
}

func (this *TestScene) Terminate() {
	this.lightEnts[0].Terminate()
	this.lightEnts[1].Terminate()
	this.lightEnts[2].Terminate()
	this.lightEnts[3].Terminate()
	this.lightEnts[4].Terminate()
}
