package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"golang.org/x/image/colornames"
)

const SIZE int = 50
const USE_INSTANCING = true

type InstancingScene struct {
	ent gohome.Entity3D
	cam gohome.Camera3D
}

func (this *InstancingScene) Init() {
	cubeMesh := gohome.Box("Box", [3]float32{1.0, 1.0, 1.0}, false)
	if USE_INSTANCING {
		gohome.ResourceMgr.LoadShaderSource("3D Instanced", gohome.ENTITY_3D_INSTANCED_SHADER_VERTEX_SOURCE_OPENGL, gohome.ENTITY_3D_NOUV_NO_SHADOWS_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
		cubeInstanced := gohome.Render.CreateInstancedMesh3D("Instanced Box")
		cubeInstanced.AddVertices(cubeMesh.GetVertices(), cubeMesh.GetIndices())
		cubeInstanced.SetNumInstances(uint32(SIZE * SIZE * SIZE))
		cubeInstanced.AddValue(gohome.VALUE_MAT4)
		cubeInstanced.SetName(0, gohome.VALUE_MAT4, "transformMatrix3D")
		cubeInstanced.Load()
		var transforms [SIZE * SIZE * SIZE]mgl32.Mat4
		for x := 0; x < SIZE; x++ {
			for y := 0; y < SIZE; y++ {
				for z := 0; z < SIZE; z++ {
					transforms[x+y*SIZE+z*SIZE*SIZE] = mgl32.Translate3D(float32(x)*2.0, float32(y)*2.0, float32(z)*2.0)
				}
			}
		}
		cubeInstanced.SetM4(0, transforms[:])
		this.ent.InitMesh(cubeInstanced)
		this.ent.Transform.Position = mgl32.Vec3{3.0, 0.0, 0.0}
		this.ent.SetShader(gohome.ResourceMgr.GetShader("3D Instanced"))
		this.ent.SetType(gohome.TYPE_3D_INSTANCED)
		gohome.RenderMgr.AddObject(&this.ent)
	} else {
		gohome.Init3DShaders()
		cubeMesh.Load()
		for x := 0; x < SIZE; x++ {
			for y := 0; y < SIZE; y++ {
				for z := 0; z < SIZE; z++ {
					var ent gohome.Entity3D
					ent.InitMesh(cubeMesh)
					ent.Transform.Position = [3]float32{float32(x) * 2.0, float32(y) * 2.0, float32(z) * 2.0}
					gohome.RenderMgr.AddObject(&ent)
				}
			}
		}
	}

	gohome.LightMgr.SetAmbientLight(colornames.Gray, 0)
	gohome.LightMgr.AddDirectionalLight(
		&gohome.DirectionalLight{
			Direction:     mgl32.Vec3{1.0, -1.0, -1.0},
			DiffuseColor:  colornames.White,
			SpecularColor: colornames.Black,
			CastsShadows:  0,
		}, 0,
	)

	var move gohome.TestCameraMovement3D
	this.cam.Init()
	move.Init(&this.cam)
	gohome.UpdateMgr.AddObject(&move)

	this.cam.Position[2] = 3.0
	gohome.RenderMgr.SetCamera3D(&this.cam, 0)
}

func (this *InstancingScene) Update(delta_time float32) {

}

func (this *InstancingScene) Terminate() {

}
