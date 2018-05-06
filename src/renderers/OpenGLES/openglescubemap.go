package renderer

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/mobile/gl"
	"image"
	"image/color"
	"log"
)

type OpenGLESCubeMap struct {
	name      string
	oglName   gl.Texture
	width     uint32
	height    uint32
	shadowMap bool
	gles      *gl.Context
}

func (this *OpenGLESCubeMap) GetName() string {
	return this.name
}

func CreateOpenGLESCubeMap(name string) *OpenGLESCubeMap {
	cubeMap := &OpenGLESCubeMap{
		name: name,
	}
	render, _ := gohome.Render.(*OpenGLESRenderer)
	cubeMap.gles = &render.gles
	return cubeMap
}

func (this *OpenGLESCubeMap) Load(data []byte, width, height int, shadowMap bool) error {
	this.width = uint32(width)
	this.height = uint32(height)
	this.shadowMap = shadowMap

	this.oglName = (*this.gles).CreateTexture()
	(*this.gles).BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)

	var i uint32
	for i = 0; i < 6; i++ {
		if shadowMap {
			(*this.gles).TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+gl.Enum(i), 0, width, height, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
		} else {
			(*this.gles).TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+gl.Enum(i), 0, width, height, gl.RGBA, gl.UNSIGNED_BYTE, nil)
		}
	}

	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	return nil
}

func (this *OpenGLESCubeMap) LoadFromImage(img image.Image) error {
	return nil
}

func (this *OpenGLESCubeMap) Bind(unit uint32) {
	(*this.gles).ActiveTexture(gl.Enum(gl.TEXTURE0 + unit))
	(*this.gles).BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
}
func (this *OpenGLESCubeMap) Unbind(unit uint32) {
	(*this.gles).ActiveTexture(gl.Enum(gl.TEXTURE0 + unit))
	(*this.gles).BindTexture(gl.TEXTURE_CUBE_MAP, gl.Texture{0})

}
func (this *OpenGLESCubeMap) GetWidth() int {
	return int(this.width)
}
func (this *OpenGLESCubeMap) GetHeight() int {
	return int(this.height)
}
func (this *OpenGLESCubeMap) Terminate() {
	(*this.gles).DeleteTexture(this.oglName)
}
func (this *OpenGLESCubeMap) SetFiltering(filtering uint32) {
	(*this.gles).BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
	var filter int
	if filtering == gohome.FILTERING_NEAREST {
		filter = gl.NEAREST
	} else if filtering == gohome.FILTERING_LINEAR {
		filter = gl.LINEAR
	} else {
		filter = gl.NEAREST
	}
	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, filter)
	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, filter)

	(*this.gles).BindTexture(gl.TEXTURE_CUBE_MAP, gl.Texture{0})
}
func (this *OpenGLESCubeMap) SetWrapping(wrapping uint32) {
	(*this.gles).BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)
	defer (*this.gles).BindTexture(gl.TEXTURE_CUBE_MAP, gl.Texture{0})
	var wrap int
	if wrapping == gohome.WRAPPING_REPEAT {
		wrap = gl.REPEAT
	} else if wrapping == gohome.WRAPPING_CLAMP_TO_BORDER {
		log.Println("CLAMP_TO_BORDER is not supported by OpenGLES")
		return
	} else if wrapping == gohome.WRAPPING_CLAMP_TO_EDGE {
		wrap = gl.CLAMP_TO_EDGE
	} else if wrapping == gohome.WRAPPING_MIRRORED_REPEAT {
		wrap = gl.MIRRORED_REPEAT
	} else {
		wrap = gl.CLAMP_TO_EDGE
	}
	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, wrap)
	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, wrap)
	(*this.gles).TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, wrap)
}
func (this *OpenGLESCubeMap) SetBorderColor(col color.Color) {
	// gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)

	// borderColor := gohome.ColorToVec4(col)

	// gl.TexParameterfv(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BORDER_COLOR, &borderColor[0])

	// gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	log.Println("Border Color is not supported by OpenGLES")
}
func (this *OpenGLESCubeMap) SetBorderDepth(depth float32) {
	// gl.BindTexture(gl.TEXTURE_CUBE_MAP, this.oglName)

	// col := [4]float32{depth, depth, depth, depth}

	// gl.TexParameterfv(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BORDER_COLOR, &col[0])

	// gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	log.Println("Border Depth is not supported by OpenGLES")
}
