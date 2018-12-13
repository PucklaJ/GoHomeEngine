package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/go-gl/gl/all-core/gl"
	"image"
	"image/color"
)

var currentlyBoundRT *OpenGLRenderTexture
var screenFramebuffer int32

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
	prevRT       *OpenGLRenderTexture
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
				if render.HasFunctionAvailable("FRAMEBUFFER_TEXTURE") {
					gl.GetError()
					gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, oglcubemap.oglName, 0)
					handleOpenGLError("RenderTexture", this.Name, "glFramebufferTexture with depthBuffer and CubeMap")
				} else {
					gl.GetError()
					gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
					handleOpenGLError("RenderTexture", this.Name, "glFramebufferTexture2D with depthBuffer and CubeMap")
				}
			} else {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, ogltex.bindingPoint(), ogltex.oglName, 0)
				handleOpenGLError("RenderTexture", this.Name, "glFramebufferTexture2D with depthBuffer and TEXTURE2D")
			}
		} else {
			if cubeMap {
				if render.HasFunctionAvailable("FRAMEBUFFER_TEXTURE") {
					gl.GetError()
					gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, oglcubemap.oglName, 0)
					handleOpenGLError("RenderTexture", this.Name, "glFramebufferTexture with CubeMap")
				} else {
					gl.GetError()
					gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
					handleOpenGLError("RenderTexture", this.Name, "glFramebufferTexture2D with CubeMap")
				}
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

func (this *OpenGLRenderTexture) loadRenderBuffer(width, height uint32) {
	if this.depthBuffer {
		gl.GetError()
		gl.GenRenderbuffers(1, &this.rbo)
		handleOpenGLError("RenderTexture", this.Name, "glGenRenderbuffers")
		gl.BindRenderbuffer(gl.RENDERBUFFER, this.rbo)
		handleOpenGLError("RenderTexture", this.Name, "glBindRenderbuffer")
		if this.multiSampled {
			gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, 8, gl.DEPTH24_STENCIL8, int32(width), int32(height))
			handleOpenGLError("RenderTexture", this.Name, "glRenderbufferStorageMultisample")
		} else {
			gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(width), int32(height))
			handleOpenGLError("RenderTexture", this.Name, "glRenderbufferStorage")
		}
		gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, this.rbo)
		handleOpenGLError("RenderTexture", this.Name, "glFramebufferRenderbuffer")
		gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
		handleOpenGLError("RenderTexture", this.Name, "glBindRenderbuffer with 0")
	}
}

func (this *OpenGLRenderTexture) Create(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) {
	if textures == 0 {
		textures = 1
	}

	render, _ := gohome.Render.(*OpenGLRenderer)

	this.Name = name
	this.shadowMap = shadowMap
	this.multiSampled = multiSampled && render.HasFunctionAvailable("MULTISAMPLE")
	this.depthBuffer = depthBuffer && !shadowMap
	this.cubeMap = cubeMap

	gl.GetError()
	gl.GenFramebuffers(1, &this.fbo)
	handleOpenGLError("RenderTexture", this.Name, "glGenFramebuffers")

	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer")

	this.loadRenderBuffer(width, height)
	this.loadTextures(width, height, textures, cubeMap)
	if shadowMap {
		gl.DrawBuffer(gl.NONE)
		handleOpenGLError("RenderTexture", this.Name, "glDrawBuffer")
		gl.ReadBuffer(gl.NONE)
		handleOpenGLError("RenderTexture", this.Name, "glReadBuffer")
	}
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
	if currentlyBoundRT == nil {
		gl.GetIntegerv(gl.DRAW_FRAMEBUFFER_BINDING, &screenFramebuffer)
	}
	this.prevRT = currentlyBoundRT
	currentlyBoundRT = this
	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer in SetAsTarget")
	this.prevViewport = gohome.Render.GetViewport()
	gohome.Render.SetViewport(this.viewport)
}

func (this *OpenGLRenderTexture) UnsetAsTarget() {
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

func (this *OpenGLRenderTexture) Blit(rtex gohome.RenderTexture) {
	var ortex *OpenGLRenderTexture
	if rtex != nil {
		ortex = rtex.(*OpenGLRenderTexture)
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
		handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer with GL_DRAW_FRAMEBUFFER in Blit")
		width = int32(this.prevViewport.Width)
		height = int32(this.prevViewport.Height)
		x = int32(this.prevViewport.X)
		y = int32(this.prevViewport.Y)
	}
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, this.fbo)
	handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer with GL_READ_FRAMEBUFFER in Blit")
	gl.BlitFramebuffer(0, 0, int32(this.GetWidth()), int32(this.GetHeight()), x, y, width, height, gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT|gl.STENCIL_BUFFER_BIT, gl.NEAREST)
	handleOpenGLError("RenderTexture", this.Name, "glBlitFramebuffer")

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer with GL_READ_FRAMEBUFFER and 0 in Blit")
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	handleOpenGLError("RenderTexture", this.Name, "glBindFramebuffer with GL_DRAW_FRAMEBUFFER and 0 in Blit")
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
	if len(this.textures) == 0 {
		return 0
	} else {
		return this.textures[0].GetWidth()
	}
}

func (this *OpenGLRenderTexture) GetHeight() int {
	if len(this.textures) == 0 {
		return 0
	} else {
		return this.textures[0].GetHeight()
	}
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

func (this *OpenGLRenderTexture) GetKeyColor() color.Color {
	if len(this.textures) == 0 {
		return nil
	}
	return this.textures[0].GetKeyColor()
}

func (this *OpenGLRenderTexture) GetModColor() color.Color {
	if len(this.textures) == 0 {
		return nil
	}
	return this.textures[0].GetModColor()
}

func (this *OpenGLRenderTexture) SetKeyColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetKeyColor(col)
	}
}

func (this *OpenGLRenderTexture) SetModColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetModColor(col)
	}
}

func (this *OpenGLRenderTexture) GetData() ([]byte, int, int) {
	if len(this.textures) == 0 {
		return nil, 0, 0
	}
	if tex, ok := this.textures[0].(*OpenGLTexture); ok {
		if !tex.multiSampled {
			return tex.GetData()
		} else {
			if gohome.Render.HasFunctionAvailable("BLIT_FRAMEBUFFER") {
				rtex := CreateOpenGLRenderTexture("Temp", uint32(this.GetWidth()), uint32(this.GetHeight()), 1, false, false, false, false)
				this.Blit(rtex)
				data, width, height := rtex.GetData()
				rtex.Terminate()
				return data, width, height
			}
		}
	}

	return this.textures[0].GetData()
}
