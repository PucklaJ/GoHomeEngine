package gohome

import (
	"image/color"
)

type Renderer interface {
	Init() error
	Terminate()
	ClearScreen(c color.Color)
	LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (Shader, error)
	CreateTexture(name string, multiSampled bool) Texture
	CreateMesh2D(name string) Mesh2D
	CreateMesh3D(name string) Mesh3D
	CreateRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) RenderTexture
	CreateCubeMap(name string) CubeMap
	CreateInstancedMesh3D(name string) InstancedMesh3D
	SetWireFrame(b bool)
	SetViewport(viewport Viewport)
	GetViewport() Viewport
	SetNativeResolution(width, height uint32)
	GetNativeResolution() (uint32, uint32)
	OnResize(newWidth, newHeight uint32)
	PreRender()
	AfterRender()

	RenderBackBuffer()

	SetBacckFaceCulling(b bool)
	GetMaxTextures() int32
	NextTextureUnit() uint32
	DecrementTextureUnit(amount uint32)
	FilterShaderFiles(name, file, shader_type string) string
}

var Render Renderer
