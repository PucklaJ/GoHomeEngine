package loader

import (
	"bufio"
	"io"
	"os"
)

type MTLLoader struct {
	Materials       []OBJMaterial
	currentMaterial *OBJMaterial
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
				this.currentMaterial.DiffuseColor = toGohomeColor(process3Floats(tokens[1:]))
			} else if tokens[0] == "Ks" {
				this.checkCurrentMaterial()
				this.currentMaterial.SpecularColor = toGohomeColor(process3Floats(tokens[1:]))
			}
		}
	}
}
