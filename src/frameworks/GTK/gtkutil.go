package framework

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/loaders/obj"
)

func loadLevelOBJ(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	return loader.LoadLevelOBJ(rsmgr, name, path, preloaded, loadToGPU)
}

func loadLevelOBJString(rsmgr *gohome.ResourceManager, name, contents, fileName string, preloaded, loadToGPU bool) *gohome.Level {
	return loader.LoadLevelOBJString(rsmgr, name, contents, fileName, preloaded, loadToGPU)
}
