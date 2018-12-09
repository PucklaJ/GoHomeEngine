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
	LoadGeneratedShader3D(0)
	if ResourceMgr.GetShader(ENTITY3D_SHADER_NAME) == nil {
		LoadGeneratedShader3D(SHADER_FLAG_NO_SHADOWS)
		if ResourceMgr.GetShader("3D NoShadows") == nil {
			LoadGeneratedShader3D(SHADER_FLAG_NO_SHADOWS | SHADER_FLAG_NOUV)
			if ResourceMgr.GetShader("3D NoUV NoShadows") != nil {
				ResourceMgr.SetShader(ENTITY3D_SHADER_NAME, "3D NoUV NoShadows")
			}
		} else {
			ResourceMgr.SetShader(ENTITY3D_SHADER_NAME, "3D NoShadows")
		}
	}
}

func Init2DShaders() {
	ResourceMgr.LoadShaderSource(SPRITE2D_SHADER_NAME, SPRITE_2D_SHADER_VERTEX_SOURCE_OPENGL, SPRITE_2D_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
}

func InitDefaultValues() {
	Init3DShaders()
	Init2DShaders()
}

var MainLop MainLoop
