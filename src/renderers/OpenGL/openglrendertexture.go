package renderer

import (
	// "fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/gl/all-core/gl"
	"image"
	"image/color"
)

type OpenGLRenderTexture struct {
	Name         string
	fbo          uint32
	rbo          uint32
	multiSampled bool
	depthBuffer  bool
	shadowMap    bool
	cubeMap      bool
	textures     []gohome.Texture
	prevViewport gohome.Viewport
	viewport     gohome.Viewport
}

func CreateOpenGLRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) *OpenGLRenderTexture {
	rt := &OpenGLRenderTexture{}

	rt.Create(name, width, height, textures, depthBuffer, multiSampled, shadowMap, cubeMap)

	return rt
}

func (this *OpenGLRenderTexture) loadTextures(width, height, textures uint32, cubeMap bool) {
	var i uint32
	render, _ := gohome.Render.(*OpenGLRenderer)
	for i = 0; i < textures; i++ {
		var ogltex *OpenGLTexture
		var oglcubemap *OpenGLCubeMap
		var texture gohome.Texture
		if cubeMap {
			oglcubemap = CreateOpenGLCubeMap(this.Name)
			texture = oglcubemap
		} else {
			ogltex = CreateOpenGLTexture(this.Name, this.multiSampled)
			texture = ogltex
		}
		texture.Load(nil, int(width), int(height), this.shadowMap)
		if cubeMap {
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, oglcubemap.oglName)
		} else {
			gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
		}
		if this.shadowMap {
			if cubeMap {
				if render.hasFunctionAvailable("FRAMEBUFFER_TEXTURE") {
					gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, oglcubemap.oglName, 0)
				} else {
					gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
				}
			} else {
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, ogltex.bindingPoint(), ogltex.oglName, 0)
			}
		} else {
			if cubeMap {
				if render.hasFunctionAvailable("FRAMEBUFFER_TEXTURE") {
					gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, oglcubemap.oglName, 0)
				} else {
					gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
				}
			} else {
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, ogltex.bindingPoint(), ogltex.oglName, 0)
			}
		}
		if !cubeMap {
			texture.SetFiltering(gohome.FILTERING_LINEAR)
		}
		if cubeMap {
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
		} else {
			gl.BindTexture(ogltex.bindingPoint(), 0)
		}
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

func (this *OpenGLRenderTexture) Create(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) {
	if textures == 0 {
		textures = 1
	}

	render, _ := gohome.Render.(*OpenGLRenderer)

	this.Name = name
	this.shadowMap = shadowMap
	this.multiSampled = multiSampled && render.hasFunctionAvailable("MULTISAMPLE")
	this.depthBuffer = depthBuffer && !shadowMap
	this.cubeMap = cubeMap

	gl.GenFramebuffers(1, &this.fbo)

	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)

	this.loadRenderBuffer(width, height)
	this.loadTextures(width, height, textures, cubeMap)
	if shadowMap {
		gl.DrawBuffer(gl.NONE)
		gl.ReadBuffer(gl.NONE)
	}
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "RenderTexture", this.Name, "Framebuffer is not complete")
		return
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	this.viewport = gohome.Viewport{
		0,
		0, 0,
		int(width), int(height),
	}
}

func (this *OpenGLRenderTexture) Load(data []byte, width, height int, shadowMap bool) error {
	return &OpenGLError{errorString: "The Load method of RenderTexture is not used!"}
}

func (ogltex *OpenGLRenderTexture) LoadFromImage(img image.Image) error {
	return &OpenGLError{errorString: "The LoadFromImage method of RenderTexture is not used!"}
}

func (this *OpenGLRenderTexture) GetName() string {
	return this.Name
}

func (this *OpenGLRenderTexture) SetAsTarget() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	gohome.Render.ClearScreen(gohome.Render.GetBackgroundColor())
	this.prevViewport = gohome.Render.GetViewport()
	gohome.Render.SetViewport(this.viewport)
}

func (this *OpenGLRenderTexture) UnsetAsTarget() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gohome.Render.SetViewport(this.prevViewport)
}

func (this *OpenGLRenderTexture) Blit(rtex gohome.RenderTexture) {
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
	if index < uint32(len(this.textures)) {
		this.textures[index].Bind(unit)
	}
}

func (this *OpenGLRenderTexture) UnbindIndex(index, unit uint32) {
	if index < uint32(len(this.textures)) {
		this.textures[index].Unbind(unit)
	}
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
		this.Create(this.Name, width, height, textures, this.depthBuffer, this.multiSampled, this.shadowMap, this.cubeMap)
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

func (this *OpenGLRenderTexture) SetBorderColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderColor(col)
	}
}

func (this *OpenGLRenderTexture) SetBorderDepth(depth float32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderDepth(depth)
	}
}
