package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Projection interface {
	CalculateProjectionMatrix()
	GetProjectionMatrix() mgl32.Mat4
	Update(newViewport Viewport)
}

type Ortho2DProjection struct {
	Left   float32
	Right  float32
	Bottom float32
	Top    float32

	oldLeft   float32
	oldRight  float32
	oldBottom float32
	oldTop    float32

	projectionMatrix mgl32.Mat4
}

func (o2Dp *Ortho2DProjection) valuesChanged() bool {
	return o2Dp.Left != o2Dp.oldLeft || o2Dp.Right != o2Dp.oldRight || o2Dp.Bottom != o2Dp.oldBottom || o2Dp.Top != o2Dp.oldTop
}

func (o2Dp *Ortho2DProjection) CalculateProjectionMatrix() {
	if o2Dp.valuesChanged() {
		o2Dp.projectionMatrix = mgl32.Ortho2D(o2Dp.Left, o2Dp.Right, o2Dp.Bottom, o2Dp.Top)
	} else {
		return
	}

	o2Dp.oldLeft = o2Dp.Left
	o2Dp.oldRight = o2Dp.Right
	o2Dp.oldBottom = o2Dp.Bottom
	o2Dp.oldTop = o2Dp.Top
}

func (o2Dp *Ortho2DProjection) GetProjectionMatrix() mgl32.Mat4 {
	return o2Dp.projectionMatrix
}

func (o2Dp *Ortho2DProjection) Update(newViewport Viewport) {
	o2Dp.Left = 0.0
	o2Dp.Right = float32(newViewport.Width)
	o2Dp.Top = 0.0
	o2Dp.Bottom = float32(newViewport.Height)
	o2Dp.CalculateProjectionMatrix()
}

type PerspectiveProjection struct {
	Width     float32
	Height    float32
	FOV       float32
	NearPlane float32
	FarPlane  float32

	oldWidth     float32
	oldHeight    float32
	oldFOV       float32
	oldNearPlane float32
	oldFarPlane  float32

	projectionMatrix mgl32.Mat4
}

func (pp *PerspectiveProjection) valuesChanged() bool {
	return pp.Width != pp.oldWidth || pp.Height != pp.oldHeight || pp.FOV != pp.oldFOV || pp.NearPlane != pp.oldNearPlane || pp.FarPlane != pp.oldFarPlane
}

func (pp *PerspectiveProjection) CalculateProjectionMatrix() {
	if pp.valuesChanged() {
		pp.projectionMatrix = mgl32.Perspective(mgl32.DegToRad(pp.FOV), pp.Width/pp.Height, pp.NearPlane, pp.FarPlane)
	} else {
		return
	}

	pp.oldWidth = pp.Width
	pp.oldHeight = pp.Height
	pp.oldFOV = pp.FOV
	pp.oldNearPlane = pp.NearPlane
	pp.oldFarPlane = pp.FarPlane
}

func (pp *PerspectiveProjection) Update(newViewport Viewport) {
	pp.Width = float32(newViewport.Width)
	pp.Height = float32(newViewport.Height)
	pp.CalculateProjectionMatrix()
}

func (pp *PerspectiveProjection) GetProjectionMatrix() mgl32.Mat4 {
	return pp.projectionMatrix
}

type IdentityProjection struct {
}

func (IdentityProjection) CalculateProjectionMatrix() {

}

func (IdentityProjection) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Ident4()
}

func (IdentityProjection) Update(newViewport Viewport) {

}
