package gohome

import (
	// "fmt"
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
	if !ml.Init(fw, r, ww, wh, wt, start_scene) {
		ml.Quit()
	}
}

func (this *MainLoop) Init(fw Framework, r Renderer, ww, wh uint32, wt string, start_scene Scene) bool {
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

func (this *MainLoop) DoStuff() {
	this.InitWindowAndRenderer()
	this.InitManagers()
	this.SetupStartScene()

	this.Loop()
	this.Quit()
}

func (this *MainLoop) SetupStartScene() {
	if this.startScene != nil {
		SceneMgr.SwitchScene(this.startScene)
	} else {
		log.Println("Please specify a start scene!")
	}
}

func (this *MainLoop) InitWindowAndRenderer() {
	var err error
	if Framew != nil {
		if err = Framew.CreateWindow(this.windowWidth, this.windowHeight, this.windowTitle); err != nil {
			log.Println("Error creating window:", err)
			return
		}
	} else {
		log.Fatalln("Framework is nil!")
	}

	if Render != nil {
		if err = Render.Init(); err != nil {
			log.Println("Error initializing Renderer:", err)
			return
		}
	}

}

func (MainLoop) InitManagers() {
	ResourceMgr.Init()
	UpdateMgr.Init()
	RenderMgr.Init()
	LightMgr.Init()
	SceneMgr.Init()
	InputMgr.Init()
	FPSLimit.Init()
}

func (this *MainLoop) Loop() {
	for !Framew.WindowClosed() {
		FPSLimit.StartMeasurement()

		this.InnerLoop()

		FPSLimit.EndMeasurement()
		FPSLimit.LimitFPS()
	}
}

func (MainLoop) InnerLoop() {
	Framew.PollEvents()
	UpdateMgr.Update(FPSLimit.DeltaTime)
	LightMgr.Update()
	InputMgr.Update(FPSLimit.DeltaTime)
	Render.ClearScreen(colornames.Blue, 1.0)
	RenderMgr.Update()
	Framew.WindowSwap()
	Framew.Update()
}

func (MainLoop) Quit() {
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
