package physics2d

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/colornames"
)

type PhysicsDebugDraw2D struct {
	gohome.NilRenderObject
	mgr *PhysicsManager2D
}

func (this *PhysicsDebugDraw2D) GetType() gohome.RenderType {
	return gohome.TYPE_2D_NORMAL
}

func (this *PhysicsDebugDraw2D) Render() {
	gohome.DrawColor = colornames.Purple
	gohome.PointSize = 3.0
	gohome.Filled = true

	w := this.mgr.World
	for b := w.GetBodyList(); b != nil; b = b.GetNext() {
		xf := b.GetTransform()
		for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
			this.DrawFixture(f, &xf)
		}
	}
}

func (this *PhysicsDebugDraw2D) DrawFixture(f *box2d.B2Fixture, xf *box2d.B2Transform) {
	pos := ToPixelCoordinates(xf.P)
	switch f.GetType() {
	case box2d.B2Shape_Type.E_circle:
		gohome.DrawCircle2D(pos, float32(f.GetShape().GetRadius()))
	case box2d.B2Shape_Type.E_polygon:
		mat := mgl32.Translate2D(pos[0], pos[1]).Mul3(mgl32.Rotate2D(-float32(xf.Q.GetAngle())).Mat3())
		polygon := f.GetShape().(*box2d.B2PolygonShape)
		vertices := make([]mgl32.Vec2, polygon.M_count)
		for i := 0; i < polygon.M_count; i++ {
			vertices[i] = mat.Mul3x1(ToPixelDirection(polygon.M_vertices[i]).Vec3(1.0)).Vec2()
		}
		gohome.DrawPolygon2D(vertices[:]...)
	}
}
