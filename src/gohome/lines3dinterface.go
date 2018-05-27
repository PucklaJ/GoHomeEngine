package gohome

type Lines3DInterface interface {
	Init()
	AddLines(lines []Line3D)
	GetLines() []Line3D
	Load()
	Render()
	Terminate()
}
