package renderer

import (
	// "fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/mobile/gl"
	"image/color"
	"log"
	"strconv"
)

type OpenGLESTexture struct {
	width   int
	height  int
	name    string
	oglName gl.Texture
	gles    *gl.Context
}

func (ogltex *OpenGLESTexture) bindingPoint() gl.Enum {
	return gl.TEXTURE_2D
}

func CreateOpenGLESTexture(name string) *OpenGLESTexture {
	tex := &OpenGLESTexture{
		name: name,
	}
	render, _ := gohome.Render.(*OpenGLESRenderer)
	tex.gles = &render.gles
	return tex
}

func printOGLESTexture2DError(ogltex *OpenGLESTexture, data []byte, width, height int) {
	err := (*ogltex.gles).GetError()
	if err != gl.NO_ERROR {
		var errString string

		if err == gl.INVALID_VALUE {
			if width < 0 {
				errString = "width is less than 0 "
			} else if width > (*ogltex.gles).GetInteger(gl.MAX_TEXTURE_SIZE) {
				errString = "width is too large: " + strconv.Itoa((*ogltex.gles).GetInteger(gl.MAX_TEXTURE_SIZE)) + " "
			}
			if height < 0 {
				errString = "height is less than 0"
			} else if height > (*ogltex.gles).GetInteger(gl.MAX_TEXTURE_SIZE) {
				errString = "height is too large: " + strconv.Itoa((*ogltex.gles).GetInteger(gl.MAX_TEXTURE_SIZE))
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

		log.Println("Error loading texture data of", ogltex.name, ":", err, errString)
	}
}

func (ogltex *OpenGLESTexture) Load(data []byte, width, height int, shadowMap bool) error {
	ogltex.width = width
	ogltex.height = height

	ogltex.oglName = (*ogltex.gles).CreateTexture()

	(*ogltex.gles).BindTexture(ogltex.bindingPoint(), ogltex.oglName)

	(*ogltex.gles).TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	(*ogltex.gles).TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	(*ogltex.gles).TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	(*ogltex.gles).TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	(*ogltex.gles).TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	(*ogltex.gles).TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_MIN_LOD, -0.4)
	(*ogltex.gles).TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_MAX_LOD, -0.4)
	(*ogltex.gles).TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_BASE_LEVEL, 0)
	(*ogltex.gles).TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MAX_LEVEL, 10)

	(*ogltex.gles).GetError()
	if shadowMap {
		(*ogltex.gles).TexImage2D(ogltex.bindingPoint(), 0, ogltex.width, ogltex.height, gl.RGBA, gl.UNSIGNED_BYTE, data)
	} else {
		(*ogltex.gles).TexImage2D(ogltex.bindingPoint(), 0, ogltex.width, ogltex.height, gl.RGBA, gl.UNSIGNED_BYTE, data)
	}
	printOGLESTexture2DError(ogltex, data, width, height)
	(*ogltex.gles).GenerateMipmap(ogltex.bindingPoint())

	(*ogltex.gles).BindTexture(ogltex.bindingPoint(), gl.Texture{0})

	return nil

}

func (ogltex *OpenGLESTexture) Bind(unit uint32) {
	(*ogltex.gles).ActiveTexture(gl.Enum(gl.TEXTURE0 + unit))
	(*ogltex.gles).BindTexture(ogltex.bindingPoint(), ogltex.oglName)
}
func (ogltex *OpenGLESTexture) Unbind(unit uint32) {
	(*ogltex.gles).ActiveTexture(gl.Enum(gl.TEXTURE0 + unit))
	(*ogltex.gles).BindTexture(ogltex.bindingPoint(), gl.Texture{0})
}
func (ogltex *OpenGLESTexture) GetWidth() int {
	return ogltex.width
}
func (ogltex *OpenGLESTexture) GetHeight() int {
	return ogltex.height
}

func (ogltex *OpenGLESTexture) Terminate() {
	(*ogltex.gles).DeleteTexture(ogltex.oglName)
}

func (ogltex *OpenGLESTexture) SetFiltering(filtering uint32) {
	(*ogltex.gles).BindTexture(ogltex.bindingPoint(), ogltex.oglName)
	var filter int
	if filtering == gohome.FILTERING_NEAREST {
		filter = gl.NEAREST
	} else if filtering == gohome.FILTERING_LINEAR {
		filter = gl.LINEAR
	} else {
		filter = gl.NEAREST
	}
	(*ogltex.gles).TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, filter)
	(*ogltex.gles).TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MAG_FILTER, filter)

	(*ogltex.gles).BindTexture(ogltex.bindingPoint(), gl.Texture{0})
}

func (ogltex *OpenGLESTexture) SetWrapping(wrapping uint32) {
	(*ogltex.gles).BindTexture(ogltex.bindingPoint(), ogltex.oglName)
	defer (*ogltex.gles).BindTexture(ogltex.bindingPoint(), gl.Texture{0})
	var wrap int
	if wrapping == gohome.WRAPPING_REPEAT {
		wrap = gl.REPEAT
	} else if wrapping == gohome.WRAPPING_CLAMP_TO_BORDER {
		log.Println("CLAMP_TO_BORDER is not supported by OpenGLES")
	} else if wrapping == gohome.WRAPPING_CLAMP_TO_EDGE {
		wrap = gl.CLAMP_TO_EDGE
	} else if wrapping == gohome.WRAPPING_MIRRORED_REPEAT {
		wrap = gl.MIRRORED_REPEAT
	} else {
		wrap = gl.CLAMP_TO_EDGE
	}
	(*ogltex.gles).TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_S, wrap)
	(*ogltex.gles).TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_T, wrap)
	(*ogltex.gles).TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_R, wrap)

}

func (ogltex *OpenGLESTexture) SetBorderColor(col color.Color) {
	// (*ogltex.gles).BindTexture(ogltex.bindingPoint(), ogltex.oglName)

	// borderColor := gohome.ColorToVec4(col)

	// (*ogltex.gles).TexParameterfv(ogltex.bindingPoint(), gl, &borderColor[0])

	// gl.BindTexture(ogltex.bindingPoint(), 0)
	log.Println("Border color is not supported by OpenGLES")
}

func (ogltex *OpenGLESTexture) SetBorderDepth(depth float32) {
	// gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)

	// col := [4]float32{depth, depth, depth, depth}

	// gl.TexParameterfv(ogltex.bindingPoint(), gl.TEXTURE_BORDER_COLOR, &col[0])

	// gl.BindTexture(ogltex.bindingPoint(), 0)
	log.Println("Border depth is not supported by OpenGLES")
}

func (ogltex *OpenGLESTexture) GetName() string {
	return ogltex.name
}
