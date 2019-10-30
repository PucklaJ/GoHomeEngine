package gohome

// A 2D mesh consisting of geometry used for rendering 2D
type Mesh2D interface {
	// Add vertices and indices to the mesh
	AddVertices(vertices []Mesh2DVertex, indices []uint32)
	// Loads vertices and indices to the GPU
	Load()
	// Call the draw method on the data
	Render()
	// Cleans everything up
	Terminate()
}

// An implementation of Mesh2D that does nothing
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
