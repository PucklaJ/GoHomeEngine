package gohome

// A scene that can be switched
type Scene interface {
	// Init is called when you switch scenes
	Init()
	// Update is called every frame
	Update(delta_time float32)
	// Terminate is called on the current scene when the scene is switched
	Terminate()
}

// An implementation of Scene that does nothing
type NilScene struct {
}

func (this *NilScene) Init() {

}

func (this *NilScene) Update(delta_time float32) {

}

func (this *NilScene) Terminate() {

}

// The manager that handles the scene switching
type SceneManager struct {
	currentScene Scene
}

func (scmgr *SceneManager) Init() {

}

// Switch to another scene
func (scmgr *SceneManager) SwitchScene(scn Scene) {
	if scmgr.currentScene != nil {
		UpdateMgr.RemoveObject(scmgr.currentScene)
		scmgr.currentScene.Terminate()
	}
	if scn != nil {
		scmgr.currentScene = scn
		scmgr.currentScene.Init()
		if scn == scmgr.currentScene {
			UpdateMgr.AddObject(scmgr.currentScene)
		}
	}
	UpdateMgr.BreakUpdateLoop()
}

// Returns the current scene
func (scmgr *SceneManager) GetCurrentScene() Scene {
	return scmgr.currentScene
}

// Calls Terminate on the current scene
func (scmgr *SceneManager) Terminate() {
	if scmgr.currentScene != nil {
		scmgr.currentScene.Terminate()
	}
}

// The SceneManager that should be used for everything
var SceneMgr SceneManager
