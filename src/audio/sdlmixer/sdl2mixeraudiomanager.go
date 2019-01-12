package sdlmixer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/go-sdl2/mix"
	"github.com/PucklaMotzer09/go-sdl2/sdl"
	"strconv"
	"time"
)

type MixerAudioManager struct {
	gohome.NilAudioManager
	volume float32

	sounds []*MixerSound
	musics []*MixerMusic
}

func (this *MixerAudioManager) Init() {
	audioManager = this
	if err := mix.Init(mix.INIT_MP3 | mix.INIT_OGG); err != nil {
		gohome.ErrorMgr.Error("AudioManager", "SDL_mixer", "Failed to initialise: "+err.Error())
		return
	}

	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE); err != nil {
		gohome.ErrorMgr.Error("AudioManager", "SDL_mixer", "Failed to open audio: "+err.Error())
		return
	}

	this.volume = 1.0
}

func (this *MixerAudioManager) CreateSound(name string, samples []byte, format uint8, sampleRate uint32) gohome.Sound {
	return &MixerSound{}
}

func (this *MixerAudioManager) CreateMusic(name string, samples []byte, format uint8, sampleRate uint32) gohome.Music {
	return &MixerMusic{}
}

func (this *MixerAudioManager) LoadSound(name, path string) gohome.Sound {
	var file *sdl.RWops
	for _, filepath := range gohome.MUSIC_SOUND_PATHS {
		file = sdl.RWFromFile(filepath+path, "rb")
		if file != nil {
			break
		} else {
			file = sdl.RWFromFile(filepath+gohome.GetFileFromPath(path), "rb")
			if file != nil {
				break
			}
		}
	}

	if file == nil {
		err := sdl.GetError()
		if err != nil {
			gohome.ErrorMgr.Error("Sound", name, "Failed to open file: "+err.Error())
		} else {
			gohome.ErrorMgr.Error("Sound", name, "Failed to open file")
		}
		return nil
	}

	var sound MixerSound
	sound.channel = -2
	sound.volume = 1.0

	var err error
	sound.chunk, err = mix.LoadWAVRW(file, true)
	if err != nil {
		gohome.ErrorMgr.Error("Sound", name, "Failed to load file: "+err.Error())
		return nil
	}

	this.sounds = append(this.sounds, &sound)

	return &sound
}

func (this *MixerAudioManager) LoadMusic(name, path string) gohome.Music {
	var file *sdl.RWops
	for _, filepath := range gohome.MUSIC_SOUND_PATHS {
		file = sdl.RWFromFile(filepath+path, "rb")
		if file != nil {
			break
		} else {
			file = sdl.RWFromFile(filepath+gohome.GetFileFromPath(path), "rb")
			if file != nil {
				break
			}
		}
	}

	if file == nil {
		err := sdl.GetError()
		if err != nil {
			gohome.ErrorMgr.Error("Music", name, "Failed to open file: "+err.Error())
		} else {
			gohome.ErrorMgr.Error("Music", name, "Failed to open file")
		}
		return nil
	}

	var music MixerMusic
	music.volume = 1.0

	var err error
	music.music, err = mix.LoadMUSRW(file, 1)
	if err != nil {
		gohome.ErrorMgr.Error("Music", name, "Failed to load file: "+err.Error())
		return nil
	}

	this.musics = append(this.musics, &music)

	return &music
}

func (this *MixerAudioManager) SetVolume(vol float32) {
	this.volume = vol
	for _, s := range this.sounds {
		s.actuallySetVolume()
	}
	for _, m := range this.musics {
		m.actuallySetVolume()
	}
}

func (this *MixerAudioManager) GetVolume() float32 {
	return this.volume
}

func (this *MixerAudioManager) Terminate() {
	for _, s := range this.sounds {
		s.chunk.Free()
		s.freed = true
	}
	for _, m := range this.musics {
		m.music.Free()
		m.freed = true
	}
	mix.CloseAudio()
	mix.Quit()
}

type MixerSound struct {
	chunk   *mix.Chunk
	channel int
	paused  bool
	volume  float32
	freed   bool
	looping bool
}

func (this *MixerSound) Play(loop bool) {
	this.looping = loop
	var loops int
	if loop {
		loops = -1
	} else {
		loops = 0
	}
	this.channel, _ = this.chunk.Play(-1, loops)
	this.paused = false
}
func (this *MixerSound) Pause() {
	if this.channel != -2 {
		this.paused = true
		mix.Pause(this.channel)
	}
}
func (this *MixerSound) Resume() {
	if !this.paused {
		this.Play(this.looping)
	}
	if this.channel != -2 {
		this.paused = false
		mix.Resume(this.channel)
	}
}
func (this *MixerSound) Stop() {
	if this.channel != -2 {
		mix.HaltChannel(this.channel)
		this.channel = -2
	}
}
func (this *MixerSound) Terminate() {
	if this.freed {
		return
	}
	this.chunk.Free()
	this.channel = -2

	for i := 0; i < len(audioManager.sounds); i++ {
		if audioManager.sounds[i] == this {
			audioManager.sounds = append(audioManager.sounds[:i], audioManager.sounds[i+1:]...)
		}
	}

	this.freed = true
}
func (this *MixerSound) IsPlaying() bool {
	if this.channel == -2 || this.paused {
		return false
	}
	return mix.Playing(this.channel) != 0
}
func (this *MixerSound) GetPlayingDuration() time.Duration {
	return time.Millisecond * 0
}
func (this *MixerSound) GetDuration() time.Duration {
	ms := this.chunk.LengthInMs()
	dur, _ := time.ParseDuration(strconv.FormatInt(int64(ms), 10) + "ms")
	return dur
}
func (this *MixerSound) SetVolume(vol float32) {
	this.volume = vol
	this.actuallySetVolume()
}

func (this *MixerSound) actuallySetVolume() {
	vol := int(this.volume * audioManager.volume * mix.MAX_VOLUME)
	this.chunk.Volume(vol)
}

func (this *MixerSound) GetVolume() float32 {
	return float32(this.chunk.Volume(-1)) / mix.MAX_VOLUME
}

type MixerMusic struct {
	music  *mix.Music
	paused bool
	volume float32
	freed  bool
}

func (this *MixerMusic) Play(loop bool) {
	var loops int
	if loop {
		loops = -1
	} else {
		loops = 1
	}
	this.music.Play(loops)
	this.actuallySetVolume()
	this.paused = false
}
func (this *MixerMusic) Pause() {
	mix.PauseMusic()
	this.paused = true
}
func (this *MixerMusic) Resume() {
	mix.ResumeMusic()
	this.paused = false
}
func (this *MixerMusic) Stop() {
	mix.HaltMusic()
}
func (this *MixerMusic) Terminate() {
	if this.freed {
		return
	}
	this.music.Free()
	for i := 0; i < len(audioManager.musics); i++ {
		if audioManager.musics[i] == this {
			audioManager.musics = append(audioManager.musics[:i], audioManager.musics[i+1:]...)
		}
	}
	this.freed = true
}
func (this *MixerMusic) IsPlaying() bool {
	if this.paused {
		return false
	}
	return mix.PlayingMusic()
}
func (this *MixerMusic) GetPlayingDuration() time.Duration {
	return time.Millisecond * 0
}
func (this *MixerMusic) GetDuration() time.Duration {
	return time.Millisecond * 0
}
func (this *MixerMusic) SetVolume(vol float32) {
	this.volume = vol
	this.actuallySetVolume()
}

func (this *MixerMusic) actuallySetVolume() {
	vol := int(this.volume * audioManager.volume * mix.MAX_VOLUME)
	mix.VolumeMusic(vol)
}

func (this *MixerMusic) GetVolume() float32 {
	return this.volume
}

var audioManager *MixerAudioManager
