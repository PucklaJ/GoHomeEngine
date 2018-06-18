package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
)

type TransformableObject3D struct {
	Position mgl32.Vec3
	Scale    mgl32.Vec3
	Rotation mgl32.Vec3

	oldPosition mgl32.Vec3
	oldScale    mgl32.Vec3
	oldRotation mgl32.Vec3

	transformMatrix      mgl32.Mat4
	camNotRelativeMatrix mgl32.Mat4
}

func (tobj *TransformableObject3D) valuesChanged() bool {
	return (tobj.Position != tobj.oldPosition || tobj.Scale != tobj.oldScale || tobj.Rotation != tobj.oldRotation)
}

func (tobj *TransformableObject3D) CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int) {
	var cam3d *Camera3D = nil
	if rmgr != nil {
		if notRelativeToCamera != -1 && len(rmgr.camera3Ds) > notRelativeToCamera {
			cam3d = rmgr.camera3Ds[notRelativeToCamera]
		}
		if cam3d != nil {
			cam3d.CalculateViewMatrix()
		}
	}

	if tobj.valuesChanged() {
		// T RX RY RZ S
		T := mgl32.Translate3D(tobj.Position[0], tobj.Position[1], tobj.Position[2])
		RX := mgl32.QuatRotate(mgl32.DegToRad(tobj.Rotation[0]), [3]float32{1.0, 0.0, 0.0}).Mat4()
		RY := mgl32.QuatRotate(mgl32.DegToRad(tobj.Rotation[1]), [3]float32{0.0, 1.0, 0.0}).Mat4()
		RZ := mgl32.QuatRotate(mgl32.DegToRad(tobj.Rotation[2]), [3]float32{0.0, 0.0, 1.0}).Mat4()
		S := mgl32.Scale3D(tobj.Scale[0], tobj.Scale[1], tobj.Scale[2])
		tobj.transformMatrix = T.Mul4(RX).Mul4(RY).Mul4(RZ).Mul4(S)

		tobj.oldPosition = tobj.Position
		tobj.oldScale = tobj.Scale
		tobj.oldRotation = tobj.Rotation
	}
	if cam3d != nil {
		tobj.camNotRelativeMatrix = cam3d.GetInverseViewMatrix().Mul4(tobj.transformMatrix)
	} else {
		tobj.camNotRelativeMatrix = tobj.transformMatrix
	}
}

func (tobj *TransformableObject3D) GetTransformMatrix() mgl32.Mat4 {
	return tobj.camNotRelativeMatrix
}

func (tobj *TransformableObject3D) SetTransformMatrix(rmgr *RenderManager) {
	rmgr.setTransformMatrix3D(tobj.GetTransformMatrix())
}

func DefaultTransformableObject3D() *TransformableObject3D {
	transform := TransformableObject3D{
		Scale: [3]float32{1.0, 1.0, 1.0},
	}

	return &transform
}
