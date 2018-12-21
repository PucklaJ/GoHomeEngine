package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"sync"
)

const (
	ENTITY_3D_INSTANCED_SHADER_NAME                 string = "3D Instanced"
	ENTITY_3D_INSTANCED_NOUV_SHADER_NAME            string = "3D Instanced NoUV"
	ENTITY_3D_INSTANCED_NO_SHADOWS_SHADER_NAME      string = "3D Instanced NoShadows"
	ENTITY_3D_INSTANCED_NOUV_NO_SHADOWS_SHADER_NAME string = "3D Instanced NoUV NoShadows"
	ENTITY_3D_INSTANCED_SIMPLE_SHADER_NAME          string = "3D Instanced Simple"
)

type InstancedEntity3D struct {
	NilRenderObject
	Name                        string
	Model3D                     *InstancedModel3D
	Visible                     bool
	NotRelativeToCamera         int
	RenderLast                  bool
	DepthTesting                bool
	StopUpdatingInstancedValues bool

	Shader     Shader
	RenderType RenderType

	Transforms        []*TransformableObjectInstanced3D
	transformMatrices []mgl32.Mat4
	transformsChanged bool
}

func (this *InstancedEntity3D) commonInit(numInstances uint32) {
	this.Visible = true
	this.NotRelativeToCamera = -1
	this.RenderLast = false
	this.DepthTesting = true
	this.RenderType = TYPE_3D_INSTANCED | TYPE_CASTS_SHADOWS
	if this.Model3D == nil {
		return
	}

	this.Model3D.AddValueFront(VALUE_MAT4)
	this.Model3D.SetName(0, VALUE_MAT4, "transformMatrix3D")
	if this.Model3D.LoadedToGPU() {
		this.Model3D.SetNumInstances(numInstances)
	} else {
		this.Model3D.SetNumInstances(numInstances)
		this.Model3D.Load()
	}

	this.configureShader()

	this.Transforms = make([]*TransformableObjectInstanced3D, this.Model3D.GetNumInstances())
	this.transformMatrices = make([]mgl32.Mat4, this.Model3D.GetNumInstances())
	for i, t := range this.Transforms {
		this.Transforms[i] = &TransformableObjectInstanced3D{}
		t = this.Transforms[i]
		t.Scale = [3]float32{1.0, 1.0, 1.0}
		t.Rotation = mgl32.QuatRotate(0.0, mgl32.Vec3{0.0, 1.0, 0.0})
		t.SetTransformMatrixPointer(&this.transformMatrices[i])
	}
}

func (this *InstancedEntity3D) configureShaderFlags() uint32 {
	var flags uint32 = SHADER_FLAG_INSTANCED
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

func (this *InstancedEntity3D) configureShader() {
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

func (this *InstancedEntity3D) InitMesh(mesh InstancedMesh3D, numInstances uint32) {
	this.Model3D = &InstancedModel3D{
		Name: mesh.GetName(),
	}
	this.Model3D.AddMesh3D(mesh)
	this.Name = mesh.GetName()
	this.commonInit(numInstances)
}

func (this *InstancedEntity3D) InitModel(model *InstancedModel3D, numInstances uint32) {
	this.Model3D = model
	if model != nil {
		this.Name = model.Name
	}
	this.commonInit(numInstances)
}

func (this *InstancedEntity3D) GetTransformableObject() TransformableObject {
	return this
}

func (this *InstancedEntity3D) GetShader() Shader {
	return this.Shader
}

func (this *InstancedEntity3D) SetShader(s Shader) {
	this.Shader = s
}

func (this *InstancedEntity3D) SetType(rtype RenderType) {
	this.RenderType = rtype
}

func (this *InstancedEntity3D) GetType() RenderType {
	return this.RenderType
}

func (this *InstancedEntity3D) Render() {
	if !this.StopUpdatingInstancedValues {
		this.UpdateInstancedValues()
	}
	if this.Model3D != nil {
		if this.transformsChanged {
			this.Model3D.SetM4(0, this.transformMatrices)
			this.transformsChanged = false
		}
		this.Model3D.Render()
	}
}

func (this *InstancedEntity3D) Terminate() {
	if this.Model3D != nil {
		this.Model3D.Terminate()
	}
}

func (this *InstancedEntity3D) IsVisible() bool {
	return this.Visible
}

func (this *InstancedEntity3D) SetVisible() {
	this.Visible = true
}

func (this *InstancedEntity3D) SetInvisible() {
	this.Visible = false
}

func (this *InstancedEntity3D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}

func (this *InstancedEntity3D) RendersLast() bool {
	return this.RenderLast
}

func (this *InstancedEntity3D) HasDepthTesting() bool {
	return this.DepthTesting
}

func (this *InstancedEntity3D) UpdateInstancedValues() {
	this.CalculateTransformMatrix(&RenderMgr, this.NotRelativeToCamera)
}

func (this *InstancedEntity3D) SetNumInstances(n uint32) {
	prev := this.Model3D.GetNumInstances()
	this.Model3D.SetNumInstances(n)
	if prev != n {
		if n > prev {
			this.Transforms = append(this.Transforms, make([]*TransformableObjectInstanced3D, n-prev)...)
			this.transformMatrices = append(this.transformMatrices, make([]mgl32.Mat4, n-prev)...)
			for i := uint32(0); i < n; i++ {
				if i >= prev {
					this.Transforms[i] = &TransformableObjectInstanced3D{}
					this.Transforms[i].Scale = [3]float32{1.0, 1.0, 1.0}
					this.Transforms[i].Rotation = mgl32.QuatRotate(0.0, [3]float32{0.0, 1.0, 0.0})
				}
				this.Transforms[i].SetTransformMatrixPointer(&this.transformMatrices[i])
			}
		} else {
			this.Transforms = this.Transforms[:n]
			this.transformMatrices = this.transformMatrices[:n]
		}
		if !this.StopUpdatingInstancedValues {
			this.UpdateInstancedValues()
		}
	}
}

func (this *InstancedEntity3D) SetNumUsedInstances(n uint32) {
	this.Model3D.SetNumUsedInstances(n)
}

func (this *InstancedEntity3D) CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int) {
	var wg sync.WaitGroup
	wg.Add(int(this.Model3D.GetNumInstances()))
	for _, t := range this.Transforms {
		go func(_t *TransformableObjectInstanced3D) {
			if _t.valuesChanged() {
				this.transformsChanged = true
			}
			_t.CalculateTransformMatrix(&RenderMgr, notRelativeToCamera)
			wg.Done()
		}(t)
	}
	wg.Wait()
}

func (this *InstancedEntity3D) SetTransformMatrix(rmgr *RenderManager) {

}
