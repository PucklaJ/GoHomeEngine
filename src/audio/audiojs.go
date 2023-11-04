//go:build js && (!linux || !windows || !darwin) && !android
// +build js
// +build !linux !windows !darwin
// +build !android

package audio

import (
	"github.com/PucklaJ/GoHomeEngine/src/audio/jsaudio"
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
)

func InitAudio() {
	var manager jsaudio.JSAudioManager
	manager.Init()
	gohome.AudioMgr = &manager
}
