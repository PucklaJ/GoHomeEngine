package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/image/colornames"
)

type SpriteAnimationScene struct {
	mage1 gohome.Sprite2D
	mage1Anim gohome.Tweenset
	mage2 gohome.Sprite2D
	mage2Anim gohome.Tweenset
	mage3 gohome.Sprite2D
	mage3Anim gohome.Tweenset
	andromalius gohome.Sprite2D
	andromaliusAnim gohome.Tweenset
	shadow gohome.Sprite2D
	shadowAnim gohome.Tweenset
}

func (this *SpriteAnimationScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.PreloadTexture("Mage1","mage-1-85x94.png")
	gohome.ResourceMgr.PreloadTexture("Mage2","mage-2-122x110.png")
	gohome.ResourceMgr.PreloadTexture("Mage3","mage-3-87x110.png")
	gohome.ResourceMgr.PreloadTexture("Andromalius","andromalius-57x88.png")
	gohome.ResourceMgr.PreloadTexture("Shadow","shadow-80x70.png")
	gohome.ResourceMgr.LoadPreloadedResources()

	gohome.Render.SetBackgroundColor(colornames.White)

	this.mage1.Init("Mage1")
	this.mage1Anim = gohome.SpriteAnimation2D(this.mage1.Texture.GetWidth(),this.mage1.Texture.GetHeight(),4,2,0.1,true)
	this.mage1Anim.SetParent(&this.mage1)
	this.mage1Anim.Start()
	gohome.UpdateMgr.AddObject(&this.mage1Anim)
	gohome.RenderMgr.AddObject(&this.mage1)

	this.mage2.Init("Mage2")
	this.mage2Anim = gohome.SpriteAnimation2D(this.mage2.Texture.GetWidth(),this.mage2.Texture.GetHeight(),4,2,0.1,true)
	this.mage2Anim.SetParent(&this.mage2)
	this.mage2Anim.Start()
	gohome.UpdateMgr.AddObject(&this.mage2Anim)
	gohome.RenderMgr.AddObject(&this.mage2)
	this.mage2.Transform.Position[0] = 85.0

	this.mage3.Init("Mage3")
	this.mage3Anim = gohome.SpriteAnimation2D(this.mage3.Texture.GetWidth(),this.mage3.Texture.GetHeight(),4,2,0.1,true)
	this.mage3Anim.SetParent(&this.mage3)
	this.mage3Anim.Start()
	gohome.UpdateMgr.AddObject(&this.mage3Anim)
	gohome.RenderMgr.AddObject(&this.mage3)
	this.mage3.Transform.Position[0] = 85.0+122.0

	this.andromalius.Init("Andromalius")
	this.andromaliusAnim = gohome.SpriteAnimation2DOffset(this.andromalius.Texture.GetWidth(),this.andromalius.Texture.GetHeight(),4,1,0,2*88,4*57,0,0.1,true)
	this.andromaliusAnim.SetParent(&this.andromalius)
	this.andromaliusAnim.LoopBackwards = true
	this.andromaliusAnim.Start()
	gohome.UpdateMgr.AddObject(&this.andromaliusAnim)
	gohome.RenderMgr.AddObject(&this.andromalius)
	this.andromalius.Transform.Position[0] = 85.0+122.0+87.0

	this.shadow.Init("Shadow")
	this.shadowAnim = gohome.SpriteAnimation2DRegions([]gohome.TextureRegion{
		{[2]float32{0,0},[2]float32{80,70}},
		{[2]float32{80,0},[2]float32{160,70}},
		{[2]float32{160,0},[2]float32{240,70}},
		{[2]float32{240,0},[2]float32{320,70}},

		{[2]float32{0,70},[2]float32{80,140}},
		{[2]float32{80,70},[2]float32{160,140}},
		{[2]float32{160,70},[2]float32{240,140}},
		{[2]float32{240,70},[2]float32{320,140}},

		{[2]float32{0,140},[2]float32{80,210}},
		{[2]float32{80,140},[2]float32{160,210}},
		{[2]float32{160,140},[2]float32{240,210}},
		{[2]float32{240,140},[2]float32{320,210}},

		{[2]float32{0,210},[2]float32{80,280}},
		{[2]float32{80,210},[2]float32{160,280}},

		{[2]float32{0,280},[2]float32{80,350}},
		{[2]float32{80,280},[2]float32{160,350}},
		{[2]float32{160,280},[2]float32{240,350}},
		{[2]float32{240,280},[2]float32{320,350}},
	},0.1,true)
	this.shadowAnim.SetParent(&this.shadow)
	this.shadowAnim.LoopBackwards = true
	this.shadowAnim.Start()
	gohome.UpdateMgr.AddObject(&this.shadowAnim)
	gohome.RenderMgr.AddObject(&this.shadow)
	this.shadow.Transform.Position[0] = 85.0+122.0+87.0+57.0

	gohome.RenderMgr.RenderToScreenFirst = true
}

func (this *SpriteAnimationScene) Update(delta_time float32) {

}

func (this *SpriteAnimationScene) Terminate() {

}
