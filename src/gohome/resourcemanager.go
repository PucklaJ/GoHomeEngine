package gohome

import (
	"github.com/PucklaMotzer09/tmx"
	"github.com/blezek/tga"
)

// The manager that handles all resources
type ResourceManager struct {
	textures          map[string]Texture
	shaders           map[string]Shader
	Models            map[string]*Model3D
	Levels            map[string]*Level
	fonts             map[string]*Font
	musics            map[string]Music
	sounds            map[string]Sound
	tmxmaps           map[string]*tmx.Map
	resourceFileNames map[string]string

	// Wether models can have the same name
	LoadModelsWithSameName bool
}

// Initialises all values of the manager
func (rsmgr *ResourceManager) Init() {
	rsmgr.textures = make(map[string]Texture)
	rsmgr.shaders = make(map[string]Shader)
	rsmgr.Models = make(map[string]*Model3D)
	rsmgr.Levels = make(map[string]*Level)
	rsmgr.fonts = make(map[string]*Font)
	rsmgr.musics = make(map[string]Music)
	rsmgr.sounds = make(map[string]Sound)
	rsmgr.tmxmaps = make(map[string]*tmx.Map)
	rsmgr.resourceFileNames = make(map[string]string)

	tga.RegisterFormat()
	rsmgr.LoadModelsWithSameName = false
}

// Cleans everything up
func (rsmgr *ResourceManager) Terminate() {
	for k, v := range rsmgr.shaders {
		v.Terminate()
		delete(rsmgr.shaders, k)
	}

	for k, v := range rsmgr.textures {
		v.Terminate()
		delete(rsmgr.textures, k)
	}

	for k, v := range rsmgr.Models {
		v.Terminate()
		delete(rsmgr.Models, k)
	}

	for k := range rsmgr.Levels {
		delete(rsmgr.Levels, k)
	}

	for k, v := range rsmgr.sounds {
		v.Terminate()
		delete(rsmgr.sounds, k)
	}

	for k, v := range rsmgr.musics {
		v.Terminate()
		delete(rsmgr.musics, k)
	}

	for k := range rsmgr.tmxmaps {
		delete(rsmgr.tmxmaps, k)
	}
}

func (rsmgr *ResourceManager) deleteResourceFileName(name string) {
	for k := range rsmgr.resourceFileNames {
		if rsmgr.resourceFileNames[k] == name {
			delete(rsmgr.resourceFileNames, k)
			return
		}
	}
}

// The ResourceManager that should be used for everything
var ResourceMgr ResourceManager
