package gohome

const (
	SHAPE_2D_SHADER_NAME string = "Shape2D"
)

// A 2D shape as a RenderObject
type Shape2D struct {
	NilRenderObject

	// The name of the shape
	Name string

	shapeInterface      Shape2DInterface
	// The transform of the shape
	Transform           *TransformableObject2D
	// Wether the shape is visible
	Visible             bool
	// The index of the camera to which this shape is not relative to
	NotRelativeToCamera int
	// The depth of this shape (0-255)
	Depth               uint8
}

// Initialises everything with default values
func (this *Shape2D) Init() {
	if ResourceMgr.GetShader(GetShaderName2D(SHADER_TYPE_SHAPE2D, 0)) == nil {
		LoadGeneratedShader2D(SHADER_TYPE_SHAPE2D, 0)
	}
	this.shapeInterface = Render.CreateShape2DInterface(this.Name)
	this.shapeInterface.Init()
	this.Transform = &TransformableObject2D{
		Size:  [2]float32{1.0, 1.0},
		Scale: [2]float32{1.0, 1.0},
	}
	this.Visible = true
	this.NotRelativeToCamera = -1
}

// Add points to the shape
func (this *Shape2D) AddPoints(points []Shape2DVertex) {
	this.shapeInterface.AddPoints(points)
}

// Add lines to the shape
func (this *Shape2D) AddLines(lines []Line2D) {
	this.shapeInterface.AddLines(lines)
}

// Add triangles to the shape
func (this *Shape2D) AddTriangles(tris []Triangle2D) {
	this.shapeInterface.AddTriangles(tris)
}

// Set the draw mode (POINTS,LINES,TRIANGLES)
func (this *Shape2D) SetDrawMode(mode uint8) {
	this.shapeInterface.SetDrawMode(mode)
}

// Set the point size
func (this *Shape2D) SetPointSize(size float32) {
	this.shapeInterface.SetPointSize(size)
}

// Set the width of the lines
func (this *Shape2D) SetLineWidth(width float32) {
	this.shapeInterface.SetLineWidth(width)
}

// Returns the shader used for rendering this shape
func (this *Shape2D) GetShader() Shader {
	return ResourceMgr.GetShader(SHAPE_2D_SHADER_NAME)
}

// Returns the transformable object used for this shape
func (this *Shape2D) GetTransformableObject() TransformableObject {
	return this.Transform
}

// Cleans everything up
func (this *Shape2D) Terminate() {
	this.shapeInterface.Terminate()
}

// Returns wether this shape is visible
func (this *Shape2D) IsVisible() bool {
	return this.Visible
}

// Returns the index of the camera to which this shape is not relative to
func (this *Shape2D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}

// Returns the render type of this shape
func (this *Shape2D) GetType() RenderType {
	return TYPE_2D_NORMAL
}

// Calls the draw method on the shape
func (this *Shape2D) Render() {
	shader := RenderMgr.CurrentShader
	if shader != nil {
		shader.SetUniformF(DEPTH_UNIFORM_NAME, convertDepth(this.Depth))
	}
	this.shapeInterface.Render()
}

// Loads the data tot the GPU
func (this *Shape2D) Load() {
	this.shapeInterface.Load()
}
