package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	VERTEX       uint8 = 0
	FRAGMENT     uint8 = 1
	GEOMETRY     uint8 = 2
	TESSELLETION uint8 = 3
	EVELUATION   uint8 = 4
	COMPUTE      uint8 = 5
)

type Shader interface {
	AddShader(shader_type uint8, src string) error
	Link() error
	Setup() error
	Terminate()
	Use()
	Unuse()
	SetUniformV2(name string, value mgl32.Vec2) error
	SetUniformV3(name string, value mgl32.Vec3) error
	SetUniformV4(name string, value mgl32.Vec4) error
	SetUniformF(name string, value float32) error
	SetUniformI(name string, value int32) error
	SetUniformUI(name string, value uint32) error
	SetUniformB(name string, value uint8) error
	SetUniformM2(name string, value mgl32.Mat2) error
	SetUniformM3(name string, value mgl32.Mat3) error
	SetUniformM4(name string, value mgl32.Mat4) error
	SetUniformMaterial(mat Material) error
	SetUniformLights(lightCollectionIndex int32) error
	GetName() string
	AddAttribute(name string, location uint32)
}
