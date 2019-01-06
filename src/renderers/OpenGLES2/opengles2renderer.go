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

type OpenGLES2Renderer struct {
	CurrentTextureUnit uint32

	availableFunctions map[string]bool
	backBufferMesh     *OpenGLES2Mesh2D
	backgroundColor    color.Color
	version            uint8
}

func (this *OpenGLES2Renderer) createBackBufferMesh() {
	this.backBufferMesh = CreateOpenGLES2Mesh2D("BackBufferMesh")

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

func (this *OpenGLES2Renderer) Init() error {
	version := gl.GetString(gl.VERSION)
	versioni := this.GetVersioni()
	if version == "" {
		version = strconv.FormatUint(uint64(versioni), 10)
	}
	gohome.ErrorMgr.Log("Renderer", "OpenGLES2\t", "Version: "+version)
	if versioni < 21 {
		gohome.ErrorMgr.Warning("Renderer", "OpenGLES2", "You don't have a graphics card or your graphics card is not supported! Minimum: OpenGL 2.1")
	}

	this.CurrentTextureUnit = 0

	this.availableFunctions = make(map[string]bool)
	this.gatherAvailableFunctions()

	if !this.HasFunctionAvailable("VERTEX_ID") || !this.HasFunctionAvailable("MULTISAMPLE") {
		this.createBackBufferMesh()
	}

	return nil
}

func (this *OpenGLES2Renderer) AfterInit() {
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.BLEND)
	gl.ClearDepthf(2.0)
	gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)

	gl.FrontFace(gl.CCW)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.DEPTH_TEST)
}

func (this *OpenGLES2Renderer) SetWireFrame(b bool) {
	gohome.ErrorMgr.Warning("Renderer", "OpenGLES2", "WireframeMode does not work")
}

func (this *OpenGLES2Renderer) Terminate() {
	if this.backBufferMesh != nil {
		this.backBufferMesh.Terminate()
	}
}

func (*OpenGLES2Renderer) ClearScreen(c color.Color) {
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

func (*OpenGLES2Renderer) CreateTexture(name string, multiSampled bool) gohome.Texture {
	return CreateOpenGLES2Texture(name)
}

func (*OpenGLES2Renderer) CreateMesh2D(name string) gohome.Mesh2D {
	return CreateOpenGLES2Mesh2D(name)
}

func (*OpenGLES2Renderer) CreateMesh3D(name string) gohome.Mesh3D {
	return CreateOpenGLES2Mesh3D(name)
}

func (*OpenGLES2Renderer) CreateRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) gohome.RenderTexture {
	return CreateOpenGLES2RenderTexture(name, width, height, textures, depthBuffer, multiSampled, shadowMap, cubeMap)
}

func (*OpenGLES2Renderer) CreateCubeMap(name string) gohome.CubeMap {
	return CreateOpenGLES2CubeMap(name)
}

func (*OpenGLES2Renderer) CreateInstancedMesh3D(name string) gohome.InstancedMesh3D {
	return CreateOpenGLES2InstancedMesh3D(name)
}

func (*OpenGLES2Renderer) CreateLines3DInterface(name string) gohome.Lines3DInterface {
	return &OpenGLES2Lines3DInterface{
		Name: name,
	}
}

func (this *OpenGLES2Renderer) LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (gohome.Shader, error) {
	var shader *OpenGLES2Shader
	var err error

	shader, err = CreateOpenGLES2Shader(name)
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

func (this *OpenGLES2Renderer) RenderBackBuffer() {
	this.backBufferMesh.Render()
}

func (this *OpenGLES2Renderer) SetViewport(viewport gohome.Viewport) {
	gl.Viewport(int32(viewport.X), int32(viewport.Y), int32(viewport.Width), int32(viewport.Height))
}

func (this *OpenGLES2Renderer) GetViewport() gohome.Viewport {
	var data [4]int32
	gl.GetIntegerv(gl.VIEWPORT, data[:])

	return gohome.Viewport{
		0,
		int(data[0]), int(data[1]),
		int(data[2]), int(data[3]),
		false,
	}
}

func (this *OpenGLES2Renderer) SetNativeResolution(width, height uint32) {
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
func (this *OpenGLES2Renderer) GetNativeResolution() mgl32.Vec2 {
	return [2]float32{float32(gohome.RenderMgr.BackBufferMS.GetWidth()), float32(gohome.RenderMgr.BackBufferMS.GetHeight())}
}
func (this *OpenGLES2Renderer) OnResize(newWidth, newHeight uint32) {
	gl.Viewport(0, 0, int32(newWidth), int32(newHeight))
}

func (this *OpenGLES2Renderer) PreRender() {
	this.CurrentTextureUnit = 1
}
func (this *OpenGLES2Renderer) AfterRender() {
	this.CurrentTextureUnit = 1
}

func (this *OpenGLES2Renderer) SetBacckFaceCulling(b bool) {
	if b {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (this *OpenGLES2Renderer) GetMaxTextures() int32 {
	var data [1]int32
	gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS, data[:])
	return data[0]
}

func (this *OpenGLES2Renderer) NextTextureUnit() uint32 {
	val := this.CurrentTextureUnit
	this.CurrentTextureUnit++
	return val
}

func (this *OpenGLES2Renderer) DecrementTextureUnit(amount uint32) {
	this.CurrentTextureUnit -= amount
}

func (this *OpenGLES2Renderer) GetVersioni() uint8 {
	return 20
}

func (this *OpenGLES2Renderer) gatherAvailableFunctions() {

}

func (this *OpenGLES2Renderer) HasFunctionAvailable(function string) bool {
	return false
}

func (this *OpenGLES2Renderer) SetBackgroundColor(bgColor color.Color) {
	this.backgroundColor = bgColor
}

func (this *OpenGLES2Renderer) GetBackgroundColor() color.Color {
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

func (this *OpenGLES2Renderer) CreateShape2DInterface(name string) gohome.Shape2DInterface {
	return &OpenGLES2Shape2DInterface{
		Name: name,
	}
}

func (this *OpenGLES2Renderer) SetDepthTesting(b bool) {
	if b {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func (this *OpenGLES2Renderer) GetName() string {
	return "OpenGLES2"
}
