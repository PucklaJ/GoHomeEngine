package gohome

const (
	SPRITE2D_SHADER_NAME string = "2D"
	SPRITE2D_MESH_NAME   string = "SPRITE2D_MESH"
)

var sprite2DMesh Mesh2D = nil

type Sprite2D struct {
	Texture
	Visible             bool
	NotRelativeToCamera int
	Flip                uint8

	Shader
	RenderType
	transform TransformableObject
	Transform *TransformableObject2D
}

func createSprite2DMesh() {
	sprite2DMesh = Render.CreateMesh2D(SPRITE2D_MESH_NAME)

	vertices := []Mesh2DVertex{
		/*X,Y
		  U,V
		*/
		Mesh2DVertex{-0.5, 0.5, // LEFT-DOWN
			0.0, 1.0},

		Mesh2DVertex{0.5, 0.5, // RIGHT-DOWN
			1.0, 1.0},

		Mesh2DVertex{0.5, -0.5, // RIGHT-UP
			1.0, 0.0},

		Mesh2DVertex{-0.5, -0.5, // LEFT-UP
			0.0, 0.0},
	}

	indices := []uint32{
		0, 1, 2, // LEFT-TRI
		2, 3, 0, // RIGHT-TRI
	}

	sprite2DMesh.AddVertices(vertices, indices)
	sprite2DMesh.Load()
}

func (spr *Sprite2D) commonInit() {
	spr.Transform = &TransformableObject2D{}

	spr.Transform.Scale = [2]float32{1.0, 1.0}
	if spr.Texture != nil {
		spr.Transform.Size = [2]float32{float32(spr.Texture.GetWidth()), float32(spr.Texture.GetHeight())}
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
}

func (spr *Sprite2D) Init(texName string) {
	spr.Texture = ResourceMgr.GetTexture(texName)
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

func (spr *Sprite2D) Render() {
	if spr.Texture != nil {
		RenderMgr.CurrentShader.SetUniformI("flip", int32(spr.Flip))
		spr.Texture.Bind(0)
		sprite2DMesh.Render()
		spr.Texture.Unbind(0)
	}
}

func (spr *Sprite2D) Terminate() {
	spr.Texture.Terminate()
}

func (spr *Sprite2D) IsVisible() bool {
	return spr.Visible
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
