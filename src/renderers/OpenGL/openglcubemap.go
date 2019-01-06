package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/go-gl/gl/all-core/gl"
	"image"
	"image/color"
	"unsafe"
)

type OpenGLCubeMap struct {
	name      string
	oglName   uint32
	width     uint32
	height    uint32
	shadowMap bool
}

func (this *OpenGLCubeMap) GetName() string {
	return this.name
}

func CreateOpenGLCubeMap(name string) *OpenGLCubeMap {
	cubeMap := &OpenGLCubeMap{
		name: name,
	}
	return cubeMap
}

func (this *OpenGLCubeMap) Load(data []byte, width, height int, shadowMap bool) {
	this.width = uint32(width)
	this.height = uint32(height)
	this.shadowMap = shadowMap

	gl.GenTextures(1, &this.oglName)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)

	var i uint32
	for i = 0; i < 6; i++ {
		if shadowMap {
			gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i, 0, gl.DEPTH_COMPONENT, int32(width), int32(height), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
		} else {
			gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
		}
	}

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
}

func (this *OpenGLCubeMap) LoadFromImage(img image.Image) {
}

func (this *OpenGLCubeMap) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
}
func (this *OpenGLCubeMap) Unbind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)

}
func (this *OpenGLCubeMap) GetWidth() int {
	return int(this.width)
}
func (this *OpenGLCubeMap) GetHeight() int {
	return int(this.height)
}
func (this *OpenGLCubeMap) Terminate() {
	gl.DeleteTextures(1, &this.oglName)
}
func (this *OpenGLCubeMap) SetFiltering(filtering uint32) {
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
	var filter int32
	if filtering == gohome.FILTERING_NEAREST {
		filter = gl.NEAREST
	} else if filtering == gohome.FILTERING_LINEAR {
		filter = gl.LINEAR
	} else {
		filter = gl.NEAREST
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, filter)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
}
func (this *OpenGLCubeMap) SetWrapping(wrapping uint32) {
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
	var wrap int32
	if wrapping == gohome.WRAPPING_REPEAT {
		wrap = gl.REPEAT
	} else if wrapping == gohome.WRAPPING_CLAMP_TO_BORDER {
		wrap = gl.CLAMP_TO_BORDER
	} else if wrapping == gohome.WRAPPING_CLAMP_TO_EDGE {
		wrap = gl.CLAMP_TO_EDGE
	} else if wrapping == gohome.WRAPPING_MIRRORED_REPEAT {
		wrap = gl.MIRRORED_REPEAT
	} else {
		wrap = gl.CLAMP_TO_EDGE
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, wrap)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, wrap)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, wrap)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)

}
func (this *OpenGLCubeMap) SetBorderColor(col color.Color) {
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)

	borderColor := gohome.ColorToVec4(col)

	gl.TexParameterfv(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BORDER_COLOR, &borderColor[0])

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
}
func (this *OpenGLCubeMap) SetBorderDepth(depth float32) {
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)

	col := [4]float32{depth, depth, depth, depth}

	gl.TexParameterfv(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BORDER_COLOR, &col[0])

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
}

func (this *OpenGLCubeMap) GetKeyColor() color.Color {
	return nil
}

func (this *OpenGLCubeMap) GetModColor() color.Color {
	return nil
}

func (this *OpenGLCubeMap) SetKeyColor(col color.Color) {

}

func (this *OpenGLCubeMap) SetModColor(col color.Color) {

}

func (this *OpenGLCubeMap) GetData() (data []byte, width int, height int) {
	width = this.GetWidth()
	height = this.GetHeight()
	data = make([]byte, width*height*4)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP_POSITIVE_X, this.oglName)
	gl.GetTexImage(gl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&data[0]))
	gl.BindTexture(gl.TEXTURE_CUBE_MAP_POSITIVE_X, 0)
	return
}
