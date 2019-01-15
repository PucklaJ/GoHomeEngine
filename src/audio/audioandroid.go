// +build android
// +build !linux !windows !darwin
// +build !js

package audio

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/audio/sdlmixer"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

func InitAudio() {
	var manager sdlmixer.MixerAudioManager
	manager.Init()
	gohome.AudioMgr = &manager
}
