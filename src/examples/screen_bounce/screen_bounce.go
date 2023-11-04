package main

import (
	"math"
	"math/rand"

	"github.com/PucklaJ/GoHomeEngine/src/audio"
	framework "github.com/PucklaJ/GoHomeEngine/src/frameworks/SAMURE"
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	renderer "github.com/PucklaJ/GoHomeEngine/src/renderers/OpenGL"
	"github.com/PucklaJ/mathgl/mgl32"
)

const (
	Speed = 300.0

	MinX = 0.0
	MinY = 0.0
	MaxX = 1920.0 * 2.0
	MaxY = 1080.0
)

type ScreenBounceScene struct {
	fuzzy  gohome.Sprite2D
	bounce gohome.Sound

	DirX, DirY float32
	radius     float32
}

func (s *ScreenBounceScene) Init() {
	audio.InitAudio()

	gohome.ResourceMgr.LoadTexture("fuzzy", "Ten 13.png")
	gohome.ResourceMgr.LoadSound("bounce", "BounceYoFrankie.wav")

	s.fuzzy.Init("fuzzy")
	s.bounce = gohome.ResourceMgr.GetSound("bounce")

	gohome.RenderMgr.AddObject(&s.fuzzy)

	s.fuzzy.Transform.Scale[0] = 0.3
	s.fuzzy.Transform.Scale[1] = 0.3
	s.radius = s.fuzzy.Transform.Scale[0] * s.fuzzy.Transform.Size[0] * 0.5

	s.fuzzy.Transform.Origin = mgl32.Vec2{0.5, 0.5}
	s.fuzzy.Transform.Position[0] = float32(math.Min(math.Max(float64(rand.Float32()*MaxX), float64(s.radius)), MaxX+float64(s.radius)))
	s.fuzzy.Transform.Position[1] = float32(math.Min(math.Max(float64(rand.Float32()*MaxY), float64(s.radius)), MaxY+float64(s.radius)))

	dir := mgl32.Vec2{
		(rand.Float32() * 2.0) - 1.0,
		(rand.Float32() * 2.0) - 1.0,
	}.Normalize()

	s.DirX = dir.X()
	s.DirY = dir.Y()
}

func (s *ScreenBounceScene) Update(delta_time float32) {
	s.fuzzy.Transform.Rotation += delta_time * Speed
	s.fuzzy.Transform.Position[0] += s.DirX * delta_time * Speed
	s.fuzzy.Transform.Position[1] += s.DirY * delta_time * Speed

	if (s.fuzzy.Transform.Position[0] + s.radius) > MaxX {
		s.fuzzy.Transform.Position[0] = MaxX - s.radius
		s.DirX *= -1.0
		s.bounce.Play(false)
	} else if s.fuzzy.Transform.Position[0]-s.radius < 0.0 {
		s.fuzzy.Transform.Position[0] = s.radius
		s.DirX *= -1.0
		s.bounce.Play(false)
	}

	if (s.fuzzy.Transform.Position[1] + s.radius) > MaxY {
		s.fuzzy.Transform.Position[1] = MaxY - s.radius
		s.DirY *= -1.0
		s.bounce.Play(false)
	} else if s.fuzzy.Transform.Position[1]-s.radius < 0.0 {
		s.fuzzy.Transform.Position[1] = s.radius
		s.DirY *= -1.0
		s.bounce.Play(false)
	}
}

func (s *ScreenBounceScene) Terminate() {
}

func main() {
	gohome.MainLop.Run(&framework.SAMUREFramework{}, &renderer.OpenGLRenderer{}, 1920, 1080, "ScreenBounce", &ScreenBounceScene{})
}
