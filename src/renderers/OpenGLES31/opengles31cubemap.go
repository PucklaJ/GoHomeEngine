package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles31"
	"image"
	"image/color"
)

type OpenGLES31CubeMap struct {
	name      string
	oglName   uint32
	width     uint32
	height    uint32
	shadowMap bool
}

func (this *OpenGLES31CubeMap) GetName() string {
	return this.name
}

func CreateOpenGLES31CubeMap(name string) *OpenGLES31CubeMap {
	cubeMap := &OpenGLES31CubeMap{
		name: name,
	}
	return cubeMap
}

func (this *OpenGLES31CubeMap) Load(data []byte, width, height int, shadowMap bool) {
	this.width = uint32(width)
	this.height = uint32(height)
	this.shadowMap = shadowMap

	var tex [1]uint32
	gl.GenTextures(1, tex[:])
	this.oglName = tex[0]
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

func (this *OpenGLES31CubeMap) LoadFromImage(img image.Image) {
}

func (this *OpenGLES31CubeMap) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
}
func (this *OpenGLES31CubeMap) Unbind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)

}
func (this *OpenGLES31CubeMap) GetWidth() int {
	return int(this.width)
}
func (this *OpenGLES31CubeMap) GetHeight() int {
	return int(this.height)
}
func (this *OpenGLES31CubeMap) Terminate() {
	var tex [1]uint32
	tex[0] = this.oglName
	gl.DeleteTextures(1, tex[:])
}
func (this *OpenGLES31CubeMap) SetFiltering(filtering uint32) {
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
func (this *OpenGLES31CubeMap) SetWrapping(wrapping uint32) {
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
	var wrap int32
	if wrapping == gohome.WRAPPING_REPEAT {
		wrap = gl.REPEAT
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
func (this *OpenGLES31CubeMap) SetBorderColor(col color.Color) {
	gohome.ErrorMgr.Warning("CubeMap", this.name, "SetBorderColor does not work in OpenGLES 3.1")
}
func (this *OpenGLES31CubeMap) SetBorderDepth(depth float32) {
	gohome.ErrorMgr.Warning("CubeMap", this.name, "SetBorderDepth does not work in OpenGLES 3.1")
}

func (this *OpenGLES31CubeMap) GetKeyColor() color.Color {
	return nil
}

func (this *OpenGLES31CubeMap) GetModColor() color.Color {
	return nil
}

func (this *OpenGLES31CubeMap) SetKeyColor(col color.Color) {

}

func (this *OpenGLES31CubeMap) SetModColor(col color.Color) {

}

func (this *OpenGLES31CubeMap) GetData() (data []byte, width int, height int) {
	width = this.GetWidth()
	height = this.GetHeight()
	gohome.ErrorMgr.Warning("CubeMap", this.name, "GetData does not work in OpenGLES 3.1")
	return
}
