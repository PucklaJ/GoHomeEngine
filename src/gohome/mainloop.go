package gohome

import (
	"golang.org/x/image/colornames"
	"log"
	"runtime"
)

type MainLoop struct {
}

func (ml MainLoop) Run(fw Framework, r Renderer, ww, wh uint32, wt string, start_scene Scene) {
	runtime.LockOSThread()
	if ml.init(fw, r, ww, wh, wt, start_scene) {
		ml.loop()
	}
	ml.quit()
}

func (MainLoop) init(fw Framework, r Renderer, ww, wh uint32, wt string, start_scene Scene) bool {
	var err error

	Framew = fw
	Render = r

	if err = Framew.Init(); err != nil {
		log.Fatalln("Error Initializing Framew:", err)
		return false
	}

	if err = Framew.CreateWindow(ww, wh, wt); err != nil {
		log.Fatalln("Error creating window:", err)
		return false
	}

	if err = Render.Init(); err != nil {
		log.Fatalln("Error initializing Renderer:", err)
		return false
	}

	ResourceMgr.Init()
	UpdateMgr.Init()
	RenderMgr.Init()
	LightMgr.Init()
	SceneMgr.Init()
	InputMgr.Init()
	FPSLimit.Init()

	if start_scene != nil {
		SceneMgr.SwitchScene(start_scene)
	} else {
		log.Println("Please specify a start scene!")
	}

	return true
}

func (MainLoop) loop() bool {
	for !Framew.WindowClosed() {
		FPSLimit.StartMeasurement()

		Framew.PollEvents()
		UpdateMgr.Update(FPSLimit.DeltaTime)
		InputMgr.Update(FPSLimit.DeltaTime)
		Render.ClearScreen(colornames.Black, 1.0)
		RenderMgr.Update()
		Framew.WindowSwap()
		Framew.Update()

		FPSLimit.EndMeasurement()
		FPSLimit.LimitFPS()
	}
	return true
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
