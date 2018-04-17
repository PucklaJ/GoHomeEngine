package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Mesh3DVertex [3 + 3 + 2 + 3]float32

type Mesh3D interface {
	AddVertices(vertices []Mesh3DVertex, indices []uint32)
	Load()
	Render()
	Terminate()
	SetMaterial(mat *Material)
	GetMaterial() *Material
	GetName() string
	GetNumVertices() uint32
	GetNumIndices() uint32
	GetVertices() []Mesh3DVertex
	GetIndices() []uint32
	CalculateTangents()
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
