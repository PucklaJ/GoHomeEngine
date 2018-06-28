package renderer

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/exp/f32"
	"encoding/binary"
)

type OpenGLESLines3DInterface struct {
	Name       string
	vbo        gl.Buffer
	vao        gl.VertexArray
	canUseVaos bool
	loaded     bool

	lines       []gohome.Line3D
	numVertices uint32

	gles 		*gl.Context
}

func (this *OpenGLESLines3DInterface) Init() {
	render := gohome.Render.(*OpenGLESRenderer)
	this.gles = render.GetContext()
	if _,ok := (*render.GetContext()).(gl.Context3); ok {
		this.canUseVaos = true
	} else {
		this.canUseVaos = false
	}
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
	(*this.gles).BindBuffer(gl.ARRAY_BUFFER,this.vbo)

	(*this.gles).VertexAttribPointer(gl.Attrib{0},3,gl.FLOAT,false,int(gohome.LINE3D_VERTEX_SIZE),0)
	(*this.gles).VertexAttribPointer(gl.Attrib{1},4,gl.FLOAT,false,int(gohome.LINE3D_VERTEX_SIZE),3*4)

	(*this.gles).EnableVertexAttribArray(gl.Attrib{0})
	(*this.gles).EnableVertexAttribArray(gl.Attrib{1})
}

func (this *OpenGLESLines3DInterface) toByteArray() ([]byte) {
	var verticesBytes []byte

	verticesFloats := make([]float32, this.numVertices*gohome.LINE3D_VERTEX_SIZE/4*2)
	var index uint32
	for i := 0; uint32(i) < this.numVertices/2; i++ {
		for k := 0; k<2;k++ {
			for j := 0; uint32(j) < gohome.LINE3D_VERTEX_SIZE/4; j++ {
				verticesFloats[index+uint32(j)] = this.lines[i][k][j]
			}
			index += gohome.LINE3D_VERTEX_SIZE/4
		}
	}

	verticesBytes = f32.Bytes(binary.LittleEndian, verticesFloats...)

	return verticesBytes
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

	this.vbo = (*this.gles).CreateBuffer()
	if this.canUseVaos {
		this.vao = (*this.gles).CreateVertexArray()
	}

	(*this.gles).BindBuffer(gl.ARRAY_BUFFER,this.vbo)
	(*this.gles).BufferData(gl.ARRAY_BUFFER,this.toByteArray(),gl.STATIC_DRAW)
	(*this.gles).BindBuffer(gl.ARRAY_BUFFER,gl.Buffer{0})

	if this.canUseVaos {
		(*this.gles).BindVertexArray(this.vao)
		this.attributePointer()
		(*this.gles).BindVertexArray(gl.VertexArray{0})
	}

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

	if this.canUseVaos {
		(*this.gles).BindVertexArray(this.vao)
	} else {
		this.attributePointer()
	}
	(*this.gles).GetError()
	(*this.gles).DrawArrays(gl.LINES,0,int(this.numVertices))
	handleOpenGLESError("Lines3DInterface", this.Name, "RenderError: ")
	if this.canUseVaos {
		(*this.gles).BindVertexArray(gl.VertexArray{0})
	} else {
		(*this.gles).BindBuffer(gl.ARRAY_BUFFER,gl.Buffer{0})
	}

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *OpenGLESLines3DInterface) Terminate() {
	defer (*this.gles).DeleteBuffer(this.vbo)
	if this.canUseVaos {
		defer (*this.gles).DeleteVertexArray(this.vao)
	}
	this.numVertices = 0
	this.loaded = false
	this.lines = this.lines[:0]
}

func CreateOpenGLESLines3DInterface(name string) *OpenGLESLines3DInterface {
	return &OpenGLESLines3DInterface{Name: name}
}
