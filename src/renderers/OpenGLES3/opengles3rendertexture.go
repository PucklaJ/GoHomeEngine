package renderer

import (
	"image"
	"image/color"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	gl "github.com/PucklaJ/android-go/gles3"
)

var currentlyBoundRT *OpenGLES3RenderTexture
var screenFramebuffer int32

type OpenGLES3RenderTexture struct {
	Name         string
	fbo          uint32
	rbo          uint32
	depthBuffer  bool
	shadowMap    bool
	cubeMap      bool
	textures     []gohome.Texture
	prevViewport gohome.Viewport
	viewport     gohome.Viewport
	prevRT       *OpenGLES3RenderTexture
}

func CreateOpenGLES3RenderTexture(name string, width, height, textures int, depthBuffer, shadowMap, cubeMap bool) *OpenGLES3RenderTexture {
	rt := &OpenGLES3RenderTexture{}

	rt.Create(name, width, height, textures, depthBuffer, false, shadowMap, cubeMap)

	return rt
}

func (this *OpenGLES3RenderTexture) loadTextures(width, height, textures int, cubeMap bool) {
	for i := 0; i < textures; i++ {
		var ogltex *OpenGLES3Texture
		var oglcubemap *OpenGLES3CubeMap
		var texture gohome.Texture
		if cubeMap {
			oglcubemap = CreateOpenGLES3CubeMap(this.Name)
			texture = oglcubemap
		} else {
			ogltex = CreateOpenGLES3Texture(this.Name)
			texture = ogltex
		}
		texture.Load(nil, int(width), int(height), this.shadowMap)
		if cubeMap {
			gl.GetError()
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, oglcubemap.oglName)
			handleOpenGLES3Error("RenderTexture", this.Name, "Binding cubemap")
		} else {
			gl.GetError()
			gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
			handleOpenGLES3Error("RenderTexture", this.Name, "Binding texture 2d")
		}
		if this.shadowMap {
			if cubeMap {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
				handleOpenGLES3Error("RenderTexture", this.Name, "glFramebufferTexture2D with depthBuffer and CubeMap")
			} else {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, ogltex.bindingPoint(), ogltex.oglName, 0)
				handleOpenGLES3Error("RenderTexture", this.Name, "glFramebufferTexture2D with depthBuffer and TEXTURE2D")
			}
		} else {
			if cubeMap {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+uint32(i), gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
				handleOpenGLES3Error("RenderTexture", this.Name, "glFramebufferTexture2D with CubeMap")
			} else {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+uint32(i), ogltex.bindingPoint(), ogltex.oglName, 0)
				handleOpenGLES3Error("RenderTexture", this.Name, "glFramebufferTexture2D with TEXTURE2D")
			}
		}
		if !cubeMap {
			texture.SetFiltering(gohome.FILTERING_LINEAR)
		}
		if cubeMap {
			gl.GetError()
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
			handleOpenGLES3Error("RenderTexture", this.Name, "glBindTexture with CubeMap")
		} else {
			gl.GetError()
			gl.BindTexture(ogltex.bindingPoint(), 0)
			handleOpenGLES3Error("RenderTexture", this.Name, "glBindTexture with TEXTURE2D")
		}
		this.textures = append(this.textures, texture)
	}
}

func (this *OpenGLES3RenderTexture) loadRenderBuffer(width, height int) {
	if this.depthBuffer {
		gl.GetError()
		var buf [1]uint32
		gl.GenRenderbuffers(1, buf[:])
		this.rbo = buf[0]
		handleOpenGLES3Error("RenderTexture", this.Name, "glGenRenderbuffers")
		gl.BindRenderbuffer(gl.RENDERBUFFER, this.rbo)
		handleOpenGLES3Error("RenderTexture", this.Name, "glBindRenderbuffer")
		gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(width), int32(height))
		handleOpenGLES3Error("RenderTexture", this.Name, "glRenderbufferStorage")
		gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, this.rbo)
		handleOpenGLES3Error("RenderTexture", this.Name, "glFramebufferRenderbuffer")
		gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
		handleOpenGLES3Error("RenderTexture", this.Name, "glBindRenderbuffer with 0")
	}
}

func (this *OpenGLES3RenderTexture) Create(name string, width, height, textures int, depthBuffer, multiSampled, shadowMap, cubeMap bool) {
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
	handleOpenGLES3Error("RenderTexture", this.Name, "glGenFramebuffers")

	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	handleOpenGLES3Error("RenderTexture", this.Name, "glBindFramebuffer")

	this.loadRenderBuffer(width, height)
	this.loadTextures(width, height, textures, cubeMap)
	if shadowMap {
		var none [1]uint32
		none[0] = gl.NONE
		gl.DrawBuffers(1, none[:])
		handleOpenGLES3Error("RenderTexture", this.Name, "glDrawBuffer")
		gl.ReadBuffer(gl.NONE)
		handleOpenGLES3Error("RenderTexture", this.Name, "glReadBuffer")
	}
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		handleOpenGLES3Error("RenderTexture", this.Name, "glCheckFramebufferStatus")
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "RenderTexture", this.Name, "Framebuffer is not complete")
		gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(screenFramebuffer))
		currentlyBoundRT = this.prevRT
		return
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	handleOpenGLES3Error("RenderTexture", this.Name, "glBindFramebuffer with 0")

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

func (this *OpenGLES3RenderTexture) Load(data []byte, width, height int, shadowMap bool) {
}

func (ogltex *OpenGLES3RenderTexture) LoadFromImage(img image.Image) {
}

func (this *OpenGLES3RenderTexture) GetName() string {
	return this.Name
}

func (this *OpenGLES3RenderTexture) SetAsTarget() {
	if currentlyBoundRT == nil {
		var data [1]int32
		gl.GetIntegerv(gl.DRAW_FRAMEBUFFER_BINDING, data[:])
		screenFramebuffer = data[0]
	}
	this.prevRT = currentlyBoundRT
	currentlyBoundRT = this
	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	handleOpenGLES3Error("RenderTexture", this.Name, "glBindFramebuffer in SetAsTarget")
	this.prevViewport = gohome.Render.GetViewport()
	gohome.Render.SetViewport(this.viewport)
}

func (this *OpenGLES3RenderTexture) UnsetAsTarget() {
	if this.prevRT != nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, this.prevRT.fbo)
		currentlyBoundRT = this.prevRT
	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(screenFramebuffer))
		currentlyBoundRT = nil
	}
	handleOpenGLES3Error("RenderTexture", this.Name, "glBindFramebuffer in UnsetAsTarget")
	gohome.Render.SetViewport(this.prevViewport)
}

func (this *OpenGLES3RenderTexture) Blit(rtex gohome.RenderTexture) {
	var ortex *OpenGLES3RenderTexture
	if rtex != nil {
		ortex = rtex.(*OpenGLES3RenderTexture)
	}
	var width int32
	var height int32
	var x int32
	var y int32
	if rtex != nil {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, ortex.fbo)
		width = int32(rtex.GetWidth())
		height = int32(rtex.GetHeight())
		x = 0
		y = 0
	} else {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, uint32(screenFramebuffer))
		handleOpenGLES3Error("RenderTexture", this.Name, "glBindFramebuffer with GL_DRAW_FRAMEBUFFER in Blit")
		width = int32(this.prevViewport.Width)
		height = int32(this.prevViewport.Height)
		x = int32(this.prevViewport.X)
		y = int32(this.prevViewport.Y)
	}
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, this.fbo)
	handleOpenGLES3Error("RenderTexture", this.Name, "glBindFramebuffer with GL_READ_FRAMEBUFFER in Blit")

	gl.BlitFramebuffer(0, 0, int32(this.GetWidth()), int32(this.GetHeight()), x, y, width, height, gl.COLOR_BUFFER_BIT, gl.NEAREST)
	handleOpenGLES3Error("RenderTexture", this.Name, "glBlitFramebuffer")

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	handleOpenGLES3Error("RenderTexture", this.Name, "glBindFramebuffer with GL_READ_FRAMEBUFFER and 0 in Blit")
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	handleOpenGLES3Error("RenderTexture", this.Name, "glBindFramebuffer with GL_DRAW_FRAMEBUFFER and 0 in Blit")
}

func (this *OpenGLES3RenderTexture) Bind(unit uint32) {
	this.BindIndex(0, unit)
}

func (this *OpenGLES3RenderTexture) Unbind(unit uint32) {
	this.UnbindIndex(0, unit)
}

func (this *OpenGLES3RenderTexture) BindIndex(index int, unit uint32) {
	if index < len(this.textures) {
		this.textures[index].Bind(unit)
	}
}

func (this *OpenGLES3RenderTexture) UnbindIndex(index int, unit uint32) {
	if index < len(this.textures) {
		this.textures[index].Unbind(unit)
	}
}

func (this *OpenGLES3RenderTexture) GetWidth() int {
	if len(this.textures) == 0 {
		return 0
	} else {
		return this.textures[0].GetWidth()
	}
}

func (this *OpenGLES3RenderTexture) GetHeight() int {
	if len(this.textures) == 0 {
		return 0
	} else {
		return this.textures[0].GetHeight()
	}
}

func (this *OpenGLES3RenderTexture) Terminate() {
	var buf [1]uint32
	buf[0] = this.fbo
	gl.DeleteFramebuffers(1, buf[:])
	if this.depthBuffer {
		buf[0] = this.rbo
		gl.DeleteRenderbuffers(1, buf[:])
	}
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].Terminate()
	}
	this.textures = this.textures[:0]
}

func (this *OpenGLES3RenderTexture) ChangeSize(width, height int) {
	if this.GetWidth() != width || this.GetHeight() != height {
		textures := len(this.textures)
		this.Terminate()
		this.Create(this.Name, width, height, textures, this.depthBuffer, false, this.shadowMap, this.cubeMap)
	}
}

func (this *OpenGLES3RenderTexture) SetFiltering(filtering int) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetFiltering(filtering)
	}
}

func (this *OpenGLES3RenderTexture) SetWrapping(wrapping int) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetWrapping(wrapping)
	}
}

func (this *OpenGLES3RenderTexture) SetBorderColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderColor(col)
	}
}

func (this *OpenGLES3RenderTexture) SetBorderDepth(depth float32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderDepth(depth)
	}
}

func (this *OpenGLES3RenderTexture) GetKeyColor() color.Color {
	if len(this.textures) == 0 {
		return nil
	}
	return this.textures[0].GetKeyColor()
}

func (this *OpenGLES3RenderTexture) GetModColor() color.Color {
	if len(this.textures) == 0 {
		return nil
	}
	return this.textures[0].GetModColor()
}

func (this *OpenGLES3RenderTexture) SetKeyColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetKeyColor(col)
	}
}

func (this *OpenGLES3RenderTexture) SetModColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetModColor(col)
	}
}

func (this *OpenGLES3RenderTexture) GetData() ([]byte, int, int) {
	if len(this.textures) == 0 {
		return nil, 0, 0
	}
	if tex, ok := this.textures[0].(*OpenGLES3Texture); ok {
		return tex.GetData()
	}

	return this.textures[0].GetData()
}
