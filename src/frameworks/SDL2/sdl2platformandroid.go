// +build android

package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

func (this *SDL2Framework) GetAudioManager() gohome.AudioManager {
	return &gohome.NilAudioManager{}
}
