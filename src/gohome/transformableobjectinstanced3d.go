package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type TransformableObjectInstanced3D struct {
	Position mgl32.Vec3
	Scale    mgl32.Vec3
	Rotation mgl32.Quat

	oldPosition mgl32.Vec3
	oldScale    mgl32.Vec3
	oldRotation mgl32.Quat

	transformMatrix      mgl32.Mat4
	finalTransformMatrix *mgl32.Mat4

	Parent                  TweenableObject3D
	IgnoreParentRotation    bool
	IgnoreParentScale       bool
	oldParentTransform      mgl32.Mat4
	oldIgnoreParentRotation bool
	oldIgnoreParentScale    bool
}

func (tobj *TransformableObjectInstanced3D) valuesChanged() bool {
	return tobj.Position != tobj.oldPosition || tobj.Scale != tobj.oldScale || tobj.Rotation != tobj.oldRotation || tobj.IgnoreParentRotation != tobj.oldIgnoreParentRotation || tobj.IgnoreParentScale != tobj.oldIgnoreParentScale || tobj.getParentTransform() != tobj.oldParentTransform
}

func (tobj *TransformableObjectInstanced3D) CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int) {
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
		// T QR S
		T := mgl32.Translate3D(tobj.Position[0], tobj.Position[1], tobj.Position[2])
		QR := tobj.Rotation.Mat4()
		S := mgl32.Scale3D(tobj.Scale[0], tobj.Scale[1], tobj.Scale[2])
		tobj.transformMatrix = T.Mul4(QR).Mul4(S)

		// Parent
		ptransform := tobj.getParentTransform()
		tobj.transformMatrix = ptransform.Mul4(tobj.transformMatrix)

		tobj.oldPosition = tobj.Position
		tobj.oldScale = tobj.Scale
		tobj.oldRotation = tobj.Rotation
		tobj.oldParentTransform = ptransform
	}
	if cam3d != nil {
		*tobj.finalTransformMatrix = cam3d.GetInverseViewMatrix().Mul4(tobj.transformMatrix)
	} else {
		*tobj.finalTransformMatrix = tobj.transformMatrix
	}
}

func (tobj *TransformableObjectInstanced3D) GetTransformMatrix() mgl32.Mat4 {
	return *tobj.finalTransformMatrix
}

func (tobj *TransformableObjectInstanced3D) SetTransformMatrix(rmgr *RenderManager) {
	rmgr.setTransformMatrix3D(tobj.GetTransformMatrix())
}

func (tobj *TransformableObjectInstanced3D) getParentTransform() mgl32.Mat4 {
	if tobj.Parent != nil {
		if ptobj, ok := tobj.Parent.(TweenableObject3D); ok {
			if transform := ptobj.GetTransform3D(); transform != nil {
				if !tobj.IgnoreParentRotation && !tobj.IgnoreParentScale {
					var nrc int = -1
					if ent, ok := tobj.Parent.(*Entity3D); ok {
						nrc = ent.NotRelativeCamera()
					}
					transform.CalculateTransformMatrix(&RenderMgr, nrc)
					return transform.GetTransformMatrix()
				} else {
					// T QR S
					T := mgl32.Translate3D(transform.Position[0], transform.Position[1], transform.Position[2])
					var QR, S mgl32.Mat4
					if tobj.IgnoreParentRotation {
						QR = mgl32.Ident4()
					} else {
						QR = transform.Rotation.Mat4()
					}
					if tobj.IgnoreParentScale {
						S = mgl32.Ident4()
					} else {
						S = mgl32.Scale3D(transform.Scale[0], transform.Scale[1], transform.Scale[2])
					}
					pmat := transform.getParentTransform().Mul4(T.Mul4(QR).Mul4(S))
					return pmat
				}
			}
		}
	}

	return mgl32.Ident4()
}

func (tobj *TransformableObjectInstanced3D) GetPosition() mgl32.Vec3 {
	ptransform := tobj.getParentTransform()
	return ptransform.Mul4x1(tobj.Position.Vec4(1.0)).Vec3()
}

func (tobj *TransformableObjectInstanced3D) SetTransformMatrixPointer(tmp *mgl32.Mat4) {
	tobj.finalTransformMatrix = tmp
}

func DefaultTransformableObjectInstanced3D() *TransformableObjectInstanced3D {
	transform := TransformableObjectInstanced3D{
		Scale:    [3]float32{1.0, 1.0, 1.0},
		Rotation: mgl32.QuatRotate(0.0, [3]float32{0.0, 1.0, 0.0}),
	}

	return &transform
}
