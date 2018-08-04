package gohome

const (
	DRAW_MODE_POINTS    uint8 = iota
	DRAW_MODE_LINES     uint8 = iota
	DRAW_MODE_TRIANGLES uint8 = iota
)

type Shape2DInterface interface {
	Init()
	AddPoints(points []Shape2DVertex)
	AddLines(lines []Line2D)
	AddTriangles(tris []Triangle2D)
	GetPoints() []Shape2DVertex
	SetDrawMode(mode uint8)
	SetPointSize(size float32)
	SetLineWidth(width float32)
	Load()
	Render()
	Terminate()
}
