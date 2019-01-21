package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
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
	SetUniformLights(lightCollectionIndex int)
	GetName() string
	AddAttribute(name string, location uint32)
}

type NilShader struct {
}

func (*NilShader) AddShader(shader_type uint8, src string) error {
	return nil
}
func (*NilShader) Link() error {
	return nil
}
func (*NilShader) Setup() error {
	return nil
}
func (*NilShader) Terminate() {

}
func (*NilShader) Use() {

}
func (*NilShader) Unuse() {

}
func (*NilShader) SetUniformV2(name string, value mgl32.Vec2) {

}
func (*NilShader) SetUniformV3(name string, value mgl32.Vec3) {

}
func (*NilShader) SetUniformV4(name string, value mgl32.Vec4) {

}
func (*NilShader) SetUniformIV2(name string, value []int32) {

}
func (*NilShader) SetUniformIV3(name string, value []int32) {

}
func (*NilShader) SetUniformIV4(name string, value []int32) {

}
func (*NilShader) SetUniformF(name string, value float32) {

}
func (*NilShader) SetUniformI(name string, value int32) {

}
func (*NilShader) SetUniformUI(name string, value uint32) {

}
func (*NilShader) SetUniformB(name string, value uint8) {

}
func (*NilShader) SetUniformM2(name string, value mgl32.Mat2) {

}
func (*NilShader) SetUniformM3(name string, value mgl32.Mat3) {

}
func (*NilShader) SetUniformM4(name string, value mgl32.Mat4) {

}
func (*NilShader) SetUniformMaterial(mat Material) {

}
func (*NilShader) SetUniformLights(lightCollectionIndex int) {

}
func (*NilShader) GetName() string {
	return ""
}
func (*NilShader) AddAttribute(name string, location uint32) {

}
