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
	shapeInterface Shape3DInterface

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
	this.shapeInterface = Render.CreateShape3DInterface(this.Name)
	this.shapeInterface.Init()
	this.Transform = &TransformableObject3D{
		Scale:    [3]float32{1.0, 1.0, 1.0},
		Rotation: mgl32.QuatRotate(0.0, mgl32.Vec3{0.0, 1.0, 0.0}),
	}
	this.transform = this.Transform
	this.Visible = true
	this.NotRelativeToCamera = -1
	this.rtype = TYPE_3D_NORMAL
}

func (this *Shape3D) AddPoint(point Shape3DVertex) {
	this.shapeInterface.AddPoints([]Shape3DVertex{point})
}

func (this *Shape3D) AddPoints(points []Shape3DVertex) {
	this.shapeInterface.AddPoints(points)
}

func (this *Shape3D) AddLine(line Line3D) {
	this.shapeInterface.AddPoints(line[:])
}

func (this *Shape3D) AddLines(lines []Line3D) {
	for _, l := range lines {
		this.AddLine(l)
	}
}

func (this *Shape3D) AddTriangle(tri Triangle3D) {
	this.shapeInterface.AddPoints(tri[:])
}

func (this *Shape3D) AddTriangles(tris []Triangle3D) {
	for _, t := range tris {
		this.AddTriangle(t)
	}
}

func (this *Shape3D) SetDrawMode(drawMode uint8) {
	this.shapeInterface.SetDrawMode(drawMode)
}

func (this *Shape3D) SetPointSize(size float32) {
	this.shapeInterface.SetPointSize(size)
}

func (this *Shape3D) SetLineWidth(width float32) {
	this.shapeInterface.SetLineWidth(width)
}

func (this *Shape3D) Load() {
	this.shapeInterface.Load()
}

func (this *Shape3D) SetColor(col color.Color) {
	points := this.shapeInterface.GetPoints()
	for _, p := range points {
		p.SetColor(col)
	}
}

func (this *Shape3D) GetPoints() []Shape3DVertex {
	return this.shapeInterface.GetPoints()
}

func (this *Shape3D) Render() {
	this.shapeInterface.Render()
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
	this.shapeInterface.Terminate()
}
