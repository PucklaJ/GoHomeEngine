package renderer

import (
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
)

type WebGLShape3DInterface struct {
	Name       string
	vbo        *js.Object
	vao        *js.Object
	canUseVaos bool
	loaded     bool

	points      []gohome.Shape3DVertex
	drawMode    int
	numVertices int
	lineWidth   float64
}

func (this *WebGLShape3DInterface) Init() {
	render := gohome.Render.(*WebGLRenderer)
	this.canUseVaos = render.HasFunctionAvailable("VERTEX_ARRAY")
	this.loaded = false
	this.lineWidth = 1.0
}

func (this *WebGLShape3DInterface) AddPoints(points []gohome.Shape3DVertex) {
	if this.loaded {
		gohome.ErrorMgr.Warning("Shape3DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
		return
	}

	this.points = append(this.points, points...)
}

func (this *WebGLShape3DInterface) GetPoints() []gohome.Shape3DVertex {
	return this.points
}

func (this *WebGLShape3DInterface) SetDrawMode(drawMode uint8) {
	switch drawMode {
	case gohome.DRAW_MODE_POINTS:
		this.drawMode = gl.POINTS
	case gohome.DRAW_MODE_LINES:
		this.drawMode = gl.LINES
	case gohome.DRAW_MODE_TRIANGLES:
		this.drawMode = gl.TRIANGLES
	default:
		this.drawMode = gl.POINTS
	}
}

func (this *WebGLShape3DInterface) SetPointSize(size float32) {
}

func (this *WebGLShape3DInterface) SetLineWidth(width float32) {
	this.lineWidth = float64(width)
}

func (this *WebGLShape3DInterface) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, gohome.SHAPE3DVERTEXSIZE, 0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, gohome.SHAPE3DVERTEXSIZE, 3*4)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *WebGLShape3DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = len(this.points)
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	this.vbo = gl.CreateBuffer()
	if this.canUseVaos {
		this.vao = gl.CreateVertexArray()
	}

	vertexBuffer := gohome.Shape3DVerticesToFloatArray(this.points)

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(gohome.SHAPE3DVERTEXSIZE*this.numVertices), vertexBuffer, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	if this.canUseVaos {
		gl.BindVertexArray(this.vao)
		this.attributePointer()
		gl.BindVertexArray(nil)
	}

	this.loaded = true
}

func (this *WebGLShape3DInterface) Render() {
	hasLoaded := this.loaded
	if !hasLoaded {
		this.Load()
	}

	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	gl.LineWidth(this.lineWidth)

	if this.canUseVaos {
		gl.BindVertexArray(this.vao)
	} else {
		this.attributePointer()
	}

	gl.GetError()
	gl.DrawArrays(this.drawMode, 0, this.numVertices)
	handleWebGLError("Shape3DInterface", this.Name, "RenderError: ")
	if this.canUseVaos {
		gl.BindVertexArray(nil)
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	}

	gl.LineWidth(1.0)

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *WebGLShape3DInterface) Terminate() {
	gl.DeleteBuffer(this.vbo)
	if this.canUseVaos {
		gl.DeleteVertexArray(this.vao)
	}
	this.numVertices = 0
	this.loaded = false
	this.points = this.points[:0]
}
