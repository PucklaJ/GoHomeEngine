package gohome

import (
	"github.com/PucklaMotzer09/GLSLGenerator"
	"strings"
)

var (
	SetFragTexCoordModuleVertex3D = glslgen.Module{
		Name: "setFragTexCoord",
		Body: "fragTexCoord = texCoord;",
	}

	SimpleMaterialModule3D = glslgen.Module{
		Structs: []glslgen.Struct{
			glslgen.Struct{
				Name: "Material",
				Variables: []glslgen.Variable{
					glslgen.Variable{"bool", "", "DiffuseTextureLoaded"},
					glslgen.Variable{"float", "highp", "transparency"},
				},
			},
		},
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"Material", "", "material"},
		},
		Name: "simpleMaterialModule",
		Body: `globalColor.a *= material.transparency;`,
	}

	ShadowMapDiffuseTextureModule = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"sampler2D", "highp", "materialdiffuseTexture"},
		},
		Name: "diffuseTextureCalculation",
		Body: `if(material.DiffuseTextureLoaded)
				  globalColor = texture2D(materialdiffuseTexture,fragTexCoord);
			   else
				  globalColor = vec4(1.0);`,
	}

	ShadowMapFinishModuleFragment = glslgen.Module{
		Name: "finish",
		Body: `if(globalColor.a < ALPHA_DISCARD_PADDING)
				discard;`,
	}
)

func GenerateShaderShadowMap(flags uint32) (n, v, f string) {
	startFlags := flags
	if !Render.HasFunctionAvailable("INSTANCED") {
		flags &= ^SHADER_FLAG_INSTANCED
	}

	var vertex glslgen.VertexGenerator
	var fragment glslgen.FragmentGenerator

	if strings.Contains(Render.GetName(), "OpenGLES") {
		vertex.SetVersion("100")
		fragment.SetVersion("100")
	} else {
		vertex.SetVersion("110")
		fragment.SetVersion("110")
	}

	vertex.AddAttributes(Attributes3D)
	if flags&SHADER_FLAG_INSTANCED != 0 {
		vertex.AddAttributes(AttributesInstanced3D)
	}
	vertex.AddOutputs(Outputs2D)
	vertex.AddGlobals(GlobalsVertex2D)
	vertex.AddModule(InitTexCoordModule2D)
	vertex.AddModule(UniformModuleMatricesDefault3D)
	if flags&SHADER_FLAG_INSTANCED == 0 {
		vertex.AddModule(UniformModuleMatricesNormal3D)
	}
	vertex.AddModule(CalculatePositionModule3D)
	vertex.AddModule(FinishTexCoordModule2D)

	fragment.AddMakros(MakrosFragment2D)
	fragment.AddGlobals(GlobalsFragment2D)
	fragment.AddInputs(Outputs2D)
	fragment.AddModule(ShadowMapDiffuseTextureModule)
	fragment.AddModule(SimpleMaterialModule3D)
	fragment.AddModule(ShadowMapFinishModuleFragment)

	n = SHADOWMAP_SHADER_NAME
	if startFlags&SHADER_FLAG_INSTANCED != 0 {
		n += " Instanced"
	}
	v, f = vertex.String(), fragment.String()

	return
}

func LoadGeneratedShaderShadowMap(flags uint32) Shader {
	n, v, f := GenerateShaderShadowMap(flags)
	return ls(n, v, f)
}
