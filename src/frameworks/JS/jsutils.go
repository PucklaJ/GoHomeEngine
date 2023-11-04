package framework

import (
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
)

func jskeyCodeTogohomeKey(keyCode int) gohome.Key {
	switch keyCode {
	case 8:
		return gohome.KeyBackspace
	case 9:
		return gohome.KeyTab
	case 13:
		return gohome.KeyEnter
	case 16:
		return gohome.KeyLeftShift
	case 17:
		return gohome.KeyLeftControl
	case 18:
		return gohome.KeyLeftAlt
	case 19:
		return gohome.KeyPause
	case 20:
		return gohome.KeyCapsLock
	case 27:
		return gohome.KeyEscape
	case 33:
		return gohome.KeyPageUp
	case 34:
		return gohome.KeyPageDown
	case 35:
		return gohome.KeyEnd
	case 36:
		return gohome.KeyHome
	case 37:
		return gohome.KeyLeft
	case 38:
		return gohome.KeyUp
	case 39:
		return gohome.KeyRight
	case 40:
		return gohome.KeyDown
	case 45:
		return gohome.KeyInsert
	case 46:
		return gohome.KeyDelete
	case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
		return gohome.Key0 + gohome.Key(keyCode-48)
	case 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90:
		return gohome.KeyA + gohome.Key(keyCode-65)
	case 91:
		return gohome.KeyLeftSuper
	case 92:
		return gohome.KeyRightSuper
	case 96, 97, 98, 99, 100, 101, 102, 103, 104, 105:
		return gohome.KeyKP0 + gohome.Key(keyCode-96)
	case 106:
		return gohome.KeyKPMultiply
	case 107:
		return gohome.KeyKPAdd
	case 109:
		return gohome.KeyKPSubtract
	case 110:
		return gohome.KeyKPDecimal
	case 111:
		return gohome.KeyKPDivide
	case 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123:
		return gohome.KeyF1 + gohome.Key(keyCode-112)
	case 144:
		return gohome.KeyNumLock
	case 145:
		return gohome.KeyScrollLock
	case 186:
		return gohome.KeySemicolon
	case 187:
		return gohome.KeyEqual
	case 188:
		return gohome.KeyComma
	case 189:
		return gohome.KeyMinus
	case 190:
		return gohome.KeyPeriod
	case 191:
		return gohome.KeySlash
	case 192:
		return gohome.KeyGraveAccent
	case 219:
		return gohome.KeyLeftBracket
	case 220:
		return gohome.KeyBackslash
	case 221:
		return gohome.KeyRightBracket
	case 222:
		return gohome.KeyApostrophe
	}

	return gohome.KeyUnknown
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
