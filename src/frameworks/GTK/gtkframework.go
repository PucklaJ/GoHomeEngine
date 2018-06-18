package framework

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"
import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
	"io"
	"os"
	"time"
	"unsafe"
)

type GTKFramework struct {
	startOtherThanPaint time.Time
	endOtherThanPaint   time.Time
}

func (this *GTKFramework) Init(ml *gohome.MainLoop) error {
	argv := os.Args
	args := len(argv)
	var argvCString []*C.char

	for i := 0; i < args; i++ {
		argvCString = append(argvCString, C.CString(argv[i]))
		defer C.free(unsafe.Pointer(argvCString[i]))
	}

	C.initialise(C.int(args), &argvCString[0])

	ml.InitWindowAndRenderer()
	gohome.Render.AfterInit()
	ml.InitManagers()
	ml.SetupStartScene()

	C.loop()

	return nil
}
func (this *GTKFramework) Update() {

}
func (this *GTKFramework) Terminate() {

}
func (this *GTKFramework) PollEvents() {

}
func (this *GTKFramework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	titleCS := C.CString(title)
	defer C.free(unsafe.Pointer(titleCS))
	if int(C.createWindow(C.uint(windowWidth), C.uint(windowHeight), titleCS)) == 0 {
		return nil
	}
	return nil
}
func (this *GTKFramework) WindowClosed() bool {
	return false
}
func (this *GTKFramework) WindowSwap() {

}
func (this *GTKFramework) WindowGetSize() mgl32.Vec2 {
	var v [2]C.float
	var v1 mgl32.Vec2

	C.windowGetSize(&v[0], &v[1])

	v1[0] = float32(v[0])
	v1[1] = float32(v[1])

	return v1
}
func (this *GTKFramework) WindowSetFullscreen(b bool) {

}
func (this *GTKFramework) WindowIsFullscreen() bool {
	return false
}

func (this *GTKFramework) CurserShow() {

}
func (this *GTKFramework) CursorHide() {

}
func (this *GTKFramework) CursorDisable() {

}
func (this *GTKFramework) CursorShown() bool {
	return true
}
func (this *GTKFramework) CursorHidden() bool {
	return false
}
func (this *GTKFramework) CursorDisabled() bool {
	return false
}

func (this *GTKFramework) OpenFile(file string) (io.ReadCloser, error) {
	return os.Open(file)
}
func (this *GTKFramework) LoadLevel(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	return nil
}

func (this *GTKFramework) ShowYesNoDialog(title, message string) uint8 {
	return gohome.DIALOG_CANCELLED
}

func (this *GTKFramework) OnResize(callback func(newWidth, newHeight uint32)) {

}
func (this *GTKFramework) OnMove(callback func(newPosX, newPosY uint32)) {

}
func (this *GTKFramework) OnClose(callback func()) {

}
func (this *GTKFramework) OnFocus(callback func(focused bool)) {

}

func (this *GTKFramework) StartTextInput() {

}
func (this *GTKFramework) GetTextInput() string {
	return ""
}
func (this *GTKFramework) EndTextInput() {

}
