package gohome

import (
	"image"
	"image/color"
)

const (
	FILTERING_NEAREST                uint32 = iota
	FILTERING_LINEAR                 uint32 = iota
	FILTERING_NEAREST_MIPMAP_NEAREST uint32 = iota
	FILTERING_LINEAR_MIPMAP_LINEAR   uint32 = iota

	WRAPPING_REPEAT          uint32 = iota
	WRAPPING_CLAMP_TO_BORDER uint32 = iota
	WRAPPING_CLAMP_TO_EDGE   uint32 = iota
	WRAPPING_MIRRORED_REPEAT uint32 = iota
)

type Texture interface {
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
