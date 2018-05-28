package renderer

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/gl/all-core/gl"
)

type OpenGLLines3DInterface struct {
	Name       string
	vbo        uint32
	vao        uint32
	canUseVaos bool
	loaded     bool

	lines       []gohome.Line3D
	numVertices uint32
}

func (this *OpenGLLines3DInterface) Init() {
	render := gohome.Render.(*OpenGLRenderer)
	this.canUseVaos = render.hasFunctionAvailable("VERTEX_ARRAY")
	this.loaded = false
}

func (this *OpenGLLines3DInterface) AddLines(lines []gohome.Line3D) {
	if this.loaded {
		gohome.ErrorMgr.Warning("Lines3DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
	}

	this.lines = append(this.lines, lines...)
}

func (this *OpenGLLines3DInterface) GetLines() []gohome.Line3D {
	return this.lines
}

func (this *OpenGLLines3DInterface) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(gohome.LINE3D_VERTEX_SIZE), gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, int32(gohome.LINE3D_VERTEX_SIZE), gl.PtrOffset(3*4))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *OpenGLLines3DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = uint32(2 * len(this.lines))
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Lines3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	gl.GenBuffers(1, &this.vbo)
	if this.canUseVaos {
		gl.GenVertexArrays(1, &this.vao)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(gohome.LINE3D_VERTEX_SIZE*this.numVertices), gl.Ptr(this.lines), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	if this.canUseVaos {
		gl.BindVertexArray(this.vao)
		this.attributePointer()
		gl.BindVertexArray(0)
	}

	this.loaded = true
}

func (this *OpenGLLines3DInterface) Render() {
	hasLoaded := this.loaded
	if !hasLoaded {
		this.Load()
	}

	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Lines3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	if this.canUseVaos {
		gl.BindVertexArray(this.vao)
	} else {
		this.attributePointer()
	}
	gl.GetError()
	gl.DrawArrays(gl.LINES, 0, int32(this.numVertices))
	handleOpenGLError("Lines3DInterface", this.Name, "RenderError: ")
	if this.canUseVaos {
		gl.BindVertexArray(0)
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *OpenGLLines3DInterface) Terminate() {
	defer gl.DeleteBuffers(1, &this.vbo)
	if this.canUseVaos {
		defer gl.DeleteVertexArrays(1, &this.vao)
	}
	this.numVertices = 0
	this.loaded = false
}
