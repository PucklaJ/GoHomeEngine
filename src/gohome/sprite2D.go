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

	Shader
	RenderType
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

func (spr *Sprite2D) Init(texName string, transform *TransformableObject2D) {
	spr.Texture = ResourceMgr.GetTexture(texName)

	if spr.Texture != nil {
		transform.Scale = [2]float32{1.0, 1.0}
		transform.Size = [2]float32{float32(spr.Texture.GetWidth()), float32(spr.Texture.GetHeight())}
		transform.RotationPoint = [2]float32{0.5, 0.5}
	}

	if sprite2DMesh == nil {
		createSprite2DMesh()
	}

	spr.Visible = true
	spr.NotRelativeToCamera = -1
	spr.RenderType = TYPE_2D_NORMAL
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
	spr.Texture.Bind(0)
	sprite2DMesh.Render()
	spr.Texture.Unbind(0)
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
