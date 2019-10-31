package gohome

import (
	"image"
	"image/color"
	"sync"
)

// The different filtering and wrapping methods based on OpenGL
const (
	FILTERING_NEAREST                = iota
	FILTERING_LINEAR                 = iota
	FILTERING_NEAREST_MIPMAP_NEAREST = iota
	FILTERING_LINEAR_MIPMAP_LINEAR   = iota

	WRAPPING_REPEAT          = iota
	WRAPPING_CLAMP_TO_BORDER = iota
	WRAPPING_CLAMP_TO_EDGE   = iota
	WRAPPING_MIRRORED_REPEAT = iota
)

// A texture or an image in memory
type Texture interface {
	// Creates the texture from pixels and its dimensions
	Load(data []byte, width, height int, shadowMap bool)
	// Creates the texture from an image
	LoadFromImage(img image.Image)
	// Binds the texture to a binding point
	Bind(unit uint32)
	// Unbinds the texture (unit needs to be the same as in Bind)
	Unbind(unit uint32)
	// Returns the width of the texture in pixels
	GetWidth() int
	// Returns the height of the texture in pixels
	GetHeight() int
	// Returns the key color that will be ignored when rendering
	GetKeyColor() color.Color
	// Returns the modulate color that will be multiplied with the texture's color
	GetModColor() color.Color
	// Cleans everything up
	Terminate()
	// Sets the filtering method
	SetFiltering(filtering int)
	// Sets the wrapping method
	SetWrapping(wrapping int)
	// Sets the border color used with WRAPPING_CLAMP_TO_BORDER
	SetBorderColor(col color.Color)
	// Sets the border depth used for shadow maps
	SetBorderDepth(depth float32)
	// Sets the key color
	SetKeyColor(col color.Color)
	// Sets the modulate color
	SetModColor(col color.Color)
	// Returns the name of this texture
	GetName() string
	// Returns the pixels and its dimensions
	GetData() ([]byte, int, int)
}

// Converts a texture to an image
func TextureToImage(tex Texture, flipX, flipY bool) image.Image {
	var wg sync.WaitGroup
	data, width, height := tex.GetData()
	if data == nil || len(data) == 0 {
		return nil
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	wg.Add(width * height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			go func(_x, _y int) {
				var img_x, img_y int

				if flipX {
					img_x = width - _x - 1
				} else {
					img_x = _x
				}

				if flipY {
					img_y = height - _y - 1
				} else {
					img_y = _y
				}

				img.SetRGBA(
					img_x, img_y,
					color.RGBA{
						R: data[(_x+_y*width)*4+0],
						G: data[(_x+_y*width)*4+1],
						B: data[(_x+_y*width)*4+2],
						A: data[(_x+_y*width)*4+3],
					},
				)
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()

	return img
}

// Converts pixel data to a color
func GetColorFromData(x, y int, data []byte, width int) color.Color {
	return &Color{
		R: data[(x+y*width)*4+0],
		G: data[(x+y*width)*4+1],
		B: data[(x+y*width)*4+2],
		A: data[(x+y*width)*4+3],
	}
}

// An implementation of Texture that does nothing
type NilTexture struct {
}

func (*NilTexture) Load(data []byte, width, height int, shadowMap bool) {
}
func (*NilTexture) LoadFromImage(img image.Image) {
}
func (*NilTexture) Bind(unit uint32) {

}
func (*NilTexture) Unbind(unit uint32) {

}
func (*NilTexture) GetWidth() int {
	return 0
}
func (*NilTexture) GetHeight() int {
	return 0
}
func (*NilTexture) GetKeyColor() color.Color {
	return nil
}
func (*NilTexture) GetModColor() color.Color {
	return nil
}
func (*NilTexture) Terminate() {

}
func (*NilTexture) SetFiltering(filtering int) {

}
func (*NilTexture) SetWrapping(wrapping int) {

}
func (*NilTexture) SetBorderColor(col color.Color) {

}
func (*NilTexture) SetBorderDepth(depth float32) {

}
func (*NilTexture) SetKeyColor(col color.Color) {

}
func (*NilTexture) SetModColor(col color.Color) {

}
func (*NilTexture) GetName() string {
	return ""
}
func (*NilTexture) GetData() ([]byte, int, int) {
	var data []byte
	return data, 0, 0
}
