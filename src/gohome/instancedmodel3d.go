package gohome

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
