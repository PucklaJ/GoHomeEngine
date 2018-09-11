package gohome

import (
	"image/color"
	"math"
)

type ButtonCallback func(btn *Button)

var focusedButton *Button
var ButtonFont string = "Button"
var ButtonFontSize uint32 = 24

type Button struct {
	Sprite2D

	PressCallback ButtonCallback
	EnterCallback ButtonCallback
	LeaveCallback ButtonCallback
	Entered       bool
	Text          string

	EnterModColor color.Color
	PressModColor color.Color
	Text2D        Text2D
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

	this.Text2D.Init(ButtonFont, ButtonFontSize, this.Text)
	this.Text2D.Transform.Origin = [2]float32{0.0, 0.0}
	this.Text2D.NotRelativeToCamera = 0
	this.Text2D.Color = Color{0, 0, 0, 255}
	this.Text2D.Transform.Origin = [2]float32{0.5, 0.5}

	RenderMgr.AddObject(&this.Text2D)
}

func (this *Button) Update(delta_time float32) {
	if focusedButton == nil || focusedButton == this {
		mpos := InputMgr.Mouse.ToScreenPosition()
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

	this.Text2D.Text = this.Text
	size := this.Transform.Size
	size[0], size[1] = size[0]*this.Transform.Scale[0], size[1]*this.Transform.Scale[1]
	offset := [2]float32{size[0] * this.Transform.Origin[0], size[1] * this.Transform.Origin[1]}
	this.Text2D.Transform.Position = this.Transform.Position.Sub(offset).Add(size.Mul(0.5))
	this.Text2D.NotRelativeToCamera = this.NotRelativeToCamera
}

func (this *Button) Terminate() {
	if focusedButton == this {
		focusedButton = nil
	}

	UpdateMgr.RemoveObject(this)
	RenderMgr.RemoveObject(this)
	RenderMgr.RemoveObject(&this.Text2D)
	this.Text2D.Terminate()
}

var focusedSlider *Slider

type SliderCallback func(sld *Slider)

type Slider struct {
	Long   Sprite2D
	Circle Sprite2D

	Entered       bool
	Value         float32
	StepSize      float32
	EnterModColor color.Color
	PressModColor color.Color

	EnterCallback        SliderCallback
	LeaveCallback        SliderCallback
	ValueChangedCallback SliderCallback

	clickedPos   [2]float32
	clickedValue float32
}

func (this *Slider) Init(pos [2]float32, longTex, circleTex string) {
	this.Long.Init(longTex)
	this.Circle.Init(circleTex)

	if this.Long.Texture == nil {
		rt := Render.CreateRenderTexture("SliderLongTexture", 200, 10, 1, false, false, false, false)
		rt.SetAsTarget()
		Render.ClearScreen(Color{255, 255, 255, 255})
		rt.UnsetAsTarget()
		this.Long.InitTexture(rt)
	}

	this.Long.Transform.Position = pos

	if this.Circle.Texture == nil {
		rt := Render.CreateRenderTexture("SliderCircleTexture", 20, 20, 1, false, false, false, false)
		rt.SetAsTarget()
		Render.ClearScreen(Color{200, 200, 200, 255})
		rt.UnsetAsTarget()
		this.Circle.InitTexture(rt)
	}

	this.Circle.Transform.Origin = [2]float32{0.5, 0.5}

	this.EnterModColor = Color{180, 180, 180, 255}
	this.PressModColor = Color{128, 128, 128, 255}

	this.Long.NotRelativeToCamera = 0
	this.Circle.NotRelativeToCamera = 0

	RenderMgr.AddObject(&this.Long)
	RenderMgr.AddObject(&this.Circle)
	UpdateMgr.AddObject(this)
}

func (this *Slider) setPositionByValue() {
	this.Circle.Transform.Position = this.Long.Transform.Position.Add([2]float32{this.Value * this.Long.Transform.Size[0] * this.Long.Transform.Scale[0], this.Long.Transform.Size[1] * this.Long.Transform.Scale[1] / 2.0})
}

func (this *Slider) Update(delta_time float32) {
	oldValue := this.Value
	this.setPositionByValue()

	if focusedSlider == nil || focusedSlider == this {
		mpos := InputMgr.Mouse.ToScreenPosition()
		size := this.Circle.Transform.Size.MulVec(this.Circle.Transform.Scale)
		pos := this.Circle.Transform.Position.Sub(this.Circle.Transform.Origin.MulVec(size))
		prevEntered := this.Entered
		this.Entered = mpos[0] > pos[0] && mpos[1] > pos[1] && mpos[0] <= pos[0]+size[0] && mpos[1] <= pos[1]+size[1]

		lsize := this.Long.Transform.Size.MulVec(this.Long.Transform.Scale)
		lpos := this.Long.Transform.Position.Sub(this.Long.Transform.Origin.MulVec(lsize))
		lentered := mpos[0] > lpos[0] && mpos[1] > lpos[1] && mpos[0] <= lpos[0]+lsize[0] && mpos[1] <= lpos[1]+lsize[1]

		if InputMgr.IsPressed(MouseButtonLeft) {
			if focusedSlider == nil && lentered {
				this.Value = mpos.Sub(this.Long.Transform.Position.Sub(this.Long.Transform.Origin.MulVec(this.Long.Transform.Size.MulVec(this.Long.Transform.Scale)))).X() / (this.Long.Transform.Size[0] * this.Long.Transform.Scale[0])
				focusedSlider = this
			}

			if focusedSlider == this {
				this.Circle.Texture.SetModColor(this.PressModColor)
				if InputMgr.JustPressed(MouseButtonLeft) {
					this.clickedPos = mpos
					this.clickedValue = this.Value
				} else {
					deltapos := mpos.Sub(this.clickedPos)
					deltax := deltapos.X()
					this.Value = this.clickedValue + (deltax / this.Long.Transform.Size[0] * this.Long.Transform.Scale[0])
				}
			}
		} else {
			if this.Entered {
				this.Circle.Texture.SetModColor(this.EnterModColor)
				if !prevEntered {
					focusedSlider = this
					if this.EnterCallback != nil {
						this.EnterCallback(this)
					}
				}
			} else {
				this.Circle.Texture.SetModColor(nil)
				if focusedSlider == this {
					focusedSlider = nil
				}
				if prevEntered {
					if this.LeaveCallback != nil {
						this.LeaveCallback(this)
					}
				}
			}
		}
	}

	if this.Value < 0.0 {
		this.Value = 0.0
	} else if this.Value > 1.0 {
		this.Value = 1.0
	}

	if this.StepSize != 0.0 && this.StepSize > 0.0 && this.StepSize <= 1.0 {
		val := this.Value / this.StepSize
		val = float32(math.Ceil(float64(val)))
		this.Value = this.StepSize * val
	}

	if this.Value != oldValue && this.ValueChangedCallback != nil {
		this.ValueChangedCallback(this)
	}
}

func (this *Slider) Terminate() {
	if focusedSlider == this {
		focusedSlider = nil
	}

	RenderMgr.RemoveObject(&this.Long)
	RenderMgr.RemoveObject(&this.Circle)
	UpdateMgr.RemoveObject(this)
}
