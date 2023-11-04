package gohome

import (
	"strings"
)

var (
	AttributesShape3D = []glslgen.Variable{
		glslgen.Variable{"vec3", "highp", "vertex"},
		glslgen.Variable{"vec4", "highp", "color"},
	}

	UniformModuleMatricesDefault3D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"mat4", "highp", "viewMatrix3D"},
			glslgen.Variable{"mat4", "highp", "projectionMatrix3D"},
		},
		Name: "uniformModuleMatricesDefault",
	}

	UniformModuleMatricesNormal3D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"mat4", "highp", "transformMatrix3D"},
		},
		Name: "uniformModuleMatricesNormal",
	}
)

func generateShaderShape3D() (n, v, f string) {
	var vertex glslgen.VertexGenerator
	var fragment glslgen.FragmentGenerator

	if Render.GetName() == "WebGL" {
		vertex.SetVersion("WebGL")
		fragment.SetVersion("WebGL")
	} else if strings.Contains(Render.GetName(), "OpenGLES") {
		vertex.SetVersion("100")
		fragment.SetVersion("100")
	} else {
		vertex.SetVersion("110")
		fragment.SetVersion("110")
	}

	vertex.AddAttributes(AttributesShape3D)
	vertex.AddOutputs(OutputsShape2D)
	vertex.AddModule(UniformModuleMatricesDefault3D)
	vertex.AddModule(UniformModuleMatricesNormal3D)
	vertex.AddModule(CalculatePositionModule3D)
	vertex.AddModule(FinishFragColorModule2D)

	fragment.AddMakros(MakrosFragment2D)
	fragment.AddInputs(OutputsShape2D)
	fragment.AddModule(FragColorModuleFragment2D)

	n = SHAPE3D_SHADER_NAME
	v, f = vertex.String(), fragment.String()

	return
}
