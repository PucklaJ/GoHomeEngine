package gohome

import (
	"github.com/golang/freetype"
)

var (
	FONT_PATHS = [4]string{
		"",
		"fonts/",
		"assets/",
		"assets/fonts/",
	}
)

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

func (rsmgr *ResourceManager) checkPreloadedFont(name string) bool {
	for i := 0; i < len(rsmgr.preloader.preloadedFonts); i++ {
		f := rsmgr.preloader.preloadedFonts[i]
		if f.Name == name {
			ErrorMgr.Warning("Font", name, "Has already been preloaded")
			return false
		}
	}

	return true
}

func (rsmgr *ResourceManager) PreloadFont(name, path string) {
	if rsmgr.checkFont(name, path) && rsmgr.checkPreloadedFont(name) {
		rsmgr.preloader.preloadedFonts = append(rsmgr.preloader.preloadedFonts, preloadedFont{name, path})
	}
}

func (rsmgr *ResourceManager) LoadFont(name, path string) *Font {
	return rsmgr.loadFont(name, path, false)
}

func (rsmgr *ResourceManager) loadFont(name, path string, preloaded bool) *Font {
	if !preloaded {
		if !rsmgr.checkFont(name, path) {
			return nil
		}
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

	if !preloaded {
		rsmgr.fonts[name] = &font
		rsmgr.resourceFileNames[path] = name
		ErrorMgr.Log("Font", name, "Finished Loading!")
		return &font
	} else {
		rsmgr.preloader.preloadedFontChan <- preloadedFontData{
			preloadedFont{
				name,
				path,
			},
			font,
		}
	}

	return nil
}

func (rsmgr *ResourceManager) DeleteFont(name string) {
	if _, ok := rsmgr.fonts[name]; ok {
		delete(rsmgr.fonts, name)
		rsmgr.deleteResourceFileName(name)
		ErrorMgr.Log("Font", name, "Deleted!")
	} else {
		ErrorMgr.Warning("Font", name, "Couldn't delete! It hasn't been loaded!")
	}
}
