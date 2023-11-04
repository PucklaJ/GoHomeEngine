package framework

import (
	"fmt"
	"os"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/PucklaJ/mathgl/mgl32"

	samure "github.com/PucklaJ/samurai-render-go"
	samureGL "github.com/PucklaJ/samurai-render-go/backends/opengl"
)

type SAMUREFramework struct {
	Ctx     samure.Context
	gl      samureGL.Backend
	running bool

	CurrentOutputGeo    samure.Rect
	CurrentLayerSurface samure.LayerSurface
}

type SAMUREApp struct {
	f *SAMUREFramework
}

func (f *SAMUREFramework) Init(ml *gohome.MainLoop) error {
	cfg := samure.CreateContextConfig(&SAMUREApp{
		f: f,
	})

	cfg.KeyboardInteraction = true
	cfg.GL.MajorVersion = 3
	cfg.GL.MinorVersion = 3

	var err error
	f.Ctx, err = samure.CreateContextWithBackend(cfg, &f.gl)
	if err != nil {
		return err
	}

	f.CurrentOutputGeo = f.Ctx.Output(0).Geo()

	f.gl.MakeContextCurrent()

	ml.InitWindowAndRenderer()
	ml.InitManagers()
	gohome.Render.AfterInit()
	ml.SetupStartScene()

	f.running = true
	for f.running {
		gohome.FPSLimit.StartMeasurement()

		f.PollEvents()
		f.Ctx.Update(float64(gohome.FPSLimit.DeltaTime))
		gohome.UpdateMgr.Update(gohome.FPSLimit.DeltaTime)
		gohome.LightMgr.Update()
		gohome.InputMgr.Update(gohome.FPSLimit.DeltaTime)

		for i := 0; i < f.Ctx.LenOutputs(); i++ {
			f.Ctx.RenderOutput(f.Ctx.Output(i))
		}

		gohome.FPSLimit.EndMeasurement()
		gohome.FPSLimit.LimitFPS()
	}

	ml.Quit()

	return nil
}

func (f *SAMUREFramework) Update() {
	f.Ctx.Update(float64(gohome.FPSLimit.DeltaTime))
}

func (f *SAMUREFramework) Terminate() {
	f.Ctx.Destroy()
}

func (f *SAMUREFramework) PollEvents() {
	f.Ctx.ProcessEvents()
}

func (f *SAMUREFramework) CreateWindow(windowWidth, windowHeight int, title string) error {
	// The "window" is the layer surface
	return nil
}

func (f *SAMUREFramework) WindowClosed() bool {
	return f.running
}

func (f *SAMUREFramework) WindowSwap() {
}

func (f *SAMUREFramework) WindowSetSize(size mgl32.Vec2) {

}

func (f *SAMUREFramework) WindowGetSize() mgl32.Vec2 {
	return mgl32.Vec2{
		float32(f.CurrentOutputGeo.W),
		float32(f.CurrentOutputGeo.H),
	}
}

func (f *SAMUREFramework) WindowSetFullscreen(b bool) {

}

func (f *SAMUREFramework) WindowIsFullscreen() bool {
	return true
}

func (f *SAMUREFramework) MonitorGetSize() mgl32.Vec2 {
	return mgl32.Vec2{
		float32(f.CurrentOutputGeo.W),
		float32(f.CurrentOutputGeo.H),
	}
}

func (f *SAMUREFramework) CurserShow() {

}

func (f *SAMUREFramework) CursorHide() {

}

func (f *SAMUREFramework) CursorDisable() {

}

func (f *SAMUREFramework) CursorShown() bool {
	return true
}

func (f *SAMUREFramework) CursorHidden() bool {
	return false
}

func (f *SAMUREFramework) CursorDisabled() bool {
	return false
}

func (f *SAMUREFramework) OpenFile(filePath string) (gohome.File, error) {
	return os.Open(filePath)
}

func (f *SAMUREFramework) LoadLevel(name, path string, loadToGPU bool) *gohome.Level {
	return nil
}

func (f *SAMUREFramework) LoadLevelString(name, contents, fileName string, loadToGPU bool) *gohome.Level {
	return nil
}

func (f *SAMUREFramework) ShowYesNoDialog(title, message string) uint8 {
	return gohome.DIALOG_ERROR
}

func (f *SAMUREFramework) Log(a ...interface{}) {
	fmt.Println(a...)
}

func (f *SAMUREFramework) OnResize(callback func(newWidth, newHeight int)) {

}

func (f *SAMUREFramework) OnMove(callback func(newPosX, newPosY int)) {

}

func (f *SAMUREFramework) OnClose(callback func()) {

}

func (f *SAMUREFramework) OnFocus(callback func(focused bool)) {

}

func (f *SAMUREFramework) StartTextInput() {

}

func (f *SAMUREFramework) GetTextInput() string {
	return ""
}

func (f *SAMUREFramework) EndTextInput() {

}

func (a *SAMUREApp) OnEvent(ctx samure.Context, event interface{}) {
	switch e := event.(type) {
	case samure.EventKeyboardKey:
		if e.Key == samure.KeyEsc {
			a.f.running = false
		}
	}
}

func (a *SAMUREApp) OnRender(ctx samure.Context, layerSurface samure.LayerSurface, outputGeo samure.Rect) {
	a.f.CurrentOutputGeo = outputGeo
	a.f.CurrentLayerSurface = layerSurface

	gohome.RenderMgr.Update()
}

func (*SAMUREApp) OnUpdate(ctx samure.Context, deltaTime float64) {

}
