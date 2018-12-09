package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/image/colornames"
)

type SpriteAnimationScene struct {
	mage1           gohome.Sprite2D
	mage1Anim       gohome.Tweenset
	mage2           gohome.Sprite2D
	mage2Anim       gohome.Tweenset
	mage3           gohome.Sprite2D
	mage3Anim       gohome.Tweenset
	andromalius     gohome.Sprite2D
	andromaliusAnim gohome.Tweenset
	shadow          gohome.Sprite2D
	shadowAnim      gohome.Tweenset
	acid1           gohome.Sprite2D
	acid1Anim       gohome.Tweenset
	acid2           gohome.Sprite2D
	acid2Anim       gohome.Tweenset
	acid3           gohome.Sprite2D
	acid3Anim       gohome.Tweenset
	natsu           gohome.Sprite2D
	natsuAnim       gohome.Tweenset
}

func (this *SpriteAnimationScene) Init() {
	gohome.Init2DShaders()
	gohome.ResourceMgr.PreloadTexture("Mage1", "mage-1-85x94.png")
	gohome.ResourceMgr.PreloadTexture("Mage2", "mage-2-122x110.png")
	gohome.ResourceMgr.PreloadTexture("Mage3", "mage-3-87x110.png")
	gohome.ResourceMgr.PreloadTexture("Andromalius", "andromalius-57x88.png")
	gohome.ResourceMgr.PreloadTexture("Shadow", "shadow-80x70.png")
	gohome.ResourceMgr.PreloadTexture("Acid", "acid2-14x67.png")
	gohome.ResourceMgr.PreloadTexture("Natsu", "Natsu.png")
	gohome.ResourceMgr.LoadPreloadedResources()

	gohome.Render.SetBackgroundColor(colornames.White)

	this.mage1.Init("Mage1")
	this.mage1Anim = gohome.SpriteAnimation2D(this.mage1.Texture, 4, 2, 0.1)
	this.mage1Anim.SetParent(&this.mage1)
	this.mage1Anim.Start()
	gohome.UpdateMgr.AddObject(&this.mage1Anim)
	gohome.RenderMgr.AddObject(&this.mage1)
	this.mage1Anim.Loop = true

	this.mage2.Init("Mage2")
	this.mage2Anim = gohome.SpriteAnimation2D(this.mage2.Texture, 4, 2, 0.1)
	this.mage2Anim.SetParent(&this.mage2)
	this.mage2Anim.Start()
	gohome.UpdateMgr.AddObject(&this.mage2Anim)
	gohome.RenderMgr.AddObject(&this.mage2)
	this.mage2.Transform.Position[0] = 85.0
	this.mage2Anim.Loop = true

	this.mage3.Init("Mage3")
	this.mage3Anim = gohome.SpriteAnimation2D(this.mage3.Texture, 4, 2, 0.1)
	this.mage3Anim.SetParent(&this.mage3)
	this.mage3Anim.Start()
	gohome.UpdateMgr.AddObject(&this.mage3Anim)
	gohome.RenderMgr.AddObject(&this.mage3)
	this.mage3.Transform.Position[0] = 85.0 + 122.0
	this.mage3Anim.Loop = true

	this.andromalius.Init("Andromalius")
	this.andromaliusAnim = gohome.SpriteAnimation2DOffset(this.andromalius.Texture, 4, 1, 0, 2*88, 4*57, 0, 0.1)
	this.andromaliusAnim.SetParent(&this.andromalius)
	this.andromaliusAnim.LoopBackwards = true
	this.andromaliusAnim.Start()
	gohome.UpdateMgr.AddObject(&this.andromaliusAnim)
	gohome.RenderMgr.AddObject(&this.andromalius)
	this.andromalius.Transform.Position[0] = 85.0 + 122.0 + 87.0
	this.andromaliusAnim.Loop = true

	this.shadow.Init("Shadow")
	shadowAnim1 := gohome.SpriteAnimation2DOffset(this.shadow.Texture, 4, 3, 0, 0, 0, 2*70, 0.1)
	shadowAnim2 := gohome.SpriteAnimation2DOffset(this.shadow.Texture, 2, 1, 0, 3*70, 2*80, 1*70, 0.1)
	shadowAnim3 := gohome.SpriteAnimation2DOffset(this.shadow.Texture, 4, 1, 0, 4*70, 0, 0, 0.1)
	this.shadowAnim = shadowAnim1.Merge(shadowAnim2).Merge(shadowAnim3)
	this.shadowAnim.SetParent(&this.shadow)
	this.shadowAnim.Start()
	gohome.UpdateMgr.AddObject(&this.shadowAnim)
	gohome.RenderMgr.AddObject(&this.shadow)
	this.shadow.Transform.Position[0] = 85.0 + 122.0 + 87.0 + 57.0
	this.shadowAnim.Loop = true
	this.shadowAnim.LoopBackwards = true
	this.andromalius.Texture.SetKeyColor(colornames.Black)

	this.acid1.Init("Acid")
	this.acid1Anim = gohome.SpriteAnimation2DOffset(this.acid1.Texture, 8, 1, 0, 0, 0, 2*67, 0.2)
	this.acid1Anim.SetParent(&this.acid1)
	this.acid1Anim.Start()
	gohome.UpdateMgr.AddObject(&this.acid1Anim)
	gohome.RenderMgr.AddObject(&this.acid1)
	this.acid1.Transform.Position[1] = 110.0
	this.acid1Anim.Loop = true
	this.acid1Anim.LoopBackwards = true
	this.acid1.Texture.SetKeyColor(colornames.Darkviolet)

	this.acid2.Init("Acid")
	this.acid2Anim = gohome.SpriteAnimation2DOffset(this.acid2.Texture, 8, 1, 0, 2*67, 0, 0, 0.2)
	this.acid2Anim.SetParent(&this.acid2)
	this.acid2Anim.Start()
	gohome.UpdateMgr.AddObject(&this.acid2Anim)
	gohome.RenderMgr.AddObject(&this.acid2)
	this.acid2.Transform.Position[0] = 14.0
	this.acid2.Transform.Position[1] = 110.0
	this.acid2Anim.Loop = true
	this.acid2Anim.LoopBackwards = true

	this.acid3.Init("Acid")
	this.acid3Anim = gohome.SpriteAnimation2DOffset(this.acid3.Texture, 8, 1, 0, 1*67, 0, 1*67, 0.2)
	this.acid3Anim = this.acid3Anim.Merge(this.acid2Anim).Merge(this.acid1Anim)
	this.acid3Anim.SetParent(&this.acid3)
	this.acid3Anim.Start()
	gohome.UpdateMgr.AddObject(&this.acid3Anim)
	gohome.RenderMgr.AddObject(&this.acid3)
	this.acid3.Transform.Position[0] = 14.0 * 2
	this.acid3.Transform.Position[1] = 110.0
	this.acid3Anim.Loop = true
	this.acid3Anim.LoopBackwards = false

	this.natsu.Init("Natsu")
	this.natsu.Texture.SetKeyColor(gohome.Color{0, 128, 0, 255})
	this.natsuAnim = gohome.SpriteAnimation2DRegions([]gohome.TextureRegion{
		{[2]float32{8, 204}, [2]float32{8 + 45, 204 + 82}},
		{[2]float32{53, 204}, [2]float32{53 + 39, 204 + 82}},
		{[2]float32{97, 205}, [2]float32{97 + 45, 205 + 81}},
		{[2]float32{149, 205}, [2]float32{149 + 47, 205 + 81}},
		{[2]float32{202, 205}, [2]float32{202 + 48, 205 + 81}},
		{[2]float32{262, 205}, [2]float32{262 + 52, 205 + 81}},
		{[2]float32{322, 205}, [2]float32{322 + 53, 205 + 81}},
		{[2]float32{388, 205}, [2]float32{388 + 45, 205 + 81}},
		{[2]float32{445, 205}, [2]float32{445 + 54, 205 + 81}},
	}, 0.2)
	this.natsuAnim.Loop = true
	this.natsuAnim.LoopBackwards = true
	this.natsuAnim.SetParent(&this.natsu)
	this.natsuAnim.Start()
	gohome.UpdateMgr.AddObject(&this.natsuAnim)
	gohome.RenderMgr.AddObject(&this.natsu)
	this.natsu.Transform.Position[1] = 110.0 + 67.0
}

func (this *SpriteAnimationScene) Update(delta_time float32) {

}

func (this *SpriteAnimationScene) Terminate() {

}
