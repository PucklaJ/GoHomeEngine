package framework

import (
	// "fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/gl"
)

type AndroidFramework struct {
	appl   app.App
	glcotx gl.Context
}

func (this *AndroidFramework) Init(ml *gohome.MainLoop) error {
	app.Main(androidFrameworkmain)
	return nil
}
func androidFrameworkmain(a app.App) {
	var androidFramework *AndroidFramework
	androidFramework = gohome.Framew.(*AndroidFramework)
	androidFramework.appl = a

	for e := range a.Events() {
		switch e := a.Filter(e).(type) {
		case lifecycle.Event:
			switch e.Crosses(lifecycle.StageVisible) {
			case lifecycle.CrossOn:
				androidFramework.glcotx, _ = e.DrawContext.(gl.Context)
				a.Send(paint.Event{})
				break
			}
			break
		case paint.Event:
			androidFramework.glcotx.ClearColor(1.0, 0.0, 0.0, 1.0)
			androidFramework.glcotx.Clear(gl.COLOR_BUFFER_BIT)
			a.Publish()
			break
		}
	}
}
func (this *AndroidFramework) Update() {

}
func (this *AndroidFramework) Terminate() {

}
func (this *AndroidFramework) PollEvents() {

}
func (this *AndroidFramework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	return nil
}
func (this *AndroidFramework) WindowClosed() bool {
	return false
}
func (this *AndroidFramework) WindowSwap() {

}
func (this *AndroidFramework) WindowGetSize() mgl32.Vec2 {
	return mgl32.Vec2{0.0, 0.0}
}
func (this *AndroidFramework) WindowSetFullscreen(b bool) {

}
func (this *AndroidFramework) WindowIsFullscreen() bool {
	return false
}
func (this *AndroidFramework) CurserShow() {

}
func (this *AndroidFramework) CursorHide() {

}
func (this *AndroidFramework) CursorDisable() {

}
func (this *AndroidFramework) CursorShown() bool {
	return true
}
func (this *AndroidFramework) CursorHidden() bool {
	return false
}
func (this *AndroidFramework) CursorDisabled() bool {
	return false
}
func (this *AndroidFramework) LoadLevel(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	return nil
}
