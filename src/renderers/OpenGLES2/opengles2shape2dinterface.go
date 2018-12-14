package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles2"
	"unsafe"
)

type OpenGLES2Shape2DInterface struct {
	Name   string
	vbo    uint32
	loaded bool

	points         []gohome.Shape2DVertex
	numVertices    uint32
	openglDrawMode uint32
	pointSize      float32
	lineWidth      float32
}

func (this *OpenGLES2Shape2DInterface) Init() {
	this.loaded = false
	this.openglDrawMode = gl.POINTS
}

func (this *OpenGLES2Shape2DInterface) checkVertices() bool {
	if this.loaded {
		gohome.ErrorMgr.Warning("Shape2DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
		return false
	}
	return true
}

func (this *OpenGLES2Shape2DInterface) AddLines(lines []gohome.Line2D) {
	if this.checkVertices() {
		for i := 0; i < len(lines); i++ {
			this.points = append(this.points, lines[i][:]...)
		}
	}
}

func (this *OpenGLES2Shape2DInterface) AddPoints(points []gohome.Shape2DVertex) {
	if this.checkVertices() {
		this.points = append(this.points, points...)
	}
}

func (this *OpenGLES2Shape2DInterface) AddTriangles(tris []gohome.Triangle2D) {
	if this.checkVertices() {
		for i := 0; i < len(tris); i++ {
			this.points = append(this.points, tris[i][:]...)
		}
	}
}

func (this *OpenGLES2Shape2DInterface) GetPoints() []gohome.Shape2DVertex {
	return this.points
}

func (this *OpenGLES2Shape2DInterface) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, gl.FALSE, (2+4)*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, gl.FALSE, (2+4)*4, gl.PtrOffset(2*4))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *OpenGLES2Shape2DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = uint32(len(this.points))
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape2DInterface", this.Name, "No Vertices have been added!")
		return
	}

	var buf [1]uint32
	gl.GenBuffers(1, buf[:])
	this.vbo = buf[0]

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int((2+4)*4*this.numVertices), unsafe.Pointer(&this.points[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	this.loaded = true
}

func (this *OpenGLES2Shape2DInterface) Render() {
	hasLoaded := this.loaded
	if !hasLoaded {
		this.Load()
	}

	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape2DInterface", this.Name, "No Vertices have been added!")
		return
	}

	this.attributePointer()

	gl.LineWidth(this.lineWidth)

	gl.GetError()
	gl.DrawArrays(this.openglDrawMode, 0, int32(this.numVertices))
	handleOpenGLError("Shape2DInterface", this.Name, "RenderError: ")

	gl.LineWidth(1.0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *OpenGLES2Shape2DInterface) Terminate() {
	var buf [1]uint32
	buf[0] = this.vbo
	defer gl.DeleteBuffers(1, buf[:])
	this.numVertices = 0
	this.loaded = false
	this.points = this.points[:0]
	this.openglDrawMode = gl.POINTS
}

func (this *OpenGLES2Shape2DInterface) SetDrawMode(mode uint8) {
	switch mode {
	case gohome.DRAW_MODE_POINTS:
		this.openglDrawMode = gl.POINTS
	case gohome.DRAW_MODE_LINES:
		this.openglDrawMode = gl.LINES
	case gohome.DRAW_MODE_TRIANGLES:
		this.openglDrawMode = gl.TRIANGLES
	default:
		this.openglDrawMode = gl.POINTS
	}
}

func (this *OpenGLES2Shape2DInterface) SetPointSize(size float32) {
	this.pointSize = size
}
func (this *OpenGLES2Shape2DInterface) SetLineWidth(width float32) {
	this.lineWidth = width
}
