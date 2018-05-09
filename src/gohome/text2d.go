package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

const (
	LINE_PADDING    int32 = 0
	FLIP_NONE       uint8 = 0
	FLIP_HORIZONTAL uint8 = 1
	FLIP_VERTICAL   uint8 = 2
	FLIP_DIAGONALLY uint8 = 3
)

type Text2D struct {
	shader        Shader
	renderType    RenderType
	font          *Font
	textures      []Texture
	renderTexture RenderTexture
	oldText       string
	transform     *TransformableObject2D

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
	var temp Texture
	if len(this.textures) == 1 {
		temp = this.textures[0]
		RenderMgr.CurrentShader.SetUniformI("flip", int32(FLIP_NONE))
	} else {
		temp = this.renderTexture
		RenderMgr.CurrentShader.SetUniformI("flip", int32(FLIP_VERTICAL))
	}
	if temp != nil {
		temp.Bind(0)
		sprite2DMesh.Render()
		temp.Unbind(0)
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
		for i := 0; i < len(this.textures); i++ {
			if this.textures[i] != nil {
				this.textures[i].Terminate()
			}
		}
		if len(this.textures) != 0 {
			this.textures = this.textures[len(this.textures)-1 : 0]
		}
		if this.renderTexture != nil {
			this.renderTexture.Terminate()
		}

		this.font.FontSize = this.FontSize
		lines := strings.Split(this.Text, "\n")
		for i := 0; i < len(lines); i++ {
			this.textures = append(this.textures, this.font.DrawString(lines[i]))
		}

		var width, height uint32 = 0, 0
		if len(this.textures) > 1 {
			for i := 0; i < len(this.textures); i++ {
				if this.textures[i] != nil {
					if uint32(this.textures[i].GetWidth()) > width {
						width = uint32(this.textures[i].GetWidth())
					}
					height += uint32(int32(this.textures[i].GetHeight()) + LINE_PADDING)
				} else {
					height += uint32(64 + LINE_PADDING)
					if 1000 > width {
						width = 1000
					}
				}
			}
			this.renderTexture = Render.CreateRenderTexture("Text2D RenderTexture", width, height, 1, false, false, false, false)
			this.renderTexturesToRenderTexture()
		} else if len(this.textures) > 0 && this.textures[0] != nil {
			width = uint32(this.textures[0].GetWidth())
			height = uint32(this.textures[0].GetHeight())
		} else {
			width = 1000
			height = 64
		}

		if this.transform != nil {
			this.transform.Size = [2]float32{float32(width), float32(height)}
		}
	}

	this.oldText = this.Text
}

func (this *Text2D) renderTexturesToRenderTexture() {
	shader := ResourceMgr.GetShader(SPRITE2D_SHADER_NAME)
	shader.Use()
	this.renderTexture.SetAsTarget()
	var projection Ortho2DProjection
	projection.Right = float32(this.renderTexture.GetWidth())
	projection.Bottom = float32(this.renderTexture.GetHeight())
	projection.Left = 0.0
	projection.Top = 0.0
	projection.CalculateProjectionMatrix()
	projectionMatrix := projection.GetProjectionMatrix()
	viewMatrix := mgl32.Ident3()

	shader.SetUniformM4("projectionMatrix2D", projectionMatrix)
	shader.SetUniformM3("viewMatrix2D", viewMatrix)
	shader.SetUniformI("texture0", 0)

	var x, y uint32 = 0, 0
	for i := 0; i < len(this.textures); i++ {
		var width, height uint32 = 1000, 64
		if this.textures[i] != nil {
			width = uint32(this.textures[i].GetWidth())
			height = uint32(this.textures[i].GetHeight())
			this.textures[i].Bind(0)
		}
		var transformMatrix TransformableObject2D
		transformMatrix.Size = [2]float32{float32(width), float32(height)}
		transformMatrix.Scale = [2]float32{1.0, 1.0}
		transformMatrix.Origin = [2]float32{0.0, 0.0}
		transformMatrix.RotationPoint = [2]float32{0.5, 0.5}
		transformMatrix.Position = [2]float32{float32(x), float32(y)}
		transformMatrix.CalculateTransformMatrix(nil, -1)
		shader.SetUniformM3("transformMatrix2D", transformMatrix.GetTransformMatrix())

		sprite2DMesh.Render()

		y += uint32(int32(height) + LINE_PADDING)
	}

	shader.Unuse()
	this.renderTexture.UnsetAsTarget()

	shader = RenderMgr.CurrentShader
	if shader != nil {
		shader.Use()
		if RenderMgr.Projection2D != nil {
			shader.SetUniformM4("projectionMatrix2D", RenderMgr.Projection2D.GetProjectionMatrix())
		} else {
			shader.SetUniformM4("projectionMatrix2D", mgl32.Ident4())
		}
		if this.transform != nil {
			shader.SetUniformM3("transformMatrix2D", this.transform.GetTransformMatrix())
		} else {
			shader.SetUniformM3("transformMatrix2D", mgl32.Ident3())
		}
		cam := RenderMgr.currentCamera2D
		if cam != nil {
			shader.SetUniformM3("viewMatrix2D", cam.GetViewMatrix())
		} else {
			shader.SetUniformM3("viewMatrix2D", mgl32.Ident3())
		}
	}

}
