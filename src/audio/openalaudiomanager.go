package audio

import "github.com/PucklaMotzer09/gohomeengine/src/gohome"

type OpenALSound struct {

}

func (this *OpenALSound) Play() {

}
func (this *OpenALSound) Pause() {

}
func (this *OpenALSound) Resume() {

}
func (this *OpenALSound) Stop() {

}
func (this *OpenALSound) Terminate() {

}

type OpenALMusic struct {

}

func (this *OpenALMusic) Play() {

}
func (this *OpenALMusic) Pause() {

}
func (this *OpenALMusic) Resume() {

}
func (this *OpenALMusic) Stop() {

}
func (this *OpenALMusic) Terminate() {

}

type OpenALAudioManager struct {

}

func (this *OpenALAudioManager) Init() {

}
func (this *OpenALAudioManager) CreateSound(name string, samples []byte, format uint8, sampleRate uint32) gohome.Sound {
	return &OpenALSound{}
}
func (this *OpenALAudioManager) CreateMusic(name string, samples []byte,format uint8, sampleRate uint32) gohome.Music {
	return &OpenALMusic{}
}
func (this *OpenALAudioManager) Terminate() {

}
