package gohome

import (
	// "fmt"
	"github.com/blezek/tga"
	"github.com/raedatoui/assimp"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	NUM_GO_ROUTINES_TEXTURE_LOADING uint32 = 10
)

type ResourceManager struct {
	textures map[string]Texture
	shaders  map[string]Shader
	models   map[string]*Model3D
	levels   map[string]*Level

	preloader
}

func (rsmgr *ResourceManager) Init() {
	rsmgr.textures = make(map[string]Texture)
	rsmgr.shaders = make(map[string]Shader)
	rsmgr.models = make(map[string]*Model3D)
	rsmgr.levels = make(map[string]*Level)

	rsmgr.preloader.Init()

	tga.RegisterFormat()
}

func loadShaderFile(path string) (string, error) {
	if path == "" {
		return "", nil
	} else {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		} else {
			return string(contents[:len(contents)]), nil
		}
	}
}

func loadShader(path, name_shader, name string) (string, bool) {
	contents, err := loadShaderFile(path)
	if err != nil {
		log.Println("Couldn't load", name, "of", name_shader, ":", err)
		return "", true
	}
	return contents, false
}

func (rsmgr *ResourceManager) LoadShader(name, vertex_path, fragment_path, geometry_path, tesselletion_control_path, eveluation_path, compute_path string) {

	shader := rsmgr.loadShader(name, vertex_path, fragment_path, geometry_path, tesselletion_control_path, eveluation_path, compute_path, false)
	if shader != nil {
		rsmgr.shaders[name] = shader
		log.Println("Finished loading shader", name, "!")
	}
}

func (rsmgr *ResourceManager) GetShader(name string) Shader {
	s := rsmgr.shaders[name]
	return s
}

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

func (rsmgr *ResourceManager) LoadTexture(name, path string) {

	tex := rsmgr.loadTexture(name, path, false)
	if tex != nil {
		rsmgr.textures[name] = tex
		log.Println("Finished loading texture", name, "W:", tex.GetWidth(), "H:", tex.GetHeight(), "!")
	}
}

func (rsmgr *ResourceManager) GetTexture(name string) Texture {
	t := rsmgr.textures[name]
	return t
}

func (rsmgr *ResourceManager) LoadLevel(name, path string) {

	level := rsmgr.loadLevel(name, path, false)

	if level != nil {
		rsmgr.levels[name] = level
		log.Println("Finished loading Level", name, "!")
	}
}

func (rsmgr *ResourceManager) processNode(node *assimp.Node, scene *assimp.Scene, level *Level, directory string, preloaded bool) {
	if node != scene.RootNode() {
		model := &Model3D{}
		model.Init(node, scene, level, directory, preloaded)
		if !preloaded {
			if _, ok := rsmgr.models[model.Name]; ok {
				log.Println("Model", model.Name, "has already been loaded! Overwritting ...")
			}
			rsmgr.models[model.Name] = model
			log.Println("Finished loading Model", model.Name, "!")
		} else {
			rsmgr.preloadedModelsChan <- model
		}

	}
	for i := 0; i < node.NumChildren(); i++ {
		rsmgr.processNode(node.Children()[i], scene, level, directory, preloaded)
	}
}

func (rsmgr *ResourceManager) GetLevel(name string) *Level {
	l := rsmgr.levels[name]
	return l
}

func (rsmgr *ResourceManager) GetModel(name string) *Model3D {
	m := rsmgr.models[name]
	return m
}

func (rsmgr *ResourceManager) Terminate() {
	for _, v := range rsmgr.shaders {
		v.Terminate()
	}

	for _, v := range rsmgr.textures {
		v.Terminate()
	}

	for _, v := range rsmgr.models {
		v.Terminate()
	}
}

func (rsmgr *ResourceManager) PreloadShader(name, vertex_path, fragment_path, geometry_path, tesselletion_control_path, eveluation_path, compute_path string) {
	rsmgr.preloader.preloadedShaders = append(rsmgr.preloader.preloadedShaders, preloadedShader{
		name,
		vertex_path,
		fragment_path,
		geometry_path,
		tesselletion_control_path,
		eveluation_path,
		compute_path,
	})
}

func (rsmgr *ResourceManager) PreloadTexture(name, path string) {
	tex := preloadedTexture{
		name,
		path,
	}
	rsmgr.preloader.preloadedTextures = append(rsmgr.preloader.preloadedTextures, tex)
}

func (rsmgr *ResourceManager) PreloadLevel(name, path string) {
	rsmgr.preloader.preloadedLevels = append(rsmgr.preloader.preloadedLevels, preloadedLevel{
		name,
		path,
	})
}

func (rsmgr *ResourceManager) loadLevel(name, path string, preloaded bool) *Level {
	if _, ok := rsmgr.levels[name]; ok {
		log.Println("The level with the name", name, "has already been loaded!")
		return nil
	}
	level := &Level{Name: name}
	var scene *assimp.Scene
	if scene = assimp.ImportFile(path, uint(assimp.Process_Triangulate|assimp.Process_FlipUVs|assimp.Process_GenNormals|assimp.Process_OptimizeMeshes)); scene == nil || (scene.Flags()&assimp.SceneFlags_Incomplete) != 0 || scene.RootNode() == nil {
		log.Println("Couldn't load level", name, "with path", path, ":", assimp.GetErrorString())
		return nil
	}

	directory := path
	if index := strings.LastIndex(directory, "/"); index != -1 {
		directory = directory[index:]
	} else {
		directory = ""
	}

	rsmgr.processNode(scene.RootNode(), scene, level, directory, preloaded)

	return level
}

func (rsmgr *ResourceManager) loadShader(name, vertex_path, fragment_path, geometry_path, tesselletion_control_path, eveluation_path, compute_path string, preloaded bool) Shader {
	_, already := rsmgr.shaders[name]
	if already {
		log.Println("Shader", name, "has already been loaded! (Use GetShader(name))")
		return nil
	}

	var contents [6]string
	var err bool
	var erro error

	contents[VERTEX], err = loadShader(vertex_path, name, "Vertex File")
	if err {
		log.Println(err)
		return nil
	}
	contents[FRAGMENT], err = loadShader(fragment_path, name, "Fragment File")
	if err {
		log.Println(err)
		return nil
	}
	contents[GEOMETRY], err = loadShader(geometry_path, name, "Geometry File")
	if err {
		log.Println(err)
		return nil
	}
	contents[TESSELLETION], err = loadShader(tesselletion_control_path, name, "Tesselletion File")
	if err {
		log.Println(err)
		return nil
	}
	contents[EVELUATION], err = loadShader(eveluation_path, name, "Eveluation File")
	if err {
		log.Println(err)
		return nil
	}
	contents[COMPUTE], err = loadShader(compute_path, name, "Compute File")
	if err {
		log.Println(err)
		return nil
	}

	var shader Shader = nil
	if !preloaded {
		shader, erro = Render.LoadShader(name, contents[VERTEX], contents[FRAGMENT], contents[GEOMETRY], contents[TESSELLETION], contents[EVELUATION], contents[COMPUTE])
		if erro != nil {
			log.Println(erro)
			return nil
		}
	} else {
		rsmgr.preloader.preloadedShaderDataChan <- preloadedShaderData{
			name,
			contents,
		}
	}

	return shader
}

func (rsmgr *ResourceManager) loadTexture(name, path string, preloaded bool) Texture {
	if _, ok := rsmgr.textures[name]; ok {
		log.Println("The texture", name, "has already been loaded")
		return nil
	}

	reader, err := os.Open(path)
	if err != nil {
		log.Println("Couldn't load texture file", name, ":", err)
		return nil
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Println("Couldn't decode texture file", name, ":", err)
		return nil
	}
	reader.Close()
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	img_data := make([]byte, width*height*4)

	var wg1 sync.WaitGroup
	var i uint32
	deltaWidth := uint32(width) / NUM_GO_ROUTINES_TEXTURE_LOADING
	wg1.Add(int(NUM_GO_ROUTINES_TEXTURE_LOADING))
	for i = 0; i < NUM_GO_ROUTINES_TEXTURE_LOADING; i++ {
		go loadImageData(&img_data, img, i*deltaWidth, (i+1)*deltaWidth, uint32(width), uint32(height), &wg1)
	}
	wg1.Wait()
	var tex Texture
	tex = Render.CreateTexture(name, false)
	if !preloaded {
		tex.Load(img_data, width, height)
	} else {
		rsmgr.preloader.preloadedTextureDataChan <- preloadedTextureData{
			tex,
			img_data,
			width,
			height,
		}
	}

	return tex
}

func (rsmgr *ResourceManager) LoadPreloadedResources() {
	rsmgr.preloader.loadPreloadedResources()
}

var ResourceMgr ResourceManager
