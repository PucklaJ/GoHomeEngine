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

func (rsmgr *ResourceManager) PreloadLevel(name, path string, loadToGPU bool) {
	level := preloadedLevel{
		name,
		path,
		loadToGPU,
		false,
	}
	if !rsmgr.checkPreloadedLevel(&level) {
		return
	}
	rsmgr.preloader.preloadedLevels = append(rsmgr.preloader.preloadedLevels, level)
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

func (rsmgr *ResourceManager) checkPreloadedLevel(level *preloadedLevel) bool {
	var alreadyLoaded = false
	if _, alreadyLoaded = rsmgr.Levels[level.Name]; alreadyLoaded && !rsmgr.LoadModelsWithSameName {
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Level", level.Name, "Has already been loaded!")
		return false
	}
	if alreadyLoaded {
		(*level).Name = getNameForAlreadyLoadedLevel(rsmgr, level.Name)
	}
	if resName, ok := rsmgr.resourceFileNames[level.Path]; ok {
		rsmgr.textures[level.Name] = rsmgr.textures[resName]
		ErrorMgr.Message(ERROR_LEVEL_WARNING, "Level", level.Name, "Has already been loaded with this or another name!")
		return false
	}
	for i := 0; i < len(rsmgr.preloadedLevels); i++ {
		if rsmgr.preloadedLevels[i].Name == level.Name {
			ErrorMgr.Message(ERROR_LEVEL_LOG, "Level", level.Name, "Has already been preloaded!")
			return false
		} else if rsmgr.preloadedLevels[i].Path == level.Path {
			ErrorMgr.Message(ERROR_LEVEL_WARNING, "Level", level.Name, "Has already been preloaded with this or another name!")
			level.fileAlreadyPreloaded = true
			return true
		}
	}

	level.fileAlreadyPreloaded = false

	return true
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
	level := rsmgr.loadLevelString(name, contents, fileName, false, loadToGPU)
	if level != nil {
		rsmgr.Levels[level.Name] = level
		ErrorMgr.Log("Level", level.Name, "Finished loading!")
	}
	return level
}

func (rsmgr *ResourceManager) LoadLevel(name, path string, loadToGPU bool) *Level {
	level := rsmgr.loadLevel(name, path, false, loadToGPU)
	if level != nil {
		rsmgr.Levels[level.Name] = level
		rsmgr.resourceFileNames[path] = level.Name
		ErrorMgr.Log("Level", level.Name, "Finished loading!")
	}
	return level
}

func (rsmgr *ResourceManager) loadLevelString(name, contents, fileName string, preloaded, loadToGPU bool) *Level {
	return Framew.LoadLevelString(rsmgr, name, contents, fileName, preloaded, loadToGPU)
}

func (rsmgr *ResourceManager) loadLevel(name, path string, preloaded, loadToGPU bool) *Level {
	if !preloaded {
		if resName, ok := rsmgr.resourceFileNames[path]; ok {
			rsmgr.Levels[name] = rsmgr.Levels[resName]
			ErrorMgr.Message(ERROR_LEVEL_WARNING, "Level", name, "Has already been loaded with this or another name!")
			return nil
		}
	}
	return Framew.LoadLevel(rsmgr, name, path, preloaded, loadToGPU)
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
