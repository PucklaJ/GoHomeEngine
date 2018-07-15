package loader

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/go-wav"
)

func LoadWAVFile(fileName string) (*wav.WavData, error) {
	/*wavInfo,err := os.Stat(fileName)
	if err != nil {
		return nil,err
	}

	file,err := gohome.Framew.OpenFile(fileName)
	if err != nil {
		return nil,err
	}

	reader,err := wav.NewReader(file,wavInfo.Size())
	if err != nil {
		return reader,err
	}

	return reader,nil*/

	file,err := gohome.Framew.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	data,err := wav.ReadWavData(file)
	if err != nil {
		return nil,err
	}

	return &data,err
}

func convert24BitTo16Bit(samples []byte,sampleCount uint32) []byte {
	s24 := samples
	newSamples := make([]byte,sampleCount*2)
	var index uint32 = 0
	for a:=0; uint32(a)<sampleCount*3; a+=3 {
		newSamples[index+0] = s24[a+1+0]
		newSamples[index+1] = s24[a+1+1]
		index +=2
	}
	return newSamples
}

func ReadAllSamples(data *wav.WavData) ([]byte,error) {

	samples := data.Data

	if data.BitsPerSample == 24 {
		samples = convert24BitTo16Bit(samples,uint32(len(samples)*8/24))
	}

	return samples,nil
}

func GetAudioFormat(data *wav.WavData) uint8 {
	numChannels := data.NumChannels
	bitsPerSample := data.BitsPerSample

	switch numChannels {
	case 1:
		switch bitsPerSample {
		case 8:
			return gohome.AUDIO_FORMAT_MONO8
		case 16,24:
			return gohome.AUDIO_FORMAT_MONO16
		}
	case 2:
		switch bitsPerSample {
		case 8:
			return gohome.AUDIO_FORMAT_STEREO8
		case 16,24:
			return gohome.AUDIO_FORMAT_STEREO16
		}
	}

	return gohome.AUDIO_FORMAT_UNKNOWN
}