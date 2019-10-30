package gohome

import (
	"runtime"
)

// The main struct of the engine handling the start and end
type MainLoop struct {
	windowWidth  int
	windowHeight int
	windowTitle  string
	startScene   Scene
}

// The method that starts everything, should be called in the main function
func (ml *MainLoop) Run(fw Framework, r Renderer, ww, wh int, wt string, start_scene Scene) {
	runtime.LockOSThread()
	if !ml.Init(fw, r, ww, wh, wt, start_scene) {
		ml.Quit()
	}
}

// Initialises the values and calls init on the framework
func (this *MainLoop) Init(fw Framework, r Renderer, ww, wh int, wt string, start_scene Scene) bool {

	Framew = fw
	Render = r
	this.windowWidth = ww
	this.windowHeight = wh
	this.windowTitle = wt
	this.startScene = start_scene
	if err := Framew.Init(this); err != nil {
		ErrorMgr.MessageError(ERROR_LEVEL_FATAL, "Framework", "Initialisation", err)
		return false
	}

	return true

}

// Initialises the managers and starts the loop. Will be called from the framework
func (this *MainLoop) DoStuff() {
	this.InitWindowAndRenderer()
	this.InitManagers()
	Render.AfterInit()
	this.SetupStartScene()
	this.Loop()
	this.Quit()
}

// Sets up the start scene
func (this *MainLoop) SetupStartScene() {
	if this.startScene != nil {
		SceneMgr.SwitchScene(this.startScene)
	} else {
		ErrorMgr.Message(ERROR_LEVEL_ERROR, "Scene", "", "Please specify a start scene!")
	}
}

// Initialises the window
func (this *MainLoop) InitWindow() bool {
	var err error
	if Framew != nil {
		if err = Framew.CreateWindow(this.windowWidth, this.windowHeight, this.windowTitle); err != nil {
			ErrorMgr.MessageError(ERROR_LEVEL_FATAL, "WindowCreation", "", err)
			return false
		}
	} else {
		ErrorMgr.Message(ERROR_LEVEL_FATAL, "WindowCreation", "", "Framework is nil!")
		return false
	}
	return true
}

// Initialises the renderer
func (this *MainLoop) InitRenderer() {
	var err error
	if Render != nil {
		if err = Render.Init(); err != nil {
			ErrorMgr.MessageError(ERROR_LEVEL_FATAL, "RendererInitialisation", "", err)
			return
		}
	}
}

// Initialises the window and the renderer
func (this *MainLoop) InitWindowAndRenderer() {
	if this.InitWindow() {
		this.InitRenderer()
	}
}

// Initialises all managers
func (MainLoop) InitManagers() {
	ErrorMgr.Init()
	ResourceMgr.Init()
	UpdateMgr.Init()
	RenderMgr.Init()
	LightMgr.Init()
	SceneMgr.Init()
	InputMgr.Init()
	FPSLimit.Init()
}

// One iteration of the loop
func (this *MainLoop) LoopOnce() {
	FPSLimit.StartMeasurement()
	this.InnerLoop()
	FPSLimit.EndMeasurement()
	FPSLimit.LimitFPS()
}

// Calls LoopOnce as long as the window is open
func (this *MainLoop) Loop() {
	for !Framew.WindowClosed() {
		this.LoopOnce()
	}
}

// Will be called in LoopOnce
func (MainLoop) InnerLoop() {
	Framew.PollEvents()
	UpdateMgr.Update(FPSLimit.DeltaTime)
	LightMgr.Update()
	InputMgr.Update(FPSLimit.DeltaTime)
	RenderMgr.Update()
	Framew.Update()
}

func (MainLoop) terminateSprite2DMesh() {
	sprite2DMesh.Terminate()
	sprite2DMesh = nil
}

// Terminates all resources of the engine
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
	defer ErrorMgr.Terminate()
	if AudioMgr != nil {
		defer AudioMgr.Terminate()
	}
}

// Initialises the 3D shaders
func Init3DShaders() {
	if shader := LoadGeneratedShader3D(SHADER_TYPE_3D, 0); shader == nil {
		if shader = LoadGeneratedShader3D(SHADER_TYPE_3D, SHADER_FLAG_NO_SHADOWS); shader == nil {
			if shader = LoadGeneratedShader3D(SHADER_TYPE_3D, SHADER_FLAG_NO_SHADOWS|SHADER_FLAG_NOUV); shader != nil {
				ResourceMgr.SetShader(ENTITY3D_SHADER_NAME, shader.GetName())
			}
		} else {
			ResourceMgr.SetShader(ENTITY3D_SHADER_NAME, shader.GetName())
		}
	}
}

// Initialises the 2D shaders
func Init2DShaders() {
	LoadGeneratedShader2D(SHADER_TYPE_SPRITE2D, 0)
	LoadGeneratedShader2D(SHADER_TYPE_SHAPE2D, 0)
	LoadGeneratedShader2D(SHADER_TYPE_TEXT2D, 0)
}

// Initialises all shaders
func InitDefaultValues() {
	Init3DShaders()
	Init2DShaders()
}

// The MainLoop that should be used for everything
var MainLop MainLoop
