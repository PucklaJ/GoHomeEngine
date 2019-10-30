package gohome

import (
	"github.com/golang/freetype"
)

var (
	// The relative paths in which is search for fonts
	FONT_PATHS = [4]string{
		"",
		"fonts/",
		"assets/",
		"assets/fonts/",
	}
)

// Returns the font width name
func (rsmgr *ResourceManager) GetFont(name string) *Font {
	font := rsmgr.fonts[name]
	return font
}

func (rsmgr *ResourceManager) checkFont(name, path string) bool {
	if resName, ok := rsmgr.resourceFileNames[path]; ok {
		rsmgr.fonts[name] = rsmgr.fonts[resName]
		ErrorMgr.Message(ERROR_LEVEL_WARNING, "Font", name, "Has already been loaded with this or another name!")
		return false
	}
	if _, ok := rsmgr.fonts[name]; ok {
		ErrorMgr.Message(ERROR_LEVEL_WARNING, "Font", name, "Has already been loaded!")
		return false
	}

	return true
}

// Loads the font from path and stores it in name
func (rsmgr *ResourceManager) LoadFont(name, path string) *Font {
	if !rsmgr.checkFont(name, path) {
		return nil
	}

	reader, _, err := OpenFileWithPaths(path, FONT_PATHS[:])
	if err != nil {
		ErrorMgr.Error("Font", name, "Couldn't load: "+err.Error())
		return nil
	}

	data, err := ReadAll(reader)
	if err != nil {
		ErrorMgr.Error("Font", name, "Couldn't read: "+err.Error())
		return nil
	}

	ttf, err := freetype.ParseFont([]byte(data))
	if err != nil {
		ErrorMgr.Error("Font", name, "Couldn't parse: "+err.Error())
		return nil
	}

	var font Font
	font.Init(ttf)

	rsmgr.fonts[name] = &font
	rsmgr.resourceFileNames[path] = name
	ErrorMgr.Log("Font", name, "Finished Loading!")
	return &font

	return nil
}

// Deletes the font with name from the manager
func (rsmgr *ResourceManager) DeleteFont(name string) {
	if _, ok := rsmgr.fonts[name]; ok {
		delete(rsmgr.fonts, name)
		rsmgr.deleteResourceFileName(name)
		ErrorMgr.Log("Font", name, "Deleted!")
	} else {
		ErrorMgr.Warning("Font", name, "Couldn't delete! It hasn't been loaded!")
	}
}
