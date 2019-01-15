package loader

import (
	"errors"
)

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
