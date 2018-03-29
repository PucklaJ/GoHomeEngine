package gohome

import (
	// "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/assimp"
	"log"
	"sync"
	"unsafe"
)

const (
	NUM_GO_ROUTINES_MESH_VERTICES_LOADING uint32 = 10
	NUM_GO_ROUTINES_TANGENTS_CALCULATING  uint32 = 10
	MESH3DVERTEX_SIZE                     uint32 = 3*4 + 3*4 + 2*4 + 3*4 // 3*sizeof(float32)+3*sizeof(float32)+2*sizeof(float32)+3*sizeof(float32)
)

type Mesh3DVertex [3 + 3 + 2 + 3]float32

type Mesh3D interface {
	AddVerticesAssimp(mesh *assimp.Mesh, node *assimp.Node, scene *assimp.Scene, level *Level, directory string, preloaded bool)
	AddVertices(vertices []Mesh3DVertex, indices []uint32)
	Load()
	Render()
	Terminate()
	SetMaterial(mat *Material)
	GetMaterial() *Material
	GetName() string
	GetNumVertices() uint32
	GetNumIndices() uint32
	calculateTangents()
}

type OpenGLMesh3D struct {
	vertices    []Mesh3DVertex
	indices     []uint32
	numVertices uint32
	numIndices  uint32

	vao    uint32
	buffer uint32

	Name     string
	Material *Material

	tangentsCalculated bool
}

func loadVertices(vertices *[]Mesh3DVertex, mesh *assimp.Mesh, start_index, end_index, max_index uint32, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start_index; i < uint32(mesh.NumVertices()) && i < end_index; i++ {
		var texCoords mgl32.Vec2
		if len(mesh.TextureCoords(0)) > 0 {
			texCoords[0] = mesh.TextureCoords(0)[i].X()
			texCoords[1] = mesh.TextureCoords(0)[i].Y()
		} else {
			texCoords[0] = 0.0
			texCoords[1] = 0.0
		}
		vertex := Mesh3DVertex{
			/* X,Y,Z,
			   NX,NY,NZ,
			   U,V,
			   TX,TY,TZ,
			*/
			mesh.Vertices()[i].X(), mesh.Vertices()[i].Y(), mesh.Vertices()[i].Z(),
			mesh.Normals()[i].X(), mesh.Normals()[i].Y(), mesh.Vertices()[i].Z(),
			texCoords[0], texCoords[1],
			0.0, 0.0, 0.0,
		}
		(*vertices)[i] = vertex
	}
}

func (oglm *OpenGLMesh3D) calculateTangentsRoutine(startIndex, maxIndex uint32, wg *sync.WaitGroup) {
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

func (oglm *OpenGLMesh3D) calculateTangents() {
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

	var i uint32
	for i = 0; i < NUM_GO_ROUTINES_TANGENTS_CALCULATING*2; i++ {
		wg.Add(1)
		go oglm.calculateTangentsRoutine(i*deltaIndex, i*deltaIndex+deltaIndex, &wg)
		if i*deltaIndex+deltaIndex >= uint32(len(oglm.indices)) {
			break
		}
	}

	wg.Wait()

	oglm.tangentsCalculated = true
}

func (oglm *OpenGLMesh3D) AddVerticesAssimp(mesh *assimp.Mesh, node *assimp.Node, scene *assimp.Scene, level *Level, directory string, preloaded bool) {
	vertices := make([]Mesh3DVertex, mesh.NumVertices())
	var indices []uint32
	var wg sync.WaitGroup
	var i uint32
	deltaIndex := uint32(mesh.NumVertices()) / NUM_GO_ROUTINES_MESH_VERTICES_LOADING
	wg.Add(int(NUM_GO_ROUTINES_MESH_VERTICES_LOADING))
	for i = 0; i < NUM_GO_ROUTINES_MESH_VERTICES_LOADING; i++ {
		go loadVertices(&vertices, mesh, i*deltaIndex, (i+1)*deltaIndex, uint32(mesh.NumVertices()), &wg)
	}

	for i = 0; i < uint32(mesh.NumFaces()); i++ {
		face := mesh.Faces()[i]
		faceIndices := face.CopyIndices()
		indices = append(indices, faceIndices...)
	}

	wg.Wait()

	oglm.AddVertices(vertices, indices)

	mat := &Material{}
	mat.Init(scene.Materials()[mesh.MaterialIndex()], scene, directory, preloaded)
	oglm.Material = mat
}

func (oglm *OpenGLMesh3D) AddVertices(vertices []Mesh3DVertex, indices []uint32) {
	oglm.vertices = append(oglm.vertices, vertices...)
	oglm.indices = append(oglm.indices, indices...)
}

func CreateOpenGLMesh3D(name string) *OpenGLMesh3D {
	mesh := OpenGLMesh3D{
		Name:               name,
		tangentsCalculated: false,
	}

	return &mesh
}

func (oglm *OpenGLMesh3D) deleteElements() {
	oglm.vertices = append(oglm.vertices[:0], oglm.vertices[len(oglm.vertices):]...)
	oglm.indices = append(oglm.indices[:0], oglm.indices[len(oglm.indices):]...)
}

func (oglm *OpenGLMesh3D) Load() {

	oglm.numVertices = uint32(len(oglm.vertices))
	oglm.numIndices = uint32(len(oglm.indices))

	if oglm.numVertices == 0 || oglm.numIndices == 0 {
		log.Println("No vertices or indices have been added for mesh", oglm.Name, "!")
		return
	}

	var verticesSize uint32 = oglm.numVertices * MESH3DVERTEX_SIZE
	var indicesSize uint32 = oglm.numIndices * INDEX_SIZE

	oglm.calculateTangents()

	gl.GenVertexArrays(1, &oglm.vao)
	gl.GenBuffers(1, &oglm.buffer)

	gl.BindBuffer(gl.ARRAY_BUFFER, oglm.buffer)
	gl.BufferData(gl.ARRAY_BUFFER, int(verticesSize)+int(indicesSize), nil, gl.STATIC_DRAW)

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(verticesSize), unsafe.Pointer(&oglm.vertices[0]))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglm.buffer)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, int(verticesSize), int(indicesSize), unsafe.Pointer(&oglm.indices[0]))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.BindVertexArray(oglm.vao)

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

	gl.BindVertexArray(0)

	oglm.deleteElements()
}

func (oglm *OpenGLMesh3D) Render() {
	if oglm.numVertices == 0 || oglm.numIndices == 0 {
		return
	}
	if RenderMgr.currentShader != nil {
		if err := RenderMgr.currentShader.SetUniformMaterial(*oglm.Material); err != nil {
			// fmt.Println("Error:", err)
		}
	}
	gl.BindVertexArray(oglm.vao)

	gl.DrawElements(gl.TRIANGLES, int32(oglm.numIndices), gl.UNSIGNED_INT, gl.PtrOffset(int(oglm.numVertices)*int(MESH3DVERTEX_SIZE)))

	gl.BindVertexArray(0)
}

func (oglm *OpenGLMesh3D) Terminate() {
	defer gl.DeleteVertexArrays(1, &oglm.vao)
	defer gl.DeleteBuffers(1, &oglm.buffer)
}

func VertexPosIndex(which int) int {
	return which
}

func VertexNormalIndex(which int) int {
	return 3 + which
}

func VertexTexCoordIndex(which int) int {
	return 2*3 + which
}

func (oglm *OpenGLMesh3D) SetMaterial(mat *Material) {
	oglm.Material = mat
}

func (oglm *OpenGLMesh3D) GetMaterial() *Material {
	if oglm.Material == nil {
		oglm.Material = &Material{}
	}
	return oglm.Material
}

func (oglm *OpenGLMesh3D) GetNumVertices() uint32 {
	return oglm.numVertices
}
func (oglm *OpenGLMesh3D) GetNumIndices() uint32 {
	return oglm.numIndices
}

func Box(name string, size mgl32.Vec3) Mesh3D {

	boxMesh := Render.CreateMesh3D(name)

	vertices := make([]Mesh3DVertex, 24)
	indices := make([]uint32, 36)

	/***** Positions ****/

	// BOTTOM_LEFT_FRONT
	vertices[0][VertexPosIndex(0)] = -size[0] / 2.0
	vertices[0][VertexPosIndex(1)] = -size[1] / 2.0
	vertices[0][VertexPosIndex(2)] = size[2] / 2.0

	for i := 0; i < 3; i++ {
		vertices[8][VertexPosIndex(i)] = vertices[0][VertexPosIndex(i)]
		vertices[20][VertexPosIndex(i)] = vertices[0][VertexPosIndex(i)]
	}

	// BOTTOM_RIGHT_FRONT
	vertices[1][VertexPosIndex(0)] = size[0] / 2.0
	vertices[1][VertexPosIndex(1)] = -size[1] / 2.0
	vertices[1][VertexPosIndex(2)] = size[2] / 2.0

	for i := 0; i < 3; i++ {
		vertices[9][VertexPosIndex(i)] = vertices[1][VertexPosIndex(i)]
		vertices[12][VertexPosIndex(i)] = vertices[1][VertexPosIndex(i)]
	}

	// BOTTOM_RIGHT_BACK
	vertices[2][VertexPosIndex(0)] = size[0] / 2.0
	vertices[2][VertexPosIndex(1)] = -size[1] / 2.0
	vertices[2][VertexPosIndex(2)] = -size[2] / 2.0

	for i := 0; i < 3; i++ {
		vertices[13][VertexPosIndex(i)] = vertices[2][VertexPosIndex(i)]
		vertices[17][VertexPosIndex(i)] = vertices[2][VertexPosIndex(i)]
	}

	// BOTTOM_LEFT_BACK
	vertices[3][VertexPosIndex(0)] = -size[0] / 2.0
	vertices[3][VertexPosIndex(1)] = -size[1] / 2.0
	vertices[3][VertexPosIndex(2)] = -size[2] / 2.0

	for i := 0; i < 3; i++ {
		vertices[16][VertexPosIndex(i)] = vertices[3][VertexPosIndex(i)]
		vertices[21][VertexPosIndex(i)] = vertices[3][VertexPosIndex(i)]
	}

	// TOP_LEFT_FRONT
	vertices[4][VertexPosIndex(0)] = -size[0] / 2.0
	vertices[4][VertexPosIndex(1)] = size[1] / 2.0
	vertices[4][VertexPosIndex(2)] = size[2] / 2.0

	for i := 0; i < 3; i++ {
		vertices[11][VertexPosIndex(i)] = vertices[4][VertexPosIndex(i)]
		vertices[23][VertexPosIndex(i)] = vertices[4][VertexPosIndex(i)]
	}

	// TOP_RIGHT_FRONT
	vertices[5][VertexPosIndex(0)] = size[0] / 2.0
	vertices[5][VertexPosIndex(1)] = size[1] / 2.0
	vertices[5][VertexPosIndex(2)] = size[2] / 2.0

	for i := 0; i < 3; i++ {
		vertices[10][VertexPosIndex(i)] = vertices[5][VertexPosIndex(i)]
		vertices[15][VertexPosIndex(i)] = vertices[5][VertexPosIndex(i)]
	}

	// TOP_RIGHT_BACK
	vertices[6][VertexPosIndex(0)] = size[0] / 2.0
	vertices[6][VertexPosIndex(1)] = size[1] / 2.0
	vertices[6][VertexPosIndex(2)] = -size[2] / 2.0

	for i := 0; i < 3; i++ {
		vertices[14][VertexPosIndex(i)] = vertices[6][VertexPosIndex(i)]
		vertices[18][VertexPosIndex(i)] = vertices[6][VertexPosIndex(i)]
	}

	// TOP_LEFT_BACK
	vertices[7][VertexPosIndex(0)] = -size[0] / 2.0
	vertices[7][VertexPosIndex(1)] = size[1] / 2.0
	vertices[7][VertexPosIndex(2)] = -size[2] / 2.0

	for i := 0; i < 3; i++ {
		vertices[19][VertexPosIndex(i)] = vertices[7][VertexPosIndex(i)]
		vertices[22][VertexPosIndex(i)] = vertices[7][VertexPosIndex(i)]
	}

	/****** Indices ******/

	// BOTTOM
	indices[0] = 0
	indices[1] = 3
	indices[2] = 2
	indices[3] = 2
	indices[4] = 1
	indices[5] = 0

	// TOP
	indices[6] = 4
	indices[7] = 5
	indices[8] = 6
	indices[9] = 6
	indices[10] = 7
	indices[11] = 4

	// FRONT
	indices[12] = 8
	indices[13] = 9
	indices[14] = 10
	indices[15] = 10
	indices[16] = 11
	indices[17] = 8

	// RIGHT
	indices[18] = 12
	indices[19] = 13
	indices[20] = 14
	indices[21] = 14
	indices[22] = 15
	indices[23] = 12

	// BACK
	indices[24] = 17
	indices[25] = 16
	indices[26] = 19
	indices[27] = 19
	indices[28] = 18
	indices[29] = 17

	// LEFT
	indices[30] = 21
	indices[31] = 20
	indices[32] = 23
	indices[33] = 23
	indices[34] = 22
	indices[35] = 21

	/****** Normals ******/

	// BOTTOM
	for i := 0; i < 4; i++ {
		vertices[i][VertexNormalIndex(0)] = 0.0
		vertices[i][VertexNormalIndex(1)] = -1.0
		vertices[i][VertexNormalIndex(2)] = 0.0
	}

	// TOP
	for i := 4; i < 8; i++ {
		vertices[i][VertexNormalIndex(0)] = 0.0
		vertices[i][VertexNormalIndex(1)] = 1.0
		vertices[i][VertexNormalIndex(2)] = 0.0
	}

	// FRONT
	for i := 12; i < 18; i++ {
		vertices[indices[i]][VertexNormalIndex(0)] = 0.0
		vertices[indices[i]][VertexNormalIndex(1)] = 0.0
		vertices[indices[i]][VertexNormalIndex(2)] = 1.0
	}

	// RIGHT
	for i := 18; i < 24; i++ {
		vertices[indices[i]][VertexNormalIndex(0)] = 1.0
		vertices[indices[i]][VertexNormalIndex(1)] = 0.0
		vertices[indices[i]][VertexNormalIndex(2)] = 0.0
	}

	// BACK
	for i := 24; i < 30; i++ {
		vertices[indices[i]][VertexNormalIndex(0)] = 0.0
		vertices[indices[i]][VertexNormalIndex(1)] = 0.0
		vertices[indices[i]][VertexNormalIndex(2)] = -1.0
	}

	// LEFT
	for i := 30; i < 36; i++ {
		vertices[indices[i]][VertexNormalIndex(0)] = -1.0
		vertices[indices[i]][VertexNormalIndex(1)] = 0.0
		vertices[indices[i]][VertexNormalIndex(2)] = 0.0
	}

	/****** UV *****/

	// BOT
	vertices[0][VertexTexCoordIndex(0)] = 0.0
	vertices[0][VertexTexCoordIndex(1)] = 1.0
	vertices[1][VertexTexCoordIndex(0)] = 1.0 / 3.0
	vertices[1][VertexTexCoordIndex(1)] = 1.0
	vertices[2][VertexTexCoordIndex(0)] = 1.0 / 3.0
	vertices[2][VertexTexCoordIndex(1)] = 0.5
	vertices[3][VertexTexCoordIndex(0)] = 0.0
	vertices[3][VertexTexCoordIndex(1)] = 0.5

	// TOP
	vertices[4][VertexTexCoordIndex(0)] = 1.0 / 3.0
	vertices[4][VertexTexCoordIndex(1)] = 1.0
	vertices[5][VertexTexCoordIndex(0)] = 2.0 / 3.0
	vertices[5][VertexTexCoordIndex(1)] = 1.0
	vertices[6][VertexTexCoordIndex(0)] = 2.0 / 3.0
	vertices[6][VertexTexCoordIndex(1)] = 0.5
	vertices[7][VertexTexCoordIndex(0)] = 1.0 / 3.0
	vertices[7][VertexTexCoordIndex(1)] = 0.5

	// FRONT
	vertices[8][VertexTexCoordIndex(0)] = 2.0 / 3.0
	vertices[8][VertexTexCoordIndex(1)] = 1.0
	vertices[9][VertexTexCoordIndex(0)] = 1.0
	vertices[9][VertexTexCoordIndex(1)] = 1.0
	vertices[10][VertexTexCoordIndex(0)] = 1.0
	vertices[10][VertexTexCoordIndex(1)] = 0.5
	vertices[11][VertexTexCoordIndex(0)] = 2.0 / 3.0
	vertices[11][VertexTexCoordIndex(1)] = 0.5

	// RIGHT
	vertices[12][VertexTexCoordIndex(0)] = 2.0 / 3.0
	vertices[12][VertexTexCoordIndex(1)] = 0.5
	vertices[13][VertexTexCoordIndex(0)] = 1.0
	vertices[13][VertexTexCoordIndex(1)] = 0.5
	vertices[14][VertexTexCoordIndex(0)] = 1.0
	vertices[14][VertexTexCoordIndex(1)] = 0.0
	vertices[15][VertexTexCoordIndex(0)] = 2.0 / 3.0
	vertices[15][VertexTexCoordIndex(1)] = 0.0

	// BACK
	vertices[16][VertexTexCoordIndex(0)] = 2.0 / 3.0
	vertices[16][VertexTexCoordIndex(1)] = 0.5
	vertices[17][VertexTexCoordIndex(0)] = 1.0 / 3.0
	vertices[17][VertexTexCoordIndex(1)] = 0.5
	vertices[18][VertexTexCoordIndex(0)] = 1.0 / 3.0
	vertices[18][VertexTexCoordIndex(1)] = 0.0
	vertices[19][VertexTexCoordIndex(0)] = 2.0 / 3.0
	vertices[19][VertexTexCoordIndex(1)] = 0.0

	// LEFT
	vertices[20][VertexTexCoordIndex(0)] = 1.0 / 3.0
	vertices[20][VertexTexCoordIndex(1)] = 0.5
	vertices[21][VertexTexCoordIndex(0)] = 0.0
	vertices[21][VertexTexCoordIndex(1)] = 0.5
	vertices[22][VertexTexCoordIndex(0)] = 0.0
	vertices[22][VertexTexCoordIndex(1)] = 0.0
	vertices[23][VertexTexCoordIndex(0)] = 1.0 / 3.0
	vertices[23][VertexTexCoordIndex(1)] = 0.0

	boxMesh.AddVertices(vertices, indices)
	boxMesh.Load()
	mat := &Material{}
	mat.InitDefault()
	boxMesh.SetMaterial(mat)

	return boxMesh
}

func Plane(name string, size mgl32.Vec2, textures uint32) Mesh3D {
	// xAxis := up.Cross([3]float32{0.0, 0.0, 1.0})
	// yAxis := up
	// zAxis := up.Cross([3]float32{1.0, 0.0, 0.0}).Mul(-1.0)

	if textures == 0 {
		textures = 1
	}

	mesh := Render.CreateMesh3D(name)
	vertices := make([]Mesh3DVertex, 4)
	indices := make([]uint32, 6)

	// Positions

	vertices[0][VertexPosIndex(0)] = -size[0] / 2.0
	vertices[0][VertexPosIndex(1)] = 0.0
	vertices[0][VertexPosIndex(2)] = size[1] / 2.0

	vertices[1][VertexPosIndex(0)] = size[0] / 2.0
	vertices[1][VertexPosIndex(1)] = 0.0
	vertices[1][VertexPosIndex(2)] = size[1] / 2.0

	vertices[2][VertexPosIndex(0)] = size[0] / 2.0
	vertices[2][VertexPosIndex(1)] = 0.0
	vertices[2][VertexPosIndex(2)] = -size[1] / 2.0

	vertices[3][VertexPosIndex(0)] = -size[0] / 2.0
	vertices[3][VertexPosIndex(1)] = 0.0
	vertices[3][VertexPosIndex(2)] = -size[1] / 2.0

	// for i := 0; i < 4; i++ {
	// 	// vec := xAxis.Mul(vertices[i][VertexPosIndex(0)]).Add(yAxis.Mul(vertices[i][VertexPosIndex(1)])).Add(zAxis.Mul(vertices[i][VertexPosIndex(2)]))
	// 	vec := mgl32.Vec3{vertices[i][0], vertices[i][1], vertices[i][2]}.Cross(up)
	// 	vertices[i][VertexPosIndex(0)] = vec[0]
	// 	vertices[i][VertexPosIndex(1)] = vec[1]
	// 	vertices[i][VertexPosIndex(2)] = vec[2]

	// }

	// Normals

	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			vertices[i][VertexNormalIndex(j)] = [3]float32{0.0, 1.0, 0.0}[j] //up[j] //yAxis[j]
		}
	}

	// TexCoords

	one := float32(textures)

	vertices[0][VertexTexCoordIndex(0)] = 0.0
	vertices[0][VertexTexCoordIndex(1)] = one

	vertices[1][VertexTexCoordIndex(0)] = one
	vertices[1][VertexTexCoordIndex(1)] = one

	vertices[2][VertexTexCoordIndex(0)] = one
	vertices[2][VertexTexCoordIndex(1)] = 0.0

	vertices[3][VertexTexCoordIndex(0)] = 0.0
	vertices[3][VertexTexCoordIndex(1)] = 0.0

	// Indices

	indices[0] = 0
	indices[1] = 1
	indices[2] = 2
	indices[3] = 2
	indices[4] = 3
	indices[5] = 0

	mesh.AddVertices(vertices, indices)
	mesh.Load()
	mat := &Material{}
	mat.InitDefault()
	mesh.SetMaterial(mat)

	return mesh
}

func (oglm *OpenGLMesh3D) GetName() string {
	return oglm.Name
}
