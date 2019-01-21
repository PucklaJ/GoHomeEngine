package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/go-gl/gl/all-core/gl"
)

type OpenGLShape3DInterface struct {
	Name       string
	vbo        uint32
	vao        uint32
	canUseVaos bool
	loaded     bool

	points      []gohome.Shape3DVertex
	drawMode    uint32
	numVertices int
	pointSize   float32
	lineWidth   float32
}

func (this *OpenGLShape3DInterface) Init() {
	render := gohome.Render.(*OpenGLRenderer)
	this.canUseVaos = render.HasFunctionAvailable("VERTEX_ARRAY")
	this.loaded = false
}

func (this *OpenGLShape3DInterface) AddPoints(points []gohome.Shape3DVertex) {
	if this.loaded {
		gohome.ErrorMgr.Warning("Shape3DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
		return
	}

	this.points = append(this.points, points...)
}

func (this *OpenGLShape3DInterface) GetPoints() []gohome.Shape3DVertex {
	return this.points
}

func (this *OpenGLShape3DInterface) SetDrawMode(drawMode uint8) {
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

func (this *OpenGLShape3DInterface) SetPointSize(size float32) {
	this.pointSize = size
}

func (this *OpenGLShape3DInterface) SetLineWidth(width float32) {
	this.lineWidth = width
}

func (this *OpenGLShape3DInterface) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, gohome.SHAPE3DVERTEXSIZE, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, gohome.SHAPE3DVERTEXSIZE, gl.PtrOffset(3*4))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *OpenGLShape3DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = len(this.points)
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	gl.GenBuffers(1, &this.vbo)
	if this.canUseVaos {
		gl.GenVertexArrays(1, &this.vao)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(gohome.SHAPE3DVERTEXSIZE*this.numVertices), gl.Ptr(this.points), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	if this.canUseVaos {
		gl.BindVertexArray(this.vao)
		this.attributePointer()
		gl.BindVertexArray(0)
	}

	this.loaded = true
}

func (this *OpenGLShape3DInterface) Render() {
	hasLoaded := this.loaded
	if !hasLoaded {
		this.Load()
	}

	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	gl.PointSize(this.pointSize)
	gl.LineWidth(this.lineWidth)

	if this.canUseVaos {
		gl.BindVertexArray(this.vao)
	} else {
		this.attributePointer()
	}

	gl.GetError()
	gl.DrawArrays(this.drawMode, 0, int32(this.numVertices))
	handleOpenGLError("Shape3DInterface", this.Name, "RenderError: ")
	if this.canUseVaos {
		gl.BindVertexArray(0)
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	gl.PointSize(1.0)
	gl.LineWidth(1.0)

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *OpenGLShape3DInterface) Terminate() {
	defer gl.DeleteBuffers(1, &this.vbo)
	if this.canUseVaos {
		defer gl.DeleteVertexArrays(1, &this.vao)
	}
	this.numVertices = 0
	this.loaded = false
	this.points = this.points[:0]
}
