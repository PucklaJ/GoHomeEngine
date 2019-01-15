package gohome

type Lines3DInterface interface {
	Init()
	AddLines(lines []Line3D)
	GetLines() []Line3D
	Load()
	Render()
	Terminate()
}

type NilLines3DInterface struct {
}

func (*NilLines3DInterface) Init() {

}
func (*NilLines3DInterface) AddLines(lines []Line3D) {

}
func (*NilLines3DInterface) GetLines() []Line3D {
	var lines []Line3D
	return lines
}
func (*NilLines3DInterface) Load() {

}
func (*NilLines3DInterface) Render() {

}
func (*NilLines3DInterface) Terminate() {

}
