package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"image/color"
	"math"
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
	MATERIAL_TRANSPARENCY_UNIFORM_NAME            string = "transparency"

	MAX_SPECULAR_EXPONENT float64 = 50.0
	MIN_SPECULAR_EXPONENT float64 = 5.0
)

// A RGBA color
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// Converts a color to a Vec3 using its RGB values
func (this *Color) ToVec3() mgl32.Vec3 {
	return mgl32.Vec3{
		float32(this.R) / 255.0,
		float32(this.G) / 255.0,
		float32(this.B) / 255.0,
	}
}

// Converts a color to a Vec4 using its RGBA values
func (this *Color) ToVec4() mgl32.Vec4 {
	return mgl32.Vec4{
		float32(this.R) / 255.0,
		float32(this.G) / 255.0,
		float32(this.B) / 255.0,
		float32(this.A) / 255.0,
	}
}

// Returns the RGBA values as uint32 (0-255)
func (this Color) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(float32(this.R) / 255.0 * float32(0xffff)), uint32(float32(this.G) / 255.0 * float32(0xffff)), uint32(float32(this.B) / 255.0 * float32(0xffff)), uint32(float32(this.A) / 255.0 * float32(0xffff))
}

// A Material having properties that define the look of 3D geometry
type Material struct {
	// The name of the material
	Name          string
	// The diffuse color of the material
	DiffuseColor  color.Color
	// The specular color of the material
	SpecularColor color.Color

	// The diffuse texture of the material
	DiffuseTexture  Texture
	// The specular texture of the material
	SpecularTexture Texture
	// The normal map of the material
	NormalMap       Texture

	// Defines how much specular light should be applied (0.0-1.0)
	Shinyness    float32
	// The transparency or alpha value of the material
	Transparency float32

	// Used to tell the shader if this material has a diffuse texture
	DiffuseTextureLoaded  uint8
	// Used to tell the shader if this material has a specular texture
	SpecularTextureLoaded uint8
	// Used to tell the shader if this material has a normal map
	NormalMapLoaded       uint8
}

// Initialises some default values
func (mat *Material) InitDefault() {
	mat.Name = "Default"
	mat.DiffuseColor = &Color{255, 255, 255, 255}
	mat.SpecularColor = &Color{255, 255, 255, 255}

	mat.Shinyness = 0.5
	mat.Transparency = 1.0
}

// Loads the textures from the resource manager
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

// Sets both colors of the material
func (mat *Material) SetColors(diffuse, specular color.Color) {
	mat.DiffuseColor = diffuse
	mat.SpecularColor = specular
}

// Sets the shinyness of the material
func (mat *Material) SetShinyness(specularExponent float32) {
	MISE := MIN_SPECULAR_EXPONENT
	MASE := MAX_SPECULAR_EXPONENT
	y := math.Max(float64(specularExponent), 0.0)

	highNumber := -1.0 * math.Log((y-MISE)/MASE+1.0) / 3.0
	x := math.Pow(math.E, highNumber)

	mat.Shinyness = float32(math.Max(x, 0.0))
}

// Converts a color to a Vec3
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

// Converts a color to a Vec4
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
