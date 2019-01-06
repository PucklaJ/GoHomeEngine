package gohome

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strconv"
	"sync"
	"time"
)

const (
	NUM_GO_ROUTINES_TEXTURE_LOADING uint32 = 10
)

var (
	TEXTURE_PATHS = []string{
		"",
		"textures/",
		"assets/",
		"assets/textures/",
	}
)

func loadImageData(img_data *[]byte, img image.Image, start_width, end_width, max_width, max_height uint32, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	var y uint32
	var x uint32
	var r, g, b, a uint32
	var color color.Color
	for x = start_width; x < max_width && x < end_width; x++ {
		for y = 0; y < max_height; y++ {
			color = img.At(int(x), int(y))
			r, g, b, a = color.RGBA()
			(*img_data)[(x+y*max_width)*4+0] = byte(float64(r) / float64(0xffff) * float64(255.0))
			(*img_data)[(x+y*max_width)*4+1] = byte(float64(g) / float64(0xffff) * float64(255.0))
			(*img_data)[(x+y*max_width)*4+2] = byte(float64(b) / float64(0xffff) * float64(255.0))
			(*img_data)[(x+y*max_width)*4+3] = byte(float64(a) / float64(0xffff) * float64(255.0))
		}
	}
}

func (rsmgr *ResourceManager) LoadTexture(name, path string) Texture {
	start := time.Now()

	if resName, ok := rsmgr.resourceFileNames[path]; ok {
		rsmgr.textures[name] = rsmgr.textures[resName]
		ErrorMgr.Message(ERROR_LEVEL_WARNING, "Texture", name, "Has already been loaded with this or another name!")
		return nil
	}

	if _, ok := rsmgr.textures[name]; ok {
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Texture", name, "Has already been loaded!")
		return nil
	}

	var reader io.Reader
	var err error
	reader, _, err = OpenFileWithPaths(path, TEXTURE_PATHS[:])
	if err != nil {
		ErrorMgr.Message(ERROR_LEVEL_ERROR, "Texture", name, "Couldn't load file: "+err.Error())
		return nil
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		ErrorMgr.Message(ERROR_LEVEL_ERROR, "Texture", name, "Couldn't decode file: "+err.Error())
		return nil
	}

	tex := Render.CreateTexture(name, false)

	tex.LoadFromImage(img)

	if tex != nil {
		rsmgr.textures[name] = tex
		rsmgr.resourceFileNames[path] = name
		end := time.Now()
		sec := end.Sub(start).Seconds()
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Texture", name, "Finished loading! W: "+strconv.Itoa(tex.GetWidth())+" H: "+strconv.Itoa(tex.GetHeight())+" T: "+strconv.FormatFloat(sec, 'f', 3, 64)+" s")
	}
	return tex
}

func (rsmgr *ResourceManager) GetTexture(name string) Texture {
	t := rsmgr.textures[name]
	return t
}

func (rsmgr *ResourceManager) SetTexture(name string, name1 string) {
	s := rsmgr.textures[name1]
	if s == nil {
		ErrorMgr.Message(ERROR_LEVEL_ERROR, "Texture", name, "Couldn't set to "+name1+" (It is nil)")
		return
	}
	rsmgr.textures[name] = s
	ErrorMgr.Message(ERROR_LEVEL_LOG, "Texture", name, "Set to "+name1)
}

func (rsmgr *ResourceManager) DeleteTexture(name string) {
	if _, ok := rsmgr.textures[name]; ok {
		rsmgr.textures[name].Terminate()
		delete(rsmgr.textures, name)
		rsmgr.deleteResourceFileName(name)
		ErrorMgr.Log("Texture", name, "Deleted!")
	} else {
		ErrorMgr.Warning("Texture", name, "Couldn't delete! It hasn't been loaded!")
	}
}
