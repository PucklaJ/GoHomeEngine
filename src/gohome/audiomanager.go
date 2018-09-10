package gohome

import "time"

type Sound interface {
	Play(loop bool)
	Pause()
	Resume()
	Stop()
	Terminate()
	IsPlaying() bool
	GetPlayingDuration() time.Duration
	GetDuration() time.Duration
	SetVolume(vol float32)
	GetVolume() float32
}

type Music interface {
	Play(loop bool)
	Pause()
	Resume()
	Stop()
	Terminate()
	IsPlaying() bool
	GetPlayingDuration() time.Duration
	GetDuration() time.Duration
	SetVolume(vol float32)
	GetVolume() float32
}

const (
	AUDIO_FORMAT_STEREO8  uint8 = iota
	AUDIO_FORMAT_STEREO16 uint8 = iota
	AUDIO_FORMAT_MONO8    uint8 = iota
	AUDIO_FORMAT_MONO16   uint8 = iota
	AUDIO_FORMAT_UNKNOWN  uint8 = iota
)

type AudioManager interface {
	Init()
	CreateSound(name string, samples []byte, format uint8, sampleRate uint32) Sound
	CreateMusic(name string, samples []byte, format uint8, sampleRate uint32) Music
	SetVolume(vol float32)
	GetVolume() float32
	Terminate()
}

type NilSound struct {
}

func (*NilSound) Play(loop bool) {

}
func (*NilSound) Pause() {

}
func (*NilSound) Resume() {

}
func (*NilSound) Stop() {

}
func (*NilSound) Terminate() {

}
func (*NilSound) IsPlaying() bool {
	return false
}
func (*NilSound) GetPlayingDuration() time.Duration {
	return time.Second * 0
}
func (*NilSound) GetDuration() time.Duration {
	return time.Second * 0
}

func (*NilSound) SetVolume(vol float32) {

}

func (*NilSound) GetVolume() float32 {
	return 1.0
}

type NilMusic struct {
}

func (*NilMusic) Play(loop bool) {

}
func (*NilMusic) Pause() {

}
func (*NilMusic) Resume() {

}
func (*NilMusic) Stop() {

}
func (*NilMusic) Terminate() {

}
func (*NilMusic) IsPlaying() bool {
	return false
}
func (*NilMusic) GetPlayingDuration() time.Duration {
	return time.Second * 0
}
func (*NilMusic) GetDuration() time.Duration {
	return time.Second * 0
}

func (*NilMusic) SetVolume(vol float32) {

}

func (*NilMusic) GetVolume() float32 {
	return 1.0
}

type NilAudioManager struct {
}

func (*NilAudioManager) Init() {

}
func (*NilAudioManager) CreateSound(name string, samples []byte, format uint8, sampleRate uint32) Sound {
	return &NilSound{}
}
func (*NilAudioManager) CreateMusic(name string, samples []byte, format uint8, sampleRate uint32) Music {
	return &NilMusic{}
}
func (*NilAudioManager) SetVolume(vol float32) {

}
func (*NilAudioManager) GetVolume() float32 {
	return 1.0
}
func (*NilAudioManager) Terminate() {

}
