package loader

import (
	"bytes"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type OBJVertex struct {
	Position     [3]float32
	TextureCoord [2]float32
	Normal       [3]float32
}

func (this *OBJVertex) Equals(other *OBJVertex) bool {
	for i := 0; i < 3; i++ {
		if this.Position[i] != other.Position[i] {
			return false
		}
		if this.Normal[i] != other.Normal[i] {
			return false
		}
	}
	for i := 0; i < 2; i++ {
		if this.TextureCoord[i] != other.TextureCoord[i] {
			return false
		}
	}

	return true
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
	Vertices []OBJVertex
	Indices  []uint32
	Material *OBJMaterial
}

type OBJModel struct {
	Name   string
	Meshes []OBJMesh
}

type OBJError struct {
	errorString string
}

func (this *OBJError) Error() string {
	return this.errorString
}

type MTLLoader struct {
	Materials       []OBJMaterial
	currentMaterial *OBJMaterial
}

type OBJLoader struct {
	Models         []OBJModel
	MaterialLoader MTLLoader

	positions        [][3]float32
	texCoords        [][2]float32
	normals          [][3]float32
	faceMethod       uint8
	normalsLoaded    bool
	texCoordsLoaded  bool
	materialPaths    []string
	openMaterialFile func(path string) (*gohome.File, error)
	directory        string
}

func (this *OBJLoader) Load(path string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	this.directory = getPathFromFile(path)
	return this.LoadReader(reader)
}

func (this *OBJLoader) LoadReader(reader io.ReadCloser) error {
	var prevOverFlow []byte
	var err error
	var line string

	for err != io.EOF || len(prevOverFlow) != 0 {
		line, prevOverFlow, err = readLine(reader, prevOverFlow)
		if err != nil && err != io.EOF {
			return err
		}
		if line != "" {
			this.processTokens(toTokens(line))
		}
	}
	reader.Close()

	return nil
}

func (this *OBJLoader) LoadString(contents string) {
	var curChar int = 0
	var finished = false
	var line string = ""
	for !finished {
		line, curChar, finished = readLineString(contents, curChar)
		if finished {
			log.Println("OBJLoader:", this.directory)
			log.Println("Positions:", len(this.positions))
			log.Println("Normals:", len(this.normals))
			log.Println("UVs:", len(this.texCoords))
			return
		}
		if line != "" {
			this.processTokens(toTokens(line))
		}
	}
}

func (this *OBJLoader) SetMaterialPaths(paths []string) {
	this.materialPaths = paths
}

func (this *OBJLoader) SetOpenMaterialFile(function func(path string) (*gohome.File, error)) {
	this.openMaterialFile = function
}

func (this *OBJLoader) SetDirectory(dir string) {
	this.directory = dir
}

func readLine(reader io.Reader, prevOverFlow []byte) (string, []byte, error) {
	var line bytes.Buffer
	var buffer [10]byte
	var overflow []byte
	var breakOut bool = false
	var prevOverFlowRead bool = len(prevOverFlow) == 0
	var n int = 0
	var err error = nil
	var endOfFile bool = false

	for !endOfFile {
		if !prevOverFlowRead {
			n = len(prevOverFlow)
			for i := 0; i < n; i++ {
				buffer[i] = prevOverFlow[i]
			}
			prevOverFlowRead = true
		} else {
			n, err = reader.Read(buffer[:])
		}
		if err == io.EOF {
			endOfFile = true
		}
		for i := 0; i < n; i++ {
			if buffer[i] == '\n' || buffer[i] == '\r' {
				if i < n-1 {
					if buffer[i+1] == '\n' || buffer[i+1] == '\r' {
						if i+2 < n {
							overflow = buffer[i+2 : n]
						}
					} else {
						overflow = buffer[i+1 : n]
					}
				}
				breakOut = true
				break
			} else {
				if err1 := line.WriteByte(buffer[i]); err1 != nil {
					return "", nil, err1
				}
			}
		}
		if breakOut {
			break
		}
	}
	lineString := line.String()
	line.Reset()
	return lineString, overflow, err
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
	var tokens []bytes.Buffer
	var tokensString []string

	for i := 0; i < len(line); i++ {
		curByte = line[i]
		if curByte == ' ' {
			readToken = false
		} else {
			if readToken {
				tokens[len(tokens)-1].WriteByte(curByte)
			} else {
				tokens = append(tokens, bytes.Buffer{})
				tokens[len(tokens)-1].WriteByte(curByte)
				readToken = true
			}
		}
	}

	tokensString = make([]string, len(tokens))
	for i := 0; i < len(tokens); i++ {
		tokensString[i] = tokens[i].String()
	}

	return tokensString
}

func getFileFromPath(path string) string {
	if index := strings.LastIndex(path, "/"); index != -1 {
		return path[index+1:]
	} else {
		return path
	}
}

func getPathFromFile(path string) string {
	if index := strings.LastIndex(path, "/"); index != -1 {
		return path[:index+1]
	} else {
		return ""
	}
}

func (this *OBJLoader) loadMaterialFile(path string) error {
	var err error
	var reader io.ReadCloser

	if len(this.materialPaths) == 0 {
		this.materialPaths = append(this.materialPaths, "")
	}
	if this.openMaterialFile == nil {
		this.openMaterialFile = func(path string) (*gohome.File, error) {
			gFile := &gohome.File{}
			osFile, err := os.Open(path)
			gFile.ReadSeeker = osFile
			gFile.Closer = osFile
			return gFile, err
		}
	}
	for i := 0; i < len(this.materialPaths); i++ {
		reader, err = this.openMaterialFile(this.directory + path)
		if err != nil {
			reader, err = this.openMaterialFile(this.directory + getFileFromPath(path))
			if err != nil {
				reader, err = this.openMaterialFile(this.materialPaths[i] + path)
				if err != nil {
					reader, err = this.openMaterialFile(this.materialPaths[i] + getFileFromPath(path))
				}
			}
		}
		if err == nil {
			err = this.MaterialLoader.LoadReader(reader)
			if err != nil {
				return err
			}
			break
		}
	}

	return err
}

func (this *OBJLoader) processTokens(tokens []string) {
	length := len(tokens)

	if length == 4 {
		if tokens[0] == "v" {
			this.positions = append(this.positions, process3Floats(tokens[1:length]))
		} else if tokens[0] == "vn" {
			this.normals = append(this.normals, process3Floats(tokens[1:length]))
		} else if tokens[0] == "f" {
			this.processFace(tokens[1:length])
		}
	} else if length == 3 {
		if tokens[0] == "vt" {
			uv := process2Floats(tokens[1:length])
			uv[1] = 1.0 - uv[1]
			this.texCoords = append(this.texCoords, uv)
		}
	} else if length == 2 {
		if tokens[0] == "mtllib" {
			if err := this.loadMaterialFile(tokens[1]); err != nil {
				log.Println("OBJLoader: Couldn't load material file", tokens[1], ":", err)
			}
		} else if tokens[0] == "o" {
			this.Models = append(this.Models, OBJModel{Name: tokens[1]})
		} else if tokens[0] == "usemtl" {
			if len(this.Models) == 0 {
				this.Models = append(this.Models, OBJModel{Name: "Default"})
			}
			this.Models[len(this.Models)-1].Meshes = append(this.Models[len(this.Models)-1].Meshes, OBJMesh{})
			this.processMaterial(tokens[1])
		}
	}
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
	this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Name = token
	this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Material = this.getMaterial(token)
}

func (this *OBJLoader) processFace(tokens []string) error {
	if len(tokens) != 3 {
		return &OBJError{"Face type not supported! Use triangles!"}
	}

	if len(this.Models) == 0 {
		this.Models = append(this.Models, OBJModel{Name: "Default"})
	}
	if len(this.Models[len(this.Models)-1].Meshes) == 0 {
		this.Models[len(this.Models)-1].Meshes = append(this.Models[len(this.Models)-1].Meshes, OBJMesh{Name: "Default"})
	}

	this.faceMethod = 0
	var elements [][]string
	elements = make([][]string, len(tokens))
	for i := 0; i < len(tokens); i++ {
		elements[i] = strings.Split(tokens[i], "/")
	}
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

	if this.faceMethod == 0 {
		return &OBJError{"Format of faces is not supported"}
	}

	vertices := this.processFaceData(elements)

	for i := 0; i < len(vertices); i++ {
		index, isNew := this.searchIndex(vertices[i])
		if isNew {
			this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Vertices = append(this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Vertices, vertices[i])
			this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Indices = append(this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Indices, uint32(len(this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Vertices)-1))
		} else {
			this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Indices = append(this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Indices, index)
		}
	}

	return nil
}

func (this *OBJLoader) searchIndex(vertex OBJVertex) (uint32, bool) {
	mesh := &this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1]

	for i := 0; i < len(mesh.Vertices); i++ {
		if vertex.Equals(&mesh.Vertices[i]) {
			return uint32(i), false
		}
	}

	return 0, true
}

func (this *OBJLoader) processTriangleFace(posIndices, normalIndices, texCoordIndices []uint32) (rv []OBJVertex) {
	rv = make([]OBJVertex, 3)
	for i := 0; i < 3; i++ {
		rv[i].Position = this.positions[posIndices[i]-1]
		if this.texCoordsLoaded {
			rv[i].TextureCoord = this.texCoords[texCoordIndices[i]-1]
		}
		if this.normalsLoaded {
			rv[i].Normal = this.normals[normalIndices[i]-1]
		}
	}
	return
}

func (this *OBJLoader) processFaceData(elements [][]string) (rv []OBJVertex) {
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
	for i := 0; i < 3; i++ {
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
	var prevOverFlow []byte
	var err error
	var line string

	for err != io.EOF || len(prevOverFlow) != 0 {
		line, prevOverFlow, err = readLine(reader, prevOverFlow)
		if err != nil && err != io.EOF {
			return err
		}
		if line != "" {
			this.processTokens(toTokens(line))
		}
	}
	reader.Close()

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
			if len(this.Materials) != 0 {
				log.Println("Texture:", this.Materials[len(this.Materials)-1].DiffuseTexture)
			}
			this.Materials = append(this.Materials, OBJMaterial{Name: tokens[1] + addAllTokens(tokens, 2)})
			this.currentMaterial = &this.Materials[len(this.Materials)-1]
		} else if tokens[0] == "map_Kd" {
			this.checkCurrentMaterial()
			log.Println("Diffuse:", tokens)
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
