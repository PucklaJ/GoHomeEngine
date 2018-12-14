package gohome

import (
	"golang.org/x/image/colornames"
	"image/color"
	"strings"
)

const (
	LINE_PADDING    int32 = 0
	FLIP_NONE       uint8 = 0
	FLIP_HORIZONTAL uint8 = 1
	FLIP_VERTICAL   uint8 = 2
	FLIP_DIAGONALLY uint8 = 3
)

const (
	TEXT_2D_SHADER_NAME string = "Text2D"

	COLOR_UNIFORM_NAME string = "color"
)

type Text2D struct {
	NilRenderObject
	shader               Shader
	renderType           RenderType
	font                 *Font
	textures             []Texture
	textureDatabase      map[string]Texture
	texturesUsedDatabase map[string]bool
	renderTexture        RenderTexture
	oldText              string
	transform            TransformableObject
	Transform            *TransformableObject2D

	Visible             bool
	NotRelativeToCamera int
	FontSize            uint32
	Text                string
	Depth               uint8
	Color               color.Color
	Flip                uint8
}

func (this *Text2D) Init(font string, fontSize uint32, str string) {
	this.font = ResourceMgr.GetFont(font)
	this.Transform = &TransformableObject2D{}
	this.Transform.Scale = [2]float32{1.0, 1.0}
	this.Transform.RotationPoint = [2]float32{0.5, 0.5}
	this.Transform.Origin = [2]float32{0.0, 0.0}
	this.transform = this.Transform
	this.textureDatabase = make(map[string]Texture)
	this.texturesUsedDatabase = make(map[string]bool)
	this.renderTexture = Render.CreateRenderTexture("Text2D RenderTexture", 10, 10, 1, false, false, false, false)

	if sprite2DMesh == nil {
		createSprite2DMesh()
	}

	this.Visible = true
	this.NotRelativeToCamera = -1
	this.FontSize = fontSize
	this.Text = str
	if this.shader = ResourceMgr.GetShader(GetShaderName2D(SHADER_TYPE_TEXT2D, 0)); this.shader == nil {
		this.shader = LoadGeneratedShader2D(SHADER_TYPE_TEXT2D, 0)
	}
	this.Depth = 0
	this.Color = colornames.White
	this.Flip = FLIP_NONE
	this.renderType = TYPE_2D_NORMAL

	this.updateText()
}

func (this *Text2D) setUniforms() {
	shader := RenderMgr.CurrentShader
	if shader != nil {
		var flip uint8
		if len(this.textures) == 1 {
			flip = FLIP_NONE
		} else {
			flip = FLIP_VERTICAL
		}
		switch flip {
		case FLIP_NONE:
			flip = this.Flip
		case FLIP_VERTICAL:
			switch this.Flip {
			case FLIP_HORIZONTAL:
				flip = FLIP_DIAGONALLY
			case FLIP_VERTICAL:
				flip = FLIP_NONE
			case FLIP_DIAGONALLY:
				flip = FLIP_HORIZONTAL
			}
		}
		shader.SetUniformI(FLIP_UNIFORM_NAME, int32(flip))
		shader.SetUniformF(DEPTH_UNIFORM_NAME, convertDepth(this.Depth))
		shader.SetUniformV4(COLOR_UNIFORM_NAME, ColorToVec4(this.Color))
	}
}

func (this *Text2D) getTexture() Texture {
	var temp Texture
	if len(this.textures) == 1 {
		temp = this.textures[0]
	} else {
		temp = this.renderTexture
	}
	return temp
}

func (this *Text2D) Render() {
	this.updateText()
	temp := this.getTexture()
	if temp != nil {
		this.setUniforms()
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

func (this *Text2D) deleteUnusedTexturesFromDatabase() {
	var texturesToDelete []string
	for k := range this.textureDatabase {
		if used, ok := this.texturesUsedDatabase[k]; !ok || !used {
			texturesToDelete = append(texturesToDelete, k)
		}
	}

	for i := 0; i < len(texturesToDelete); i++ {
		this.textureDatabase[texturesToDelete[i]].Terminate()
		delete(this.textureDatabase, texturesToDelete[i])
	}
}

func (this *Text2D) updateText() {
	this.texturesUsedDatabase = make(map[string]bool)

	if this.font == nil {
		return
	}

	if this.valuesChanged() {
		defer this.deleteUnusedTexturesFromDatabase()
		defer func() {
			this.oldText = this.Text
		}()

		if this.renderTexture != nil {
			this.renderTexture.Terminate()
		}

		if len(this.textures) > 0 {
			this.textures = this.textures[:0]
		}

		this.font.FontSize = this.FontSize

		if this.Text == "" {
			return
		}

		lines := strings.Split(this.Text, "\n")
		if len(lines) == 0 {
			return
		}
		var tempTexture Texture
		for i := 0; i < len(lines); i++ {
			if lines[i] != "" {
				texturedb, ok := this.textureDatabase[lines[i]]
				if ok && texturedb != nil {
					tempTexture = texturedb
					this.texturesUsedDatabase[lines[i]] = true
				} else {
					tempTexture = this.font.DrawString(lines[i])
					this.textureDatabase[lines[i]] = tempTexture
					this.texturesUsedDatabase[lines[i]] = true
				}

				this.textures = append(this.textures, tempTexture)
			} else {
				this.textures = append(this.textures, nil)
			}
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
					height += this.font.GetGlyphMaxHeight() + uint32(LINE_PADDING)
					if this.font.GetGlyphMaxWidth()*100 > width {
						width = this.font.GetGlyphMaxWidth() * 100
					}
				}
			}
			this.renderTexture.ChangeSize(width, height)
			this.renderTexturesToRenderTexture()
		} else if len(this.textures) > 0 && this.textures[0] != nil {
			width = uint32(this.textures[0].GetWidth())
			height = uint32(this.textures[0].GetHeight())
		} else {
			width = this.font.GetGlyphMaxWidth() * 100
			height = this.font.GetGlyphMaxHeight()
		}

		if this.Transform != nil {
			this.Transform.Size = [2]float32{float32(width), float32(height)}
		}

		this.updateUniforms()
	}
}

func (this *Text2D) renderTexturesToRenderTexture() {
	this.renderTexture.SetAsTarget()
	prevProj := RenderMgr.Projection2D
	RenderMgr.SetProjection2DToTexture(this.renderTexture)

	var x, y uint32 = 0, 0
	for i := 0; i < len(this.textures); i++ {
		height := this.font.GetGlyphMaxHeight()
		if this.textures[i] != nil {
			var spr Sprite2D
			spr.InitTexture(this.textures[i])
			spr.Transform.Position = [2]float32{float32(x), float32(y)}
			spr.NotRelativeToCamera = 0
			RenderMgr.RenderRenderObject(&spr)
			height = uint32(this.textures[i].GetHeight())
		}

		y += uint32(int32(height) + LINE_PADDING)
	}

	this.renderTexture.UnsetAsTarget()
	RenderMgr.Projection2D = prevProj
}

func (this *Text2D) updateUniforms() {
	RenderMgr.prepareRenderRenderObject(this, -1)
}

func (this *Text2D) SetTransformableObject(tobj TransformableObject) {
	this.transform = tobj
	if tobj != nil {
		this.Transform = tobj.(*TransformableObject2D)
	} else {
		this.Transform = nil
	}
}

func (this *Text2D) GetTransformableObject() TransformableObject {
	return this.transform
}

func (this *Text2D) Terminate() {
	this.texturesUsedDatabase = make(map[string]bool)
	this.deleteUnusedTexturesFromDatabase()
	if len(this.textures) > 0 {
		this.textures = this.textures[:0]
	}
}
