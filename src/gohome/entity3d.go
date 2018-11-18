package gohome

import "github.com/PucklaMotzer09/mathgl/mgl32"

const (
	ENTITY3D_SHADER_NAME                  string = "3D"
	ENTITY3D_NO_UV_SHADER_NAME            string = "3D NoUV"
	ENTITY3D_NO_UV_NO_SHADOWS_SHADER_NAME string = "3D NoUV NoShadows"
)

type Entity3D struct {
	NilRenderObject
	Name                string
	Model3D             *Model3D
	Visible             bool
	NotRelativeToCamera int
	RenderLast          bool
	DepthTesting        bool

	Shader     Shader
	RenderType RenderType

	Transform *TransformableObject3D
	transform TransformableObject
	parent    interface{}
}

func (this *Entity3D) commonInit() {
	this.Transform = &TransformableObject3D{}
	this.Transform.Scale = [3]float32{1.0, 1.0, 1.0}
	this.Transform.Rotation = mgl32.QuatRotate(0.0, mgl32.Vec3{0.0, 1.0, 0.0})

	this.transform = this.Transform

	this.Visible = true
	this.NotRelativeToCamera = -1
	this.RenderLast = false
	this.DepthTesting = true
	this.RenderType = TYPE_3D_NORMAL
	this.Shader = ResourceMgr.GetShader(ENTITY3D_SHADER_NAME)
	if this.Model3D != nil && !this.Model3D.HasUV() {
		this.Shader = ResourceMgr.GetShader(ENTITY3D_NO_UV_SHADER_NAME)
		if this.Shader == nil {
			ResourceMgr.LoadShaderSource(ENTITY3D_NO_UV_SHADER_NAME, ENTITY_3D_NOUV_SHADER_VERTEX_SOURCE_OPENGL, ENTITY_3D_NOUV_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
			this.Shader = ResourceMgr.GetShader(ENTITY3D_NO_UV_SHADER_NAME)
			if this.Shader == nil {
				ResourceMgr.LoadShaderSource(ENTITY3D_NO_UV_NO_SHADOWS_SHADER_NAME, ENTITY_3D_NOUV_SHADER_VERTEX_SOURCE_OPENGL, ENTITY_3D_NOUV_NO_SHADOWS_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
				this.Shader = ResourceMgr.GetShader(ENTITY3D_NO_UV_NO_SHADOWS_SHADER_NAME)
				if this.Shader != nil {
					ResourceMgr.SetShader(ENTITY3D_NO_UV_SHADER_NAME, ENTITY3D_NO_UV_NO_SHADOWS_SHADER_NAME)
				}
			}
		}
	}
}

func (this *Entity3D) InitName(name string) {
	this.Model3D = ResourceMgr.GetModel(name)
	this.Name = name
	this.commonInit()
}

func (this *Entity3D) InitMesh(mesh Mesh3D) {
	this.Model3D = &Model3D{
		Name: mesh.GetName(),
	}
	this.Model3D.AddMesh3D(mesh)
	this.Name = mesh.GetName()
	this.commonInit()
}

func (this *Entity3D) InitModel(model *Model3D) {
	this.Model3D = model
	if model != nil {
		this.Name = model.Name
	}
	this.commonInit()
}

func (this *Entity3D) InitLevel(level *Level) {
	if level != nil {
		if len(level.LevelObjects) != 0 {
			this.Model3D = level.LevelObjects[0].Model3D
			if this.Model3D != nil {
				this.Name = this.Model3D.Name
			}
		}
	}
	this.commonInit()
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
	if robj, ok := this.parent.(RenderObject); ok && !robj.IsVisible() {
		return false
	}
	return this.Visible
}

func (this *Entity3D) SetVisible() {
	this.Visible = true
}

func (this *Entity3D) SetInvisible() {
	this.Visible = false
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

func (this *Entity3D) GetTransform3D() *TransformableObject3D {
	return this.Transform
}

func (this *Entity3D) RendersLast() bool {
	return this.RenderLast
}

func (this *Entity3D) SetParent(parent interface{}) {
	this.parent = parent
	if tobj, ok := parent.(TweenableObject3D); ok {
		this.Transform.Parent = tobj
	}
}

func (this *Entity3D) GetParent() interface{} {
	return this.parent
}

func (this *Entity3D) HasDepthTesting() bool {
	return this.DepthTesting
}
