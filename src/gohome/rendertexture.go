package gohome

import (
	// "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"log"
)

type RenderTexture interface {
	Load(data []byte, width, height int) error // Is not used. It there just make RenderTexture able to be a Texture
	GetName() string
	SetAsTarget()
	UnsetAsTarget()
	Blit(rtex RenderTexture)
	Bind(unit uint32)
	Unbind(unit uint32)
	GetWidth() int
	GetHeight() int
	ChangeSize(width, height uint32)
	Terminate()
	SetFiltering(filtering uint32)
	SetWrapping(wrapping uint32)
}

type OpenGLRenderTexture struct {
	Name         string
	fbo          uint32
	rbo          uint32
	multiSampled bool
	depthBuffer  bool
	textures     []*OpenGLTexture
	prevViewport Viewport
	viewport     Viewport
}

func CreateOpenGLRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled bool) *OpenGLRenderTexture {
	rt := &OpenGLRenderTexture{}

	rt.Create(name, width, height, textures, depthBuffer, multiSampled)

	return rt
}

func (this *OpenGLRenderTexture) loadTextures(width, height, textures uint32) {
	var i uint32
	for i = 0; i < textures; i++ {
		texture := CreateOpenGLTexture(this.Name, this.multiSampled)
		texture.Load(nil, int(width), int(height))
		texture.Bind(0)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, texture.bindingPoint(), texture.oglName, 0)
		texture.Unbind(0)
		this.textures = append(this.textures, texture)
	}
}

func (this *OpenGLRenderTexture) loadRenderBuffer(width, height uint32) {
	if this.depthBuffer {
		gl.GenRenderbuffers(1, &this.rbo)
		gl.BindRenderbuffer(gl.RENDERBUFFER, this.rbo)
		if this.multiSampled {
			gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, 8, gl.DEPTH24_STENCIL8, int32(width), int32(height))
		} else {
			gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(width), int32(height))
		}
		gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, this.rbo)
		gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	}
}

func (this *OpenGLRenderTexture) Create(name string, width, height, textures uint32, depthBuffer, multiSampled bool) {
	if textures == 0 {
		textures = 1
	}

	this.Name = name
	this.multiSampled = multiSampled
	this.depthBuffer = depthBuffer

	gl.GenFramebuffers(1, &this.fbo)

	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)

	this.loadRenderBuffer(width, height)
	this.loadTextures(width, height, textures)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		log.Println("Error creating RenderTexture: Framebuffer is not complete")
		return
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	this.viewport = Viewport{
		0,
		0, 0,
		int(width), int(height),
	}
}

func (this *OpenGLRenderTexture) Load(data []byte, width, height int) error {
	return &OpenGLError{errorString: "The Load method of RenderTexture is not used!"}
}

func (this *OpenGLRenderTexture) GetName() string {
	return this.Name
}

func (this *OpenGLRenderTexture) SetAsTarget() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	Render.ClearScreen(&Color{0, 0, 0, 0}, 0.0)
	this.prevViewport = Render.GetViewport()
	Render.SetViewport(this.viewport)
}

func (this *OpenGLRenderTexture) UnsetAsTarget() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	Render.SetViewport(this.prevViewport)
}

func (this *OpenGLRenderTexture) Blit(rtex RenderTexture) {
	var width int32
	var height int32
	var x int32
	var y int32
	if rtex != nil {
		rtex.SetAsTarget()
		width = int32(rtex.GetWidth())
		height = int32(rtex.GetHeight())
		x = 0
		y = 0
	} else {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
		width = int32(this.prevViewport.Width)
		height = int32(this.prevViewport.Height)
		x = int32(this.prevViewport.X)
		y = int32(this.prevViewport.Y)
	}
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, this.fbo)
	gl.BlitFramebuffer(0, 0, int32(this.GetWidth()), int32(this.GetHeight()), x, y, width, height, gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT|gl.STENCIL_BUFFER_BIT, gl.NEAREST)
	if rtex != nil {
		rtex.UnsetAsTarget()
	} else {
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	}
}

func (this *OpenGLRenderTexture) Bind(unit uint32) {
	this.BindIndex(0, unit)
}

func (this *OpenGLRenderTexture) Unbind(unit uint32) {
	this.UnbindIndex(0, unit)
}

func (this *OpenGLRenderTexture) BindIndex(index, unit uint32) {
	this.textures[index].Bind(unit)
}

func (this *OpenGLRenderTexture) UnbindIndex(index, unit uint32) {
	this.textures[index].Unbind(unit)
}

func (this *OpenGLRenderTexture) GetWidth() int {
	return this.textures[0].GetWidth()
}

func (this *OpenGLRenderTexture) GetHeight() int {
	return this.textures[0].GetHeight()
}

func (this *OpenGLRenderTexture) Terminate() {
	defer gl.DeleteFramebuffers(1, &this.fbo)
	if this.depthBuffer {
		defer gl.DeleteRenderbuffers(1, &this.rbo)
	}
	for i := 0; i < len(this.textures); i++ {
		defer this.textures[i].Terminate()
	}
	this.textures = append(this.textures[:0], this.textures[len(this.textures):]...)
}

func (this *OpenGLRenderTexture) ChangeSize(width, height uint32) {
	if uint32(this.GetWidth()) != width || uint32(this.GetHeight()) != height {
		textures := uint32(len(this.textures))
		this.Terminate()
		this.Create(this.Name, width, height, textures, this.depthBuffer, this.multiSampled)
	}
}

func (this *OpenGLRenderTexture) SetFiltering(filtering uint32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetFiltering(filtering)
	}
}

func (this *OpenGLRenderTexture) SetWrapping(wrapping uint32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetWrapping(wrapping)
	}
}
