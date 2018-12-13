package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles2"
	"unsafe"
)

type OpenGLESLines3DInterface struct {
	Name   string
	vbo    uint32
	loaded bool

	lines       []gohome.Line3D
	numVertices uint32
}

func (this *OpenGLESLines3DInterface) Init() {
	this.loaded = false
}

func (this *OpenGLESLines3DInterface) AddLines(lines []gohome.Line3D) {
	if this.loaded {
		gohome.ErrorMgr.Warning("Lines3DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
		return
	}

	this.lines = append(this.lines, lines...)
}

func (this *OpenGLESLines3DInterface) GetLines() []gohome.Line3D {
	return this.lines
}

func (this *OpenGLESLines3DInterface) attributePointer() {
	offset0 := 0
	offset1 := 3 * 4

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, int32(gohome.LINE3D_VERTEX_SIZE), unsafe.Pointer(&offset0))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, gl.FALSE, int32(gohome.LINE3D_VERTEX_SIZE), unsafe.Pointer(&offset1))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *OpenGLESLines3DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = uint32(2 * len(this.lines))
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Lines3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	var buf [1]uint32
	gl.GenBuffers(1, buf[:])
	this.vbo = buf[0]

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(gohome.LINE3D_VERTEX_SIZE*this.numVertices), unsafe.Pointer(&this.lines[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	this.loaded = true
}

func (this *OpenGLESLines3DInterface) Render() {
	hasLoaded := this.loaded
	if !hasLoaded {
		this.Load()
	}

	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Lines3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	this.attributePointer()
	gl.GetError()
	gl.DrawArrays(gl.LINES, 0, int32(this.numVertices))
	handleOpenGLError("Lines3DInterface", this.Name, "RenderError: ")
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *OpenGLESLines3DInterface) Terminate() {
	var buf [1]uint32
	buf[0] = this.vbo
	defer gl.DeleteBuffers(1, buf[:])
	this.numVertices = 0
	this.loaded = false
	this.lines = this.lines[:0]
}
