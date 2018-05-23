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
	SetUniformV2(name string, value mgl32.Vec2)
	SetUniformV3(name string, value mgl32.Vec3)
	SetUniformV4(name string, value mgl32.Vec4)
	SetUniformIV2(name string, value []int32)
	SetUniformIV3(name string, value []int32)
	SetUniformIV4(name string, value []int32)
	SetUniformF(name string, value float32)
	SetUniformI(name string, value int32)
	SetUniformUI(name string, value uint32)
	SetUniformB(name string, value uint8)
	SetUniformM2(name string, value mgl32.Mat2)
	SetUniformM3(name string, value mgl32.Mat3)
	SetUniformM4(name string, value mgl32.Mat4)
	SetUniformMaterial(mat Material)
	SetUniformLights(lightCollectionIndex int32)
	GetName() string
	AddAttribute(name string, location uint32)
}
