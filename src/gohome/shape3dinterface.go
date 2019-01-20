package gohome

type Shape3DInterface interface {
	Init()
	AddLines(lines []Line3D)
	GetLines() []Line3D
	Load()
	Render()
	Terminate()
}

type NilShape3DInterface struct {
}

func (*NilShape3DInterface) Init() {

}
func (*NilShape3DInterface) AddLines(lines []Line3D) {

}
func (*NilShape3DInterface) GetLines() []Line3D {
	var lines []Line3D
	return lines
}
func (*NilShape3DInterface) Load() {

}
func (*NilShape3DInterface) Render() {

}
func (*NilShape3DInterface) Terminate() {

}
