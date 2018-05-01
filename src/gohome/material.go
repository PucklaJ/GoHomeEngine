package gohome

import (
	// "fmt"
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
	// "math"
)

const (
	MATERIAL_UNIFORM_NAME                         string = "material"
	MATERIAL_DIFFUSE_COLOR_UNIFORM_NAME           string = "diffuseColor"
	MATERIAL_SPECULAR_COLOR_UNIFORM_NAME          string = "specularColor"
	MATERIAL_SPECULAR_TEXTURE_UNIFORM_NAME        string = "specularTexture"
	MATERIAL_DIFFUSE_TEXTURE_UNIFORM_NAME         string = "diffuseTexture"
	MATERIAL_SHINYNESS_UNIFORM_NAME               string = "shinyness"
	MATERIAL_DIFFUSE_TEXTURE_LOADED_UNIFORM_NAME  string = "DiffuseTextureLoaded"
	MATERIAL_SPECULAR_TEXTURE_LOADED_UNIFORM_NAME string = "SpecularTextureLoaded"
	MATERIAL_NORMALMAP_LOADED_UNIFORM_NAME        string = "NormalMapLoaded"
	MATERIAL_NORMALMAP_UNIFORM_NAME               string = "normalMap"
)

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (this *Color) ToVec3() mgl32.Vec3 {
	return mgl32.Vec3{
		float32(this.R) / 255.0,
		float32(this.G) / 255.0,
		float32(this.B) / 255.0,
	}
}

func (this *Color) ToVec4() mgl32.Vec4 {
	return mgl32.Vec4{
		float32(this.R) / 255.0,
		float32(this.G) / 255.0,
		float32(this.B) / 255.0,
		float32(this.A) / 255.0,
	}
}

func (this Color) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(float32(this.R) / 255.0 * float32(0xffff)), uint32(float32(this.G) / 255.0 * float32(0xffff)), uint32(float32(this.B) / 255.0 * float32(0xffff)), uint32(float32(this.A) / 255.0 * float32(0xffff))
}

type Material struct {
	Name          string
	DiffuseColor  color.Color
	SpecularColor color.Color

	DiffuseTexture  Texture
	SpecularTexture Texture
	NormalMap       Texture

	Shinyness float32

	DiffuseTextureLoaded  uint8
	SpecularTextureLoaded uint8
	NormalMapLoaded       uint8
}

func (mat *Material) InitDefault() {
	mat.Name = "Default"
	mat.DiffuseColor = &Color{255, 255, 255, 255}
	mat.SpecularColor = &Color{255, 255, 255, 255}

	mat.Shinyness = 0.5
}

func (mat *Material) SetTextures(diffuse, specular, normalMap string) {
	if diffuse != "" {
		mat.DiffuseTexture = ResourceMgr.GetTexture(diffuse)
	}
	if specular != "" {
		mat.SpecularTexture = ResourceMgr.GetTexture(specular)
	}
	if normalMap != "" {
		mat.NormalMap = ResourceMgr.GetTexture(normalMap)
	}
}

func (mat *Material) SetColors(diffuse, specular color.Color) {
	mat.DiffuseColor = diffuse
	mat.SpecularColor = specular
}

func ColorToVec3(c color.Color) mgl32.Vec3 {
	if c == nil {
		return [3]float32{0.0, 0.0, 0.0}
	}
	r, g, b, _ := c.RGBA()

	var vec mgl32.Vec3

	vec[0] = float32(r) / float32(0xffff)
	vec[1] = float32(g) / float32(0xffff)
	vec[2] = float32(b) / float32(0xffff)

	return vec
}

func ColorToVec4(c color.Color) mgl32.Vec4 {
	if c == nil {
		return [4]float32{0.0, 0.0, 0.0, 0.0}
	}
	r, g, b, a := c.RGBA()

	var vec mgl32.Vec4

	vec[0] = float32(r) / float32(0xffff)
	vec[1] = float32(g) / float32(0xffff)
	vec[2] = float32(b) / float32(0xffff)
	vec[3] = float32(a) / float32(0xffff)

	return vec
}
