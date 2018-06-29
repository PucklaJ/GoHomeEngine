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
}
