package gohome

import (
	"github.com/PucklaMotzer09/GLSLGenerator"
	"strings"
)

var (
	AttributesShape3D = []glslgen.Variable{
		glslgen.Variable{"vec3", "highp", "vertex"},
		glslgen.Variable{"vec4", "highp", "color"},
	}

	UniformModuleShape3D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"mat4", "highp", "transformMatrix3D"},
			glslgen.Variable{"mat4", "highp", "viewMatrix3D"},
			glslgen.Variable{"mat4", "highp", "projectionMatrix3D"},
		},
		Name: "uniformModuleShape3D",
	}
)

func generateShaderShape3D() (n, v, f string) {
	var vertex glslgen.VertexGenerator
	var fragment glslgen.FragmentGenerator

	if strings.Contains(Render.GetName(), "OpenGLES") {
		vertex.SetVersion("100")
		fragment.SetVersion("100")
	} else {
		vertex.SetVersion("110")
		fragment.SetVersion("110")
	}

	vertex.AddAttributes(AttributesShape3D)
	vertex.AddOutputs(OutputsShape2D)
	vertex.AddModule(UniformModuleShape3D)
	vertex.AddModule(CalculatePositionModule3D)
	vertex.AddModule(FinishFragColorModule2D)

	fragment.AddMakros(MakrosFragment2D)
	fragment.AddInputs(OutputsShape2D)
	fragment.AddModule(FragColorModuleFragment2D)

	n = SHAPE3D_SHADER_NAME
	v, f = vertex.String(), fragment.String()

	return
}
