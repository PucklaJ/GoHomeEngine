package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
)

type WebGLShape2DInterface struct {
	Name   string
	vbo    *js.Object
	loaded bool

	points        []gohome.Shape2DVertex
	numVertices   uint32
	WebGLDrawMode int
	pointSize     float32
	lineWidth     float32
}

func (this *WebGLShape2DInterface) Init() {
	this.loaded = false
	this.WebGLDrawMode = gl.POINTS
}

func (this *WebGLShape2DInterface) checkVertices() bool {
	if this.loaded {
		gohome.ErrorMgr.Warning("Shape2DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
		return false
	}
	return true
}

func (this *WebGLShape2DInterface) AddLines(lines []gohome.Line2D) {
	if this.checkVertices() {
		for i := 0; i < len(lines); i++ {
			this.points = append(this.points, lines[i][:]...)
		}
	}
}

func (this *WebGLShape2DInterface) AddPoints(points []gohome.Shape2DVertex) {
	if this.checkVertices() {
		this.points = append(this.points, points...)
	}
}

func (this *WebGLShape2DInterface) AddTriangles(tris []gohome.Triangle2D) {
	if this.checkVertices() {
		for i := 0; i < len(tris); i++ {
			this.points = append(this.points, tris[i][:]...)
		}
	}
}

func (this *WebGLShape2DInterface) GetPoints() []gohome.Shape2DVertex {
	return this.points
}

func (this *WebGLShape2DInterface) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, (2+4)*4, 0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, (2+4)*4, 2*4)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *WebGLShape2DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = uint32(len(this.points))
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape2DInterface", this.Name, "No Vertices have been added!")
		return
	}

	this.vbo = gl.CreateBuffer()

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int((2+4)*4*this.numVertices), this.points, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	this.loaded = true
}

func (this *WebGLShape2DInterface) Render() {
	hasLoaded := this.loaded
	if !hasLoaded {
		this.Load()
	}

	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape2DInterface", this.Name, "No Vertices have been added!")
		return
	}

	this.attributePointer()

	gl.LineWidth(float64(this.lineWidth))

	gl.GetError()
	gl.DrawArrays(this.WebGLDrawMode, 0, int(this.numVertices))
	handleWebGLError("Shape2DInterface", this.Name, "RenderError: ")

	gl.LineWidth(1.0)

	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *WebGLShape2DInterface) Terminate() {
	gl.DeleteBuffer(this.vbo)
	this.numVertices = 0
	this.loaded = false
	this.points = this.points[:0]
	this.WebGLDrawMode = gl.POINTS
}

func (this *WebGLShape2DInterface) SetDrawMode(mode uint8) {
	switch mode {
	case gohome.DRAW_MODE_POINTS:
		this.WebGLDrawMode = gl.POINTS
	case gohome.DRAW_MODE_LINES:
		this.WebGLDrawMode = gl.LINES
	case gohome.DRAW_MODE_TRIANGLES:
		this.WebGLDrawMode = gl.TRIANGLES
	default:
		this.WebGLDrawMode = gl.POINTS
	}
}

func (this *WebGLShape2DInterface) SetPointSize(size float32) {
	this.pointSize = size
}
func (this *WebGLShape2DInterface) SetLineWidth(width float32) {
	this.lineWidth = width
}
