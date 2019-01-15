package gohome

type Mesh2D interface {
	AddVertices(vertices []Mesh2DVertex, indices []uint32)
	Load()
	Render()
	Terminate()
}

type NilMesh2D struct {
}

func (*NilMesh2D) AddVertices(vertices []Mesh2DVertex, indices []uint32) {

}
func (*NilMesh2D) Load() {

}
func (*NilMesh2D) Render() {

}
func (*NilMesh2D) Terminate() {

}
