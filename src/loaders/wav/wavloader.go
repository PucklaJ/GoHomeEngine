package loader

import (
	"github.com/cryptix/wav"
	"os"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"io"
)

func LoadWAVFile(fileName string) (*wav.Reader, error) {
	wavInfo,err := os.Stat(fileName)
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

	return reader,nil
}

func StreamWAVFile(reader *wav.Reader) ([]byte,error) {
	return reader.ReadRawSample()
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

func ReadAllSamples(reader *wav.Reader) ([]byte,error) {
	rawReader,err := reader.GetDumbReader()
	if err != nil {
		return nil,err
	}

	numSamples := reader.GetSampleCount()
	bitsPerSample := reader.GetBitsPerSample()
	bytesPerSample := bitsPerSample/8
	byteCount := uint32(bytesPerSample)*numSamples

	data := make([]byte,byteCount)

	readBytes,err := rawReader.Read(data)
	if err != nil && err != io.EOF {
		return nil,err
	}
	data = data[:readBytes]

	if reader.GetBitsPerSample() == 24 {
		data = convert24BitTo16Bit(data,reader.GetSampleCount())
	}

	return data,nil
}

func GetAudioFormat(reader *wav.Reader) uint8 {
	numChannels := reader.GetNumChannels()
	bitsPerSample := reader.GetBitsPerSample()

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