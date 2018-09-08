package physics2d

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"golang.org/x/image/colornames"
)

type PhysicsDebugDraw2D struct {
	gohome.NilRenderObject

	DrawBodies      bool
	DrawJoints      bool
	DrawAABBs       bool
	OnlyDrawStatic  bool
	OnlyDrawDynamic bool
	Visible         bool

	mgr *PhysicsManager2D

	lines     []gohome.Line2D
	triangles []gohome.Triangle2D
}

func (this *PhysicsDebugDraw2D) GetType() gohome.RenderType {
	return gohome.TYPE_2D_NORMAL
}

func (this *PhysicsDebugDraw2D) RendersLast() bool {
	return true
}

func (this *PhysicsDebugDraw2D) IsVisible() bool {
	return this.Visible
}

func (this *PhysicsDebugDraw2D) Render() {
	w := this.mgr.World
	if this.DrawBodies {
		for b := w.GetBodyList(); b != nil; b = b.GetNext() {
			if (!this.OnlyDrawDynamic || b.GetType() == box2d.B2BodyType.B2_dynamicBody) &&
				(!this.OnlyDrawStatic || b.GetType() == box2d.B2BodyType.B2_staticBody) {
				xf := b.GetTransform()
				for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
					this.DrawFixture(f, &xf, b.IsAwake())
				}
			}
		}
	}

	if this.DrawJoints {
		for j := w.GetJointList(); j != nil; j = j.GetNext() {
			this.DrawJoint(j)
		}
	}

	if this.DrawAABBs {
		brightblue := colornames.Skyblue
		brightblue.A = 120
		gohome.DrawColor = brightblue
		bp := &w.M_contactManager.M_broadPhase
		for b := w.GetBodyList(); b != nil; b = b.GetNext() {
			if !b.IsActive() {
				continue
			}
			for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
				for i := 0; i < f.M_proxyCount; i++ {
					proxy := f.M_proxies[i]
					aabb := bp.GetFatAABB(proxy.ProxyId)
					var vs [4]mgl32.Vec2
					vs[0] = ToPixelCoordinates(aabb.LowerBound)
					vs[1] = ToPixelCoordinates(box2d.B2Vec2{aabb.UpperBound.X, aabb.LowerBound.Y})
					vs[2] = ToPixelCoordinates(aabb.UpperBound)
					vs[3] = ToPixelCoordinates(box2d.B2Vec2{aabb.LowerBound.X, aabb.UpperBound.Y})
					gohome.DrawRectangle2D(vs[0], vs[1], vs[2], vs[3])
				}
			}
		}
	}

	var shapes gohome.Shape2D
	if len(this.triangles) != 0 {
		shapes.Init()
		shapes.Depth = 255
		shapes.AddTriangles(this.triangles)
		shapes.SetDrawMode(gohome.DRAW_MODE_TRIANGLES)
		shapes.Load()
		gohome.RenderMgr.RenderRenderObject(&shapes)
		shapes.Terminate()

		this.triangles = this.triangles[:0]
	}
	if len(this.lines) != 0 {
		shapes.Init()
		shapes.Depth = 255
		shapes.AddLines(this.lines)
		shapes.SetDrawMode(gohome.DRAW_MODE_LINES)
		shapes.Load()
		gohome.RenderMgr.RenderRenderObject(&shapes)
		shapes.Terminate()

		this.lines = this.lines[:0]
	}
}

func (this *PhysicsDebugDraw2D) DrawJoint(j box2d.B2JointInterface) {
	red := colornames.Red
	if j.IsActive() {
		red.A = 180
	} else {
		red.A = 120
	}

	pos1 := ToPixelCoordinates(j.GetBodyA().GetPosition())
	pos2 := ToPixelCoordinates(j.GetBodyB().GetPosition())

	var line gohome.Line2D

	switch j.GetType() {
	case box2d.B2JointType.E_ropeJoint:
		rope := j.(*box2d.B2RopeJoint)
		lca := ToPixelDirection(rope.M_localCenterA)
		lcb := ToPixelDirection(rope.M_localCenterB)
		line[0].Make(pos1.Add(lca), red)
		line[1].Make(pos2.Add(lcb), red)
		this.lines = append(this.lines, line)
	default:
		line[0].Make(pos1, red)
		line[1].Make(pos2, red)
		this.lines = append(this.lines, line)
	}

}

func (this *PhysicsDebugDraw2D) DrawFixture(f *box2d.B2Fixture, xf *box2d.B2Transform, awake bool) {
	col := colornames.Purple
	if f.IsSensor() {
		col = colornames.Blue
	}
	if awake {
		col.A = 180
	} else {
		col.A = 120
	}
	gohome.DrawColor = col
	switch f.GetType() {
	case box2d.B2Shape_Type.E_circle:
		this.DrawCircle(f, xf)
	case box2d.B2Shape_Type.E_polygon:
		this.DrawPolygon(f, xf)
	case box2d.B2Shape_Type.E_edge:
		this.DrawEdge(f, xf)
	case box2d.B2Shape_Type.E_chain:
		this.DrawChain(f, xf)
	}
}

func (this *PhysicsDebugDraw2D) DrawCircle(f *box2d.B2Fixture, xf *box2d.B2Transform) {
	radius := ScalarToPixel(f.GetShape().GetRadius())
	b2offset := f.GetShape().(*box2d.B2CircleShape).M_p
	offset := ToPixelDirection(b2offset)
	b2Pos := xf.P
	b2Pos.X += b2offset.X
	b2Pos.Y += b2offset.Y
	pos := ToPixelCoordinates(b2Pos)
	var circle2d gohome.Circle2D
	circle2d.Position = pos
	circle2d.Radius = radius
	circle2d.Col = gohome.DrawColor
	this.triangles = append(this.triangles, circle2d.ToTriangles(gohome.CircleDetail)...)
	pos2 := pos.Add(mgl32.Translate2D(-offset.X(), -offset.Y()).Mul3(mgl32.Rotate2D(-float32(xf.Q.GetAngle())).Mat3()).Mul3(mgl32.Translate2D(offset.X(), offset.Y())).Mul3(mgl32.Scale2D(radius, radius)).Mul3x1(mgl32.Vec3{1.0, 0.0, 1.0}).Vec2())
	blue := colornames.Blue
	blue.A = 150
	var line gohome.Line2D
	line[0].Make(pos, blue)
	line[1].Make(pos2, blue)
	this.lines = append(this.lines, line)
}

func (this *PhysicsDebugDraw2D) DrawPolygon(f *box2d.B2Fixture, xf *box2d.B2Transform) {
	pos := ToPixelCoordinates(xf.P)
	mat := mgl32.Translate2D(pos[0], pos[1]).Mul3(mgl32.Rotate2D(-float32(xf.Q.GetAngle())).Mat3())
	polygon := f.GetShape().(*box2d.B2PolygonShape)
	vertices := make([]mgl32.Vec2, polygon.M_count)
	for i := 0; i < polygon.M_count; i++ {
		vertices[i] = mat.Mul3x1(ToPixelDirection(polygon.M_vertices[i]).Vec3(1.0)).Vec2()
	}
	if polygon.M_count == 3 {
		var tri gohome.Triangle2D
		for i := 0; i < 3; i++ {
			tri[i].Make(vertices[i], gohome.DrawColor)
		}
		this.triangles = append(this.triangles, tri)
	} else if polygon.M_count == 4 {
		var rect gohome.Rectangle2D
		for i := 0; i < 4; i++ {
			rect[i].Make(vertices[i], gohome.DrawColor)
		}
		tris := rect.ToTriangles()
		this.triangles = append(this.triangles, tris[:]...)
	} else {
		var poly gohome.Polygon2D
		for i := 0; i < polygon.M_count; i++ {
			var vert gohome.Shape2DVertex
			vert.Make(vertices[i], gohome.DrawColor)
			poly.Points = append(poly.Points, vert)
		}
		this.triangles = append(this.triangles, poly.ToTriangles()...)
	}

	radius := float32(10.0)
	pos2 := pos.Add(mgl32.Rotate2D(-float32(xf.Q.GetAngle())).Mul2(mgl32.Scale2D(radius, radius).Mat2()).Mul2x1(mgl32.Vec2{1.0, 0.0}))
	blue := colornames.Blue
	blue.A = 150
	var line gohome.Line2D
	line[0].Make(pos, blue)
	line[1].Make(pos2, blue)
	this.lines = append(this.lines, line)
}

func (this *PhysicsDebugDraw2D) DrawEdge(f *box2d.B2Fixture, xf *box2d.B2Transform) {
	edge := f.GetShape().(*box2d.B2EdgeShape)
	pos := ToPixelCoordinates(xf.P)
	mat := mgl32.Translate2D(pos[0], pos[1]).Mul3(mgl32.Rotate2D(-float32(xf.Q.GetAngle())).Mat3())
	var numVertices uint32 = 2
	if edge.M_hasVertex0 {
		numVertices++
	}
	if edge.M_hasVertex3 {
		numVertices++
	}
	vertices := make([]mgl32.Vec2, numVertices)
	for i := 0; i < len(vertices); i++ {
		if i == 0 {
			if edge.M_hasVertex0 {
				vertices[i] = mat.Mul3x1(ToPixelDirection(edge.M_vertex0).Vec3(1.0)).Vec2()
			} else {
				vertices[i] = mat.Mul3x1(ToPixelDirection(edge.M_vertex1).Vec3(1.0)).Vec2()
			}
		} else if i == 1 {
			if edge.M_hasVertex0 {
				vertices[i] = mat.Mul3x1(ToPixelDirection(edge.M_vertex1).Vec3(1.0)).Vec2()
			} else {
				vertices[i] = mat.Mul3x1(ToPixelDirection(edge.M_vertex2).Vec3(1.0)).Vec2()
			}
		} else if i == 2 {
			if edge.M_hasVertex0 {
				vertices[i] = mat.Mul3x1(ToPixelDirection(edge.M_vertex2).Vec3(1.0)).Vec2()
			} else {
				vertices[i] = mat.Mul3x1(ToPixelDirection(edge.M_vertex3).Vec3(1.0)).Vec2()
			}
		} else {
			vertices[i] = mat.Mul3x1(ToPixelDirection(edge.M_vertex3).Vec3(1.0)).Vec2()
		}
	}

	for i := 0; i < len(vertices)-1; i++ {
		var line gohome.Line2D
		line[0].Make(vertices[i], gohome.DrawColor)
		line[1].Make(vertices[i+1], gohome.DrawColor)
		this.lines = append(this.lines, line)
	}
}

func (this *PhysicsDebugDraw2D) DrawChain(f *box2d.B2Fixture, xf *box2d.B2Transform) {
	chain := f.GetShape().(*box2d.B2ChainShape)
	pos := ToPixelCoordinates(xf.P)
	mat := mgl32.Translate2D(pos[0], pos[1]).Mul3(mgl32.Rotate2D(-float32(xf.Q.GetAngle())).Mat3())
	numVertices := chain.M_count
	vertices := make([]mgl32.Vec2, numVertices)
	for i := 0; i < len(vertices); i++ {
		vertices[i] = mat.Mul3x1(ToPixelDirection(chain.M_vertices[i]).Vec3(1.0)).Vec2()
	}
	for i := 0; i < len(vertices)-1; i++ {
		var line gohome.Line2D
		line[0].Make(vertices[i], gohome.DrawColor)
		line[1].Make(vertices[i+1], gohome.DrawColor)
		this.lines = append(this.lines, line)
	}
}
