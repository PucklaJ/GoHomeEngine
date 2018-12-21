package renderer

import (
	"image/color"
	"strconv"

	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"github.com/go-gl/gl/all-core/gl"
)

const (
	GL_MAX_TEXTURE_MAX_ANISOTROPY uint32 = 0x84FF
	GL_TEXTURE_MAX_ANISOTROPY     uint32 = 0x84FE
)

type OpenGLRenderer struct {
	BackBufferVao      uint32
	CurrentTextureUnit uint32

	availableFunctions map[string]bool
	backBufferMesh     *OpenGLMesh2D
	backgroundColor    color.Color
	version            uint8
}

func (this *OpenGLRenderer) createBackBufferMesh() {
	this.backBufferMesh = CreateOpenGLMesh2D("BackBufferMesh")

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

func (this *OpenGLRenderer) Init() error {
	if err := gl.Init(); err != nil {
		return err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	versioni := this.GetVersioni()
	if version == "" {
		version = strconv.FormatUint(uint64(versioni), 10)
	}
	gohome.ErrorMgr.Log("Renderer", "OpenGL\t", "Version: "+version)
	if versioni < 21 {
		gohome.ErrorMgr.Warning("Renderer", "OpenGL", "You don't have a graphics card or your graphics card is not supported! Minimum: OpenGL 2.1")
	}

	this.CurrentTextureUnit = 0

	this.availableFunctions = make(map[string]bool)
	this.gatherAvailableFunctions()

	if !this.HasFunctionAvailable("VERTEX_ID") || !this.HasFunctionAvailable("MULTISAMPLE") {
		this.createBackBufferMesh()
	} else {
		gl.GenVertexArrays(1, &this.BackBufferVao)
	}

	return nil
}

func (this *OpenGLRenderer) AfterInit() {
	gl.Enable(gl.MULTISAMPLE)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.DEPTH_CLAMP)
	gl.ClearDepth(2.0)
	gl.Enable(gl.BLEND)
	gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)

	gl.FrontFace(gl.CCW)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.DEPTH_TEST)
}

func (this *OpenGLRenderer) HasExtension(name string) bool {
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
	if this.BackBufferVao != 0 {
		gl.DeleteVertexArrays(1, &this.BackBufferVao)
	}
	if this.backBufferMesh != nil {
		this.backBufferMesh.Terminate()
	}
}

func (*OpenGLRenderer) ClearScreen(c color.Color) {
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

func (*OpenGLRenderer) CreateLines3DInterface(name string) gohome.Lines3DInterface {
	return &OpenGLLines3DInterface{
		Name: name,
	}
}

func (this *OpenGLRenderer) LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (gohome.Shader, error) {
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
	if this.backBufferMesh != nil {
		this.backBufferMesh.Render()
	} else {
		gl.GetError()
		gl.BindVertexArray(this.BackBufferVao)
		handleOpenGLError("Renderer", "OpenGL", "glBindVertexArray in RenderBackBuffer")
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		handleOpenGLError("Renderer", "OpenGL", "RenderError in RenderBackBuffer")
		gl.BindVertexArray(0)
		handleOpenGLError("Renderer", "OpenGL", "glBindVertexArray with 0 in RenderBackBuffer")
	}
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
		false,
	}
}

func (this *OpenGLRenderer) SetNativeResolution(width, height uint32) {
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
func (this *OpenGLRenderer) GetNativeResolution() mgl32.Vec2 {
	return [2]float32{float32(gohome.RenderMgr.BackBufferMS.GetWidth()), float32(gohome.RenderMgr.BackBufferMS.GetHeight())}
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

func (this *OpenGLRenderer) GetVersioni() uint8 {
	var major, minor, combined int32
	gl.GetIntegerv(gl.MAJOR_VERSION, &major)
	gl.GetIntegerv(gl.MINOR_VERSION, &minor)

	combined = major*10 + minor

	return uint8(combined)
}

func (this *OpenGLRenderer) gatherAvailableFunctions() {
	combined := this.GetVersioni()
	this.version = uint8(combined)
	if combined >= 30 {
		this.availableFunctions["VERTEX_ID"] = true
		this.availableFunctions["VERTEX_ARRAY"] = true
		this.availableFunctions["BLIT_FRAMEBUFFER"] = true
	}
	if combined >= 31 {
		this.availableFunctions["INSTANCED"] = true
	}
	if combined >= 32 {
		this.availableFunctions["MULTISAMPLE"] = true
		this.availableFunctions["FRAMEBUFFER_TEXTURE"] = true
		this.availableFunctions["GEOMETRY_SHADER"] = true
	}
	if combined >= 40 {
		this.availableFunctions["INDIRECT"] = true
	}
}

func (this *OpenGLRenderer) HasFunctionAvailable(function string) bool {
	v, ok := this.availableFunctions[function]
	return ok && v
}

func (this *OpenGLRenderer) FilterShaderFiles(name, file, shader_type string) string {
	return file
}

func (this *OpenGLRenderer) FilterShaderSource(name, source, shader_type string) string {
	return source
}

func (this *OpenGLRenderer) SetBackgroundColor(bgColor color.Color) {
	this.backgroundColor = bgColor
}

func (this *OpenGLRenderer) GetBackgroundColor() color.Color {
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
		case gl.STACK_OVERFLOW:
			errString = "STACK_OVERFLOW"
		case gl.STACK_UNDERFLOW:
			errString = "STACK_UNDERFLOW"
		case gl.OUT_OF_MEMORY:
			errString = "OUT_OF_MEMORY"
		case gl.INVALID_FRAMEBUFFER_OPERATION:
			errString = "INVALID_FRAMEBUFFER_OPERATION"
		case gl.CONTEXT_LOST:
			errString = "CONTEXT_LOST"
		case 0x8031:
			errString = "TABLE_TOO_LARGE"
		default:
			errString = strconv.Itoa(int(err))
		}
		gohome.ErrorMgr.Error(tag, objectName, errorPrefix+" ErrorCode: "+errString)
	}
}

func (this *OpenGLRenderer) CreateShape2DInterface(name string) gohome.Shape2DInterface {
	return &OpenGLShape2DInterface{
		Name: name,
	}
}

func (this *OpenGLRenderer) SetDepthTesting(b bool) {
	if b {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func (this *OpenGLRenderer) GetName() string {
	return "OpenGL"
}

func maxMultisampleSamples() int32 {
	var data int32
	gl.GetIntegerv(gl.MAX_INTEGER_SAMPLES, &data)
	return data
}
