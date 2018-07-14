package loader

import (
	"github.com/cryptix/wav"
	"os"
	"io"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

func LoadWAVFile(fileName string) (*wav.Reader, error) {
	wavInfo,err := os.Stat(fileName)
	if err != nil {
		return nil,err
	}

	file,err := os.Open(fileName)
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

	return data,nil
}

func GetAudioFormat(reader *wav.Reader) uint8 {
	numChannels := reader.GetNumChannels()
	bitsPerSample := reader.GetBitsPerSample()
	bitsPerChannel := bitsPerSample / numChannels

	switch numChannels {
	case 1:
		switch bitsPerChannel {
		case 8:
			return gohome.AUDIO_FORMAT_MONO8
		case 16:
			return gohome.AUDIO_FORMAT_MONO16
		}
	case 2:
		switch bitsPerChannel {
		case 8:
			return gohome.AUDIO_FORMAT_STEREO8
		case 16:
			return gohome.AUDIO_FORMAT_STEREO16
		}
	}

	return gohome.AUDIO_FORMAT_UNKNOWN
}