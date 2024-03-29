package loader

import (
	"errors"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/PucklaJ/go-wav"
)

func LoadWAVFile(fileName string) (*wav.WavData, error) {

	var file gohome.File
	var err error
	for _, path := range gohome.MUSIC_SOUND_PATHS {
		file, err = gohome.Framew.OpenFile(path + fileName)
		if err == nil {
			break
		} else {
			file, err = gohome.Framew.OpenFile(path + gohome.GetFileFromPath(fileName))
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return nil, err
	}
	files, ok := file.(gohome.FileSeeker)
	if ok {
		data, err := wav.ReadWavData(files)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}

	return nil, errors.New("Cannot convert file to io.ReadSeeker")
}

func convert24BitTo16Bit(samples []byte, sampleCount int) []byte {
	s24 := samples
	newSamples := make([]byte, sampleCount*2)
	var index = 0
	for a := 0; a < sampleCount*3; a += 3 {
		newSamples[index+0] = s24[a+1+0]
		newSamples[index+1] = s24[a+1+1]
		index += 2
	}
	return newSamples
}

func ReadAllSamples(data *wav.WavData) ([]byte, error) {

	samples := data.Data

	if data.BitsPerSample == 24 {
		samples = convert24BitTo16Bit(samples, len(samples)*8/24)
	}

	return samples, nil
}

func GetAudioFormat(data *wav.WavData) uint8 {
	numChannels := data.NumChannels
	bitsPerSample := data.BitsPerSample

	switch numChannels {
	case 1:
		switch bitsPerSample {
		case 8:
			return gohome.AUDIO_FORMAT_MONO8
		case 16, 24:
			return gohome.AUDIO_FORMAT_MONO16
		}
	case 2:
		switch bitsPerSample {
		case 8:
			return gohome.AUDIO_FORMAT_STEREO8
		case 16, 24:
			return gohome.AUDIO_FORMAT_STEREO16
		}
	}

	return gohome.AUDIO_FORMAT_UNKNOWN
}
