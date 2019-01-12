package jsaudio

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/gopherjs/gopherjs/js"
	"time"
)

type JSAudioManager struct {
	sounds []*JSAudio
	music  []*JSAudio

	volume float32
}

func (this *JSAudioManager) Init() {
	this.volume = 1.0
}

func (this *JSAudioManager) CreateSound(name string, samples []byte, format uint8, sampleRate uint32) gohome.Sound {
	return &JSAudio{}
}

func (this *JSAudioManager) CreateMusic(name string, samples []byte, format uint8, sampleRate uint32) gohome.Music {
	return &JSAudio{}
}

func (this *JSAudioManager) loadAudio(name, path string) *JSAudio {
	var rv JSAudio
	rv.audio = js.Global.Get("Audio").New(path)
	rv.audio.Set("onerror", func() {
		gohome.ErrorMgr.Error("Audio", name, "Error: "+rv.audio.Get("error").Get("message").String())
	})
	rv.volume = 1.0
	return &rv
}

func (this *JSAudioManager) LoadSound(name, path string) gohome.Sound {
	audio := this.loadAudio(name, path)
	this.sounds = append(this.sounds, audio)
	return audio
}

func (this *JSAudioManager) LoadMusic(name, path string) gohome.Music {
	audio := this.loadAudio(name, path)
	this.music = append(this.music, audio)
	return audio
}

func (this *JSAudioManager) SetVolume(vol float32) {
	this.volume = vol
	for _, s := range this.sounds {
		s.actuallySetVolume()
	}
	for _, m := range this.music {
		m.actuallySetVolume()
	}
}

func (this *JSAudioManager) GetVolume() float32 {
	return this.volume
}

func (this *JSAudioManager) Terminate() {

}

type JSAudio struct {
	audio  *js.Object
	volume float32
}

func (this *JSAudio) Play(loop bool) {
	this.audio.Set("loop", loop)
	this.audio.Call("play")
}
func (this *JSAudio) Pause() {
	this.audio.Call("pause")
}
func (this *JSAudio) Resume() {
	this.audio.Call("play")
}
func (this *JSAudio) Stop() {
	this.audio.Call("pause")
	this.audio.Set("currentTime", 0)
}
func (this *JSAudio) Terminate() {

}
func (this *JSAudio) IsPlaying() bool {
	return !this.audio.Get("paused").Bool() && !this.audio.Get("ended").Bool()
}
func (this *JSAudio) GetPlayingDuration() time.Duration {
	dur, _ := time.ParseDuration(this.audio.Get("currentTime").String() + "s")
	return dur
}
func (this *JSAudio) GetDuration() time.Duration {
	dur, _ := time.ParseDuration(this.audio.Get("duration").String() + "s")
	return dur
}
func (this *JSAudio) actuallySetVolume() {
	this.audio.Set("volume", this.volume*audioManager.volume)
}
func (this *JSAudio) SetVolume(vol float32) {
	this.volume = vol
	this.actuallySetVolume()
}
func (this *JSAudio) GetVolume() float32 {
	return this.volume
}

var audioManager JSAudioManager
