package gohome

const (
	TWEEN_TYPE_WITH_PREVIOUS  uint8 = iota
	TWEEN_TYPE_AFTER_PREVIOUS uint8 = iota
	TWEEN_TYPE_ALWAYS uint8 = iota
)

type Tween interface {
	Start(parent TweenableObject)
	Update(delta_time float32) bool
	GetType() uint8
	End()
	Reset()
}

type TweenableObject interface {
	SetTweenset(set Tweenset)
	StartTweens()
	StopTweens()
	PauseTweens()
	ResumeTweens()
}

type TweenableObject2D interface {
	GetTransform2D() *TransformableObject2D
}

type TweenableObject3D interface {
	GetTransform3D() *TransformableObject3D
}