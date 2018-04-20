package main

import (
	// "encoding/binary"
	"fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGLES"
	// "golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/gl"
)

type TestRenderObject struct {
}

func (this *TestRenderObject) Render() {
	if gohome.InputMgr.IsPressed(gohome.MouseButtonLeft) {
		gohome.Render.ClearScreen(gohome.Color{255, 0, 0, 255})
	} else {
		gohome.Render.ClearScreen(gohome.Color{0, 255, 0, 255})
	}

	triangleShader.Use()
	renderMesh()
	triangleShader.Unuse()
}

func (this *TestRenderObject) SetShader(s gohome.Shader) {

}
func (this *TestRenderObject) GetShader() gohome.Shader {
	return nil
}
func (this *TestRenderObject) SetType(rtype gohome.RenderType) {

}
func (this *TestRenderObject) GetType() gohome.RenderType {
	return gohome.TYPE_2D_NORMAL
}
func (this *TestRenderObject) IsVisible() bool {
	return true
}
func (this *TestRenderObject) NotRelativeCamera() int {
	return -1
}

type TestScene2 struct {
}

func (this *TestScene2) Init() {
	gohome.InitDefaultValues()
	gohome.FPSLimit.MaxFPS = 1000

	gohome.ResourceMgr.LoadTexture("Icon", "icon.png")
	gohome.ResourceMgr.LoadShader("Triangle", "triangleVert.glsl", "triangleFrag.glsl", "", "", "", "")

	gohome.RenderMgr.AddObject(&TestRenderObject{}, nil)

	width, height := gohome.Render.GetNativeResolution()

	fmt.Println("NativeRes:", width, height)

	gohome.RenderMgr.EnableBackBuffer = true

	render, _ := gohome.Render.(*renderer.OpenGLESRenderer)
	gles = render.GetContext()

	createMesh()
	triangleShader = gohome.ResourceMgr.GetShader("Triangle")
}

func (this *TestScene2) Update(delta_time float32) {

}

func (this *TestScene2) Terminate() {

}

var triangle gohome.Mesh2D
var triangleShader gohome.Shader
var gles gl.Context

func createMesh() {
	var vertices []gohome.Mesh2DVertex = []gohome.Mesh2DVertex{
		gohome.Mesh2DVertex{
			-0.5, -0.5,
			0.0, 0.0,
		},
		gohome.Mesh2DVertex{
			0.5, -0.5,
			1.0, 0.0,
		},
		gohome.Mesh2DVertex{
			0.0, 0.5,
			0.5, 1.0,
		},
	}
	var indices []uint32 = []uint32{
		0, 1, 2,
	}
	triangle = gohome.Render.CreateMesh2D("Triangle")
	triangle.AddVertices(vertices[:], indices[:])
	triangle.Load()
}

func renderMesh() {
	triangle.Render()
}
