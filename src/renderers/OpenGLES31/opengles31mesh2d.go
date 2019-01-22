package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles31"
	"unsafe"
)

type OpenGLES31Mesh2D struct {
	vertices    []gohome.Mesh2DVertex
	indices     []uint32
	numVertices int
	numIndices  int
	Name        string
	vbo         uint32
	ibo         uint32
	vao         uint32
}

func CreateOpenGLES31Mesh2D(name string) *OpenGLES31Mesh2D {
	mesh := OpenGLES31Mesh2D{
		Name: name,
	}

	return &mesh
}

func (oglm *OpenGLES31Mesh2D) deleteElements() {
	oglm.vertices = append(oglm.vertices[:0], oglm.vertices[len(oglm.vertices):]...)
	oglm.indices = append(oglm.indices[:0], oglm.indices[len(oglm.indices):]...)
}

func (oglm *OpenGLES31Mesh2D) AddVertices(vertices []gohome.Mesh2DVertex, indices []uint32) {
	oglm.vertices = append(oglm.vertices, vertices...)
	oglm.indices = append(oglm.indices, indices...)
}

func (oglm *OpenGLES31Mesh2D) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, gl.FALSE, gohome.MESH2DVERTEXSIZE, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, gl.FALSE, gohome.MESH2DVERTEXSIZE, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
}

func (oglm *OpenGLES31Mesh2D) Load() {
	oglm.numVertices = len(oglm.vertices)
	oglm.numIndices = len(oglm.indices)
	verticesSize := oglm.numVertices * gohome.MESH2DVERTEXSIZE
	indicesSize := oglm.numIndices * gohome.INDEXSIZE

	var buf [2]uint32
	gl.GenVertexArrays(1, buf[:])
	oglm.vao = buf[0]
	gl.GenBuffers(2, buf[:])
	oglm.vbo = buf[0]
	oglm.ibo = buf[1]

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(verticesSize), unsafe.Pointer(&oglm.vertices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(indicesSize), unsafe.Pointer(&oglm.indices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.BindVertexArray(oglm.vao)
	oglm.attributePointer()
	gl.BindVertexArray(0)

	oglm.deleteElements()

}

func (oglm *OpenGLES31Mesh2D) Render() {
	gl.BindVertexArray(oglm.vao)

	gl.GetError()
	gl.DrawElements(gl.TRIANGLES, int32(oglm.numIndices), gl.UNSIGNED_INT, nil)
	handleOpenGLES31Error("Mesh2D", oglm.Name, "RenderError: ")

	gl.BindVertexArray(0)
}

func (oglm *OpenGLES31Mesh2D) Terminate() {
	var vbuf [1]uint32
	vbuf[0] = oglm.vao
	gl.DeleteVertexArrays(1, vbuf[:])
	var buf [2]uint32
	buf[0] = oglm.vbo
	buf[1] = oglm.ibo
	gl.DeleteBuffers(2, buf[:])
}
