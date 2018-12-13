package renderer

import (
	"image/color"
	"strconv"

	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles2"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

const (
	GL_MAX_TEXTURE_MAX_ANISOTROPY uint32 = 0x84FF
	GL_TEXTURE_MAX_ANISOTROPY     uint32 = 0x84FE
)

type OpenGLESRenderer struct {
	CurrentTextureUnit uint32

	availableFunctions map[string]bool
	backBufferMesh     *OpenGLESMesh2D
	backgroundColor    color.Color
	version            uint8
}

func (this *OpenGLESRenderer) createBackBufferMesh() {
	this.backBufferMesh = CreateOpenGLESMesh2D("BackBufferMesh")

	var vertices []gohome.Mesh2DVertex = make([]gohome.Mesh2DVertex, 4)
	var indices []uint32 = make([]uint32, 6)

	vertices[0].Vertex(-1.0, -1.0)
	vertices[1].Vertex(1.0, -1.0)
	vertices[2].Vertex(1.0, 1.0)
	vertices[3].Vertex(-1.0, 1.0)

	vertices[0].TexCoord(0.0, 0.0)
	vertices[1].TexCoord(1.0, 0.0)
	vertices[2].TexCoord(1.0, 1.0)
	vertices[3].TexCoord(0.0, 1.0)

	indices[0] = 0
	indices[1] = 1
	indices[2] = 2
	indices[3] = 2
	indices[4] = 3
	indices[5] = 0

	this.backBufferMesh.AddVertices(vertices, indices)
	this.backBufferMesh.Load()
}

func (this *OpenGLESRenderer) Init() error {
	version := gl.GetString(gl.VERSION)
	versioni := this.GetVersioni()
	if version == "" {
		version = strconv.FormatUint(uint64(versioni), 10)
	}
	gohome.ErrorMgr.Log("Renderer", "OpenGLES\t", "Version: "+version)
	if versioni < 21 {
		gohome.ErrorMgr.Warning("Renderer", "OpenGLES", "You don't have a graphics card or your graphics card is not supported! Minimum: OpenGL 2.1")
	}

	this.CurrentTextureUnit = 0

	this.availableFunctions = make(map[string]bool)
	this.gatherAvailableFunctions()

	if !this.HasFunctionAvailable("VERTEX_ID") || !this.HasFunctionAvailable("MULTISAMPLE") {
		this.createBackBufferMesh()
	}

	return nil
}

func (this *OpenGLESRenderer) AfterInit() {
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.BLEND)
	gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)

	gl.FrontFace(gl.CCW)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.DEPTH_TEST)
}

func (this *OpenGLESRenderer) SetWireFrame(b bool) {
	gohome.ErrorMgr.Error("Renderer", "OpenGLES", "WireframeMode does not work")
}

func (this *OpenGLESRenderer) Terminate() {
	if this.backBufferMesh != nil {
		this.backBufferMesh.Terminate()
	}
}

func (*OpenGLESRenderer) ClearScreen(c color.Color) {
	clearColor := gohome.ColorToVec4(c)
	gl.ClearColor(clearColor[0], clearColor[1], clearColor[2], clearColor[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
}

type OpenGLError struct {
	errorString string
}

func (oerr OpenGLError) Error() string {
	return oerr.errorString
}

func (*OpenGLESRenderer) CreateTexture(name string, multiSampled bool) gohome.Texture {
	return CreateOpenGLESTexture(name)
}

func (*OpenGLESRenderer) CreateMesh2D(name string) gohome.Mesh2D {
	return CreateOpenGLESMesh2D(name)
}

func (*OpenGLESRenderer) CreateMesh3D(name string) gohome.Mesh3D {
	return CreateOpenGLESMesh3D(name)
}

func (*OpenGLESRenderer) CreateRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) gohome.RenderTexture {
	return CreateOpenGLESRenderTexture(name, width, height, textures, depthBuffer, multiSampled, shadowMap, cubeMap)
}

func (*OpenGLESRenderer) CreateCubeMap(name string) gohome.CubeMap {
	return CreateOpenGLESCubeMap(name)
}

func (*OpenGLESRenderer) CreateInstancedMesh3D(name string) gohome.InstancedMesh3D {
	return CreateOpenGLESInstancedMesh3D(name)
}

func (*OpenGLESRenderer) CreateLines3DInterface(name string) gohome.Lines3DInterface {
	return &OpenGLESLines3DInterface{
		Name: name,
	}
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

func (this *OpenGLESRenderer) RenderBackBuffer() {
	this.backBufferMesh.Render()
}

func (this *OpenGLESRenderer) SetViewport(viewport gohome.Viewport) {
	gl.Viewport(int32(viewport.X), int32(viewport.Y), int32(viewport.Width), int32(viewport.Height))
}

func (this *OpenGLESRenderer) GetViewport() gohome.Viewport {
	var data [4]int32
	gl.GetIntegerv(gl.VIEWPORT, data[:])

	return gohome.Viewport{
		0,
		int(data[0]), int(data[1]),
		int(data[2]), int(data[3]),
		false,
	}
}

func (this *OpenGLESRenderer) SetNativeResolution(width, height uint32) {
	previous := gohome.Viewport{
		X:      0,
		Y:      0,
		Width:  gohome.RenderMgr.BackBufferMS.GetWidth(),
		Height: gohome.RenderMgr.BackBufferMS.GetHeight(),
	}

	gohome.RenderMgr.BackBufferMS.ChangeSize(width, height)
	gohome.RenderMgr.BackBuffer2D.ChangeSize(width, height)
	gohome.RenderMgr.BackBuffer3D.ChangeSize(width, height)

	gohome.RenderMgr.BackBufferMS.SetFiltering(gohome.FILTERING_LINEAR)
	gohome.RenderMgr.BackBuffer2D.SetFiltering(gohome.FILTERING_LINEAR)
	gohome.RenderMgr.BackBuffer3D.SetFiltering(gohome.FILTERING_LINEAR)

	current := gohome.Viewport{
		X:      0,
		Y:      0,
		Width:  gohome.RenderMgr.BackBufferMS.GetWidth(),
		Height: gohome.RenderMgr.BackBufferMS.GetHeight(),
	}

	gohome.RenderMgr.UpdateViewports(current, previous)
}
func (this *OpenGLESRenderer) GetNativeResolution() mgl32.Vec2 {
	return [2]float32{float32(gohome.RenderMgr.BackBufferMS.GetWidth()), float32(gohome.RenderMgr.BackBufferMS.GetHeight())}
}
func (this *OpenGLESRenderer) OnResize(newWidth, newHeight uint32) {
	gl.Viewport(0, 0, int32(newWidth), int32(newHeight))
}

func (this *OpenGLESRenderer) PreRender() {
	this.CurrentTextureUnit = 1
}
func (this *OpenGLESRenderer) AfterRender() {
	this.CurrentTextureUnit = 1
}

func (this *OpenGLESRenderer) SetBacckFaceCulling(b bool) {
	if b {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (this *OpenGLESRenderer) GetMaxTextures() int32 {
	var data [1]int32
	gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS, data[:])
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

func (this *OpenGLESRenderer) GetVersioni() uint8 {
	return 20
}

func (this *OpenGLESRenderer) gatherAvailableFunctions() {
	combined := this.GetVersioni()
	this.version = uint8(combined)
	if combined >= 30 {
		this.availableFunctions["VERTEX_ID"] = true
		this.availableFunctions["VERTEX_ARRAY"] = true
		this.availableFunctions["BLIT_FRAMEBUFFER"] = true
	}
	if combined >= 31 {
		this.availableFunctions["INSTANCED"] = true
		this.availableFunctions["INDIRECT"] = true
		this.availableFunctions["MULTISAMPLE"] = true
		this.availableFunctions["FRAMEBUFFER_TEXTURE"] = true
	}
	if combined >= 32 {
		this.availableFunctions["GEOMETRY_SHADER"] = true
	}
}

func (this *OpenGLESRenderer) HasFunctionAvailable(function string) bool {
	v, ok := this.availableFunctions[function]
	return ok && v
}

func (this *OpenGLESRenderer) FilterShaderFiles(name, file, shader_type string) string {
	if name == "BackBufferShader" {
		if !this.HasFunctionAvailable("MULTISAMPLE") {
			if shader_type == "Vertex File" {
				file = "backBufferShaderNoMSVert.glsl"
			} else if shader_type == "Fragment File" {
				file = "backBufferShaderNoMSFrag.glsl"
			}
		}
	}

	return file
}

func (this *OpenGLESRenderer) FilterShaderSource(name, source, shader_type string) string {
	if name == "BackBufferShader" {
		if !this.HasFunctionAvailable("MULTISAMPLE") {
			if shader_type == "Vertex File" {
				source = gohome.BACKBUFFER_NOMS_SHADER_VERTEX_SOURCE_OPENGL
			} else if shader_type == "Fragment File" {
				source = gohome.BACKBUFFER_NOMS_SHADER_FRAGMENT_SOURCE_OPENGL
			}
		}
	}

	return source
}

func (this *OpenGLESRenderer) SetBackgroundColor(bgColor color.Color) {
	this.backgroundColor = bgColor
}

func (this *OpenGLESRenderer) GetBackgroundColor() color.Color {
	return this.backgroundColor
}

func handleOpenGLError(tag, objectName, errorPrefix string) {
	err := gl.GetError()
	if err != gl.NO_ERROR {
		var errString string
		switch err {
		case gl.INVALID_OPERATION:
			errString = "INVALID_OPERATION"
		case gl.INVALID_VALUE:
			errString = "INVALID_VALUE"
		case gl.INVALID_ENUM:
			errString = "INVALID_ENUM"
		case gl.OUT_OF_MEMORY:
			errString = "OUT_OF_MEMORY"
		case gl.INVALID_FRAMEBUFFER_OPERATION:
			errString = "INVALID_FRAMEBUFFER_OPERATION"
		case 0x8031:
			errString = "TABLE_TOO_LARGE"
		default:
			errString = strconv.Itoa(int(err))
		}
		gohome.ErrorMgr.Error(tag, objectName, errorPrefix+" ErrorCode: "+errString)
	}
}

func (this *OpenGLESRenderer) CreateShape2DInterface(name string) gohome.Shape2DInterface {
	return &OpenGLESShape2DInterface{
		Name: name,
	}
}

func (this *OpenGLESRenderer) SetDepthTesting(b bool) {
	if b {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}
