package gohome

import (
	"github.com/PucklaMotzer09/tmx"
	"io"
)

var (
	// The relative paths in which will be searched for tmx files
	TMX_MAP_PATHS = [4]string{
		"",
		"maps/",
		"assets/",
		"assets/maps/",
	}
)

// Deletes the tmx map with name from the manager
func (rsmgr *ResourceManager) DeleteTMXMap(name string) {
	_, ok := rsmgr.tmxmaps[name]
	if ok {
		delete(rsmgr.tmxmaps, name)
		rsmgr.deleteResourceFileName(name)
		ErrorMgr.Log("TMXMap", name, "Deleted!")
	} else {
		ErrorMgr.Warning("TMXMap", name, "Couldn't delete! It has not been loaded!")
	}
}

func (rsmgr *ResourceManager) checkTMXMap(name, path string) bool {
	if name1, ok := rsmgr.resourceFileNames[path]; ok {
		rsmgr.tmxmaps[name] = rsmgr.tmxmaps[name1]
		ErrorMgr.Warning("TMXMap", name, "Has alreay been loaded with this or another name!")
		return false
	}
	if _, ok := rsmgr.tmxmaps[name]; ok {
		ErrorMgr.Warning("TMXMap", name, "Has already been loaded!")
		return false
	}
	return true
}

// Loads the tmx map from path and stores it in name
func (rsmgr *ResourceManager) LoadTMXMap(name, path string) *tmx.Map {
	if !rsmgr.checkTMXMap(name, path) {
		return nil
	}

	file, fileName, err := OpenFileWithPaths(path, TMX_MAP_PATHS[:])
	if err != nil {
		ErrorMgr.MessageError(ERROR_LEVEL_ERROR, "TMXMap", name, err)
		return nil
	}

	tmx.OpenFileFunction = func(filename string) (io.ReadCloser, error) {
		return Framew.OpenFile(filename)
	}
	tmxmap, err := tmx.LoadReader(file, fileName)
	if err != nil {
		ErrorMgr.MessageError(ERROR_LEVEL_ERROR, "TMXMap", name, err)
		return nil
	}

	rsmgr.tmxmaps[name] = tmxmap
	rsmgr.resourceFileNames[path] = name

	ErrorMgr.Log("TMXMap", name, "Finished Loading!")
	return tmxmap
}

// Returns the tmx map with name
func (rsmgr *ResourceManager) GetTMXMap(name string) *tmx.Map {
	return rsmgr.tmxmaps[name]
}
