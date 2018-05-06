package gohome

type Text2D struct {
	shader     Shader
	renderType RenderType
	font       *Font
	texture    Texture
	oldText    string
	transform  *TransformableObject2D

	Visible             bool
	NotRelativeToCamera int
	FontSize            uint32
	Text                string
}

func (this *Text2D) Init(font string, fontSize uint32, str string, transform *TransformableObject2D) {
	this.font = ResourceMgr.GetFont(font)
	if transform != nil {
		transform.Scale = [2]float32{1.0, 1.0}
		transform.RotationPoint = [2]float32{0.5, 0.5}
		transform.Origin = [2]float32{0.0, 0.0}
		this.transform = transform
	}

	if sprite2DMesh == nil {
		createSprite2DMesh()
	}

	this.Visible = true
	this.NotRelativeToCamera = -1
	this.FontSize = fontSize
	this.Text = str
	this.shader = ResourceMgr.GetShader(SPRITE2D_SHADER_NAME)

	this.updateText()
}

func (this *Text2D) Render() {
	this.updateText()
	if this.texture != nil {
		this.texture.Bind(0)
		sprite2DMesh.Render()
		this.texture.Unbind(0)
	}
}
func (this *Text2D) SetShader(s Shader) {
	this.shader = s
}
func (this *Text2D) GetShader() Shader {
	return this.shader
}
func (this *Text2D) SetType(rtype RenderType) {
	this.renderType = rtype
}
func (this *Text2D) GetType() RenderType {
	return this.renderType
}
func (this *Text2D) IsVisible() bool {
	return this.Visible
}
func (this *Text2D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}
func (this *Text2D) SetFont(name string) {
	this.font = ResourceMgr.GetFont(name)
}

func (this *Text2D) valuesChanged() bool {
	return this.Text != this.oldText
}

func (this *Text2D) updateText() {
	if this.font == nil {
		return
	}

	if this.valuesChanged() {
		if this.texture != nil {
			this.texture.Terminate()
		}

		this.font.FontSize = this.FontSize
		this.texture = this.font.DrawString(this.Text)
		if this.transform != nil && this.texture != nil {
			this.transform.Size = [2]float32{float32(this.texture.GetWidth()), float32(this.texture.GetHeight())}
		}
	}

	this.oldText = this.Text
}
