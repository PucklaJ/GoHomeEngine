package openal

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	loadmp3 "github.com/PucklaMotzer09/GoHomeEngine/src/loaders/mp3"
	loadwav "github.com/PucklaMotzer09/GoHomeEngine/src/loaders/wav"
	al "github.com/PucklaMotzer09/go-openal/openal"
	"github.com/hajimehoshi/go-mp3"
	"io"
	"strconv"
	"time"
)

type OpenALSound struct {
	Name     string
	Duration time.Duration

	buffer     al.Buffer
	source     al.Source
	playing    bool
	volume     float32
	terminated bool
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
	if this.terminated {
		return
	}
	this.buffer.Delete()
	this.source.Delete()
	this.playing = false
	audioManager.removeSoundFromSlice(this)
	this.terminated = true
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
	this.setVolumeHard(vol * audioManager.volume)
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

	buffers    []al.Buffer
	source     al.Source
	playing    bool
	volume     float32
	terminated bool
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
	if this.terminated {
		return
	}
	this.source.Delete()
	for _, b := range this.buffers {
		b.Delete()
	}
	this.playing = false
	audioManager.removeMusicFromSlice(this)
	this.terminated = true
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
	this.setVolumeHard(vol * audioManager.volume)
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
	failed bool
}

func (this *OpenALAudioManager) Init() {
	audioManager = this
	defer func() {
		this.failed = !this.failed
	}()
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
		return
	}

	gohome.UpdateMgr.AddObject(this)

	this.volume = 1.0
	this.failed = true
}
func (this *OpenALAudioManager) CreateSound(name string, samples []byte, format uint8, sampleRate uint32) gohome.Sound {
	if this.failed {
		gohome.ErrorMgr.Error("Sound", name, "Couldn't create because the initialisation failed!")
		return &gohome.NilSound{}
	}
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
	if this.failed {
		gohome.ErrorMgr.Error("Music", name, "Couldn't create because the initialisation failed!")
		return &gohome.NilMusic{}
	}
	music := &OpenALMusic{}
	music.Name = name
	music.buffers = append(music.buffers, al.NewBuffer())
	if err := al.Err(); err != nil {
		gohome.ErrorMgr.Error("Music", name, "Couldn't create buffer: "+err.Error())
		return nil
	}
	music.source = al.NewSource()
	if err := al.Err(); err != nil {
		gohome.ErrorMgr.Error("Music", name, "Couldn't create source: "+err.Error())
		music.buffers[0].Delete()
		return nil
	}

	music.buffers[0].SetData(getOpenALFormat(format), samples, int32(sampleRate))
	music.source.SetBuffer(music.buffers[0])

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

const SAMPLE_SIZE_MP3 = 2 * 2
const SAMPLE_RATE_MP3 = 44100
const BUFFER_SIZE_MP3 = SAMPLE_RATE_MP3 / 10 * SAMPLE_SIZE_MP3

func readSamples(decoder *mp3.Decoder) ([BUFFER_SIZE_MP3]byte, int, error) {
	var samples [BUFFER_SIZE_MP3]byte
	var num int
	var err error
	var n int

	for num < BUFFER_SIZE_MP3 && err == nil {
		n, err = decoder.Read(samples[num:])
		num += n
	}

	return samples, num, err
}

func (this *OpenALAudioManager) createMusicMP3(name string, decoder *mp3.Decoder) gohome.Music {
	format := loadmp3.GetAudioFormat()
	sampleRate := decoder.SampleRate()

	music := &OpenALMusic{}
	music.Name = name
	music.buffers = append(music.buffers, al.NewBuffer())
	if err := al.Err(); err != nil {
		gohome.ErrorMgr.Error("Music", name, "Couldn't create buffer: "+err.Error())
		return nil
	}
	music.source = al.NewSource()
	if err := al.Err(); err != nil {
		gohome.ErrorMgr.Error("Music", name, "Couldn't create source: "+err.Error())
		music.buffers[0].Delete()
		return nil
	}

	samples, num, err := readSamples(decoder)
	if err != nil && err != io.EOF {
		gohome.ErrorMgr.Error("Music", name, err.Error())
		return nil
	}

	music.buffers[0].SetData(getOpenALFormat(format), samples[:num], int32(sampleRate))
	music.source.QueueBuffer(music.buffers[0])
	go func() {
		var err error
		var samples [BUFFER_SIZE_MP3]byte
		var num int
		var numBuffers int
		for err == nil && !music.terminated {
			buffer := al.NewBuffer()
			samples, num, err = readSamples(decoder)
			buffer.SetData(getOpenALFormat(format), samples[:num], int32(decoder.SampleRate()))
			music.source.QueueBuffer(buffer)
			music.buffers = append(music.buffers, buffer)
			numBuffers++
		}
		decoder.Close()
		if err != nil && err != io.EOF {
			gohome.ErrorMgr.Error("Music", name, err.Error())
		} else if err != nil && err == io.EOF {
			gohome.ErrorMgr.Log("Music", name, "Read all samples! B: "+strconv.FormatInt(int64(numBuffers), 10))
		}
	}()

	length := decoder.Length()

	var microSeconds int64
	switch format {
	case gohome.AUDIO_FORMAT_MONO8:
		microSeconds = int64((float64(length) * 8.0 / 8.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	case gohome.AUDIO_FORMAT_MONO16:
		microSeconds = int64((float64(length) * 8.0 / 16.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	case gohome.AUDIO_FORMAT_STEREO8:
		microSeconds = int64((float64(length) * 8.0 / 8.0 / 2.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	case gohome.AUDIO_FORMAT_STEREO16:
		microSeconds = int64((float64(length) * 8.0 / 16.0 / 2.0 * (1.0 / float64(sampleRate))) * 1000000.0)
	}

	music.Duration, _ = time.ParseDuration(strconv.Itoa(int(microSeconds)) + "µs")

	this.musics = append(this.musics, music)

	music.SetVolume(1.0)

	return music
}

func (this *OpenALAudioManager) Terminate() {
	for _, s := range this.sounds {
		s.Terminate()
	}
	for _, m := range this.musics {
		m.Terminate()
	}
	this.sounds = this.sounds[:0]
	this.musics = this.musics[:0]

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

func (this *OpenALAudioManager) LoadSound(name, path string) gohome.Sound {
	wavReader, err := loadwav.LoadWAVFile(path)
	if err != nil {
		gohome.ErrorMgr.MessageError(gohome.ERROR_LEVEL_ERROR, "Sound", name, err)
		return nil
	}
	format := loadwav.GetAudioFormat(wavReader)
	if format == gohome.AUDIO_FORMAT_UNKNOWN {
		gohome.ErrorMgr.Error("Sound", name, "The audio format is unknow: C: "+strconv.Itoa(int(wavReader.NumChannels))+" B: "+strconv.Itoa(int(wavReader.BitsPerSample)))
		return nil
	}
	samples, err := loadwav.ReadAllSamples(wavReader)
	if err != nil {
		gohome.ErrorMgr.MessageError(gohome.ERROR_LEVEL_ERROR, "Sound", name, err)
		return nil
	}
	sampleRate := wavReader.SampleRate

	sound := this.CreateSound(name, samples, format, sampleRate)

	return sound
}

func (this *OpenALAudioManager) LoadMusic(name, path string) gohome.Music {
	decoder, err := loadmp3.LoadMP3File(path)
	if err != nil {
		gohome.ErrorMgr.MessageError(gohome.ERROR_LEVEL_ERROR, "Music", name, err)
		return nil
	}

	music := this.createMusicMP3(name, decoder)
	return music
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

var audioManager *OpenALAudioManager
