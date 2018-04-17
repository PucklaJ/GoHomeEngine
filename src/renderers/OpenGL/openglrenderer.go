package renderer

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/gl/v4.1-core/gl"
	"image/color"
	"log"
)

const (
	GL_MAX_TEXTURE_MAX_ANISOTROPY uint32 = 0x84FF
	GL_TEXTURE_MAX_ANISOTROPY     uint32 = 0x84FE
)

type OpenGLRenderer struct {
	BackBufferVao      uint32
	CurrentTextureUnit uint32
}

func (this *OpenGLRenderer) Init() error {
	if err := gl.Init(); err != nil {
		return err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL Version:", version)

	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.DEPTH_CLAMP)
	gl.ClearDepth(1.0)

	gl.FrontFace(gl.CCW)
	gl.Enable(gl.CULL_FACE)

	gl.GenVertexArrays(1, &this.BackBufferVao)

	this.CurrentTextureUnit = 0

	return nil
}

func (this *OpenGLRenderer) hasExtenstion(name string) bool {
	var numExtensions int32
	gl.GetIntegerv(gl.NUM_EXTENSIONS, &numExtensions)
	for i := 0; i < int(numExtensions); i++ {
		ext := gl.GoStr(gl.GetStringi(gl.EXTENSIONS, uint32(i)))
		if ext == name {
			return true
		}
	}
	return false
}

func (this *OpenGLRenderer) SetWireFrame(b bool) {
	if b {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
}

func (this *OpenGLRenderer) Terminate() {
	gl.DeleteVertexArrays(1, &this.BackBufferVao)
}

func (*OpenGLRenderer) ClearScreen(c color.Color, alpha float32) {
	clearColor := gohome.ColorToVec3(c)
	gl.ClearColor(clearColor[0], clearColor[1], clearColor[2], alpha)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
}

type OpenGLError struct {
	errorString string
}

func (oerr OpenGLError) Error() string {
	return oerr.errorString
}

func (*OpenGLRenderer) CreateTexture(name string, multiSampled bool) gohome.Texture {
	return CreateOpenGLTexture(name, multiSampled)
}

func (*OpenGLRenderer) CreateMesh2D(name string) gohome.Mesh2D {
	return CreateOpenGLMesh2D(name)
}

func (*OpenGLRenderer) CreateMesh3D(name string) gohome.Mesh3D {
	return CreateOpenGLMesh3D(name)
}

func (*OpenGLRenderer) CreateRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) gohome.RenderTexture {
	return CreateOpenGLRenderTexture(name, width, height, textures, depthBuffer, multiSampled, shadowMap, cubeMap)
}

func (*OpenGLRenderer) CreateCubeMap(name string) gohome.CubeMap {
	return CreateOpenGLCubeMap(name)
}

func (*OpenGLRenderer) CreateInstancedMesh3D(name string) gohome.InstancedMesh3D {
	return CreateOpenGLInstancedMesh3D(name)
}

func (*OpenGLRenderer) LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (gohome.Shader, error) {
	var shader *OpenGLShader
	var err error

	shader, err = CreateOpenGLShader(name)
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

func (this *OpenGLRenderer) RenderBackBuffer() {
	gl.BindVertexArray(this.BackBufferVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}

func (this *OpenGLRenderer) SetViewport(viewport gohome.Viewport) {
	gl.Viewport(int32(viewport.X), int32(viewport.Y), int32(viewport.Width), int32(viewport.Height))
}

func (this *OpenGLRenderer) GetViewport() gohome.Viewport {
	var data [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &data[0])

	return gohome.Viewport{
		0,
		int(data[0]), int(data[1]),
		int(data[2]), int(data[3]),
	}
}

func (this *OpenGLRenderer) SetNativeResolution(width, height uint32) {
	gohome.RenderMgr.BackBufferMS.ChangeSize(width, height)
	gohome.RenderMgr.BackBuffer.ChangeSize(width, height)
	gohome.RenderMgr.BackBuffer2D.ChangeSize(width, height)
	gohome.RenderMgr.BackBuffer3D.ChangeSize(width, height)

	gohome.RenderMgr.BackBufferMS.SetFiltering(gohome.FILTERING_LINEAR)
	gohome.RenderMgr.BackBuffer.SetFiltering(gohome.FILTERING_LINEAR)
	gohome.RenderMgr.BackBuffer2D.SetFiltering(gohome.FILTERING_LINEAR)
	gohome.RenderMgr.BackBuffer3D.SetFiltering(gohome.FILTERING_LINEAR)
}
func (this *OpenGLRenderer) GetNativeResolution() (uint32, uint32) {
	return uint32(gohome.RenderMgr.BackBuffer.GetWidth()), uint32(gohome.RenderMgr.BackBuffer.GetHeight())
}
func (this *OpenGLRenderer) OnResize(newWidth, newHeight uint32) {
	gl.Viewport(0, 0, int32(newWidth), int32(newHeight))
}

func (this *OpenGLRenderer) PreRender() {
	this.CurrentTextureUnit = 1
}
func (this *OpenGLRenderer) AfterRender() {
	this.CurrentTextureUnit = 1
}

func (this *OpenGLRenderer) SetBacckFaceCulling(b bool) {
	if b {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (this *OpenGLRenderer) GetMaxTextures() int32 {
	var data int32
	gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS, &data)
	return data
}

func (this *OpenGLRenderer) NextTextureUnit() uint32 {
	val := this.CurrentTextureUnit
	this.CurrentTextureUnit++
	return val
}

func (this *OpenGLRenderer) DecrementTextureUnit(amount uint32) {
	this.CurrentTextureUnit -= amount
}
