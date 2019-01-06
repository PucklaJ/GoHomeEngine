package defaultlevel

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	loader "github.com/PucklaMotzer09/GoHomeEngine/src/loaders/obj"
	"strings"
)

type Loader struct {
	gohome.NilFramework
}

func (*Loader) LoadLevel(rsmgr *gohome.ResourceManager, name, path string, loadToGPU bool) *gohome.Level {
	extension := getFileExtension(path)
	if equalIgnoreCase(extension, "obj") {
		return loader.LoadLevelOBJ(rsmgr, name, path, loadToGPU)
	}
	gohome.ErrorMgr.Error("Level", name, "The extension "+extension+" is not supported")
	return nil
}

func (*Loader) LoadLevelString(rsmgr *gohome.ResourceManager, name, contents, fileName string, loadToGPU bool) *gohome.Level {
	return loader.LoadLevelOBJString(rsmgr, name, contents, fileName, loadToGPU)
}

func equalIgnoreCase(str1, str string) bool {
	if len(str1) != len(str) {
		return false
	}
	for i := 0; i < len(str1); i++ {
		if str1[i] != str[i] {
			if str1[i] >= 65 && str1[i] <= 90 {
				if str[i] >= 97 && str[i] <= 122 {
					if str1[i]+32 != str[i] {
						return false
					}
				} else {
					return false
				}
			} else if str1[i] >= 97 && str1[i] <= 122 {
				if str[i] >= 65 && str[i] <= 90 {
					if str1[i]-32 != str[i] {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

func getFileExtension(file string) string {
	index := strings.LastIndex(file, ".")
	if index == -1 {
		return ""
	}
	return file[index+1:]
}
