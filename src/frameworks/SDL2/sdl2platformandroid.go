// +build android

package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

var audioManager MixerAudioManager

func (this *SDL2Framework) GetAudioManager() gohome.AudioManager {
	return &audioManager
}
