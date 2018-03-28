package gohome

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/assimp"
	"image/color"
	"log"
	// "math"
)

const (
	MATERIAL_UNIFORM_NAME                         string = "material"
	MATERIAL_DIFFUSE_COLOR_UNIFORM_NAME           string = "diffuseColor"
	MATERIAL_SPECULAR_COLOR_UNIFORM_NAME          string = "specularColor"
	MATERIAL_SPECULAR_TEXTURE_UNIFORM_NAME        string = "specularTexture"
	MATERIAL_DIFFUSE_TEXTURE_UNIFORM_NAME         string = "diffuseTexture"
	MATERIAL_SHINYNESS_UNIFORM_NAME               string = "shinyness"
	MATERIAL_DIFFUSE_TEXTURE_LOADED_UNIFORM_NAME  string = "diffuseTextureLoaded"
	MATERIAL_SPECULAR_TEXTURE_LOADED_UNIFORM_NAME string = "specularTextureLoaded"
	MATERIAL_NORMALMAP_LOADED_UNIFORM_NAME        string = "normalMapLoaded"
	MATERIAL_NORMALMAP_UNIFORM_NAME               string = "normalMap"
)

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (this Color) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(float32(this.R) / 255.0 * float32(0xffff)), uint32(float32(this.G) / 255.0 * float32(0xffff)), uint32(float32(this.B) / 255.0 * float32(0xffff)), uint32(float32(this.A) / 255.0 * float32(0xffff))
}

type Material struct {
	DiffuseColor  color.Color
	SpecularColor color.Color

	DiffuseTexture  Texture
	SpecularTexture Texture
	NormalMap       Texture

	Shinyness float32

	diffuseTextureLoaded  uint8
	specularTextureLoaded uint8
	normalMapLoaded       uint8
}

func convertAssimpColor(color assimp.Color4) *Color {
	return &Color{uint8(color.R() * 255.0), uint8(color.G() * 255.0), uint8(color.B() * 255.0), uint8(color.A() * 255.0)}
}

func (mat *Material) InitDefault() {
	mat.DiffuseColor = &Color{255, 255, 255, 255}
	mat.SpecularColor = &Color{255, 255, 255, 255}

	mat.Shinyness = 0.5
}

func printProperties(material *assimp.Material) {
	for i := 0; i < material.NumProperties(); i++ {
		prop := material.Properties()[i]
		fmt.Println("Prop", i)
		fmt.Println("Index:", prop.Index())
		// fmt.Println("Key:", prop.Key()+"\x00")
		fmt.Println("Semantic:", prop.Semantic())
		fmt.Println("Type:", prop.Type())
		fmt.Println("DataLength: ", prop.DataLength())
		// fmt.Print("Data:")
		// for j := 0; j < prop.DataLength(); j++ {
		// 	fmt.Println("j:", j)
		// 	fmt.Print(prop.Data()[j])
		// }
		// fmt.Print("\n")
	}
}

func (mat *Material) Init(material *assimp.Material, scene *assimp.Scene, directory string, preloaded bool) {
	var ret assimp.Return
	var matDifColor assimp.Color4
	var matSpecColor assimp.Color4
	var matShininess float32

	// printProperties(material)

	matDifColor, ret = material.GetMaterialColor(assimp.MatKey_ColorDiffuse, 0, 0)
	if ret == assimp.Return_Failure {
		// log.Println("Couldn't return diffuse color:", assimp.GetErrorString())
		mat.DiffuseColor = &Color{255, 255, 255, 255}
	} else {
		mat.DiffuseColor = convertAssimpColor(matDifColor)
	}
	matSpecColor, ret = material.GetMaterialColor(assimp.MatKey_ColorSpecular, 0, 0)
	if ret == assimp.Return_Failure {
		// log.Println("Couldn't return specular color:", assimp.GetErrorString())
		mat.SpecularColor = &Color{255, 255, 255, 255}
	} else {
		mat.SpecularColor = convertAssimpColor(matSpecColor)
	}
	matShininess, ret = material.GetMaterialFloat(assimp.MatKey_Shininess, 0, 0)
	if ret == assimp.Return_Failure {
		// log.Println("Couldn't return shininess:", assimp.GetErrorString())
		mat.Shinyness = 0.0
	} else {
		mat.Shinyness = matShininess
	}

	diffuseTextures := material.GetMaterialTextureCount(1)
	specularTextures := material.GetMaterialTextureCount(2)
	for i := 0; i < diffuseTextures; i++ {
		texPath, _, _, _, _, _, _, ret := material.GetMaterialTexture(1, i)
		if ret == assimp.Return_Failure {
			log.Println("Couldn't return diffuse Texture")
		} else {
			if !preloaded {
				ResourceMgr.LoadTexture(texPath, directory+texPath)
				mat.DiffuseTexture = ResourceMgr.GetTexture(texPath)
			} else {
				mat.DiffuseTexture = ResourceMgr.loadTexture(texPath, directory+texPath, true)
			}

			break
		}
	}
	for i := 0; i < specularTextures; i++ {
		texPath, _, _, _, _, _, _, ret := material.GetMaterialTexture(2, i)
		if ret == assimp.Return_Failure {
			log.Println("Couldn't return specular Texture")
		} else {
			if !preloaded {
				ResourceMgr.LoadTexture(texPath, directory+texPath)
				mat.SpecularTexture = ResourceMgr.GetTexture(texPath)
			} else {
				mat.SpecularTexture = ResourceMgr.loadTexture(texPath, directory+texPath, true)
			}
			break
		}
	}
}

func (mat *Material) SetTextures(diffuse, specular, normalMap string) {
	mat.DiffuseTexture = ResourceMgr.GetTexture(diffuse)
	mat.SpecularTexture = ResourceMgr.GetTexture(specular)
	mat.NormalMap = ResourceMgr.GetTexture(normalMap)
}

func (mat *Material) SetColors(diffuse, specular color.Color) {
	mat.DiffuseColor = diffuse
	mat.SpecularColor = specular
}

func colorToVec3(c color.Color) mgl32.Vec3 {
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

func colorToVec4(c color.Color) mgl32.Vec4 {
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
