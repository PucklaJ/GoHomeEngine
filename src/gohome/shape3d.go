package gohome

import (
	"image/color"

	"github.com/PucklaJ/mathgl/mgl32"
)

const (
	SHAPE3D_SHADER_NAME string = "Shape3D"
)

// A 3D shape as a RenderObject
type Shape3D struct {
	NilRenderObject
	// The name of this shape
	Name           string
	shapeInterface Shape3DInterface

	transform TransformableObject
	// The transform of this shape
	Transform *TransformableObject3D
	// Wether this shape is visible
	Visible bool
	shader  Shader
	// The index of the camera to which this shape is relative
	NotRelativeToCamera int
	rtype               RenderType
}

// Initialises all values
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

// Add a point to the shape
func (this *Shape3D) AddPoint(point Shape3DVertex) {
	this.shapeInterface.AddPoints([]Shape3DVertex{point})
}

// Add multiple points to the shape
func (this *Shape3D) AddPoints(points []Shape3DVertex) {
	this.shapeInterface.AddPoints(points)
}

// Add a line to the shape
func (this *Shape3D) AddLine(line Line3D) {
	this.shapeInterface.AddPoints(line[:])
}

// Add multiple lines to the shape
func (this *Shape3D) AddLines(lines []Line3D) {
	for _, l := range lines {
		this.AddLine(l)
	}
}

// Add a triangle to the shape
func (this *Shape3D) AddTriangle(tri Triangle3D) {
	this.shapeInterface.AddPoints(tri[:])
}

// Add multiple triangles to the shape
func (this *Shape3D) AddTriangles(tris []Triangle3D) {
	for _, t := range tris {
		this.AddTriangle(t)
	}
}

// Set the draw mode (POINTS,LINES,TRIANGLES)
func (this *Shape3D) SetDrawMode(drawMode uint8) {
	this.shapeInterface.SetDrawMode(drawMode)
}

// Sets the point size
func (this *Shape3D) SetPointSize(size float32) {
	this.shapeInterface.SetPointSize(size)
}

// Sets the line width
func (this *Shape3D) SetLineWidth(width float32) {
	this.shapeInterface.SetLineWidth(width)
}

// Loads the data into the GPU
func (this *Shape3D) Load() {
	this.shapeInterface.Load()
}

// Sets the color of this shape
func (this *Shape3D) SetColor(col color.Color) {
	points := this.shapeInterface.GetPoints()
	for _, p := range points {
		p.SetColor(col)
	}
}

// Returns all points of the shape
func (this *Shape3D) GetPoints() []Shape3DVertex {
	return this.shapeInterface.GetPoints()
}

// Calls the draw method on the data
func (this *Shape3D) Render() {
	this.shapeInterface.Render()
}

// Sets the shader of this shape
func (this *Shape3D) SetShader(s Shader) {
	this.shader = s
}

// Returns the shader of this shape
func (this *Shape3D) GetShader() Shader {
	if this.shader == nil {
		this.shader = ResourceMgr.GetShader(SHAPE3D_SHADER_NAME)
	}
	return this.shader
}

// Sets the render type of the shape
func (this *Shape3D) SetType(rtype RenderType) {
	this.rtype = rtype
}

// Returns the render type of the shape
func (this *Shape3D) GetType() RenderType {
	return this.rtype
}

// Returns wether this shape is visible
func (this *Shape3D) IsVisible() bool {
	return this.Visible
}

// Returns the index of the camera to which this shape is not relative to
func (this *Shape3D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}

// Sets the transformable object of this shape
func (this *Shape3D) SetTransformableObject(tobj TransformableObject) {
	this.transform = tobj
	if tobj != nil {
		this.Transform = tobj.(*TransformableObject3D)
	} else {
		this.Transform = nil
	}
}

// Returns the transformable object of this shape
func (this *Shape3D) GetTransformableObject() TransformableObject {
	return this.transform
}

// Cleans everything up
func (this *Shape3D) Terminate() {
	this.shapeInterface.Terminate()
}
