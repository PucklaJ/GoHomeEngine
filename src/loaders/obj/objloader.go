package loader

import (
	"bufio"
	"errors"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Equals(this *gohome.Mesh3DVertex, other *gohome.Mesh3DVertex) bool {
	for i := 0; i < len(*this); i++ {
		if (*this)[i] != (*other)[i] {
			return false
		}
	}

	return true
}

type OBJWaitGroup struct {
	activeRoutines int
	waitChannel    chan byte
}

func (this *OBJWaitGroup) Add(i int) {
	if this.waitChannel == nil {
		this.waitChannel = make(chan byte)
	}
	this.activeRoutines += i
}

func (this *OBJWaitGroup) Done() {
	this.waitChannel <- '0'
}

func (this *OBJWaitGroup) WaitForDone() bool {
	if this.waitChannel == nil {
		return true
	}
	select {
	case <-this.waitChannel:
		this.activeRoutines--
		if this.activeRoutines == 0 {
			close(this.waitChannel)
			this.waitChannel = nil
			return true
		}
	default:
	}

	return false
}

type positionData struct {
	position [3]float32
	index    uint32
}

type normalData positionData

type texCoordData struct {
	texCoord [2]float32
	index    uint32
}

type tokenData struct {
	tokens []string
	index  uint32
}

type OBJColor [3]float32

type OBJMaterial struct {
	Name             string
	DiffuseColor     OBJColor
	SpecularColor    OBJColor
	Transperancy     float32
	SpecularExponent float32
	DiffuseTexture   string
	SpecularTexture  string
	NormalMap        string
}

type OBJMesh struct {
	Name     string
	Vertices []gohome.Mesh3DVertex
	Indices  []uint32
	Material *OBJMaterial
}

type OBJModel struct {
	Name   string
	Meshes []OBJMesh
}

type MTLLoader struct {
	Materials       []OBJMaterial
	currentMaterial *OBJMaterial
}

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

	positionIndex uint32
	normalIndex   uint32
	texCoordIndex uint32
	tokensIndex   uint32
}

func (this *OBJLoader) Load(path string) error {
	reader, filename, err := gohome.OpenFileWithPaths(path, gohome.LEVEL_PATHS[:])
	if err != nil {
		return err
	}
	this.directory = gohome.GetPathFromFile(filename)
	return this.LoadReader(reader)
}

func (this *OBJLoader) addToken(token tokenData) {
	if len(this.tokens) == 0 {
		this.tokens = make([][]string, token.index+1)
	} else if token.index+1 > uint32(len(this.tokens)) {
		this.tokens = append(this.tokens, make([][]string, token.index+1-uint32(len(this.tokens)))...)
	}
	this.tokens[token.index] = token.tokens
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
			go func(_line string, index uint32) {
				tokens := toTokens(_line)
				this.tokensChan <- tokenData{tokens, index}
				this.tokensWG.Done()
			}(line, this.tokensIndex-1)
		}

	}

	return nil
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
		gohome.Framew.Log("Process Token:", t)
		if err = this.processTokens(t); err != nil {
			return
		}
	}

	return
}

func (this *OBJLoader) LoadReader(reader io.ReadCloser) error {
	if this.DisableGoRoutines {
		return this.parseFileWithoutGoRoutines(reader)
	} else {
		return this.parseFileWithGoRoutines(reader)
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

func (this *OBJLoader) LoadString(contents string) error {
	if !this.DisableGoRoutines {
		this.openChannels()
		defer this.closeChannels()
	}

	var curChar int = 0
	var finished = false
	var line string = ""
	for !finished {
		line, curChar, finished = readLineString(contents, curChar)
		if line != "" {
			if err := this.processTokens(toTokens(line)); err != nil {
				return err
			}
		}
		if finished {
			break
		}
	}

	if !this.DisableGoRoutines {
		if err := this.waitForDataToFinish(); err != nil {
			return err
		}
	}

	return nil
}

func (this *OBJLoader) SetDirectory(dir string) {
	this.directory = dir
}

func readLine(rd *bufio.Reader) (string, error) {
	var str string
	var isPrefix = true
	var err error
	var buf []byte

	for err == nil && isPrefix {
		buf, isPrefix, err = rd.ReadLine()
		str += string(buf)
	}

	return str, err
}

func readLineString(contents string, curChar int) (string, int, bool) {
	var line string
	var finished = false
	var i int

	for i = curChar; i < len(contents); i++ {
		if i == len(contents)-1 {
			finished = true
		}
		if contents[i] == '\n' || contents[i] == '\r' {
			i++
			break
		}
		line += string(contents[i])
	}

	return line, i, finished
}

func toTokens(line string) []string {
	var curByte byte
	var readToken bool = false
	var tokens []string

	for i := 0; i < len(line); i++ {
		curByte = line[i]
		if curByte == ' ' {
			readToken = false
		} else {
			if readToken {
				tokens[len(tokens)-1] += string(curByte)
			} else {
				tokens = append(tokens, "")
				tokens[len(tokens)-1] += string(curByte)
				readToken = true
			}
		}
	}

	return tokens
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

func (this *OBJLoader) newPosition(tokens []string) {
	if this.DisableGoRoutines {
		this.positions = append(this.positions, process3Floats(tokens[1:]))
	} else {
		this.positionIndex++
		this.verticesWG.Add(1)
		go func(index uint32) {
			this.positionChan <- positionData{process3Floats(tokens[1:]), index}
			this.verticesWG.Done()
		}(this.positionIndex - 1)
	}
}

func (this *OBJLoader) newNormal(tokens []string) {
	if this.DisableGoRoutines {
		this.normals = append(this.normals, process3Floats(tokens[1:]))
	} else {
		this.normalIndex++
		this.verticesWG.Add(1)
		go func(index uint32) {
			this.normalChan <- normalData{process3Floats(tokens[1:]), index}
			this.verticesWG.Done()
		}(this.normalIndex - 1)
	}
}

func (this *OBJLoader) newTexCoord(tokens []string) {
	if this.DisableGoRoutines {
		vt := process2Floats(tokens[1:])
		vt[1] = 1.0 - vt[1]
		this.texCoords = append(this.texCoords, vt)
	} else {
		this.texCoordIndex++
		this.verticesWG.Add(1)
		go func(index uint32) {
			vt := process2Floats(tokens[1:])
			vt[1] = 1.0 - vt[1]
			this.texCoordChan <- texCoordData{vt, index}

			this.verticesWG.Done()
		}(this.texCoordIndex - 1)
	}
}

func (this *OBJLoader) newFace(tokens []string) error {
	if len(this.Models) == 0 {
		this.Models = append(this.Models, OBJModel{Name: "Default"})
		this.currentModel = &this.Models[len(this.Models)-1]
	}
	if len(this.currentModel.Meshes) == 0 {
		this.currentModel.Meshes = append(this.currentModel.Meshes, OBJMesh{Name: "Default"})
		this.currentMesh = &this.currentModel.Meshes[len(this.currentModel.Meshes)-1]
		if !this.DisableGoRoutines {
			this.waitForDataToFinish()
		}
	}
	if this.DisableGoRoutines {
		if err := this.processFace(tokens[1:]); err != nil {
			return err
		}
	} else {
		this.facesWG.Add(1)
		go func() {
			if err := this.processFace(tokens[1:]); err != nil {
				this.errorChan <- err
			}
			this.facesWG.Done()
		}()
	}

	return nil
}

func (this *OBJLoader) newMaterialFile(tokens []string) error {
	if this.DisableGoRoutines {
		if err := this.loadMaterialFile(tokens[1]); err != nil {
			return errors.New("Couldn't load material file " + tokens[1] + ": " + err.Error())
		}
	} else {
		this.materialWG.Add(1)
		go func() {
			if err := this.loadMaterialFile(tokens[1]); err != nil {
				this.errorChan <- errors.New("Couldn't load material file " + tokens[1] + ": " + err.Error())
			}
			this.materialWG.Done()
		}()
	}

	return nil
}

func (this *OBJLoader) newModel(tokens []string) error {
	if !this.DisableGoRoutines {
		if err := this.waitForDataToFinish(); err != nil {
			return err
		}
	}

	this.positions = this.positions[:0]
	this.normals = this.normals[:0]
	this.texCoords = this.texCoords[:0]
	this.positionIndex = 0
	this.normalIndex = 0
	this.texCoordIndex = 0

	this.Models = append(this.Models, OBJModel{Name: tokens[1]})
	this.currentModel = &this.Models[len(this.Models)-1]

	return nil
}

func (this *OBJLoader) newMesh(tokens []string) error {
	if !this.DisableGoRoutines {
		if err := this.waitForDataToFinish(); err != nil {
			return err
		}
	}

	if len(this.Models) == 0 {
		this.Models = append(this.Models, OBJModel{Name: "Default"})
		this.currentModel = &this.Models[len(this.Models)-1]
	}
	this.currentModel.Meshes = append(this.currentModel.Meshes, OBJMesh{})
	this.currentMesh = &this.currentModel.Meshes[len(this.currentModel.Meshes)-1]
	this.processMaterial(tokens[1])

	return nil
}

func (this *OBJLoader) processTokens(tokens []string) error {
	length := len(tokens)
	if length != 0 {
		if tokens[0] == "f" {
			if err := this.newFace(tokens); err != nil {
				return err
			}
		} else {
			if length == 4 {
				if tokens[0] == "v" {
					this.newPosition(tokens)
				} else if tokens[0] == "vn" {
					this.newNormal(tokens)
				}
			} else if length == 3 {
				if tokens[0] == "vt" {
					this.newTexCoord(tokens)
				}
			} else if length == 2 {
				if tokens[0] == "mtllib" {
					if err := this.newMaterialFile(tokens); err != nil {
						return err
					}
				} else if tokens[0] == "o" {
					if err := this.newModel(tokens); err != nil {
						return err
					}
				} else if tokens[0] == "usemtl" {
					if err := this.newMesh(tokens); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (this *OBJLoader) addPosition(pos positionData) {
	if len(this.positions) == 0 {
		this.positions = make([][3]float32, pos.index+1)
	} else if pos.index+1 > uint32(len(this.positions)) {
		this.positions = append(this.positions, make([][3]float32, pos.index+1-uint32(len(this.positions)))...)
	}

	this.positions[pos.index] = pos.position
}

func (this *OBJLoader) addNormal(norm normalData) {
	if len(this.normals) == 0 {
		this.normals = make([][3]float32, norm.index+1)
	} else if norm.index+1 > uint32(len(this.normals)) {
		this.normals = append(this.normals, make([][3]float32, norm.index+1-uint32(len(this.normals)))...)
	}

	this.normals[norm.index] = norm.position
}

func (this *OBJLoader) addTexCoord(texCoord texCoordData) {
	if len(this.texCoords) == 0 {
		this.texCoords = make([][2]float32, texCoord.index+1)
	} else if texCoord.index+1 > uint32(len(this.texCoords)) {
		this.texCoords = append(this.texCoords, make([][2]float32, texCoord.index+1-uint32(len(this.texCoords)))...)
	}

	this.texCoords[texCoord.index] = texCoord.texCoord
}

func (this *OBJLoader) addFace(face []gohome.Mesh3DVertex) {
	for i := 0; i < len(face); i++ {
		index, isNew := this.searchIndex(face[i])
		if isNew {
			this.currentMesh.Vertices = append(this.currentMesh.Vertices, face[i])
			this.currentMesh.Indices = append(this.currentMesh.Indices, uint32(len(this.currentMesh.Vertices))-1)
		} else {
			this.currentMesh.Indices = append(this.currentMesh.Indices, index)
		}
	}
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

func (this *OBJLoader) processMaterial(token string) {
	this.currentMesh.Name = token
	this.currentMesh.Material = this.getMaterial(token)
}

func (this *OBJLoader) processFace(tokens []string) error {
	if len(tokens) != 3 && len(tokens) != 4 {
		return errors.New("Face type not supported: " + strconv.FormatInt(int64(len(tokens)), 10) + "! Use triangles or quads!")
	}

	var elements [][]string
	elements = make([][]string, len(tokens))
	for i := 0; i < len(tokens); i++ {
		elements[i] = strings.Split(tokens[i], "/")
	}

	if this.faceMethod == 0 {
		if len(elements[0]) == 1 {
			this.faceMethod = 1
			this.normalsLoaded = false
			this.texCoordsLoaded = false
		} else if len(elements[0]) == 2 {
			this.faceMethod = 2
			this.normalsLoaded = false
			this.texCoordsLoaded = true
		} else if len(elements[0]) == 3 && elements[0][1] == "" {
			this.faceMethod = 3
			this.normalsLoaded = true
			this.texCoordsLoaded = false
		} else if len(elements[0]) == 3 && elements[0][1] != "" {
			this.faceMethod = 4
			this.normalsLoaded = true
			this.texCoordsLoaded = true
		}
	}

	if this.faceMethod == 0 {
		return errors.New("Format of faces is not supported")
	}

	vertices := this.processFaceData(elements)

	if this.DisableGoRoutines {
		this.addFace(vertices)
	} else {
		this.facesChan <- vertices
	}

	return nil
}

func (this *OBJLoader) searchIndex(vertex gohome.Mesh3DVertex) (uint32, bool) {
	for i := 0; i < len(this.currentMesh.Vertices); i++ {
		if Equals(&vertex, &this.currentMesh.Vertices[i]) {
			return uint32(i), false
		}
	}

	return 0, true
}

func (this *OBJLoader) processTriangleFace(posIndices, normalIndices, texCoordIndices []uint32) (rv []gohome.Mesh3DVertex) {
	rv = make([]gohome.Mesh3DVertex, 3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			rv[i][j] = this.positions[posIndices[i]-1][j]
		}
		if this.normalsLoaded {
			for j := 0; j < 3; j++ {
				ni := normalIndices[i] - 1
				normal := this.normals[ni]
				float := normal[j]
				rv[i][j+3] = float
			}
		}
		if this.texCoordsLoaded {
			for j := 0; j < 2; j++ {
				rv[i][j+3+3] = this.texCoords[texCoordIndices[i]-1][j]
			}
		}
	}
	return
}

var quadIs = [6]int{
	0, 1, 2, 2, 3, 0,
}

func (this *OBJLoader) processQuadFace(posIndices, normalIndices, texCoordIndices []uint32) (rv []gohome.Mesh3DVertex) {
	rv = make([]gohome.Mesh3DVertex, 6)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for i := 0; i < 6; i++ {
			for j := 0; j < 3; j++ {
				rv[i][j] = this.positions[posIndices[quadIs[i]]-1][j]
			}
		}
		wg.Done()
	}()
	if this.normalsLoaded {
		wg.Add(1)
		go func() {
			for i := 0; i < 6; i++ {
				for j := 0; j < 3; j++ {
					rv[i][j+3] = this.normals[normalIndices[quadIs[i]]-1][j]
				}
			}
			wg.Done()
		}()
	}
	if this.texCoordsLoaded {
		wg.Add(1)
		go func() {
			for i := 0; i < 6; i++ {
				for j := 0; j < 2; j++ {
					rv[i][j+3+3] = this.texCoords[texCoordIndices[quadIs[i]]-1][j]
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return
}

func (this *OBJLoader) processFaceData(elements [][]string) (rv []gohome.Mesh3DVertex) {
	var posIndices []uint32
	var normalIndices []uint32
	var texCoordIndices []uint32
	switch this.faceMethod {
	case 1:
		posIndices = processFaceData1(elements)
	case 2:
		posIndices, texCoordIndices = processFaceData2(elements)
	case 3:
		posIndices, normalIndices = processFaceData3(elements)
	case 4:
		posIndices, texCoordIndices, normalIndices = processFaceData4(elements)
	}
	if len(elements) == 3 {
		rv = this.processTriangleFace(posIndices, normalIndices, texCoordIndices)
	} else if len(elements) == 4 {
		rv = this.processQuadFace(posIndices, normalIndices, texCoordIndices)
	}

	return
}

func processFaceData1(elements [][]string) (rv []uint32) {
	rv = make([]uint32, len(elements))
	for i := 0; i < len(rv); i++ {
		temp, _ := strconv.ParseUint(elements[i][0], 10, 32)
		rv[i] = uint32(temp)
	}
	return
}

func processFaceData2(elements [][]string) (pos []uint32, tex []uint32) {
	rv := make([]uint32, len(elements)*2)
	pos = make([]uint32, len(elements))
	tex = make([]uint32, len(elements))
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			temp, _ := strconv.ParseUint(elements[i][j], 10, 32)
			rv[i*2+j] = uint32(temp)
		}
	}

	for i := 0; i < len(elements); i++ {
		pos[i] = rv[i*2]
		tex[i] = rv[i*2+1]
	}

	return
}

func processFaceData3(elements [][]string) (pos []uint32, norm []uint32) {
	rv := make([]uint32, len(elements)*2)
	pos = make([]uint32, len(elements))
	norm = make([]uint32, len(elements))
	var readIndex uint32
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			if j == 1 {
				readIndex = 2
			} else {
				readIndex = uint32(j)
			}
			temp, _ := strconv.ParseUint(elements[i][readIndex], 10, 32)
			rv[i*2+j] = uint32(temp)
		}
	}

	for i := 0; i < len(elements); i++ {
		pos[i] = rv[i*2]
		norm[i] = rv[i*2+1]
	}

	return
}

func processFaceData4(elements [][]string) (pos []uint32, tex []uint32, norm []uint32) {
	rv := make([]uint32, len(elements)*3)
	pos = make([]uint32, len(elements))
	tex = make([]uint32, len(elements))
	norm = make([]uint32, len(elements))
	for i := 0; i < len(elements); i++ {
		for j := 0; j < 3; j++ {
			temp, _ := strconv.ParseUint(elements[i][j], 10, 32)
			rv[i*3+j] = uint32(temp)
		}
	}

	for i := 0; i < len(elements); i++ {
		pos[i] = rv[i*3]
		tex[i] = rv[i*3+1]
		norm[i] = rv[i*3+2]
	}

	return
}

func process3Floats(tokens []string) [3]float32 {
	var rv [3]float32
	var temp float64
	var err error

	for i := 0; i < 3; i++ {
		temp, err = strconv.ParseFloat(tokens[i], 32)
		if err != nil {
			return [3]float32{0.0, 0.0, 0.0}
		}
		rv[i] = float32(temp)
	}

	return rv
}

func process2Floats(tokens []string) [2]float32 {
	var rv [2]float32
	var temp float64
	var err error

	for i := 0; i < 2; i++ {
		temp, err = strconv.ParseFloat(tokens[i], 32)
		if err != nil {
			return [2]float32{0.0, 0.0}
		}
		rv[i] = float32(temp)
	}

	return rv
}

func process1Float(tokens string) float32 {
	var rv float32
	var temp float64
	var err error

	temp, err = strconv.ParseFloat(tokens, 32)
	if err != nil {
		return 0.0
	}
	rv = float32(temp)

	return rv
}

func (this *MTLLoader) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return this.LoadReader(file)
}

func (this *MTLLoader) LoadReader(reader io.ReadCloser) error {
	var err error
	var line string

	rd := bufio.NewReader(reader)
	defer reader.Close()

	for err != io.EOF {
		line, err = readLine(rd)

		if err != nil && err != io.EOF {
			return err
		}
		if line != "" {
			this.processTokens(toTokens(line))
		}
	}

	return nil
}

func (this *MTLLoader) checkCurrentMaterial() {
	if this.currentMaterial == nil {
		this.Materials = append(this.Materials, OBJMaterial{Name: "Default"})
		this.currentMaterial = &this.Materials[len(this.Materials)-1]
	}
}

func addAllTokens(tokens []string, start int) (str string) {
	for i := start; i < len(tokens); i++ {
		if i == start {
			str += " "
		}
		str += tokens[i]
	}
	return
}

func (this *MTLLoader) processTokens(tokens []string) {
	length := len(tokens)
	if length > 0 {
		if tokens[0] == "newmtl" {
			this.Materials = append(this.Materials, OBJMaterial{Name: tokens[1] + addAllTokens(tokens, 2)})
			this.currentMaterial = &this.Materials[len(this.Materials)-1]
		} else if tokens[0] == "map_Kd" {
			this.checkCurrentMaterial()
			this.currentMaterial.DiffuseTexture = tokens[1] + addAllTokens(tokens, 2)
		} else if tokens[0] == "map_Ks" {
			this.checkCurrentMaterial()
			this.currentMaterial.SpecularTexture = tokens[1] + addAllTokens(tokens, 2)
		} else if tokens[0] == "norm" {
			this.checkCurrentMaterial()
			this.currentMaterial.NormalMap = tokens[1] + addAllTokens(tokens, 2)
		}
		if length == 2 {
			if tokens[0] == "Ns" {
				this.checkCurrentMaterial()
				this.currentMaterial.SpecularExponent = process1Float(tokens[1])
			} else if tokens[0] == "d" {
				this.checkCurrentMaterial()
				this.currentMaterial.Transperancy = process1Float(tokens[1])
			}
		} else if length == 4 {
			if tokens[0] == "Kd" {
				this.checkCurrentMaterial()
				this.currentMaterial.DiffuseColor = process3Floats(tokens[1:])
			} else if tokens[0] == "Ks" {
				this.checkCurrentMaterial()
				this.currentMaterial.SpecularColor = process3Floats(tokens[1:])
			}
		}
	}
}
