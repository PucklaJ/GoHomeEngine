// +build linux windows darwin
// +build !android
// +build !js

package audio

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/audio/openal"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

func InitAudio() {
	var manager openal.OpenALAudioManager
	manager.Init()
	gohome.AudioMgr = &manager
}
