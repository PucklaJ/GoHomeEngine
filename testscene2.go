package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
)

type TestScene2 struct {
	kratosMesh  [3]gohome.InstancedMesh3D
	kratosModel gohome.Model3D
	kratosTobjs [2]gohome.TransformableObject3D
	kratosEnt   gohome.Entity3D
}

func (this *TestScene2) toMat4Slice() []mgl32.Mat4 {
	rv := make([]mgl32.Mat4, len(this.kratosTobjs))

	for i := 0; i < len(rv); i++ {
		rv[i] = this.kratosTobjs[i].GetTransformMatrix()
	}

	return rv
}

func (this *TestScene2) Init() {
	gohome.InitDefaultValues()
	gohome.FPSLimit.MaxFPS = 1000

	gohome.ResourceMgr.PreloadLevel("Kratos", "Kratos.obj", false)
	gohome.ResourceMgr.PreloadShader("Instanced3D", "instanced3dVert.glsl", "fragment3d.glsl", "", "", "", "")
	gohome.ResourceMgr.LoadPreloadedResources()

	if gohome.ResourceMgr.GetLevel("Kratos") != nil {
		this.kratosTobjs[0].Position = [3]float32{0.0, -1.0, -5.0}
		this.kratosTobjs[0].Scale = [3]float32{1.0, 1.0, 1.0}
		this.kratosTobjs[0].Rotation = [3]float32{0.0, 0.0, 0.0}
		this.kratosTobjs[0].CalculateTransformMatrix(nil, -1)
		this.kratosTobjs[1].Position = [3]float32{1.0, -1.0, -5.0}
		this.kratosTobjs[1].Scale = [3]float32{1.0, 1.0, 1.0}
		this.kratosTobjs[1].Rotation = [3]float32{20.0, 0.0, 0.0}
		this.kratosTobjs[1].CalculateTransformMatrix(nil, -1)

		model := gohome.ResourceMgr.GetLevel("Kratos").GetModel("Kratos")
		for i := 0; model.GetMeshIndex(uint32(i)) != nil; i++ {
			mesh := model.GetMeshIndex(uint32(i))
			this.kratosMesh[i] = gohome.Render.CreateInstancedMesh3D("Kratos")
			this.kratosMesh[i].AddVertices(mesh.GetVertices(), mesh.GetIndices())
			this.kratosMesh[i].SetNumInstances(2)
			this.kratosMesh[i].AddValue(gohome.VALUE_MAT4)
			this.kratosMesh[i].Load()
			this.kratosMesh[i].SetName(0, gohome.VALUE_MAT4, "transformMatrix3D")
			this.kratosMesh[i].SetM4(0, this.toMat4Slice())
			this.kratosMesh[i].SetMaterial(mesh.GetMaterial())
			this.kratosModel.AddMesh3D(this.kratosMesh[i])
		}
		this.kratosEnt.InitModel(&this.kratosModel, nil)
		this.kratosEnt.SetShader(gohome.ResourceMgr.GetShader("Instanced3D"))
		gohome.RenderMgr.AddObject(&this.kratosEnt, nil)
	}

	gohome.LightMgr.CurrentLightCollection = -1
	gohome.RenderMgr.EnableBackBuffer = false
}

func (this *TestScene2) Update(delta_time float32) {
	for i := 0; i < len(this.kratosTobjs); i++ {
		this.kratosTobjs[i].Rotation[0] += 30.0 * delta_time
		this.kratosTobjs[i].Rotation[1] += 30.0 * delta_time
		this.kratosTobjs[i].CalculateTransformMatrix(nil, -1)
	}
	for i := 0; i < 3; i++ {
		this.kratosMesh[i].SetM4(0, this.toMat4Slice())
	}
}

func (this *TestScene2) Terminate() {

}
