package framework

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"
import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/loaders/assimp"
	"github.com/go-gl/mathgl/mgl32"
	"io"
	"os"
	"strings"
	"time"
	"unsafe"
)

type GTKFramework struct {
	startOtherThanPaint time.Time
	endOtherThanPaint   time.Time

	prevMousePos [2]int16

	isFullscreen bool
}

func (this *GTKFramework) Init(ml *gohome.MainLoop) error {
	argv := os.Args
	args := len(argv)
	var argvCString []*C.char

	for i := 0; i < args; i++ {
		argvCString = append(argvCString, C.CString(argv[i]))
		defer C.free(unsafe.Pointer(argvCString[i]))
	}
	this.isFullscreen = false
	C.initialise(C.int(args), &argvCString[0])

	ml.InitWindowAndRenderer()
	gohome.Render.AfterInit()
	ml.InitManagers()
	gohome.RenderMgr.EnableBackBuffer = false
	ml.SetupStartScene()

	C.loop()

	return nil
}
func (this *GTKFramework) Update() {
	gohome.InputMgr.Mouse.Wheel[0] = 0
	gohome.InputMgr.Mouse.Wheel[1] = 0
	gohome.InputMgr.Mouse.DPos[0] = 0
	gohome.InputMgr.Mouse.DPos[1] = 0
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
	if b {
		C.gtk_window_fullscreen(C.Window)
	} else {
		C.gtk_window_unfullscreen(C.Window)
	}
	this.isFullscreen = b
}
func (this *GTKFramework) WindowIsFullscreen() bool {
	return this.isFullscreen
}

func (this *GTKFramework) CurserShow() {
	C.windowShowCursor()
}
func (this *GTKFramework) CursorHide() {
	C.windowHideCursor()
}
func (this *GTKFramework) CursorDisable() {
	C.windowDisableCursor()
}
func (this *GTKFramework) CursorShown() bool {
	return int(C.windowCursorShown()) == 1
}
func (this *GTKFramework) CursorHidden() bool {
	return int(C.windowCursorHidden()) == 1
}
func (this *GTKFramework) CursorDisabled() bool {
	return int(C.windowCursorDisabled()) == 1
}

func (this *GTKFramework) OpenFile(file string) (io.ReadCloser, error) {
	return os.Open(file)
}

func getFileExtension(file string) string {
	index := strings.LastIndex(file, ".")
	if index == -1 {
		return ""
	}
	return file[index+1:]
}

func equalIgnoreCase(str1, str string) bool {
	if len(str1) != len(str) {
		return false
	}
	for i := 0; i < len(str1); i++ {
		if str1[i] != str[i] {
			if str1[i] >= 65 && str1[i] <= 90 {
				if str[i] >= 97 && str[i] <= 122 {
					if str1[i]+32 != str[i] {
						return false
					}
				} else {
					return false
				}
			} else if str1[i] >= 97 && str1[i] <= 122 {
				if str[i] >= 65 && str[i] <= 90 {
					if str1[i]-32 != str[i] {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

func (this *GTKFramework) LoadLevel(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	extension := getFileExtension(path)
	if equalIgnoreCase(extension, "obj") {
		return loadLevelOBJ(rsmgr, name, path, preloaded, loadToGPU)
	}
	return loader.LoadLevelAssimp(rsmgr, name, path, preloaded, loadToGPU)
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
