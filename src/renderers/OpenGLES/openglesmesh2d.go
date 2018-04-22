package renderer

import (
	"encoding/binary"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/gl"
	// "log"
)

type OpenGLESMesh2D struct {
	vertices    []gohome.Mesh2DVertex
	indices     []uint32
	numVertices uint32
	numIndices  uint32
	Name        string
	vbo         gl.Buffer
	ibo         gl.Buffer
	vao         gl.VertexArray // OpenGLES 3.0 only
	gles        *gl.Context
	isgles3     bool
}

func CreateOpenGLESMesh2D(name string) *OpenGLESMesh2D {
	mesh := OpenGLESMesh2D{
		Name:    name,
		isgles3: false,
	}

	render, _ := gohome.Render.(*OpenGLESRenderer)
	mesh.gles = &render.gles
	_, mesh.isgles3 = (*mesh.gles).(gl.Context3)

	return &mesh
}

func (oglm *OpenGLESMesh2D) deleteElements() {
	oglm.vertices = append(oglm.vertices[:0], oglm.vertices[len(oglm.vertices):]...)
	oglm.indices = append(oglm.indices[:0], oglm.indices[len(oglm.indices):]...)
}

func (oglm *OpenGLESMesh2D) AddVertices(vertices []gohome.Mesh2DVertex, indices []uint32) {
	oglm.vertices = append(oglm.vertices, vertices...)
	oglm.indices = append(oglm.indices, indices...)
}

func (oglm *OpenGLESMesh2D) attributePointer() {
	(*oglm.gles).BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	(*oglm.gles).VertexAttribPointer(gl.Attrib{0}, 2, gl.FLOAT, false, int(gohome.MESH2DVERTEX_SIZE), 0)
	(*oglm.gles).VertexAttribPointer(gl.Attrib{1}, 2, gl.FLOAT, false, int(gohome.MESH2DVERTEX_SIZE), 2*4)
	(*oglm.gles).EnableVertexAttribArray(gl.Attrib{0})
	(*oglm.gles).EnableVertexAttribArray(gl.Attrib{1})

	(*oglm.gles).BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
}

func (oglm *OpenGLESMesh2D) Load() {
	oglm.numVertices = uint32(len(oglm.vertices))
	oglm.numIndices = uint32(len(oglm.indices))
	// var verticesSize uint32 = oglm.numVertices * gohome.MESH2DVERTEX_SIZE
	// var indicesSize uint32 = oglm.numIndices * gohome.INDEX_SIZE

	var verticesFloats []float32 = make([]float32, oglm.numVertices*4)
	var index uint32 = 0
	for i := 0; i < int(oglm.numVertices); i++ {
		verticesFloats[index+0] = oglm.vertices[i][0]
		verticesFloats[index+1] = oglm.vertices[i][1]
		verticesFloats[index+2] = oglm.vertices[i][2]
		verticesFloats[index+3] = oglm.vertices[i][3]
		index += 4
	}
	indicesBytes := make([]byte, oglm.numIndices*gohome.INDEX_SIZE)
	for i := 0; i < int(oglm.numIndices); i++ {
		binary.LittleEndian.PutUint32(indicesBytes[uint32(i)*gohome.INDEX_SIZE:uint32(i)*gohome.INDEX_SIZE+3+1], oglm.indices[i])
	}

	if oglm.isgles3 {
		oglm.vao = (*oglm.gles).CreateVertexArray()
	}
	oglm.vbo = (*oglm.gles).CreateBuffer()
	oglm.ibo = (*oglm.gles).CreateBuffer()

	(*oglm.gles).BindBuffer(gl.ARRAY_BUFFER, oglm.vbo)
	(*oglm.gles).BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, verticesFloats...), gl.STATIC_DRAW)
	(*oglm.gles).BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{0})

	(*oglm.gles).BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.ibo)
	(*oglm.gles).BufferData(gl.ELEMENT_ARRAY_BUFFER, indicesBytes, gl.STATIC_DRAW)
	(*oglm.gles).BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{0})

	if oglm.isgles3 {
		(*oglm.gles).BindVertexArray(oglm.vao)
		oglm.attributePointer()
		(*oglm.gles).BindVertexArray(gl.VertexArray{0})
	}

	oglm.deleteElements()

}

func (oglm *OpenGLESMesh2D) Render() {
	if oglm.isgles3 {
		(*oglm.gles).BindVertexArray(oglm.vao)
	} else {
		oglm.attributePointer()
	}

	(*oglm.gles).DrawElements(gl.TRIANGLES, int(oglm.numIndices), gl.UNSIGNED_INT, 0)

	if oglm.isgles3 {
		(*oglm.gles).BindVertexArray(gl.VertexArray{0})
	} else {
		(*oglm.gles).BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{0})
		(*oglm.gles).BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{0})
	}
}

func (oglm *OpenGLESMesh2D) Terminate() {
	if oglm.isgles3 {
		defer (*oglm.gles).DeleteVertexArray(oglm.vao)
	}
	defer (*oglm.gles).DeleteBuffer(oglm.vbo)
	defer (*oglm.gles).DeleteBuffer(oglm.ibo)
}
