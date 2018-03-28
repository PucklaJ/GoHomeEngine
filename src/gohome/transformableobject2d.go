package gohome

import (
	// "fmt"
	"github.com/go-gl/mathgl/mgl32"
)

type TransformableObject2D struct {
	Position      mgl32.Vec2
	Size          mgl32.Vec2
	Scale         mgl32.Vec2
	Rotation      float32
	Origin        mgl32.Vec2
	RotationPoint mgl32.Vec2

	oldPosition mgl32.Vec2
	oldSize     mgl32.Vec2
	oldScale    mgl32.Vec2
	oldRotation float32

	transformMatrix      mgl32.Mat3
	camNotRelativeMatrix mgl32.Mat3
}

func (tobj *TransformableObject2D) getOrigin(i int) float32 {
	return tobj.Size[i] * tobj.Scale[i] * ((tobj.Origin[i]*2.0 - 1.0) / -2.0)
}

func (tobj *TransformableObject2D) getRotationPoint(i int) float32 {
	return tobj.Size[i] * tobj.Scale[i] * ((tobj.RotationPoint[i]*2.0 - 1.0) / -2.0)
}

func (tobj *TransformableObject2D) valuesChanged() bool {
	return (tobj.Position != tobj.oldPosition || tobj.Size != tobj.oldSize || tobj.Scale != tobj.oldScale || tobj.Rotation != tobj.oldRotation)
}

func (tobj *TransformableObject2D) CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int) {
	var cam2d *Camera2D = nil
	if notRelativeToCamera != -1 && len(rmgr.camera2Ds) > notRelativeToCamera {
		cam2d = rmgr.camera2Ds[notRelativeToCamera]
	}
	if cam2d != nil {
		cam2d.CalculateViewMatrix()
	}
	// OT T -RPT R RPT S
	if tobj.valuesChanged() {
		tobj.transformMatrix = mgl32.Translate2D(tobj.getOrigin(0), tobj.getOrigin(1)).Mul3(mgl32.Translate2D(tobj.Position[0], tobj.Position[1])).Mul3(mgl32.Translate2D(-tobj.getRotationPoint(0), -tobj.getRotationPoint(1))).Mul3(mgl32.Rotate2D(-mgl32.DegToRad(tobj.Rotation)).Mat3()).Mul3(mgl32.Translate2D(tobj.getRotationPoint(0), tobj.getRotationPoint(1))).Mul3(mgl32.Scale2D(tobj.Scale[0]*tobj.Size[0], tobj.Scale[1]*tobj.Size[1]))

		tobj.oldPosition = tobj.Position
		tobj.oldSize = tobj.Size
		tobj.oldScale = tobj.Scale
		tobj.oldRotation = tobj.Rotation
	}
	if cam2d != nil {
		tobj.camNotRelativeMatrix = cam2d.GetInverseViewMatrix().Mul3(tobj.transformMatrix)
	} else {
		tobj.camNotRelativeMatrix = tobj.transformMatrix
	}
}

func (tobj *TransformableObject2D) GetTransformMatrix() mgl32.Mat3 {
	return tobj.camNotRelativeMatrix
}

func (tobj *TransformableObject2D) SetTransformMatrix(rmgr *RenderManager) {
	rmgr.setTransformMatrix2D(tobj.GetTransformMatrix())
}
