package gohome

import (
	"github.com/PucklaJ/mathgl/mgl32"
)

func toLine3D(pos1 mgl32.Vec3, pos2 mgl32.Vec3) (line Line3D) {
	vecCol := ColorToVec4(DrawColor)

	line[0][0] = pos1[0]
	line[0][1] = pos1[1]
	line[0][2] = pos1[2]

	line[0][3] = vecCol[0]
	line[0][4] = vecCol[1]
	line[0][5] = vecCol[2]
	line[0][6] = vecCol[3]

	line[1][0] = pos2[0]
	line[1][1] = pos2[1]
	line[1][2] = pos2[2]

	line[1][3] = vecCol[0]
	line[1][4] = vecCol[1]
	line[1][5] = vecCol[2]
	line[1][6] = vecCol[3]

	return
}

func toPoint2D(point mgl32.Vec2) Shape2DVertex {
	vecCol := ColorToVec4(DrawColor)

	return Shape2DVertex{
		point[0], point[1],
		vecCol[0], vecCol[1], vecCol[2], vecCol[3],
	}
}

func toLine2D(pos1 mgl32.Vec2, pos2 mgl32.Vec2) (line Line2D) {
	vecCol := ColorToVec4(DrawColor)

	line[0][0] = pos1[0]
	line[0][1] = pos1[1]
	line[0][2] = vecCol[0]
	line[0][3] = vecCol[1]
	line[0][4] = vecCol[2]
	line[0][5] = vecCol[3]

	line[1][0] = pos2[0]
	line[1][1] = pos2[1]
	line[1][2] = vecCol[0]
	line[1][3] = vecCol[1]
	line[1][4] = vecCol[2]
	line[1][5] = vecCol[3]

	return
}

func toTriangle2D(pos1 mgl32.Vec2, pos2 mgl32.Vec2, pos3 mgl32.Vec2) (tri Triangle2D) {
	vecCol := ColorToVec4(DrawColor)

	tri[0][0] = pos1[0]
	tri[0][1] = pos1[1]
	tri[0][2] = vecCol[0]
	tri[0][3] = vecCol[1]
	tri[0][4] = vecCol[2]
	tri[0][5] = vecCol[3]

	tri[1][0] = pos2[0]
	tri[1][1] = pos2[1]
	tri[1][2] = vecCol[0]
	tri[1][3] = vecCol[1]
	tri[1][4] = vecCol[2]
	tri[1][5] = vecCol[3]

	tri[2][0] = pos3[0]
	tri[2][1] = pos3[1]
	tri[2][2] = vecCol[0]
	tri[2][3] = vecCol[1]
	tri[2][4] = vecCol[2]
	tri[2][5] = vecCol[3]

	return
}

func toRectangle2D(pos1, pos2, pos3, pos4 mgl32.Vec2) (rect Rectangle2D) {
	vecCol := ColorToVec4(DrawColor)

	rect[0][0] = pos1[0]
	rect[0][1] = pos1[1]
	rect[0][2] = vecCol[0]
	rect[0][3] = vecCol[1]
	rect[0][4] = vecCol[2]
	rect[0][5] = vecCol[3]

	rect[1][0] = pos2[0]
	rect[1][1] = pos2[1]
	rect[1][2] = vecCol[0]
	rect[1][3] = vecCol[1]
	rect[1][4] = vecCol[2]
	rect[1][5] = vecCol[3]

	rect[2][0] = pos3[0]
	rect[2][1] = pos3[1]
	rect[2][2] = vecCol[0]
	rect[2][3] = vecCol[1]
	rect[2][4] = vecCol[2]
	rect[2][5] = vecCol[3]

	rect[3][0] = pos4[0]
	rect[3][1] = pos4[1]
	rect[3][2] = vecCol[0]
	rect[3][3] = vecCol[1]
	rect[3][4] = vecCol[2]
	rect[3][5] = vecCol[3]

	return
}

func toPolygon2D(positions ...mgl32.Vec2) (poly Polygon2D) {
	vecCol := ColorToVec4(DrawColor)
	poly.Points = append(poly.Points, make([]Shape2DVertex, len(positions))...)
	for i := 0; i < len(positions); i++ {
		vertex := Shape2DVertex{
			positions[i][0], positions[i][1],
			vecCol[0], vecCol[1], vecCol[2], vecCol[3],
		}
		poly.Points[i] = vertex
	}
	return
}

func toVertex3D(pos mgl32.Vec3) (vert Shape3DVertex) {
	for i := 0; i < 3; i++ {
		vert[i] = pos[i]
	}

	vecCol := ColorToVec4(DrawColor)
	for i := 0; i < 4; i++ {
		vert[i+3] = vecCol[i]
	}
	return
}

func toTriangle3D(pos1, pos2, pos3 mgl32.Vec3) (tri Triangle3D) {
	pos := [3]mgl32.Vec3{
		pos1, pos2, pos3,
	}

	for i := 0; i < 3; i++ {
		tri[i] = toVertex3D(pos[i])
	}
	return
}

func cubeToTriangle3Ds(width, height, depth float32) (tris [6 * 2]Triangle3D) {
	const LDB = 0
	const RDB = 1
	const RDF = 2
	const LDF = 3
	const LUB = 4
	const RUB = 5
	const RUF = 6
	const LUF = 7

	p := [8]mgl32.Vec3{
		mgl32.Vec3{-width / 2.0, -height / 2.0, -depth / 2.0}, // LDB
		mgl32.Vec3{+width / 2.0, -height / 2.0, -depth / 2.0}, // RDB
		mgl32.Vec3{+width / 2.0, -height / 2.0, +depth / 2.0}, // RDF
		mgl32.Vec3{-width / 2.0, -height / 2.0, +depth / 2.0}, // LDF

		mgl32.Vec3{-width / 2.0, +height / 2.0, -depth / 2.0}, // LUB
		mgl32.Vec3{+width / 2.0, +height / 2.0, -depth / 2.0}, // RUB
		mgl32.Vec3{+width / 2.0, +height / 2.0, +depth / 2.0}, // RUF
		mgl32.Vec3{-width / 2.0, +height / 2.0, +depth / 2.0}, // LUF
	}

	tris = [6 * 2]Triangle3D{
		toTriangle3D(p[LUF], p[LDF], p[RDF]), // FRONT
		toTriangle3D(p[RDF], p[RUF], p[LUF]),

		toTriangle3D(p[RUF], p[RDF], p[RDB]), // RIGHT
		toTriangle3D(p[RDB], p[RUB], p[RUF]),

		toTriangle3D(p[RUB], p[RDB], p[LDB]), // BACK
		toTriangle3D(p[LDB], p[LUB], p[RUB]),

		toTriangle3D(p[LUB], p[LDB], p[LDF]), // LEFT
		toTriangle3D(p[LDF], p[LUF], p[LUB]),

		toTriangle3D(p[LUB], p[LUF], p[RUF]), // UP
		toTriangle3D(p[RUF], p[RUB], p[LUB]),

		toTriangle3D(p[LDF], p[LDB], p[RDB]), // DOWN
		toTriangle3D(p[RDB], p[RDF], p[LDF]),
	}

	return
}
