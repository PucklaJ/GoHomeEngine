package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type TransformableObjectInstanced3D struct {
	TransformableObject3D

	finalTransformMatrix *mgl32.Mat4
}

func (tobj *TransformableObjectInstanced3D) CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int) {
	tobj.TransformableObject3D.CalculateTransformMatrix(rmgr, notRelativeToCamera)
	*tobj.finalTransformMatrix = tobj.camNotRelativeMatrix
}

func (tobj *TransformableObjectInstanced3D) GetTransformMatrix() mgl32.Mat4 {
	return *tobj.finalTransformMatrix
}

func (tobj *TransformableObjectInstanced3D) SetTransformMatrixPointer(tmp *mgl32.Mat4) {
	tobj.finalTransformMatrix = tmp
}
