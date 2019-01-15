package gohome

import (
	"runtime"
)

type MainLoop struct {
	windowWidth  uint32
	windowHeight uint32
	windowTitle  string
	startScene   Scene
}

func (ml *MainLoop) Run(fw Framework, r Renderer, ww, wh uint32, wt string, start_scene Scene) {
	runtime.LockOSThread()
	if !ml.Init(fw, r, ww, wh, wt, start_scene) {
		ml.Quit()
	}
}

func (this *MainLoop) Init(fw Framework, r Renderer, ww, wh uint32, wt string, start_scene Scene) bool {

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

func (this *MainLoop) DoStuff() {
	this.InitWindowAndRenderer()
	this.InitManagers()
	Render.AfterInit()
	this.SetupStartScene()
	this.Loop()
	this.Quit()
}

func (this *MainLoop) SetupStartScene() {
	if this.startScene != nil {
		SceneMgr.SwitchScene(this.startScene)
	} else {
		ErrorMgr.Message(ERROR_LEVEL_ERROR, "Scene", "", "Please specify a start scene!")
	}
}

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

func (this *MainLoop) InitRenderer() {
	var err error
	if Render != nil {
		if err = Render.Init(); err != nil {
			ErrorMgr.MessageError(ERROR_LEVEL_FATAL, "RendererInitialisation", "", err)
			return
		}
	}
}

func (this *MainLoop) InitWindowAndRenderer() {
	if this.InitWindow() {
		this.InitRenderer()
	}
}

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

func (this *MainLoop) LoopOnce() {
	FPSLimit.StartMeasurement()
	this.InnerLoop()
	FPSLimit.EndMeasurement()
	FPSLimit.LimitFPS()
}

func (this *MainLoop) Loop() {
	for !Framew.WindowClosed() {
		this.LoopOnce()
	}
}

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
}

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

func Init2DShaders() {
	LoadGeneratedShader2D(SHADER_TYPE_SPRITE2D, 0)
	LoadGeneratedShader2D(SHADER_TYPE_SHAPE2D, 0)
	LoadGeneratedShader2D(SHADER_TYPE_TEXT2D, 0)
}

func InitDefaultValues() {
	Init3DShaders()
	Init2DShaders()
}

var MainLop MainLoop
