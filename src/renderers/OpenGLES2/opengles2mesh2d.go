package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles2"
	"unsafe"
)

type OpenGLES2Mesh2D struct {
	vertices    []gohome.Mesh2DVertex
	indices     []uint32
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
	oglm.indices = append(oglm.indices, indices...)
}

func (oglm *OpenGLES2Mesh2D) attributePointer() {
	offset0 := 0
	offset1 := 3 * 4

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, gl.FALSE, int32(gohome.MESH2DVERTEX_SIZE), unsafe.Pointer(&offset0))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, gl.FALSE, int32(gohome.MESH2DVERTEX_SIZE), unsafe.Pointer(&offset1))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
}

func (oglm *OpenGLES2Mesh2D) Load() {
	oglm.numVertices = uint32(len(oglm.vertices))
	oglm.numIndices = uint32(len(oglm.indices))
	var verticesSize uint32 = oglm.numVertices * gohome.MESH2DVERTEX_SIZE
	var indicesSize uint32 = oglm.numIndices * gohome.INDEX_SIZE

	var buf [1]uint32
	gl.GenBuffers(1, buf[:])
	oglm.vbo = buf[0]
	gl.GenBuffers(1, buf[:])
	oglm.ibo = buf[0]

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(verticesSize), unsafe.Pointer(&oglm.vertices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(indicesSize), unsafe.Pointer(&oglm.indices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	oglm.deleteElements()
}

func (oglm *OpenGLES2Mesh2D) Render() {
	oglm.attributePointer()

	gl.GetError()
	gl.DrawElements(gl.TRIANGLES, int32(oglm.numIndices), gl.UNSIGNED_INT, nil)
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
