package gohome

type Model3D struct {
	Name   string
	meshes []Mesh3D
	AABB   AxisAlignedBoundingBox
}

func (this *Model3D) AddMesh3D(m Mesh3D) {
	this.meshes = append(this.meshes, m)
	this.checkAABB(m)
}

func (this *Model3D) Render() {
	for i := 0; i < len(this.meshes); i++ {
		this.meshes[i].Render()
	}
}

func (this *Model3D) Terminate() {
	for i := 0; i < len(this.meshes); i++ {
		this.meshes[i].Terminate()
	}

	this.meshes = append(this.meshes[:0], this.meshes[len(this.meshes):]...)
}

func (this *Model3D) GetMesh(name string) Mesh3D {
	for i := 0; i < len(this.meshes); i++ {
		if this.meshes[i].GetName() == name {
			return this.meshes[i]
		}
	}

	return nil
}

func (this *Model3D) GetMeshIndex(index uint32) Mesh3D {
	if index > uint32(len(this.meshes)-1) {
		return nil
	} else {
		return this.meshes[index]
	}
}

func (this *Model3D) checkAABB(m Mesh3D) {
	for i := 0; i < 3; i++ {
		if m.AABB().Max[i] > this.AABB.Max[i] {
			this.AABB.Max[i] = m.AABB().Max[i]
		}
		if m.AABB().Min[i] < this.AABB.Min[i] {
			this.AABB.Min[i] = m.AABB().Min[i]
		}
	}
}

func (this *Model3D) HasUV() bool {
	for i := 0; i < len(this.meshes); i++ {
		if !this.meshes[i].HasUV() {
			return false
		}
	}
	return true
}

func (this *Model3D) Copy() *Model3D {
	var model Model3D
	model.Name = this.Name + " Copy"
	model.meshes = make([]Mesh3D, len(this.meshes))
	for i := 0; i < len(this.meshes); i++ {
		model.meshes[i] = this.meshes[i].Copy()
	}
	model.AABB = this.AABB
	return &model
}

func (this *Model3D) Load() {
	for _, m := range this.meshes {
		m.Load()
		mat := m.GetMaterial()
		if mat.DiffuseTexture != nil {
			data := preloadedTextureData[mat.DiffuseTexture]
			mat.DiffuseTexture.Load(data.data, data.width, data.height, false)
		}
		if mat.SpecularTexture != nil {
			data := preloadedTextureData[mat.SpecularTexture]
			mat.SpecularTexture.Load(data.data, data.width, data.height, false)
		}
		if mat.NormalMap != nil {
			data := preloadedTextureData[mat.NormalMap]
			mat.NormalMap.Load(data.data, data.width, data.height, false)
		}
	}
}
