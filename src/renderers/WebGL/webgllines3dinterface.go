package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
)

type WebGLLines3DInterface struct {
	Name   string
	vbo    *js.Object
	loaded bool

	lines       []gohome.Line3D
	numVertices int
}

func (this *WebGLLines3DInterface) Init() {
	this.loaded = false
}

func (this *WebGLLines3DInterface) AddLines(lines []gohome.Line3D) {
	if this.loaded {
		gohome.ErrorMgr.Warning("Lines3DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
		return
	}

	this.lines = append(this.lines, lines...)
}

func (this *WebGLLines3DInterface) GetLines() []gohome.Line3D {
	return this.lines
}

func (this *WebGLLines3DInterface) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, gohome.LINE3DVERTEXSIZE, 0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, gohome.LINE3DVERTEXSIZE, 3*4)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *WebGLLines3DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = 2 * len(this.lines)
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Lines3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	vertexBuffer := gohome.Lines3DToFloatArray(this.lines)

	this.vbo = gl.CreateBuffer()

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, gohome.LINE3DVERTEXSIZE*this.numVertices, vertexBuffer, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	this.loaded = true
}

func (this *WebGLLines3DInterface) Render() {
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
	gl.DrawArrays(gl.LINES, 0, this.numVertices)
	handleWebGLError("Lines3DInterface", this.Name, "RenderError: ")
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *WebGLLines3DInterface) Terminate() {
	gl.DeleteBuffer(this.vbo)
	this.numVertices = 0
	this.loaded = false
	this.lines = this.lines[:0]
}
