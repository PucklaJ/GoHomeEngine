package gohome

// These constants define when the Tween will be applied
const (
	// The tween starts with the previous tween
	TWEEN_TYPE_WITH_PREVIOUS  = iota
	// The tween starts after the previous tween has finished
	TWEEN_TYPE_AFTER_PREVIOUS = iota
	// The tween starts as the first tween in the list starts
	TWEEN_TYPE_ALWAYS         = iota
)

// Essentially an animation
type Tween interface {
	// Starts the tween on a given object (Sprite2D, Entity3D, etc.)
	Start(parent interface{})
	// Updates the tween and returns wether the tween has finished
	Update(delta_time float32) bool
	// Returns the type of this tween
	GetType() uint8
	// Is called when the tween has finished
	End()
	// Is called when the whole TweenSet is reset
	Reset()
	// Creates a copy of this tween
	Copy() Tween
}

// An object that can be used by tweens
type TweenableObject2D interface {
	GetTransform2D() *TransformableObject2D
}

// An object that can be used by tweens
type TweenableObject3D interface {
	GetTransform3D() *TransformableObject3D
}

// An object that is a parent of another object
type ParentObject3D interface {
	TweenableObject3D
	SetChildChannel(channel chan bool, tobj *TransformableObject3D)
}
