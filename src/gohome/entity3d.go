package gohome

const (
	ENTITY3D_SHADER_NAME string = "3D"
)

type Entity3D struct {
	Name                string
	Model3D             *Model3D
	Visible             bool
	NotRelativeToCamera int

	Shader
	RenderType

	Transform *TransformableObject3D
	transform TransformableObject
}

func (this *Entity3D) commonInit() {
	this.Transform = &TransformableObject3D{}
	this.Transform.Scale = [3]float32{1.0, 1.0, 1.0}

	this.transform = this.Transform

	this.Visible = true
	this.NotRelativeToCamera = -1
	this.RenderType = TYPE_3D_NORMAL
	this.Shader = ResourceMgr.GetShader(ENTITY3D_SHADER_NAME)
}

func (this *Entity3D) InitName(name string) {
	this.commonInit()
	this.Model3D = ResourceMgr.GetModel(name)
	this.Name = name
}

func (this *Entity3D) InitMesh(mesh Mesh3D) {
	this.commonInit()
	this.Model3D = &Model3D{
		Name: mesh.GetName(),
	}
	this.Model3D.AddMesh3D(mesh)
	this.Name = mesh.GetName()
}

func (this *Entity3D) InitModel(model *Model3D) {
	this.commonInit()
	this.Model3D = model
	if model != nil {
		this.Name = model.Name
	}
}

func (this *Entity3D) GetShader() Shader {
	return this.Shader
}

func (this *Entity3D) SetShader(s Shader) {
	this.Shader = s
}

func (this *Entity3D) SetType(rtype RenderType) {
	this.RenderType = rtype
}

func (this *Entity3D) GetType() RenderType {
	return this.RenderType
}

func (this *Entity3D) Render() {
	if this.Model3D != nil {
		this.Model3D.Render()
	}
}

func (this *Entity3D) Terminate() {
	if this.Model3D != nil {
		this.Model3D.Terminate()
	}
}

func (this *Entity3D) IsVisible() bool {
	return this.Visible
}

func (this *Entity3D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}

func (this *Entity3D) SetTransformableObject(tobj TransformableObject) {
	this.transform = tobj
	if tobj != nil {
		this.Transform = tobj.(*TransformableObject3D)
	} else {
		this.Transform = nil
	}
}

func (this *Entity3D) GetTransformableObject() TransformableObject {
	return this.transform
}
