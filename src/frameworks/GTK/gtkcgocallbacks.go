package framework

/*
	#include "includes.h"
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
