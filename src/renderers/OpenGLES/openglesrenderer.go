package renderer

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/mobile/gl"
	"image/color"
	"strconv"
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
	backBufferMesh     *OpenGLESMesh2D
	backgroundColor    color.Color
}

func (this *OpenGLESRenderer) createBackBufferMesh() {
	this.backBufferMesh = CreateOpenGLESMesh2D("BackBufferMesh")

	vertices := []gohome.Mesh2DVertex{
		/*X,Y
		  U,V
		*/
		gohome.Mesh2DVertex{-1.0, -1.0, // LEFT-DOWN
			0.0, 0.0},

		gohome.Mesh2DVertex{1.0, -1.0, // RIGHT-DOWN
			1.0, 0.0},

		gohome.Mesh2DVertex{1.0, 1.0, // RIGHT-UP
			1.0, 1.0},

		gohome.Mesh2DVertex{-1.0, 1.0, // LEFT-UP
			0.0, 1.0},
	}

	indices := []uint32{
		0, 1, 2, // LEFT-TRI
		2, 3, 0, // RIGHT-TRI
	}

	this.backBufferMesh.AddVertices(vertices, indices)
	this.backBufferMesh.Load()

	// this.backBufferMesh = CreateOpenGLESMesh2D("BackBufferMesh")

	// var vertices []gohome.Mesh2DVertex = make([]gohome.Mesh2DVertex, 4)
	// var indices []uint32 = make([]uint32, 6)

	// vertices[0].Vertex(-1.0, -1.0)
	// vertices[1].Vertex(1.0, -1.0)
	// vertices[2].Vertex(1.0, 1.0)
	// vertices[3].Vertex(-1.0, 1.0)

	// vertices[0].TexCoord(0.0, 0.0)
	// vertices[1].TexCoord(1.0, 0.0)
	// vertices[2].TexCoord(1.0, 1.0)
	// vertices[3].TexCoord(0.0, 1.0)

	// indices[0] = 0
	// indices[1] = 1
	// indices[2] = 2
	// indices[3] = 2
	// indices[4] = 3
	// indices[5] = 0

	// this.backBufferMesh.AddVertices(vertices, indices)
	// this.backBufferMesh.Load()

}

func (this *OpenGLESRenderer) Init() error {
	this.CurrentTextureUnit = 0

	version := this.gles.GetString(gl.VERSION)
	gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_LOG, "Renderer", "OpenGLES\t", "Version: "+version+" "+gl.Version())
	this.gles.Enable(gl.DEPTH_TEST)
	this.gles.Enable(gl.CULL_FACE)

	this.createBackBufferMesh()

	return nil
}
func (this *OpenGLESRenderer) Terminate() {
	if this.backBufferMesh != nil {
		this.backBufferMesh.Terminate()
	}
}
func (this *OpenGLESRenderer) ClearScreen(c color.Color) {
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
	return CreateOpenGLESTexture(name)
}
func (this *OpenGLESRenderer) CreateMesh2D(name string) gohome.Mesh2D {
	return CreateOpenGLESMesh2D(name)
}
func (this *OpenGLESRenderer) CreateMesh3D(name string) gohome.Mesh3D {
	return CreateOpenGLESMesh3D(name)
}
func (this *OpenGLESRenderer) CreateRenderTexture(name string, width, height, textures uint32, depthBuffer, multiSampled, shadowMap, cubeMap bool) gohome.RenderTexture {
	return CreateOpenGLESRenderTexture(name, width, height, textures, depthBuffer, shadowMap, cubeMap)
}
func (this *OpenGLESRenderer) CreateCubeMap(name string) gohome.CubeMap {
	return CreateOpenGLESCubeMap(name)
}
func (this *OpenGLESRenderer) CreateInstancedMesh3D(name string) gohome.InstancedMesh3D {
	return CreateOpenGLESInstancedMesh3D(name)
}
func (this *OpenGLESRenderer) SetWireFrame(b bool) {

}
func (this *OpenGLESRenderer) SetViewport(viewport gohome.Viewport) {
	this.gles.Viewport(viewport.X, viewport.Y, viewport.Width, viewport.Height)
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
	if gohome.RenderMgr.BackBuffer2D == nil || gohome.RenderMgr.BackBuffer3D == nil || gohome.RenderMgr.BackBufferMS == nil || gohome.RenderMgr.BackBuffer == nil {
		return
	}

	previous := gohome.Viewport{
		X:      0,
		Y:      0,
		Width:  gohome.RenderMgr.BackBuffer.GetWidth(),
		Height: gohome.RenderMgr.BackBuffer.GetHeight(),
	}

	gohome.RenderMgr.BackBuffer2D.ChangeSize(width, height)
	gohome.RenderMgr.BackBuffer3D.ChangeSize(width, height)
	gohome.RenderMgr.BackBufferMS.ChangeSize(width, height)
	gohome.RenderMgr.BackBuffer.ChangeSize(width, height)

	current := gohome.Viewport{
		X:      0,
		Y:      0,
		Width:  gohome.RenderMgr.BackBuffer.GetWidth(),
		Height: gohome.RenderMgr.BackBuffer.GetHeight(),
	}

	gohome.RenderMgr.UpdateViewports(current, previous)
}
func (this *OpenGLESRenderer) GetNativeResolution() (uint32, uint32) {
	var width, height uint32

	width = uint32(gohome.RenderMgr.GetBackBuffer().GetWidth())
	height = uint32(gohome.RenderMgr.GetBackBuffer().GetHeight())

	return width, height
}
func (this *OpenGLESRenderer) OnResize(newWidth, newHeight uint32) {
	if this.gles != nil {
		this.gles.Viewport(0, 0, int(newWidth), int(newHeight))
		this.SetNativeResolution(newWidth, newHeight)
	}
}
func (this *OpenGLESRenderer) PreRender() {
	this.CurrentTextureUnit = 1
}
func (this *OpenGLESRenderer) AfterRender() {
	this.CurrentTextureUnit = 1
}
func (this *OpenGLESRenderer) RenderBackBuffer() {
	this.backBufferMesh.Render()
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

func (this *OpenGLESRenderer) GetContext() gl.Context {
	return this.gles
}

func (this *OpenGLESRenderer) FilterShaderFiles(name, file, shader_type string) string {
	return file
}

func (this *OpenGLESRenderer) SetBackgroundColor(bgColor color.Color) {
	this.backgroundColor = bgColor
}

func (this *OpenGLESRenderer) GetBackgroundColor() color.Color {
	return this.backgroundColor
}

func handleOpenGLESError(tag, objectName, errorPrefix string) {
	err := gohome.Render.(*OpenGLESRenderer).GetContext().GetError()
	if err != gl.NO_ERROR {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, tag, objectName, errorPrefix+"ErrorCode: "+strconv.Itoa(int(err)))
	}
}
