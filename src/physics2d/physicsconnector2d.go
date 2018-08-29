package physics2d

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
)

type PhysicsConnector2D struct {
	Transform *gohome.TransformableObject2D
	Body      *box2d.B2Body
	Offset    mgl32.Vec2
}

func (this *PhysicsConnector2D) Init(tobj *gohome.TransformableObject2D, body *box2d.B2Body) {
	this.Transform = tobj
	this.Body = body
}

func (this *PhysicsConnector2D) Update(delta_time float32) {
	pixelPos := ToPixelCoordinates(this.Body.GetPosition())
	offset := this.Transform.Origin.Sub(mgl32.Vec2{0.5, 0.5})
	offset[0] *= this.Transform.Size[0] * this.Transform.Scale[0]
	offset[1] *= this.Transform.Size[1] * this.Transform.Scale[1]
	this.Transform.Position = pixelPos.Add(offset).Add(this.Offset)
	this.Transform.Rotation = ToPixelAngle(this.Body.GetAngle())
}
