// +build !android

package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/audio"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

var audioManager audio.OpenALAudioManager

func (this *SDL2Framework) GetAudioManager() gohome.AudioManager {
	return &audioManager
}
