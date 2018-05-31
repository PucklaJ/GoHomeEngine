package renderer

import (
	// "fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"sync"
	"unsafe"
)

const (
	NUM_GO_ROUTINES_TANGENTS_CALCULATING uint32 = 10
	MESH3DVERTEX_SIZE                    uint32 = 3*4 + 3*4 + 2*4 + 3*4 // 3*sizeof(float32)+3*sizeof(float32)+2*sizeof(float32)+3*sizeof(float32)
)

type OpenGLMesh3D struct {
	vertices    []gohome.Mesh3DVertex
	indices     []uint32
	numVertices uint32
	numIndices  uint32

	vao    uint32
	buffer uint32

	Name     string
	Material *gohome.Material

	tangentsCalculated bool
	canUseVAOs         bool
	hasUV 			   bool

	aabb gohome.AxisAlignedBoundingBox
}

func (oglm *OpenGLMesh3D) CalculateTangentsRoutine(startIndex, maxIndex uint32, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	vertices := &oglm.vertices
	indices := oglm.indices

	var p0, p1, p2 mgl32.Vec3
	var t0, t1, t2 mgl32.Vec2
	var r float32
	var deltaPos1, deltaPos2 mgl32.Vec3
	var deltaUv1, deltaUv2 mgl32.Vec2
	var tangent mgl32.Vec3
	var normal mgl32.Vec3
	var bitangent mgl32.Vec3
	var i uint32
	for i = startIndex; i < maxIndex && i < uint32(len(indices)); i += 3 {
		if i > uint32(len(indices)-3) {
			break
		}

		p0 = mgl32.Vec3{(*vertices)[indices[i]][0], (*vertices)[indices[i]][1], (*vertices)[indices[i]][2]}
		p1 = mgl32.Vec3{(*vertices)[indices[i+1]][0], (*vertices)[indices[i+1]][1], (*vertices)[indices[i+1]][2]}
		p2 = mgl32.Vec3{(*vertices)[indices[i+2]][0], (*vertices)[indices[i+2]][1], (*vertices)[indices[i+2]][2]}

		t0 = mgl32.Vec2{(*vertices)[indices[i]][6], (*vertices)[indices[i]][7]}
		t1 = mgl32.Vec2{(*vertices)[indices[i+1]][6], (*vertices)[indices[i+1]][7]}
		t2 = mgl32.Vec2{(*vertices)[indices[i+2]][6], (*vertices)[indices[i+2]][7]}

		if t0.X() == 0.0 && t0.Y() == 0.0 && t1.X() == 0.0 && t1.Y() == 0.0 && t2.X() == 0.0 && t2.Y() == 0.0 {
			oglm.hasUV = false
			continue
		}

		normal = mgl32.Vec3{(*vertices)[indices[i]][3], (*vertices)[indices[i]][4], (*vertices)[indices[i]][5]}

		deltaPos1 = p1.Sub(p0)
		deltaPos2 = p2.Sub(p0)

		deltaUv1 = t1.Sub(t0)
		deltaUv2 = t2.Sub(t0)

		r = 1.0 / (deltaUv1[0]*deltaUv2[1] - deltaUv1[1]*deltaUv2[0])

		tangent = (deltaPos1.Mul(deltaUv2[1]).Sub(deltaPos2.Mul(deltaUv1[1]))).Mul(r).Normalize()
		tangent = tangent.Sub(normal.Mul(normal.Dot(tangent))).Normalize()
		bitangent = (deltaPos2.Mul(deltaUv1[0]).Sub(deltaPos1.Mul(deltaUv2[0]))).Mul(r).Normalize()
		if normal.Cross(tangent).Dot(bitangent) < 0.0 {
			tangent = tangent.Mul(-1.0)
		}
		var j uint32
		for j = 0; j < 3; j++ {
			(*vertices)[indices[i+j]][8] = tangent[0]
			(*vertices)[indices[i+j]][9] = tangent[1]
			(*vertices)[indices[i+j]][10] = tangent[2]
		}
	}
}

func (oglm *OpenGLMesh3D) CalculateTangents() {
	if oglm.tangentsCalculated {
		return
	}
	var wg sync.WaitGroup

	deltaIndex := uint32(len(oglm.indices)) / NUM_GO_ROUTINES_TANGENTS_CALCULATING
	if deltaIndex == 0 {
		deltaIndex = uint32(len(oglm.indices)) / 3
	}
	if deltaIndex > 3 {
		deltaIndex -= deltaIndex % 3
	} else {
		deltaIndex = 3
	}
	oglm.hasUV = true
	var i uint32
	for i = 0; i < NUM_GO_ROUTINES_TANGENTS_CALCULATING*2; i++ {
		wg.Add(1)
		go oglm.CalculateTangentsRoutine(i*deltaIndex, i*deltaIndex+deltaIndex, &wg)
		if i*deltaIndex+deltaIndex >= uint32(len(oglm.indices)) {
			break
		}
	}

	wg.Wait()

	oglm.tangentsCalculated = true
}

func (oglm *OpenGLMesh3D) AddVertices(vertices []gohome.Mesh3DVertex, indices []uint32) {
	oglm.vertices = append(oglm.vertices, vertices...)
	oglm.indices = append(oglm.indices, indices...)
	oglm.checkAABB()
}

func (oglm *OpenGLMesh3D) checkAABB() {
	var max,min mgl32.Vec3 = [3]float32{oglm.vertices[0][0],oglm.vertices[0][1],oglm.vertices[0][2]}, [3]float32{oglm.vertices[0][0],oglm.vertices[0][1],oglm.vertices[0][2]}
	var current gohome.Mesh3DVertex
	for i:=0;i<len(oglm.vertices);i++ {
		current = oglm.vertices[i]
		for j:=0;j<3;j++ {
			if current[j] > max[j] {
				max[j] = current[j]
			} else if current[j] < min[j] {
				min[j] = current[j]
			}
		}
	}

	for i:=0;i<3;i++ {
		if max[i] > oglm.aabb.Max[i] {
			oglm.aabb.Max[i] = max[i]
		}
		if min[i] < oglm.aabb.Min[i] {
			oglm.aabb.Min[i] = min[i]
		}
	}
}

func CreateOpenGLMesh3D(name string) *OpenGLMesh3D {
	mesh := OpenGLMesh3D{
		Name:               name,
		tangentsCalculated: false,
	}
	render, _ := gohome.Render.(*OpenGLRenderer)
	mesh.canUseVAOs = render.hasFunctionAvailable("VERTEX_ARRAY")

	return &mesh
}

func (oglm *OpenGLMesh3D) deleteElements() {
	oglm.vertices = append(oglm.vertices[:0], oglm.vertices[len(oglm.vertices):]...)
	oglm.indices = append(oglm.indices[:0], oglm.indices[len(oglm.indices):]...)
}

func (oglm *OpenGLMesh3D) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.buffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(MESH3DVERTEX_SIZE), gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(MESH3DVERTEX_SIZE), gl.PtrOffset(3*4))
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(MESH3DVERTEX_SIZE), gl.PtrOffset(3*4+3*4))
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(MESH3DVERTEX_SIZE), gl.PtrOffset(3*4+3*4+2*4))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.buffer)
}

func (oglm *OpenGLMesh3D) Load() {

	oglm.numVertices = uint32(len(oglm.vertices))
	oglm.numIndices = uint32(len(oglm.indices))

	if oglm.numVertices == 0 || oglm.numIndices == 0 {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Mesh3D", oglm.Name, "No vertices or indices have been added!")
		return
	}

	var verticesSize uint32 = oglm.numVertices * MESH3DVERTEX_SIZE
	var indicesSize uint32 = oglm.numIndices * gohome.INDEX_SIZE

	oglm.CalculateTangents()

	if oglm.canUseVAOs {
		gl.GenVertexArrays(1, &oglm.vao)
	}
	gl.GenBuffers(1, &oglm.buffer)

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.buffer)
	gl.BufferData(gl.ARRAY_BUFFER, int(verticesSize)+int(indicesSize), nil, gl.STATIC_DRAW)

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(verticesSize), unsafe.Pointer(&oglm.vertices[0]))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.buffer)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, int(verticesSize), int(indicesSize), unsafe.Pointer(&oglm.indices[0]))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	if oglm.canUseVAOs {
		gl.BindVertexArray(oglm.vao)
		oglm.attributePointer()
		gl.BindVertexArray(0)
	}

	oglm.deleteElements()
}

func (oglm *OpenGLMesh3D) Render() {
	if oglm.numVertices == 0 || oglm.numIndices == 0 {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Mesh", oglm.Name, "No Vertices or Indices have been loaded!")
		return
	}
	if gohome.RenderMgr.CurrentShader != nil {
		gohome.RenderMgr.CurrentShader.SetUniformMaterial(*oglm.Material)
	}
	if oglm.canUseVAOs {
		gl.BindVertexArray(oglm.vao)
	} else {
		oglm.attributePointer()
	}
	gl.GetError()
	gl.DrawElements(gl.TRIANGLES, int32(oglm.numIndices), gl.UNSIGNED_INT, gl.PtrOffset(int(oglm.numVertices*MESH3DVERTEX_SIZE)))
	handleOpenGLError("Mesh3D", oglm.Name, "RenderError: ")
	if oglm.canUseVAOs {
		gl.BindVertexArray(0)
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	}
}

func (oglm *OpenGLMesh3D) Terminate() {
	if oglm.canUseVAOs {
		defer gl.DeleteVertexArrays(1, &oglm.vao)
	}
	defer gl.DeleteBuffers(1, &oglm.buffer)
}

func (oglm *OpenGLMesh3D) SetMaterial(mat *gohome.Material) {
	oglm.Material = mat
}

func (oglm *OpenGLMesh3D) GetMaterial() *gohome.Material {
	if oglm.Material == nil {
		oglm.Material = &gohome.Material{}
	}
	return oglm.Material
}

func (oglm *OpenGLMesh3D) GetNumVertices() uint32 {
	return oglm.numVertices
}
func (oglm *OpenGLMesh3D) GetNumIndices() uint32 {
	return oglm.numIndices
}

func (oglm *OpenGLMesh3D) GetVertices() []gohome.Mesh3DVertex {
	return oglm.vertices
}
func (oglm *OpenGLMesh3D) GetIndices() []uint32 {
	return oglm.indices
}

func (oglm *OpenGLMesh3D) GetName() string {
	return oglm.Name
}

func (oglm *OpenGLMesh3D) AABB() gohome.AxisAlignedBoundingBox {
	return oglm.aabb
}

func (oglm *OpenGLMesh3D) HasUV() bool {
	return oglm.hasUV
}

