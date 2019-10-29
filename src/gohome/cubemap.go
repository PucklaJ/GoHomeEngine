package gohome

import (
	"image"
	"image/color"
)

// This interface represents a cube map with six faces
type CubeMap interface {
	// Loads a cube map from data with width and height
	Load(data []byte, width, height int, shadowMap bool)
	// Loads the cube map from an image
	LoadFromImage(img image.Image)
	// Binds the cube map to unit
	Bind(unit uint32)
	// Unbinds cube map
	Unbind(unit uint32)
	// Returns the width of the cube map
	GetWidth() int
	// Returns the height of the cube map
	GetHeight() int
	// Returns the color that will be keyed
	GetKeyColor() color.Color
	// Returns the modulate color
	GetModColor() color.Color
	// Cleans everything up
	Terminate()
	// Sets filtering to the give method
	SetFiltering(filtering int)
	// Sets wrapping to the given method
	SetWrapping(wrapping int)
	// Sets the color of the border to col if used
	SetBorderColor(col color.Color)
	// Sets the depth of the border
	SetBorderDepth(depth float32)
	// Sets the key color
	SetKeyColor(col color.Color)
	// Sets the modulate color
	SetModColor(col color.Color)
	// Returns the name of this cube map
	GetName() string
	// Returns the data its with and its height
	GetData() ([]byte, int, int)
}

// An implementation of CubeMap that does nothing
type NilCubeMap struct {
}

func (*NilCubeMap) Load(data []byte, width, height int, shadowMap bool) {
}
func (*NilCubeMap) LoadFromImage(img image.Image) {
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
func (*NilCubeMap) SetFiltering(filtering int) {

}
func (*NilCubeMap) SetWrapping(wrapping int) {

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
