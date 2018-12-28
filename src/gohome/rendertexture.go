package gohome

import (
	"image"
	"image/color"
)

type RenderTexture interface {
	Load(data []byte, width, height int, shadowMap bool) error // Is not used. It there just make RenderTexture able to be a Texture
	LoadFromImage(img image.Image) error
	GetName() string
	SetAsTarget()
	UnsetAsTarget()
	Blit(rtex RenderTexture)
	Bind(unit uint32)
	Unbind(unit uint32)
	GetWidth() int
	GetHeight() int
	GetKeyColor() color.Color
	GetModColor() color.Color
	ChangeSize(width, height uint32)
	Terminate()
	SetFiltering(filtering uint32)
	SetWrapping(wrapping uint32)
	SetBorderColor(col color.Color)
	SetBorderDepth(depth float32)
	SetKeyColor(col color.Color)
	SetModColor(col color.Color)
	GetData() ([]byte, int, int)
}

type NilRenderTexture struct {
}

func (*NilRenderTexture) Load(data []byte, width, height int, shadowMap bool) error {
	return nil
}
func (*NilRenderTexture) LoadFromImage(img image.Image) error {
	return nil
}
func (*NilRenderTexture) GetName() string {
	return ""
}
func (*NilRenderTexture) SetAsTarget() {

}
func (*NilRenderTexture) UnsetAsTarget() {

}
func (*NilRenderTexture) Blit(rtex RenderTexture) {

}
func (*NilRenderTexture) Bind(unit uint32) {

}
func (*NilRenderTexture) Unbind(unit uint32) {

}
func (*NilRenderTexture) GetWidth() int {
	return 0
}
func (*NilRenderTexture) GetHeight() int {
	return 0
}
func (*NilRenderTexture) GetKeyColor() color.Color {
	return nil
}
func (*NilRenderTexture) GetModColor() color.Color {
	return nil
}
func (*NilRenderTexture) ChangeSize(width, height uint32) {

}
func (*NilRenderTexture) Terminate() {

}
func (*NilRenderTexture) SetFiltering(filtering uint32) {

}
func (*NilRenderTexture) SetWrapping(wrapping uint32) {

}
func (*NilRenderTexture) SetBorderColor(col color.Color) {

}
func (*NilRenderTexture) SetBorderDepth(depth float32) {

}
func (*NilRenderTexture) SetKeyColor(col color.Color) {

}
func (*NilRenderTexture) SetModColor(col color.Color) {

}
func (*NilRenderTexture) GetData() ([]byte, int, int) {
	var data []byte
	return data, 0, 0
}
