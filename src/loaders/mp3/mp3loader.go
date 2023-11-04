package loader

import (
	"io/ioutil"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/hajimehoshi/go-mp3"
)

func LoadMP3File(fileName string) (*mp3.Decoder, error) {
	file, err := gohome.Framew.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return decoder, err
	}
	return decoder, nil
}

func ReadAllSamples(decoder *mp3.Decoder) ([]byte, error) {
	return ioutil.ReadAll(decoder)
}

func GetAudioFormat() uint8 {
	return gohome.AUDIO_FORMAT_STEREO16
}
