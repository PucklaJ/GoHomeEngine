package gohome

import (
	"golang.org/x/image/colornames"
	"log"
	"runtime"
)

type MainLoop struct {
	windowWidth  uint32
	windowHeight uint32
	windowTitle  string
	startScene   Scene
}

func (ml MainLoop) Run(fw Framework, r Renderer, ww, wh uint32, wt string, start_scene Scene) {
	runtime.LockOSThread()
	if !ml.init(fw, r, ww, wh, wt, start_scene) {
		ml.quit()
	}
}

func (this *MainLoop) init(fw Framework, r Renderer, ww, wh uint32, wt string, start_scene Scene) bool {
	var err error

	Framew = fw
	Render = r
	this.windowWidth = ww
	this.windowHeight = wh
	this.windowTitle = wt
	this.startScene = start_scene

	if err = Framew.Init(this); err != nil {
		log.Println("Error Initializing Framework:", err)
		return false
	}

	return true

}

func (this *MainLoop) doStuff() {
	this.initWindowAndRenderer()
	this.initManagers()
	this.setupStartScene()

	this.loop()
	this.quit()
}

func (this *MainLoop) setupStartScene() {
	if this.startScene != nil {
		SceneMgr.SwitchScene(this.startScene)
	} else {
		log.Println("Please specify a start scene!")
	}
}

func (this *MainLoop) initWindowAndRenderer() {
	var err error
	if err = Framew.CreateWindow(this.windowWidth, this.windowHeight, this.windowTitle); err != nil {
		log.Println("Error creating window:", err)
		return
	}

	if err = Render.Init(); err != nil {
		log.Println("Error initializing Renderer:", err)
		return
	}
}

func (MainLoop) initManagers() {
	ResourceMgr.Init()
	UpdateMgr.Init()
	RenderMgr.Init()
	LightMgr.Init()
	SceneMgr.Init()
	InputMgr.Init()
	FPSLimit.Init()
}

func (this *MainLoop) loop() {
	for !Framew.WindowClosed() {
		FPSLimit.StartMeasurement()

		this.innerLoop()

		FPSLimit.EndMeasurement()
		FPSLimit.LimitFPS()
	}
}

func (MainLoop) innerLoop() {
	Framew.PollEvents()
	UpdateMgr.Update(FPSLimit.DeltaTime)
	LightMgr.Update()
	InputMgr.Update(FPSLimit.DeltaTime)
	Render.ClearScreen(colornames.Black, 1.0)
	RenderMgr.Update()
	Framew.WindowSwap()
	Framew.Update()
}

func (MainLoop) quit() {
	defer Framew.Terminate()
	defer Render.Terminate()
	defer ResourceMgr.Terminate()
	defer UpdateMgr.Terminate()
	defer RenderMgr.Terminate()
	defer SceneMgr.Terminate()
	if sprite2DMesh != nil {
		defer sprite2DMesh.Terminate()
	}
}

func InitDefaultValues() {
	ResourceMgr.LoadShader(ENTITY3D_SHADER_NAME, "vertex3d.glsl", "fragment3d.glsl", "", "", "", "")
	ResourceMgr.LoadShader(SPRITE2D_SHADER_NAME, "vertex1.glsl", "fragment.glsl", "", "", "", "")
	RenderMgr.SetProjection2D(&Ortho2DProjection{
		Left:   0.0,
		Right:  Framew.WindowGetSize()[0],
		Top:    0.0,
		Bottom: Framew.WindowGetSize()[1],
	})
	RenderMgr.SetProjection3D(&PerspectiveProjection{
		Width:     Framew.WindowGetSize()[0],
		Height:    Framew.WindowGetSize()[1],
		FOV:       70.0,
		NearPlane: 0.1,
		FarPlane:  1000.0,
	})
}
