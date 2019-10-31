package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

// A transform storing everything needed for the transform matrix
type TransformableObject3D struct {
	// The position in the world
	Position mgl32.Vec3
	// The scale the multiplies all vertices
	Scale    mgl32.Vec3
	// The rotation represented as a Quaternion
	Rotation mgl32.Quat

	oldPosition mgl32.Vec3
	oldScale    mgl32.Vec3
	oldRotation mgl32.Quat

	transformMatrix      mgl32.Mat4
	camNotRelativeMatrix mgl32.Mat4

	parent                  ParentObject3D
	parentChannel           chan bool
	childChannels           map[*TransformableObject3D]chan bool
	IgnoreParentRotation    bool
	IgnoreParentScale       bool
	oldParentTransform      mgl32.Mat4
	oldIgnoreParentRotation bool
	oldIgnoreParentScale    bool
}

func (tobj *TransformableObject3D) valuesChanged() bool {
	return tobj.Position != tobj.oldPosition || tobj.Scale != tobj.oldScale || tobj.Rotation != tobj.oldRotation || tobj.IgnoreParentRotation != tobj.oldIgnoreParentRotation || tobj.IgnoreParentScale != tobj.oldIgnoreParentScale || tobj.getParentTransform() != tobj.oldParentTransform
}

// Calculates the transform matrix
func (tobj *TransformableObject3D) CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int) {
	if tobj.parent != nil && RenderMgr.calculatingTransformMatricesParallel {
		<-tobj.parentChannel
	}
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
		tobj.camNotRelativeMatrix = cam3d.GetInverseViewMatrix().Mul4(tobj.transformMatrix)
	} else {
		tobj.camNotRelativeMatrix = tobj.transformMatrix
	}

	if RenderMgr.calculatingTransformMatricesParallel {
		if tobj.childChannels != nil {
			for _, ch := range tobj.childChannels {
				if ch != nil {
					ch <- true
				}
			}
		}
	}
}

// Returns the transform matrix that represents the transformation of this object
func (tobj *TransformableObject3D) GetTransformMatrix() mgl32.Mat4 {
	return tobj.camNotRelativeMatrix
}

// Sets the current transform matrix in the render manager
func (tobj *TransformableObject3D) SetTransformMatrix(rmgr *RenderManager) {
	rmgr.setTransformMatrix3D(tobj.GetTransformMatrix())
}

func (tobj *TransformableObject3D) getParentTransform() mgl32.Mat4 {
	if tobj.parent != nil {
		if ptobj, ok := tobj.parent.(ParentObject3D); ok {
			if transform := ptobj.GetTransform3D(); transform != nil {
				if !tobj.IgnoreParentRotation && !tobj.IgnoreParentScale {
					var nrc int = -1
					if ent, ok := tobj.parent.(*Entity3D); ok {
						nrc = ent.NotRelativeCamera()
					}
					if !RenderMgr.calculatingTransformMatricesParallel {
						transform.CalculateTransformMatrix(&RenderMgr, nrc)
					}
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

// The position based on the parent transform
func (tobj *TransformableObject3D) GetPosition() mgl32.Vec3 {
	ptransform := tobj.getParentTransform()
	return ptransform.Mul4x1(tobj.Position.Vec4(1.0)).Vec3()
}

// Returns itself
func (tobj *TransformableObject3D) GetTransform3D() *TransformableObject3D {
	return tobj
}

// Returns the parent of this transform
func (tobj *TransformableObject3D) GetParent() ParentObject3D {
	return tobj.parent
}

// Sets the parent of this transform to which this is relative to
func (tobj *TransformableObject3D) SetParent(parent ParentObject3D) {
	if parent == nil {
		if tobj.parent != nil {
			close(tobj.parentChannel)
			tobj.parent.SetChildChannel(nil, tobj)
			tobj.parent = nil
		}
		return
	}
	if tobj.parent != nil {
		close(tobj.parentChannel)
		tobj.parent.SetChildChannel(nil, tobj)
	}
	tobj.parent = parent
	tobj.parentChannel = make(chan bool)
	tobj.parent.SetChildChannel(tobj.parentChannel, tobj)
}

// Used for calculating the transform matrix
func (tobj *TransformableObject3D) SetChildChannel(channel chan bool, tobj1 *TransformableObject3D) {
	if tobj.childChannels == nil {
		tobj.childChannels = make(map[*TransformableObject3D]chan bool)
	}
	tobj.childChannels[tobj1] = channel
}

// Returns a transformable object with an identity matrix
func DefaultTransformableObject3D() *TransformableObject3D {
	transform := TransformableObject3D{
		Scale: [3]float32{1.0, 1.0, 1.0},
	}

	return &transform
}
