package audio

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	al "github.com/timshannon/go-openal/openal"
	"strconv"
	"time"
)

type OpenALSound struct {
	Name     string
	Duration time.Duration

	buffer  al.Buffer
	source  al.Source
	playing bool
	volume  float32
}

func (this *OpenALSound) Play(loop bool) {
	this.source.SetLooping(loop)
	this.source.Play()
	this.playing = true
}
func (this *OpenALSound) Pause() {
	this.source.Pause()
	this.playing = false
}
func (this *OpenALSound) Resume() {
	if !this.playing {
		this.Play(this.source.GetLooping())
	}
}
func (this *OpenALSound) Stop() {
	this.source.Stop()
	this.playing = false
}
func (this *OpenALSound) Terminate() {
	this.buffer.Delete()
	this.source.Delete()
	this.playing = false
	audioMgr := gohome.Framew.GetAudioManager().(*OpenALAudioManager)
	audioMgr.removeSoundFromSlice(this)
}
func (this *OpenALSound) IsPlaying() bool {
	return this.playing
}
func (this *OpenALSound) GetPlayingDuration() time.Duration {
	microSecOffset := int64(this.source.GetOffsetSeconds() * 1000000.0)
	dur, _ := time.ParseDuration(strconv.Itoa(int(microSecOffset)) + "µs")
	return dur
}
func (this *OpenALSound) GetDuration() time.Duration {
	return this.Duration
}

func (this *OpenALSound) SetVolume(vol float32) {
	audio := gohome.Framew.GetAudioManager().(*OpenALAudioManager)
	this.setVolumeHard(vol * audio.volume)
	this.volume = vol
}

func (this *OpenALSound) setVolumeHard(vol float32) {
	this.source.Setf(al.AlGain, vol)
}

func (this *OpenALSound) GetVolume() float32 {
	return this.volume
}

type OpenALMusic struct {
	Name     string
	Duration time.Duration

	buffer  al.Buffer
	source  al.Source
	playing bool
	volume  float32
}

func (this *OpenALMusic) Play(loop bool) {
	this.source.SetLooping(loop)
	this.source.Play()
	this.playing = true
}
func (this *OpenALMusic) Pause() {
	this.source.Pause()
	this.playing = false
}
func (this *OpenALMusic) Resume() {
	if !this.playing {
		this.Play(this.source.GetLooping())
	}
}
func (this *OpenALMusic) Stop() {
	this.source.Stop()
	this.playing = false
}
func (this *OpenALMusic) Terminate() {
	this.source.Delete()
	this.buffer.Delete()
	this.playing = false
	audioMgr := gohome.Framew.GetAudioManager().(*OpenALAudioManager)
	audioMgr.removeMusicFromSlice(this)
}
func (this *OpenALMusic) IsPlaying() bool {
	return this.playing
}
func (this *OpenALMusic) GetPlayingDuration() time.Duration {
	microSecOffset := int64(this.source.GetOffsetSeconds() * 1000000.0)
	dur, _ := time.ParseDuration(strconv.Itoa(int(microSecOffset)) + "µs")
	return dur
}
func (this *OpenALMusic) GetDuration() time.Duration {
	return this.Duration
}

func (this *OpenALMusic) SetVolume(vol float32) {
	audio := gohome.Framew.GetAudioManager().(*OpenALAudioManager)
	this.setVolumeHard(vol * audio.volume)
	this.volume = vol
}

func (this *OpenALMusic) setVolumeHard(vol float32) {
	this.source.Setf(al.AlGain, vol)
}

func (this *OpenALMusic) GetVolume() float32 {
	return this.volume
}

type OpenALAudioManager struct {
	device  *al.Device
	context *al.Context
	sounds  []*OpenALSound
	musics  []*OpenALMusic

	volume float32
}

func (this *OpenALAudioManager) Init() {
	this.device = al.OpenDevice("")
	if err := this.device.Err(); err != nil {
		gohome.ErrorMgr.Error("Audio", "OpenAL", "Couldn't open device: "+err.Error())
		return
	}
	this.context = this.device.CreateContext()
	if err := this.device.Err(); err != nil {
		gohome.ErrorMgr.Error("Audio", "OpenAL", "Couldn't create context: "+err.Error())
		this.device.CloseDevice()
		return
	}
	this.context.Activate()
	if err := this.device.Err(); err != nil {
		gohome.ErrorMgr.Error("Audio", "OpenAL", "Couldn't activate context: "+err.Error())
		this.context.Destroy()
		this.device.CloseDevice()
	}

	gohome.UpdateMgr.AddObject(this)

	this.volume = 1.0
}
func (this *OpenALAudioManager) CreateSound(name string, samples []byte, format uint8, sampleRate uint32) gohome.Sound {
	sound := &OpenALSound{}
	sound.Name = name
	sound.buffer = al.NewBuffer()
	if err := al.Err(); err != nil {
		gohome.ErrorMgr.Error("Sound", name, "Couldn't create buffer: "+err.Error())
		return nil
	}
	sound.source = al.NewSource()
	if err := al.Err(); err != nil {
		gohome.ErrorMgr.Error("Sound", name, "Couldn't create source: "+err.Error())
		sound.buffer.Delete()
		return nil
	}

	sound.buffer.SetData(getOpenALFormat(format), samples, int32(sampleRate))
	sound.source.SetBuffer(sound.buffer)

	var microSeconds int64
	switch format {
	case gohome.AUDIO_FORMAT_MONO8:
		microSeconds = int64((float64(len(samples)) * 8.0 / 8.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	case gohome.AUDIO_FORMAT_MONO16:
		microSeconds = int64((float64(len(samples)) * 8.0 / 16.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	case gohome.AUDIO_FORMAT_STEREO8:
		microSeconds = int64((float64(len(samples)) * 8.0 / 8.0 / 2.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	case gohome.AUDIO_FORMAT_STEREO16:
		microSeconds = int64((float64(len(samples)) * 8.0 / 16.0 / 2.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	}

	sound.Duration, _ = time.ParseDuration(strconv.Itoa(int(microSeconds)) + "µs")

	this.sounds = append(this.sounds, sound)

	sound.SetVolume(1.0)

	return sound
}
func (this *OpenALAudioManager) CreateMusic(name string, samples []byte, format uint8, sampleRate uint32) gohome.Music {
	music := &OpenALMusic{}
	music.Name = name
	music.buffer = al.NewBuffer()
	if err := al.Err(); err != nil {
		gohome.ErrorMgr.Error("Music", name, "Couldn't create buffer: "+err.Error())
		return nil
	}
	music.source = al.NewSource()
	if err := al.Err(); err != nil {
		gohome.ErrorMgr.Error("Music", name, "Couldn't create source: "+err.Error())
		music.buffer.Delete()
		return nil
	}

	music.buffer.SetData(getOpenALFormat(format), samples, int32(sampleRate))
	music.source.SetBuffer(music.buffer)

	var microSeconds int64
	switch format {
	case gohome.AUDIO_FORMAT_MONO8:
		microSeconds = int64((float64(len(samples)) * 8.0 / 8.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	case gohome.AUDIO_FORMAT_MONO16:
		microSeconds = int64((float64(len(samples)) * 8.0 / 16.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	case gohome.AUDIO_FORMAT_STEREO8:
		microSeconds = int64((float64(len(samples)) * 8.0 / 8.0 / 2.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	case gohome.AUDIO_FORMAT_STEREO16:
		microSeconds = int64((float64(len(samples)) * 8.0 / 16.0 / 2.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	}

	music.Duration, _ = time.ParseDuration(strconv.Itoa(int(microSeconds)) + "µs")

	this.musics = append(this.musics, music)

	music.SetVolume(1.0)

	return music
}
func (this *OpenALAudioManager) Terminate() {
	this.context.Destroy()
	this.device.CloseDevice()
}
func (this *OpenALAudioManager) Update(delta_time float32) {
	for i := 0; i < len(this.musics); i++ {
		if this.musics[i].IsPlaying() && !this.musics[i].source.GetLooping() {
			plPos := this.musics[i].GetPlayingDuration()
			if plPos >= this.musics[i].Duration || plPos == time.Second*0 {
				this.musics[i].playing = false
			} else {
				this.musics[i].playing = true
			}
		}
	}
	for i := 0; i < len(this.sounds); i++ {
		if this.sounds[i].IsPlaying() && !this.sounds[i].source.GetLooping() {
			plPos := this.sounds[i].GetPlayingDuration()
			if plPos >= this.sounds[i].Duration || plPos == time.Second*0 {
				this.sounds[i].playing = false
			} else {
				this.sounds[i].playing = true
			}
		}
	}
}
func (this *OpenALAudioManager) removeMusicFromSlice(music *OpenALMusic) {
	if len(this.musics) == 1 {
		this.musics = this.musics[:0]
	} else if len(this.musics) == 0 {
		return
	} else {
		var index, i uint32
		for i = 0; i < uint32(len(this.musics)); i++ {
			if this.musics[i] == music {
				index = i
				break
			}
		}
		this.musics = append(this.musics[:index], this.musics[index+1:]...)
	}
}
func (this *OpenALAudioManager) removeSoundFromSlice(sound *OpenALSound) {
	if len(this.sounds) == 1 {
		this.sounds = this.sounds[:0]
	} else if len(this.sounds) == 0 {
		return
	} else {
		var index, i uint32
		for i = 0; i < uint32(len(this.sounds)); i++ {
			if this.sounds[i] == sound {
				index = i
				break
			}
		}
		this.sounds = append(this.sounds[:index], this.sounds[index+1:]...)
	}
}
func (this *OpenALAudioManager) SetVolume(vol float32) {
	this.volume = vol

	for _, s := range this.sounds {
		s.setVolumeHard(s.volume * vol)
	}
	for _, m := range this.musics {
		m.setVolumeHard(m.volume * vol)
	}
}

func (this *OpenALAudioManager) GetVolume() float32 {
	return this.volume
}

func getOpenALFormat(gohomeformat uint8) al.Format {
	switch gohomeformat {
	case gohome.AUDIO_FORMAT_MONO8:
		return al.FormatMono8
	case gohome.AUDIO_FORMAT_MONO16:
		return al.FormatMono16
	case gohome.AUDIO_FORMAT_STEREO8:
		return al.FormatStereo8
	case gohome.AUDIO_FORMAT_STEREO16:
		return al.FormatStereo16
	}

	return 0
}
