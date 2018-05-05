package gohome

import (
	// "fmt"
	// "golang.org/x/image/colornames"
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
			log.Fatalln("Error creating window:", err)
			return
		}
	} else {
		log.Fatalln("Framework is nil!")
	}

	if Render != nil {
		if err = Render.Init(); err != nil {
			log.Fatalln("Error initializing Renderer:", err)
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
	RenderMgr.Update()
	Framew.WindowSwap()
	Framew.Update()
}

func (MainLoop) terminateSprite2DMesh() {
	sprite2DMesh.Terminate()
	sprite2DMesh = nil
}

func (this *MainLoop) Quit() {
	defer Framew.Terminate()
	defer Render.Terminate()
	defer ResourceMgr.Terminate()
	defer UpdateMgr.Terminate()
	defer RenderMgr.Terminate()
	defer SceneMgr.Terminate()
	if sprite2DMesh != nil {
		defer this.terminateSprite2DMesh()
	}
}

func InitDefaultValues() {
	ResourceMgr.LoadShader(ENTITY3D_SHADER_NAME, "vertex3d.glsl", "fragment3d.glsl", "", "", "", "")
	if ResourceMgr.GetShader(ENTITY3D_SHADER_NAME) == nil {
		ResourceMgr.LoadShader("3D No Shadows", "vertex3dNoShadows.glsl", "fragment3dNoShadows.glsl", "", "", "", "")
		if ResourceMgr.GetShader("3D No Shadows") == nil {
			ResourceMgr.LoadShader("3D Simple", "vertex3dNoShadows.glsl", "fragment3dSimple.glsl", "", "", "", "")
			if ResourceMgr.GetShader("3D Simple") != nil {
				ResourceMgr.SetShader(ENTITY3D_SHADER_NAME, "3D Simple")
			}
		} else {
			ResourceMgr.SetShader(ENTITY3D_SHADER_NAME, "3D No Shadows")
		}
	}

	ResourceMgr.LoadShader(SPRITE2D_SHADER_NAME, "vertex1.glsl", "fragment.glsl", "", "", "", "")
}
