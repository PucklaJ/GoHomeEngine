package renderer

import (
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/go-gl/gl/all-core/gl"
)

type OpenGLShape2DInterface struct {
	Name       string
	vbo        uint32
	vao        uint32
	canUseVaos bool
	loaded     bool

	points         []gohome.Shape2DVertex
	numVertices    int
	openglDrawMode uint32
	pointSize      float32
	lineWidth      float32
}

func (this *OpenGLShape2DInterface) Init() {
	render := gohome.Render.(*OpenGLRenderer)
	this.canUseVaos = render.HasFunctionAvailable("VERTEX_ARRAY")
	this.loaded = false
	this.openglDrawMode = gl.POINTS
}

func (this *OpenGLShape2DInterface) checkVertices() bool {
	if this.loaded {
		gohome.ErrorMgr.Warning("Shape2DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
		return false
	}
	return true
}

func (this *OpenGLShape2DInterface) AddLines(lines []gohome.Line2D) {
	if this.checkVertices() {
		for i := 0; i < len(lines); i++ {
			this.points = append(this.points, lines[i][:]...)
		}
	}
}

func (this *OpenGLShape2DInterface) AddPoints(points []gohome.Shape2DVertex) {
	if this.checkVertices() {
		this.points = append(this.points, points...)
	}
}

func (this *OpenGLShape2DInterface) AddTriangles(tris []gohome.Triangle2D) {
	if this.checkVertices() {
		for i := 0; i < len(tris); i++ {
			this.points = append(this.points, tris[i][:]...)
		}
	}
}

func (this *OpenGLShape2DInterface) GetPoints() []gohome.Shape2DVertex {
	return this.points
}

func (this *OpenGLShape2DInterface) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, (2+4)*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, (2+4)*4, gl.PtrOffset(2*4))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *OpenGLShape2DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = len(this.points)
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape2DInterface", this.Name, "No Vertices have been added!")
		return
	}

	gl.GenBuffers(1, &this.vbo)
	if this.canUseVaos {
		gl.GenVertexArrays(1, &this.vao)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int((2+4)*4*this.numVertices), gl.Ptr(this.points), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	if this.canUseVaos {
		gl.BindVertexArray(this.vao)
		this.attributePointer()
		gl.BindVertexArray(0)
	}

	this.loaded = true
}

func (this *OpenGLShape2DInterface) Render() {
	hasLoaded := this.loaded
	if !hasLoaded {
		this.Load()
	}

	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape2DInterface", this.Name, "No Vertices have been added!")
		return
	}

	if this.canUseVaos {
		gl.BindVertexArray(this.vao)
	} else {
		this.attributePointer()
	}

	gl.PointSize(this.pointSize)
	gl.LineWidth(this.lineWidth)

	gl.GetError()
	gl.DrawArrays(this.openglDrawMode, 0, int32(this.numVertices))
	handleOpenGLError("Shape2DInterface", this.Name, "RenderError: ")

	gl.PointSize(1.0)
	gl.LineWidth(1.0)

	if this.canUseVaos {
		gl.BindVertexArray(0)
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *OpenGLShape2DInterface) Terminate() {
	defer gl.DeleteBuffers(1, &this.vbo)
	if this.canUseVaos {
		defer gl.DeleteVertexArrays(1, &this.vao)
	}
	this.numVertices = 0
	this.loaded = false
	this.points = this.points[:0]
	this.openglDrawMode = gl.POINTS
}

func (this *OpenGLShape2DInterface) SetDrawMode(mode uint8) {
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

func (this *OpenGLShape2DInterface) SetPointSize(size float32) {
	this.pointSize = size
}
func (this *OpenGLShape2DInterface) SetLineWidth(width float32) {
	this.lineWidth = width
}
