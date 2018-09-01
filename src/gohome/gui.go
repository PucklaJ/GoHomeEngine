package gohome

import (
	"image/color"
)

type ButtonCallback func(btn *Button)

var focusedButton *Button

type Button struct {
	Sprite2D

	PressCallback ButtonCallback
	EnterCallback ButtonCallback
	LeaveCallback ButtonCallback
	Entered       bool

	EnterModColor color.Color
	PressModColor color.Color
}

func (this *Button) Init(pos [2]float32, texture string) {
	this.Sprite2D.Init(texture)
	if this.Texture == nil {
		rt := Render.CreateRenderTexture("ButtonTexture", 200, 100, 1, false, false, false, false)
		rt.SetAsTarget()
		Render.ClearScreen(Color{255, 255, 255, 255})
		rt.UnsetAsTarget()
		this.Texture = rt
		this.Transform.Size[0], this.Transform.Size[1] = 200, 100
	}
	this.Transform.Position = pos
	this.NotRelativeToCamera = 0

	UpdateMgr.AddObject(this)
	RenderMgr.AddObject(this)

	this.EnterModColor = Color{200, 200, 200, 255}
	this.PressModColor = Color{128, 128, 128, 255}
}

func (this *Button) Update(delta_time float32) {
	if focusedButton == nil || focusedButton == this {
		mpos := [2]float32{float32(InputMgr.Mouse.Pos[0]), float32(InputMgr.Mouse.Pos[1])}
		size := this.Transform.Size
		size[0] *= this.Transform.Scale[0]
		size[1] *= this.Transform.Scale[1]
		pos := this.Transform.Position
		pos[0] -= this.Transform.Origin[0] * size[0]
		pos[1] -= this.Transform.Origin[1] * size[1]
		prevEntered := this.Entered
		this.Entered = mpos[0] > pos[0] && mpos[1] > pos[1] && mpos[0] <= pos[0]+size[0] && mpos[1] <= pos[1]+size[1]
		if !InputMgr.IsPressed(MouseButtonLeft) {
			if this.Entered {
				this.Texture.SetModColor(this.EnterModColor)
				if !prevEntered && this.EnterCallback != nil {
					this.EnterCallback(this)
				}
				if InputMgr.WasPressed(MouseButtonLeft) {
					if this.PressCallback != nil && focusedButton == this {
						this.PressCallback(this)
					}
				} else if focusedButton == nil {
					focusedButton = this
				}
			} else {
				this.Texture.SetModColor(nil)
				if focusedButton == this {
					focusedButton = nil
				}
				if prevEntered && this.LeaveCallback != nil {
					this.LeaveCallback(this)
				}
			}
		} else if this.Entered {
			this.Texture.SetModColor(this.PressModColor)
		}
	}
}

func (this *Button) Terminate() {
	UpdateMgr.RemoveObject(this)
	RenderMgr.RemoveObject(this)
}
