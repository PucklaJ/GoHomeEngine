package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"math"
)

const (
	NEAR_LEFT_DOWN  = 0
	NEAR_RIGHT_DOWN = 1
	NEAR_RIGHT_UP   = 2
	NEAR_LEFT_UP    = 3
	FAR_LEFT_DOWN   = 4
	FAR_RIGHT_DOWN  = 5
	FAR_RIGHT_UP    = 6
	FAR_LEFT_UP     = 7
)

// A projection that is used to correctly display the objects on the screen
type Projection interface {
	// Calculates the projection matrix used in the shader
	CalculateProjectionMatrix()
	// Returns the projection matrix
	GetProjectionMatrix() mgl32.Mat4
	// Updates the projection with the new viewport
	Update(newViewport Viewport)
	// Returns the view frustum
	GetFrustum() [8]mgl32.Vec3
}

// A 2-dimensional orthogonal projection
type Ortho2DProjection struct {
	// The left most value
	Left   float32
	// The right most value
	Right  float32
	// The bottom mose value
	Bottom float32
	// The top most value
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

func (o2Dp *Ortho2DProjection) GetFrustum() [8]mgl32.Vec3 {
	return [8]mgl32.Vec3{}
}

// A 3-dimensional perspective projection
type PerspectiveProjection struct {
	// The width of the viewport
	Width     float32
	// The height of the viewport
	Height    float32
	// The field of view
	FOV       float32
	// The near plane at which objects start to clip away
	NearPlane float32
	// The far plane at which the objects start to clip away
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

func (pp *PerspectiveProjection) GetFrustum() [8]mgl32.Vec3 {
	var farPlaneHalfWidth, nearPlaneHalfWidth float32
	var farPlaneHalfHeight, nearPlaneHalfHeight float32
	var centerFarPlane, centerNearPlane mgl32.Vec3
	var points [8]mgl32.Vec3

	forward := mgl32.Vec3{0.0, 0.0, -1.0}
	up := mgl32.Vec3{0.0, 1.0, 0.0}
	down := up.Mul(-1.0)
	right := mgl32.Vec3{1.0, 0.0, 0.0}
	left := right.Mul(-1.0)

	farPlaneHalfWidth = float32(math.Tan(float64(pp.FOV)/180.0*math.Pi) * float64(pp.FarPlane))
	nearPlaneHalfWidth = float32(math.Tan(float64(pp.FOV)/180.0*math.Pi) * float64(pp.NearPlane))
	farPlaneHalfHeight = farPlaneHalfWidth / (pp.Width / pp.Height)
	nearPlaneHalfHeight = nearPlaneHalfWidth / (pp.Width / pp.Height)

	centerFarPlane = forward.Mul(pp.FarPlane)
	centerNearPlane = forward.Mul(pp.NearPlane)

	points[NEAR_LEFT_DOWN] = centerNearPlane.Add(left.Mul(nearPlaneHalfWidth)).Add(down.Mul(nearPlaneHalfHeight))
	points[NEAR_RIGHT_DOWN] = centerNearPlane.Add(right.Mul(nearPlaneHalfWidth)).Add(down.Mul(nearPlaneHalfHeight))
	points[NEAR_RIGHT_UP] = centerNearPlane.Add(right.Mul(nearPlaneHalfWidth)).Add(up.Mul(nearPlaneHalfHeight))
	points[NEAR_LEFT_UP] = centerNearPlane.Add(left.Mul(nearPlaneHalfWidth)).Add(up.Mul(nearPlaneHalfHeight))
	points[FAR_LEFT_DOWN] = centerFarPlane.Add(left.Mul(farPlaneHalfWidth)).Add(down.Mul(farPlaneHalfHeight))
	points[FAR_RIGHT_DOWN] = centerFarPlane.Add(right.Mul(farPlaneHalfWidth)).Add(down.Mul(farPlaneHalfHeight))
	points[FAR_RIGHT_UP] = centerFarPlane.Add(right.Mul(farPlaneHalfWidth)).Add(up.Mul(farPlaneHalfHeight))
	points[FAR_LEFT_UP] = centerFarPlane.Add(left.Mul(farPlaneHalfWidth)).Add(up.Mul(farPlaneHalfHeight))

	return points
}

// A projection doing nothing
type IdentityProjection struct {
}

func (IdentityProjection) CalculateProjectionMatrix() {

}

func (IdentityProjection) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Ident4()
}

func (IdentityProjection) Update(newViewport Viewport) {

}

func (IdentityProjection) GetFrustum() [8]mgl32.Vec3 {
	return [8]mgl32.Vec3{}
}

// A 3-dimensional orthogonal projection
type Ortho3DProjection struct {
	// The left most point
	Left   float32
	// The right most point
	Right  float32
	// The bottom most point
	Bottom float32
	// The top most point
	Top    float32
	// The near plane
	Near   float32
	// The far plane
	Far    float32

	oldLeft   float32
	oldRight  float32
	oldBottom float32
	oldTop    float32
	oldNear   float32
	oldFar    float32

	projectionMatrix mgl32.Mat4
}

func (this *Ortho3DProjection) valuesChanged() bool {
	return this.Left != this.oldLeft || this.Right != this.oldRight || this.Bottom != this.oldBottom || this.Top != this.oldTop || this.Near != this.oldNear || this.Far != this.oldFar
}

func (this *Ortho3DProjection) CalculateProjectionMatrix() {
	if this.valuesChanged() {
		this.projectionMatrix = mgl32.Ortho(this.Left, this.Right, this.Bottom, this.Top, this.Near, this.Far)
	} else {
		return
	}

	this.oldLeft = this.Left
	this.oldRight = this.Right
	this.oldBottom = this.Bottom
	this.oldTop = this.Top
	this.oldNear = this.Near
	this.oldFar = this.Far
}

func (this *Ortho3DProjection) GetProjectionMatrix() mgl32.Mat4 {
	return this.projectionMatrix
}

func (this *Ortho3DProjection) Update(newViewport Viewport) {
	this.Left = 0.0
	this.Right = float32(newViewport.Width)
	this.Top = 0.0
	this.Bottom = float32(newViewport.Height)
	this.CalculateProjectionMatrix()
}

func (this *Ortho3DProjection) GetFrustum() [8]mgl32.Vec3 {
	return [8]mgl32.Vec3{}
}
