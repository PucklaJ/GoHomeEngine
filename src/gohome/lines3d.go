package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"image/color"
)

const (
	LINES3D_SHADER_NAME string = "Lines3DShader"
)

type Lines3D struct {
	NilRenderObject
	Name           string
	linesInterface Lines3DInterface

	transform           TransformableObject
	Transform           *TransformableObject3D
	Visible             bool
	shader              Shader
	NotRelativeToCamera int
	rtype               RenderType
}

func (this *Lines3D) Init() {
	if ResourceMgr.GetShader(LINES3D_SHADER_NAME) == nil {
		ResourceMgr.LoadShaderSource(LINES3D_SHADER_NAME, LINES_3D_SHADER_VERTEX_SOURCE_OPENGL, LINES_3D_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
	}
	this.linesInterface = Render.CreateLines3DInterface(this.Name)
	this.linesInterface.Init()
	this.Transform = &TransformableObject3D{
		Scale:    [3]float32{1.0, 1.0, 1.0},
		Rotation: mgl32.QuatRotate(0.0, mgl32.Vec3{0.0, 1.0, 0.0}),
	}
	this.transform = this.Transform
	this.Visible = true
	this.NotRelativeToCamera = -1
	this.rtype = TYPE_3D_NORMAL
}

func (this *Lines3D) AddLine(line Line3D) {
	this.linesInterface.AddLines([]Line3D{line})
}

func (this *Lines3D) AddLines(lines []Line3D) {
	this.linesInterface.AddLines(lines)
}

func (this *Lines3D) Load() {
	this.linesInterface.Load()
}

func (this *Lines3D) SetColor(col color.Color) {
	for i := 0; i < len(this.linesInterface.GetLines()); i++ {
		this.linesInterface.GetLines()[i].SetColor(col)
	}
}

func (this *Lines3D) GetLines() []Line3D {
	return this.linesInterface.GetLines()
}

func (this *Lines3D) Render() {
	this.linesInterface.Render()
}
func (this *Lines3D) SetShader(s Shader) {
	this.shader = s
}
func (this *Lines3D) GetShader() Shader {
	if this.shader == nil {
		this.shader = ResourceMgr.GetShader(LINES3D_SHADER_NAME)
	}
	return this.shader
}
func (this *Lines3D) SetType(rtype RenderType) {
	this.rtype = rtype
}
func (this *Lines3D) GetType() RenderType {
	return this.rtype
}
func (this *Lines3D) IsVisible() bool {
	return this.Visible
}
func (this *Lines3D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}
func (this *Lines3D) SetTransformableObject(tobj TransformableObject) {
	this.transform = tobj
	if tobj != nil {
		this.Transform = tobj.(*TransformableObject3D)
	} else {
		this.Transform = nil
	}
}
func (this *Lines3D) GetTransformableObject() TransformableObject {
	return this.transform
}

func (this *Lines3D) Terminate() {
	this.linesInterface.Terminate()
}
