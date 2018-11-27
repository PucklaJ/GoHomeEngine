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
