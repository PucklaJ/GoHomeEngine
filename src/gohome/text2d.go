package gohome

import (
	"golang.org/x/image/colornames"
	"image/color"
	"strings"
)

// The values used for the flip
const (
	LINE_PADDING    = 0
	FLIP_NONE       = 0
	FLIP_HORIZONTAL = 1
	FLIP_VERTICAL   = 2
	FLIP_DIAGONALLY = 3
)

const (
	TEXT_2D_SHADER_NAME = "Text2D"

	COLOR_UNIFORM_NAME = "color"
)

// A text/string rendered to the screen
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
	// The transform of the object
	Transform            *TransformableObject2D

	// Wether this object is visible
	Visible             bool
	// The index of the camera to which this object is not relative to
	NotRelativeToCamera int
	// The size of the font
	FontSize            int
	// The text that will be displayed
	Text                string
	// The depth of the object (0-255)
	Depth               uint8
	// The color of the text
	Color               color.Color
	// The flip used for rendering
	Flip                uint8
}

// Initialises the object with a font a font size and a text
func (this *Text2D) Init(font string, fontSize int, str string) {
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

// Sets the shader of this object
func (this *Text2D) SetShader(s Shader) {
	this.shader = s
}

// Returns the shader of this object
func (this *Text2D) GetShader() Shader {
	return this.shader
}

// Sets the render type of this object
func (this *Text2D) SetType(rtype RenderType) {
	this.renderType = rtype
}

// Returns the render type of this object
func (this *Text2D) GetType() RenderType {
	return this.renderType
}

// Returns wether this object is visible
func (this *Text2D) IsVisible() bool {
	return this.Visible
}

// Returns the index of the camera to which this object is not relative to
func (this *Text2D) NotRelativeCamera() int {
	return this.NotRelativeToCamera
}

// Sets the font based on the name
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

		var width, height = 0, 0
		if len(this.textures) > 1 {
			for i := 0; i < len(this.textures); i++ {
				if this.textures[i] != nil {
					if this.textures[i].GetWidth() > width {
						width = this.textures[i].GetWidth()
					}
					height += this.textures[i].GetHeight() + LINE_PADDING
				} else {
					height += this.font.GetGlyphMaxHeight() + LINE_PADDING
					if this.font.GetGlyphMaxWidth()*100 > width {
						width = this.font.GetGlyphMaxWidth() * 100
					}
				}
			}
			this.renderTexture.ChangeSize(width, height)
			this.renderTexturesToRenderTexture()
		} else if len(this.textures) > 0 && this.textures[0] != nil {
			width = this.textures[0].GetWidth()
			height = this.textures[0].GetHeight()
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

	var y = 0
	for i := 0; i < len(this.textures); i++ {
		height := this.font.GetGlyphMaxHeight()
		if this.textures[i] != nil {
			var spr Sprite2D
			spr.InitTexture(this.textures[i])
			spr.Transform.Position = [2]float32{0, float32(y)}
			spr.NotRelativeToCamera = 0
			RenderMgr.RenderRenderObject(&spr)
			height = this.textures[i].GetHeight()
		}

		y += height + LINE_PADDING
	}

	this.renderTexture.UnsetAsTarget()
	RenderMgr.Projection2D = prevProj
}

func (this *Text2D) updateUniforms() {
	RenderMgr.prepareRenderRenderObject(this, -1)
}

// Sets the transformable object of this object
func (this *Text2D) SetTransformableObject(tobj TransformableObject) {
	this.transform = tobj
	if tobj != nil {
		this.Transform = tobj.(*TransformableObject2D)
	} else {
		this.Transform = nil
	}
}

// Returns the transformable object of this object
func (this *Text2D) GetTransformableObject() TransformableObject {
	return this.transform
}

// Cleans everything up (does not delete the font)
func (this *Text2D) Terminate() {
	this.texturesUsedDatabase = make(map[string]bool)
	this.deleteUnusedTexturesFromDatabase()
	if len(this.textures) > 0 {
		this.textures = this.textures[:0]
	}
}
