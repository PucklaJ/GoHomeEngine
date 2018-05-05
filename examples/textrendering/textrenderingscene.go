package main

import (
	"github.com/PucklaMotzer09/freetypeparser"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"log"
)

type TextRenderingScene struct {
}

func (this *TextRenderingScene) Init() {
	gohome.Init2DShaders()

	font, err := ftparser.ParseFile("/usr/share/fonts/truetype/abyssinica/AbyssinicaSIL-R.ttf", 48)
	if err != nil {
		log.Println("Error parsing font file:", err)
		return
	}

	var textureAtlasSpr gohome.Sprite2D
	var textureAtlas gohome.Texture
	var textureAtlasTobj gohome.TransformableObject2D
	textureAtlas = gohome.Render.CreateTexture("TextureAtlas", false)
	textureAtlas.LoadFromImage(font.TextureAtlas)
	textureAtlasSpr.InitTexture(textureAtlas, &textureAtlasTobj)

	log.Println("TextureAtlas: W:", textureAtlas.GetWidth(), "H:", textureAtlas.GetHeight())

	gohome.RenderMgr.AddObject(&textureAtlasSpr, &textureAtlasTobj)
}

func (this *TextRenderingScene) Update(delta_time float32) {

}

func (this *TextRenderingScene) Terminate() {

}
