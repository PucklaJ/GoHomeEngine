package gohome

import (
	"strconv"
)

var (
	LEVEL_PATHS = [6]string{
		"",
		"models/",
		"levels/",
		"assets/",
		"assets/models/",
		"assets/levels/",
	}
	MATERIAL_PATHS = [8]string{
		"",
		"models/",
		"levels/",
		"assets/",
		"assets/models/",
		"assets/levels/",
		"materials/",
		"assets/materials/",
	}
)

func (rsmgr *ResourceManager) GetLevel(name string) *Level {
	l := rsmgr.Levels[name]
	return l
}

func (rsmgr *ResourceManager) SetLevel(name string, name1 string) {
	s := rsmgr.Levels[name1]
	if s == nil {
		ErrorMgr.Message(ERROR_LEVEL_ERROR, "Level", name, "Couldn't set to "+name1+" (It is nil)")
		return
	}
	rsmgr.Levels[name] = s
	ErrorMgr.Message(ERROR_LEVEL_LOG, "Level", name, "Set to "+name1)
}

func getNameForAlreadyLoadedLevel(rsmgr *ResourceManager, name string) string {
	var alreadyLoaded = true
	var count = 1
	var newName string
	for alreadyLoaded {
		newName = name + strconv.FormatInt(int64(count), 10)
		_, alreadyLoaded = rsmgr.Levels[newName]
		count++
	}
	return newName
}

func (rsmgr *ResourceManager) DeleteLevel(name string) {
	if _, ok := rsmgr.Levels[name]; ok {
		delete(rsmgr.Levels, name)
		rsmgr.deleteResourceFileName(name)
		ErrorMgr.Log("Level", name, "Deleted!")
	} else {
		ErrorMgr.Warning("Level", name, "Couldn't delete! It hasn't been loaded!")
	}
}

func (rsmgr *ResourceManager) LoadLevelString(name, contents, fileName string, loadToGPU bool) *Level {
	level := Framew.LoadLevelString(name, contents, fileName, loadToGPU)
	if level != nil {
		rsmgr.Levels[level.Name] = level
		ErrorMgr.Log("Level", level.Name, "Finished loading!")
	}
	return level
}

func (rsmgr *ResourceManager) LoadLevel(name, path string, loadToGPU bool) *Level {
	if resName, ok := rsmgr.resourceFileNames[path]; ok {
		rsmgr.Levels[name] = rsmgr.Levels[resName]
		ErrorMgr.Message(ERROR_LEVEL_WARNING, "Level", name, "Has already been loaded with this or another name!")
		return nil
	}

	level := Framew.LoadLevel(name, path, loadToGPU)

	if level != nil {
		rsmgr.Levels[level.Name] = level
		rsmgr.resourceFileNames[path] = level.Name
		ErrorMgr.Log("Level", level.Name, "Finished loading!")
	}
	return level
}

func (rsmgr *ResourceManager) GetModel(name string) *Model3D {
	m := rsmgr.Models[name]
	return m
}

func (rsmgr *ResourceManager) DeleteModel(name string) {
	if _, ok := rsmgr.Models[name]; ok {
		rsmgr.Models[name].Terminate()
		delete(rsmgr.Models, name)
		rsmgr.deleteResourceFileName(name)
		ErrorMgr.Log("Model", name, "Deleted!")
	} else {
		ErrorMgr.Warning("Model", name, "Couldn't delete! It hasn't been loaded!")
	}
}
