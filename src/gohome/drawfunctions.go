package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"image/color"
)

// The color in which the shapes will be drawn
var DrawColor color.Color = Color{255, 255, 255, 255}
// The point size which will be used
var PointSize float32 = 1.0
// The line width which will be used
var LineWidth float32 = 1.0
// Wether the shapes should be filled
var Filled bool = true
// The detail of the curcle. The higher the more detailed
var CircleDetail int = 30

// Draws a point to point
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

// Draws a 2D line from pos1 to pos2
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

// Draws a 2D triangle between pos1, pos2 and pos3
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

// Draws a 2D rectangle between pos1,pos2,pos3 and pos4
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

// Draws a 2D circle with pos as the middle with radius
func DrawCircle2D(pos mgl32.Vec2, radius float32) {
	circle := Circle2D{
		pos,
		radius,
		DrawColor,
	}

	var robj Shape2D
	robj.Init()
	if Filled {
		robj.AddTriangles(circle.ToTriangles(CircleDetail))
		robj.SetDrawMode(DRAW_MODE_TRIANGLES)
	} else {
		robj.AddLines(circle.ToLines(CircleDetail))
		robj.SetDrawMode(DRAW_MODE_LINES)
	}

	robj.Load()
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

// Draws a polygon with positions
func DrawPolygon2D(positions ...mgl32.Vec2) {
	if len(positions) < 3 {
		ErrorMgr.Error("Polygon2D", "Draw", "Cannot draw polygon with less than 3 vertices")
		return
	}

	poly := toPolygon2D(positions...)
	var robj Shape2D
	robj.Init()
	if Filled {
		robj.AddTriangles(poly.ToTriangles())
		robj.SetDrawMode(DRAW_MODE_TRIANGLES)
	} else {
		robj.AddLines(poly.ToLines())
		robj.SetDrawMode(DRAW_MODE_LINES)
	}
	robj.Load()
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

// Draws a texture to [x,y]
func DrawTexture(tex Texture, x, y int) {
	DrawTextureAdv(tex, x, y, tex.GetWidth(), tex.GetHeight(), TextureRegion{
		Min: [2]float32{0.0, 0.0},
		Max: [2]float32{float32(tex.GetWidth()), float32(tex.GetHeight())},
	},
		FLIP_NONE)
}

// Draws a specific region of the texture with a flip
func DrawTextureAdv(tex Texture, x, y, width, height int, texReg TextureRegion, flip uint8) {
	var spr Sprite2D
	spr.InitTexture(tex)
	spr.Transform.Position[0] = float32(x)
	spr.Transform.Position[1] = float32(y)
	spr.Transform.Size[0] = float32(width)
	spr.Transform.Size[1] = float32(height)
	spr.TextureRegion = texReg
	spr.Flip = flip

	RenderMgr.RenderRenderObject(&spr)
}

// Draws a 3D point
func DrawPoint3D(pos mgl32.Vec3) {
	point := toVertex3D(pos)
	var robj Shape3D
	robj.Init()
	robj.AddPoint(point)
	robj.Load()
	robj.SetDrawMode(DRAW_MODE_POINTS)
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

// Draws a 3D line from pos1 to pos2
func DrawLine3D(pos1, pos2 mgl32.Vec3) {
	line := toLine3D(pos1, pos2)
	var robj Shape3D
	robj.Init()
	robj.AddLine(line)
	robj.Load()
	robj.SetDrawMode(DRAW_MODE_LINES)
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

// Draws a triangle between pos1,pos2 and pos3
func DrawTriangle3D(pos1, pos2, pos3 mgl32.Vec3) {
	tri := toTriangle3D(pos1, pos2, pos3)
	var robj Shape3D
	robj.Init()
	robj.AddTriangle(tri)
	robj.Load()
	robj.SetDrawMode(DRAW_MODE_TRIANGLES)
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}

// Draws a cube with pos as the middle and with,height and depth
// pitch, yaw and roll defines the rotation in degrees
func DrawCube(pos mgl32.Vec3, width, height, depth, pitch, yaw, roll float32) {
	tris := cubeToTriangle3Ds(width, height, depth)
	var robj Shape3D
	robj.Init()
	robj.AddTriangles(tris[:])
	robj.Load()
	robj.SetDrawMode(DRAW_MODE_TRIANGLES)
	robj.Transform.Position = pos
	pitch, yaw, roll = mgl32.DegToRad(pitch), mgl32.DegToRad(yaw), mgl32.DegToRad(roll)

	robj.Transform.Rotation = mgl32.QuatRotate(pitch, [3]float32{1.0, 0.0, 0.0}).Mul(mgl32.QuatRotate(yaw, [3]float32{0.0, 1.0, 0.0})).Mul(mgl32.QuatRotate(roll, [3]float32{0.0, 0.0, 1.0}))
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}
