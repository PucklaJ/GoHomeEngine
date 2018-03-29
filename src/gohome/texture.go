package gohome

import (
	// "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"image/color"
	"log"
	"strconv"
	"unsafe"
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

type OpenGLTexture struct {
	width        int
	height       int
	name         string
	oglName      uint32
	multiSampled bool
}

func (ogltex *OpenGLTexture) bindingPoint() uint32 {
	if ogltex.multiSampled {
		return gl.TEXTURE_2D_MULTISAMPLE
	} else {
		return gl.TEXTURE_2D
	}
}

func CreateOpenGLTexture(name string, multiSampled bool) *OpenGLTexture {
	tex := &OpenGLTexture{
		name:         name,
		multiSampled: multiSampled,
	}
	return tex
}

func printOGLTexture2DError(ogltex *OpenGLTexture, data []byte, width, height int) {
	err := gl.GetError()
	if err != gl.NO_ERROR {
		var errString string

		if err == gl.INVALID_VALUE {
			if width < 0 {
				errString = "width is less than 0 "
			} else if width > gl.MAX_TEXTURE_SIZE {
				errString = "width is too large: " + strconv.Itoa(gl.MAX_TEXTURE_SIZE) + " "
			}
			if height < 0 {
				errString = "height is less than 0"
			} else if height > gl.MAX_TEXTURE_SIZE {
				errString = "height is too large: " + strconv.Itoa(gl.MAX_TEXTURE_SIZE)
			}
			if errString == "" {
				errString = "Invalid Value"
			}
		} else if err == gl.INVALID_ENUM {
			if ogltex.bindingPoint() != gl.TEXTURE_2D {
				errString = "target should be TEXTURE_2D"
			}
			if errString == "" {
				errString = "Invalid Enum"
			}

		} else if err == gl.INVALID_OPERATION {
			errString = "Invalid Operation"
		}

		log.Println("Error loading texture data:", err, errString)
	}
}

func (ogltex *OpenGLTexture) Load(data []byte, width, height int, shadowMap bool) error {
	ogltex.width = width
	ogltex.height = height

	gl.GenTextures(1, &ogltex.oglName)

	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)

	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_LOD_BIAS, -0.4)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_BASE_LEVEL, 0)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MAX_LEVEL, 10)
	if oglRenderer, ok := Render.(*OpenGLRenderer); ok {
		if oglRenderer.hasExtenstion("GL_EXT_texture_filter_anisotropic") {
			gl.TexParameterf(ogltex.bindingPoint(), GL_TEXTURE_MAX_ANISOTROPY, 4.0)
		}
	}

	if ogltex.multiSampled {
		gl.TexImage2DMultisample(ogltex.bindingPoint(), 8, gl.RGBA, int32(ogltex.width), int32(ogltex.height), true)
	} else {
		var ptr unsafe.Pointer
		if data == nil {
			ptr = gl.Ptr(nil)
		} else {
			ptr = gl.Ptr(&data[0])
		}
		gl.GetError()
		if shadowMap {
			gl.TexImage2D(ogltex.bindingPoint(), 0, gl.DEPTH_COMPONENT, int32(ogltex.width), int32(ogltex.height), 0, gl.DEPTH_COMPONENT, gl.FLOAT, ptr)
		} else {
			gl.TexImage2D(ogltex.bindingPoint(), 0, gl.RGBA, int32(ogltex.width), int32(ogltex.height), 0, gl.RGBA, gl.UNSIGNED_BYTE, ptr)
		}
		printOGLTexture2DError(ogltex, data, width, height)
	}
	gl.GenerateMipmap(ogltex.bindingPoint())

	gl.BindTexture(ogltex.bindingPoint(), 0)

	return nil

}

func toTextureUnit(unit uint32) uint32 {
	switch unit {
	case 0:
		return gl.TEXTURE0
	case 1:
		return gl.TEXTURE1
	case 2:
		return gl.TEXTURE2
	case 3:
		return gl.TEXTURE3
	case 4:
		return gl.TEXTURE4
	case 5:
		return gl.TEXTURE5
	case 6:
		return gl.TEXTURE6
	default:
		return gl.TEXTURE0 + unit
	}
}

func (ogltex *OpenGLTexture) Bind(unit uint32) {
	newUnit := toTextureUnit(unit)
	gl.ActiveTexture(newUnit)
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
}
func (ogltex *OpenGLTexture) Unbind(unit uint32) {
	newUnit := toTextureUnit(unit)
	gl.ActiveTexture(newUnit)
	gl.BindTexture(ogltex.bindingPoint(), 0)
}
func (ogltex *OpenGLTexture) GetWidth() int {
	return ogltex.width
}
func (ogltex *OpenGLTexture) GetHeight() int {
	return ogltex.height
}

func (ogltex *OpenGLTexture) Terminate() {
	gl.DeleteTextures(1, &ogltex.oglName)
}

func (ogltex *OpenGLTexture) SetFiltering(filtering uint32) {
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
	var filter int32
	if filtering == FILTERING_NEAREST {
		filter = gl.NEAREST
	} else if filtering == FILTERING_LINEAR {
		filter = gl.LINEAR
	} else {
		filter = gl.NEAREST
	}
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MAG_FILTER, filter)

	gl.BindTexture(ogltex.bindingPoint(), 0)
}

func (ogltex *OpenGLTexture) SetWrapping(wrapping uint32) {
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
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
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_S, wrap)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_T, wrap)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_R, wrap)

	gl.BindTexture(ogltex.bindingPoint(), 0)
}

func (ogltex *OpenGLTexture) SetBorderColor(col color.Color) {
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)

	borderColor := colorToVec4(col)

	gl.TexParameterfv(ogltex.bindingPoint(), gl.TEXTURE_BORDER_COLOR, &borderColor[0])

	gl.BindTexture(ogltex.bindingPoint(), 0)
}

func (ogltex *OpenGLTexture) SetBorderDepth(depth float32) {
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)

	col := [4]float32{depth, depth, depth, depth}

	gl.TexParameterfv(ogltex.bindingPoint(), gl.TEXTURE_BORDER_COLOR, &col[0])

	gl.BindTexture(ogltex.bindingPoint(), 0)
}

func (ogltex *OpenGLTexture) GetName() string {
	return ogltex.name
}
