package gohome

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"image/color"
)

type CubeMap interface {
	Load(data []byte, width, height int, shadowMap bool) error
	Bind(unit uint32)
	Unbind(unit uint32)
	GetWidth() int
	GetHeight() int
	Terminate()
	SetFiltering(filtering uint32)
	SetWrapping(wrapping uint32)
	SetBorderColor(col color.Color)
	SetBorderDepth(depth float32)
	GetName() string
}

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

func (this *OpenGLCubeMap) Load(data []byte, width, height int, shadowMap bool) error {
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

	return nil
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
	if filtering == FILTERING_NEAREST {
		filter = gl.NEAREST
	} else if filtering == FILTERING_LINEAR {
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
	if wrapping == WRAPPING_REPEAT {
		wrap = gl.REPEAT
	} else if wrapping == WRAPPING_CLAMP_TO_BORDER {
		wrap = gl.CLAMP_TO_BORDER
	} else if wrapping == WRAPPING_CLAMP_TO_EDGE {
		wrap = gl.CLAMP_TO_EDGE
	} else if wrapping == WRAPPING_MIRRORED_REPEAT {
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

	borderColor := colorToVec4(col)

	gl.TexParameterfv(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BORDER_COLOR, &borderColor[0])

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
}
func (this *OpenGLCubeMap) SetBorderDepth(depth float32) {
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)

	col := [4]float32{depth, depth, depth, depth}

	gl.TexParameterfv(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BORDER_COLOR, &col[0])

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
}
