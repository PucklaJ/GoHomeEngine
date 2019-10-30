package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"image/color"
)

// This interface handles every low level rendering operation
type Renderer interface {
	// Initialises the renderer
	Init() error
	// Gets called after the initialisation of the engine
	AfterInit()
	// Cleans everything up
	Terminate()
	// Clears the screen with the given color
	ClearScreen(c color.Color)
	// Loads a shader given the contents of shaders
	LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (Shader, error)
	// Creates a texture without data
	CreateTexture(name string, multiSampled bool) Texture
	// Creates a Mesh2D
	CreateMesh2D(name string) Mesh2D
	// Creates a Mesh3D
	CreateMesh3D(name string) Mesh3D
	// Creates a RenderTexture from the given parameters
	CreateRenderTexture(name string, width, height, textures int, depthBuffer, multiSampled, shadowMap, cubeMap bool) RenderTexture
	// Creates a cube map
	CreateCubeMap(name string) CubeMap
	// Creates an instanced mesh 3d
	CreateInstancedMesh3D(name string) InstancedMesh3D
	// Creates a shape 3d interface
	CreateShape3DInterface(name string) Shape3DInterface
	// Creates a shape 2d interface
	CreateShape2DInterface(name string) Shape2DInterface
	// Enables or disables wire frame render mode
	SetWireFrame(b bool)
	// Sets the current viewport for the GPU
	SetViewport(viewport Viewport)
	// Returns the current viewport of the GPU
	GetViewport() Viewport
	// Sets the resolution of the back buffer
	SetNativeResolution(width, height int)
	// Returns the resolution of the back buffer
	GetNativeResolution() mgl32.Vec2
	// Gets called when the window resizes
	OnResize(newWidth, newHeight int)
	// Gets called before rendering a RenderObject
	PreRender()
	// Gets called after rendering a RenderObject
	AfterRender()
	// Sets the clear color
	SetBackgroundColor(bgColor color.Color)
	// Returns the clear color
	GetBackgroundColor() color.Color
	// Returns the name of the renderer
	GetName() string

	// Calls the draw methods of the back buffer
	RenderBackBuffer()

	// Enable or disable back face culling
	SetBacckFaceCulling(b bool)
	// Enable or disable depth testing
	SetDepthTesting(b bool)
	// Returns the number maximum textures supported by the GPU
	GetMaxTextures() int
	// Increments the texture unit used for textures
	NextTextureUnit() uint32
	// Decrements the texture unit used for textures
	DecrementTextureUnit(amount uint32)
	// Returns wether the given function is supported by the hardware
	HasFunctionAvailable(name string) bool

	// Returns a InstancedMesh3D created from an already loaded Mesh3D
	InstancedMesh3DFromLoadedMesh3D(mesh Mesh3D) InstancedMesh3D
}

// The Renderer that should be used for everything
var Render Renderer

// An implementation of Renderer that does nothing
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
func (*NilRenderer) CreateRenderTexture(name string, width, height, textures int, depthBuffer, multiSampled, shadowMap, cubeMap bool) RenderTexture {
	return &NilRenderTexture{}
}
func (*NilRenderer) CreateCubeMap(name string) CubeMap {
	return &NilCubeMap{}
}
func (*NilRenderer) CreateInstancedMesh3D(name string) InstancedMesh3D {
	return &NilInstancedMesh3D{}
}
func (*NilRenderer) CreateShape3DInterface(name string) Shape3DInterface {
	return &NilShape3DInterface{}
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
func (*NilRenderer) SetNativeResolution(width, height int) {

}
func (*NilRenderer) GetNativeResolution() mgl32.Vec2 {
	return [2]float32{0.0, 0.0}
}
func (*NilRenderer) OnResize(newWidth, newHeight int) {

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
