package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles2"
	"image"
	"image/color"
)

var currentlyBoundRT *OpenGLES2RenderTexture
var screenFramebuffer int32

type OpenGLES2RenderTexture struct {
	Name         string
	fbo          uint32
	rbo          uint32
	depthBuffer  bool
	shadowMap    bool
	cubeMap      bool
	textures     []gohome.Texture
	prevViewport gohome.Viewport
	viewport     gohome.Viewport
	prevRT       *OpenGLES2RenderTexture
}

func CreateOpenGLES2RenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) *OpenGLES2RenderTexture {
	rt := &OpenGLES2RenderTexture{}

	rt.Create(name, width, height, textures, depthBuffer, multiSampled, shadowMap, cubeMap)

	return rt
}

func (this *OpenGLES2RenderTexture) loadTextures(width, height, textures uint32, cubeMap bool) {
	var i uint32
	for i = 0; i < textures; i++ {
		var ogltex *OpenGLES2Texture
		var oglcubemap *OpenGLES2CubeMap
		var texture gohome.Texture
		if cubeMap {
			oglcubemap = CreateOpenGLES2CubeMap(this.Name)
			texture = oglcubemap
		} else {
			ogltex = CreateOpenGLES2Texture(this.Name)
			texture = ogltex
		}
		texture.Load(nil, int(width), int(height), this.shadowMap)
		if cubeMap {
			gl.GetError()
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, oglcubemap.oglName)
			handleOpenGLError("RenderTexture", this.Name, "Binding cubemap")
		} else {
			gl.GetError()
			gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
			handleOpenGLError("RenderTexture", this.Name, "Binding texture 2d")
		}
		if this.shadowMap {
			if cubeMap {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
				handleOpenGLError("RenderTexture", this.Name, "glFramebufferTexture2D with depthBuffer and CubeMap")
			} else {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, ogltex.bindingPoint(), ogltex.oglName, 0)
				handleOpenGLError("RenderTexture", this.Name, "glFramebufferTexture2D with depthBuffer and TEXTURE2D")
			}
		} else {
			if cubeMap {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
				handleOpenGLError("RenderTexture", this.Name, "glFramebufferTexture2D with CubeMap")
			} else {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, ogltex.bindingPoint(), ogltex.oglName, 0)
				handleOpenGLError("RenderTexture", this.Name, "glFramebufferTexture2D with TEXTURE2D")
			}
		}
		if !cubeMap {
			texture.SetFiltering(gohome.FILTERING_LINEAR)
		}
		if cubeMap {
			gl.GetError()
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
			handleOpenGLError("RenderTexture", this.Name, "glBindTexture with CubeMap")
		} else {
			gl.GetError()
			gl.BindTexture(ogltex.bindingPoint(), 0)
			handleOpenGLError("RenderTexture", this.Name, "glBindTexture with TEXTURE2D")
		}
		this.textures = append(this.textures, texture)
	}
}

func (this *OpenGLES2RenderTexture) loadRenderBuffer(width, height uint32) {
	if this.depthBuffer {
		gl.GetError()
		var buf [1]uint32
		gl.GenRenderbuffers(1, buf[:])
		this.rbo = buf[0]
		handleOpenGLError("RenderTexture", this.Name, "glGenRenderbuffers")
		gl.BindRenderbuffer(gl.RENDERBUFFER, this.rbo)
		handleOpenGLError("RenderTexture", this.Name, "glBindRenderbuffer")
		gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT16, int32(width), int32(height))
		handleOpenGLError("RenderTexture", this.Name, "glRenderbufferStorage")
		gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, this.rbo)
		handleOpenGLError("RenderTexture", this.Name, "glFramebufferRenderbuffer")
		gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
		handleOpenGLError("RenderTexture", this.Name, "glBindRenderbuffer with 0")
	}
}

func (this *OpenGLES2RenderTexture) Create(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) {
	if textures == 0 {
		textures = 1
	}

	this.Name = name
	this.shadowMap = shadowMap
	this.depthBuffer = depthBuffer && !shadowMap
	this.cubeMap = cubeMap

	gl.GetError()
	var buf [1]uint32
	gl.GenFramebuffers(1, buf[:])
	this.fbo = buf[0]
	handleOpenGLError("RenderTexture", this.Name, "glGenFramebuffers")

	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer")

	this.loadRenderBuffer(width, height)
	this.loadTextures(width, height, textures, cubeMap)
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		handleOpenGLError("RenderTexture", this.Name, "glCheckFramebufferStatus")
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "RenderTexture", this.Name, "Framebuffer is not complete")
		return
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer with 0")

	this.viewport = gohome.Viewport{
		0,
		0, 0,
		int(width), int(height),
		false,
	}

	this.SetAsTarget()
	gohome.Render.ClearScreen(gohome.Color{0, 0, 0, 0})
	this.UnsetAsTarget()
}

func (this *OpenGLES2RenderTexture) Load(data []byte, width, height int, shadowMap bool) error {
	return &OpenGLError{errorString: "The Load method of RenderTexture is not used!"}
}

func (ogltex *OpenGLES2RenderTexture) LoadFromImage(img image.Image) error {
	return &OpenGLError{errorString: "The LoadFromImage method of RenderTexture is not used!"}
}

func (this *OpenGLES2RenderTexture) GetName() string {
	return this.Name
}

func (this *OpenGLES2RenderTexture) SetAsTarget() {
	if currentlyBoundRT == nil {
		var buf [1]int32
		gl.GetIntegerv(gl.FRAMEBUFFER_BINDING, buf[:])
		screenFramebuffer = buf[0]
	} else {
	}
	this.prevRT = currentlyBoundRT
	currentlyBoundRT = this
	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer in SetAsTarget")
	this.prevViewport = gohome.Render.GetViewport()
	gohome.Render.SetViewport(this.viewport)
}

func (this *OpenGLES2RenderTexture) UnsetAsTarget() {
	if this.prevRT != nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, this.prevRT.fbo)
		currentlyBoundRT = this.prevRT
	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(screenFramebuffer))
		currentlyBoundRT = nil
	}
	handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer in UnsetAsTarget")
	gohome.Render.SetViewport(this.prevViewport)
}

func (this *OpenGLES2RenderTexture) Blit(rtex gohome.RenderTexture) {
	gohome.ErrorMgr.Error("RenderTexture", this.Name, "BlitFramebuffer does not work in OpenGLES2 2.0")
}

func (this *OpenGLES2RenderTexture) Bind(unit uint32) {
	this.BindIndex(0, unit)
}

func (this *OpenGLES2RenderTexture) Unbind(unit uint32) {
	this.UnbindIndex(0, unit)
}

func (this *OpenGLES2RenderTexture) BindIndex(index, unit uint32) {
	if index < uint32(len(this.textures)) {
		this.textures[index].Bind(unit)
	}
}

func (this *OpenGLES2RenderTexture) UnbindIndex(index, unit uint32) {
	if index < uint32(len(this.textures)) {
		this.textures[index].Unbind(unit)
	}
}

func (this *OpenGLES2RenderTexture) GetWidth() int {
	if len(this.textures) == 0 {
		return 0
	} else {
		return this.textures[0].GetWidth()
	}
}

func (this *OpenGLES2RenderTexture) GetHeight() int {
	if len(this.textures) == 0 {
		return 0
	} else {
		return this.textures[0].GetHeight()
	}
}

func (this *OpenGLES2RenderTexture) Terminate() {
	var fbuf [1]uint32
	fbuf[0] = this.fbo
	defer gl.DeleteFramebuffers(1, fbuf[:])
	if this.depthBuffer {
		var rbuf [1]uint32
		rbuf[0] = this.rbo
		defer gl.DeleteRenderbuffers(1, rbuf[:])
	}
	for i := 0; i < len(this.textures); i++ {
		defer this.textures[i].Terminate()
	}
	this.textures = append(this.textures[:0], this.textures[len(this.textures):]...)
}

func (this *OpenGLES2RenderTexture) ChangeSize(width, height uint32) {
	if uint32(this.GetWidth()) != width || uint32(this.GetHeight()) != height {
		textures := uint32(len(this.textures))
		this.Terminate()
		this.Create(this.Name, width, height, textures, this.depthBuffer, false, this.shadowMap, this.cubeMap)
	}
}

func (this *OpenGLES2RenderTexture) SetFiltering(filtering uint32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetFiltering(filtering)
	}
}

func (this *OpenGLES2RenderTexture) SetWrapping(wrapping uint32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetWrapping(wrapping)
	}
}

func (this *OpenGLES2RenderTexture) SetBorderColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderColor(col)
	}
}

func (this *OpenGLES2RenderTexture) SetBorderDepth(depth float32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderDepth(depth)
	}
}

func (this *OpenGLES2RenderTexture) GetKeyColor() color.Color {
	if len(this.textures) == 0 {
		return nil
	}
	return this.textures[0].GetKeyColor()
}

func (this *OpenGLES2RenderTexture) GetModColor() color.Color {
	if len(this.textures) == 0 {
		return nil
	}
	return this.textures[0].GetModColor()
}

func (this *OpenGLES2RenderTexture) SetKeyColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetKeyColor(col)
	}
}

func (this *OpenGLES2RenderTexture) SetModColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetModColor(col)
	}
}

func (this *OpenGLES2RenderTexture) GetData() ([]byte, int, int) {
	if len(this.textures) == 0 {
		return nil, 0, 0
	}
	if _, ok := this.textures[0].(*OpenGLES2Texture); ok {
		if gohome.Render.HasFunctionAvailable("BLIT_FRAMEBUFFER") {
			rtex := CreateOpenGLES2RenderTexture("Temp", uint32(this.GetWidth()), uint32(this.GetHeight()), 1, false, false, false, false)
			this.Blit(rtex)
			data, width, height := rtex.GetData()
			rtex.Terminate()
			return data, width, height
		}
	}

	return this.textures[0].GetData()
}
