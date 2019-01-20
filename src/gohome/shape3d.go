package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"image/color"
)

const (
	SHAPE3D_SHADER_NAME string = "Shape3D"
)

type Shape3D struct {
	NilRenderObject
	Name           string
	linesInterface Shape3DInterface

	transform           TransformableObject
	Transform           *TransformableObject3D
	Visible             bool
	shader              Shader
	NotRelativeToCamera int
	rtype               RenderType
}

func (this *Shape3D) Init() {
	if ResourceMgr.GetShader(SHAPE3D_SHADER_NAME) == nil {
		LoadGeneratedShader3D(SHADER_TYPE_SHAPE3D, 0)
	}
	this.linesInterface = Render.CreateShape3DInterface(this.Name)
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

func (this *Shape3D) AddLine(line Line3D) {
	this.linesInterface.AddLines([]Line3D{line})
}

func (this *Shape3D) AddLines(lines []Line3D) {
	this.linesInterface.AddLines(lines)
}

func (this *Shape3D) Load() {
	this.linesInterface.Load()
}

func (this *Shape3D) SetColor(col color.Color) {
	for i := 0; i < len(this.linesInterface.GetLines()); i++ {
		this.linesInterface.GetLines()[i].SetColor(col)
	}
}

func (this *Shape3D) GetLines() []Line3D {
	return this.linesInterface.GetLines()
}

func (this *Shape3D) Render() {
	this.linesInterface.Render()
}
func (this *Shape3D) SetShader(s Shader) {
	this.shader = s
}
func (this *Shape3D) GetShader() Shader {
	if this.shader == nil {
		this.shader = ResourceMgr.GetShader(SHAPE3D_SHADER_NAME)
	}
	return this.shader
}
func (this *Shape3D) SetType(rtype RenderType) {
	this.rtype = rtype
}
func (this *Shape3D) GetType() RenderType {
	return this.rtype
}
func (this *Shape3D) IsVisible() bool {
	return this.Visible
}
func (this *Shape3D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}
func (this *Shape3D) SetTransformableObject(tobj TransformableObject) {
	this.transform = tobj
	if tobj != nil {
		this.Transform = tobj.(*TransformableObject3D)
	} else {
		this.Transform = nil
	}
}
func (this *Shape3D) GetTransformableObject() TransformableObject {
	return this.transform
}

func (this *Shape3D) Terminate() {
	this.linesInterface.Terminate()
}
