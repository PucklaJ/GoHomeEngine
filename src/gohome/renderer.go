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

	RenderBackBuffer()

	SetBacckFaceCulling(b bool)
	SetDepthTesting(b bool)
	GetMaxTextures() int32
	NextTextureUnit() uint32
	DecrementTextureUnit(amount uint32)
	FilterShaderFiles(name, file, shader_type string) string
	FilterShaderSource(name, source, shader_type string) string
}

var Render Renderer
