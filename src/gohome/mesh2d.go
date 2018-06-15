package gohome

type Mesh2D interface {
	AddVertices(vertices []Mesh2DVertex, indices []uint32)
	Load()
	Render()
	Terminate()
}
