package physics2d

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type PhysicsConnector2D struct {
	Transform *gohome.TransformableObject2D
	Body      *box2d.B2Body
	Offset    mgl32.Vec2

	pmgr *PhysicsManager2D
}

func (this *PhysicsConnector2D) Init(tobj *gohome.TransformableObject2D, body *box2d.B2Body, pmgr *PhysicsManager2D) {
	this.Transform = tobj
	this.Body = body
	pmgr.connectors = append(pmgr.connectors, this)
	this.pmgr = pmgr
}

func (this *PhysicsConnector2D) Update() {
	pixelPos := ToPixelCoordinates(this.Body.GetPosition())
	offset := this.Transform.Origin.Sub(mgl32.Vec2{0.5, 0.5})
	offset[0] *= this.Transform.Size[0] * this.Transform.Scale[0]
	offset[1] *= this.Transform.Size[1] * this.Transform.Scale[1]
	this.Transform.Position = pixelPos.Add(offset).Add(this.Offset)
	this.Transform.Rotation = ToPixelAngle(this.Body.GetAngle())
}

func (this *PhysicsConnector2D) Terminate() {
	for i := 0; i < len(this.pmgr.connectors); i++ {
		if this.pmgr.connectors[i] == this {
			this.pmgr.connectors = append(this.pmgr.connectors[:i], this.pmgr.connectors[i+1:]...)
			return
		}
	}
}
