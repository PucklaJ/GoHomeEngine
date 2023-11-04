package renderer

import (
	"unsafe"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/go-gl/gl/all-core/gl"
)

type OpenGLMesh2D struct {
	vertices    []gohome.Mesh2DVertex
	indices     []uint32
	numVertices int
	numIndices  int
	Name        string
	vbo         uint32
	ibo         uint32
	vao         uint32
	canUseVAOs  bool
}

func CreateOpenGLMesh2D(name string) *OpenGLMesh2D {
	mesh := OpenGLMesh2D{
		Name: name,
	}

	render, _ := gohome.Render.(*OpenGLRenderer)
	mesh.canUseVAOs = render.HasFunctionAvailable("VERTEX_ARRAY")

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

func (oglm *OpenGLMesh2D) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, gohome.MESH2DVERTEXSIZE, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, gohome.MESH2DVERTEXSIZE, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
}

func (oglm *OpenGLMesh2D) Load() {
	oglm.numVertices = len(oglm.vertices)
	oglm.numIndices = len(oglm.indices)
	verticesSize := oglm.numVertices * gohome.MESH2DVERTEXSIZE
	indicesSize := oglm.numIndices * gohome.INDEXSIZE

	if oglm.canUseVAOs {
		gl.GenVertexArrays(1, &oglm.vao)
	}
	gl.GenBuffers(1, &oglm.vbo)
	gl.GenBuffers(1, &oglm.ibo)

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(verticesSize), unsafe.Pointer(&oglm.vertices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(indicesSize), unsafe.Pointer(&oglm.indices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	if oglm.canUseVAOs {
		gl.BindVertexArray(oglm.vao)
		oglm.attributePointer()
		gl.BindVertexArray(0)
	}

	oglm.deleteElements()

}

func (oglm *OpenGLMesh2D) Render() {
	if oglm.canUseVAOs {
		gl.BindVertexArray(oglm.vao)
	} else {
		oglm.attributePointer()
	}

	gl.GetError()
	gl.DrawElements(gl.TRIANGLES, int32(oglm.numIndices), gl.UNSIGNED_INT, nil)
	handleOpenGLError("Mesh2D", oglm.Name, "RenderError: ")

	if oglm.canUseVAOs {
		gl.BindVertexArray(0)
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	}
}

func (oglm *OpenGLMesh2D) Terminate() {
	if oglm.canUseVAOs {
		defer gl.DeleteVertexArrays(1, &oglm.vao)
	}
	defer gl.DeleteBuffers(1, &oglm.vbo)
	defer gl.DeleteBuffers(1, &oglm.ibo)
}
