// +build js
// +build !linux !windows !darwin
// +build !android

package audio

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/audio/jsaudio"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

func InitAudio() {
	var manager jsaudio.JSAudioManager
	manager.Init()
	gohome.AudioMgr = &manager
}
