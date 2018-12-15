package gohome

const (
	SPRITE2D_SHADER_NAME               string = "2D"
	SPRITE2D_MESH_NAME                 string = "SPRITE2D_MESH"
	FLIP_UNIFORM_NAME                  string = "flip"
	TEXTURE_REGION_UNIFORM_NAME        string = "textureRegion"
	DEPTH_UNIFORM_NAME                 string = "depth"
	ENABLE_KEY_UNIFORM_NAME            string = "enableKey"
	KEY_COLOR_UNIFORM_NAME             string = "keyColor"
	ENABLE_MOD_UNIFORM_NAME            string = "enableMod"
	MOD_COLOR_UNIFORM_NAME             string = "modColor"
	ENABLE_TEXTURE_REGION_UNIFORM_NAME string = "enableTextureRegion"
)

var sprite2DMesh Mesh2D = nil

type Sprite2D struct {
	NilRenderObject
	Texture             Texture
	Visible             bool
	NotRelativeToCamera int
	Flip                uint8
	Name                string

	Shader        Shader
	RenderType    RenderType
	transform     TransformableObject
	Transform     *TransformableObject2D
	TextureRegion TextureRegion
	Depth         uint8
}

func createSprite2DMesh() {
	sprite2DMesh = Render.CreateMesh2D(SPRITE2D_MESH_NAME)

	vertices := []Mesh2DVertex{
		/*X,Y
		  U,V
		*/
		{-0.5, -0.5, // LEFT-DOWN
			0.0, 0.0},

		{0.5, -0.5, // RIGHT-DOWN
			1.0, 0.0},

		{0.5, 0.5, // RIGHT-UP
			1.0, 1.0},

		{-0.5, 0.5, // LEFT-UP
			0.0, 1.0},
	}

	indices := []uint32{
		0, 3, 2, // LEFT-TRI
		2, 1, 0, // RIGHT-TRI
	}

	sprite2DMesh.AddVertices(vertices, indices)
	sprite2DMesh.Load()
}

func (spr *Sprite2D) commonInit() {
	spr.Transform = &TransformableObject2D{}

	spr.Transform.Scale = [2]float32{1.0, 1.0}
	if spr.Texture != nil {
		spr.Transform.Size = [2]float32{float32(spr.Texture.GetWidth()), float32(spr.Texture.GetHeight())}
		spr.TextureRegion.Min = [2]float32{0.0, 0.0}
		spr.TextureRegion.Max = spr.Transform.Size
		spr.Name = spr.Texture.GetName()
	}
	spr.Transform.RotationPoint = [2]float32{0.5, 0.5}
	spr.Transform.Origin = [2]float32{0.0, 0.0}

	spr.transform = spr.Transform

	if sprite2DMesh == nil {
		createSprite2DMesh()
	}

	spr.Visible = true
	spr.NotRelativeToCamera = -1
	spr.RenderType = TYPE_2D_NORMAL
	spr.Flip = FLIP_NONE
	spr.Shader = ResourceMgr.GetShader(SPRITE2D_SHADER_NAME)
	if spr.Shader == nil {
		spr.Shader = LoadGeneratedShader2D(SHADER_TYPE_SPRITE2D, 0)
	}
}

func (spr *Sprite2D) Init(texName string) {
	spr.Texture = ResourceMgr.GetTexture(texName)
	spr.Name = texName
	spr.commonInit()
}

func (spr *Sprite2D) InitTexture(texture Texture) {
	spr.Texture = texture
	spr.commonInit()
}

func (spr *Sprite2D) SetShader(s Shader) {
	spr.Shader = s
}

func (spr *Sprite2D) GetShader() Shader {
	if spr.Shader == nil {
		spr.Shader = ResourceMgr.GetShader(SPRITE2D_SHADER_NAME)
	}
	return spr.Shader
}

func (spr *Sprite2D) SetType(rtype RenderType) {
	spr.RenderType = rtype
}

func (spr *Sprite2D) GetType() RenderType {
	return spr.RenderType
}

func convertDepth(depth uint8) float32 {
	return (1.0-float32(depth)/255.0)*2.0 - 1.0
}

func (spr *Sprite2D) setUniforms() {
	shader := RenderMgr.CurrentShader

	if shader != nil {
		shader.SetUniformI(FLIP_UNIFORM_NAME, int32(spr.Flip))
		shader.SetUniformV4(TEXTURE_REGION_UNIFORM_NAME, spr.TextureRegion.Normalize(spr.Texture).Vec4())
		shader.SetUniformF(DEPTH_UNIFORM_NAME, convertDepth(spr.Depth))
		if spr.Texture.GetKeyColor() != nil {
			shader.SetUniformI(ENABLE_KEY_UNIFORM_NAME, 1)
			shader.SetUniformV3(KEY_COLOR_UNIFORM_NAME, ColorToVec3(spr.Texture.GetKeyColor()))
		} else {
			shader.SetUniformI(ENABLE_KEY_UNIFORM_NAME, 0)
		}
		if spr.Texture.GetModColor() != nil {
			shader.SetUniformI(ENABLE_MOD_UNIFORM_NAME, 1)
			shader.SetUniformV4(MOD_COLOR_UNIFORM_NAME, ColorToVec4(spr.Texture.GetModColor()))
		} else {
			shader.SetUniformI(ENABLE_MOD_UNIFORM_NAME, 0)
		}
		shader.SetUniformI(ENABLE_TEXTURE_REGION_UNIFORM_NAME, 1)
	}
}

func (spr *Sprite2D) Render() {
	if spr.Texture != nil {
		spr.setUniforms()
		spr.Texture.Bind(0)
		sprite2DMesh.Render()
		spr.Texture.Unbind(0)
	} else {
		ErrorMgr.Error("Sprite2D", spr.Name, "Couldn't render: The Texture is nil")
	}
}

func (spr *Sprite2D) Terminate() {
	if spr.Texture != nil {
		spr.Texture.Terminate()
	}
}

func (spr *Sprite2D) IsVisible() bool {
	return spr.Visible
}

func (spr *Sprite2D) SetVisible() {
	spr.Visible = true
}

func (spr *Sprite2D) SetInvisible() {
	spr.Visible = false
}

func (spr *Sprite2D) NotRelativeCamera() int {
	return spr.NotRelativeToCamera
}

func (spr *Sprite2D) SetTransformableObject(tobj TransformableObject) {
	spr.transform = tobj
	if tobj != nil {
		spr.Transform = tobj.(*TransformableObject2D)
	} else {
		spr.Transform = nil
	}
}

func (spr *Sprite2D) GetTransformableObject() TransformableObject {
	return spr.transform
}

func (spr *Sprite2D) GetTransform2D() *TransformableObject2D {
	return spr.Transform
}
