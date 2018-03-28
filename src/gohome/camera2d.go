package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera2D struct {
	Position mgl32.Vec2
	Zoom     float32
	Rotation float32
	Origin   mgl32.Vec2

	oldPosition mgl32.Vec2
	oldZoom     float32
	oldRotation float32

	viewMatrix        mgl32.Mat3
	inverseViewMatrix mgl32.Mat3
}

func (cam *Camera2D) valuesChanged() bool {
	return cam.Position != cam.oldPosition || cam.Zoom != cam.oldZoom || cam.Rotation != cam.oldRotation
}

func (cam *Camera2D) CalculateViewMatrix() {
	if cam.valuesChanged() {
		// -OT S R OT T
		windowSize := Framew.WindowGetSize()
		ot := mgl32.Translate2D(-windowSize[0]*cam.Origin[0], -windowSize[1]*cam.Origin[1])
		not := mgl32.Translate2D(windowSize[0]*cam.Origin[0], windowSize[1]*cam.Origin[1])
		cam.viewMatrix = not.Mul3(mgl32.Scale2D(cam.Zoom, cam.Zoom)).Mul3(mgl32.Rotate2D(mgl32.DegToRad(cam.Rotation)).Mat3()).Mul3(ot).Mul3(mgl32.Translate2D(-cam.Position[0], -cam.Position[1]))
		cam.inverseViewMatrix = cam.viewMatrix.Inv()
	} else {
		return
	}

	cam.oldPosition = cam.Position
	cam.oldZoom = cam.Zoom
	cam.oldRotation = cam.Rotation
}

func (cam *Camera2D) GetViewMatrix() mgl32.Mat3 {
	return cam.viewMatrix
}

func (cam *Camera2D) GetInverseViewMatrix() mgl32.Mat3 {
	return cam.inverseViewMatrix
}

func (cam *Camera2D) AddPositionRotated(pos mgl32.Vec2) {
	mat := mgl32.Rotate2D(mgl32.DegToRad(-cam.Rotation))
	x := mat.At(0, 0)*pos[0] + mat.At(0, 1)*pos[1]
	y := mat.At(1, 0)*pos[0] + mat.At(1, 1)*pos[1]
	pos[0] = x
	pos[1] = y

	cam.Position[0] += pos[0]
	cam.Position[1] += pos[1]
}
