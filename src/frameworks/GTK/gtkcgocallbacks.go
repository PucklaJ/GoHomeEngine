package framework

/*
	#include "includes.h"
	#include <gdk/gdk.h>
*/
import "C"
import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"log"
	"time"
)

//export gtkgo_quit
func gtkgo_quit() {
	C.gtk_main_quit()
}

var firstFrameDone bool = false

//export gtkgo_gl_area_render
func gtkgo_gl_area_render(area *C.GtkGLArea, context *C.GdkGLContext) {
	gtkframework := gohome.Framew.(*GTKFramework)
	gtkframework.endOtherThanPaint = time.Now()
	if firstFrameDone {
		addTime := float32(gtkframework.endOtherThanPaint.Sub(gtkframework.startOtherThanPaint).Seconds())
		gohome.FPSLimit.AddTime(addTime)
	}
	gohome.MainLop.LoopOnce()
	gtkframework.startOtherThanPaint = time.Now()
	firstFrameDone = true
}

//export gtkgo_gl_area_realize
func gtkgo_gl_area_realize(area *C.GtkGLArea, err int) {
	if err == 1 {
		log.Println("Error:", C.GoString(C.ErrorString))
		return
	}
}

//export gtkgo_gl_area_key_press
func gtkgo_gl_area_key_press(widget *C.GtkWidget, event *C.GdkEventKey) {
	gohomeKey := gdkkeysymTogohomekey(event.keyval)
	gohome.InputMgr.PressKey(gohomeKey)
}

//export gtkgo_gl_area_key_release
func gtkgo_gl_area_key_release(widget *C.GtkWidget, event *C.GdkEventKey) {
	gohomeKey := gdkkeysymTogohomekey(event.keyval)
	gohome.InputMgr.ReleaseKey(gohomeKey)
}

func gdkkeysymTogohomekey(key C.guint) gohome.Key {
	switch key {
	case C.GDK_KEY_A, C.GDK_KEY_a:
		return gohome.KeyA
		break
	case C.GDK_KEY_B, C.GDK_KEY_b:
		return gohome.KeyB
		break
	case C.GDK_KEY_C, C.GDK_KEY_c:
		return gohome.KeyC
		break
	case C.GDK_KEY_D, C.GDK_KEY_d:
		return gohome.KeyD
		break
	case C.GDK_KEY_E, C.GDK_KEY_e:
		return gohome.KeyE
		break
	case C.GDK_KEY_F, C.GDK_KEY_f:
		return gohome.KeyF
		break
	case C.GDK_KEY_G, C.GDK_KEY_g:
		return gohome.KeyG
		break
	case C.GDK_KEY_H, C.GDK_KEY_h:
		return gohome.KeyH
		break
	case C.GDK_KEY_I, C.GDK_KEY_i:
		return gohome.KeyI
		break
	case C.GDK_KEY_J, C.GDK_KEY_j:
		return gohome.KeyJ
		break
	case C.GDK_KEY_K, C.GDK_KEY_k:
		return gohome.KeyK
		break
	case C.GDK_KEY_L, C.GDK_KEY_l:
		return gohome.KeyL
		break
	case C.GDK_KEY_M, C.GDK_KEY_m:
		return gohome.KeyM
		break
	case C.GDK_KEY_N, C.GDK_KEY_n:
		return gohome.KeyN
		break
	case C.GDK_KEY_O, C.GDK_KEY_o:
		return gohome.KeyO
		break
	case C.GDK_KEY_P, C.GDK_KEY_p:
		return gohome.KeyP
		break
	case C.GDK_KEY_Q, C.GDK_KEY_q:
		return gohome.KeyQ
		break
	case C.GDK_KEY_R, C.GDK_KEY_r:
		return gohome.KeyR
		break
	case C.GDK_KEY_S, C.GDK_KEY_s:
		return gohome.KeyS
		break
	case C.GDK_KEY_T, C.GDK_KEY_t:
		return gohome.KeyT
		break
	case C.GDK_KEY_U, C.GDK_KEY_u:
		return gohome.KeyU
		break
	case C.GDK_KEY_V, C.GDK_KEY_v:
		return gohome.KeyV
		break
	case C.GDK_KEY_W, C.GDK_KEY_w:
		return gohome.KeyW
		break
	case C.GDK_KEY_X, C.GDK_KEY_x:
		return gohome.KeyX
		break
	case C.GDK_KEY_Y, C.GDK_KEY_y:
		return gohome.KeyY
		break
	case C.GDK_KEY_Z, C.GDK_KEY_z:
		return gohome.KeyZ
		break
	case C.GDK_KEY_0, C.GDK_KEY_equal, C.GDK_KEY_braceright:
		return gohome.Key0
	case C.GDK_KEY_1, C.GDK_KEY_exclam, C.GDK_KEY_onesuperior:
		return gohome.Key1
	case C.GDK_KEY_2, C.GDK_KEY_quotedbl, C.GDK_KEY_twosuperior:
		return gohome.Key2
	case C.GDK_KEY_3, C.GDK_KEY_section, C.GDK_KEY_threesuperior:
		return gohome.Key3
	case C.GDK_KEY_4, C.GDK_KEY_dollar, C.GDK_KEY_onequarter:
		return gohome.Key4
	case C.GDK_KEY_5, C.GDK_KEY_percent, C.GDK_KEY_onehalf:
		return gohome.Key5
	case C.GDK_KEY_6, C.GDK_KEY_ampersand, C.GDK_KEY_notsign:
		return gohome.Key6
	case C.GDK_KEY_7, C.GDK_KEY_slash, C.GDK_KEY_braceleft:
		return gohome.Key7
	case C.GDK_KEY_8, C.GDK_KEY_parenleft, C.GDK_KEY_bracketleft:
		return gohome.Key8
	case C.GDK_KEY_9, C.GDK_KEY_parenright, C.GDK_KEY_bracketright:
		return gohome.Key9
	case C.GDK_KEY_Shift_R:
		return gohome.KeyRightShift
	case C.GDK_KEY_Shift_L:
		return gohome.KeyLeftShift
	case C.GDK_KEY_Alt_L:
		return gohome.KeyLeftAlt
	case C.GDK_KEY_Alt_R:
		return gohome.KeyRightAlt
	case C.GDK_KEY_ISO_Level3_Shift:
		return gohome.KeyRightSuper
	case C.GDK_KEY_Control_L:
		return gohome.KeyLeftControl
	case C.GDK_KEY_Control_R:
		return gohome.KeyRightControl
	case C.GDK_KEY_Caps_Lock:
		return gohome.KeyCapsLock
	case C.GDK_KEY_Tab:
		return gohome.KeyTab
	case C.GDK_KEY_F1:
		return gohome.KeyF1
	case C.GDK_KEY_F2:
		return gohome.KeyF2
	case C.GDK_KEY_F3:
		return gohome.KeyF3
	case C.GDK_KEY_F4:
		return gohome.KeyF4
	case C.GDK_KEY_F5:
		return gohome.KeyF5
	case C.GDK_KEY_F6:
		return gohome.KeyF6
	case C.GDK_KEY_F7:
		return gohome.KeyF7
	case C.GDK_KEY_F8:
		return gohome.KeyF8
	case C.GDK_KEY_F9:
		return gohome.KeyF9
	case C.GDK_KEY_F10:
		return gohome.KeyF10
	case C.GDK_KEY_F11:
		return gohome.KeyF11
	case C.GDK_KEY_F13:
		return gohome.KeyF13
	case C.GDK_KEY_F14:
		return gohome.KeyF14
	case C.GDK_KEY_F15:
		return gohome.KeyF15
	case C.GDK_KEY_F16:
		return gohome.KeyF16
	case C.GDK_KEY_F17:
		return gohome.KeyF17
	case C.GDK_KEY_F18:
		return gohome.KeyF18
	case C.GDK_KEY_F19:
		return gohome.KeyF19
	case C.GDK_KEY_F20:
		return gohome.KeyF20
	case C.GDK_KEY_F21:
		return gohome.KeyF21
	case C.GDK_KEY_F22:
		return gohome.KeyF22
	case C.GDK_KEY_F23:
		return gohome.KeyF23
	case C.GDK_KEY_F24:
		return gohome.KeyF24
	case C.GDK_KEY_F25:
		return gohome.KeyF25
	case C.GDK_KEY_Escape:
		return gohome.KeyEscape
	default:
		return gohome.KeyUnknown
		break
	}

	return gohome.KeyUnknown
}

//export gtkgo_gl_area_button_press
func gtkgo_gl_area_button_press(widget *C.GtkWidget, event *C.GdkEventButton) {
	gohomeKey := gdkbuttonTogohomekey(event.button)
	gohome.InputMgr.PressKey(gohomeKey)

	if gohomeKey == gohome.MouseButtonLeft {
		gohome.InputMgr.Touch(0)
	}
}

//export gtkgo_gl_area_button_release
func gtkgo_gl_area_button_release(widget *C.GtkWidget, event *C.GdkEventButton) {
	gohomeKey := gdkbuttonTogohomekey(event.button)
	gohome.InputMgr.ReleaseKey(gohomeKey)

	if gohomeKey == gohome.MouseButtonLeft {
		gohome.InputMgr.ReleaseTouch(0)
	}
}

func gdkbuttonTogohomekey(button C.guint) gohome.Key {
	switch button {
	case 1:
		return gohome.MouseButtonLeft
	case 2:
		return gohome.MouseButtonMiddle
	case 3:
		return gohome.MouseButtonRight
	default:
		return gohome.MouseButton1 - 1 + gohome.Key(button)
	}

	return gohome.KeyUnknown
}

//export gtkgo_gl_area_motion_notify
func gtkgo_gl_area_motion_notify(widget *C.GtkWidget, event *C.GdkEventMotion) {
	framew := gohome.Framew.(*GTKFramework)
	gohome.InputMgr.Mouse.Pos[0] = int16(event.x)
	gohome.InputMgr.Mouse.Pos[1] = int16(event.y)
	gohome.InputMgr.Mouse.DPos[0] = int16(event.x) - framew.prevMousePos[0]
	gohome.InputMgr.Mouse.DPos[1] = int16(event.y) - framew.prevMousePos[1]
	framew.prevMousePos[0] = gohome.InputMgr.Mouse.Pos[0]
	framew.prevMousePos[1] = gohome.InputMgr.Mouse.Pos[1]

	inputTouch := gohome.InputMgr.Touches[0]
	inputTouch.Pos = gohome.InputMgr.Mouse.Pos
	inputTouch.DPos = gohome.InputMgr.Mouse.DPos
	inputTouch.PPos = framew.prevMousePos
	inputTouch.ID = 0
	gohome.InputMgr.Touches[0] = inputTouch
}

//export gtkgo_gl_area_scroll
func gtkgo_gl_area_scroll(widget *C.GtkWidget, event *C.GdkEventScroll, gdkevent *C.GdkEvent) {
	/*
		one of GDK_SCROLL_UP, GDK_SCROLL_DOWN, GDK_SCROLL_LEFT, GDK_SCROLL_RIGHT or GDK_SCROLL_SMOOTH).
	*/
	var wheel [2]int8
	switch event.direction {
	case C.GDK_SCROLL_UP:
		wheel[0] = 0
		wheel[1] = 1
	case C.GDK_SCROLL_DOWN:
		wheel[0] = 0
		wheel[1] = -1
	case C.GDK_SCROLL_LEFT:
		wheel[0] = -1
		wheel[1] = 0
	case C.GDK_SCROLL_RIGHT:
		wheel[0] = 1
		wheel[1] = 0
	case C.GDK_SCROLL_SMOOTH:
		var sxc, syc C.gdouble
		C.gdk_event_get_scroll_deltas(gdkevent, &sxc, &syc)
		wheel[0] += int8(sxc)
		wheel[1] += int8(syc)
	default:
		break
	}

	gohome.InputMgr.Mouse.Wheel = wheel
}
