package physics2d

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
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
		ScalarToPixel(vec[0]),
		ScalarToPixel(vec[1]),
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

type PhysicsManager2D struct {
	World box2d.B2World
}

func (this *PhysicsManager2D) Init(gravity mgl32.Vec2) {
	this.World = box2d.MakeB2World(ToBox2DDirection(gravity))
	nw, nh := gohome.Render.GetNativeResolution()
	WORLD_SIZE[0] = float32(nw)
	WORLD_SIZE[1] = float32(nh)
	gohome.UpdateMgr.AddObject(this)
}

func (this *PhysicsManager2D) Update(delta_time float32) {
	this.World.Step(float64(delta_time), VELOCITY_ITERATIONS, POSITION_ITERATIONS)
}
