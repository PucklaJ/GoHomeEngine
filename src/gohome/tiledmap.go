package gohome

import (
	"github.com/elliotmr/tmx"
)

type TiledMap struct {
	Sprite2D
	tmx.Map
}

func (this *TiledMap) Init(tmxmapname string) {

}
