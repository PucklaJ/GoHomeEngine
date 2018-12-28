package gohome

import (
	"image"
	"image/color"
)

type CubeMap interface {
	Load(data []byte, width, height int, shadowMap bool) error
	LoadFromImage(img image.Image) error
	Bind(unit uint32)
	Unbind(unit uint32)
	GetWidth() int
	GetHeight() int
	GetKeyColor() color.Color
	GetModColor() color.Color
	Terminate()
	SetFiltering(filtering uint32)
	SetWrapping(wrapping uint32)
	SetBorderColor(col color.Color)
	SetBorderDepth(depth float32)
	SetKeyColor(col color.Color)
	SetModColor(col color.Color)
	GetName() string
	GetData() ([]byte, int, int)
}

type NilCubeMap struct {
}

func (*NilCubeMap) Load(data []byte, width, height int, shadowMap bool) error {
	return nil
}
func (*NilCubeMap) LoadFromImage(img image.Image) error {
	return nil
}
func (*NilCubeMap) Bind(unit uint32) {

}
func (*NilCubeMap) Unbind(unit uint32) {

}
func (*NilCubeMap) GetWidth() int {
	return 0
}
func (*NilCubeMap) GetHeight() int {
	return 0
}
func (*NilCubeMap) GetKeyColor() color.Color {
	return nil
}
func (*NilCubeMap) GetModColor() color.Color {
	return nil
}
func (*NilCubeMap) Terminate() {

}
func (*NilCubeMap) SetFiltering(filtering uint32) {

}
func (*NilCubeMap) SetWrapping(wrapping uint32) {

}
func (*NilCubeMap) SetBorderColor(col color.Color) {

}
func (*NilCubeMap) SetBorderDepth(depth float32) {

}
func (*NilCubeMap) SetKeyColor(col color.Color) {

}
func (*NilCubeMap) SetModColor(col color.Color) {

}
func (*NilCubeMap) GetName() string {
	return ""
}
func (*NilCubeMap) GetData() ([]byte, int, int) {
	var data []byte
	return data, 0, 0
}
