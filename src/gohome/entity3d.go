package gohome

import (
	"github.com/PucklaJ/mathgl/mgl32"
)

const (
	ENTITY3D_SHADER_NAME                  string = "3D"
	ENTITY3D_NO_UV_SHADER_NAME            string = "3D NoUV"
	ENTITY3D_NO_UV_NO_SHADOWS_SHADER_NAME string = "3D NoUV NoShadows"
)

// A 3D RenderObject with a 3D Model
type Entity3D struct {
	NilRenderObject
	// The name of the Entity
	Name string
	// The 3D Model of the Entity
	Model3D *Model3D
	// Wether it is visible
	Visible bool
	// The index of the camera to which it is not relative
	// or -1 if it relative to every camera
	NotRelativeToCamera int
	// Wether it should render after everyting else
	RenderLast bool
	// Wether the depth test is enabled
	DepthTesting bool

	// The shader that will be used on this 3D Model
	Shader Shader
	// The render type of the Entity
	RenderType RenderType

	// The transform of the Entity
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
	this.RenderType = TYPE_3D_NORMAL | TYPE_CASTS_SHADOWS
	this.configureShader()
}

func (this *Entity3D) configureShaderFlags() uint32 {
	var flags uint32 = 0
	if !this.Model3D.HasUV() {
		flags |= SHADER_FLAG_NOUV
	}
	if LightMgr.CurrentLightCollection == -1 {
		flags |= SHADER_FLAG_NO_LIGHTING
	}
	if this.Model3D.HasUV() {
		var hasDif, hasSpec, hasNorm = false, false, false
		for i := 0; i < len(this.Model3D.meshes); i++ {
			m := this.Model3D.meshes[i]
			mat := m.GetMaterial()
			if mat.DiffuseColor != nil {
				hasDif = true
			}
			if mat.SpecularTexture != nil {
				hasSpec = true
			}
			if mat.NormalMap != nil {
				hasNorm = true
			}
		}
		if !hasDif {
			flags |= SHADER_FLAG_NO_DIFTEX
		}
		if !hasSpec {
			flags |= SHADER_FLAG_NO_SPECTEX
		}
		if !hasNorm {
			flags |= SHADER_FLAG_NO_NORMAP
		}
	}

	return flags
}

func (this *Entity3D) configureShader() {
	if this.Model3D == nil {
		return
	}
	flags := this.configureShaderFlags()
	name := GetShaderName3D(flags)
	this.Shader = ResourceMgr.GetShader(name)
	if this.Shader == nil {
		LoadGeneratedShader3D(SHADER_TYPE_3D, flags)
		this.Shader = ResourceMgr.GetShader(name)
		if this.Shader == nil {
			flags |= SHADER_FLAG_NO_SHADOWS
			name = GetShaderName3D(flags)
			this.Shader = ResourceMgr.GetShader(name)
			if this.Shader == nil {
				LoadGeneratedShader3D(SHADER_TYPE_3D, flags)
				this.Shader = ResourceMgr.GetShader(name)
			}
		}
	}
}

// Initialises the Entity with the 3D Model name
func (this *Entity3D) InitName(name string) {
	this.Model3D = ResourceMgr.GetModel(name)
	this.Name = name
	this.commonInit()
}

// Initialises the Entity with mesh
func (this *Entity3D) InitMesh(mesh Mesh3D) {
	this.Model3D = &Model3D{
		Name: mesh.GetName(),
	}
	this.Model3D.AddMesh3D(mesh)
	this.Name = mesh.GetName()
	this.commonInit()
}

// Initialises the Entity with model
func (this *Entity3D) InitModel(model *Model3D) {
	this.Model3D = model
	if model != nil {
		this.Name = model.Name
	}
	this.commonInit()
}

// Initialises the Entity using the first model of level
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

// Returns the shader of this Entity
func (this *Entity3D) GetShader() Shader {
	return this.Shader
}

// Sets the shader of this entity
func (this *Entity3D) SetShader(s Shader) {
	this.Shader = s
}

// Sets the render type of this entity
func (this *Entity3D) SetType(rtype RenderType) {
	this.RenderType = rtype
}

// Returns the render type of this entity
func (this *Entity3D) GetType() RenderType {
	return this.RenderType
}

// Renders the entity (a lot of values need to be set up before
// calling this method, use RenderMgr.RenderRenderObject if you want
// to render a specific RenderObject)
func (this *Entity3D) Render() {
	if this.Model3D != nil {
		this.Model3D.Render()
	}
}

// Cleans up the 3D Model
func (this *Entity3D) Terminate() {
	if this.Model3D != nil {
		this.Model3D.Terminate()
	}
}

// Returns wether the Entity is visible
func (this *Entity3D) IsVisible() bool {
	if robj, ok := this.parent.(RenderObject); ok && !robj.IsVisible() {
		return false
	}
	return this.Visible
}

// Sets the Entity to be visible
func (this *Entity3D) SetVisible() {
	this.Visible = true
}

// Sets the Entity to be invisible
func (this *Entity3D) SetInvisible() {
	this.Visible = false
}

// Returns the index to which camera it is not relative
// or -1 if it is relative to every camera
func (this *Entity3D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}

// Sets the transformable object of this Entity
func (this *Entity3D) SetTransformableObject(tobj TransformableObject) {
	this.transform = tobj
	if tobj != nil {
		this.Transform = tobj.(*TransformableObject3D)
	} else {
		this.Transform = nil
	}
}

// Returns the tranformable object of the Entity
func (this *Entity3D) GetTransformableObject() TransformableObject {
	return this.transform
}

// Returns the Transform of the Entity
func (this *Entity3D) GetTransform3D() *TransformableObject3D {
	return this.Transform
}

// Used for calculating the transformation matrices in go routines
func (this *Entity3D) SetChildChannel(channel chan bool, tobj *TransformableObject3D) {
	this.Transform.SetChildChannel(channel, tobj)
}

// Wether this Entity renders last
func (this *Entity3D) RendersLast() bool {
	return this.RenderLast
}

// Sets the parent of this Entity
func (this *Entity3D) SetParent(parent interface{}) {
	this.parent = parent
	if tobj, ok := parent.(ParentObject3D); ok {
		this.Transform.SetParent(tobj)
	} else {
		this.Transform.SetParent(nil)
	}
}

// Returns the parent of this entity
func (this *Entity3D) GetParent() interface{} {
	return this.parent
}

// Returns wether the depth test is enabled on this object
func (this *Entity3D) HasDepthTesting() bool {
	return this.DepthTesting
}
