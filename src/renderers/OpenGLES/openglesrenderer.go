package renderer

import (
	"fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/mobile/gl"
	"image/color"
)

type OpenGLESError struct {
	errorString string
}

func (this *OpenGLESError) Error() string {
	return this.errorString
}

type OpenGLESRenderer struct {
	gles               gl.Context
	CurrentTextureUnit uint32
	backBufferVao      gl.VertexArray
}

func (this *OpenGLESRenderer) Init() error {
	this.CurrentTextureUnit = 0

	fmt.Println("Version:", gl.Version())

	this.backBufferVao = this.gles.CreateVertexArray()

	return nil
}
func (this *OpenGLESRenderer) Terminate() {

}
func (this *OpenGLESRenderer) ClearScreen(c color.Color, alpha float32) {
	col := gohome.ColorToVec4(c)
	this.gles.ClearColor(col.X(), col.Y(), col.Z(), col.W())
	this.gles.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
}
func (this *OpenGLESRenderer) LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (gohome.Shader, error) {
	var shader *OpenGLESShader
	var err error

	shader, err = CreateOpenGLESShader(name)
	if err != nil {
		return nil, err
	}

	if vertex_contents != "" {
		err = shader.AddShader(gohome.VERTEX, vertex_contents)
		if err != nil {
			return nil, err
		}
	}
	if fragment_contents != "" {
		err = shader.AddShader(gohome.FRAGMENT, fragment_contents)
		if err != nil {
			return nil, err
		}
	}
	if geometry_contents != "" {
		err = shader.AddShader(gohome.GEOMETRY, geometry_contents)
		if err != nil {
			return nil, err
		}
	}
	if tesselletion_control_contents != "" {
		err = shader.AddShader(gohome.TESSELLETION, tesselletion_control_contents)
		if err != nil {
			return nil, err
		}
	}
	if eveluation_contents != "" {
		err = shader.AddShader(gohome.EVELUATION, eveluation_contents)
		if err != nil {
			return nil, err
		}
	}
	if compute_contents != "" {
		err = shader.AddShader(gohome.COMPUTE, compute_contents)
		if err != nil {
			return nil, err
		}
	}

	err = shader.Link()
	if err != nil {
		return nil, err
	}
	err = shader.Setup()
	if err != nil {
		return nil, err
	}

	return shader, nil
}
func (this *OpenGLESRenderer) CreateTexture(name string, multiSampled bool) gohome.Texture {
	return nil
}
func (this *OpenGLESRenderer) CreateMesh2D(name string) gohome.Mesh2D {
	return nil
}
func (this *OpenGLESRenderer) CreateMesh3D(name string) gohome.Mesh3D {
	return nil
}
func (this *OpenGLESRenderer) CreateRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) gohome.RenderTexture {
	return nil
}
func (this *OpenGLESRenderer) CreateCubeMap(name string) gohome.CubeMap {
	return nil
}
func (this *OpenGLESRenderer) CreateInstancedMesh3D(name string) gohome.InstancedMesh3D {
	return nil
}
func (this *OpenGLESRenderer) SetWireFrame(b bool) {

}
func (this *OpenGLESRenderer) SetViewport(viewport gohome.Viewport) {

}
func (this *OpenGLESRenderer) GetViewport() gohome.Viewport {
	var data [4]int32

	this.gles.GetIntegerv(data[:], gl.VIEWPORT)

	return gohome.Viewport{
		X:      int(data[0]),
		Y:      int(data[1]),
		Width:  int(data[2]),
		Height: int(data[3]),
	}
}
func (this *OpenGLESRenderer) SetNativeResolution(width, height uint32) {

}
func (this *OpenGLESRenderer) GetNativeResolution() (uint32, uint32) {
	return 0, 0
}
func (this *OpenGLESRenderer) OnResize(newWidth, newHeight uint32) {

}
func (this *OpenGLESRenderer) PreRender() {
	this.CurrentTextureUnit = 1
}
func (this *OpenGLESRenderer) AfterRender() {
	this.CurrentTextureUnit = 1
}
func (this *OpenGLESRenderer) RenderBackBuffer() {
	// this.gles.BindVertexArray(this.backBufferVao)
	// this.gles.DrawArrays(gl.TRIANGLES, 0, 6)
	// this.gles.BindVertexArray(gl.VertexArray{0})
}
func (this *OpenGLESRenderer) SetBacckFaceCulling(b bool) {
	if b {
		this.gles.Enable(gl.CULL_FACE)
	} else {
		this.gles.Disable(gl.CULL_FACE)
	}
}
func (this *OpenGLESRenderer) GetMaxTextures() int32 {
	var data [1]int32
	this.gles.GetIntegerv(data[:], gl.MAX_TEXTURE_IMAGE_UNITS)
	return data[0]
}
func (this *OpenGLESRenderer) NextTextureUnit() uint32 {
	val := this.CurrentTextureUnit
	this.CurrentTextureUnit++
	return val
}
func (this *OpenGLESRenderer) DecrementTextureUnit(amount uint32) {
	this.CurrentTextureUnit -= amount
}

func (this *OpenGLESRenderer) SetOpenGLESContex(context gl.Context) {
	this.gles = context
}
