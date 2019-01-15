package loader

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"image/color"
)

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
	DiffuseColor     color.Color
	SpecularColor    color.Color
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
