package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
	"image"
	"image/color"
)

var currentlyBoundRT *WebGLRenderTexture

type WebGLRenderTexture struct {
	Name         string
	fbo          *js.Object
	rbo          *js.Object
	depthBuffer  bool
	shadowMap    bool
	cubeMap      bool
	textures     []gohome.Texture
	prevViewport gohome.Viewport
	viewport     gohome.Viewport
	prevRT       *WebGLRenderTexture
}

func CreateWebGLRenderTexture(name string, width, height, textures int, depthBuffer, shadowMap, cubeMap bool) *WebGLRenderTexture {
	rt := &WebGLRenderTexture{}

	rt.Create(name, width, height, textures, depthBuffer, false, shadowMap, cubeMap)

	return rt
}

func (this *WebGLRenderTexture) loadTextures(width, height, textures int, cubeMap bool) {
	for i := 0; i < textures; i++ {
		var ogltex *WebGLTexture
		var oglcubemap *WebGLCubeMap
		var texture gohome.Texture
		if cubeMap {
			oglcubemap = CreateWebGLCubeMap(this.Name)
			texture = oglcubemap
		} else {
			ogltex = CreateWebGLTexture(this.Name)
			texture = ogltex
		}
		texture.Load(nil, int(width), int(height), this.shadowMap)
		if cubeMap {
			gl.GetError()
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, oglcubemap.oglName)
			handleWebGLError("RenderTexture", this.Name, "Binding cubemap")
		} else {
			gl.GetError()
			gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
			handleWebGLError("RenderTexture", this.Name, "Binding texture 2d")
		}
		if this.shadowMap {
			if cubeMap {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
				handleWebGLError("RenderTexture", this.Name, "glFramebufferTexture2D with depthBuffer and CubeMap")
			} else {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, ogltex.bindingPoint(), ogltex.oglName, 0)
				handleWebGLError("RenderTexture", this.Name, "glFramebufferTexture2D with depthBuffer and TEXTURE2D")
			}
		} else {
			if cubeMap {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
				handleWebGLError("RenderTexture", this.Name, "glFramebufferTexture2D with CubeMap")
			} else {
				gl.GetError()
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, ogltex.bindingPoint(), ogltex.oglName, 0)
				handleWebGLError("RenderTexture", this.Name, "glFramebufferTexture2D with TEXTURE2D")
			}
		}
		if !cubeMap {
			texture.SetFiltering(gohome.FILTERING_LINEAR)
		}
		if cubeMap {
			gl.GetError()
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, nil)
			handleWebGLError("RenderTexture", this.Name, "glBindTexture with CubeMap")
		} else {
			gl.GetError()
			gl.BindTexture(ogltex.bindingPoint(), nil)
			handleWebGLError("RenderTexture", this.Name, "glBindTexture with TEXTURE2D")
		}
		this.textures = append(this.textures, texture)
	}
}

func (this *WebGLRenderTexture) loadRenderBuffer(width, height int) {
	if this.depthBuffer {
		gl.GetError()
		this.rbo = gl.CreateRenderbuffer()
		handleWebGLError("RenderTexture", this.Name, "glGenRenderbuffers")
		gl.BindRenderbuffer(gl.RENDERBUFFER, this.rbo)
		handleWebGLError("RenderTexture", this.Name, "glBindRenderbuffer")

		gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT16, width, height)
		handleWebGLError("RenderTexture", this.Name, "glRenderbufferStorage")

		gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, this.rbo)
		handleWebGLError("RenderTexture", this.Name, "glFramebufferRenderbuffer")
		gl.BindRenderbuffer(gl.RENDERBUFFER, nil)
		handleWebGLError("RenderTexture", this.Name, "glBindRenderbuffer with 0")
	}
}

func (this *WebGLRenderTexture) Create(name string, width, height, textures int, depthBuffer, multiSampled, shadowMap, cubeMap bool) {
	if textures == 0 {
		textures = 1
	}

	render, _ := gohome.Render.(*WebGLRenderer)

	this.Name = name
	this.shadowMap = shadowMap
	this.depthBuffer = depthBuffer && !shadowMap
	this.cubeMap = cubeMap

	gl.GetError()
	this.fbo = gl.CreateFramebuffer()
	handleWebGLError("RenderTexture", this.Name, "glGenFramebuffers")

	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	handleWebGLError("RenderTexture", this.Name, "glBindFramebuffer")

	this.loadRenderBuffer(int(width), int(height))
	this.loadTextures(int(width), int(height), int(textures), cubeMap)

	if shadowMap {
		if render.HasFunctionAvailable("DRAW_BUFFERS") {
			var buffers [1]int
			buffers[0] = gl.NONE
			gl.DrawBuffers(buffers[:])
			handleWebGLError("RenderTexture", this.Name, "glDrawBuffer")
			gl.ReadBuffer(gl.NONE)
			handleWebGLError("RenderTexture", this.Name, "glReadBuffer")
		}
	}
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		handleWebGLError("RenderTexture", this.Name, "glCheckFramebufferStatus")
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "RenderTexture", this.Name, "Framebuffer is not complete")
		gl.BindFramebuffer(gl.FRAMEBUFFER, nil)
		currentlyBoundRT = this.prevRT
		return
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, nil)
	handleWebGLError("RenderTexture", this.Name, "glBindFramebuffer with 0")

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

func (this *WebGLRenderTexture) Load(data []byte, width, height int, shadowMap bool) {
}

func (ogltex *WebGLRenderTexture) LoadFromImage(img image.Image) {
}

func (this *WebGLRenderTexture) GetName() string {
	return this.Name
}

func (this *WebGLRenderTexture) SetAsTarget() {
	this.prevRT = currentlyBoundRT
	currentlyBoundRT = this
	gl.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	handleWebGLError("RenderTexture", this.Name, "glBindFramebuffer in SetAsTarget")
	this.prevViewport = gohome.Render.GetViewport()
	gohome.Render.SetViewport(this.viewport)
}

func (this *WebGLRenderTexture) UnsetAsTarget() {
	if this.prevRT != nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, this.prevRT.fbo)
		currentlyBoundRT = this.prevRT
	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, nil)
		currentlyBoundRT = nil
	}
	handleWebGLError("RenderTexture", this.Name, "glBindFramebuffer in UnsetAsTarget")
	gohome.Render.SetViewport(this.prevViewport)
}

func (this *WebGLRenderTexture) Blit(rtex gohome.RenderTexture) {
	var ortex *WebGLRenderTexture
	if rtex != nil {
		ortex = rtex.(*WebGLRenderTexture)
	}
	var width int
	var height int
	var x int
	var y int
	if rtex != nil {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, ortex.fbo)
		width = rtex.GetWidth()
		height = rtex.GetHeight()
		x = 0
		y = 0
	} else {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, nil)
		handleWebGLError("RenderTexture", this.Name, "glBindFramebuffer with GL_DRAW_FRAMEBUFFER in Blit")
		width = this.prevViewport.Width
		height = this.prevViewport.Height
		x = this.prevViewport.X
		y = this.prevViewport.Y
	}
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, this.fbo)
	handleWebGLError("RenderTexture", this.Name, "glBindFramebuffer with GL_READ_FRAMEBUFFER in Blit")
	gl.BlitFramebuffer(0, 0, this.GetWidth(), this.GetHeight(), x, y, width, height, gl.COLOR_BUFFER_BIT, gl.NEAREST)
	handleWebGLError("RenderTexture", this.Name, "glBlitFramebuffer")

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, nil)
	handleWebGLError("RenderTexture", this.Name, "glBindFramebuffer with GL_READ_FRAMEBUFFER and 0 in Blit")
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, nil)
	handleWebGLError("RenderTexture", this.Name, "glBindFramebuffer with GL_DRAW_FRAMEBUFFER and 0 in Blit")
}

func (this *WebGLRenderTexture) Bind(unit uint32) {
	this.BindIndex(0, unit)
}

func (this *WebGLRenderTexture) Unbind(unit uint32) {
	this.UnbindIndex(0, unit)
}

func (this *WebGLRenderTexture) BindIndex(index int, unit uint32) {
	if index < len(this.textures) {
		this.textures[index].Bind(unit)
	}
}

func (this *WebGLRenderTexture) UnbindIndex(index int, unit uint32) {
	if index < len(this.textures) {
		this.textures[index].Unbind(unit)
	}
}

func (this *WebGLRenderTexture) GetWidth() int {
	if len(this.textures) == 0 {
		return 0
	} else {
		return this.textures[0].GetWidth()
	}
}

func (this *WebGLRenderTexture) GetHeight() int {
	if len(this.textures) == 0 {
		return 0
	} else {
		return this.textures[0].GetHeight()
	}
}

func (this *WebGLRenderTexture) Terminate() {
	gl.DeleteFramebuffer(this.fbo)
	if this.depthBuffer {
		gl.DeleteRenderbuffer(this.rbo)
	}
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].Terminate()
	}
	this.textures = append(this.textures[:0], this.textures[len(this.textures):]...)
}

func (this *WebGLRenderTexture) ChangeSize(width, height int) {
	if this.GetWidth() != width || this.GetHeight() != height {
		textures := len(this.textures)
		this.Terminate()
		this.Create(this.Name, width, height, textures, this.depthBuffer, false, this.shadowMap, this.cubeMap)
	}
}

func (this *WebGLRenderTexture) SetFiltering(filtering int) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetFiltering(filtering)
	}
}

func (this *WebGLRenderTexture) SetWrapping(wrapping int) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetWrapping(wrapping)
	}
}

func (this *WebGLRenderTexture) SetBorderColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderColor(col)
	}
}

func (this *WebGLRenderTexture) SetBorderDepth(depth float32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderDepth(depth)
	}
}

func (this *WebGLRenderTexture) GetKeyColor() color.Color {
	if len(this.textures) == 0 {
		return nil
	}
	return this.textures[0].GetKeyColor()
}

func (this *WebGLRenderTexture) GetModColor() color.Color {
	if len(this.textures) == 0 {
		return nil
	}
	return this.textures[0].GetModColor()
}

func (this *WebGLRenderTexture) SetKeyColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetKeyColor(col)
	}
}

func (this *WebGLRenderTexture) SetModColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetModColor(col)
	}
}

func (this *WebGLRenderTexture) GetData() ([]byte, int, int) {
	if len(this.textures) == 0 {
		return nil, 0, 0
	}
	if tex, ok := this.textures[0].(*WebGLTexture); ok {
		return tex.GetData()
	}

	return this.textures[0].GetData()
}
