package gohome

type textureData struct {
	data   []byte
	width  int
	height int
}

var preloadedTextureData map[Texture]textureData
