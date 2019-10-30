package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

// A Model3D that renders instanced
type InstancedModel3D struct {
	// The name of this object
	Name   string
	meshes []InstancedMesh3D
	// The bounding box of this Model
	AABB   AxisAlignedBoundingBox
}

// Adds a mesh to the drawn meshes
func (this *InstancedModel3D) AddMesh3D(m InstancedMesh3D) {
	this.meshes = append(this.meshes, m)
	this.checkAABB(m)
}

// Calls the Render method on all meshes
func (this *InstancedModel3D) Render() {
	for i := 0; i < len(this.meshes); i++ {
		this.meshes[i].Render()
	}
}

// Cleans everything up
func (this *InstancedModel3D) Terminate() {
	for i := 0; i < len(this.meshes); i++ {
		this.meshes[i].Terminate()
	}

	this.meshes = append(this.meshes[:0], this.meshes[len(this.meshes):]...)
}

// Returns the mesh with the given name
func (this *InstancedModel3D) GetMesh(name string) InstancedMesh3D {
	for i := 0; i < len(this.meshes); i++ {
		if this.meshes[i].GetName() == name {
			return this.meshes[i]
		}
	}

	return nil
}

// Returns the mesh with the given index
func (this *InstancedModel3D) GetMeshIndex(index int) InstancedMesh3D {
	if index > len(this.meshes)-1 {
		return nil
	} else {
		return this.meshes[index]
	}
}

func (this *InstancedModel3D) checkAABB(m InstancedMesh3D) {
	for i := 0; i < 3; i++ {
		if m.AABB().Max[i] > this.AABB.Max[i] {
			this.AABB.Max[i] = m.AABB().Max[i]
		}
		if m.AABB().Min[i] < this.AABB.Min[i] {
			this.AABB.Min[i] = m.AABB().Min[i]
		}
	}
}

// Returns wether some Meshes have UV values
func (this *InstancedModel3D) HasUV() bool {
	for i := 0; i < len(this.meshes); i++ {
		if !this.meshes[i].HasUV() {
			return false
		}
	}
	return true
}

// Adds an instanced value
func (this *InstancedModel3D) AddValue(valueType int) {
	for _, m := range this.meshes {
		m.AddValue(valueType)
	}
}

// Adds an instanced value to the front of the array
func (this *InstancedModel3D) AddValueFront(valueType int) {
	for _, m := range this.meshes {
		m.AddValueFront(valueType)
	}
}

// Sets the name of an instanced value
func (this *InstancedModel3D) SetName(index, valueType int, value string) {
	for _, m := range this.meshes {
		m.SetName(index, valueType, value)
	}
}

// Sets the float value of index to value, value must be of size num instances
func (this *InstancedModel3D) SetF(index int, value []float32) {
	for _, m := range this.meshes {
		m.SetF(index, value)
	}
}

// Sets the Vec2 value of index to value, value must be of size num instances
func (this *InstancedModel3D) SetV2(index int, value []mgl32.Vec2) {
	for _, m := range this.meshes {
		m.SetV2(index, value)
	}
}

// Sets the Vec3 value of index to value, value must be of size num instances
func (this *InstancedModel3D) SetV3(index int, value []mgl32.Vec3) {
	for _, m := range this.meshes {
		m.SetV3(index, value)
	}
}

// Sets the Vec4 value of index to value, value must be of size num instances
func (this *InstancedModel3D) SetV4(index int, value []mgl32.Vec4) {
	for _, m := range this.meshes {
		m.SetV4(index, value)
	}
}

// Sets the Mat2 value of index to value, value must be of size num instances
func (this *InstancedModel3D) SetM2(index int, value []mgl32.Mat2) {
	for _, m := range this.meshes {
		m.SetM2(index, value)
	}
}

// Sets the Mat3 value of index to value, value must be of size num instances
func (this *InstancedModel3D) SetM3(index int, value []mgl32.Mat3) {
	for _, m := range this.meshes {
		m.SetM3(index, value)
	}
}

// Sets the Mat4 value of index to value, value must be of size num instances
func (this *InstancedModel3D) SetM4(index int, value []mgl32.Mat4) {
	for _, m := range this.meshes {
		m.SetM4(index, value)
	}
}

// Loads all the data to the GPU
func (this *InstancedModel3D) Load() {
	for _, m := range this.meshes {
		m.Load()
	}
}

// Sets the number of buffered instances
func (this *InstancedModel3D) SetNumInstances(n int) {
	for _, m := range this.meshes {
		m.SetNumInstances(n)
	}
}

// Returns the number of buffered instances
func (this *InstancedModel3D) GetNumInstances() int {
	if len(this.meshes) == 0 {
		return 0
	} else {
		return this.meshes[0].GetNumInstances()
	}
}

// Sets the number of drawn instances
func (this *InstancedModel3D) SetNumUsedInstances(n int) {
	for _, m := range this.meshes {
		m.SetNumUsedInstances(n)
	}
}

// Returns the number of drawn instances
func (this *InstancedModel3D) GetNumUsedInstances() int {
	if len(this.meshes) == 0 {
		return 0
	} else {
		return this.meshes[0].GetNumUsedInstances()
	}
}

// Returns wether Load has been called
func (this *InstancedModel3D) LoadedToGPU() bool {
	for _, m := range this.meshes {
		if !m.LoadedToGPU() {
			return false
		}
	}
	return true
}

// Creates an instanced model from a Model3D
func InstancedModel3DFromModel3D(m *Model3D) (im *InstancedModel3D) {
	im = &InstancedModel3D{}
	im.meshes = make([]InstancedMesh3D, len(m.meshes))
	for i := 0; i < len(m.meshes); i++ {
		im.meshes[i] = InstancedMesh3DFromMesh3D(m.meshes[i])
	}
	im.Name = m.Name
	im.AABB = m.AABB
	return
}
