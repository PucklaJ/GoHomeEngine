package renderer

import (
	"github.com/PucklaMotzer09/webgl"
	"image/color"
	"strconv"

	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"github.com/gopherjs/gopherjs/js"
)

var gl *webgl.Context

type WebGLRenderer struct {
	CurrentTextureUnit uint32

	availableFunctions map[string]bool
	backBufferMesh     *WebGLMesh2D
	backgroundColor    color.Color
	version            uint8
}

func (this *WebGLRenderer) createBackBufferMesh() {
	this.backBufferMesh = CreateWebGLMesh2D("BackBufferMesh")

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

func (this *WebGLRenderer) Init() error {
	document := js.Global.Get("document")
	canvas := document.Call("getElementById", "gohome_canvas")

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false

	ctx, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		return err
	}
	gl = ctx

	version := gl.GetString(gl.VERSION)
	versioni := gl.GetVersioni()
	if version == "" {
		version = strconv.FormatUint(uint64(versioni), 10)
	}
	gohome.ErrorMgr.Log("Renderer", "WebGL\t", "Version: "+version)

	this.CurrentTextureUnit = 0

	this.availableFunctions = make(map[string]bool)
	this.gatherAvailableFunctions()

	this.createBackBufferMesh()

	return nil
}

func (this *WebGLRenderer) AfterInit() {
	gl.DepthFunc(gl.LEQUAL)
	gl.ClearDepth(2.0)
	gl.Enable(gl.BLEND)
	gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)

	gl.FrontFace(gl.CCW)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.DEPTH_TEST)
}

func (this *WebGLRenderer) HasExtension(name string) bool {
	ext := gl.GetExtension(name)
	return ext != js.Undefined
}

func (this *WebGLRenderer) SetWireFrame(b bool) {
	gohome.ErrorMgr.Warning("Renderer", "WebGL", "SetWireFrame does not work")
}

func (this *WebGLRenderer) Terminate() {
	if this.backBufferMesh != nil {
		this.backBufferMesh.Terminate()
	}
}

func (*WebGLRenderer) ClearScreen(c color.Color) {
	clearColor := gohome.ColorToVec4(c)
	gl.ClearColor(clearColor[0], clearColor[1], clearColor[2], clearColor[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
}

type WebGLError struct {
	errorString string
}

func (oerr WebGLError) Error() string {
	return oerr.errorString
}

func (*WebGLRenderer) CreateTexture(name string, multiSampled bool) gohome.Texture {
	return CreateWebGLTexture(name)
}

func (*WebGLRenderer) CreateMesh2D(name string) gohome.Mesh2D {
	return CreateWebGLMesh2D(name)
}

func (*WebGLRenderer) CreateMesh3D(name string) gohome.Mesh3D {
	return CreateWebGLMesh3D(name)
}

func (*WebGLRenderer) CreateRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) gohome.RenderTexture {
	return CreateWebGLRenderTexture(name, width, height, textures, depthBuffer, shadowMap, cubeMap)
}

func (*WebGLRenderer) CreateCubeMap(name string) gohome.CubeMap {
	return CreateWebGLCubeMap(name)
}

func (*WebGLRenderer) CreateInstancedMesh3D(name string) gohome.InstancedMesh3D {
	return CreateWebGLInstancedMesh3D(name)
}

func (*WebGLRenderer) CreateLines3DInterface(name string) gohome.Lines3DInterface {
	return &WebGLLines3DInterface{
		Name: name,
	}
}

func (this *WebGLRenderer) LoadShader(name, vertex_contents, fragment_contents, geometry_contents, tesselletion_control_contents, eveluation_contents, compute_contents string) (gohome.Shader, error) {
	var shader *WebGLShader
	var err error

	shader, err = CreateWebGLShader(name)
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

func (this *WebGLRenderer) RenderBackBuffer() {
	this.backBufferMesh.Render()
}

func (this *WebGLRenderer) SetViewport(viewport gohome.Viewport) {
	gl.Viewport(viewport.X, viewport.Y, viewport.Width, viewport.Height)
}

func (this *WebGLRenderer) GetViewport() gohome.Viewport {
	data := gl.GetIntegerv(gl.VIEWPORT)

	return gohome.Viewport{
		0,
		data[0], data[1],
		data[2], data[3],
		false,
	}
}

func (this *WebGLRenderer) SetNativeResolution(width, height uint32) {
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
func (this *WebGLRenderer) GetNativeResolution() mgl32.Vec2 {
	return [2]float32{float32(gohome.RenderMgr.BackBufferMS.GetWidth()), float32(gohome.RenderMgr.BackBufferMS.GetHeight())}
}
func (this *WebGLRenderer) OnResize(newWidth, newHeight uint32) {
	gl.Viewport(0, 0, int(newWidth), int(newHeight))
}

func (this *WebGLRenderer) PreRender() {
	this.CurrentTextureUnit = 1
}
func (this *WebGLRenderer) AfterRender() {
	this.CurrentTextureUnit = 1
}

func (this *WebGLRenderer) SetBacckFaceCulling(b bool) {
	if b {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (this *WebGLRenderer) GetMaxTextures() int32 {
	data := gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS)
	return int32(data[0])
}

func (this *WebGLRenderer) NextTextureUnit() uint32 {
	val := this.CurrentTextureUnit
	this.CurrentTextureUnit++
	return val
}

func (this *WebGLRenderer) DecrementTextureUnit(amount uint32) {
	this.CurrentTextureUnit -= amount
}

func (this *WebGLRenderer) GetVersioni() uint8 {
	return uint8(gl.GetVersioni())
}

func (this *WebGLRenderer) gatherAvailableFunctions() {
	versioni := gl.GetVersioni()
	if versioni == 20 {
		this.availableFunctions["VERTEX_ID"] = true
		this.availableFunctions["VERTEX_ARRAY"] = true
		this.availableFunctions["INSTANCED"] = true
		this.availableFunctions["DRAW_BUFFERS"] = true
		this.availableFunctions["BLIT_FRAMEBUFFER"] = true
	}
}

func (this *WebGLRenderer) HasFunctionAvailable(function string) bool {
	v, ok := this.availableFunctions[function]
	return ok && v
}

func (this *WebGLRenderer) SetBackgroundColor(bgColor color.Color) {
	this.backgroundColor = bgColor
}

func (this *WebGLRenderer) GetBackgroundColor() color.Color {
	return this.backgroundColor
}

func handleWebGLError(tag, objectName, errorPrefix string) bool {
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
		return false
	}
	return true
}

func (this *WebGLRenderer) CreateShape2DInterface(name string) gohome.Shape2DInterface {
	return &WebGLShape2DInterface{
		Name: name,
	}
}

func (this *WebGLRenderer) SetDepthTesting(b bool) {
	if b {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func (this *WebGLRenderer) GetName() string {
	return "WebGL"
}

func maxMultisampleSamples() int32 {
	return 0
}
