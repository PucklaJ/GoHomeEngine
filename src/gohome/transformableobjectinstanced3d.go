package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

// A transformable object used for instancing
type TransformableObjectInstanced3D struct {
	TransformableObject3D

	finalTransformMatrix *mgl32.Mat4
}

// Calculates the transform matrix
func (tobj *TransformableObjectInstanced3D) CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int) {
	tobj.TransformableObject3D.CalculateTransformMatrix(rmgr, notRelativeToCamera)
	*tobj.finalTransformMatrix = tobj.camNotRelativeMatrix
}

// Returns the transform matrix that represents this transformation
func (tobj *TransformableObjectInstanced3D) GetTransformMatrix() mgl32.Mat4 {
	return *tobj.finalTransformMatrix
}

// Sets a pointer to a matrix so that it can be used for instancing
func (tobj *TransformableObjectInstanced3D) SetTransformMatrixPointer(tmp *mgl32.Mat4) {
	tobj.finalTransformMatrix = tmp
}
