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
)

type AudioManager interface {
	Init()
	CreateSound(name string,fileName string) Sound
	CreateMusic(name string, samples []byte,format uint8, sampleRate uint32) Music
	Terminate()
}
