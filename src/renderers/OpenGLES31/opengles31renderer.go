package renderer

import (
	"image/color"
	"strconv"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	gl "github.com/PucklaJ/android-go/gles31"
	"github.com/PucklaJ/mathgl/mgl32"
)

type OpenGLES31Renderer struct {
	BackBufferVao      uint32
	CurrentTextureUnit uint32

	availableFunctions map[string]bool
	backBufferMesh     *OpenGLES31Mesh2D
	backgroundColor    color.Color
	version            int
}

func (this *OpenGLES31Renderer) createBackBufferMesh() {
	this.backBufferMesh = CreateOpenGLES31Mesh2D("BackBufferMesh")

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

func (this *OpenGLES31Renderer) Init() error {
	version := gl.GetString(gl.VERSION)
	versioni := this.GetVersioni()
	if version == "" {
		version = strconv.FormatUint(uint64(versioni), 10)
	}
	gohome.ErrorMgr.Log("Renderer", "OpenGLES31\t", "Version: "+version)
	if versioni < 21 {
		gohome.ErrorMgr.Warning("Renderer", "OpenGLES31", "You don't have a graphics card or your graphics card is not supported! Minimum: OpenGLES 2.0")
	}

	this.CurrentTextureUnit = 0

	this.availableFunctions = make(map[string]bool)
	this.gatherAvailableFunctions()

	if !this.HasFunctionAvailable("VERTEX_ID") || !this.HasFunctionAvailable("MULTISAMPLE") {
		this.createBackBufferMesh()
	} else {
		var buf [1]uint32
		gl.GenVertexArrays(1, buf[:])
		this.BackBufferVao = buf[0]
	}

	return nil
}

func (this *OpenGLES31Renderer) AfterInit() {
	gl.DepthFunc(gl.LEQUAL)
	gl.ClearDepthf(2.0)
	gl.Enable(gl.BLEND)
	gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)

	gl.FrontFace(gl.CCW)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.DEPTH_TEST)
}

func (this *OpenGLES31Renderer) HasExtension(name string) bool {
	var numExtensions [1]int32
	gl.GetIntegerv(gl.NUM_EXTENSIONS, numExtensions[:])
	for i := 0; i < int(numExtensions[0]); i++ {
		ext := gl.GetStringi(gl.EXTENSIONS, uint32(i))
		if ext == name {
			return true
		}
	}
	return false
}

func (this *OpenGLES31Renderer) SetWireFrame(b bool) {
	gohome.ErrorMgr.Warning("Renderer", "OpenGLES31", "SetWireFrame does not work in OpenGLES 3.0")
}

func (this *OpenGLES31Renderer) Terminate() {
	if this.BackBufferVao != 0 {
		var buf [1]uint32
		buf[0] = this.BackBufferVao
		gl.DeleteVertexArrays(1, buf[:])
	}
	if this.backBufferMesh != nil {
		this.backBufferMesh.Terminate()
	}
}

func (*OpenGLES31Renderer) ClearScreen(c color.Color) {
	clearColor := gohome.ColorToVec4(c)
	gl.ClearColor(clearColor[0], clearColor[1], clearColor[2], clearColor[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
}

type OpenGLES31Error struct {
	errorString string
}

func (oerr OpenGLES31Error) Error() string {
	return oerr.errorString
}

func (*OpenGLES31Renderer) CreateTexture(name string, multiSampled bool) gohome.Texture {
	return CreateOpenGLES31Texture(name)
}

func (*OpenGLES31Renderer) CreateMesh2D(name string) gohome.Mesh2D {
	return CreateOpenGLES31Mesh2D(name)
}

func (*OpenGLES31Renderer) CreateMesh3D(name string) gohome.Mesh3D {
	return CreateOpenGLES31Mesh3D(name)
}

func (*OpenGLES31Renderer) CreateRenderTexture(name string, width, height, textures int, depthBuffer, multiSampled, shadowMap, cubeMap bool) gohome.RenderTexture {
	return CreateOpenGLES31RenderTexture(name, width, height, textures, depthBuffer, shadowMap, cubeMap)
}

func (*OpenGLES31Renderer) CreateCubeMap(name string) gohome.CubeMap {
	return CreateOpenGLES31CubeMap(name)
}

func (*OpenGLES31Renderer) CreateInstancedMesh3D(name string) gohome.InstancedMesh3D {
	return CreateOpenGLES31InstancedMesh3D(name)
}

func (*OpenGLES31Renderer) CreateShape3DInterface(name string) gohome.Shape3DInterface {
	return &OpenGLES31Shape3DInterface{
		Name: name,
	}
}

func (this *OpenGLES31Renderer) LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (gohome.Shader, error) {
	var shader *OpenGLES31Shader
	var err error

	shader, err = CreateOpenGLES31Shader(name)
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

func (this *OpenGLES31Renderer) RenderBackBuffer() {
	if this.backBufferMesh != nil {
		this.backBufferMesh.Render()
	} else {
		gl.GetError()
		gl.BindVertexArray(this.BackBufferVao)
		handleOpenGLES31Error("Renderer", "OpenGLES31", "glBindVertexArray in RenderBackBuffer")
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		handleOpenGLES31Error("Renderer", "OpenGLES31", "RenderError in RenderBackBuffer")
		gl.BindVertexArray(0)
		handleOpenGLES31Error("Renderer", "OpenGLES31", "glBindVertexArray with 0 in RenderBackBuffer")
	}
}

func (this *OpenGLES31Renderer) SetViewport(viewport gohome.Viewport) {
	gl.Viewport(int32(viewport.X), int32(viewport.Y), int32(viewport.Width), int32(viewport.Height))
}

func (this *OpenGLES31Renderer) GetViewport() gohome.Viewport {
	var data [4]int32
	gl.GetIntegerv(gl.VIEWPORT, data[:])

	return gohome.Viewport{
		0,
		int(data[0]), int(data[1]),
		int(data[2]), int(data[3]),
		false,
	}
}

func (this *OpenGLES31Renderer) SetNativeResolution(width, height int) {
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
func (this *OpenGLES31Renderer) GetNativeResolution() mgl32.Vec2 {
	return [2]float32{float32(gohome.RenderMgr.BackBufferMS.GetWidth()), float32(gohome.RenderMgr.BackBufferMS.GetHeight())}
}
func (this *OpenGLES31Renderer) OnResize(newWidth, newHeight int) {
	gl.Viewport(0, 0, int32(newWidth), int32(newHeight))
}

func (this *OpenGLES31Renderer) PreRender() {
	this.CurrentTextureUnit = 1
}
func (this *OpenGLES31Renderer) AfterRender() {
	this.CurrentTextureUnit = 1
}

func (this *OpenGLES31Renderer) SetBacckFaceCulling(b bool) {
	if b {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (this *OpenGLES31Renderer) GetMaxTextures() int {
	var data [1]int32
	gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS, data[:])
	return int(data[0])
}

func (this *OpenGLES31Renderer) NextTextureUnit() uint32 {
	val := this.CurrentTextureUnit
	this.CurrentTextureUnit++
	return val
}

func (this *OpenGLES31Renderer) DecrementTextureUnit(amount uint32) {
	this.CurrentTextureUnit -= amount
}

func (this *OpenGLES31Renderer) GetVersioni() int {
	var major, minor, combined [1]int32
	gl.GetIntegerv(gl.MAJOR_VERSION, major[:])
	gl.GetIntegerv(gl.MINOR_VERSION, minor[:])

	combined[0] = major[0]*10 + minor[0]

	return int(combined[0])
}

func (this *OpenGLES31Renderer) gatherAvailableFunctions() {
	this.availableFunctions["VERTEX_ID"] = true
	this.availableFunctions["VERTEX_ARRAY"] = true
	this.availableFunctions["BLIT_FRAMEBUFFER"] = true
	this.availableFunctions["BLIT_FRAMEBUFFER_SCREEN"] = true
	this.availableFunctions["INSTANCED"] = true
	this.availableFunctions["INDIRECT"] = true
}

func (this *OpenGLES31Renderer) HasFunctionAvailable(function string) bool {
	v, ok := this.availableFunctions[function]
	return ok && v
}

func (this *OpenGLES31Renderer) SetBackgroundColor(bgColor color.Color) {
	this.backgroundColor = bgColor
}

func (this *OpenGLES31Renderer) GetBackgroundColor() color.Color {
	return this.backgroundColor
}

func handleOpenGLES31Error(tag, objectName, errorPrefix string) {
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

func (this *OpenGLES31Renderer) CreateShape2DInterface(name string) gohome.Shape2DInterface {
	return &OpenGLES31Shape2DInterface{
		Name: name,
	}
}

func (this *OpenGLES31Renderer) SetDepthTesting(b bool) {
	if b {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func (this *OpenGLES31Renderer) GetName() string {
	return "OpenGLES31"
}
