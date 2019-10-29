package gohome

import "time"

// This interface represents short chunks of audio
type Sound interface {
	// Plays the audio chunk
	// loops sets wether it should loop
	Play(loop bool)
	// Pauses this Sound
	Pause()
	// Resumes it if it has been paused
	Resume()
	// Stops it so that it will start from the
	// beginning the next time you play it
	Stop()
	// Cleans up all data related to this resource
	Terminate()
	// Returns wether it is currently playing/ not paused
	IsPlaying() bool
	// Returns how far the Sound has progressed in playing
	GetPlayingDuration() time.Duration
	// Returns the total duration of this sound
	GetDuration() time.Duration
	// Sets the volume of this sound (0.0 - 1.0)
	SetVolume(vol float32)
	// Returns the currently set volume
	GetVolume() float32
}

// This interface represents a longer chunk of audio
type Music interface {
	// Plays the audio chunk
	// loops sets wether it should loop
	Play(loop bool)
	// Pauses this Music
	Pause()
	// Resumes it if it has been paused
	Resume()
	// Stops it so that it will start from the
	// beginning the next time you play it
	Stop()
	// Cleans up all data related to this resource
	Terminate()
	// Returns wether it is currently playing/ not paused
	IsPlaying() bool
	// Returns how far the Music has progressed in playing
	GetPlayingDuration() time.Duration
	// Returns the total duration of this Music
	GetDuration() time.Duration
	// Sets the volume of this sound (0.0 - 1.0)
	SetVolume(vol float32)
	// Returns the currently set volume
	GetVolume() float32
}

const (
	AUDIO_FORMAT_STEREO8  uint8 = iota
	AUDIO_FORMAT_STEREO16 uint8 = iota
	AUDIO_FORMAT_MONO8    uint8 = iota
	AUDIO_FORMAT_MONO16   uint8 = iota
	AUDIO_FORMAT_UNKNOWN  uint8 = iota
)

// This interface handles everything audio related
type AudioManager interface {
	// Initialises the AudioManager
	Init()
	// Creates a new Sound object from samples and the given format and sampleRate
	CreateSound(name string, samples []byte, format uint8, sampleRate int) Sound
	// Creates a new Music object from samples and the given format and sampleRate
	CreateMusic(name string, samples []byte, format uint8, sampleRate int) Music
	// Loads a Sound object from a file (.wav)
	LoadSound(name, path string) Sound
	// Loads a new Music object from a file (.mp3)
	LoadMusic(name, path string) Music
	// Sets the master volume of the game
	SetVolume(vol float32)
	// Gets the master volume of the game
	GetVolume() float32
	// Cleans up all resources
	Terminate()
}

// The AudioManager that should be used for everything
var AudioMgr AudioManager

// An implementation of the Sound interface that does nothing
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

// An implementation of the Music interface that does nothing
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

// An implementation of the AudioManager interface that does nothing
type NilAudioManager struct {
}

func (*NilAudioManager) Init() {

}
func (*NilAudioManager) CreateSound(name string, samples []byte, format uint8, sampleRate int) Sound {
	return &NilSound{}
}
func (*NilAudioManager) CreateMusic(name string, samples []byte, format uint8, sampleRate int) Music {
	return &NilMusic{}
}
func (*NilAudioManager) SetVolume(vol float32) {

}
func (*NilAudioManager) GetVolume() float32 {
	return 1.0
}
func (*NilAudioManager) Terminate() {

}

func (*NilAudioManager) LoadSound(name, path string) Sound {
	return &NilSound{}
}

func (*NilAudioManager) LoadMusic(name, path string) Music {
	return &NilMusic{}
}
