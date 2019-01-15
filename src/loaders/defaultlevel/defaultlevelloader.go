package defaultlevel

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	loader "github.com/PucklaMotzer09/GoHomeEngine/src/loaders/obj"
)

type Loader struct {
	gohome.NilFramework
}

func (*Loader) LoadLevel(name, path string, loadToGPU bool) *gohome.Level {
	extension := gohome.GetFileExtension(path)
	if gohome.EqualIgnoreCase(extension, "obj") {
		return loader.LoadLevelOBJ(name, path, loadToGPU)
	}
	gohome.ErrorMgr.Error("Level", name, "The extension "+extension+" is not supported")
	return nil
}

func (*Loader) LoadLevelString(name, contents, fileName string, loadToGPU bool) *gohome.Level {
	return loader.LoadLevelOBJString(name, contents, fileName, loadToGPU)
}
