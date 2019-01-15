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

type NilShape2DInterface struct {
}

func (*NilShape2DInterface) Init() {

}
func (*NilShape2DInterface) AddPoints(points []Shape2DVertex) {

}
func (*NilShape2DInterface) AddLines(lines []Line2D) {

}
func (*NilShape2DInterface) AddTriangles(tris []Triangle2D) {

}
func (*NilShape2DInterface) GetPoints() []Shape2DVertex {
	var points []Shape2DVertex
	return points
}
func (*NilShape2DInterface) SetDrawMode(mode uint8) {

}
func (*NilShape2DInterface) SetPointSize(size float32) {

}
func (*NilShape2DInterface) SetLineWidth(width float32) {

}
func (*NilShape2DInterface) Load() {

}
func (*NilShape2DInterface) Render() {

}
func (*NilShape2DInterface) Terminate() {

}
