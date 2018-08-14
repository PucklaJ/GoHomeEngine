package physics2d

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/elliotmr/tmx"
	"github.com/go-gl/mathgl/mgl32"
	"strconv"
	"strings"
)

var (
	PIXEL_PER_METER     float32 = 100.0
	WORLD_SIZE          mgl32.Vec2
	VELOCITY_ITERATIONS uint32 = 6
	POSITION_ITERATIONS uint32 = 2
)

func ScalarToPixel(v float64) float32 {
	return float32(v) * PIXEL_PER_METER
}

func ScalarToBox2D(v float32) float64 {
	return float64(v / PIXEL_PER_METER)
}

func Vec2ToB2Vec2(vec mgl32.Vec2) box2d.B2Vec2 {
	return box2d.B2Vec2{
		ScalarToBox2D(vec[0]),
		ScalarToBox2D(vec[1]),
	}
}

func B2Vec2ToVec2(vec box2d.B2Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		ScalarToPixel(vec.X),
		ScalarToPixel(vec.Y),
	}
}

func ToPixelDirection(vec box2d.B2Vec2) mgl32.Vec2 {
	vec1 := B2Vec2ToVec2(vec)
	vec1[1] = -vec1[1]
	return vec1
}

func ToBox2DDirection(vec mgl32.Vec2) box2d.B2Vec2 {
	vec[1] = -vec[1]
	return Vec2ToB2Vec2(vec)
}

func ToPixelCoordinates(vec box2d.B2Vec2) mgl32.Vec2 {
	vec1 := B2Vec2ToVec2(vec)
	vec1[1] = WORLD_SIZE[1] - vec1[1]
	return vec1
}

func ToBox2DCoordinates(vec mgl32.Vec2) box2d.B2Vec2 {
	vec[1] = WORLD_SIZE[1] - vec[1]
	return Vec2ToB2Vec2(vec)
}

func ToBox2DAngle(angle float32) float64 {
	return float64(mgl32.DegToRad(angle))
}

func ToPixelAngle(angle float64) float32 {
	return mgl32.RadToDeg(float32(angle))
}

type PhysicsManager2D struct {
	World box2d.B2World
}

func (this *PhysicsManager2D) Init(gravity mgl32.Vec2) {
	this.World = box2d.MakeB2World(ToBox2DDirection(gravity))
	nw, nh := gohome.Render.GetNativeResolution()
	WORLD_SIZE[0] = float32(nw)
	WORLD_SIZE[1] = float32(nh)

	gohome.ErrorMgr.Log("Physics", "Box2D", "Initialized!")
}

func (this *PhysicsManager2D) Update(delta_time float32) {
	this.World.Step(float64(delta_time), int(VELOCITY_ITERATIONS), int(POSITION_ITERATIONS))
}

func (this *PhysicsManager2D) CreateDynamicBox(pos mgl32.Vec2, size mgl32.Vec2) *box2d.B2Body {
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	bodyDef.Position = ToBox2DCoordinates(pos)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(ScalarToBox2D(size[0])/2.0, ScalarToBox2D(size[1])/2.0)
	body := this.World.CreateBody(&bodyDef)
	body.CreateFixture(&shape, 1.0)
	return body
}

func (this *PhysicsManager2D) CreateStaticBox(pos mgl32.Vec2, size mgl32.Vec2) *box2d.B2Body {
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_staticBody
	bodyDef.Position = ToBox2DCoordinates(pos)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(ScalarToBox2D(size[0])/2.0, ScalarToBox2D(size[1])/2.0)
	body := this.World.CreateBody(&bodyDef)
	body.CreateFixture(&shape, 1.0)
	return body
}

func (this *PhysicsManager2D) CreateDynamicCircle(pos mgl32.Vec2, radius float32) *box2d.B2Body {
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	bodyDef.Position = ToBox2DCoordinates(pos)
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(ScalarToBox2D(radius))
	body := this.World.CreateBody(&bodyDef)
	body.CreateFixture(&shape, 1.0)
	return body
}

func (this *PhysicsManager2D) CreateStaticCircle(pos mgl32.Vec2, radius float32) *box2d.B2Body {
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_staticBody
	bodyDef.Position = ToBox2DCoordinates(pos)
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(ScalarToBox2D(radius))
	body := this.World.CreateBody(&bodyDef)
	body.CreateFixture(&shape, 1.0)
	return body
}

func (this *PhysicsManager2D) GetDebugDraw() PhysicsDebugDraw2D {
	return PhysicsDebugDraw2D{
		mgr:        this,
		DrawBodies: true,
		DrawJoints: true,
		DrawAABBs:  false,
	}
}

func (this *PhysicsManager2D) LayerToCollision(tiledmap *gohome.TiledMap, layerName string) {
	layers := tiledmap.Layers
	for i := 0; i < len(layers); i++ {
		l := layers[i]
		if !strings.Contains(l.Name, layerName) {
			continue
		}
		objs := l.Objects
		if len(objs) == 0 {
			continue
		}
		var lx, ly float64 = 0.0, 0.0
		if l.OffsetX != nil {
			lx = *l.OffsetX
		}
		if l.OffsetY != nil {
			ly = *l.OffsetY
		}
		for j := 0; j < len(objs); j++ {
			o := objs[j]
			if o.Ellipse != nil {
				if o.Width != nil && o.Height != nil {
					this.CreateEllipse(lx+o.X, ly+o.Y, *o.Width, *o.Height)
				}
			} else if o.Polygon != nil {
				this.CreatePolygon(lx+o.X, ly+o.Y, o.Polygon)
			} else if o.Polyline != nil {
				this.CreatePolyline(lx+o.X, ly+o.Y, o.Polyline)
			} else if o.Point == nil && o.Text == nil && o.GID == nil {
				if o.Width != nil && o.Height != nil {
					this.CreateRectangle(lx+o.X, ly+o.Y, *o.Width, *o.Height)
				}
			}
		}
	}
}

func (this *PhysicsManager2D) CreateEllipse(X, Y, Width, Height float64) {
	radius := float32((Width + Height) / 2.0 / 2.0)
	pos := [2]float32{float32(X) + radius, float32(Y) + radius}

	this.CreateStaticCircle(pos, radius)
}

func (this *PhysicsManager2D) CreateRectangle(X, Y, Width, Height float64) {
	size := mgl32.Vec2{float32(Width), float32(Height)}
	pos := mgl32.Vec2{float32(X), float32(Y)}.Add(size.Mul(0.5))

	this.CreateStaticBox(pos, size)
}

func (this *PhysicsManager2D) CreatePolygon(X, Y float64, poly *tmx.Polygon) {
	points := strings.Split(poly.Points, " ")
	if len(points) > 8 {
		gohome.ErrorMgr.Error("Physics", "Box2D", "Couldn't create collision polygon: It has more than 8 vertices")
		return
	}
	vertices := make([]mgl32.Vec2, len(points))
	b2vertices := make([]box2d.B2Vec2, len(points))
	for i := 0; i < len(points); i++ {
		point := strings.Split(points[i], ",")
		x, _ := strconv.ParseFloat(point[0], 32)
		y, _ := strconv.ParseFloat(point[1], 32)
		vertices[i][0] = float32(x)
		vertices[i][1] = float32(y)
	}
	for i := 0; i < len(vertices); i++ {
		b2vertices[i] = ToBox2DDirection(vertices[i])
	}

	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_staticBody
	bodyDef.Position = ToBox2DCoordinates([2]float32{float32(X), float32(Y)})
	shape := box2d.MakeB2PolygonShape()

	shape.Set(b2vertices, len(b2vertices))
	body := this.World.CreateBody(&bodyDef)
	body.CreateFixture(&shape, 1.0)
}

func (this *PhysicsManager2D) CreatePolyline(X, Y float64, line *tmx.Polyline) {
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_staticBody
	bodyDef.Position = ToBox2DCoordinates([2]float32{float32(X), float32(Y)})
	shape := box2d.MakeB2ChainShape()

	points := strings.Split(line.Points, " ")
	vertices := make([]mgl32.Vec2, len(points))
	b2vertices := make([]box2d.B2Vec2, len(points))
	for i := 0; i < len(points); i++ {
		point := strings.Split(points[i], ",")
		x, _ := strconv.ParseFloat(point[0], 32)
		y, _ := strconv.ParseFloat(point[1], 32)
		vertices[i][0] = float32(x)
		vertices[i][1] = float32(y)
	}
	for i := 0; i < len(vertices); i++ {
		b2vertices[i] = ToBox2DDirection(vertices[i])
	}
	shape.CreateChain(b2vertices, len(b2vertices))
	body := this.World.CreateBody(&bodyDef)
	body.CreateFixture(&shape, 1.0)
}
