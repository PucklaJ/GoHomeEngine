//go:build (linux || windows || darwin) && !android && !js
// +build linux windows darwin
// +build !android
// +build !js

package audio

import (
	"github.com/PucklaJ/GoHomeEngine/src/audio/openal"
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
)

func InitAudio() {
	var manager openal.OpenALAudioManager
	manager.Init()
	gohome.AudioMgr = &manager
}
