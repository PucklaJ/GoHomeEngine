package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
)

var DrawColor color.Color = Color{255, 255, 255, 255}
var PointSize float32 = 1.0
var LineWidth float32 = 1.0
var Filled bool = true

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

func DrawLine3D(pos1, pos2 mgl32.Vec3) {
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

func DrawLine2D(pos1, pos2 mgl32.Vec2) {
	line := toLine2D(pos1, pos2)
	var robj Shape2D
	robj.Init()
	robj.AddLines([]Line2D{line})
	robj.Load()
	robj.SetDrawMode(DRAW_MODE_LINES)
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

func DrawTriangle2D(pos1, pos2, pos3 mgl32.Vec2) {
	tri := toTriangle2D(pos1, pos2, pos3)
	var robj Shape2D
	robj.Init()
	if Filled {
		robj.AddTriangles([]Triangle2D{tri})
		robj.SetDrawMode(DRAW_MODE_TRIANGLES)
	} else {
		robj.AddLines(tri.ToLines())
		robj.SetDrawMode(DRAW_MODE_LINES)
	}
	robj.Load()
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

func DrawRectangle2D(pos1, pos2, pos3, pos4 mgl32.Vec2) {
	rect := toRectangle2D(pos1, pos2, pos3, pos4)
	var robj Shape2D
	robj.Init()
	if Filled {
		tris := rect.ToTriangles()
		robj.AddTriangles(tris[:])
		robj.SetDrawMode(DRAW_MODE_TRIANGLES)
	} else {
		robj.AddLines(rect.ToLines())
		robj.SetDrawMode(DRAW_MODE_LINES)
	}
	robj.Load()
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

func DrawCircle2D(pos mgl32.Vec2, radius float32) {
	circle := Circle2D{
		pos,
		radius,
		DrawColor,
	}

	var robj Shape2D
	robj.Init()
	if Filled {
		robj.AddTriangles(circle.ToTriangles(40))
		robj.SetDrawMode(DRAW_MODE_TRIANGLES)
	} else {
		robj.AddLines(circle.ToLines(40))
		robj.SetDrawMode(DRAW_MODE_LINES)
	}

	robj.Load()
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}
