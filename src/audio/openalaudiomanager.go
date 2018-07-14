package audio

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/phf/go-openal/alc"
	"strconv"
	"github.com/phf/go-openal/al"
)

type OpenALSound struct {
	Name string

	buffer al.Buffer
	source al.Source
	playing bool
}

func (this *OpenALSound) Play() {
	al.PlaySources([]al.Source{this.source})
	this.playing = true
}
func (this *OpenALSound) Pause() {
	al.PauseSources([]al.Source{this.source})
	this.playing = false
}
func (this *OpenALSound) Resume() {
	this.Play()
}
func (this *OpenALSound) Stop() {
	al.StopSources([]al.Source{this.source})
	this.playing = false
}
func (this *OpenALSound) Terminate() {
	al.DeleteBuffer(this.buffer)
	al.DeleteSource(this.source)
}
func (this *OpenALSound) IsPlaying() bool {
	return this.playing
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
	device *alc.Device
	context *alc.Context
}

func (this *OpenALAudioManager) Init() {
	this.device = alc.OpenDevice("")
	if err := this.device.GetError(); err != alc.NoError {
		gohome.ErrorMgr.Error("Audio", "OpenAL", "Couldn't open device: "+strconv.Itoa(int(err)))
		return
	}
	this.context = this.device.CreateContext()
	if err := this.device.GetError(); err != alc.NoError {
		gohome.ErrorMgr.Error("Audio","OpenAL","Couldn't create context: " + strconv.Itoa(int(err)))
		this.device.CloseDevice()
		return
	}
	this.context.Activate()
	if err := this.device.GetError(); err != alc.NoError {
		gohome.ErrorMgr.Error("Audio","OpenAL", "Couldn't activate context: " + strconv.Itoa(int(err)))
		this.context.Destroy()
		this.device.CloseDevice()
	}

}
func (this *OpenALAudioManager) CreateSound(name string, samples []byte, format uint8, sampleRate uint32) gohome.Sound {
	sound := &OpenALSound{}
	sound.Name = name
	sound.buffer = al.NewBuffer()
	if err := al.GetError(); err != al.NoError {
		gohome.ErrorMgr.Error("Sound",name,"Couldn't create buffer: " + strconv.Itoa(int(err)))
		return nil
	}
	sound.source = al.NewSource()
	if err := al.GetError(); err != al.NoError {
		gohome.ErrorMgr.Error("Sound",name,"Couldn't create source: " + strconv.Itoa(int(err)))
		al.DeleteBuffer(sound.buffer)
		return nil
	}

	sound.buffer.SetData(getOpenALFormat(format),samples,int32(sampleRate))
	sound.source.SetBuffer(sound.buffer)

	return sound
}
func (this *OpenALAudioManager) CreateMusic(name string, samples []byte,format uint8, sampleRate uint32) gohome.Music {
	return &OpenALMusic{}
}
func (this *OpenALAudioManager) Terminate() {
	this.context.Destroy()
	this.device.CloseDevice()
}

func getOpenALFormat(gohomeformat uint8) int32 {
	switch gohomeformat {
	case gohome.AUDIO_FORMAT_MONO8: return al.FormatMono8
	case gohome.AUDIO_FORMAT_MONO16: return al.FormatMono16
	case gohome.AUDIO_FORMAT_STEREO8: return al.FormatStereo8
	case gohome.AUDIO_FORMAT_STEREO16: return al.FormatStereo16
	}

	return 0
}
