package gohome

import (
	"strconv"
	"sync"
)

type preloadedTexture struct {
	Name string
	Path string
}

type preloadedShader struct {
	Name                      string
	VertexShader              string
	FragmentShader            string
	GeometryShader            string
	TesselletionControlShader string
	EveluationShader          string
	ComputeShader             string
}

type preloadedLevel struct {
	Name      string
	Path      string
	LoadToGPU bool
}

type preloadedTextureData struct {
	Tex      Texture
	img_data []byte
	width    int
	height   int
}

type preloadedShaderData struct {
	name     string
	contents [6]string
}

type preloadedLevelObject struct {
	Lvl    *Level
	Lvlobj LevelObject
}

type PreloadedMesh struct {
	Mesh      Mesh3D
	LoadToGPU bool
}

type preloader struct {
	preloadedTextures []preloadedTexture
	preloadedShaders  []preloadedShader
	preloadedLevels   []preloadedLevel

	preloadedShaderDataChan  chan preloadedShaderData
	preloadedLevelsChan      chan *Level
	PreloadedModelsChan      chan *Model3D
	PreloadedMeshesChan      chan PreloadedMesh
	preloadedTextureDataChan chan preloadedTextureData
	exitChan                 chan bool
	exitLevelsChan           chan bool
	exitTexturesChan         chan bool
	exitShadersChan          chan bool

	preloadedTexturesToFinish []preloadedTextureData
	preloadedShadersToFinish  []preloadedShaderData
	PreloadedMeshesToFinish   []PreloadedMesh
}

func (this *preloader) Init() {
	this.preloadedShaderDataChan = make(chan preloadedShaderData)
	this.preloadedLevelsChan = make(chan *Level)
	this.PreloadedModelsChan = make(chan *Model3D)
	this.PreloadedMeshesChan = make(chan PreloadedMesh)
	this.preloadedTextureDataChan = make(chan preloadedTextureData)
	this.exitChan = make(chan bool)
	this.exitLevelsChan = make(chan bool)
	this.exitTexturesChan = make(chan bool)
	this.exitShadersChan = make(chan bool)
}

func (this *preloader) loadPreloadedLevel(lvl *preloadedLevel, wg *sync.WaitGroup) {
	defer wg.Done()

	name := lvl.Name
	path := lvl.Path

	level := ResourceMgr.loadLevel(name, path, true, lvl.LoadToGPU)
	if level != nil {
		this.preloadedLevelsChan <- level
	}
}

func (this *preloader) loadPreloadedLevels() {
	if len(this.preloadedLevels) == 0 {

	} else {
		var wg1 sync.WaitGroup
		wg1.Add(len(this.preloadedLevels))
		for i := 0; i < len(this.preloadedLevels); i++ {
			go this.loadPreloadedLevel(&this.preloadedLevels[i], &wg1)
		}
		wg1.Wait()
	}

	go func() {
		this.exitLevelsChan <- true
	}()

}

func (this *preloader) loadPreloadedShader(s *preloadedShader, wg *sync.WaitGroup) {
	defer wg.Done()
	name := s.Name
	vertex_path := s.VertexShader
	fragment_path := s.FragmentShader
	geometry_path := s.GeometryShader
	tesselletion_control_path := s.TesselletionControlShader
	eveluation_path := s.EveluationShader
	compute_path := s.ComputeShader
	ResourceMgr.loadShader(name, vertex_path, fragment_path, geometry_path, tesselletion_control_path, eveluation_path, compute_path, true)
}

func (this *preloader) loadPreloadedShaders() {
	if len(this.preloadedShaders) == 0 {

	} else {
		var wg1 sync.WaitGroup
		wg1.Add(len(this.preloadedShaders))
		for i := 0; i < len(this.preloadedShaders); i++ {
			go this.loadPreloadedShader(&this.preloadedShaders[i], &wg1)
		}
		wg1.Wait()
	}

	go func() {
		this.exitShadersChan <- true
	}()
}

func (this *preloader) loadPreloadedTexture(tex *preloadedTexture, wg *sync.WaitGroup) {
	defer wg.Done()

	name := tex.Name
	path := tex.Path

	ResourceMgr.LoadTextureFunction(name, path, true)
}

func (this *preloader) loadPreloadedTextures() {
	if len(this.preloadedTextures) == 0 {

	} else {
		var wg1 sync.WaitGroup
		wg1.Add(len(this.preloadedTextures))
		for i := 0; i < len(this.preloadedTextures); i++ {
			go this.loadPreloadedTexture(&this.preloadedTextures[i], &wg1)
		}
		wg1.Wait()
	}

	go func() {
		this.exitTexturesChan <- true
	}()
}

func (this *preloader) finish(wg *sync.WaitGroup) {
	defer wg.Done()

	var done bool = false

	for true {
		select {
		case lvl := <-this.preloadedLevelsChan:
			ResourceMgr.Levels[lvl.Name] = lvl
			ErrorMgr.Message(ERROR_LEVEL_LOG, "Level", lvl.Name, "Finished loading!")
		case tex := <-this.preloadedTextureDataChan:
			this.preloadedTexturesToFinish = append(this.preloadedTexturesToFinish, tex)
		case shader := <-this.preloadedShaderDataChan:
			this.preloadedShadersToFinish = append(this.preloadedShadersToFinish, shader)
		case mesh := <-this.PreloadedMeshesChan:
			this.PreloadedMeshesToFinish = append(this.PreloadedMeshesToFinish, mesh)
		case model := <-this.PreloadedModelsChan:
			ResourceMgr.Models[model.Name] = model
			ErrorMgr.Message(ERROR_LEVEL_LOG, "Model", model.Name, "Finished loading!")
		case <-this.exitChan:
			done = true
		default:
		}
		if done {
			break
		}
	}
}

func (this *preloader) checkExit(wg *sync.WaitGroup) {
	defer wg.Done()

	var LevelsExit, texturesExit, shadersExit bool = false, false, false
	var done bool = false

	for true {
		select {
		case <-this.exitLevelsChan:
			LevelsExit = true
		case <-this.exitTexturesChan:
			texturesExit = true
		case <-this.exitShadersChan:
			shadersExit = true
		default:
			if LevelsExit && texturesExit && shadersExit {
				this.exitChan <- true
				done = true
			}
		}
		if done {
			break
		}
	}

}

func (this *preloader) finishTextures() {
	for i := 0; i < len(this.preloadedTexturesToFinish); i++ {
		tex := this.preloadedTexturesToFinish[i]
		tex.Tex.Load(tex.img_data, tex.width, tex.height, false)
		ResourceMgr.textures[tex.Tex.GetName()] = tex.Tex
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Texture", tex.Tex.GetName(), "Finished loading! W: "+strconv.Itoa(tex.width)+" H: "+strconv.Itoa(tex.height))
	}
}

func (this *preloader) finishShaders() {
	for i := 0; i < len(this.preloadedShadersToFinish); i++ {
		shader := this.preloadedShadersToFinish[i]
		s, err := Render.LoadShader(shader.name, shader.contents[VERTEX], shader.contents[FRAGMENT], shader.contents[GEOMETRY], shader.contents[TESSELLETION], shader.contents[EVELUATION], shader.contents[COMPUTE])
		if s != nil {
			ResourceMgr.shaders[shader.name] = s
			ErrorMgr.Message(ERROR_LEVEL_LOG, "Shader", s.GetName(), "Finished loading!")
		} else {
			ErrorMgr.MessageError(ERROR_LEVEL_ERROR, "Shader", s.GetName(), err)
		}
	}
}

func (this *preloader) finishMeshes() {
	for i := 0; i < len(this.PreloadedMeshesToFinish); i++ {
		mesh := this.PreloadedMeshesToFinish[i]
		if mesh.LoadToGPU {
			mesh.Mesh.Load()
		}
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Mesh", mesh.Mesh.GetName(), "Finished loading! V: "+strconv.Itoa(int(mesh.Mesh.GetNumVertices()))+" I: "+strconv.Itoa(int(mesh.Mesh.GetNumIndices())))
	}
}

func (this *preloader) finishData() {
	this.finishTextures()
	this.finishShaders()
	this.finishMeshes()
}

func (this *preloader) loadPreloadedResources() {
	var wg sync.WaitGroup
	wg.Add(2)

	go this.checkExit(&wg)
	go this.finish(&wg)

	go this.loadPreloadedLevels()
	go this.loadPreloadedShaders()
	go this.loadPreloadedTextures()

	wg.Wait()

	this.finishData()
}
