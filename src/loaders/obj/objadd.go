package loader

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
)

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

func (this *OBJLoader) addToken(token tokenData) {
	if len(this.tokens) == 0 {
		this.tokens = make([][]string, token.index+1)
	} else if token.index+1 > uint32(len(this.tokens)) {
		this.tokens = append(this.tokens, make([][]string, token.index+1-uint32(len(this.tokens)))...)
	}
	this.tokens[token.index] = token.tokens
}
