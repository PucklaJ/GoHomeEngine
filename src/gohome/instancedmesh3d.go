package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

const (
	VALUE_FLOAT uint32 = iota
	VALUE_VEC2  uint32 = iota
	VALUE_VEC3  uint32 = iota
	VALUE_VEC4  uint32 = iota
	VALUE_MAT2  uint32 = iota
	VALUE_MAT3  uint32 = iota
	VALUE_MAT4  uint32 = iota
)

/*
type Mesh3D interface {
	AddVertices(vertices []Mesh3DVertex, indices []uint32)
	Load()
	Render()
	Terminate()
	SetMaterial(mat *Material)
	GetMaterial() *Material
	GetName() string
	GetNumVertices() uint32
	GetNumIndices() uint32
	GetVertices() []Mesh3DVertex
	GetIndices() []uint32
	CalculateTangents()
	HasUV() bool
	AABB() AxisAlignedBoundingBox
	Copy() Mesh3D
}
*/

type InstancedMesh3D interface {
	AddVertices(vertices []Mesh3DVertex, indices []uint32)
	Load()
	Render()
	Terminate()
	SetMaterial(mat *Material)
	GetMaterial() *Material
	GetName() string
	GetNumVertices() uint32
	GetNumIndices() uint32
	GetVertices() []Mesh3DVertex
	GetIndices() []uint32
	CalculateTangents()
	HasUV() bool
	AABB() AxisAlignedBoundingBox
	Copy() Mesh3D
	LoadedToGPU() bool
	SetNumInstances(n uint32)
	GetNumInstances() uint32
	SetNumUsedInstances(n uint32)
	GetNumUsedInstances() uint32
	AddValue(valueType uint32)
	AddValueFront(valueType uint32)
	SetF(index uint32, value []float32)
	SetV2(index uint32, value []mgl32.Vec2)
	SetV3(index uint32, value []mgl32.Vec3)
	SetV4(index uint32, value []mgl32.Vec4)
	SetM2(index uint32, value []mgl32.Mat2)
	SetM3(index uint32, value []mgl32.Mat3)
	SetM4(index uint32, value []mgl32.Mat4)
	SetName(index uint32, value_type uint32, value string)
}

func InstancedMesh3DFromMesh3D(mesh Mesh3D) (imesh InstancedMesh3D) {
	if !mesh.LoadedToGPU() {
		imesh = Render.CreateInstancedMesh3D(mesh.GetName())
		imesh.AddVertices(mesh.GetVertices(), mesh.GetIndices())
	} else {
		imesh = Render.InstancedMesh3DFromLoadedMesh3D(mesh)
	}

	mat := mesh.GetMaterial()
	mats := *mat
	imesh.SetMaterial(&mats)

	return
}

type NilInstancedMesh3D struct {
}

func (*NilInstancedMesh3D) AddVertices(vertices []Mesh3DVertex, indices []uint32) {

}
func (*NilInstancedMesh3D) Load() {

}
func (*NilInstancedMesh3D) Render() {

}
func (*NilInstancedMesh3D) Terminate() {

}
func (*NilInstancedMesh3D) SetMaterial(mat *Material) {

}
func (*NilInstancedMesh3D) GetMaterial() *Material {
	var mat Material
	mat.InitDefault()
	return &mat
}
func (*NilInstancedMesh3D) GetName() string {
	return ""
}
func (*NilInstancedMesh3D) GetNumVertices() uint32 {
	return 0
}
func (*NilInstancedMesh3D) GetNumIndices() uint32 {
	return 0
}
func (*NilInstancedMesh3D) GetVertices() []Mesh3DVertex {
	var verts []Mesh3DVertex
	return verts
}
func (*NilInstancedMesh3D) GetIndices() []uint32 {
	var inds []uint32
	return inds
}
func (*NilInstancedMesh3D) CalculateTangents() {

}
func (*NilInstancedMesh3D) HasUV() bool {
	return true
}
func (*NilInstancedMesh3D) AABB() AxisAlignedBoundingBox {
	return AxisAlignedBoundingBox{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 0, 0},
	}
}
func (*NilInstancedMesh3D) Copy() Mesh3D {
	return &NilMesh3D{}
}
func (*NilInstancedMesh3D) LoadedToGPU() bool {
	return true
}
func (*NilInstancedMesh3D) SetNumInstances(n uint32) {

}
func (*NilInstancedMesh3D) GetNumInstances() uint32 {
	return 0
}
func (*NilInstancedMesh3D) SetNumUsedInstances(n uint32) {

}
func (*NilInstancedMesh3D) GetNumUsedInstances() uint32 {
	return 0
}
func (*NilInstancedMesh3D) AddValue(valueType uint32) {

}
func (*NilInstancedMesh3D) AddValueFront(valueType uint32) {

}
func (*NilInstancedMesh3D) SetF(index uint32, value []float32) {

}
func (*NilInstancedMesh3D) SetV2(index uint32, value []mgl32.Vec2) {

}
func (*NilInstancedMesh3D) SetV3(index uint32, value []mgl32.Vec3) {

}
func (*NilInstancedMesh3D) SetV4(index uint32, value []mgl32.Vec4) {

}
func (*NilInstancedMesh3D) SetM2(index uint32, value []mgl32.Mat2) {

}
func (*NilInstancedMesh3D) SetM3(index uint32, value []mgl32.Mat3) {

}
func (*NilInstancedMesh3D) SetM4(index uint32, value []mgl32.Mat4) {

}
func (*NilInstancedMesh3D) SetName(index uint32, value_type uint32, value string) {

}
