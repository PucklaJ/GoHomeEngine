package gohome

import (
	"github.com/elliotmr/tmx"
	"image/color"
)

type TiledMap struct {
	Sprite2D
	*tmx.Map
}

func HEXToUint4(str byte) uint8 {
	switch str {
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	case '3':
		return 3
	case '4':
		return 4
	case '5':
		return 5
	case '6':
		return 6
	case '7':
		return 7
	case '8':
		return 8
	case '9':
		return 9
	case 'a', 'A':
		return 10
	case 'b', 'B':
		return 11
	case 'c', 'C':
		return 12
	case 'd', 'D':
		return 13
	case 'e', 'E':
		return 14
	case 'f', 'F':
		return 15
	}

	return 0
}

func HEXToUint8(str string) uint8 {
	var first, second uint8
	first = HEXToUint4(str[0])
	second = HEXToUint4(str[1])
	return (first << 4) | second
}

func HEXToColor(str string) color.Color {
	col := &Color{}
	col.R = HEXToUint8(str[1:3])
	col.G = HEXToUint8(str[3:5])
	col.B = HEXToUint8(str[5:7])
	col.A = 255
	return col
}

func (this *TiledMap) Init(tmxmapname string) {
	tmxmap := ResourceMgr.GetTMXMap(tmxmapname)

	if tmxmap == nil {
		ErrorMgr.Warning("TiledMap", tmxmapname, "Has not been loaded!")
		return
	}

	this.Map = tmxmap
	this.Texture = Render.CreateRenderTexture("TMXMapTexture", this.Width*this.TileWidth, this.Height*this.TileHeight, 1, false, false, false, false)
	rt := this.Texture.(RenderTexture)
	rt.SetAsTarget()
	var backCol color.Color
	back := this.BackgroundColor
	if back != nil {
		backCol = HEXToColor(*back)
	} else {
		backCol = Color{0, 0, 0, 255}
	}
	Render.ClearScreen(backCol)
	rt.UnsetAsTarget()

	this.Sprite2D.commonInit()
}
