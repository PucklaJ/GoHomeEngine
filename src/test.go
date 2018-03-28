package ignore

import (
	"fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/colornames"
	"runtime"
)

const NUM_KRATOS int = 20
const NUM_KRATOS_ROWS int = 5

type TestScene struct {
	kratos         [NUM_KRATOS]gohome.Entity3D
	kratosTobj     [NUM_KRATOS]gohome.TransformableObject3D
	M4             gohome.Entity3D
	M4Tobj         gohome.TransformableObject3D
	M41            gohome.Entity3D
	M41Tobj        gohome.TransformableObject3D
	cam3d          gohome.Camera3D
	level          *gohome.Level
	planeEnt       gohome.Entity3D
	planeEntTobj   gohome.TransformableObject3D
	lightCube      gohome.Entity3D
	lightCubeTobj  gohome.TransformableObject3D
	barrel         gohome.Entity3D
	barrelTobj     gohome.TransformableObject3D
	barrelSpr      gohome.Sprite2D
	barrelSprTobj  gohome.TransformableObject2D
	barrelSpr1     gohome.Sprite2D
	barrelSprTobj1 gohome.TransformableObject2D
}

func (this *TestScene) Init() {
	gohome.InitDefaultValues(nil)

	gohome.LightMgr.AmbientLight = &gohome.Color{20, 20, 20, 255}

	gohome.ResourceMgr.LoadLevel("Kratos", "Kratos.obj", nil)
	gohome.ResourceMgr.LoadLevel("M4", "Arma_M4.obj", nil)
	gohome.ResourceMgr.LoadLevel("Barrel", "barrel.obj", nil)
	gohome.ResourceMgr.LoadTexture("BarrelTexture", "barrel.png", nil)
	gohome.ResourceMgr.LoadTexture("BarrelNormalMap", "barrelNormal.png", nil)
	gohome.ResourceMgr.LoadTexture("Kratos_legs_s.tga", "Kratos_legs_s.tga", nil)
	gohome.ResourceMgr.LoadTexture("Kratos_torso_s.tga", "Kratos_torso_s.tga", nil)
	gohome.ResourceMgr.LoadTexture("Kratos_head_s.tga", "Kratos_head_s.tga", nil)
	gohome.ResourceMgr.LoadTexture("Kratos_legs_n.tga", "Kratos_legs_n.tga", nil)
	gohome.ResourceMgr.LoadTexture("Kratos_torso_n.tga", "Kratos_torso_n.tga", nil)
	gohome.ResourceMgr.LoadTexture("Kratos_head_n.tga", "Kratos_head_n.tga", nil)

	this.M4.InitName("Cylinder047.001_Untitled.006", &this.M4Tobj)
	this.M41.InitName("Cylinder047.001_Untitled.006", &this.M41Tobj)
	this.M4.RelativeToCamera = false
	this.M41.RelativeToCamera = false
	this.M41.Model3D.GetMeshIndex(0).GetMaterial().SetColors(colornames.White, colornames.White)
	this.M41.Model3D.GetMeshIndex(0).GetMaterial().SetTextures("Kratos_head_d.tga", "", "")
	this.planeEnt.InitMesh(gohome.Plane("Plane", [2]float32{float32(NUM_KRATOS_ROWS)*2.0 + 2.0, float32(NUM_KRATOS/NUM_KRATOS_ROWS)*2.0 + 2.0}), &this.planeEntTobj)
	this.planeEnt.Model3D.GetMeshIndex(0).GetMaterial().SetTextures("Kratos_torso_d.tga", "", "Kratos_torso_n.tga")
	this.planeEntTobj.Position = [3]float32{float32(NUM_KRATOS_ROWS), 0.0, -float32(NUM_KRATOS / NUM_KRATOS_ROWS)}
	this.barrel.InitName("Barrel", &this.barrelTobj)
	this.barrel.Model3D.GetMeshIndex(0).GetMaterial().SetTextures("BarrelTexture", "", "BarrelNormalMap")
	this.barrel.Model3D.GetMeshIndex(0).GetMaterial().Shinyness = 0.7
	this.barrelTobj.Scale = mgl32.Vec3{0.2, 0.2, 0.2}
	this.barrelTobj.Position = mgl32.Vec3{5.0, 0.0, 0.0}
	this.barrelSpr.Init("BarrelNormalMap", &this.barrelSprTobj)
	this.barrelSprTobj.Scale = [2]float32{0.3, 0.3}
	this.barrelSpr1.Init("BarrelTexture", &this.barrelSprTobj1)
	this.barrelSprTobj1.Scale = [2]float32{0.3, 0.3}
	this.barrelSprTobj1.Position[1] += this.barrelSprTobj1.Size[1] * this.barrelSprTobj1.Scale[1]

	for x := 0; x < NUM_KRATOS_ROWS; x++ {
		for z := 0; z < NUM_KRATOS/NUM_KRATOS_ROWS; z++ {
			i := z + x*NUM_KRATOS/NUM_KRATOS_ROWS
			if i > NUM_KRATOS-1 {
				break
			}
			this.kratos[i].InitName("defaultobject", &this.kratosTobj[i])
			this.kratos[i].Model3D.GetMeshIndex(0).GetMaterial().SetTextures("Kratos_legs_d.tga", "Kratos_legs_s.tga", "Kratos_legs_n.tga")
			this.kratos[i].Model3D.GetMeshIndex(1).GetMaterial().SetTextures("Kratos_torso_d.tga", "Kratos_torso_s.tga", "Kratos_torso_n.tga")
			this.kratos[i].Model3D.GetMeshIndex(2).GetMaterial().SetTextures("Kratos_head_d.tga", "Kratos_head_s.tga", "Kratos_head_n.tga")
			this.kratos[i].Model3D.GetMeshIndex(0).GetMaterial().Shinyness = 1.0
			this.kratos[i].Model3D.GetMeshIndex(1).GetMaterial().Shinyness = 1.0
			this.kratos[i].Model3D.GetMeshIndex(2).GetMaterial().Shinyness = 1.0
			this.kratosTobj[i].Position[0] = float32(x) * 2.0
			this.kratosTobj[i].Position[2] = float32(z) * -2.0
			gohome.RenderMgr.AddObject(&this.kratos[i], &this.kratosTobj[i])
		}
	}

	gohome.RenderMgr.AddObject(&this.M4, &this.M4Tobj)
	gohome.RenderMgr.AddObject(&this.M41, &this.M41Tobj)
	gohome.RenderMgr.AddObject(&this.planeEnt, &this.planeEntTobj)
	gohome.RenderMgr.AddObject(&this.barrel, &this.barrelTobj)
	gohome.RenderMgr.AddObject(&this.barrelSpr, &this.barrelSprTobj)
	gohome.RenderMgr.AddObject(&this.barrelSpr1, &this.barrelSprTobj1)

	this.M4Tobj.Position = [3]float32{0.2, -0.2, -0.2}
	this.M4Tobj.Rotation[1] = 180.0
	this.M41Tobj.Position = [3]float32{-0.2, -0.2, -0.2}
	this.M41Tobj.Rotation[1] = 180.0

	tst := &gohome.TestCameraMovement3D{}
	tst.Init(&this.cam3d)
	gohome.UpdateMgr.AddObject(tst)

	gohome.RenderMgr.SetCamera3D(&this.cam3d)

	dl := gohome.DirectionalLight{
		Direction:    mgl32.Vec3{0.0, -1.0, 0.0},
		DiffuseColor: colornames.Khaki,
	}

	pl := gohome.PointLight{
		Position:      mgl32.Vec3{1.0, 1.0, -1.0},
		DiffuseColor:  colornames.White,
		SpecularColor: colornames.White,
		Attentuation: gohome.Attentuation{
			Constant: 2.0,
		},
	}
	// pl1 := gohome.PointLight{
	// 	Position:      mgl32.Vec3{3.0, 1.0, -1.0},
	// 	DiffuseColor:  colornames.White,
	// 	SpecularColor: colornames.White,
	// 	Attentuation: gohome.Attentuation{
	// 		Constant: 2.0,
	// 	},
	// }
	// pl2 := gohome.PointLight{
	// 	Position:      mgl32.Vec3{-4.0, 1.0, -1.0},
	// 	DiffuseColor:  colornames.White,
	// 	SpecularColor: colornames.White,
	// 	Attentuation: gohome.Attentuation{
	// 		Constant: 2.0,
	// 	},
	// }
	// pl3 := gohome.PointLight{
	// 	Position:      mgl32.Vec3{1.0, 2.0, -3.0},
	// 	DiffuseColor:  colornames.White,
	// 	SpecularColor: colornames.White,
	// 	Attentuation: gohome.Attentuation{
	// 		Constant: 2.0,
	// 	},
	// }

	gohome.LightMgr.AddDirectionalLight(&dl)
	gohome.LightMgr.AddPointLight(&pl)
	// gohome.LightMgr.AddPointLight(&pl1)
	// gohome.LightMgr.AddPointLight(&pl2)
	// gohome.LightMgr.AddPointLight(&pl3)

	this.lightCube.InitMesh(gohome.Box("Box", [3]float32{1.0, 1.0, 1.0}), &this.lightCubeTobj)
	this.lightCubeTobj.Scale = mgl32.Vec3{0.2, 0.2, 0.2}
	this.lightCubeTobj.Position = pl.Position

	gohome.RenderMgr.AddObject(&this.lightCube, &this.lightCubeTobj)

	gohome.FPSLimit.MaxFPS = 100
}

var elapsed_time float32 = 0.0

func (this *TestScene) Update(delta_time float32) {
	if gohome.InputMgr.JustPressed(gohome.KeyB) {
		viewMat := this.cam3d.GetViewMatrix()
		invViewMat := this.cam3d.GetInverseViewMatrix()

		fmt.Println("ViewMatrix:")
		fmt.Println(viewMat)
		fmt.Println("InverseViewMatrix:")
		fmt.Println(invViewMat)
		fmt.Println("Combined:")
		fmt.Println(viewMat.Mul4(invViewMat))
	}
	if gohome.InputMgr.JustPressed(gohome.KeyV) {
		this.M4.RelativeToCamera = !this.M4.RelativeToCamera
		this.M41.RelativeToCamera = !this.M41.RelativeToCamera
	}
	elapsed_time += delta_time

	this.barrelTobj.Rotation[1] = elapsed_time * 10.0
	for i := 0; i < NUM_KRATOS; i++ {
		this.kratosTobj[i].Rotation[1] = elapsed_time * 10.0
	}
}

func (this *TestScene) Terminate() {

}

func main() {
	runtime.LockOSThread()
	gohome.MainLoop{}.Run(&gohome.GLFWFramework{}, &gohome.OpenGLRenderer{}, 1280, 720, "GoHomeEngine", &TestScene{})
}
