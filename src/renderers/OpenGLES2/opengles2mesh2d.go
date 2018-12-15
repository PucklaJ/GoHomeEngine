package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles2"
	"unsafe"
)

type OpenGLES2Mesh2D struct {
	vertices    []gohome.Mesh2DVertex
	indices     []uint8
	numVertices uint32
	numIndices  uint32
	Name        string
	vbo         uint32
	ibo         uint32
}

func CreateOpenGLES2Mesh2D(name string) *OpenGLES2Mesh2D {
	mesh := OpenGLES2Mesh2D{
		Name: name,
	}

	return &mesh
}

func (oglm *OpenGLES2Mesh2D) deleteElements() {
	oglm.vertices = append(oglm.vertices[:0], oglm.vertices[len(oglm.vertices):]...)
	oglm.indices = append(oglm.indices[:0], oglm.indices[len(oglm.indices):]...)
}

func (oglm *OpenGLES2Mesh2D) AddVertices(vertices []gohome.Mesh2DVertex, indices []uint32) {
	oglm.vertices = append(oglm.vertices, vertices...)
	index := len(oglm.indices)
	oglm.indices = append(oglm.indices, make([]uint8, len(indices))...)
	for id, i := range indices {
		oglm.indices[index+id] = uint8(i)
	}
}

func (oglm *OpenGLES2Mesh2D) attributePointer() {
	gl.GetError()
	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	handleOpenGLError("Mesh2D", oglm.Name, "glBindBuffer vbo in attributePointer: ")
	gl.VertexAttribPointer(0, 2, gl.FLOAT, gl.FALSE, int32(gohome.MESH2DVERTEX_SIZE), gl.PtrOffset(0))
	handleOpenGLError("Mesh2D", oglm.Name, "glVertexAttribPointer 0 in attributePointer: ")
	gl.VertexAttribPointer(1, 2, gl.FLOAT, gl.FALSE, int32(gohome.MESH2DVERTEX_SIZE), gl.PtrOffset(2*4))
	handleOpenGLError("Mesh2D", oglm.Name, "glVertexAttribPointer 1 in attributePointer: ")
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
	handleOpenGLError("Mesh2D", oglm.Name, "glBindBuffer ibo in attributePointer: ")
}

func (oglm *OpenGLES2Mesh2D) Load() {
	oglm.numVertices = uint32(len(oglm.vertices))
	oglm.numIndices = uint32(len(oglm.indices))
	var verticesSize uint32 = oglm.numVertices * gohome.MESH2DVERTEX_SIZE
	var indicesSize uint32 = oglm.numIndices

	var buf [1]uint32
	gl.GenBuffers(1, buf[:])
	handleOpenGLError("Mesh2D", oglm.Name, "glGenBuffers VBO: ")
	oglm.vbo = buf[0]
	gl.GenBuffers(1, buf[:])
	handleOpenGLError("Mesh2D", oglm.Name, "glGenBuffers IBO: ")
	oglm.ibo = buf[0]

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	handleOpenGLError("Mesh2D", oglm.Name, "glBindBuffer vbo in Load: ")
	gl.BufferData(gl.ARRAY_BUFFER, int(verticesSize), unsafe.Pointer(&oglm.vertices[0]), gl.STATIC_DRAW)
	handleOpenGLError("Mesh2D", oglm.Name, "glBufferData vbo in Load: ")
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
	handleOpenGLError("Mesh2D", oglm.Name, "glBindBuffer ibo in Load: ")
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(indicesSize), unsafe.Pointer(&oglm.indices[0]), gl.STATIC_DRAW)
	handleOpenGLError("Mesh2D", oglm.Name, "glBufferData ibo in Load: ")
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	oglm.deleteElements()
}

func (oglm *OpenGLES2Mesh2D) Render() {
	oglm.attributePointer()

	gl.GetError()
	gl.DrawElements(gl.TRIANGLES, int32(oglm.numIndices), gl.UNSIGNED_BYTE, nil)
	handleOpenGLError("Mesh2D", oglm.Name, "RenderError: ")

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (oglm *OpenGLES2Mesh2D) Terminate() {
	var buf [2]uint32
	buf[0] = oglm.vbo
	buf[1] = oglm.ibo

	defer gl.DeleteBuffers(2, buf[:])
}
