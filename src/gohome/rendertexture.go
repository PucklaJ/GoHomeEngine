package gohome

import (
	"image"
	"image/color"
)

// A Texture to which you can render
type RenderTexture interface {
	// Is not used. It there just make RenderTexture able to be a Texture
	Load(data []byte, width, height int, shadowMap bool)
	// Is not used. It there just make RenderTexture able to be a Texture
	LoadFromImage(img image.Image)
	// Returns the name of the texture
	GetName() string
	// Renders everything after that to this render texture
	SetAsTarget()
	// Renders everything to the previously set render target or to the screen
	UnsetAsTarget()
	// Copies the contents of this render texture to rtex
	Blit(rtex RenderTexture)
	// Binds this texture to unit
	Bind(unit uint32)
	// Unbinds this texture
	Unbind(unit uint32)
	// Returns the width of the texture in pixels
	GetWidth() int
	// Returns the height of the texture in pixels
	GetHeight() int
	// Returns the key color
	GetKeyColor() color.Color
	// Returns the modulate color
	GetModColor() color.Color
	// Recreates the texture with a new size
	ChangeSize(width, height int)
	// Cleans everything up
	Terminate()
	// Sets the filter method used for this texture
	SetFiltering(filtering int)
	// Sets the wrapping method used for this texture
	SetWrapping(wrapping int)
	// Sets the border color for this texture
	SetBorderColor(col color.Color)
	// Sets the border depth for this texture
	SetBorderDepth(depth float32)
	// Sets the key color which tells the texture which color should be ignored
	SetKeyColor(col color.Color)
	// Sets the modulate color
	SetModColor(col color.Color)
	// Returns the pixels of the texture as a byte array
	GetData() ([]byte, int, int)
}

// An implementation of RenderTexture that does nothing
type NilRenderTexture struct {
}

func (*NilRenderTexture) Load(data []byte, width, height int, shadowMap bool) {
}
func (*NilRenderTexture) LoadFromImage(img image.Image) {
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
func (*NilRenderTexture) ChangeSize(width, height int) {

}
func (*NilRenderTexture) Terminate() {

}
func (*NilRenderTexture) SetFiltering(filtering int) {

}
func (*NilRenderTexture) SetWrapping(wrapping int) {

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
