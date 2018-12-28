package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

func jskeyCodeTogohomeKey(keyCode int) gohome.Key {
	switch keyCode {
	default:
		return gohome.KeyUnknown
	}
}

func jsmouseButtonTogohomeKey(button int) gohome.Key {
	switch button {
	case 0:
		return gohome.MouseButtonLeft
	case 1:
		return gohome.MouseButtonMiddle
	case 2:
		return gohome.MouseButtonRight
	default:
		return gohome.KeyUnknown
	}
}
