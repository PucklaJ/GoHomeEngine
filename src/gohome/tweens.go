package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
)

type TweenPosition2D struct {
	Destination mgl32.Vec2
	Time float32
	TweenType uint8

	transform *TransformableObject2D
	velocity mgl32.Vec2
	elapsedTime float32
}

func (this *TweenPosition2D) Start(parent interface{}) {
	parent2D,ok := parent.(TweenableObject2D)
	if ok {
		this.transform = parent2D.GetTransform2D()
	} else {
		this.transform = nil
	}

	if this.transform != nil {
		this.velocity = this.Destination.Sub(this.transform.Position).Mul(1.0/this.Time)
	}
	this.elapsedTime = 0.0
}
func (this *TweenPosition2D) Update(delta_time float32) bool {
	if this.transform == nil {
		return true
	}
	this.elapsedTime += delta_time

	this.transform.Position = this.transform.Position.Add(this.velocity.Mul(delta_time))

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenPosition2D) GetType() uint8 {
	return this.TweenType
}
func (this *TweenPosition2D) End() {
	if this.transform != nil {
		this.transform.Position = this.Destination
	}
}
func (this *TweenPosition2D) Reset() {
	if this.transform != nil {
		this.transform.Position = this.Destination.Sub(this.velocity.Mul(this.Time))
		this.elapsedTime = 0.0
	}
}

type TweenRotation2D struct {
	Destination float32
	Time float32
	TweenType uint8

	transform *TransformableObject2D
	velocity float32
	elapsedTime float32
}

func (this *TweenRotation2D) Start(parent interface{}) {
	parent2D,ok := parent.(TweenableObject2D)
	if ok {
		this.transform = parent2D.GetTransform2D()
	} else {
		this.transform = nil
	}

	if this.transform != nil {
		this.velocity = (this.Destination - this.transform.Rotation)  / this.Time
	}
	this.elapsedTime = 0.0
}
func (this *TweenRotation2D) Update(delta_time float32) bool {
	if this.transform == nil {
		return true
	}
	this.elapsedTime += delta_time

	this.transform.Rotation += this.velocity * delta_time

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenRotation2D) GetType() uint8 {
	return this.TweenType
}
func (this *TweenRotation2D) End() {
	if this.transform != nil {
		this.transform.Rotation = this.Destination
	}
}
func (this *TweenRotation2D) Reset() {
	if this.transform != nil {
		this.transform.Rotation = this.Destination - this.velocity * this.Time
		this.elapsedTime = 0.0
	}
}

type TweenWait struct {
	Time float32
	TweenType uint8

	elapsedTime float32
}

func (this *TweenWait) Start(parent interface{}) {
	this.elapsedTime = 0.0
}
func (this *TweenWait) Update(delta_time float32) bool {
	this.elapsedTime += delta_time

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenWait) GetType() uint8 {
	return this.TweenType
}
func (this *TweenWait) End() {

}
func (this *TweenWait) Reset() {
	this.elapsedTime = 0.0
}

type BlinkableObject interface {
	SetVisible()
	SetInvisible()
	IsVisible() bool
}

type TweenBlink struct {
	Amount uint32
	Time float32
	TweenType uint8

	timeForOneBlink float32
	elapsedTime float32
	oneBlinkElapsedTime float32
	previousVisible bool
	parent BlinkableObject
}

func (this *TweenBlink) Start(parent interface{}) {
	this.elapsedTime = 0.0
	this.timeForOneBlink = this.Time/float32(this.Amount)
	if parent != nil {
		this.parent = parent.(BlinkableObject)
		if this.parent != nil {
			this.previousVisible = this.parent.IsVisible()
		}
	}
}
func (this *TweenBlink) Update(delta_time float32) bool {
	if this.parent == nil {
		return true
	}

	this.elapsedTime += delta_time
	this.oneBlinkElapsedTime += delta_time

	if this.oneBlinkElapsedTime >= this.timeForOneBlink/2.0 {
		if this.parent.IsVisible() {
			this.parent.SetInvisible()
		} else {
			this.parent.SetVisible()
		}
		this.oneBlinkElapsedTime = 0.0
	}

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenBlink) GetType() uint8 {
	return this.TweenType
}
func (this *TweenBlink) End() {
	if this.parent != nil {
		if this.previousVisible {
			this.parent.SetVisible()
		} else {
			this.parent.SetInvisible()
		}
	}
}
func (this *TweenBlink) Reset() {
	this.elapsedTime = 0.0
	this.oneBlinkElapsedTime = 0.0
	if this.parent != nil {
		if this.previousVisible {
			this.parent.SetVisible()
		} else {
			this.parent.SetInvisible()
		}
	}
}

