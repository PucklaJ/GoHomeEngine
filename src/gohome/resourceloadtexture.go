package gohome

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strconv"
	"sync"
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
	tex := rsmgr.LoadTextureFunction(name, path, false)
	if tex != nil {
		rsmgr.textures[name] = tex
		rsmgr.resourceFileNames[path] = name
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Texture", name, "Finished loading! W: "+strconv.Itoa(tex.GetWidth())+" H: "+strconv.Itoa(tex.GetHeight()))
	}
	return tex
}

func (rsmgr *ResourceManager) GetTexture(name string) Texture {
	t := rsmgr.textures[name]
	return t
}

func (rsmgr *ResourceManager) PreloadTexture(name, path string) {

	tex := preloadedTexture{
		name,
		path,
		false,
	}
	if !rsmgr.checkPreloadedTexture(&tex) {
		return
	}
	rsmgr.preloader.preloadedTextures = append(rsmgr.preloader.preloadedTextures, tex)
}

func (rsmgr *ResourceManager) LoadTextureFunction(name, path string, preloaded bool) Texture {
	if !preloaded {
		if resName, ok := rsmgr.resourceFileNames[path]; ok {
			rsmgr.textures[name] = rsmgr.textures[resName]
			ErrorMgr.Message(ERROR_LEVEL_WARNING, "Texture", name, "Has already been loaded with this or another name!")
			return nil
		}
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

	var img_data []byte
	var width, height int

	tex := Render.CreateTexture(name, false)

	if preloaded {
		width = img.Bounds().Size().X
		height = img.Bounds().Size().Y
		img_data = make([]byte, width*height*4)

		var wg1 sync.WaitGroup
		var i float32
		deltaWidth := float32(width) / float32(NUM_GO_ROUTINES_TEXTURE_LOADING)
		wg1.Add(int(NUM_GO_ROUTINES_TEXTURE_LOADING + 1))
		for i = 0; i <= float32(NUM_GO_ROUTINES_TEXTURE_LOADING); i++ {
			go loadImageData(&img_data, img, uint32(i*deltaWidth), uint32((i+1)*deltaWidth), uint32(width), uint32(height), &wg1)
		}
		wg1.Wait()
	} else {
		tex.LoadFromImage(img)
	}

	if tex == nil {
		return nil
	}
	if preloaded {
		rsmgr.preloader.preloadedTextureDataChan <- preloadedTextureData{
			tex,
			img_data,
			width,
			height,
			path,
		}
	}

	return tex
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

func (rsmgr *ResourceManager) checkPreloadedTexture(texture *preloadedTexture) bool {
	if _, ok := rsmgr.textures[texture.Name]; ok {
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Texture", texture.Name, "Has already been loaded!")
		return false
	}
	if resName, ok := rsmgr.resourceFileNames[texture.Path]; ok {
		rsmgr.textures[texture.Name] = rsmgr.textures[resName]
		ErrorMgr.Message(ERROR_LEVEL_WARNING, "Texture", texture.Name, "Has already been loaded with this or another name!")
		return false
	}
	for i := 0; i < len(rsmgr.preloadedTextures); i++ {
		if rsmgr.preloadedTextures[i].Name == texture.Name {
			ErrorMgr.Message(ERROR_LEVEL_LOG, "Texture", texture.Name, "Has already been preloaded!")
			return false
		} else if rsmgr.preloadedTextures[i].Path == texture.Path {
			ErrorMgr.Message(ERROR_LEVEL_WARNING, "Texture", texture.Name, "Has already been preloaded with this or another name!")
			texture.fileAlreadyPreloaded = true
			return true
		}
	}

	texture.fileAlreadyPreloaded = false

	return true
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
