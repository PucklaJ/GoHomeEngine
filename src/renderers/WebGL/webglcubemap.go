package renderer

import (
	"image"
	"image/color"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
)

type WebGLCubeMap struct {
	name      string
	oglName   *js.Object
	width     int
	height    int
	shadowMap bool
}

func (this *WebGLCubeMap) GetName() string {
	return this.name
}

func CreateWebGLCubeMap(name string) *WebGLCubeMap {
	cubeMap := &WebGLCubeMap{
		name: name,
	}
	return cubeMap
}

func (this *WebGLCubeMap) Load(data []byte, width, height int, shadowMap bool) {
	this.width = width
	this.height = height
	this.shadowMap = shadowMap

	this.oglName = gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)

	for i := 0; i < 6; i++ {
		if shadowMap {
			gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i, 0, gl.DEPTH_COMPONENT, width, height, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
		} else {
			gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
		}
	}

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
}

func (this *WebGLCubeMap) LoadFromImage(img image.Image) {
}

func (this *WebGLCubeMap) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + int(unit))
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
}
func (this *WebGLCubeMap) Unbind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + int(unit))
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, nil)

}
func (this *WebGLCubeMap) GetWidth() int {
	return int(this.width)
}
func (this *WebGLCubeMap) GetHeight() int {
	return int(this.height)
}
func (this *WebGLCubeMap) Terminate() {
	gl.DeleteTexture(this.oglName)
}
func (this *WebGLCubeMap) SetFiltering(filtering int) {
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
	var filter int
	if filtering == gohome.FILTERING_NEAREST {
		filter = gl.NEAREST
	} else if filtering == gohome.FILTERING_LINEAR {
		filter = gl.LINEAR
	} else {
		filter = gl.NEAREST
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, filter)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, nil)
}
func (this *WebGLCubeMap) SetWrapping(wrapping int) {
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
	var wrap int
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

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, nil)

}
func (this *WebGLCubeMap) SetBorderColor(col color.Color) {
	gohome.ErrorMgr.Warning("CubeMap", this.name, "SetBorderColor does not work in WebGL")
}
func (this *WebGLCubeMap) SetBorderDepth(depth float32) {
	gohome.ErrorMgr.Warning("CubeMap", this.name, "SetBorderDepth does not work in WebGL")
}

func (this *WebGLCubeMap) GetKeyColor() color.Color {
	return nil
}

func (this *WebGLCubeMap) GetModColor() color.Color {
	return nil
}

func (this *WebGLCubeMap) SetKeyColor(col color.Color) {

}

func (this *WebGLCubeMap) SetModColor(col color.Color) {

}

func (this *WebGLCubeMap) GetData() (data []byte, width int, height int) {
	gohome.ErrorMgr.Warning("CubeMap", this.name, "GetData does not work in WebGL")
	width = this.GetWidth()
	height = this.GetHeight()
	return
}
