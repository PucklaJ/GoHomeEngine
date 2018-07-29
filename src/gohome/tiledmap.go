package gohome

import (
	"github.com/elliotmr/tmx"
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
)

type sprite2DConfiguration struct {
	TextureName string
	Flip        uint8
	Region      TextureRegion
}

type TiledMap struct {
	Sprite2D
	*tmx.Map

	layers []Texture
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
	rt.SetFiltering(FILTERING_NEAREST)
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

	this.generateTextures()
}

func getTextureRegionFromID(ts *tmx.TileSet, id uint32) TextureRegion {
	tilesWidth := ts.Columns

	x := id % tilesWidth
	y := (id - x) / tilesWidth

	var region TextureRegion
	region.Min[0] = float32(x*(ts.TileWidth+ts.Spacing) + ts.Margin)
	region.Min[1] = float32(y*(ts.TileHeight+ts.Spacing) + ts.Margin)
	region.Max[0] = region.Min[0] + float32(ts.TileWidth)
	region.Max[1] = region.Min[1] + float32(ts.TileHeight)

	return region
}

func getFlip(tile tmx.TileInstance) uint8 {
	if tile.FlippedHorizontally() {
		return FLIP_HORIZONTAL
	} else if tile.FlippedVertically() {
		return FLIP_VERTICAL
	} else if tile.FlippedDiagonally() {
		return FLIP_DIAGONALLY
	} else {
		return FLIP_NONE
	}
}

func (this *TiledMap) getSprite2DConfiguration(tile tmx.TileInstance) sprite2DConfiguration {
	gid := tile.GID()
	var ts *tmx.TileSet
	var config sprite2DConfiguration

	for i := 0; i < len(this.TileSets); i++ {
		t := this.TileSets[i]
		if gid >= t.FirstGID && gid <= t.FirstGID+(t.TileCount-1) {
			ts = t
			break
		}
	}

	if ts == nil {
		return config
	}

	id := gid - ts.FirstGID
	config.Region = getTextureRegionFromID(ts, id)
	config.TextureName = ts.Name
	config.Flip = getFlip(tile)

	return config
}

func (this *TiledMap) renderConfiguration(config sprite2DConfiguration, pos mgl32.Vec2, rt RenderTexture) {
	if config.TextureName == "" {
		return
	}

	texture := ResourceMgr.GetTexture(config.TextureName)
	if texture == nil {
		ErrorMgr.Error("TiledMap", config.TextureName, "Couldn't get texture from tile!")
		return
	}

	texture.SetFiltering(FILTERING_NEAREST)

	var spr Sprite2D
	spr.InitTexture(texture)

	spr.TextureRegion = config.Region
	spr.Transform.Size[0] = float32(config.Region.Width())
	spr.Transform.Size[1] = float32(config.Region.Height())
	spr.Transform.Position = pos
	spr.Flip = config.Flip

	RenderMgr.RenderRenderObject(&spr)
}

func (this *TiledMap) loadTileLayer(l *tmx.Layer) {
	data := l.Data

	iter, err := data.Iter()
	if err != nil {
		ErrorMgr.Error("TiledMap", l.Name, "Couldn't get the TileIterator!")
		return
	}

	texture := Render.CreateRenderTexture(l.Name+" RenderTexture", this.Width*this.TileWidth, this.Height*this.TileHeight, 1, false, false, false, false)
	if texture == nil {
		ErrorMgr.Error("TiledMap", l.Name, "Couldn't create the RenderTexture for the layer")
		return
	}
	texture.SetFiltering(FILTERING_NEAREST)

	var pos mgl32.Vec2
	texture.SetAsTarget()
	for iter.Next() {
		tile := iter.Get()
		if tile.GID() == 0 {
			continue
		}
		config := this.getSprite2DConfiguration(tile)
		counter := iter.GetIndex()
		pos[0] = float32((counter % this.Width) * this.TileWidth)
		pos[1] = float32(((counter - (counter % this.Width)) / this.Width) * this.TileHeight)

		this.renderConfiguration(config, pos, texture)
	}
	texture.UnsetAsTarget()

	if l.Opacity != nil {
		alpha := uint8(*l.Opacity * 255.0)
		texture.SetModColor(Color{255, 255, 255, alpha})
	}

	this.layers = append(this.layers, texture)
}

func (this *TiledMap) loadImageLayer(l *tmx.Layer) {

}

func (this *TiledMap) loadObjectGroup(l *tmx.Layer) {
	texture := Render.CreateRenderTexture(l.Name+" RenderTexture", this.Width*this.TileWidth, this.Height*this.TileHeight, 1, false, false, false, false)

	var pos mgl32.Vec2
	objs := l.Objects
	texture.SetAsTarget()
	for i := 0; i < len(objs); i++ {
		o := objs[i]
		if o.GID != nil {
			config := this.getSprite2DConfiguration(tmx.TileInstance(*o.GID))
			pos[0] = float32(o.X)
			pos[1] = float32(o.Y) - float32(this.TileHeight)

			this.renderConfiguration(config, pos, texture)
		}
	}
	texture.UnsetAsTarget()

	texture.SetFiltering(FILTERING_NEAREST)

	this.layers = append(this.layers, texture)
}

func (this *TiledMap) generateTextures() {
	for i := 0; i < len(this.TileSets); i++ {
		t := this.TileSets[i]
		img := t.Image
		if img == nil {
			continue
		}
		if t.TileCount == 0 {
			if img.Height != nil && img.Width != nil {
				rows := (uint32(*img.Height) - 2*t.Margin + t.Spacing) / (t.Spacing + t.TileHeight)
				if t.Columns == 0 {
					t.Columns = (uint32(*img.Width) - 2*t.Margin + t.Spacing) / (t.Spacing + t.TileWidth)
				}
				t.TileCount = rows * t.Columns
			}
		}

		ResourceMgr.PreloadTexture(t.Name, img.Source)
	}
	ResourceMgr.LoadPreloadedResources()

	for i := 0; i < len(this.TileSets); i++ {
		img := this.TileSets[i].Image
		if img == nil || img.Trans == nil {
			continue
		}
		tex := ResourceMgr.GetTexture(this.TileSets[i].Name)
		if tex == nil {
			continue
		}
		keyCol := *img.Trans
		if keyCol[0] != '#' {
			keyCol = "#" + keyCol
		}
		col := HEXToColor(keyCol)
		tex.SetKeyColor(col)
	}

	rt := this.Texture.(RenderTexture)
	prev := RenderMgr.Projection2D
	RenderMgr.SetProjection2DToTexture(rt)
	for i := 0; i < len(this.Layers); i++ {
		l := this.Layers[i]
		if l.Data != nil {
			this.loadTileLayer(l)
		} else if len(l.Objects) > 0 {
			this.loadObjectGroup(l)
		} else {
			this.loadImageLayer(l)
		}
	}
	rt.SetAsTarget()

	for i := 0; i < len(this.layers); i++ {
		l := this.layers[i]
		var spr Sprite2D
		spr.InitTexture(l)
		RenderMgr.RenderRenderObject(&spr)
	}

	rt.UnsetAsTarget()
	RenderMgr.Projection2D = prev
}
