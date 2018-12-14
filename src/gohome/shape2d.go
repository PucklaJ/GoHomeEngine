package gohome

const (
	SHAPE_2D_SHADER_NAME string = "Shape2D"
)

type Shape2D struct {
	NilRenderObject

	Name string

	shapeInterface      Shape2DInterface
	Transform           *TransformableObject2D
	Visible             bool
	NotRelativeToCamera int
	Depth               uint8
}

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

func (this *Shape2D) AddPoints(points []Shape2DVertex) {
	this.shapeInterface.AddPoints(points)
}

func (this *Shape2D) AddLines(lines []Line2D) {
	this.shapeInterface.AddLines(lines)
}

func (this *Shape2D) AddTriangles(tris []Triangle2D) {
	this.shapeInterface.AddTriangles(tris)
}

func (this *Shape2D) SetDrawMode(mode uint8) {
	this.shapeInterface.SetDrawMode(mode)
}

func (this *Shape2D) SetPointSize(size float32) {
	this.shapeInterface.SetPointSize(size)
}

func (this *Shape2D) SetLineWidth(width float32) {
	this.shapeInterface.SetLineWidth(width)
}

func (this *Shape2D) GetShader() Shader {
	return ResourceMgr.GetShader(SHAPE_2D_SHADER_NAME)
}

func (this *Shape2D) GetTransformableObject() TransformableObject {
	return this.Transform
}

func (this *Shape2D) Terminate() {
	this.shapeInterface.Terminate()
}

func (this *Shape2D) IsVisible() bool {
	return this.Visible
}

func (this *Shape2D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}

func (this *Shape2D) GetType() RenderType {
	return TYPE_2D_NORMAL
}

func (this *Shape2D) Render() {
	shader := RenderMgr.CurrentShader
	if shader != nil {
		shader.SetUniformF(DEPTH_UNIFORM_NAME, convertDepth(this.Depth))
	}
	this.shapeInterface.Render()
}

func (this *Shape2D) Load() {
	this.shapeInterface.Load()
}
