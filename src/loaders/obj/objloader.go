package loader

import (
	"bytes"
	// "fmt"
	"io"
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

type OBJMaterial struct {
}

type OBJMesh struct {
	Name     string
	Vertices []OBJVertex
	Indices  []uint32
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

type OBJLoader struct {
	Models []OBJModel

	positions       [][3]float32
	texCoords       [][2]float32
	normals         [][3]float32
	faceMethod      uint8
	normalsLoaded   bool
	texCoordsLoaded bool
}

func (this *OBJLoader) Load(path string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
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
			this.texCoords = append(this.texCoords, process2Floats(tokens[1:length]))
		}
	} else if length == 2 {
		if tokens[0] == "o" {
			this.Models = append(this.Models, OBJModel{Name: tokens[1]})
		} else if tokens[0] == "usemtl" {
			this.Models[len(this.Models)-1].Meshes = append(this.Models[len(this.Models)-1].Meshes, OBJMesh{})
			this.processMaterial(tokens[1])
		}
	}
}

func (this *OBJLoader) processMaterial(token string) {
	this.Models[len(this.Models)-1].Meshes[len(this.Models[len(this.Models)-1].Meshes)-1].Name = token
}

func (this *OBJLoader) processFace(tokens []string) error {
	if len(this.Models) == 0 {
		this.Models = append(this.Models, OBJModel{Name: "Default"})
	}
	if len(this.Models[len(this.Models)-1].Meshes) == 0 {
		this.Models[len(this.Models)-1].Meshes = append(this.Models[len(this.Models)-1].Meshes, OBJMesh{Name: "Default"})
	}

	this.faceMethod = 0
	var elements [3][]string
	for i := 0; i < 3; i++ {
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

	for i := 0; i < 3; i++ {
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

func (this *OBJLoader) processFaceData(elements [3][]string) [3]OBJVertex {
	var rv [3]OBJVertex
	var posIndices [3]uint32
	var normalIndices [3]uint32
	var texCoordIndices [3]uint32

	switch this.faceMethod {
	case 1:
		posIndices[0], posIndices[1], posIndices[2] = processFaceData1(elements)
	case 2:
		posIndices[0], texCoordIndices[0], posIndices[1], texCoordIndices[1], posIndices[2], texCoordIndices[2] = processFaceData2(elements)
	case 3:
		posIndices[0], normalIndices[0], posIndices[1], normalIndices[1], posIndices[2], normalIndices[2] = processFaceData3(elements)
	case 4:
		posIndices[0], texCoordIndices[0], normalIndices[0], posIndices[1], texCoordIndices[1], normalIndices[1], posIndices[2], texCoordIndices[2], normalIndices[2] = processFaceData4(elements)
	}
	for i := 0; i < 3; i++ {
		rv[i].Position = this.positions[posIndices[i]-1]
		if this.texCoordsLoaded {
			rv[i].TextureCoord = this.texCoords[texCoordIndices[i]-1]
		}
		if this.normalsLoaded {
			rv[i].Normal = this.normals[normalIndices[i]-1]
		}
	}

	return rv
}

func processFaceData1(elements [3][]string) (uint32, uint32, uint32) {
	var rv [3]uint32
	for i := 0; i < 3; i++ {
		temp, _ := strconv.ParseUint(elements[i][0], 10, 32)
		rv[i] = uint32(temp)
	}
	return rv[0], rv[1], rv[2]
}

func processFaceData2(elements [3][]string) (uint32, uint32, uint32, uint32, uint32, uint32) {
	var rv [6]uint32
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			temp, _ := strconv.ParseUint(elements[i][j], 10, 32)
			rv[i*2+j] = uint32(temp)
		}
	}
	return rv[0], rv[1], rv[2], rv[3], rv[4], rv[5]
}

func processFaceData3(elements [3][]string) (uint32, uint32, uint32, uint32, uint32, uint32) {
	var rv [6]uint32
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
	return rv[0], rv[1], rv[2], rv[3], rv[4], rv[5]
}

func processFaceData4(elements [3][]string) (uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32) {
	var rv [9]uint32
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			temp, _ := strconv.ParseUint(elements[i][j], 10, 32)
			rv[i*3+j] = uint32(temp)
		}
	}
	return rv[0], rv[1], rv[2], rv[3], rv[4], rv[5], rv[6], rv[7], rv[8]
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
