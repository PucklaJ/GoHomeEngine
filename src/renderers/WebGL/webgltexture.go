package renderer

import (
	"image"
	"image/color"
	"strconv"
	"sync"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
)

type WebGLTexture struct {
	width   int
	height  int
	name    string
	oglName *js.Object

	keyColor color.Color
	modColor color.Color
}

func (ogltex *WebGLTexture) bindingPoint() int {
	return gl.TEXTURE_2D
}

func CreateWebGLTexture(name string) *WebGLTexture {
	tex := &WebGLTexture{
		name: name,
	}
	return tex
}

func printOGLTexture2DError(ogltex *WebGLTexture, data []byte, width, height int) {
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

		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Texture", ogltex.GetName(), "Couldn't load data: "+strconv.Itoa(int(err))+": "+errString)
	}
}

func isPowerOfTwo(width, height int) bool {
	var pow = 1
	var widthis, heightis bool
	for i := 0; i < 11; i++ {
		if !widthis && width == pow {
			widthis = true
		}
		if !heightis && height == pow {
			heightis = true
		}
		pow *= 2
	}

	return widthis && heightis
}

func (ogltex *WebGLTexture) Load(data []byte, width, height int, shadowMap bool) {
	ogltex.width = width
	ogltex.height = height

	ogltex.oglName = gl.CreateTexture()

	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)

	if isPowerOfTwo(width, height) {
		gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_T, gl.REPEAT)
	} else {
		gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	var dataobj *js.Object
	if data == nil {
		dataobj = nil
	} else {
		dataobj = js.InternalObject(data).Get("$array")
	}

	gl.GetError()
	gl.TexImage2D(ogltex.bindingPoint(), 0, gl.RGBA, ogltex.width, ogltex.height, 0, gl.RGBA, gl.UNSIGNED_BYTE, dataobj)
	printOGLTexture2DError(ogltex, data, width, height)

	gl.BindTexture(ogltex.bindingPoint(), nil)
}

func loadImageData(img_data *[]byte, img image.Image, start_width, end_width, max_width, max_height int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	var r, g, b, a uint32
	var color color.Color
	for x := start_width; x < max_width && x < end_width; x++ {
		for y := 0; y < max_height; y++ {
			color = img.At(int(x), int(y))
			r, g, b, a = color.RGBA()
			(*img_data)[(x+y*max_width)*4+0] = byte(float64(r) / float64(0xffff) * float64(255.0))
			(*img_data)[(x+y*max_width)*4+1] = byte(float64(g) / float64(0xffff) * float64(255.0))
			(*img_data)[(x+y*max_width)*4+2] = byte(float64(b) / float64(0xffff) * float64(255.0))
			(*img_data)[(x+y*max_width)*4+3] = byte(float64(a) / float64(0xffff) * float64(255.0))
		}
	}
}

func (ogltex *WebGLTexture) LoadFromImage(img image.Image) {

	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	img_data := make([]byte, width*height*4)

	var wg1 sync.WaitGroup
	var i float32
	deltaWidth := float32(width) / float32(gohome.NUM_GO_ROUTINES_TEXTURE_LOADING)
	wg1.Add(int(gohome.NUM_GO_ROUTINES_TEXTURE_LOADING + 1))
	for i = 0; i <= float32(gohome.NUM_GO_ROUTINES_TEXTURE_LOADING); i++ {
		go loadImageData(&img_data, img, int(i*deltaWidth), int((i+1)*deltaWidth), width, height, &wg1)
	}
	wg1.Wait()

	ogltex.Load(img_data, width, height, false)
}

func toTextureUnit(unit uint32) int {
	return gl.TEXTURE0 + int(unit)
}

func (ogltex *WebGLTexture) Bind(unit uint32) {
	newUnit := toTextureUnit(unit)
	gl.GetError()
	gl.ActiveTexture(newUnit)
	handleWebGLError("Texture", ogltex.name, "glActiveTexture in Bind with "+strconv.Itoa(int(unit)))
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
	handleWebGLError("Texture", ogltex.name, "glBindTexture in Bind")
}
func (ogltex *WebGLTexture) Unbind(unit uint32) {
	newUnit := toTextureUnit(unit)
	gl.GetError()
	gl.ActiveTexture(newUnit)
	handleWebGLError("Texture", ogltex.name, "glActiveTexture in Unbind with "+strconv.Itoa(int(unit)))
	gl.BindTexture(ogltex.bindingPoint(), nil)
	handleWebGLError("Texture", ogltex.name, "glBindTexture in Unbind")
}
func (ogltex *WebGLTexture) GetWidth() int {
	return ogltex.width
}
func (ogltex *WebGLTexture) GetHeight() int {
	return ogltex.height
}

func (ogltex *WebGLTexture) Terminate() {
	gl.DeleteTexture(ogltex.oglName)
}

func (ogltex *WebGLTexture) SetFiltering(filtering int) {
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
	var filter int
	if filtering == gohome.FILTERING_NEAREST {
		filter = gl.NEAREST
	} else if filtering == gohome.FILTERING_LINEAR {
		filter = gl.LINEAR
	} else {
		filter = gl.NEAREST
	}
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MAG_FILTER, filter)

	gl.BindTexture(ogltex.bindingPoint(), nil)
}

func (ogltex *WebGLTexture) SetWrapping(wrapping int) {
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
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
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_S, wrap)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_T, wrap)

	gl.BindTexture(ogltex.bindingPoint(), nil)
}

func (ogltex *WebGLTexture) SetBorderColor(col color.Color) {
	gohome.ErrorMgr.Warning("Texture", ogltex.name, "SetBorderColor does not work in WebGL")
}

func (ogltex *WebGLTexture) SetBorderDepth(depth float32) {
	gohome.ErrorMgr.Warning("Texture", ogltex.name, "SetBorderDepth does not work in WebGL")
}

func (ogltex *WebGLTexture) SetKeyColor(col color.Color) {
	ogltex.keyColor = col
}

func (ogltex *WebGLTexture) GetKeyColor() color.Color {
	return ogltex.keyColor
}

func (ogltex *WebGLTexture) SetModColor(col color.Color) {
	ogltex.modColor = col
}

func (ogltex *WebGLTexture) GetModColor() color.Color {
	return ogltex.modColor
}

func (ogltex *WebGLTexture) GetName() string {
	return ogltex.name
}

func (ogltex *WebGLTexture) GetData() (data []byte, width int, height int) {
	width = ogltex.GetWidth()
	height = ogltex.GetHeight()
	gohome.ErrorMgr.Warning("Texture", ogltex.name, "GetData does not work in WebGL")
	return
}
