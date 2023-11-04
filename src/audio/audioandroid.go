//go:build android && (!linux || !windows || !darwin) && !js
// +build android
// +build !linux !windows !darwin
// +build !js

package audio

import (
	"github.com/PucklaJ/GoHomeEngine/src/audio/sdlmixer"
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
)

func InitAudio() {
	var manager sdlmixer.MixerAudioManager
	manager.Init()
	gohome.AudioMgr = &manager
}
