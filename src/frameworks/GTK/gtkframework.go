package framework

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK/gtk"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
	"io"
	"os"
	"strings"
	"time"
)

type GTKFramework struct {
	startOtherThanPaint time.Time
	endOtherThanPaint   time.Time

	prevMousePos [2]int16

	isFullscreen bool

	UseWholeWindowAsGLArea bool
}

func (this *GTKFramework) InitStuff(ml *gohome.MainLoop) {
	ml.InitWindow()
	gtk.GetWindow().ToContainer().Add(gtk.GetGLArea().ToWidget())
	ml.InitRenderer()
	ml.InitManagers()
	gohome.Render.AfterInit()
	gohome.RenderMgr.EnableBackBuffer = false
	gohome.RenderMgr.RenderToScreenFirst = true
}

func (this *GTKFramework) Init(ml *gohome.MainLoop) error {

	this.isFullscreen = false
	gtk.OnRender = gtkgo_gl_area_render
	gtk.OnMotion = gtkgo_gl_area_motion_notify
	gtk.OnUseWholeScreen = useWholeWindowAsGLArea
	gtk.Init()
	if this.UseWholeWindowAsGLArea {
		this.InitStuff(ml)
	}
	gohome.RenderMgr.EnableBackBuffer = false
	ml.SetupStartScene()

	gtk.Main()

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
	return gtk.CreateWindow(windowWidth, windowHeight, title)
}
func (this *GTKFramework) WindowClosed() bool {
	return false
}
func (this *GTKFramework) WindowSwap() {

}

func (this *GTKFramework) WindowSetSize(size mgl32.Vec2) {
	gtk.WindowSetSize(size)
}

func (this *GTKFramework) WindowGetSize() mgl32.Vec2 {
	return gtk.WindowGetSize()
}
func (this *GTKFramework) WindowSetFullscreen(b bool) {
	gtk.WindowSetFullscreen(b)
	this.isFullscreen = b
}
func (this *GTKFramework) WindowIsFullscreen() bool {
	return this.isFullscreen
}

func (this *GTKFramework) CurserShow() {
	gtk.CursorShow()
}
func (this *GTKFramework) CursorHide() {
	gtk.CursorHide()
}
func (this *GTKFramework) CursorDisable() {
	gtk.CursorDisable()
}
func (this *GTKFramework) CursorShown() bool {
	return gtk.CursorShown()
}
func (this *GTKFramework) CursorHidden() bool {
	return gtk.CursorHidden()
}
func (this *GTKFramework) CursorDisabled() bool {
	return gtk.CursorDisabled()
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
	// extension := getFileExtension(path)
	// if equalIgnoreCase(extension, "obj") {
	// 	return loadLevelOBJ(rsmgr, name, path, preloaded, loadToGPU)
	// }
	// return loader.LoadLevelAssimp(rsmgr, name, path, preloaded, loadToGPU)
	return loadLevelOBJ(rsmgr, name, path, preloaded, loadToGPU)
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
func (this *GTKFramework) GetGtkWindow() gtk.Window {
	return gtk.GetWindow()
}
