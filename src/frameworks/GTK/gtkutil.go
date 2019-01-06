package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/loaders/obj"
)

func loadLevelOBJ(rsmgr *gohome.ResourceManager, name, path string, loadToGPU bool) *gohome.Level {
	return loader.LoadLevelOBJ(rsmgr, name, path, loadToGPU)
}

func loadLevelOBJString(rsmgr *gohome.ResourceManager, name, contents, fileName string, loadToGPU bool) *gohome.Level {
	return loader.LoadLevelOBJString(rsmgr, name, contents, fileName, loadToGPU)
}
