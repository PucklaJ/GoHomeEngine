package gohome

import (
	"image"
	"image/color"
	"sync"
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
	GetData() ([]byte, int, int)
}

func TextureToImage(tex Texture, flipX, flipY bool) image.Image {
	var wg sync.WaitGroup
	data, width, height := tex.GetData()
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

func GetColorFromData(x, y int, data []byte, width int) color.Color {
	return &Color{
		R: data[(x+y*width)*4+0],
		G: data[(x+y*width)*4+1],
		B: data[(x+y*width)*4+2],
		A: data[(x+y*width)*4+3],
	}
}
