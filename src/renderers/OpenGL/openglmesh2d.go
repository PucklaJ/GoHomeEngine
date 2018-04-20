package renderer

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/gl/v4.1-core/gl"
	"unsafe"
)

type OpenGLMesh2D struct {
	vertices    []gohome.Mesh2DVertex
	indices     []uint32
	numVertices uint32
	numIndices  uint32
	Name        string
	vbo         uint32
	ibo         uint32
	vao         uint32
}

func CreateOpenGLMesh2D(name string) *OpenGLMesh2D {
	mesh := OpenGLMesh2D{
		Name: name,
	}

	gl.GenVertexArrays(1, &mesh.vao)
	gl.GenBuffers(1, &mesh.vbo)
	gl.GenBuffers(1, &mesh.ibo)

	return &mesh
}

func (oglm *OpenGLMesh2D) deleteElements() {
	oglm.vertices = append(oglm.vertices[:0], oglm.vertices[len(oglm.vertices):]...)
	oglm.indices = append(oglm.indices[:0], oglm.indices[len(oglm.indices):]...)
}

func (oglm *OpenGLMesh2D) AddVertices(vertices []gohome.Mesh2DVertex, indices []uint32) {
	oglm.vertices = append(oglm.vertices, vertices...)
	oglm.indices = append(oglm.indices, indices...)
}

func (oglm *OpenGLMesh2D) Load() {
	oglm.numVertices = uint32(len(oglm.vertices))
	oglm.numIndices = uint32(len(oglm.indices))
	var verticesSize uint32 = oglm.numVertices * gohome.MESH2DVERTEX_SIZE
	var indicesSize uint32 = oglm.numIndices * gohome.INDEX_SIZE

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(verticesSize), unsafe.Pointer(&oglm.vertices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(indicesSize), unsafe.Pointer(&oglm.indices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.BindVertexArray(oglm.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, int32(gohome.MESH2DVERTEX_SIZE), gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, int32(gohome.MESH2DVERTEX_SIZE), gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)

	gl.BindVertexArray(0)

	oglm.deleteElements()

}

func (oglm *OpenGLMesh2D) Render() {
	gl.BindVertexArray(oglm.vao)

	gl.DrawElements(gl.TRIANGLES, int32(oglm.numIndices), gl.UNSIGNED_INT, nil)

	gl.BindVertexArray(0)
}

func (oglm *OpenGLMesh2D) Terminate() {
	defer gl.DeleteVertexArrays(1, &oglm.vao)
	defer gl.DeleteBuffers(1, &oglm.vbo)
	defer gl.DeleteBuffers(1, &oglm.ibo)
}