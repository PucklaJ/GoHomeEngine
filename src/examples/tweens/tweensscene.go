package main

import "github.com/PucklaMotzer09/gohomeengine/src/gohome"

type TweensScene struct {
	gopher gohome.Sprite2D
}

func (this *TweensScene) Init() {
	gohome.Init2DShaders()

	gohome.ResourceMgr.LoadTexture("Gopher","gopher.png")

	this.gopher.Init("Gopher")
	this.gopher.Transform.Origin = [2]float32{0.5,0.5}
	nw,nh := gohome.Render.GetNativeResolution()
	this.gopher.Transform.Position = [2]float32{0.0,0.0}
	this.gopher.Transform.Scale = [2]float32{0.5,0.5}

	this.gopher.SetTweenset(gohome.Tweenset{
		Tweens: []gohome.Tween{
			&gohome.TweenPosition2D{Destination:[2]float32{float32(nw),0.0},Time:5.0,TweenType:gohome.TWEEN_TYPE_AFTER_PREVIOUS},
			&gohome.TweenPosition2D{Destination:[2]float32{float32(nw),float32(nh)},Time:5.0,TweenType:gohome.TWEEN_TYPE_AFTER_PREVIOUS,},
			&gohome.TweenPosition2D{Destination:[2]float32{0.0,float32(nh)},Time:5.0,TweenType:gohome.TWEEN_TYPE_AFTER_PREVIOUS,},
			&gohome.TweenPosition2D{Destination:[2]float32{0.0,0.0},Time:5.0,TweenType:gohome.TWEEN_TYPE_AFTER_PREVIOUS,},
		},
		Loop: true,
	})
	this.gopher.StartTweens()
	gohome.RenderMgr.AddObject(&this.gopher)
}

func (this *TweensScene) Update(delta_time float32) {

}

func (this *TweensScene) Terminate() {

}
