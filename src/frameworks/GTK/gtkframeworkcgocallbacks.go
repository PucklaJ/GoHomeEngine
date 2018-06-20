package framework

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"time"
)

var firstFrameDone bool = false

func gtkgo_gl_area_render() {
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

func gtkgo_gl_area_motion_notify(x, y int16) {
	framew := gohome.Framew.(*GTKFramework)
	gohome.InputMgr.Mouse.Pos[0] = x
	gohome.InputMgr.Mouse.Pos[1] = y
	gohome.InputMgr.Mouse.DPos[0] = x - framew.prevMousePos[0]
	gohome.InputMgr.Mouse.DPos[1] = y - framew.prevMousePos[1]
	framew.prevMousePos[0] = gohome.InputMgr.Mouse.Pos[0]
	framew.prevMousePos[1] = gohome.InputMgr.Mouse.Pos[1]

	inputTouch := gohome.InputMgr.Touches[0]
	inputTouch.Pos = gohome.InputMgr.Mouse.Pos
	inputTouch.DPos = gohome.InputMgr.Mouse.DPos
	inputTouch.PPos = framew.prevMousePos
	inputTouch.ID = 0
	gohome.InputMgr.Touches[0] = inputTouch
}

func useWholeWindowAsGLArea() bool {
	framew := gohome.Framew.(*GTKFramework)
	return framew.UseWholeWindowAsGLArea
}
