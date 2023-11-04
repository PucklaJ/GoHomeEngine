package gohome

import (
	"strings"
)

// 2D

var (
	Attributes2D = []glslgen.Variable{
		glslgen.Variable{"vec2", "highp", "vertex"},
		glslgen.Variable{"vec2", "highp", "texCoord"},
	}

	AttributesVertexOnly2D = []glslgen.Variable{
		glslgen.Variable{"vec2", "highp", "vertex"},
	}

	AttributesShape2D = []glslgen.Variable{
		glslgen.Variable{"vec2", "highp", "vertex"},
		glslgen.Variable{"vec4", "highp", "color"},
	}

	Outputs2D = []glslgen.Variable{
		glslgen.Variable{"vec2", "highp", "fragTexCoord"},
	}

	OutputsShape2D = []glslgen.Variable{
		glslgen.Variable{"vec4", "highp", "fragColor"},
	}

	GlobalsVertex2D = []glslgen.Variable{
		glslgen.Variable{"vec2", "highp", "globalTexCoord"},
	}

	GlobalsFragment2D = []glslgen.Variable{
		glslgen.Variable{"vec4", "highp", "globalColor"},
	}

	MakrosVertex2D = []glslgen.Makro{
		glslgen.Makro{"FLIP_NONE", "0"},
		glslgen.Makro{"FLIP_HORIZONTAL", "1"},
		glslgen.Makro{"FLIP_VERTICAL", "2"},
		glslgen.Makro{"FLIP_DIAGONALLY", "3"},
	}

	MakrosFragment2D = []glslgen.Makro{
		glslgen.Makro{"KEY_COLOR_PADDING", "0.1"},
		glslgen.Makro{"ALPHA_DISCARD_PADDING", "0.1"},
	}

	InitModuleVertex2D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"mat3", "highp", "transformMatrix2D"},
			glslgen.Variable{"mat4", "highp", "projectionMatrix2D"},
			glslgen.Variable{"mat3", "highp", "viewMatrix2D"},
		},
		Name: "init",
		Body: "gl_Position = projectionMatrix2D *vec4(vec2(viewMatrix2D*transformMatrix2D*vec3(vertex,1.0)),0.0,1.0);",
	}

	InitTexCoordModule2D = glslgen.Module{
		Name: "initTexCoord",
		Body: "globalTexCoord = texCoord;",
	}

	DepthModuleVertex2D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"float", "highp", "depth"},
		},
		Name: "initDepth",
		Body: "gl_Position.z = depth;",
	}

	FlipModuleVertex2D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"int", "", "flip"},
		},
		Functions: []glslgen.Function{
			glslgen.Function{
				"vec2 flipTexCoord(vec2 tc)",
				`vec2 flippedTexCoord;


				if(flip == FLIP_NONE)
				{
					flippedTexCoord = tc;
				}
				else if(flip == FLIP_HORIZONTAL)
				{
					flippedTexCoord.x = 1.0-tc.x;
					flippedTexCoord.y = tc.y;
				}
				else if(flip == FLIP_VERTICAL)
				{
					flippedTexCoord.x = tc.x;
					flippedTexCoord.y = 1.0-tc.y;
				}
				else if(flip == FLIP_DIAGONALLY)
				{
					flippedTexCoord.x = 1.0-tc.x;
					flippedTexCoord.y = 1.0-tc.y;
				}
				else
				{
					flippedTexCoord = tc;
				}

				return flippedTexCoord;`,
			},
		},
		Name: "flipModule",
		Body: "globalTexCoord = flipTexCoord(globalTexCoord);",
	}

	TextureRegionModule = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"vec4", "highp", "textureRegion"},
			glslgen.Variable{"bool", "", "enableTextureRegion"},
		},
		Functions: []glslgen.Function{
			glslgen.Function{
				"vec2 textureRegionToTexCoord(vec2 tc)",
				`if(!enableTextureRegion)
				return tc;

			// X: 0 ->      0 -> Min X
			// Y: 0 ->      0 -> Min Y
			// Z: WIDTH ->  1 -> Max X
			// W: HEIGHT -> 1 -> Max Y

			vec2 newTexCoord = tc;
			newTexCoord.x = newTexCoord.x * (textureRegion.z-textureRegion.x)+textureRegion.x;
			newTexCoord.y = newTexCoord.y * (textureRegion.w-textureRegion.y)+textureRegion.y;

			return newTexCoord;`,
			},
		},
		Name: "textureRegionModule",
		Body: "globalTexCoord = textureRegionToTexCoord(globalTexCoord);",
	}

	FinishTexCoordModule2D = glslgen.Module{
		Name: "finishTexCoord",
		Body: "fragTexCoord = globalTexCoord;",
	}

	FinishFragColorModule2D = glslgen.Module{
		Name: "finishFragColor",
		Body: "fragColor = color;",
	}

	InitModuleFragment2D = glslgen.Module{
		Name: "init",
		Body: "globalColor = vec4(1.0,1.0,1.0,1.0);",
	}

	TextureModuleFragment2D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"sampler2D", "highp", "texture0"},
		},
		Name: "textureModule",
		Body: "globalColor = texture2D(texture0,fragTexCoord);",
	}

	KeyColorModuleFragment2D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"vec3", "highp", "keyColor"},
			glslgen.Variable{"bool", "", "enableKey"},
		},
		Name: "keyColorModule",
		Body: `if(enableKey)
		{
			if(globalColor.r >= keyColor.r - KEY_COLOR_PADDING && globalColor.r <= keyColor.r + KEY_COLOR_PADDING &&
			   globalColor.g >= keyColor.g - KEY_COLOR_PADDING && globalColor.g <= keyColor.g + KEY_COLOR_PADDING &&
			   globalColor.b >= keyColor.b - KEY_COLOR_PADDING && globalColor.b <= keyColor.b + KEY_COLOR_PADDING)
			{
			   discard;
			}
		}`,
	}

	ModColorModuleFragment2D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"vec4", "highp", "modColor"},
			glslgen.Variable{"bool", "", "enableMod"},
		},
		Name: "modColorModule",
		Body: `if(enableMod)
		{
			globalColor *= modColor;
		}`,
	}

	DepthMapModuleFragment2D = glslgen.Module{
		Name: "depthMapModule",
		Body: "globalColor = vec4(globalColor.r,globalColor.r,globalColor.r,globalColor.a);",
	}

	FinishColorModuleFragment2D = glslgen.Module{
		Name: "finishColor",
		Body: `gl_FragColor = globalColor;
		if(gl_FragColor.a < ALPHA_DISCARD_PADDING)
			discard;`,
	}

	FragColorModuleFragment2D = glslgen.Module{
		Name: "applyFragColor",
		Body: `gl_FragColor = fragColor;
		if(gl_FragColor.a < ALPHA_DISCARD_PADDING)
			discard;`,
	}

	ColorModuleFragment2D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"vec4", "highp", "color"},
		},
		Name: "colorModule",
		Body: "globalColor *= color;",
	}
)

const (
	SHADER_TYPE_SPRITE2D uint8 = 1
	SHADER_TYPE_SHAPE2D  uint8 = 2
	SHADER_TYPE_TEXT2D   uint8 = 3

	SHADER_TYPE_3D      uint8 = 4
	SHADER_TYPE_SHAPE3D uint8 = 5
)

func GenerateShader2D(shader_type uint8, flags uint32) (n, v, f string) {
	startFlags := flags
	flags &= ^SHADER_FLAG_INSTANCED // Instanced is not implemented yet
	if flags&SHADER_FLAG_NO_TEXTURE != 0 {
		flags &= ^(SHADER_FLAG_NO_TEXTURE_REGION | SHADER_FLAG_NO_FLIP)
	}
	if shader_type == SHADER_TYPE_TEXT2D {
		flags &= ^SHADER_FLAG_NO_TEXTURE_REGION
	}

	var vertex glslgen.VertexGenerator
	var fragment glslgen.FragmentGenerator

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

	fragment.AddMakros(MakrosFragment2D)

	if shader_type == SHADER_TYPE_SPRITE2D || shader_type == SHADER_TYPE_TEXT2D {
		if flags&SHADER_FLAG_NO_FLIP == 0 {
			vertex.AddMakros(MakrosVertex2D)
		}
		if flags&SHADER_FLAG_NO_TEXTURE == 0 {
			vertex.AddAttributes(Attributes2D)
			vertex.AddOutputs(Outputs2D)
			vertex.AddGlobals(GlobalsVertex2D)
		} else {
			vertex.AddAttributes(AttributesVertexOnly2D)
		}

		vertex.AddModule(InitModuleVertex2D)
		if flags&SHADER_FLAG_NO_TEXTURE == 0 {
			vertex.AddModule(InitTexCoordModule2D)
		}
		if flags&SHADER_FLAG_NO_DEPTH == 0 {
			vertex.AddModule(DepthModuleVertex2D)
		}
		if flags&SHADER_FLAG_NO_FLIP == 0 {
			vertex.AddModule(FlipModuleVertex2D)
		}
		if flags&SHADER_FLAG_NO_TEXTURE_REGION == 0 {
			vertex.AddModule(TextureRegionModule)
		}
		if flags&SHADER_FLAG_NO_TEXTURE == 0 {
			vertex.AddModule(FinishTexCoordModule2D)
		}

		if flags&SHADER_FLAG_NO_TEXTURE == 0 {
			fragment.AddInputs(Outputs2D)
		}
		fragment.AddGlobals(GlobalsFragment2D)
		fragment.AddModule(InitModuleFragment2D)
		if flags&SHADER_FLAG_NO_TEXTURE == 0 {
			fragment.AddModule(TextureModuleFragment2D)
		}
	}

	if shader_type == SHADER_TYPE_SPRITE2D {
		if flags&SHADER_FLAG_NO_KEYCOLOR == 0 {
			fragment.AddModule(KeyColorModuleFragment2D)
		}
		if flags&SHADER_FLAG_NO_MODCOLOR == 0 {
			fragment.AddModule(ModColorModuleFragment2D)
		}
		if flags&SHADER_FLAG_DEPTHMAP != 0 {
			fragment.AddModule(DepthMapModuleFragment2D)
		}

	} else if shader_type == SHADER_TYPE_SHAPE2D {
		vertex.AddAttributes(AttributesShape2D)
		vertex.AddOutputs(OutputsShape2D)
		vertex.AddModule(InitModuleVertex2D)
		vertex.AddModule(DepthModuleVertex2D)
		vertex.AddModule(FinishFragColorModule2D)

		fragment.AddInputs(OutputsShape2D)
		fragment.AddModule(FragColorModuleFragment2D)
	} else { // SHADER_TYPE_TEXT2D
		fragment.AddModule(ColorModuleFragment2D)
	}

	if shader_type == SHADER_TYPE_SPRITE2D || shader_type == SHADER_TYPE_TEXT2D {
		fragment.AddModule(FinishColorModuleFragment2D)
	}

	v, f = vertex.String(), fragment.String()
	n = GetShaderName2D(shader_type, startFlags)

	return
}

func GetShaderName2D(shader_type uint8, flags uint32) string {
	var n string
	switch shader_type {
	case SHADER_TYPE_SPRITE2D:
		n = SPRITE2D_SHADER_NAME
	case SHADER_TYPE_SHAPE2D:
		n = SHAPE_2D_SHADER_NAME
	default:
		n = TEXT_2D_SHADER_NAME
	}
	for i := uint32(0); i < NUM_FLAGS_2D; i++ {
		if flags&(1<<i) != 0 {
			n += " " + FLAG_NAMES_2D[i]
		}
	}
	return n
}

func LoadGeneratedShader2D(shader_type uint8, flags uint32) Shader {
	n, v, f := GenerateShader2D(shader_type, flags)
	return ls(n, v, f)
}
