package main

import "github.com/PucklaMotzer09/GoHomeEngine/src/gohome"

type MultipleTouchesScene struct {
	targets [6]gohome.Sprite2D
}

func (this *MultipleTouchesScene) Init() {
	gohome.Init2DShaders()

	gohome.ResourceMgr.LoadTexture("Target", "target.png")
	gohome.ResourceMgr.GetTexture("Target").SetFiltering(gohome.FILTERING_NEAREST)

	for i := 0; i < 6; i++ {
		this.targets[i].Init("Target")
		this.targets[i].Transform.Scale[0] = 3.0
		this.targets[i].Transform.Scale[1] = 3.0
		this.targets[i].Transform.Origin = [2]float32{0.5, 0.5}
		this.targets[i].Visible = false
		gohome.RenderMgr.AddObject(&this.targets[i])
	}
}

func (this *MultipleTouchesScene) Update(delta_time float32) {
	for i := 0; i < 6; i++ {
		if gohome.InputMgr.IsTouched(uint8(i)) {
			this.targets[i].Transform.Position[0] = float32(gohome.InputMgr.Touches[uint8(i)].Pos[0])
			this.targets[i].Transform.Position[1] = float32(gohome.InputMgr.Touches[uint8(i)].Pos[1])
			this.targets[i].Visible = true
		} else {
			this.targets[i].Visible = false
		}
	}
}

func (this *MultipleTouchesScene) Terminate() {

}
