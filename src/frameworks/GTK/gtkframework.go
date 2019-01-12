package framework

import (
	"errors"
	"github.com/PucklaMotzer09/GoHomeEngine/src/loaders/defaultlevel"
	"os"
	"time"

	"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/GTK/gtk"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type GTKFramework struct {
	defaultlevel.Loader
	startOtherThanPaint time.Time
	endOtherThanPaint   time.Time

	prevMousePos [2]int16

	isFullscreen bool

	UseWholeWindowAsGLArea bool
	UseExternalWindow      bool
	MenuBarFix             bool
}

func (this *GTKFramework) InitStuff(ml *gohome.MainLoop) {
	this.DefaultWindowCreation(ml)
	this.AfterWindowCreation(ml)
}

func (this *GTKFramework) DefaultWindowCreation(ml *gohome.MainLoop) {
	ml.InitWindow()
	gtk.GetWindow().ToContainer().Add(gtk.GetGLArea().ToWidget())
}

func (this *GTKFramework) AfterWindowCreation(ml *gohome.MainLoop) {
	ml.InitRenderer()
	ml.InitManagers()
	gohome.Render.AfterInit()
}

func (this *GTKFramework) Init(ml *gohome.MainLoop) error {

	this.isFullscreen = false
	gtk.OnRender = gtkgo_gl_area_render
	gtk.OnMotion = gtkgo_gl_area_motion_notify
	gtk.OnUseWholeScreen = useWholeWindowAsGLArea
	gtk.Init()
	if this.UseWholeWindowAsGLArea && !this.UseExternalWindow {
		this.InitStuff(ml)
	} else {
		gohome.ErrorMgr.Init()
	}
	ml.SetupStartScene()
	if gtk.GetGLArea().ToWidget().GetParent().IsNULL() {
		return errors.New("GLArea has not been added!")
	}

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
	gtk.MainQuit()
}

func (this *GTKFramework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	return gtk.CreateWindow(windowWidth, windowHeight, title)
}
func (this *GTKFramework) WindowClosed() bool {
	return false
}

func (this *GTKFramework) WindowSetSize(size mgl32.Vec2) {
	gtk.WindowSetSize(size)
}

func (this *GTKFramework) WindowGetSize() mgl32.Vec2 {
	w, h := gtk.GetGLArea().ToWidget().GetSize()
	return [2]float32{float32(w), float32(h)}
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

func (this *GTKFramework) OpenFile(file string) (gohome.File, error) {
	return os.Open(file)
}

func (this *GTKFramework) ShowYesNoDialog(title, message string) uint8 {
	return gohome.DIALOG_CANCELLED
}

func (this *GTKFramework) GetGtkWindow() gtk.Window {
	return gtk.GetWindow()
}

func (this *GTKFramework) MonitorGetSize() mgl32.Vec2 {
	return this.WindowGetSize()
}

func (this *GTKFramework) InitExternalDefault(window *gtk.Window, glarea *gtk.GLArea) {
	if window != nil {
		window.ConfigureParameters()
		window.ConnectSignals()
		gtk.SetWindow(*window)
	}
	if glarea != nil {
		glarea.Configure()
		gtk.SetGLArea(*glarea)
	}
	this.AfterWindowCreation(&gohome.MainLop)
}
