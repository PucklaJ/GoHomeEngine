package gohome

import "github.com/go-gl/mathgl/mgl32"

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