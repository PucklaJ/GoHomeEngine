package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

// The different types of values that can be used for instancing
const (
	VALUE_FLOAT = iota
	VALUE_VEC2  = iota
	VALUE_VEC3  = iota
	VALUE_VEC4  = iota
	VALUE_MAT2  = iota
	VALUE_MAT3  = iota
	VALUE_MAT4  = iota
)

// A mesh that will be drawn using instancing
type InstancedMesh3D interface {
	// Adds vertices and indices to the mesh
	AddVertices(vertices []Mesh3DVertex, indices []uint32)
	// Loads the vertices and indices to the GPU
	Load()
	// Calls the draw method for this mesh
	Render()
	// Cleans everything up
	Terminate()
	// Sets the material of this mesh
	SetMaterial(mat *Material)
	// Returns the material of this mesh
	GetMaterial() *Material
	// Returns the name of this mesh
	GetName() string
	// Returns the number of vertices of this mesh
	GetNumVertices() int
	// Returns the number of indices of this mesh
	GetNumIndices() int
	// Returns all the vertices of this mesh
	GetVertices() []Mesh3DVertex
	// Returns all the indices of this mesh
	GetIndices() []uint32
	// Calculates the tangents that will be used for a normal map
	CalculateTangents()
	// Returns wether the vertices have UV values
	// Checks if the uv values are zero
	HasUV() bool
	// Returns the bounding box from the lowest vertex to the highest vertex
	AABB() AxisAlignedBoundingBox
	// Creates a new mesh from this mesh
	Copy() Mesh3D
	// Returns wether Load has been called
	LoadedToGPU() bool
	// Sets the number of buffered instances
	SetNumInstances(n int)
	// Returns the number of buffered instances
	GetNumInstances() int
	// Sets number of instances that will be drawn
	SetNumUsedInstances(n int)
	// Returns the number of instances that will be drawn
	GetNumUsedInstances() int
	// Adds a value that will be used for instancing
	AddValue(valueType int)
	// Adds a value that will be used for instancing to the front
	AddValueFront(valueType int)
	// Sets the float values of index, value needs to be of size num instances
	SetF(index int, value []float32)
	// Sets the Vec2 values of index, value needs to be of size num instances
	SetV2(index int, value []mgl32.Vec2)
	// Sets the vec3 values of index, value needs to be of size num instances	
	SetV3(index int, value []mgl32.Vec3)
	// Sets the vec4 values of index, value needs to be of size num instances	
	SetV4(index int, value []mgl32.Vec4)
	// Sets the Mat2 values of index, value needs to be of size num instances	
	SetM2(index int, value []mgl32.Mat2)
	// Sets the Mat3 values of index, value needs to be of size num instances	
	SetM3(index int, value []mgl32.Mat3)
	// Sets the Mat4 values of index, value needs to be of size num instances	
	SetM4(index int, value []mgl32.Mat4)
	// Sets the name of a value	
	SetName(index int, value_type int, value string)
}

// Creates a instanced mesh from a Mesh3D
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

// An implementation of InstancedMesh3D that does nothing
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
func (*NilInstancedMesh3D) GetNumVertices() int {
	return 0
}
func (*NilInstancedMesh3D) GetNumIndices() int {
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
func (*NilInstancedMesh3D) SetNumInstances(n int) {

}
func (*NilInstancedMesh3D) GetNumInstances() int {
	return 0
}
func (*NilInstancedMesh3D) SetNumUsedInstances(n int) {

}
func (*NilInstancedMesh3D) GetNumUsedInstances() int {
	return 0
}
func (*NilInstancedMesh3D) AddValue(valueType int) {

}
func (*NilInstancedMesh3D) AddValueFront(valueType int) {

}
func (*NilInstancedMesh3D) SetF(index int, value []float32) {

}
func (*NilInstancedMesh3D) SetV2(index int, value []mgl32.Vec2) {

}
func (*NilInstancedMesh3D) SetV3(index int, value []mgl32.Vec3) {

}
func (*NilInstancedMesh3D) SetV4(index int, value []mgl32.Vec4) {

}
func (*NilInstancedMesh3D) SetM2(index int, value []mgl32.Mat2) {

}
func (*NilInstancedMesh3D) SetM3(index int, value []mgl32.Mat3) {

}
func (*NilInstancedMesh3D) SetM4(index int, value []mgl32.Mat4) {

}
func (*NilInstancedMesh3D) SetName(index int, value_type int, value string) {

}
