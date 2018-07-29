package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type Tweens3DScene struct {
	gopher gohome.Entity3D
	tset   gohome.Tweenset
	cam    gohome.Camera3D
}

func (this *Tweens3DScene) Init() {
	gohome.Init3DShaders()

	gohome.ResourceMgr.LoadLevel("Gopher","gopher.obj",true)

	this.gopher.InitLevel(gohome.ResourceMgr.GetLevel("Gopher"))

	this.gopher.Transform.Position = [3]float32{-5.0, 0.0, 0.0}
	this.gopher.Transform.Scale = [3]float32{0.75, 0.75, 0.75}
	this.tset = gohome.Tweenset{
		Tweens: []gohome.Tween{
			&gohome.TweenBlink{Amount:80,Time:20.0,TweenType:gohome.TWEEN_TYPE_ALWAYS},
			&gohome.TweenRotation3D{Destination:[3]float32{0.0,360.0,0.0},Time:20.0,TweenType:gohome.TWEEN_TYPE_ALWAYS},
			&gohome.TweenPosition3D{Destination:[3]float32{5.0,0.0,0.0},Time:10.0,TweenType:gohome.TWEEN_TYPE_AFTER_PREVIOUS},
			&gohome.TweenPosition3D{Destination:[3]float32{-5.0,0.0,0.0},Time:10.0,TweenType:gohome.TWEEN_TYPE_AFTER_PREVIOUS},
		},
		Loop: true,
	}
	this.tset.SetParent(&this.gopher)
	this.tset.Start()
	gohome.UpdateMgr.AddObject(&this.tset)
	gohome.RenderMgr.AddObject(&this.gopher)

	gohome.LightMgr.DisableLighting()

	this.cam.Init()
	this.cam.Position[2] = 10.0
	gohome.RenderMgr.SetCamera3D(&this.cam,0)
}

func (this *Tweens3DScene) Update(delta_time float32) {

}

func (this *Tweens3DScene) Terminate() {

}