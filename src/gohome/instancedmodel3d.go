package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type InstancedModel3D struct {
	Name   string
	meshes []InstancedMesh3D
	AABB   AxisAlignedBoundingBox
}

func (this *InstancedModel3D) AddMesh3D(m InstancedMesh3D) {
	this.meshes = append(this.meshes, m)
	this.checkAABB(m)
}

func (this *InstancedModel3D) Render() {
	for i := 0; i < len(this.meshes); i++ {
		this.meshes[i].Render()
	}
}

func (this *InstancedModel3D) Terminate() {
	for i := 0; i < len(this.meshes); i++ {
		this.meshes[i].Terminate()
	}

	this.meshes = append(this.meshes[:0], this.meshes[len(this.meshes):]...)
}

func (this *InstancedModel3D) GetMesh(name string) InstancedMesh3D {
	for i := 0; i < len(this.meshes); i++ {
		if this.meshes[i].GetName() == name {
			return this.meshes[i]
		}
	}

	return nil
}

func (this *InstancedModel3D) GetMeshIndex(index uint32) InstancedMesh3D {
	if index > uint32(len(this.meshes)-1) {
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

func (this *InstancedModel3D) HasUV() bool {
	for i := 0; i < len(this.meshes); i++ {
		if !this.meshes[i].HasUV() {
			return false
		}
	}
	return true
}

func (this *InstancedModel3D) AddValue(valueType uint32) {
	for _, m := range this.meshes {
		m.AddValue(valueType)
	}
}

func (this *InstancedModel3D) AddValueFront(valueType uint32) {
	for _, m := range this.meshes {
		m.AddValueFront(valueType)
	}
}

func (this *InstancedModel3D) SetName(index, valueType uint32, value string) {
	for _, m := range this.meshes {
		m.SetName(index, valueType, value)
	}
}

func (this *InstancedModel3D) SetF(index uint32, value []float32) {
	for _, m := range this.meshes {
		m.SetF(index, value)
	}
}

func (this *InstancedModel3D) SetV2(index uint32, value []mgl32.Vec2) {
	for _, m := range this.meshes {
		m.SetV2(index, value)
	}
}

func (this *InstancedModel3D) SetV3(index uint32, value []mgl32.Vec3) {
	for _, m := range this.meshes {
		m.SetV3(index, value)
	}
}

func (this *InstancedModel3D) SetV4(index uint32, value []mgl32.Vec4) {
	for _, m := range this.meshes {
		m.SetV4(index, value)
	}
}

func (this *InstancedModel3D) SetM2(index uint32, value []mgl32.Mat2) {
	for _, m := range this.meshes {
		m.SetM2(index, value)
	}
}

func (this *InstancedModel3D) SetM3(index uint32, value []mgl32.Mat3) {
	for _, m := range this.meshes {
		m.SetM3(index, value)
	}
}

func (this *InstancedModel3D) SetM4(index uint32, value []mgl32.Mat4) {
	for _, m := range this.meshes {
		m.SetM4(index, value)
	}
}

func (this *InstancedModel3D) Load() {
	for _, m := range this.meshes {
		m.Load()
	}
}

func (this *InstancedModel3D) SetNumInstances(n uint32) {
	for _, m := range this.meshes {
		m.SetNumInstances(n)
	}
}

func (this *InstancedModel3D) GetNumInstances() uint32 {
	if len(this.meshes) == 0 {
		return 0
	} else {
		return this.meshes[0].GetNumInstances()
	}
}

func (this *InstancedModel3D) SetNumUsedInstances(n uint32) {
	for _, m := range this.meshes {
		m.SetNumUsedInstances(n)
	}
}

func (this *InstancedModel3D) GetNumUsedInstances() uint32 {
	if len(this.meshes) == 0 {
		return 0
	} else {
		return this.meshes[0].GetNumUsedInstances()
	}
}

func (this *InstancedModel3D) LoadedToGPU() bool {
	for _, m := range this.meshes {
		if !m.LoadedToGPU() {
			return false
		}
	}
	return true
}

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
