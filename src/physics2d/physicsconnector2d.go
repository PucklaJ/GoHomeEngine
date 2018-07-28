package physics2d

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
)

type PhysicsConnector2D struct {
	Transform *gohome.TransformableObject2D
	Body      *box2d.B2Body
}

func (this *PhysicsConnector2D) Init(tobj *gohome.TransformableObject2D, body *box2d.B2Body) {
	this.Transform = tobj
	this.Body = body
}

func (this *PhysicsConnector2D) Update(delta_time float32) {
	this.Transform.Position = ToPixelCoordinates(this.Body.GetPosition())
	this.Transform.Rotation = float32(this.Body.GetAngle())
}
