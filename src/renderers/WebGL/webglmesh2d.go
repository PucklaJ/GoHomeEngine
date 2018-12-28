package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
)

type WebGLMesh2D struct {
	vertices    []gohome.Mesh2DVertex
	indices     []uint8
	numVertices uint32
	numIndices  uint32
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
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, int(gohome.MESH2DVERTEX_SIZE), 0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, int(gohome.MESH2DVERTEX_SIZE), 2*4)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
}

func (oglm *WebGLMesh2D) Load() {
	oglm.numVertices = uint32(len(oglm.vertices))
	oglm.numIndices = uint32(len(oglm.indices))
	var verticesSize uint32 = oglm.numVertices * gohome.MESH2DVERTEX_SIZE
	var indicesSize uint32 = oglm.numIndices * 1

	oglm.vbo = gl.CreateBuffer()
	oglm.ibo = gl.CreateBuffer()

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(verticesSize), oglm.vertices, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(indicesSize), oglm.indices, gl.STATIC_DRAW)
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
