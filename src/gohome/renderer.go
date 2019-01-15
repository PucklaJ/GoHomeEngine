package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"image/color"
)

type Renderer interface {
	Init() error
	AfterInit()
	Terminate()
	ClearScreen(c color.Color)
	LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (Shader, error)
	CreateTexture(name string, multiSampled bool) Texture
	CreateMesh2D(name string) Mesh2D
	CreateMesh3D(name string) Mesh3D
	CreateRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) RenderTexture
	CreateCubeMap(name string) CubeMap
	CreateInstancedMesh3D(name string) InstancedMesh3D
	CreateLines3DInterface(name string) Lines3DInterface
	CreateShape2DInterface(name string) Shape2DInterface
	SetWireFrame(b bool)
	SetViewport(viewport Viewport)
	GetViewport() Viewport
	SetNativeResolution(width, height uint32)
	GetNativeResolution() mgl32.Vec2
	OnResize(newWidth, newHeight uint32)
	PreRender()
	AfterRender()
	SetBackgroundColor(bgColor color.Color)
	GetBackgroundColor() color.Color
	GetName() string

	RenderBackBuffer()

	SetBacckFaceCulling(b bool)
	SetDepthTesting(b bool)
	GetMaxTextures() int32
	NextTextureUnit() uint32
	DecrementTextureUnit(amount uint32)
	HasFunctionAvailable(name string) bool

	InstancedMesh3DFromLoadedMesh3D(mesh Mesh3D) InstancedMesh3D
}

var Render Renderer

type NilRenderer struct {
}

func (*NilRenderer) Init() error {
	return nil
}
func (*NilRenderer) AfterInit() {

}
func (*NilRenderer) Terminate() {

}
func (*NilRenderer) ClearScreen(c color.Color) {

}
func (*NilRenderer) LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (Shader, error) {
	return &NilShader{}, nil
}
func (*NilRenderer) CreateTexture(name string, multiSampled bool) Texture {
	return &NilTexture{}
}
func (*NilRenderer) CreateMesh2D(name string) Mesh2D {
	return &NilMesh2D{}
}
func (*NilRenderer) CreateMesh3D(name string) Mesh3D {
	return &NilMesh3D{}
}
func (*NilRenderer) CreateRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) RenderTexture {
	return &NilRenderTexture{}
}
func (*NilRenderer) CreateCubeMap(name string) CubeMap {
	return &NilCubeMap{}
}
func (*NilRenderer) CreateInstancedMesh3D(name string) InstancedMesh3D {
	return &NilInstancedMesh3D{}
}
func (*NilRenderer) CreateLines3DInterface(name string) Lines3DInterface {
	return &NilLines3DInterface{}
}
func (*NilRenderer) CreateShape2DInterface(name string) Shape2DInterface {
	return &NilShape2DInterface{}
}
func (*NilRenderer) SetWireFrame(b bool) {

}
func (*NilRenderer) SetViewport(viewport Viewport) {

}
func (*NilRenderer) GetViewport() Viewport {
	return Viewport{
		0, 0, 0, 0, 0, false,
	}
}
func (*NilRenderer) SetNativeResolution(width, height uint32) {

}
func (*NilRenderer) GetNativeResolution() mgl32.Vec2 {
	return [2]float32{0.0, 0.0}
}
func (*NilRenderer) OnResize(newWidth, newHeight uint32) {

}
func (*NilRenderer) PreRender() {

}
func (*NilRenderer) AfterRender() {

}
func (*NilRenderer) SetBackgroundColor(bgColor color.Color) {

}
func (*NilRenderer) GetBackgroundColor() color.Color {
	return nil
}
func (*NilRenderer) GetName() string {
	return ""
}

func (*NilRenderer) RenderBackBuffer() {

}

func (*NilRenderer) SetBacckFaceCulling(b bool) {

}
func (*NilRenderer) SetDepthTesting(b bool) {

}
func (*NilRenderer) GetMaxTextures() int32 {
	return 0
}
func (*NilRenderer) NextTextureUnit() uint32 {
	return 0
}
func (*NilRenderer) DecrementTextureUnit(amount uint32) {

}
func (*NilRenderer) FilterShaderFiles(name, file, shader_type string) string {
	return file
}
func (*NilRenderer) FilterShaderSource(name, source, shader_type string) string {
	return source
}
func (*NilRenderer) HasFunctionAvailable(name string) bool {
	return false
}

func (*NilRenderer) InstancedMesh3DFromLoadedMesh3D(mesh Mesh3D) InstancedMesh3D {
	return &NilInstancedMesh3D{}
}
