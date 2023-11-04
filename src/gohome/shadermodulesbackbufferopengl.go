package gohome

import (
	"strings"

	glslgen "github.com/PucklaJ/GLSLGenerator"
)

var (
	GlobalsBBMS = []glslgen.Variable{
		glslgen.Variable{"vec2", "highp", "vertices[6]"},
		glslgen.Variable{"vec2", "highp", "texCoords[6]"},
	}

	SetValuesModuleBBMS = glslgen.Module{
		Name: "setValues",
		Body: `vertices[0] = vec2(-1.0,-1.0);
		vertices[1] = vec2(1.0,-1.0);
		vertices[2] = vec2(1.0,1.0);
		vertices[3] = vec2(1.0,1.0);
		vertices[4] = vec2(-1.0,1.0);
		vertices[5] = vec2(-1.0,-1.0);

		texCoords[0] = vec2(0.0,0.0);
		texCoords[1] = vec2(1.0,0.0);
		texCoords[2] = vec2(1.0,1.0);
		texCoords[3] = vec2(1.0,1.0);
		texCoords[4] = vec2(0.0,1.0);
		texCoords[5] = vec2(0.0,0.0);`,
	}

	SetGLPositionBBMS = glslgen.Module{
		Name: "setGLPosition",
		Body: "gl_Position = vec4(vertices[gl_VertexID],0.0,1.0);",
	}

	SetGLPositionBBNOMS = glslgen.Module{
		Name: "setGLPosition",
		Body: "gl_Position = vec4(vertex,0.0,1.0);",
	}

	TextureMSModule = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"sampler2DMS", "highp", "texture0"},
		},
		Functions: []glslgen.Function{
			glslgen.Function{
				"vec4 fetchColor()",
				`vec4 color = vec4(0.0);
				ivec2 texCoords = ivec2(fragTexCoord * textureSize(texture0));

				for(int i = 0;i<8;i++)
				{
					color += texelFetch(texture0,texCoords,i);
				}
				color /= 8.0;

				return color;`,
			},
		},
		Name: "textureMSModule",
		Body: "globalColor = fetchColor();",
	}

	SetFragTexCoordModuleBBMS = glslgen.Module{
		Name: "setFragTexCoord",
		Body: "fragTexCoord = texCoords[gl_VertexID];",
	}
)

const (
	SHADER_FLAG_NO_MS uint32 = (1 << 0)
)

func GenerateShaderBackBuffer(flags uint32) (n, v, f string) {
	var vertex glslgen.VertexGenerator
	var fragment glslgen.FragmentGenerator

	if flags&SHADER_FLAG_NO_MS == 0 {
		if Render.GetName() == "WebGL" {
			vertex.SetVersion("WebGL")
			fragment.SetVersion("WebGL")
		} else if strings.Contains(Render.GetName(), "OpenGLES") {
			vertex.SetVersion("300 es")
			fragment.SetVersion("300 es")
		} else {
			vertex.SetVersion("150")
			fragment.SetVersion("150")
		}

		vertex.AddOutputs(Outputs2D)
		vertex.AddGlobals(GlobalsBBMS)
		vertex.AddModule(SetValuesModuleBBMS)
		vertex.AddModule(SetGLPositionBBMS)
		vertex.AddModule(DepthModuleVertex2D)
		vertex.AddModule(SetFragTexCoordModuleBBMS)

		fragment.AddMakros(MakrosFragment2D)
		fragment.AddInputs(Outputs2D)
		fragment.AddGlobals(GlobalsFragment2D)
		fragment.AddModule(InitModuleFragment2D)
		fragment.AddModule(TextureMSModule)
		fragment.AddModule(FinishColorModuleFragment2D)
	} else {
		if Render.GetName() == "WebGL" {
			vertex.SetVersion("WebGL")
			fragment.SetVersion("WebGL")
		} else if strings.Contains(Render.GetName(), "OpenGLES") {
			vertex.SetVersion("100")
			fragment.SetVersion("100")
		} else {
			vertex.SetVersion(ShaderVersion)
			fragment.SetVersion(ShaderVersion)
		}

		vertex.AddAttributes(Attributes2D)
		vertex.AddOutputs(Outputs2D)
		vertex.AddGlobals(GlobalsVertex2D)
		vertex.AddModule(InitTexCoordModule2D)
		vertex.AddModule(SetGLPositionBBNOMS)
		vertex.AddModule(DepthModuleVertex2D)
		vertex.AddModule(FinishTexCoordModule2D)

		fragment.AddMakros(MakrosFragment2D)
		fragment.AddInputs(Outputs2D)
		fragment.AddGlobals(GlobalsFragment2D)
		fragment.AddModule(TextureModuleFragment2D)
		fragment.AddModule(FinishColorModuleFragment2D)
	}

	n = "BackBufferShader"
	v, f = vertex.String(), fragment.String()
	return
}
