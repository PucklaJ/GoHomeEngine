package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

type TweensScene struct {
	gopher gohome.Sprite2D
	tset   gohome.Tweenset
	cam    gohome.Camera2D
}

func (this *TweensScene) Init() {
	gohome.Init2DShaders()

	gohome.ResourceMgr.LoadTexture("Gopher", "gopher.png")

	this.gopher.Init("Gopher")
	this.gopher.Transform.Origin = [2]float32{0.5, 0.5}
	nr := gohome.Render.GetNativeResolution()
	this.gopher.Transform.Position = [2]float32{0.0, 0.0}
	this.gopher.Transform.Scale = [2]float32{0.8, 0.8}
	this.tset = gohome.Tweenset{
		Tweens: []gohome.Tween{
			&gohome.TweenBlink{Amount: 80, Time: 20.0, TweenType: gohome.TWEEN_TYPE_ALWAYS},
			&gohome.TweenPosition2D{Destination: [2]float32{nr[0], 0.0}, Time: 5.0, TweenType: gohome.TWEEN_TYPE_AFTER_PREVIOUS},
			&gohome.TweenScale2D{Destination: [2]float32{0.8 * 2.0, 0.8 * 2.0}, Time: 5.0, TweenType: gohome.TWEEN_TYPE_WITH_PREVIOUS},
			&gohome.TweenPosition2D{Destination: nr, Time: 5.0, TweenType: gohome.TWEEN_TYPE_AFTER_PREVIOUS},
			&gohome.TweenRotation2D{Destination: 180.0, Time: 5.0, TweenType: gohome.TWEEN_TYPE_WITH_PREVIOUS},
			&gohome.TweenPosition2D{Destination: [2]float32{0.0, nr[1]}, Time: 5.0, TweenType: gohome.TWEEN_TYPE_AFTER_PREVIOUS},
			&gohome.TweenPosition2D{Destination: [2]float32{0.0, 0.0}, Time: 5.0, TweenType: gohome.TWEEN_TYPE_AFTER_PREVIOUS},
			&gohome.TweenScale2D{Destination: [2]float32{0.8, 0.8}, Time: 5.0, TweenType: gohome.TWEEN_TYPE_WITH_PREVIOUS},
			&gohome.TweenRotation2D{Destination: 360.0, Time: 5.0, TweenType: gohome.TWEEN_TYPE_WITH_PREVIOUS},
		},
		Loop: true,
	}
	this.tset.SetParent(&this.gopher)
	this.tset.Start()
	gohome.UpdateMgr.AddObject(&this.tset)
	gohome.RenderMgr.AddObject(&this.gopher)

	this.cam.Zoom = 0.5
	this.cam.Position = nr.Div(-2.0)
	gohome.RenderMgr.SetCamera2D(&this.cam, 0)
}

func (this *TweensScene) Update(delta_time float32) {

}

func (this *TweensScene) Terminate() {

}
