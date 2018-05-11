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
}

func (this *Entity3D) commonInit(tobj *TransformableObject3D) {
	if tobj != nil {
		tobj.Scale = [3]float32{1.0, 1.0, 1.0}
	}
	this.Visible = true
	this.NotRelativeToCamera = -1
	this.RenderType = TYPE_3D_NORMAL
	this.Shader = ResourceMgr.GetShader(ENTITY3D_SHADER_NAME)
}

func (this *Entity3D) InitName(name string, tobj *TransformableObject3D) {
	this.commonInit(tobj)
	this.Model3D = ResourceMgr.GetModel(name)
	this.Name = name
}

func (this *Entity3D) InitMesh(mesh Mesh3D, tobj *TransformableObject3D) {
	this.commonInit(tobj)
	this.Model3D = &Model3D{
		Name: mesh.GetName(),
	}
	this.Model3D.AddMesh3D(mesh)
	this.Name = mesh.GetName()
}

func (this *Entity3D) InitModel(model *Model3D, tobj *TransformableObject3D) {
	this.commonInit(tobj)
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
