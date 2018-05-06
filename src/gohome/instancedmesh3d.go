package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
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
	CalculateTangents()
	SetNumInstances(n uint32)
	GetNumInstances() uint32
	GetVertices() []Mesh3DVertex
	GetIndices() []uint32
	AddValue(valueType uint32)
	SetF(index uint32, value []float32)
	SetV2(index uint32, value []mgl32.Vec2)
	SetV3(index uint32, value []mgl32.Vec3)
	SetV4(index uint32, value []mgl32.Vec4)
	SetM2(index uint32, value []mgl32.Mat2)
	SetM3(index uint32, value []mgl32.Mat3)
	SetM4(index uint32, value []mgl32.Mat4)
	SetName(index uint32, value_type uint32, value string)
}
