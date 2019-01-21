package loader

import (
	"bufio"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"io"
)

type OBJLoader struct {
	Models            []OBJModel
	MaterialLoader    MTLLoader
	DisableGoRoutines bool

	positions       [][3]float32
	normals         [][3]float32
	texCoords       [][2]float32
	faceMethod      uint8
	normalsLoaded   bool
	texCoordsLoaded bool
	directory       string

	tokens [][]string

	currentModel *OBJModel
	currentMesh  *OBJMesh

	verticesWG OBJWaitGroup
	facesWG    OBJWaitGroup
	materialWG OBJWaitGroup
	tokensWG   OBJWaitGroup

	positionChan chan positionData
	normalChan   chan normalData
	texCoordChan chan texCoordData
	facesChan    chan []gohome.Mesh3DVertex
	errorChan    chan error
	tokensChan   chan tokenData

	positionIndex int
	normalIndex   int
	texCoordIndex int
	tokensIndex   int
}

func (this *OBJLoader) Load(path string) error {
	reader, filename, err := gohome.OpenFileWithPaths(path, gohome.LEVEL_PATHS[:])
	if err != nil {
		return err
	}
	this.directory = gohome.GetPathFromFile(filename)
	return this.LoadReader(reader)
}

func (this *OBJLoader) LoadReader(reader io.ReadCloser) error {
	if this.DisableGoRoutines {
		return this.parseFileWithoutGoRoutines(reader)
	} else {
		return this.parseFileWithGoRoutines(reader)
	}
}

func (this *OBJLoader) SetDirectory(dir string) {
	this.directory = dir
}

func (this *OBJLoader) parseFileWithoutGoRoutines(reader io.ReadCloser) error {
	var err error
	var line string
	var rd *bufio.Reader

	rd = bufio.NewReader(reader)
	defer reader.Close()

	for err != io.EOF {
		line, err = readLine(rd)
		if err != nil && err != io.EOF {
			return err
		}
		if line != "" {
			if err = this.processTokens(toTokens(line)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *OBJLoader) parseFileWithGoRoutines(reader io.ReadCloser) (err error) {
	this.openChannels()
	defer this.closeChannels()
	defer func() {
		var err1 error
		if err1 = this.waitForDataToFinish(); err1 != nil {
			if err == nil {
				err = err1
			}
		}
	}()

	this.readTokens(reader)
	for _, t := range this.tokens {
		if err = this.processTokens(t); err != nil {
			return
		}
	}

	return
}

func (this *OBJLoader) readTokens(reader io.ReadCloser) error {
	this.tokensChan = make(chan tokenData)
	defer close(this.tokensChan)
	defer this.waitForTokens()

	var err error
	var line string
	var rd *bufio.Reader

	rd = bufio.NewReader(reader)
	defer reader.Close()

	for err != io.EOF {
		line, err = readLine(rd)

		if err != nil && err != io.EOF {
			return err
		}

		if line != "" {
			this.tokensIndex++
			this.tokensWG.Add(1)
			go func(_line string, index int) {
				tokens := toTokens(_line)
				this.tokensChan <- tokenData{tokens, index}
				this.tokensWG.Done()
			}(line, this.tokensIndex-1)
		}

	}

	return nil
}

func (this *OBJLoader) waitForTokens() {
	for !this.tokensWG.WaitForDone() {
		select {
		case token := <-this.tokensChan:
			this.addToken(token)
		default:
		}
	}
}

func (this *OBJLoader) openChannels() {
	if !this.DisableGoRoutines {
		this.positionChan = make(chan positionData)
		this.normalChan = make(chan normalData)
		this.texCoordChan = make(chan texCoordData)
		this.facesChan = make(chan []gohome.Mesh3DVertex)
		this.errorChan = make(chan error)
	}
}

func (this *OBJLoader) closeChannels() {
	close(this.positionChan)
	close(this.normalChan)
	close(this.texCoordChan)
	close(this.facesChan)
	close(this.errorChan)
}

func (this *OBJLoader) loadMaterialFile(path string) error {
	var err error
	var reader io.ReadCloser

	reader, _, err = gohome.OpenFileWithPaths(path, append([]string{this.directory}, gohome.MATERIAL_PATHS[:]...))
	if err == nil {
		err = this.MaterialLoader.LoadReader(reader)
	}

	return err
}

func (this *OBJLoader) waitForDataToFinish() (err1 error) {
	for {
		select {
		case pos := <-this.positionChan:
			this.addPosition(pos)
		case norm := <-this.normalChan:
			this.addNormal(norm)
		case texCoord := <-this.texCoordChan:
			this.addTexCoord(texCoord)
		case face := <-this.facesChan:
			this.addFace(face)
		case err := <-this.errorChan:
			if err1 == nil {
				err1 = err
			}
		default:
		}

		if this.verticesWG.WaitForDone() && this.materialWG.WaitForDone() && this.facesWG.WaitForDone() {
			break
		}
	}

	return
}

func (this *OBJLoader) getMaterial(name string) *OBJMaterial {
	for i := 0; i < len(this.MaterialLoader.Materials); i++ {
		if this.MaterialLoader.Materials[i].Name == name {
			return &this.MaterialLoader.Materials[i]
		}
	}

	return nil
}

func (this *OBJLoader) searchIndex(vertex gohome.Mesh3DVertex) (uint32, bool) {
	for i := 0; i < len(this.currentMesh.Vertices); i++ {
		if vertex.Equals(&this.currentMesh.Vertices[i]) {
			return uint32(i), false
		}
	}

	return 0, true
}
