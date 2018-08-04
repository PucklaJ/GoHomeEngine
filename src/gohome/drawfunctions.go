package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
)

var DrawColor color.Color = Color{255, 255, 255, 255}
var PointSize float32 = 1.0
var LineWidth float32 = 1.0

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

func DrawLine3D(pos1 mgl32.Vec3, pos2 mgl32.Vec3) {
	line := toLine3D(pos1, pos2)
	var robj Lines3D
	robj.Init()
	robj.AddLine(line)
	robj.Load()
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

func DrawPoint2D(point mgl32.Vec2) {
	point2D := toPoint2D(point)
	var robj Shape2D
	robj.Init()
	robj.SetPointSize(PointSize)
	robj.SetLineWidth(LineWidth)
	robj.AddPoints([]Shape2DVertex{point2D})
	robj.Load()
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

func DrawLine2D(pos1 mgl32.Vec2, pos2 mgl32.Vec2) {
	line := toLine2D(pos1, pos2)
	var robj Shape2D
	robj.Init()
	robj.AddLines([]Line2D{line})
	robj.Load()
	robj.SetDrawMode(DRAW_MODE_LINES)
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

func DrawTriangle2D(pos1 mgl32.Vec2, pos2 mgl32.Vec2, pos3 mgl32.Vec2) {
	tri := toTriangle2D(pos1, pos2, pos3)
	var robj Shape2D
	robj.Init()
	robj.AddTriangles([]Triangle2D{tri})
	robj.Load()
	robj.SetDrawMode(DRAW_MODE_TRIANGLES)
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}
