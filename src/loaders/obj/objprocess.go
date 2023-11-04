package loader

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
)

func (this *OBJLoader) processTriangleFace(posIndices, normalIndices, texCoordIndices []int) (rv []gohome.Mesh3DVertex) {
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

func (this *OBJLoader) processQuadFace(posIndices, normalIndices, texCoordIndices []int) (rv []gohome.Mesh3DVertex) {
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

func (this *OBJLoader) processFaceData(elements [][]string) (rv []gohome.Mesh3DVertex) {
	var posIndices []int
	var normalIndices []int
	var texCoordIndices []int
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
