package gohome

const (
	MESH2DVERTEX_SIZE uint32 = 2 * 2 * 4 // 2*2*sizeof(float32)
	INDEX_SIZE        uint32 = 4         // sizeof(uint32)
)

type Mesh2DVertex [4]float32

func (m *Mesh2DVertex) Vertex(x, y float32) {
	m[0] = x
	m[1] = y
}

func (m *Mesh2DVertex) TexCoord(u, v float32) {
	m[2] = u
	m[3] = v
}

type Mesh2D interface {
	AddVertices(vertices []Mesh2DVertex, indices []uint32)
	Load()
	Render()
	Terminate()
}
