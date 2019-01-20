package gohome

type Shape3DInterface interface {
	Init()
	AddPoints(points []Shape3DVertex)
	GetPoints() []Shape3DVertex
	SetDrawMode(drawMode uint8)
	SetPointSize(size float32)
	SetLineWidth(size float32)
	Load()
	Render()
	Terminate()
}

type NilShape3DInterface struct {
}

func (*NilShape3DInterface) Init() {

}
func (*NilShape3DInterface) AddPoints(points []Shape3DVertex) {

}
func (*NilShape3DInterface) GetPoints() []Shape3DVertex {
	var points []Shape3DVertex
	return points
}

func (*NilShape3DInterface) SetDrawMode(drawMode uint8) {

}

func (*NilShape3DInterface) SetPointSize(size float32) {

}

func (*NilShape3DInterface) SetLineWidth(size float32) {

}

func (*NilShape3DInterface) Load() {

}
func (*NilShape3DInterface) Render() {

}
func (*NilShape3DInterface) Terminate() {

}
