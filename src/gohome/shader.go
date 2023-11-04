package gohome

import (
	"github.com/PucklaJ/mathgl/mgl32"
)

// The different parts of a shader
const (
	VERTEX       uint8 = 0
	FRAGMENT     uint8 = 1
	GEOMETRY     uint8 = 2
	TESSELLETION uint8 = 3
	EVELUATION   uint8 = 4
	COMPUTE      uint8 = 5
)

// A shader controls how a mesh is rendered
type Shader interface {
	// Adds a shader of shader_type with its source code to this shader
	AddShader(shader_type uint8, src string) error
	// Links all shaders
	Link() error
	// Sets all values up
	Setup() error
	// Cleans everything up
	Terminate()
	// Use this shader for the following draw calls
	Use()
	// Don't use this shader anymore
	Unuse()
	// Sets the value of a uniform Vec2
	SetUniformV2(name string, value mgl32.Vec2)
	// Sets the value of a uniform Vec3
	SetUniformV3(name string, value mgl32.Vec3)
	// Sets the value of a uniform Vec4
	SetUniformV4(name string, value mgl32.Vec4)
	// Sets the value of a uniform ivec2
	SetUniformIV2(name string, value []int32)
	// Sets the value of a uniform ivec3
	SetUniformIV3(name string, value []int32)
	// Sets the value of a uniform ivec4
	SetUniformIV4(name string, value []int32)
	// Sets the value of a uniform float
	SetUniformF(name string, value float32)
	// Sets the value of a uniform int
	SetUniformI(name string, value int32)
	// Sets the value of a uniform unsigned int
	SetUniformUI(name string, value uint32)
	// Sets the value of a uniform bool
	SetUniformB(name string, value uint8)
	// Sets the value of a uniform Mat2
	SetUniformM2(name string, value mgl32.Mat2)
	// Sets the value of a uniform Mat3
	SetUniformM3(name string, value mgl32.Mat3)
	// Sets the value of a uniform Mat4
	SetUniformM4(name string, value mgl32.Mat4)
	// Sets the value of a uniform material
	SetUniformMaterial(mat Material)
	// Sets the value of all uniforms of lights
	SetUniformLights(lightCollectionIndex int)
	// Returns the name of this shader
	GetName() string
	// Adds a vertex attribute
	AddAttribute(name string, location uint32)
}

// An implementation of Shader that does nothing
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
