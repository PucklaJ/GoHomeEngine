package gohome

var (
	MUSIC_SOUND_PATHS = [8]string{
		"",
		"sounds/",
		"sound/",
		"music/",
		"assets/",
		"assets/sounds/",
		"assets/sound/",
		"assets/music/",
	}
)

func (rsmgr *ResourceManager) GetSound(name string) Sound {
	return rsmgr.sounds[name]
}

func (rsmgr *ResourceManager) GetMusic(name string) Music {
	return rsmgr.musics[name]
}

func (rsmgr *ResourceManager) checkMusic(name, path string) bool {
	if resName, ok := rsmgr.resourceFileNames[path]; ok {
		rsmgr.musics[name] = rsmgr.musics[resName]
		ErrorMgr.Message(ERROR_LEVEL_WARNING, "Music", name, "Has already been loaded with this or another name!")
		return false
	}
	if _, ok := rsmgr.musics[name]; ok {
		ErrorMgr.Message(ERROR_LEVEL_WARNING, "Music", name, "Has already been loaded!")
		return false
	}

	return true
}

func (rsmgr *ResourceManager) checkPreloadedMusic(name string) bool {
	for i := 0; i < len(rsmgr.preloader.preloadedMusics); i++ {
		f := rsmgr.preloader.preloadedMusics[i]
		if f.Name == name {
			ErrorMgr.Warning("Music", name, "Has already been preloaded")
			return false
		}
	}

	return true
}

func (rsmgr *ResourceManager) PreloadMusic(name, path string) {
	if rsmgr.checkMusic(name, path) && rsmgr.checkPreloadedMusic(name) {
		rsmgr.preloadedMusics = append(rsmgr.preloadedMusics, preloadedMusic{name, path})
	}
}

func (rsmgr *ResourceManager) LoadMusic(name, path string) Music {
	return rsmgr.loadMusic(name, path, false)
}

func (rsmgr *ResourceManager) loadMusic(name, path string, preloaded bool) Music {
	if !preloaded {
		if !rsmgr.checkMusic(name, path) {
			return nil
		}
	}

	music := Framew.GetAudioManager().LoadMusic(name, path)

	if music != nil {
		if !preloaded {
			rsmgr.musics[name] = music
			rsmgr.resourceFileNames[path] = name
			ErrorMgr.Log("Music", name, "Finished Loading!")
			return music
		} else {
			rsmgr.preloadedMusicChan <- preloadedMusicData{
				preloadedMusic{
					name,
					path,
				},
				music,
			}
		}
	}

	return nil
}

func (rsmgr *ResourceManager) checkSound(name, path string) bool {
	if resName, ok := rsmgr.resourceFileNames[path]; ok {
		rsmgr.sounds[name] = rsmgr.sounds[resName]
		ErrorMgr.Warning("Sound", name, "Has already been loaded with this or another name!")
		return false
	}
	if _, ok := rsmgr.sounds[name]; ok {
		ErrorMgr.Warning("Sound", name, "Has already been loaded!")
		return false
	}

	return true
}

func (rsmgr *ResourceManager) checkPreloadedSound(name string) bool {
	for i := 0; i < len(rsmgr.preloader.preloadedSounds); i++ {
		f := rsmgr.preloader.preloadedSounds[i]
		if f.Name == name {
			ErrorMgr.Warning("Sound", name, "Has already been preloaded")
			return false
		}
	}

	return true
}

func (rsmgr *ResourceManager) PreloadSound(name, path string) {
	if rsmgr.checkSound(name, path) && rsmgr.checkPreloadedSound(name) {
		rsmgr.preloadedSounds = append(rsmgr.preloadedSounds, preloadedSound{name, path})
	}
}

func (rsmgr *ResourceManager) LoadSound(name, path string) Sound {
	return rsmgr.loadSound(name, path, false)
}

func (rsmgr *ResourceManager) loadSound(name, path string, preloaded bool) Sound {
	if !preloaded {
		if !rsmgr.checkSound(name, path) {
			return nil
		}
	}

	sound := Framew.GetAudioManager().LoadSound(name, path)

	if sound != nil {
		if !preloaded {
			rsmgr.sounds[name] = sound
			rsmgr.resourceFileNames[path] = name
			ErrorMgr.Log("Sound", name, "Finished Loading!")
			return sound
		} else {
			rsmgr.preloadedSoundChan <- preloadedSoundData{
				preloadedSound{
					name,
					path,
				},
				sound,
			}
		}
	}

	return nil
}

func (rsmgr *ResourceManager) DeleteSound(name string) {
	sound, ok := rsmgr.sounds[name]
	if ok {
		sound.Terminate()
		delete(rsmgr.sounds, name)
		rsmgr.deleteResourceFileName(name)
		ErrorMgr.Log("Sound", name, "Deleted!")
	} else {
		ErrorMgr.Warning("Sound", name, "Couldn't delete! It has not been loaded!")
	}
}

func (rsmgr *ResourceManager) DeleteMusic(name string) {
	music, ok := rsmgr.musics[name]
	if ok {
		music.Terminate()
		delete(rsmgr.musics, name)
		rsmgr.deleteResourceFileName(name)
		ErrorMgr.Log("Music", name, "Deleted!")
	} else {
		ErrorMgr.Warning("Music", name, "Couldn't delete! It has not been loaded!")
	}
}
