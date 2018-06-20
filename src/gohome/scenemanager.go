package gohome

type Scene interface {
	Init()
	Update(delta_time float32)
	Terminate()
}

type SceneManager struct {
	currentScene Scene
}

func (scmgr *SceneManager) Init() {

}

func (scmgr *SceneManager) SwitchScene(scn Scene) {
	if scmgr.currentScene != nil {
		UpdateMgr.RemoveObject(scmgr.currentScene)
		scmgr.currentScene.Terminate()
	}
	if scn != nil {
		scmgr.currentScene = scn
		scmgr.currentScene.Init()
		UpdateMgr.AddObject(scmgr.currentScene)
	}
}

func (scmgr *SceneManager) GetCurrentScene() Scene {
	return scmgr.currentScene
}

func (scmgr *SceneManager) Terminate() {
	if scmgr.currentScene != nil {
		scmgr.currentScene.Terminate()
	}
}

var SceneMgr SceneManager
