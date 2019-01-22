package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
)

type WebGLMesh2D struct {
	vertices    []gohome.Mesh2DVertex
	indices     []uint8
	numVertices int
	numIndices  int
	Name        string
	vbo         *js.Object
	ibo         *js.Object
}

func CreateWebGLMesh2D(name string) *WebGLMesh2D {
	mesh := WebGLMesh2D{
		Name: name,
	}

	return &mesh
}

func (oglm *WebGLMesh2D) deleteElements() {
	oglm.vertices = append(oglm.vertices[:0], oglm.vertices[len(oglm.vertices):]...)
	oglm.indices = append(oglm.indices[:0], oglm.indices[len(oglm.indices):]...)
}

func (oglm *WebGLMesh2D) AddVertices(vertices []gohome.Mesh2DVertex, indices []uint32) {
	oglm.vertices = append(oglm.vertices, vertices...)
	index := len(oglm.indices)
	oglm.indices = append(oglm.indices, make([]uint8, len(indices))...)
	for id, i := range indices {
		oglm.indices[index+id] = uint8(i)
	}
}

func (oglm *WebGLMesh2D) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, gohome.MESH2DVERTEXSIZE, 0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, gohome.MESH2DVERTEXSIZE, 2*4)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
}

func (oglm *WebGLMesh2D) Load() {
	oglm.numVertices = len(oglm.vertices)
	oglm.numIndices = len(oglm.indices)
	verticesSize := oglm.numVertices * gohome.MESH2DVERTEXSIZE
	indicesSize := oglm.numIndices * 1

	floatBuffer := gohome.Mesh2DVerticesToFloatArray(oglm.vertices)

	oglm.vbo = gl.CreateBuffer()
	oglm.ibo = gl.CreateBuffer()

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.GetError()
	gl.BufferData(gl.ARRAY_BUFFER, verticesSize, floatBuffer, gl.STATIC_DRAW)
	handleWebGLError("Mesh2D", oglm.Name, "glBufferData VBO: ")
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indicesSize, oglm.indices, gl.STATIC_DRAW)
	handleWebGLError("Mesh2D", oglm.Name, "glBufferData IBO: ")
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, nil)

	oglm.deleteElements()
}

func (oglm *WebGLMesh2D) Render() {
	oglm.attributePointer()

	gl.GetError()
	gl.DrawElements(gl.TRIANGLES, int(oglm.numIndices), gl.UNSIGNED_BYTE, 0)
	handleWebGLError("Mesh2D", oglm.Name, "RenderError: ")

	gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, nil)
}

func (oglm *WebGLMesh2D) Terminate() {
	gl.DeleteBuffer(oglm.vbo)
	gl.DeleteBuffer(oglm.ibo)
}
