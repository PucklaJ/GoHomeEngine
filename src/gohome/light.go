package gohome

import (
	// "fmt"
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
	"strconv"
)

const (
	AMBIENT_LIGHT_UNIFORM_NAME          string = "ambientLight"
	POINT_LIGHTS_UNIFORM_NAME           string = "pointLights"
	DIRECTIONAL_LIGHTS_UNIFORM_NAME     string = "directionalLights"
	SPOT_LIGHTS_UNIFORM_NAME            string = "spotLights"
	POSITION_UNIFORM_NAME               string = "position"
	DIRECTION_UNIFORM_NAME              string = "direction"
	DIFFUSE_COLOR_UNIFORM_NAME          string = "diffuseColor"
	SPECULAR_COLOR_UNIFORM_NAME         string = "specularColor"
	INNERCUTOFF_UNIFORM_NAME            string = "innerCutOff"
	OUTERCUTOFF_UNIFORM_NAME            string = "outerCutOff"
	ATTENTUATION_UNIFORM_NAME           string = "attentuation"
	ATTENTUATION_CONSTANT_UNIFORM_NAME  string = "constant"
	ATTENTUATION_LINEAR_UNIFORM_NAME    string = "linear"
	ATTENTUATION_QUADRATIC_UNIFORM_NAME string = "quadratic"
	NUM_POINT_LIGHTS_UNIFORM_NAME       string = "numPointLights"
	NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME string = "numDirectionalLights"
	NUM_SPOT_LIGHTS_UNIFORM_NAME        string = "numSpotLights"
)

type Attentuation struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

func (a Attentuation) SetUniforms(s Shader, variableName string, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformF(variableName+"["+strconv.Itoa(int(arrayIndex))+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_CONSTANT_UNIFORM_NAME, a.Constant); err != nil {
		return err
	}
	if err = s.SetUniformF(variableName+"["+strconv.Itoa(int(arrayIndex))+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_LINEAR_UNIFORM_NAME, a.Linear); err != nil {
		return err
	}
	if err = s.SetUniformF(variableName+"["+strconv.Itoa(int(arrayIndex))+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_QUADRATIC_UNIFORM_NAME, a.Quadratic); err != nil {
		return err
	}

	return nil
}

type PointLight struct {
	Position mgl32.Vec3

	DiffuseColor  color.Color
	SpecularColor color.Color

	Attentuation
}

func (pl PointLight) SetUniforms(s Shader, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+POSITION_UNIFORM_NAME, pl.Position); err != nil {
		return err
	}
	if err = s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, colorToVec3(pl.DiffuseColor)); err != nil {
		return err
	}
	if err = s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, colorToVec3(pl.SpecularColor)); err != nil {
		return err
	}
	if err = pl.Attentuation.SetUniforms(s, POINT_LIGHTS_UNIFORM_NAME, arrayIndex); err != nil {
		return err
	}
	return nil
}

type DirectionalLight struct {
	Direction mgl32.Vec3

	DiffuseColor  color.Color
	SpecularColor color.Color
}

func (pl DirectionalLight) SetUniforms(s Shader, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIRECTION_UNIFORM_NAME, pl.Direction); err != nil {
		return err
	}
	if err = s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, colorToVec3(pl.DiffuseColor)); err != nil {
		return err
	}
	if err = s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, colorToVec3(pl.SpecularColor)); err != nil {
		return err
	}

	return nil
}

type SpotLight struct {
	Position  mgl32.Vec3
	Direction mgl32.Vec3

	DiffuseColor  color.Color
	SpecularColor color.Color

	InnerCutOff float32
	OuterCutOff float32

	Attentuation
}

func (pl SpotLight) SetUniforms(s Shader, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+POSITION_UNIFORM_NAME, pl.Position); err != nil {
		return err
	}
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIRECTION_UNIFORM_NAME, pl.Direction); err != nil {
		return err
	}
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, colorToVec3(pl.DiffuseColor)); err != nil {
		return err
	}
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, colorToVec3(pl.SpecularColor)); err != nil {
		return err
	}
	if err = s.SetUniformF(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+INNERCUTOFF_UNIFORM_NAME, pl.InnerCutOff); err != nil {
		return err
	}
	if err = s.SetUniformF(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+OUTERCUTOFF_UNIFORM_NAME, pl.OuterCutOff); err != nil {
		return err
	}
	if err = pl.Attentuation.SetUniforms(s, SPOT_LIGHTS_UNIFORM_NAME, arrayIndex); err != nil {
		return err
	}
	return nil
}

type LightCollection struct {
	AmbientLight color.Color

	PointLights       []*PointLight
	DirectionalLights []*DirectionalLight
	SpotLights        []*SpotLight
}

func (this *LightCollection) AddPointLight(pl *PointLight) {
	this.PointLights = append(this.PointLights, pl)
}

func (this *LightCollection) AddDirectionalLight(pl *DirectionalLight) {
	this.DirectionalLights = append(this.DirectionalLights, pl)
}

func (this *LightCollection) AddSpotLight(pl *SpotLight) {
	this.SpotLights = append(this.SpotLights, pl)
}

type LightManager struct {
	lightCollections       []LightCollection
	CurrentLightCollection int32
}

func (this *LightManager) Init() {
	this.lightCollections = make([]LightCollection, 1)
	this.CurrentLightCollection = 0
}

func (this *LightManager) SetAmbientLight(color color.Color, lightCollectionIndex uint32) {
	if len(this.lightCollections) == 0 {
		this.lightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.lightCollections)-1) < lightCollectionIndex {
		this.lightCollections = append(this.lightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.lightCollections)-1))...)
	}
	this.lightCollections[lightCollectionIndex].AmbientLight = color
}

func (this *LightManager) AddPointLight(pl *PointLight, lightCollectionIndex uint32) {
	if len(this.lightCollections) == 0 {
		this.lightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.lightCollections)-1) < lightCollectionIndex {
		this.lightCollections = append(this.lightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.lightCollections)-1))...)
	}
	this.lightCollections[lightCollectionIndex].AddPointLight(pl)
}

func (this *LightManager) AddDirectionalLight(pl *DirectionalLight, lightCollectionIndex uint32) {
	if len(this.lightCollections) == 0 {
		this.lightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.lightCollections)-1) < lightCollectionIndex {
		this.lightCollections = append(this.lightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.lightCollections)-1))...)
	}
	this.lightCollections[lightCollectionIndex].AddDirectionalLight(pl)
}

func (this *LightManager) AddSpotLight(pl *SpotLight, lightCollectionIndex uint32) {
	if len(this.lightCollections) == 0 {
		this.lightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.lightCollections)-1) < lightCollectionIndex {
		this.lightCollections = append(this.lightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.lightCollections)-1))...)
	}
	this.lightCollections[lightCollectionIndex].AddSpotLight(pl)
}

var LightMgr LightManager
