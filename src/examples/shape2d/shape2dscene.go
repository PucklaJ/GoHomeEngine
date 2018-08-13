package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"math"
)

type Shape2DScene struct {
	gohome.NilRenderObject
}

func (this *Shape2DScene) Init() {
	gohome.RenderMgr.AddObject(this)

	gohome.PointSize = 3.0
	gohome.Filled = false
}

var elapsedTime float64

func (this *Shape2DScene) Update(delta_time float32) {
	elapsedTime += float64(delta_time)
	gohome.DrawColor = gohome.Color{
		uint8(255.0 * (math.Sin(elapsedTime)/2.0 + 0.5)),
		255,
		255,
		255,
	}

	if gohome.InputMgr.JustPressed(gohome.MouseButtonLeft) {
		gohome.Filled = !gohome.Filled
	}

}

func (this *Shape2DScene) Render() {
	gohome.DrawPoint2D([2]float32{320.0/2.0 + 30.0, 240.0/2.0 + 30.0})
	gohome.DrawLine2D([2]float32{320.0 / 2.0, 240.0/2.0 + 100.0}, [2]float32{320.0/2.0 + 100.0, 240.0 / 2.0})
	gohome.DrawTriangle2D([2]float32{320.0, 240.0 + 120.0}, [2]float32{320.0 + 160.0, 240.0 + 120.0}, [2]float32{320.0, 240.0})

	gohome.DrawRectangle2D(
		[2]float32{320.0 / 4.0, 480.0 - 240.0/2.0},
		[2]float32{320.0 / 2.0, 480.0 - 240.0/2.0},
		[2]float32{320.0 / 2.0, 480.0 - 240.0},
		[2]float32{320.0 / 4.0, 480.0 - 240.0},
	)

	gohome.DrawCircle2D([2]float32{320.0 + 160.0, 240.0 - 120.0}, 100.0)

	gohome.DrawPolygon2D([2]float32{160.0 / 3.0, (480 - 120.0) / 3.0},
		[2]float32{320.0 / 3.0, 300.0 / 3.0},
		[2]float32{(320.0 + 160.0) / 3.0, (480.0 - 120.0) / 3.0},
		[2]float32{(320.0 + 160.0) / 3.0, 120.0 / 3.0},
		[2]float32{320.0 / 3.0, 180.0 / 3.0},
		[2]float32{160.0 / 3.0, 120.0 / 3.0})
}

func (this *Shape2DScene) Terminate() {
}
