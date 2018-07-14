package gohome

type Sound interface {
	Play()
	Pause()
	Resume()
	Stop()
	Terminate()
}

type Music interface {
	Play()
	Pause()
	Resume()
	Stop()
	Terminate()
}

const (
	AUDIO_FORMAT_STEREO8 uint8 = iota
	AUDIO_FORMAT_STEREO16 uint8 = iota
	AUDIO_FORMAT_MONO8 uint8 = iota
	AUDIO_FORMAT_MONO16 uint8 = iota
	AUDIO_FORMAT_UNKNOWN uint8 = iota
)

type AudioManager interface {
	Init()
	CreateSound(name string, samples []byte, format uint8, sampleRate uint32) Sound
	CreateMusic(name string, samples []byte,format uint8, sampleRate uint32) Music
	Terminate()
}

type NilSound struct {

}

func (*NilSound) Play() {

}

func (*NilSound) Pause() {

}

func (*NilSound) Resume() {

}

func (*NilSound) Stop() {

}

func (*NilSound) Terminate() {

}

type NilMusic struct {

}

func (*NilMusic) Play() {

}

func (*NilMusic) Pause() {

}

func (*NilMusic) Resume() {

}

func (*NilMusic) Stop() {

}

func (*NilMusic) Terminate() {

}

type NilAudioManager struct {

}

func (*NilAudioManager) Init() {

}

func (*NilAudioManager) CreateSound(name string, samples []byte, format uint8, sampleRate uint32) Sound {
	return &NilSound{}
}

func (*NilAudioManager) CreateMusic(name string,samples []byte, format uint8, sampleRate uint32) Music {
	return &NilMusic{}
}

func (*NilAudioManager) Terminate() {

}