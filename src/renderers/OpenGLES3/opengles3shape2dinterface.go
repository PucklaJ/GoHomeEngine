package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles3"
	"unsafe"
)

type OpenGLES3Shape2DInterface struct {
	Name       string
	vbo        uint32
	vao        uint32
	canUseVaos bool
	loaded     bool

	points         []gohome.Shape2DVertex
	numVertices    int
	openglDrawMode uint32
	lineWidth      float32
}

func (this *OpenGLES3Shape2DInterface) Init() {
	render := gohome.Render.(*OpenGLES3Renderer)
	this.canUseVaos = render.HasFunctionAvailable("VERTEX_ARRAY")
	this.loaded = false
	this.openglDrawMode = gl.POINTS
}

func (this *OpenGLES3Shape2DInterface) checkVertices() bool {
	if this.loaded {
		gohome.ErrorMgr.Warning("Shape2DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
		return false
	}
	return true
}

func (this *OpenGLES3Shape2DInterface) AddLines(lines []gohome.Line2D) {
	if this.checkVertices() {
		for i := 0; i < len(lines); i++ {
			this.points = append(this.points, lines[i][:]...)
		}
	}
}

func (this *OpenGLES3Shape2DInterface) AddPoints(points []gohome.Shape2DVertex) {
	if this.checkVertices() {
		this.points = append(this.points, points...)
	}
}

func (this *OpenGLES3Shape2DInterface) AddTriangles(tris []gohome.Triangle2D) {
	if this.checkVertices() {
		for i := 0; i < len(tris); i++ {
			this.points = append(this.points, tris[i][:]...)
		}
	}
}

func (this *OpenGLES3Shape2DInterface) GetPoints() []gohome.Shape2DVertex {
	return this.points
}

func (this *OpenGLES3Shape2DInterface) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, gl.FALSE, (2+4)*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, gl.FALSE, (2+4)*4, gl.PtrOffset(2*4))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *OpenGLES3Shape2DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = len(this.points)
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape2DInterface", this.Name, "No Vertices have been added!")
		return
	}

	var buf [1]uint32
	gl.GenBuffers(1, buf[:])
	this.vbo = buf[0]
	if this.canUseVaos {
		gl.GenVertexArrays(1, buf[:])
		this.vao = buf[0]
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int((2+4)*4*this.numVertices), unsafe.Pointer(&this.points[0][0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	if this.canUseVaos {
		gl.BindVertexArray(this.vao)
		this.attributePointer()
		gl.BindVertexArray(0)
	}

	this.loaded = true
}

func (this *OpenGLES3Shape2DInterface) Render() {
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

	gl.LineWidth(this.lineWidth)

	gl.GetError()
	gl.DrawArrays(this.openglDrawMode, 0, int32(this.numVertices))
	handleOpenGLES3Error("Shape2DInterface", this.Name, "RenderError: ")

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
func (this *OpenGLES3Shape2DInterface) Terminate() {
	var buf [1]uint32
	buf[0] = this.vbo
	gl.DeleteBuffers(1, buf[:])
	if this.canUseVaos {
		buf[0] = this.vao
		gl.DeleteVertexArrays(1, buf[:])
	}
	this.numVertices = 0
	this.loaded = false
	this.points = this.points[:0]
	this.openglDrawMode = gl.POINTS
}

func (this *OpenGLES3Shape2DInterface) SetDrawMode(mode uint8) {
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

func (this *OpenGLES3Shape2DInterface) SetPointSize(size float32) {
}
func (this *OpenGLES3Shape2DInterface) SetLineWidth(width float32) {
	this.lineWidth = width
}
