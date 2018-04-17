package gohome

type Model3D struct {
	Name   string
	meshes []Mesh3D
}

func (this *Model3D) AddMesh3D(m Mesh3D) {
	this.meshes = append(this.meshes, m)
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
