package main

import (
	"math"
	"math/rand"

	framework "github.com/PucklaJ/GoHomeEngine/src/frameworks/SAMURE"
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	renderer "github.com/PucklaJ/GoHomeEngine/src/renderers/OpenGL"
	"github.com/PucklaJ/mathgl/mgl32"
	"golang.org/x/image/colornames"
)

const (
	Radius = 100.0
	Speed  = 500.0

	MinX = 0.0
	MinY = 0.0
	MaxX = 1920.0 * 2.0
	MaxY = 1080.0
)

type ScreenBounceScene struct {
	gohome.NilRenderObject

	X, Y       float32
	DirX, DirY float32
}

func (s *ScreenBounceScene) Init() {
	gohome.FPSLimit.MaxFPS = math.MaxInt
	gohome.RenderMgr.EnableBackBuffer = false
	gohome.Render.SetBackgroundColor(gohome.Color{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	})

	gohome.RenderMgr.AddObject(s)

	s.X = 0.0
	s.Y = 0.0

	s.X = float32(math.Min(math.Max(float64(rand.Float32()*MaxX), float64(Radius)), MaxX+Radius))
	s.Y = float32(math.Min(math.Max(float64(rand.Float32()*MaxY), float64(Radius)), MaxY+Radius))

	dir := mgl32.Vec2{
		(rand.Float32() * 2.0) - 1.0,
		(rand.Float32() * 2.0) - 1.0,
	}.Normalize()

	s.DirX = dir.X()
	s.DirY = dir.Y()
}

func (s *ScreenBounceScene) Update(delta_time float32) {
	s.X += s.DirX * delta_time * Speed
	s.Y += s.DirY * delta_time * Speed

	if (s.X + Radius) > MaxX {
		s.X = MaxX - Radius
		s.DirX *= -1.0
	} else if s.X-Radius < 0.0 {
		s.X = Radius
		s.DirX *= -1.0
	}

	if (s.Y + Radius) > MaxY {
		s.Y = MaxY - Radius
		s.DirY *= -1.0
	} else if s.Y-Radius < 0.0 {
		s.Y = Radius
		s.DirY *= -1.0
	}
}

func (s *ScreenBounceScene) Render() {
	samure := gohome.Framew.(*framework.SAMUREFramework)

	x, y := samure.CurrentOutputGeo.RelX(float64(s.X)), samure.CurrentOutputGeo.RelY(float64(s.Y))

	gohome.Filled = true
	gohome.DrawColor = colornames.Lime
	gohome.DrawCircle2D(mgl32.Vec2{float32(x), float32(y)}, Radius)

	gohome.Filled = false
	gohome.LineWidth = 20.0
	gohome.DrawColor = colornames.Blue
	gohome.DrawCircle2D(mgl32.Vec2{float32(x), float32(y)}, Radius)
}

func (s *ScreenBounceScene) Terminate() {
}

func main() {
	gohome.MainLop.Run(&framework.SAMUREFramework{}, &renderer.OpenGLRenderer{}, 1920, 1080, "ScreenBounce", &ScreenBounceScene{})
}
